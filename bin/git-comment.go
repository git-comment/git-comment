package main

import (
	"fmt"
	gc "git_comment"
	gx "git_comment/exec"
	gg "git_comment/git"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	buildVersion string
	app          = kp.New("git-comment", "Add comments to commits and diffs within git repositories")
	message      = app.Flag("message", "comment content").Short('m').String()
	amendID      = app.Flag("amend", "ID of a comment to amend").String()
	deleteID     = app.Flag("delete", "ID of a comment to delete").String()
	commit       = app.Flag("commit", "ID of a commit to annotate").Short('c').String()
	update       = app.Flag("update", "Upgrade repository to use current version of git-comment").Bool()
	fileref      = app.Arg("file:line", "File and line number to annotate").String()
	markDeleted  = app.Flag("mark-deleted-line", "Add comment to the deleted version of the file and line number").Bool()
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	if *update {
		gc.VersionUpdate(pwd, buildVersion)
		return
	}
	gx.FatalIfError(app, gc.VersionCheck(pwd, buildVersion), "version")
	if len(*deleteID) > 0 {
		app.FatalIfError(gc.DeleteComment(pwd, *deleteID).Failure, "git")
		fmt.Println("Comment deleted")
	} else {
		editComment(pwd)
	}
}

func editComment(pwd string) {
	parsedCommit := gx.FatalIfError(app, gg.ResolvedCommit(pwd, commit), "git")
	if len(*message) == 0 {
		*message = gx.GetMessageFromEditor(app, pwd)
	}
	if len(*amendID) > 0 {
		id := gx.FatalIfError(app, gc.UpdateComment(pwd, *amendID, *message), "git")
		fmt.Printf("[%v] Comment updated\n", (*id.(*string))[:7])
	} else {
		ref := gc.CreateFileRef(*fileref, *markDeleted)
		id := gx.FatalIfError(app, gc.CreateComment(pwd, parsedCommit.(*string), ref, *message), "git")
		hash := *(id.(*string))
		fmt.Printf("[%v] Comment created\n", hash[:7])
	}
}
