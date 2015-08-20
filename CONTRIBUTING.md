# Contribution Guide

## Development Dependencies

Development dependencies are listed in the `INSTALL` file and include [Go](http://golang.org), [libgit2](https://libgit2.github.com), [make](https://www.gnu.com/software/make), and various libraries available via `go get`. The command `make bootstrap` installs the library dependencies.

The manual is written using [pod2man](http://perldoc.perl.org/pod2man.html), which should be available on most GNU/Linux and OS X distributions by default. Changes should be documented with friendliness in mind.

## Submitting a Change

0. Fork and clone the repo

1. Install the requirements listed in the `INSTALL` file

2. Run `make bootstrap` to get the library dependencies

3. Make your changes, adding tests where necessary. `git-comment` is tested using [assert](https://github.com/stvp/assert).

4. Ensure all tests pass by running `make test`

5. Commit your changes, writing a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)

6. Push to your fork and [submit a pull request](https://github.com/kattrali/git-comment/compare/)

## Useful References

* [Git Internals Guide](http://www.git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain)
* [libgit2 Documentation](https://libgit2.github.com) (and it's [Go bindings](http://godoc.org/github.com/libgit2/git2go))
* [gitcore-tutorial(7)](https://www.kernel.org/pub/software/scm/git/docs/gitcore-tutorial.html)
* [gitrepository-layout(5)](https://www.kernel.org/pub/software/scm/git/docs/gitrepository-layout.html)
* [gitcli(7)](https://www.kernel.org/pub/software/scm/git/docs/gitcli.html)
