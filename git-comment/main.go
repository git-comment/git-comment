package main

import (
	"fmt"
	gg "git"
	"github.com/kylef/result.go/src/result"
	kp "gopkg.in/alecthomas/kingpin.v2"
	gc "libgitcomment"
	"os"
)

var (
	buildVersion string
	app          = kp.New("git-comment", "Add comments to commits and diffs within git repositories")
	message      = app.Flag("message", "comment content").Short('m').String()
	amendID      = app.Flag("amend", "ID of a comment to amend").String()
	deleteID     = app.Flag("delete", "ID of a comment to delete").String()
	commit       = app.Flag("commit", "ID of a commit to annotate").Short('c').String()
	author       = app.Flag("author", "Override the comment author").String()
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
	fatalIfError(app, gc.VersionCheck(pwd, buildVersion), "version")
	if len(*deleteID) > 0 {
		app.FatalIfError(gc.DeleteComment(pwd, *deleteID).Failure, "git")
		fmt.Println("Comment deleted")
	} else {
		editComment(pwd)
	}
}

func editComment(pwd string) {
	resolved := fatalIfError(app, gg.ResolvedCommit(pwd, *commit), "git")
	parsedCommit := resolved.(*string)
	if len(*message) == 0 {
		*message = getMessageFromEditor(app, pwd)
	}
	if len(*amendID) > 0 {
		id := fatalIfError(app, gc.UpdateComment(pwd, *amendID, *author, *message), "git")
		fmt.Printf("[%v] Comment updated\n", (*id.(*string))[:7])
	} else {
		ref := gc.CreateFileRef(*fileref, *markDeleted)
		id := fatalIfError(app, gc.CreateComment(pwd, *parsedCommit, *author, *message, ref), "git")

		hash := *(id.(*string))
		fmt.Printf("[%v] Comment created\n", hash[:7])
	}
}

// Return the success value, otherwise kill the app with
// the error code specified
func fatalIfError(app *kp.Application, r result.Result, code string) interface{} {
	app.FatalIfError(r.Failure, code)
	return r.Success
}
