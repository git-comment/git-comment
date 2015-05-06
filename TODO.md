## Components

- git-comment: creates/edits/deletes annotations
- git-comment-log: lists annotations by commit
- git-comment-grep: searches annotation content

## Comment file format

Components:

- author
- timestamp
- committer
- committer timestamp
- ID
- commit
- filepath and line
- content

## Things that need to work

- lookup by commit
- lookup by ID
- comment content search
- list comments by tree
  * sort by date
  * show diff lines in log
- pre-comment/post-comment hooks
  * disallow "invalid" comments
  * push comments to external services
- export from gh/bb/other
  * support setting a custom creation time
- documentation for each binary

## Finding a comment by commit
* nest comment refs under commit hashes
  - refs/notes/[commit hash]/[comment hash]
* include hashes of all comments for a commit in a ref
  - refs/notes/[commit hash] contains delimited comment hashes
