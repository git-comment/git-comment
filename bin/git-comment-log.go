package main

import (
	// "errors"
	"fmt"
	gitc "git_comment"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

const (
	errorPrefix = "git-comment-log"
)

var (
	buildVersion string
	app          = kp.New("git-comment-log", "List git commit comments")
	revision     = app.Arg("revision range", "Filter comments to comments on commits from the specified range").String()
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, errorPrefix)
	comments, err := gitc.CommentsOnCommit(pwd, revision)
	app.FatalIfError(err, errorPrefix)
	// concat all output
	// check length
	// check $LINES
	// if length > $LINES open pager with content
	// otherwise print content
	for i := 0; i < len(comments); i++ {
		comment := comments[i]
		fmt.Println(*comment.ID)
		fmt.Println(comment.Serialize())
	}
}
