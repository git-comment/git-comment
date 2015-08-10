# git-comment User Guide

## Contents

* Getting Started
* Adding Comments
* Viewing Comments
* Sharing Comments
* Searching for Comments
* Hooks

## Getting Started

### Installation
* Pre-built binaries
* Source

#### Identity
* Name, email, signature

### Migrating from Hosted Services
* Import from GitHub
* Import from Bitbucket

## Adding Comments

### Command-line interface

The core binary can add comments to commits, optionally with a file and
line reference. Creating a comment without a supplied message opens the
default git editor.

The name and email used as the comment author and committer identities are
shared from git environment variables `GIT_AUTHOR_IDENT` and
`GIT_COMMITTER_IDENT` respectively. See `git help var` to learn more.

```
git comment [-m <msg>] [--amend <comment>] [-c <commit>] [<filepath:line>]
git comment --delete <comment>
git comment --help
git comment --version
```

Comment text can be any number of lines, or use any formatting syntax,
though plain text formats like markdown and textile ensure the best
readability for command-line and web-based interfaces.

`git-comment` supports prepopulating comment content from a file based
on configuration option `comment.template` or
`$HOME/.gitcommenttemplate` if available in that order.

### Editor integration
### Web interface

`git-comment-web` starts a web server hosting a friendly interface for
editing comments on diffs

```
git comment-web [<revision range>] [--port <port>]
git comment-web --help
git comment-web --version
```

## Viewing Comments

### Command-line interface

`git-comment-log` prints comments and associated diffs by commit or tree.

```
git comment-log [--pretty <format>] [<revision range>]
git comment-log --help
git comment-log --version

```

### Editor integration
### Web interface

## Sharing Comments

### Merge (Central Remote) Workflow

Comments can be pushed to a central server using

    git push <remote> 'refs/comments/*'

Or fetched via

    git fetch origin '+refs/comments/*:refs/remotes/<remote>/comments/*'

Alternately, a remote can be configured to push and fetch comments by
default. The `git-comment` suite includes commands for remote
configuration and comment pruning.

```
git comment-remote config <remote>
git comment-remote delete <remote> <comment>
git comment-remote --help
git comment-remote --version
```

Using `git comment-remote config` adds fetch and push refspecs to a
remote for comments. After use, using `git fetch` or `git push` will
fetch or push new comments to the remote by default.

Deleting the reference to a remote comment renders it inaccessible by
other users who have not yet fetched it. `git comment-remote delete`
deletes the remote reference. Note that other users who have already
fetched the comment could repush it unless blocked by a push hook on the
remote side.

### No Server Workflow

### Creating a patch
### Applying a patch
### Viewing comments from a patch

## Searching for Comments

`git-comment-grep` prints comments containing text.

```
git comment-grep find <text>
git comment-grep index
git comment-grep --help
git comment-grep --version
```

## Hooks

`git-comment` supports running git hooks before and after creating a
comment, named `pre-comment` and `post-comment` respectively. These
files are shell scripts which can be configured to cancel comment
creation by exiting with a non-zero status code.
