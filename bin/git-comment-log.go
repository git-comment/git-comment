package main

import (
	"fmt"
	gitc "git_comment"
	ex "git_comment/exec"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"os/exec"
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
	app.FatalIfError(err, "pwd")
	showComments(pwd)
}

func showComments(pwd string) {
	comments, err := gitc.CommentsOnCommit(pwd, revision)
	app.FatalIfError(err, "git")
	lineCount := ex.CalculateLineCount()
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
				cmd, writer, err = ex.ExecPager(pwd)
				app.FatalIfError(err, "pager")
			}
			if len(content) > 0 {
				_, err = writer.Write(content)
				content = []byte{}
			} else {
				_, err = writer.Write(formatted)
			}
			app.FatalIfError(err, "writer")
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
		return []byte(fmt.Sprintf("comment %v\n%v\n", *comment.ID, comment.Serialize()))
	}
	return []byte(comment.Serialize())
}
