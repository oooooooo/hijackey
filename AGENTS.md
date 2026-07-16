# hijackey

Single-file Go CLI. Runs another program under a pty and remaps single-byte
keystrokes before they reach it (e.g. `hijackey space=d b=u leaf`). See
README.md for user-facing usage and behavior.

## Layout

- `main.go` — everything: arg parsing (`parseArgs`, `splitMapping`,
  `keyByte`), pty setup (`run`), and the byte-remapping copy loop
  (`copyRemapped`). No other packages.
- `main_test.go` — unit tests for the pieces above that don't need a real
  pty/tty (arg parsing, key resolution, byte remapping).
- `main.md` — stray duplicate of `main.go`'s content, unrelated to the
  project; not doc source, don't treat it as such.
- `Taskfile.yml` — `build` and `test` tasks, run via the `task` CLI.

## Build / test

```bash
task build
task test
```

## Key invariant

`parseArgs` treats any token shaped like `word=word` (no leading `-`) as a
key mapping, wherever it appears in argv; the command name is the first
token that *isn't* shaped that way. This is what lets mappings precede the
command (`hijackey space=d b=u leaf`), which in turn is what makes the
alias-friendly usage in README.md ("Shell alias" section) work — a plain
`alias leaf='hijackey space=d b=u leaf'` relies on the shell appending trailing
args, so the command name must be allowed to show up after the mappings.
Keep this property if you touch arg parsing.

## Platform

Unix only (macOS/Linux) — depends on `github.com/creack/pty`, which has no
Windows implementation.
