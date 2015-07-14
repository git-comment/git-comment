package main

import (
	"errors"
	"fmt"
	gc "git_comment"
	ex "git_comment/exec"
	"github.com/kylef/result.go/src/result"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	editorFailed      = "Failed to launch preferred editor"
	noMessageProvided = "Aborting comment, no message provided"
)

var (
	buildVersion   string
	app            = kp.New("git-comment", "Add comments to commits and diffs within git repositories")
	message        = app.Flag("message", "comment content").Short('m').String()
	amendID        = app.Flag("amend", "ID of a comment to amend").String()
	deleteID       = app.Flag("delete", "ID of a comment to delete").String()
	remoteToConfig = app.Flag("configure-remote", "remote to configure for fetch and pushing comments").String()
	commit         = app.Flag("commit", "ID of a commit to annotate").Short('c').String()
	fileref        = app.Arg("file:line", "File and line number to annotate").String()
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	if len(*remoteToConfig) > 0 {
		app.FatalIfError(gc.ConfigureRemoteForComments(pwd, *remoteToConfig).Failure, "git")
		fmt.Printf("Remote '%v' updated\n", *remoteToConfig)
	} else if len(*deleteID) > 0 {
		app.FatalIfError(gc.DeleteComment(pwd, *deleteID).Failure, "git")
		fmt.Println("Comment deleted")
	} else {
		editComment(pwd)
	}
}

func editComment(pwd string) {
	parsedCommit := fatalIfError(gc.ResolvedCommit(pwd, commit), "git")
	if len(*message) == 0 {
		*message = getMessageFromEditor(pwd)
	}
	if len(*amendID) > 0 {
		id := fatalIfError(gc.UpdateComment(pwd, *amendID, *message), "git")
		fmt.Printf("[%v] Comment updated\n", (*id.(*string))[:7])
	} else {
		id := fatalIfError(gc.CreateComment(pwd,
			parsedCommit.(*string),
			gc.CreateFileRef(*fileref), *message), "git")
		hash := *(id.(*string))
		fmt.Printf("[%v] Comment created\n", hash[:7])
	}
}

func fatalIfError(r result.Result, code string) interface{} {
	app.FatalIfError(r.Failure, code)
	return r.Success
}

func getMessageFromEditor(pwd string) string {
	editor := gc.ConfiguredEditor(pwd)
	file, err := ioutil.TempFile("", "gitc")
	app.FatalIfError(err, "io")
	path := file.Name()
	file.Write([]byte(gc.DefaultMessageTemplate))
	file.Close()
	err = ex.ExecCommand(editor, path)
	app.FatalIfError(err, "io")
	content, err := ioutil.ReadFile(path)
	os.Remove(path)
	app.FatalIfError(err, "io")
	return sanitizeMessage(string(content))
}

func sanitizeMessage(message string) string {
	reg, err := regexp.Compile("(?m)^#.*$")
	app.FatalIfError(err, "regex")
	stripped := reg.ReplaceAllString(message, "")
	content := strings.TrimSpace(stripped)
	if len(content) == 0 {
		app.FatalIfError(errors.New(noMessageProvided), "git-comment")
	}
	return content
}
