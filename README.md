# Multiverse

[![Gitter](https://badges.gitter.im/multiverse-vcs/community.svg)](https://gitter.im/multiverse-vcs/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/multiverse-vcs/go-multiverse)](https://goreportcard.com/report/github.com/multiverse-vcs/go-multiverse)
[![codecov](https://codecov.io/gh/multiverse-vcs/go-multiverse/branch/master/graph/badge.svg?token=Y6UBYBD56P)](https://codecov.io/gh/multiverse-vcs/go-multiverse)

> A decentralized version control system for peer-to-peer software development.

*This project is under heavy development. Changes may cause compatibility issues with previous versions.*

### Features

- ***Peer-to-peer*** - self host your repositories
- ***Works offline*** - integrated local code viewer
- ***Private & Secure*** - all communications are encrypted

### Getting Started

[Read the manual](https://www.multiverse-vcs.com/docs/).

### Building

Go 1.16 or higher is required.

```bash
$ git clone https://github.com/multiverse-vcs/go-multiverse.git
$ cd go-multiverse
$ make install
```

### Usage

```
USAGE:
   multi [global options] command [command options] [arguments...]
   
COMMANDS:
   branch   List, create, or delete branches
   commit   Record changes
   daemon   Starts a client
   import   Import a repo
   init     Create a repo
   log      Print repo history
   merge    Merge commits
   status   Print changes
   tag      List, create, or delete tags
   help, h  Shows a list of commands or help for one command
```

### Contributing

Found a bug or have a feature request? [Open an issue](https://github.com/multiverse-vcs/go-multiverse/issues/new).

Multiverse follows the [Contributor Covenant](https://contributor-covenant.org/version/2/0/code_of_conduct/) Code of Conduct.

### License

GNU Affero General Public License v3.0
