# tuckit

Clone a git repository into your base directory — and optionally register it with
[zoxide](https://github.com/ajeetdsouza/zoxide) so you can jump to it instantly.

`tuckit` keeps all your clones tucked under one directory. Give it a URL (or let it
prompt you), and it clones into `<base_dir>/<repo-name>`.

## Install

```sh
go install github.com/alain-sarti/tuckit@latest
```

Requires [git](https://git-scm.com/) on your `PATH`. [zoxide](https://github.com/ajeetdsouza/zoxide)
is optional — needed only for the "Add to zoxide?" step.

## Configure

`tuckit` needs to know your base directory. Set it one of two ways (the environment
variable wins if both are present):

**Environment variable:**

```sh
export TUCKIT_BASE_DIR=~/Dev
```

**Config file** at `~/.config/tuckit/config.toml`:

```toml
base_dir = "~/Dev"
```

A leading `~` is expanded to your home directory in both cases.

## Usage

Clone by URL:

```sh
tuckit github.com/alain-sarti/tuckit
```

Short URLs get `https://` prepended automatically; SCP-like SSH (`git@…`) and
URLs that already carry a scheme (`https://`, `http://`, `ssh://`, `git://`) are
used as-is.

Run without a URL to be prompted for one:

```sh
tuckit
```

On success `tuckit` prints the clone path, then asks whether to add it to zoxide.
If the target path is already occupied — a file, or a non-empty directory —
`tuckit` stops with an error rather than touching it.

## Acknowledgements

Parts of `tuckit` were written with the help of [Claude Code](https://www.anthropic.com/claude-code).

## License

[MIT](LICENSE)
