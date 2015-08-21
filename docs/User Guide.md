# git-comment User Guide

## Contents

* Command-line Interface
  - Adding Comments
  - Viewing Comments
  - Sharing Comments
  - Searching for Comments
* Web Interface
  - Adding Comments
  - Viewing Comments
  - Searching for Comments
* Editor Interface
  - Adding Comments
  - Viewing Comments
* Configuration
  - Identity
  - Hooks

## Command-line interface

### Adding Comments

The core binary can add comments to commits, optionally with a file and
line reference. Creating a comment without a supplied message opens the
default git editor.

The name and email used as the comment author and committer identities are
shared from git environment variables `GIT_AUTHOR_IDENT` and
`GIT_COMMITTER_IDENT` respectively. See `git help var` to learn more. To
override the author identity, use the `--author` flag.

```
git comment [-m <msg>] [--amend <comment>] [-c <commit>]
            [--author=<author>] [<filepath:line>]
git comment --delete <comment>
git comment --help
git comment --version
```

Comment text can be any number of lines, or use any formatting syntax,
though plain text formats like markdown and textile ensure the best
readability for command-line and web-based interfaces.

`git-comment` supports prepopulating a comment's content from a file based
on the configuration option `comment.template` or
`$HOME/.gitcommenttemplate` if available in that order.

### Viewing Comments

`git-comment-log` prints comments and associated diffs by commit or tree.

```
git comment-log [--pretty <format>] [<revision range>]
git comment-log --help
git comment-log --version

```

### Sharing Comments

#### Merge (Central Remote) Workflow

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

#### Patch (No Remote) Workflow

##### Creating a patch

TODO

##### Applying a patch
TODO

### Searching for Comments

`git-comment-grep` prints comments containing text.

```
git comment-grep find <text>
git comment-grep index
git comment-grep --help
git comment-grep --version
```


## Web Interface

`git-comment-web` starts a web server hosting a friendly interface for
editing comments on diffs

```
git comment-web [--port <port>]
git comment-web --help
git comment-web --version
```

### Adding Comments

### Viewing Comments

### Searching for Comments



## Configuration

### Identity

The comment author name and email is shared from the user's git 
configuration properties `user.name` and `user.email`.

### Hooks

`git-comment` supports running git hooks before and after creating a
comment, named `pre-comment` and `post-comment` respectively. These
executable files can be configured to cancel comment creation by exiting
with a non-zero status code.
