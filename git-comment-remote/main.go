package main

import (
	"fmt"
	"github.com/kylef/result.go/src/result"
	kp "gopkg.in/alecthomas/kingpin.v2"
	gc "libgitcomment"
	"os"
)

var (
	buildVersion  string
	app           = kp.New("git-comment-remote", "Helper commands for the merge workflow")
	configCmd     = app.Command("config", "Configure remote to fetch and push comments by default")
	configRemote  = configCmd.Arg("remote", "Remote to configure").Required().String()
	deleteCmd     = app.Command("delete", "Delete remote copy of a comment")
	deleteRemote  = deleteCmd.Arg("remote", "Remote from which to delete comment").Required().String()
	deleteComment = deleteCmd.Arg("comment", "Comment to delete").Required().String()
)

func main() {
	app.Version(buildVersion)
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	fatalIfError(app, gc.VersionCheck(pwd, buildVersion), "version")
	switch kp.MustParse(app.Parse(os.Args[1:])) {
	case "config":
		app.FatalIfError(gc.ConfigureRemoteForComments(pwd, *configRemote).Failure, "git")
		fmt.Printf("Remote '%v' updated\n", *configRemote)
	case "delete":
		app.FatalIfError(gc.DeleteRemoteComment(pwd, *deleteRemote, *deleteComment).Failure, "git")
		fmt.Printf("Remote comment reference deleted\n")
	}
}

// Return the success value, otherwise kill the app with
// the error code specified
func fatalIfError(app *kp.Application, r result.Result, code string) interface{} {
	app.FatalIfError(r.Failure, code)
	return r.Success
}
