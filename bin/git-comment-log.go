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
	commentFull          = "%C"
	commentShort         = "%c"
	commitFull           = "%H"
	commitShort          = "%h"
	filePath             = "%f"
	lineNumber           = "%L"
	authorName           = "%an"
	authorEmail          = "%ae"
	authorDateISO8601    = "%ad"
	authorDateUnix       = "%aU"
	committerName        = "%cn"
	committerEmail       = "%ce"
	committerDateISO8601 = "%cd"
	committerDateUnix    = "%cU"
	bodyContent          = "%b"
	titleLine            = "%t"
	newLine              = "%n"
	dividerLine          = "%d"
)

const (
	Short         = "short"
	Full          = "full"
	Raw           = "raw"
	ShortFormat   = "[%h] %c %an <%ae>\n%t\n\n"
	FullFormat    = "commit  %H\ncomment %C\nAuthor: %an <%ae>\n%b\n\n"
	RawFormat     = "comment %C\n%v\n\n"
	ISO8601       = "2015-06-21T18:24:18Z"
	formatPrefix  = "format:"
	invalidFormat = "Unknown pretty format."
)

var (
	buildVersion string
	termWidth    uint16
	termHeight   uint16
	app          = kp.New("git-comment-log", "List git commit comments")
	pretty       = app.Flag("pretty", "Pretty-print the comments in a format such as short, full, raw, or custom placeholders.").String()
	revision     = app.Arg("revision range", "Filter comments to comments on commits from the specified range").String()
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	termHeight, termWidth = ex.CalculateDimensions()
	showComments(pwd)
}

func showComments(pwd string) {
	comments, err := gitc.CommentsOnCommit(pwd, revision)
	app.FatalIfError(err, "git")
	var usePager bool = termHeight == 0
	var content []byte
	var writer io.WriteCloser
	var cmd *exec.Cmd
	for i := 0; i < len(comments); i++ {
		comment := comments[i]
		formatted := formattedContent(comment)
		if !usePager {
			content = append(content, formatted...)
			lines := strings.Split(string(content), "\n")
			usePager = len(lines) > int(termHeight-1)
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
	if *pretty == Short || len(*pretty) == 0 {
		return substituteVariables(ShortFormat, comment)
	} else if *pretty == Full {
		return substituteVariables(FullFormat, comment)
	} else if *pretty == Raw {
		format := string(substituteVariables(RawFormat, comment))
		return []byte(fmt.Sprintf(format, *comment.ID))
	} else if strings.HasPrefix(*pretty, formatPrefix) {
		return substituteVariables((*pretty)[len(formatPrefix):], comment)
	}
	app.FatalUsage(invalidFormat)
	return []byte{}
}

func substituteVariables(format string, comment *gitc.Comment) []byte {
	var path = ""
	var line = ""
	if len(comment.FileRef.Path) > 0 {
		path = comment.FileRef.Path
	}
	if comment.FileRef.Line > 0 {
		line = fmt.Sprintf("%v", comment.FileRef.Line)
	}
	format = strings.Replace(format, authorName, comment.Author.Name, -1)
	format = strings.Replace(format, authorEmail, comment.Author.Email, -1)
	format = strings.Replace(format, authorDateISO8601, comment.CreateTime.Format(ISO8601), -1)
	format = strings.Replace(format, authorDateUnix, fmt.Sprintf("%v", comment.CreateTime.Unix()), -1)
	format = strings.Replace(format, committerName, comment.Amender.Name, -1)
	format = strings.Replace(format, committerEmail, comment.Amender.Email, -1)
	format = strings.Replace(format, committerDateISO8601, comment.AmendTime.Format(ISO8601), -1)
	format = strings.Replace(format, committerDateUnix, fmt.Sprintf("%v", comment.AmendTime.Unix()), -1)
	format = strings.Replace(format, commentFull, *comment.ID, -1)
	format = strings.Replace(format, commentShort, (*comment.ID)[:7], -1)
	format = strings.Replace(format, commitFull, *comment.Commit, -1)
	format = strings.Replace(format, commitShort, (*comment.Commit)[:7], -1)
	format = strings.Replace(format, bodyContent, comment.Content, -1)
	format = strings.Replace(format, titleLine, strings.Split(comment.Content, "\n")[0], -1)
	format = strings.Replace(format, filePath, path, -1)
	format = strings.Replace(format, lineNumber, line, -1)
	format = strings.Replace(format, newLine, "\n", -1)
	format = strings.Replace(format, dividerLine, strings.Repeat("-", int(termWidth)), -1)
	return []byte(format)
}
