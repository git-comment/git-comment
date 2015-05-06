# NAME

  git-comment: annotations for git commits

# SYNOPSIS

    git comment [-m <msg>] [--amend <ID>] [<commit>] [<filepath:line>]
    git comment -d <ID>

# DESCRIPTION

  Adds, removes, or amends comments attached to git commits, without
changing the commits themselves.

# OPTIONS

    --amend <ID>
       Edit the existing comment with a given ID. No modifications are
made if the a commit is specified but no comment with a given ID is
attached.

    -d <ID>, --delete=<ID>
       Remove comment with a given ID

    -m <msg>, --message=<msg>
       Use the given message instead of prompting for message content.

# ARGUMENTS

    <commit>
       A commit hash

    <filepath:line>
      A reference to a file and line number, to make the annotation more
specific.

# DISCUSSION

  Comments are blobs containing extra information about a commit and
potentially a changed line within the commit.

# AUTHORS

  git-comment was written and is maintained by Delisa Mason.

# SEE ALSO

  git-comment-log(1), git-comment-grep(1)
