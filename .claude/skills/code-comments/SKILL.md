---
name: code-comments
description: >-
  How to write source code comments in this repo. Use whenever adding or
  editing a comment in .go (or any source) file. Comments must explain
  *why*, never *what*; if the why isn't actually known, ask the user
  instead of inventing one.
---

# code-comments

## Default: no comment

Well-named identifiers already say *what* the code does. Don't add a
comment unless the *why* is genuinely non-obvious: a hidden constraint, a
subtle invariant, a workaround for a specific bug, or behavior that would
surprise a future reader.

## When you do write one, it must explain why

Bad — restates what the code already says:

```go
// increment i by one
i++
```

Good — explains a non-obvious why:

```go
// SIGWINCH must be forwarded before the first Setsize call, or the child
// inherits the pty's stale size from before the resize.
```

## If you don't know the why: ask, don't guess

Never invent a plausible-sounding rationale to fill a comment. If you're
about to write a why-comment (or asked to add one) and you're not actually
sure why the constraint/workaround/behavior exists, stop and ask the user
instead of fabricating a justification. A wrong "why" is worse than no
comment — it actively misleads the next reader, including future sessions.

## Notes

- Applies to comments you write or edit going forward; don't rewrite
  existing comments unprompted just to re-justify them.
- Pairs with [[conventional-commits]] — commit body explains why a change
  was made, this explains why the resulting code is shaped the way it is.
