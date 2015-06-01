# git-comment

Annotations for git commits. Open and distributed collaboration around
code and other version-controlled text and resources.

## why?

Comments on git commits are siloed into various web services or backed
up externally, when all repo history can be stored together and
accessible offline.


## use cases

* Storing git comment history in every local copy
* Unified and open format for web services to store comments
* Reviewing and annotating code while offline
* Attaching release notes and other longform metadata to branches or
  tags
* Viewing comments and associated context diffs offline
* Pre- and post-comment hooks for validation with web services
* Searching for comments by commit or text while offline
* Exporting comments from web services into local backups

To meet these requirements, `git-comment` ships with a few helpful
components:

### `git-comment`

The core binary can add comments to commits, optionally with a file and
line reference.

```
git comment [-m <msg>] [--amend <comment>] [<commit>] [<filepath:line>]
git comment --delete <comment>
git comment --help
git comment -v
```

### `git-comment-log`

View comments and associated diffs by commit or tree.

```
git comment-log [<revision range>]
git comment-log --help
git comment-log -v
```

### `git-comment-grep`

Print comments matching a pattern.

```
git comment-grep <pattern>
git comment-grep --help
git comment-grep -v
```

### `pre_comment` and `post_comment` git hooks

Execute commands before or after creating comments, aborting the
operation when the scripts fail.

## how?

Comments are regular git objects, stored in a format similar to tags,
with the addition of file and line references and a flag for deletion. A
reference for the comment is added in refs/comments for lookup by
commit. An example comment object would look something like:

```
commit 0155eb4229851634a0f03eb265b69f5a2d56f341
file src/example.txt:12
author Delisa Mason <name@example.com>
created 1243040974 -0900
amender Delisa Mason <name@example.com>
amended 1243040974 -0900

Too many levels of indentation here.
```

Comment text can be any number of lines, or use any formatting syntax,
though plain text formats like markdown and textile ensure the best
readability for command-line and web-based interfaces.

## license

Copyright (c) 2015, Delisa Mason <delisam@acm.org>. All rights reserved.

`git-comment` is licensed under the BSD 2-clause license, and detailed in the `LICENSE` file.

## contributing

`git-comment` is written in [Go](http://golang.org) and tested using [assert](https://github.com/stvp/assert).

Dependencies are listed in the `Makefile` and can be installed by running `make bootstrap`. The default command installs `libgit2` via [Homebrew](http://brew.sh), but it can be substituted for any other suitable package manager or installation method.

The manual is written using [pod2man](http://perldoc.perl.org/pod2man.html), which should be available on most GNU/Linux and OS X distributions by default. Changes should be documented with friendliness in mind.
