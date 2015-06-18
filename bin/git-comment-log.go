package main

import (
	// "errors"
	"fmt"
	gitc "git_comment"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	errorPrefix = "git-comment-log"
)

var (
	buildVersion string
	app          = kp.New("git-comment-log", "List git commit comments")
	revision     = app.Arg("revision range", "Filter comments to comments on commits from the specified range").String()
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, errorPrefix)
	comments, err := gitc.CommentsOnCommit(pwd, revision)
	app.FatalIfError(err, errorPrefix)
	var lineCount int = 0
	if lineEnv := getEnv("LINES"); lineEnv != nil {
		lines, _ := strconv.ParseInt(*lineEnv, 10, 0)
		lineCount = int(lines)
	}
	var usePager bool = lineCount == 0
	var content []byte
	var writer io.WriteCloser
	var cmd *exec.Cmd
	for i := 0; i < len(comments); i++ {
		comment := comments[i]
		formatted := formattedContent(comment)
		if !usePager {
			content = append(content, formatted...)
			lines := strings.Split(string(content), "\n")
			usePager = len(lines) > lineCount-1
		}
		if usePager {
			if writer == nil {
				cmd, writer, err = execPager(pwd)
				app.FatalIfError(err, errorPrefix)
			}
			_, err = writer.Write(formatted)
			app.FatalIfError(err, errorPrefix)
		}
	}
	if writer != nil {
		writer.Close()
		cmd.Wait()
	} else {
		fmt.Println(string(content))
	}
}

func formattedContent(comment *gitc.Comment) []byte {
	if comment.ID != nil && len(*comment.ID) > 0 {
		return []byte(fmt.Sprintf("comment %v\n%v", *comment.ID, comment.Serialize()))
	}
	return []byte(comment.Serialize())
}

func execPager(pwd string) (*exec.Cmd, io.WriteCloser, error) {
	pager := gitc.ConfiguredPager(pwd)
	cmd := exec.Command(*pager)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return cmd, pipe, nil
}

func getEnv(name string) *string {
	if env := os.Getenv(name); len(env) > 0 {
		return &env
	}
	return nil
}
