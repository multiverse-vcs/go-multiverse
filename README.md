# Multiverse

[![Go Report Card](https://goreportcard.com/badge/github.com/multiverse-vcs/go-multiverse)](https://goreportcard.com/report/github.com/multiverse-vcs/go-multiverse)

Multiverse is a decentralized version control system that enables peer-to-peer software development.

It builds on top of [IPFS](https://ipfs.io) to provide a secure, resilient, and censorship resistent developer network.

### Building

To build Multiverse from source run the following.

```bash
$ make
```

Then copy the executable to your local bin directory.

```bash
$ sudo mv multi /usr/local/bin/
```

### Usage

```
Usage:
  multi [command]

Available Commands:
  branch      List, create, or delete branches.
  checkout    Copy changes from a commit to the local repo.
  commit      Record changes to the local repo.
  daemon      Run a persistent Multiverse node.
  help        Help about any command
  import      Import a repo from an external VCS.
  init        Create a new empty repo or copy an existing repo.
  log         Print change history.
  merge       Merge changes from a peer into the local repo.
  status      Print status of the local repo.
  switch      Change to a different branch.

Flags:
  -h, --help   help for multi

Use "multi [command] --help" for more information about a command.
```

### Contributing

Found a bug or have a feature request? [Open an issue](https://github.com/multiverse-vcs/go-multiverse/issues/new).

Multiverse follows the [Contributor Covenant](https://contributor-covenant.org/version/2/0/code_of_conduct/) Code of Conduct.

### License

GNU Affero General Public License v3.0