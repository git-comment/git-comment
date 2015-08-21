# The Big TODO List

## Integrations

### Import/export from web services

Synchronize comments with external services.

* Fetch comments from the service
* Push comments to the service
* Push and fetch only as needed, avoiding duplication automatically

Services:

* [ ] [GitHub](https://developer.github.com/v3/pulls/comments/#list-comments-in-a-repository)
* [ ] [Bitbucket](https://confluence.atlassian.com/display/BITBUCKET/pullrequests+Resource+1.0)
* [ ] [GitLab](http://doc.gitlab.com/ce/api/notes.html)

### Editor plugins

An editor plugin to comment while viewing a diff and view existing comments.

* Add comments
* View comments
* Fold comments
* View gutter icon indicating a comment exists at a line
* Update the git-comment user guide

Editors:

* [ ] Vim
* [ ] Sublime Text
* [ ] Notepad++
* [ ] TextMate
* [ ] Xcode

### Continuous Integration and Deployment

Configure running the test suite on various platforms, along with generating binary
distributions.

* [ ] [Windows](http://www.appveyor.com)
* [ ] Debian
* [ ] RPM
* [ ] OS X
* [ ] OS X via Homebrew

## Documentation

### Usage tutorials for creating, sharing, and viewing comments

* [ ] Merge flow
* [ ] Patch flow

## Minor

### Support pre-/post- comment hooks for integrations

Optionally run an executable file before and after creating a comment, to
optimize actions such as search indexing and synchronization with external
services.

### Create and apply comments from patches

    $ git comment-patch format master > comments.patch
    $ git comment-patch apply comments.patch

### Advanced filters for `git-comment-log`

* [ ] Filter by author
* [ ] Filter by date
* [ ] Filter by file

### Automatically index comments for search

* [ ] Index after comment creation
* [ ] Index after fetching comments from a remote
* [ ] Index after applying patch

## Major

### Support signed comments

Use a format similar to `git tag`:

    -s, --sign
      Make a GPG-signed comment, using the default e-mail address's key.

    -u <key-id>, --local-user=<key-id>
      Make a GPG-signed comment, using the given key.

Here be dragons.

### Better comment storage

Change comment storage format to avoid reference bloat, perhaps with real
trees. The current implementation requires one reference per comment, making
`git remote show` less than useful.

### Comment grouping unit

Write a specification for comment grouping as a part of a unified
merge request/issue tracking system.
