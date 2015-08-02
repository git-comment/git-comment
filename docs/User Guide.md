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

### Configuration

#### Identity
* Name, email, signature

#### Formatting
* Comment template

### Migrating from Hosted Services
* Import from GitHub
* Import from Bitbucket

## Adding Comments

### Command-line interface

The core binary can add comments to commits, optionally with a file and
line reference. It includes a helper command (`--configure-remote`) for
fetching and pushing comments by default with other refs. Creating a
comment without a supplied message opens the default git editor.

```
git comment [-m <msg>] [--amend <comment>] [-c <commit>] [<filepath:line>]
git comment --delete <comment>
git comment --configure-remote <remote>
git comment --help
git comment --version
```

Comment text can be any number of lines, or use any formatting syntax,
though plain text formats like markdown and textile ensure the best
readability for command-line and web-based interfaces.

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

### Central Remote Workflow
* Push/pull comments by default

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
