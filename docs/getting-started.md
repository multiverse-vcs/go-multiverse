# Getting Started

Multiverse uses peer-to-peer network protocols to exchange data. 

Instead of using a centralized server, every participant runs a peer node.

Start your peer node and keep it running in a separate terminal.

```bash
$ multi daemon
```

With the peer node up and running you are ready to create your first repository.

Create an empty directory and initialize the repository.

```bash
$ mkdir my_project
$ cd my_project
$ multi init
```

After you have made some changes check the repository status.

Notice that all files are tracked by default.

```bash
$ multi status
```

Once you are happy with your work, commit the changes to the repository.

A commit contains a snapshot of the files in your repository.

```bash
# describe your changes with an optional message
$ multi commit --message "add initial code"
```

To share your code with others create a new repository on the peer node.

Ensure there is no confidential info. Everything shared on the main network is public.

```bash
$ multi repo create my_project
```

If the repository was created successfully its remote path will be printed.

```bash
# your remote path will be different than this one
12D3KooWFRfidCtkUkViUMTnoEoVtzDLmdCix8XUmVCoZcATLixG/my_project
```

A remote path is how your repository is uniquely identified on the Multiverse network.

The remote path consists of your unique peer identifier followed by the repository name.

To push changes to the new remote repository, first add the remote to the repository.

```bash
# replace the remote path with your unique remote path
multi remote create origin 12D3KooWFRfidCtkUkViUMTnoEoVtzDLmdCix8XUmVCoZcATLixG/my_project
```

Next set the branch remote to the one created above.

```bash
multi branch set remote origin
```

Finally push the changes.

```bash
multi push
```
