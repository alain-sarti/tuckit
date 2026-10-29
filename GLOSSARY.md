# tuckit

tuckit clones git repositories into one flat directory. Each clone gets a folder named after its repository.

## Language

**Base directory**:
The configured directory that every target path sits directly under. It is a setting, not a precondition, so it does not have to exist before the first clone.
_Avoid_: base_dir (outside config), root, workspace, clone dir

**Repository URL**:
The address of the repository to clone, in whatever form the user gives it.
_Avoid_: URL alone, remote, source

**Short URL**:
A repository URL that has no scheme and isn't in SCP-like SSH form (`git@host:path`). tuckit treats every short URL as HTTPS, including ones that look like `host:port/path`.
_Avoid_: short form, bare URL

**Repository name**:
The last path segment of a repository URL, without a trailing `.git`. It is the only thing that identifies a clone in the base directory. Host and owner play no part.
_Avoid_: repo name, project name

**Target path**:
The base directory joined with the repository name. A clone is written here and nowhere else.
_Avoid_: clone path, target, destination

**Occupied**:
A target path is occupied when something already there would make the clone fail, meaning a file or a non-empty directory. A missing path or an empty directory is not occupied.
_Avoid_: conflict, unavailable, taken

**Zoxide registration**:
An optional convenience step after a successful clone that adds the target path to zoxide. It is never part of the clone, so declining it, skipping it or having it fail leaves the clone a success.
_Avoid_: zoxide add, zoxide step
