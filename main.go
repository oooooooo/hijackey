// Command hijackey runs another program under a pseudo-terminal and remaps
// single-byte keystrokes before they reach it, e.g.:
//
//	hijackey space=d b=u leaf
//
// presses 'd' as if space were pressed, and 'b' as if 'u' were pressed.
package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

// version is overridden at release-build time via -ldflags
// "-X main.version=vX.Y.Z"; local/dev builds report "dev".
var version = "dev"

const usage = "usage: hijackey [srcKey=dstKey ...] <command> [args...]"

const helpText = usage + `

Hijack single-character keystrokes before they reach <command>, running
under a pty.

A key is a single ASCII character (d, U, ...) or a named key: space, tab,
enter, backspace, esc, eq (=).

Flags:
  --version, -v   print version and exit
  --help, -h      print this help and exit
`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "hijackey:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cmdArgs, mapping, err := parseArgs(args)
	if err != nil {
		return err
	}
	if len(cmdArgs) > 0 && (cmdArgs[0] == "--version" || cmdArgs[0] == "-v") {
		fmt.Println("hijackey", version)
		return nil
	}
	if len(cmdArgs) > 0 && (cmdArgs[0] == "--help" || cmdArgs[0] == "-h") {
		fmt.Print(helpText)
		return nil
	}
	if len(cmdArgs) == 0 {
		return fmt.Errorf("%s", usage)
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("start %s: %w", cmdArgs[0], err)
	}
	defer ptmx.Close()

	if err := syncSize(ptmx); err != nil {
		return fmt.Errorf("sync terminal size: %w", err)
	}
	winch := make(chan os.Signal, 1)
	signal.Notify(winch, syscall.SIGWINCH)
	go func() {
		for range winch {
			_ = syncSize(ptmx)
		}
	}()
	defer signal.Stop(winch)

	stdinFd := int(os.Stdin.Fd())
	restore, err := term.MakeRaw(stdinFd)
	if err != nil {
		return fmt.Errorf("enter raw mode: %w", err)
	}
	defer func() { _ = term.Restore(stdinFd, restore) }()

	go io.Copy(os.Stdout, ptmx)
	go copyRemapped(ptmx, os.Stdin, mapping)

	err = cmd.Wait()
	if exitErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
	}
	return err
}

func syncSize(ptmx *os.File) error {
	size, err := pty.GetsizeFull(os.Stdout)
	if err != nil {
		return err
	}
	return pty.Setsize(ptmx, size)
}

// copyRemapped copies bytes from src to dst, substituting any byte found in
// mapping. Multi-byte sequences (arrow keys, mouse reports, ...) pass through
// untouched unless one of their individual bytes happens to be a mapped key.
func copyRemapped(dst io.Writer, src io.Reader, mapping map[byte]byte) {
	buf := make([]byte, 4096)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			out := make([]byte, n)
			for i, b := range buf[:n] {
				if mapped, ok := mapping[b]; ok {
					out[i] = mapped
				} else {
					out[i] = b
				}
			}
			if _, werr := dst.Write(out); werr != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

// parseArgs splits args into the target command (with its own arguments) and
// the srcKey=dstKey mapping directives. A token is treated as a mapping
// directive when it has the shape `word=word` (no leading '-'), wherever it
// appears; anything else is passed through as an argument to the target
// command, so the command name is simply the first non-mapping token.
func parseArgs(args []string) (cmdArgs []string, mapping map[byte]byte, err error) {
	mapping = make(map[byte]byte)
	for _, arg := range args {
		src, dst, ok := splitMapping(arg)
		if !ok {
			cmdArgs = append(cmdArgs, arg)
			continue
		}
		srcByte, err := keyByte(src)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", arg, err)
		}
		dstByte, err := keyByte(dst)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", arg, err)
		}
		if _, dup := mapping[srcByte]; dup {
			return nil, nil, fmt.Errorf("%s: duplicate mapping for key %q", arg, src)
		}
		mapping[srcByte] = dstByte
	}
	return cmdArgs, mapping, nil
}

func splitMapping(arg string) (src, dst string, ok bool) {
	if strings.HasPrefix(arg, "-") {
		return "", "", false
	}
	src, dst, found := strings.Cut(arg, "=")
	if !found || src == "" || dst == "" {
		return "", "", false
	}
	return src, dst, true
}

var namedKeys = map[string]byte{
	"space":     0x20,
	"tab":       0x09,
	"enter":     0x0d,
	"backspace": 0x7f,
	"esc":       0x1b,
	"eq":        '=',
}

// keyByte resolves a key spec to the single raw byte it produces in a
// terminal's raw input stream: either a named key (space, tab, enter,
// backspace, esc, eq) or a single ASCII character typed literally.
func keyByte(s string) (byte, error) {
	if b, ok := namedKeys[strings.ToLower(s)]; ok {
		return b, nil
	}
	if len(s) == 1 && s[0] < 0x80 {
		return s[0], nil
	}
	return 0, fmt.Errorf("unsupported key %q (use a single ASCII character or one of: space, tab, enter, backspace, esc, eq)", s)
}
