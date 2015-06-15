package main

import (
	"errors"
	"fmt"
	gitc "git_comment"
	goopt "github.com/droundy/goopt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	tooManyArguments  = "Too many arguments provided"
	noMessageProvided = "Aborting comment, no message provided"
	editorFailed      = "Failed to launch preferred editor"
	defaultMessage    = "\n# Enter comment content\n# Lines beginning with '#' will be stripped"
)

var buildVersion string
var message = goopt.String([]string{"-m", "--message"}, "", "comment message")
var amendID = goopt.String([]string{"--amend"}, "", "ID of a comment to amend.")
var deleteID = goopt.String([]string{"--delete"}, "", "ID of a comment to delete")
var printVersion = goopt.Flag([]string{"-v", "--version"}, []string{}, "Show the version number", "")
var remoteToConfig = goopt.String([]string{"--configure-remote"}, "", "remote to configure for fetching and pushing comments")

func main() {
	goopt.Description = func() string {
		return "Add comments to commits and files within git repositories"
	}
	goopt.Version = buildVersion
	goopt.Summary = "Annotate git commits"
	goopt.Parse(nil)
	pwd, err := os.Getwd()
	handleError(err)
	if len(goopt.Args) > 2 {
		handleInputError(errors.New(tooManyArguments))
	} else if *printVersion {
		fmt.Println(buildVersion)
	} else if len(*remoteToConfig) > 0 {
		handleError(gitc.ConfigureRemoteForComments(pwd, *remoteToConfig))
		fmt.Printf("Remote '%v' updated\n", *remoteToConfig)
	} else if len(*deleteID) > 0 {
		handleError(gitc.DeleteComment(pwd, *deleteID))
		fmt.Println("Comment deleted")
	} else {
		editComment(pwd)
	}
}

func handleInputError(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(goopt.Help())
		os.Exit(1)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func editComment(pwd string) {
	var commit *string = nil
	if len(*message) == 0 {
		*message = getMessageFromEditor(pwd)
	}
	var fileref = ""
	if len(goopt.Args) > 1 {
		fileref = goopt.Args[1]
	}
	if len(goopt.Args) > 0 {
		commit = &goopt.Args[0]
	}
	if len(*amendID) > 0 {
		id, err := gitc.UpdateComment(pwd, *amendID, *message)
		handleError(err)
		fmt.Printf("[%v] Comment updated\n", (*id)[:7])
	} else {
		id, err := gitc.CreateComment(pwd, commit, gitc.CreateFileRef(fileref), *message)
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
