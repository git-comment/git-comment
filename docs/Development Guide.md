# git-comment Developer Guide

## Contents

* Comment Specification
* Comment Reference Specification
* Building an Editor Integration
* Building Import Scripts
* Building Export Scripts
* Reference

## Comment Specification

Comments are regular git objects, stored in a format similar to tags,
with the addition of file and line references and a flag for deletion. A
reference for the comment is added in refs/comments for lookup by
commit. An example comment object would look something like:

```
commit 0155eb4229851634a0f03eb265b69f5a2d56f341
file src/example.txt:12
author Delisa Mason <name@example.com> 1243040974 -0900
amender Delisa Mason <name@example.com> 1243040974 -0900

Too many levels of indentation here.
```

The valid fields are (in this order):

* `commit` : full hash of the attached commit
* `file` : path (relative to repository root) and line number of the
  comment or an empty string
* `author` : the writer of the comment's name, email, and time of
  authorship
* `amender` : person who checked the comment into git

After these fields, there is a single empty line, then the content of
the comment, which can be any number of lines and contain any UTF-8
characters.

## Comment Reference Specification

A comment reference has the format `refs/comments/[first four character of commit hash]/[rest of commit]/[comment hash]`, like:

```
refs/comments/0155/eb4229851634a0f03eb265b69f5a2d56f341/f9da8cdd40bbce4c7bd1aa4e46608107184bd91c
```

The contents of the reference are the comment hash.

## Building an Editor Integration

## Building Import Scripts

## Building Export Scripts

## Reference

* [Git Internals Guide](http://www.git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain)
* [libgit2 Documentation](https://libgit2.github.com) (and it's [Go bindings](http://godoc.org/github.com/libgit2/git2go))
* [gitcore-tutorial(7)](https://www.kernel.org/pub/software/scm/git/docs/gitcore-tutorial.html)
* [gitrepository-layout(5)](https://www.kernel.org/pub/software/scm/git/docs/gitrepository-layout.html)
* [gitcli(7)](https://www.kernel.org/pub/software/scm/git/docs/gitcli.html)
