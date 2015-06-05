## Components

- git-comment: creates/edits/deletes annotations
- git-comment-log: lists annotations by commit
- git-comment-grep: searches annotation content
- import/export scripts

## Things that need to work

- lookup by commit
- lookup by ID
- save git-comment version in use
  * per comment?
  * per repo?
  * validate before making changes/reading comments
- comment content search
- list comments by tree
  * sort by date
  * show diff lines in log
- pre-comment/post-comment hooks
  * disallow "invalid" comments
  * change comments before persisting
- export from gh/bb/other
  * support setting a custom creation time, author (but not committer?)
  * API support/scraping
    - https://confluence.atlassian.com/display/BITBUCKET/pullrequests+Resource
    - https://developer.github.com/v3/repos/comments/#list-comments-for-a-single-commit
    - http://doc.gitlab.com/ce/api/notes.html
- documentation for each binary
- signed comments, a la `git tag`:

```
 -s, --sign
    Make a GPG-signed comment, using the default e-mail address's
key.

 -u <key-id>, --local-user=<key-id>
    Make a GPG-signed comment, using the given key.
```

## Finding a comment by commit
* nest comment refs under commit hashes
  - refs/notes/[commit hash]/[comment hash]
* include hashes of all comments for a commit in a ref
  - refs/comments/[prefix]/[rest of commit hash] contains newline delimited comment hashes
    * Pros:
      - Easy to find a comment by commit hash
      - New comments can have refs added easily
    * Cons:
      - Deleting a comment is slightly intensive, as the ref file must
        be scanned for the comment ID
      - No inherent comment ordering ?
