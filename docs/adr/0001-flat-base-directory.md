# Flat base directory keyed by repository name

tuckit puts every clone directly under the base directory, in a folder named after the repository name. Host and owner are dropped. This is the layout the maintainer already uses, and it keeps target paths short.

## Considered Options

- **Flat (chosen).** `~/Dev/tuckit`.
- **Host and owner tree, as `ghq` does it.** `~/Dev/github.com/alain-sarti/tuckit`. Two repositories with the same name never collide, but every path is three levels deeper.

## Consequences

Two repositories with the same name collide. After `github.com/alain-sarti/tuckit` is cloned, `gitlab.com/someone-else/tuckit` finds its target path occupied and tuckit refuses. That refusal is intended, not a bug.

Switching layouts later would mean moving every existing clone and fixing the zoxide entries that point at the old paths.
