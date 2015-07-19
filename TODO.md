## Components

- git-comment: creates/edits/deletes annotations
- git-comment-log: lists annotations by commit
- git-comment-grep: searches annotation content
- import/export scripts
- editor plugins

## Things that need to work

- [x] lookup by commit
- [x] lookup by ID
- [ ] save git-comment version in use
  * per comment?
  * per repo?
  * validate before making changes/reading comments
- [ ] comment content search
- [x] list comments by tree
  * [x] sort by date
  * [x] show diff lines before/after in log
- [ ] pre-comment/post-comment hooks
  * [ ] disallow "invalid" comments
  * [ ] change comments before persisting
- [ ] export from web services
  * support setting a custom creation time, author (but not committer?)
  * API support/scraping
    - [ ] [Bitbucket](https://confluence.atlassian.com/display/BITBUCKET/pullrequests+Resource)
    - [ ] [Github](https://developer.github.com/v3/repos/comments/#list-comments-for-a-single-commit)
    - [ ] [Gitlab / Gitorious](http://doc.gitlab.com/ce/api/notes.html)
- [ ] documentation for each binary
- [x] hard limit for number of comments on a commit (2^12?)
- [ ] editor integration
  * [ ] Reference plugin for vim
  * [ ] Reference plugin for xcode ?
- [ ] signed comments, a la `git tag`:

```
 -s, --sign
    Make a GPG-signed comment, using the default e-mail address's
key.

 -u <key-id>, --local-user=<key-id>
    Make a GPG-signed comment, using the given key.
```

## Finding a comment by commit
* include hashes of all comments for a commit in a ref
  - refs/comments/[prefix]/[rest of commit hash]/[comment id] contains newline delimited comment hashes
    * Pros:
      - Easy to find a comment by commit hash
      - New comments can have refs added easily
    * Cons:
      - No inherent comment ordering
