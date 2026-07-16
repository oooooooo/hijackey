# hijackey

Hijack keystrokes for CLIs that won't let you rebind them.

```bash
hijackey space=d b=u leaf README.md
hijackey q=esc some-cli
```

While [leaf](https://github.com/rivolink/leaf) runs, `space` is `d` and `b` is `u`.

## Install

Unix only (macOS/Linux), via [creack/pty](https://github.com/creack/pty).

```bash
brew tap oooooooo/hijackey https://github.com/oooooooo/hijackey.git
brew install oooooooo/hijackey/hijackey
```

Or

```bash
git clone https://github.com/oooooooo/hijackey.git
cd hijackey
# https://taskfile.dev/docs/installation
task build
```

## Usage

```text
hijackey [srcKey=dstKey ...] <command> [args...]
```

A key is a single ASCII character (`d`, `U`, ...) or a named key: `space`, `tab`, `enter`, `backspace`, `esc`, `eq` (`=`).

`hijackey --version` (or `-v`) prints the version and exits; `--help` (or `-h`) prints this usage and exits.

## Alias

```bash
alias leaf='hijackey space=d b=u leaf'
leaf --width 80 README.md
```

## Limitations

No multi-byte sequences (arrow keys, function keys, mouse reports) — only single raw bytes. Remapping `esc` can break arrow keys, which also start with ESC.

## License

[MIT](LICENSE)
