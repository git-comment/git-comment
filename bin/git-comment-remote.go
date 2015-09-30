package main

import (
	"fmt"
	"os"

	gc "github.com/git-comment/git-comment/libgitcomment"
	gx "github.com/git-comment/git-comment/libgitcomment/exec"
	kp "gopkg.in/alecthomas/kingpin.v2"
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
	gx.FatalIfError(app, gc.VersionCheck(pwd, buildVersion), "version")
	switch kp.MustParse(app.Parse(os.Args[1:])) {
	case "config":
		app.FatalIfError(gc.ConfigureRemoteForComments(pwd, *configRemote).Failure, "git")
		fmt.Printf("Remote '%v' updated\n", *configRemote)
	case "delete":
		app.FatalIfError(gc.DeleteRemoteComment(pwd, *deleteRemote, *deleteComment).Failure, "git")
		fmt.Printf("Remote comment reference deleted\n")
	}
}
