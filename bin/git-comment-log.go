package main

import (
	// "errors"
	//gitc "git_comment"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	buildVersion string
	app          = kp.New("git-comment-log", "List git commit annotations")
	revision     = app.Arg("revision range", "Filter comments to comments on commits from the specified range").String()
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
}
