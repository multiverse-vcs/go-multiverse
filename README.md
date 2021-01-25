# Multiverse

[![Go Report Card](https://goreportcard.com/badge/github.com/multiverse-vcs/go-multiverse)](https://goreportcard.com/report/github.com/multiverse-vcs/go-multiverse)

A decentralized version control system for peer-to-peer software development.

> This project is currently alpha quality. We are looking for feedback from the open source community.

### Features

- ***Peer-to-peer*** - self host your repositories
- ***Works offline*** - integrated local code viewer
- ***Secure*** - all communications are encrypted

### Install

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
   init     Create a repo
   merge    Merge commits
   status   Print changes
   tag      List, create, or delete tags
   help, h  Shows a list of commands or help for one command
```

### Getting Started

[Read the manual](https://www.multiverse-vcs.com/docs/).

### Contributing

Found a bug or have a feature request? [Open an issue](https://github.com/multiverse-vcs/go-multiverse/issues/new).

Multiverse follows the [Contributor Covenant](https://contributor-covenant.org/version/2/0/code_of_conduct/) Code of Conduct.

### License

GNU Affero General Public License v3.0
