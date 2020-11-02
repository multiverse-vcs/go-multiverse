# Multiverse

[![Go Report Card](https://goreportcard.com/badge/github.com/yondero/go-multiverse)](https://goreportcard.com/report/github.com/yondero/go-multiverse)

Multiverse is a decentralized version control system that enables peer-to-peer software development.

It builds on top of [IPFS](https://ipfs.io) to provide a secure, resilient, and censorship resistent developer network.

### Usage

```

Usage:
  multi [command]

Available Commands:
  checkout    Checkout files from a different commit.
  commit      Record changes in the local repo.
  daemon      Run a persistent Multiverse node.
  help        Help about any command
  ignore      Ignore changes to files matching pattern.
  init        Create a new empty repo or copy an existing repo.
  log         Print change history.
  merge       Merge changes from a peer into the local repo.
  publish     Announce a new version to peers.
  status      Print status of the local repo.

Flags:
  -h, --help   help for multi

Use "multi [command] --help" for more information about a command.
```

### Contributing

Found a bug or have a feature request? [Open an issue](https://github.com/yondero/go-multiverse/issues/new).

Multiverse follows the [Contributor Covenant](https://contributor-covenant.org/version/2/0/code_of_conduct/) Code of Conduct.

### License

GNU Affero General Public License v3.0