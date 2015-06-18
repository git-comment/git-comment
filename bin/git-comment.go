package main

import (
	"errors"
	"fmt"
	gitc "git_comment"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	noMessageProvided = "Aborting comment, no message provided"
	editorFailed      = "Failed to launch preferred editor"
	defaultMessage    = "\n# Enter comment content\n# Lines beginning with '#' will be stripped"
	errorPrefix       = "bin"
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
	handleError(err)
	if len(*remoteToConfig) > 0 {
		handleError(gitc.ConfigureRemoteForComments(pwd, *remoteToConfig))
		fmt.Printf("Remote '%v' updated\n", *remoteToConfig)
	} else if len(*deleteID) > 0 {
		handleError(gitc.DeleteComment(pwd, *deleteID))
		fmt.Println("Comment deleted")
	} else {
		editComment(pwd)
	}
}

func handleError(err error) {
	app.FatalIfError(err, errorPrefix)
}

func editComment(pwd string) {
	var commit *string = nil
	if len(*message) == 0 {
		*message = getMessageFromEditor(pwd)
	}
	if len(*amendID) > 0 {
		id, err := gitc.UpdateComment(pwd, *amendID, *message)
		handleError(err)
		fmt.Printf("[%v] Comment updated\n", (*id)[:7])
	} else {
		id, err := gitc.CreateComment(pwd, commit, gitc.CreateFileRef(*fileref), *message)
		handleError(err)
		fmt.Printf("[%v] Comment created\n", (*id)[:7])
	}
}

func getMessageFromEditor(pwd string) string {
	editor := gitc.ConfiguredEditor(pwd)
	file, err := ioutil.TempFile("", "gitc")
	handleError(err)
	path := file.Name()
	file.Write([]byte(defaultMessage))
	file.Close()
	execCommand(*editor, path)
	content, err := ioutil.ReadFile(path)
	os.Remove(path)
	handleError(err)
	return sanitizeMessage(string(content))
}

func sanitizeMessage(message string) string {
	reg, err := regexp.Compile("(?m)^#.*$")
	handleError(err)
	stripped := reg.ReplaceAllString(message, "")
	content := strings.TrimSpace(stripped)
	if len(content) == 0 {
		handleError(errors.New(noMessageProvided))
	}
	return content
}

func execCommand(program string, args ...string) {
	cmd := exec.Command(program, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	handleError(cmd.Run())
}
