package main

import (
	"errors"
	"fmt"
	gitc "git_comment"
	ex "git_comment/exec"
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
	commit         = app.Arg("commit", "ID of a commit to annotate").String()
	fileref        = app.Arg("file:line", "File and line number to annotate").String()
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	if len(*remoteToConfig) > 0 {
		app.FatalIfError(gitc.ConfigureRemoteForComments(pwd, *remoteToConfig), "git")
		fmt.Printf("Remote '%v' updated\n", *remoteToConfig)
	} else if len(*deleteID) > 0 {
		app.FatalIfError(gitc.DeleteComment(pwd, *deleteID), "git")
		fmt.Println("Comment deleted")
	} else {
		editComment(pwd)
	}
}

func editComment(pwd string) {
	var commit *string = nil
	if len(*message) == 0 {
		*message = getMessageFromEditor(pwd)
	}
	if len(*amendID) > 0 {
		id, err := gitc.UpdateComment(pwd, *amendID, *message)
		app.FatalIfError(err, "git")
		fmt.Printf("[%v] Comment updated\n", (*id)[:7])
	} else {
		id, err := gitc.CreateComment(pwd, commit, gitc.CreateFileRef(*fileref), *message)
		app.FatalIfError(err, "git")
		fmt.Printf("[%v] Comment created\n", (*id)[:7])
	}
}

func getMessageFromEditor(pwd string) string {
	editor := gitc.ConfiguredEditor(pwd)
	file, err := ioutil.TempFile("", "gitc")
	app.FatalIfError(err, "io")
	path := file.Name()
	file.Write([]byte(gitc.DefaultMessageTemplate))
	file.Close()
	err = ex.ExecCommand(*editor, path)
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
