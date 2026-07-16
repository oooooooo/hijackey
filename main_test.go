package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	cmdArgs, mapping, err := parseArgs([]string{"leaf", "d=space", "b=u"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(cmdArgs, []string{"leaf"}) {
		t.Errorf("cmdArgs = %v, want [leaf]", cmdArgs)
	}
	want := map[byte]byte{'d': 0x20, 'b': 'u'}
	if !reflect.DeepEqual(mapping, want) {
		t.Errorf("mapping = %v, want %v", mapping, want)
	}
}

func TestParseArgsMappingBeforeCommand(t *testing.T) {
	cmdArgs, mapping, err := parseArgs([]string{"space=d", "b=u", "leaf", "main.md"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(cmdArgs, []string{"leaf", "main.md"}) {
		t.Errorf("cmdArgs = %v, want [leaf main.md]", cmdArgs)
	}
	want := map[byte]byte{0x20: 'd', 'b': 'u'}
	if !reflect.DeepEqual(mapping, want) {
		t.Errorf("mapping = %v, want %v", mapping, want)
	}
}

func TestParseArgsPassesThroughCommandArgs(t *testing.T) {
	cmdArgs, mapping, err := parseArgs([]string{"leaf", "notes.md", "--watch", "d=space"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(cmdArgs, []string{"leaf", "notes.md", "--watch"}) {
		t.Errorf("cmdArgs = %v, want [leaf notes.md --watch]", cmdArgs)
	}
	if mapping['d'] != 0x20 {
		t.Errorf("mapping[d] = %v, want space", mapping['d'])
	}
}

func TestParseArgsDuplicateMapping(t *testing.T) {
	_, _, err := parseArgs([]string{"leaf", "d=space", "d=u"})
	if err == nil {
		t.Fatal("expected error for duplicate mapping, got nil")
	}
}

func TestParseArgsInvalidKey(t *testing.T) {
	_, _, err := parseArgs([]string{"leaf", "foo=bar"})
	if err == nil {
		t.Fatal("expected error for multi-character key, got nil")
	}
}

func TestParseArgsEq(t *testing.T) {
	cmdArgs, mapping, err := parseArgs([]string{"eq=d", "leaf"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cmdArgs, []string{"leaf"}) {
		t.Errorf("cmdArgs = %v, want [leaf]", cmdArgs)
	}
	if mapping['='] != 'd' {
		t.Errorf("mapping[=] = %v, want d", mapping['='])
	}
}

func TestKeyByte(t *testing.T) {
	cases := map[string]byte{
		"space":     0x20,
		"Space":     0x20,
		"tab":       0x09,
		"enter":     0x0d,
		"backspace": 0x7f,
		"esc":       0x1b,
		"eq":        '=',
		"d":         'd',
		"U":         'U',
	}
	for in, want := range cases {
		got, err := keyByte(in)
		if err != nil {
			t.Errorf("keyByte(%q) error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("keyByte(%q) = %#x, want %#x", in, got, want)
		}
	}
}

func TestKeyByteInvalid(t *testing.T) {
	for _, in := range []string{"", "ab", "space2", "return", "escape"} {
		if _, err := keyByte(in); err == nil {
			t.Errorf("keyByte(%q) expected error, got nil", in)
		}
	}
}

func TestRunVersionFlag(t *testing.T) {
	for _, flag := range []string{"--version", "-v"} {
		if err := run([]string{flag}); err != nil {
			t.Errorf("run([%q]) error: %v", flag, err)
		}
	}
}

func TestRunHelpFlag(t *testing.T) {
	for _, flag := range []string{"--help", "-h"} {
		if err := run([]string{flag}); err != nil {
			t.Errorf("run([%q]) error: %v", flag, err)
		}
	}
}

func TestCopyRemapped(t *testing.T) {
	mapping := map[byte]byte{'d': 0x20, 'b': 'u'}
	src := bytes.NewReader([]byte("db hello\n"))
	var dst bytes.Buffer
	copyRemapped(&dst, src, mapping)
	want := " u hello\n"
	if dst.String() != want {
		t.Errorf("copyRemapped output = %q, want %q", dst.String(), want)
	}
}
