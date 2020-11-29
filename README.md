# Multiverse

[![Go Report Card](https://goreportcard.com/badge/github.com/multiverse-vcs/go-multiverse)](https://goreportcard.com/report/github.com/multiverse-vcs/go-multiverse)

Multiverse is a decentralized version control system that enables peer-to-peer software development.

### About

Multiverse differs from other version control systems in that it is completely *decentralized*.

Repositories are shared using a globally distributed commit graph instead of *centralized* servers.

### Install

Go 1.14 or higher is required.

```bash
$ git clone https://github.com/multiverse-vcs/go-multiverse.git
$ cd go-multiverse
$ make install
```

### Usage

```
branch    Add, remove, or list branches
checkout  Checkout committed files
clone     Copy an existing repo
commit    Record repo changes
fetch     Fetch changes from peers
init      Initialize a new repo
log       Print change history
status    Print repo status
swap      Exchange data with peers
switch    Change branches
help, h   Shows a list of commands or help for one command
```

### Contributing

Found a bug or have a feature request? [Open an issue](https://github.com/multiverse-vcs/go-multiverse/issues/new).

Multiverse follows the [Contributor Covenant](https://contributor-covenant.org/version/2/0/code_of_conduct/) Code of Conduct.

### License

GNU Affero General Public License v3.0