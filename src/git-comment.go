package main

import (
	"errors"
	"fmt"
	gitc "git_comment"
	goopt "github.com/droundy/goopt"
	"os"
)

var buildVersion string
var message = goopt.String([]string{"-m", "--message"}, "", "comment message")
var amendID = goopt.String([]string{"--amend"}, "", "ID of a comment to amend. `--message` is required")
var deleteID = goopt.String([]string{"--delete"}, "", "ID of a comment to delete")
var printVersion = goopt.Flag([]string{"-v", "--version"}, []string{}, "Show the version number", "")

func main() {
	goopt.Parse(nil)
	var err error
	pwd, osErr := os.Getwd()
	handleError(osErr)
	if len(goopt.Args) > 2 {
		handleError(errors.New("Too many arguments provided"))
	} else if *printVersion {
		fmt.Println(buildVersion)
	} else if len(*deleteID) > 0 {
		err = gitc.DeleteComment(pwd, *deleteID)
		handleError(err)
	} else if len(*message) > 0 {
		var commit = ""
		var fileref = ""
		if len(goopt.Args) > 1 {
			fileref = goopt.Args[1]
		}
		if len(goopt.Args) > 0 {
			commit = goopt.Args[0]
		}
		if len(*amendID) > 0 {
			err = gitc.UpdateComment(pwd, *amendID, *message)
		} else {
			_, err = gitc.CreateComment(pwd, commit, gitc.CreateFileRef(fileref), *message)
		}
		handleError(err)
	} else {
		fmt.Println(goopt.Help())
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(goopt.Help())
		os.Exit(1)
	}
}
