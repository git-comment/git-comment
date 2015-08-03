# git-comment

Comments for git commits and files. Open and distributed collaboration
around code and other version-controlled text and resources.

## Why?

Comments on git commits are siloed into various web services or backed
up externally, when all repo history can be stored together and
accessible offline.

## Use Cases

* Storing git comment history in every local copy
* Unified and open format for web services to store comments
* Reviewing and annotating code while offline
* Viewing comments and associated context diffs offline
* Pre- and post-comment hooks for validation with web services
* Searching for comments by commit or text while offline
* Exporting comments from web services into local backups

To meet these requirements, `git-comment` ships with a few helpful
components:

* `git-comment`: adds comments
* `git-comment-log`: prints comments inline with diffs
* `git-comment-grep`: searches comment content for text
* `git-comment-web`: launches a web server hosting a friendly web UI for
  comment editing
* `git-comment-remote`: helpful tools for working with a remote server
  with git comment, like configuring remotes to push and pull comments
  by default, indexing comments for search after push, and deleting remote
  comments
* `git-comment-patch`: formats and applies comment patch files, for
  working with a fully decentralized flow

More information and usage is available in the manual or the User Guide.

### Import/Export scripts

Retrieve all commits from external services including GitHub and
BitBucket. Check the `scripts` directory.

### Editor integrations

A reference plugin for vim can be found [here](). Submissions for other
editor integrations are encouraged!

## Installation

### Binaries

Binaries for several platforms are available on the [downloads
page](https://github.com/kattrali/git-comment/releases)

### Source

Instructions for source installation are provided in the `INSTALL` file.

## Help

* User Guide
* `#git-comment` on freenode
* Search open and closed issues for similar problems
* Open an issue

## License

Copyright (c) 2015, Delisa Mason <delisam@acm.org>. All rights reserved.

`git-comment` uses the BSD license as detailed in the `LICENSE` file.

## Contributing

The `CONTRIBUTING.md` file details the setup process for building
`git-comment` from source and submitting a change.

