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
	"time"
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
	var usePager bool = termHeight == 0
	var content []byte
	var writer io.WriteCloser
	var cmd *exec.Cmd
	diff, err := gitc.DiffCommits(pwd, *revision)
	app.FatalIfError(err, "diff")
	pageContent := func(data string) {
		content = append(content, []byte(data)...)
		if !usePager {
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
				app.FatalIfError(err, "writer")
			}
		}
	}
	for _, file := range diff.Files {
		pageContent(formattedFilePath(file))
		for _, line := range file.Lines {
			pageContent(formattedLine(line))
			for _, comment := range line.Comments {
				pageContent(formattedComment(comment))
			}
		}
	}
	if !usePager {
		fmt.Println(string(content))
	}

	if writer != nil {
		writer.Close()
		cmd.Wait()
	}
}

func formattedLineNumber(number int) string {
	const lineNumberMax = 5
	var line string
	if number < 0 {
		line = " "
	} else {
		line = fmt.Sprintf("%d", number)
	}
	for len(line) < lineNumberMax {
		line = fmt.Sprintf(" %v", line)
	}
	return line
}

func formattedLineNumbers(oldNum, newNum int) string {
	var newLine string
	oldLine := formattedLineNumber(oldNum)
	if oldNum == newNum {
		newLine = formattedLineNumber(-1)
	} else {
		newLine = formattedLineNumber(newNum)
	}
	return fmt.Sprintf("%v %v", oldLine, newLine)
}

func formattedFilePath(file *gitc.DiffFile) string {
	var path string
	if file.OldPath == file.NewPath {
		path = file.OldPath
	} else if len(file.OldPath) > 0 && len(file.NewPath) > 0 {
		path = fmt.Sprintf("%v -> %v", file.OldPath, file.NewPath)
	} else if len(file.OldPath) > 0 {
		path = file.OldPath
	} else {
		path = file.NewPath
	}
	return fmt.Sprintf("%v\n", path)
}

func formattedLinePrefix(line *gitc.DiffLine) string {
	switch line.Type {
	case gitc.DiffAdd, gitc.DiffAddNewline:
		return "+"
	case gitc.DiffRemove, gitc.DiffRemoveNewline:
		return "-"
	default:
		return " "
	}
}

func formattedLineContent(line *gitc.DiffLine) string {
	switch line.Type {
	case gitc.DiffAddNewline, gitc.DiffRemoveNewline:
		return "↵"
	default:
		return line.Content
	}
}

func formattedLine(line *gitc.DiffLine) string {
	prefix := formattedLinePrefix(line)
	number := formattedLineNumbers(line.OldLineNumber, line.NewLineNumber)
	content := formattedLineContent(line)
	return fmt.Sprintf("%v %v %v", prefix, number, content)
}

func formattedComment(comment *gitc.Comment) string {
	if *pretty == Short || len(*pretty) == 0 {
		return substituteVariables(ShortFormat, comment)
	} else if *pretty == Full {
		return substituteVariables(FullFormat, comment)
	} else if *pretty == Raw {
		format := string(substituteVariables(RawFormat, comment))
		return fmt.Sprintf(format, *comment.ID)
	} else if strings.HasPrefix(*pretty, formatPrefix) {
		return substituteVariables((*pretty)[len(formatPrefix):], comment)
	}
	app.FatalUsage(invalidFormat)
	return ""
}

func substituteVariables(format string, comment *gitc.Comment) string {
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
	format = strings.Replace(format, authorDateISO8601, comment.CreateTime.Format(time.RFC3339), -1)
	format = strings.Replace(format, authorDateUnix, fmt.Sprintf("%v", comment.CreateTime.Unix()), -1)
	format = strings.Replace(format, committerName, comment.Amender.Name, -1)
	format = strings.Replace(format, committerEmail, comment.Amender.Email, -1)
	format = strings.Replace(format, committerDateISO8601, comment.AmendTime.Format(time.RFC3339), -1)
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
	return format
}
