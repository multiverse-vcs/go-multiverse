# Multiverse

[![Go Report Card](https://goreportcard.com/badge/github.com/multiverse-vcs/go-multiverse)](https://goreportcard.com/report/github.com/multiverse-vcs/go-multiverse)

Multiverse is a version control system that enables peer-to-peer software development.
It differs from other version control systems in that it is completely decentralized.
Repositories are shared using a globally distributed commit graph instead of centralized servers.

### Install

Go 1.14 or higher is required.

```bash
$ git clone https://github.com/multiverse-vcs/go-multiverse.git
$ cd go-multiverse
$ make install
```

### Concepts

Multiverse uses content identifiers or [CIDs](https://docs.ipfs.io/concepts/content-addressing/#identifier-formats) to locate and retrieve content.
A CID can refer to a specific commit, tag, or branch of a repository. Only trust shared CIDs from secure sources.
Below is an example of using a CID to clone a repository.

```bash
$ multi clone bafyreid6o5dhbnxthplfnsuuyqyhxs6vi22xep62f3szk2krq2r7pyaiim project
```

When you are ready to share your changes you will need to connect to other peers in the network.
Multiverse uses multiple networking protocols to find peers on both global and local networks.
An internet connection is **never** required to share changes with peers on your local network.
Use the swap command to connect with peers and share changes.

```bash
$ multi swap
```

### Usage

```
USAGE:
   multi [global options] command [command options] [arguments...]

COMMANDS:
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
