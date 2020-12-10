# Multiverse

[![Go Report Card](https://goreportcard.com/badge/github.com/multiverse-vcs/go-multiverse)](https://goreportcard.com/report/github.com/multiverse-vcs/go-multiverse)

Multiverse is a decentralized version control system that enables peer-to-peer software development.

### Features

- Offline first. Sync changes and continue working offline.
- Own your data. Host your code on a local IPFS node.
- Work with peers locally. No internet connection required.

### Install

Go 1.14 or higher is required.

```bash
$ git clone https://github.com/multiverse-vcs/go-multiverse.git
$ cd go-multiverse
$ make install
```

### Quick Start

A local IPFS node is required to use Multiverse. Instructions for setting up IPFS can be found [here](https://ipfs.io/#install).

#### Creating a new repo

All files in your project directory are tracked by default.

```bash
$ mkdir my_project
$ cd my_project
$ multi init
```

#### Ignoring files

Create a `.multignore` in your project root and add a new line for each rule.

```
# ignore all exe files
*.exe

# ignore the bin directory
bin
```

#### Viewing repo status

Check your changes before recording them.

```bash
$ multi status
```

#### Recording changes

Commits can include an optional message to describe what has been changed.

```bash
$ multi commit -m "changed some files"
```

After a commit has been recorded its [CID](https://docs.ipfs.io/concepts/content-addressing/#identifier-formats) will be printed.

```bash
$ bafyreick5vs3ayhxxotfeo2udpbmnzj3ekmhmofncppbfqtpwr6bs44vvy
```

#### Sharing changes

Changes will be pinned to your local IPFS node. *Warning all shared content will be public*.

```bash
$ multi push
```

#### Merging changes

Sharing changes requires an external communication channel such as an instant messanger or email. *Only trust CIDs from secure sources*.

```bash
$ multi pull bafyreick5vs3ayhxxotfeo2udpbmnzj3ekmhmofncppbfqtpwr6bs44vvy
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
   history   Print change history
   init      Initialize a new repo
   pull      Merge changes into the current branch
   push      Upload changes to a remote
   status    Print repo status
   switch    Change branches
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### Contributing

Found a bug or have a feature request? [Open an issue](https://github.com/multiverse-vcs/go-multiverse/issues/new).

Multiverse follows the [Contributor Covenant](https://contributor-covenant.org/version/2/0/code_of_conduct/) Code of Conduct.

### License

GNU Affero General Public License v3.0
