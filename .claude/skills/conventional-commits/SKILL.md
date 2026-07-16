---
name: conventional-commits
description: >-
  Write git commit messages for this repo. Use whenever creating a commit.
  Enforces Conventional Commits format for the subject line and requires the
  body to explain *why*, not just restate the diff.
---

# conventional-commits

## When

Every time a commit is created in this repo (including automated release
commits made by CI, e.g. `.github/workflows/ci.yml`'s `chore: release ...`
step).

## Subject line

`type(scope): short summary`, imperative mood, no trailing period, under
~70 chars.

Types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `ci`, `build`,
`perf`. Scope is optional — omit it if the whole repo/tool is affected
(this is a single-file CLI, so scope is often unnecessary).

## Body — always explain why

The subject says *what* changed; the body must say *why* — the motivating
problem, constraint, prior behavior, or decision. A body that only restates
the diff in prose is not acceptable; if you can't say why, the subject line
alone is enough and the body should be omitted rather than padded.

Good:

```text
fix(pty): forward SIGWINCH to child on resize

Without this the child kept using the size it started with, so
`less`/`vim` never learned the terminal had been resized.
```

Bad — restates the diff, no reasoning:

```text
fix(pty): forward SIGWINCH to child on resize

Added a signal.Notify call for SIGWINCH and forward it to the child.
```

## Notes

- Trivial changes (formatting, typo fixes) can skip the body, but the
  subject line itself must still be accurate — don't stretch "why" prose
  onto something that has none.
- This overrides the built-in `conventional-commits` skill for this repo
  only; the "always include why" body requirement is specific to this
  project and not assumed elsewhere.
