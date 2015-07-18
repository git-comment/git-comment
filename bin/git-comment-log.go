package main

import (
	"fmt"
	gitc "git_comment"
	gite "git_comment/exec"
	gitl "git_comment/log"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	buildVersion string
	app          = kp.New("git-comment-log", "List git commit comments")
	pretty       = app.Flag("pretty", "Pretty-print the comments in a format such as short, full, raw, or custom placeholders.").String()
	noPager      = app.Flag("nopager", "Disable pager").Bool()
	noColor      = app.Flag("nocolor", "Disable color").Bool()
	lineNumbers  = app.Flag("line-numbers", "Show line numbers").Bool()
	revision     = app.Arg("revision range", "Filter comments to comments on commits from the specified range").String()
	formatter    *gitl.Formatter
	termHeight   uint16
	termWidth    uint16
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	configureFormatter()
	showComments(pwd)
}

func configureFormatter() {
	var useColor bool
	if !*noColor {
		if wd, err := os.Getwd(); err == nil {
			useColor = gitc.ConfiguredBool(wd, "color.pager", false)
		}
	}
	termHeight, termWidth = gite.CalculateDimensions()
	formatter = gitl.NewFormatter(*pretty, *lineNumbers, useColor, termWidth)
}

func showComments(pwd string) {
	var usePager bool = termHeight == 0 && !*noPager
	var content []byte
	var writer io.WriteCloser
	var cmd *exec.Cmd
	var err error
	diff := gite.FatalIfError(app, gitc.DiffCommits(pwd, *revision), "diff")
	pageContent := func(data string) {
		content = append(content, []byte(data)...)
		if !usePager {
			lines := strings.Split(string(content), "\n")
			usePager = !*noPager && uint16(len(lines)) > termHeight-1
		}
		if usePager {
			if writer == nil {
				cmd, writer, err = gite.ExecPager(pwd)
				app.FatalIfError(err, "pager")
			}
			if len(content) > 0 {
				_, err = writer.Write(content)
				content = []byte{}
				app.FatalIfError(err, "writer")
			}
		}
	}
	for _, file := range diff.(*gitc.Diff).Files {
		var printedFileHeader = false
		var afterComment = false
		beforeBuffer := make([]*gitc.DiffLine, 0)
		afterBuffer := make([]*gitc.DiffLine, 0)
		for _, line := range file.Lines {
			if len(line.Comments) > 0 {
				for _, line := range afterBuffer {
					pageContent(formatter.FormatLine(line))
				}
				afterBuffer = make([]*gitc.DiffLine, 0)
				if !printedFileHeader {
					pageContent(formatter.FormatFilePath(file))
					printedFileHeader = true
				}
				for _, line := range beforeBuffer {
					pageContent(formatter.FormatLine(line))
				}
				pageContent(formatter.FormatLine(line))
				for _, comment := range line.Comments {
					pageContent(formatter.FormatComment(comment))
				}
				beforeBuffer = make([]*gitc.DiffLine, 0)
				afterComment = true
			} else {
				if afterComment {
					afterBuffer = append(afterBuffer, line)
					if len(afterBuffer) == 5 {
						for _, line := range afterBuffer {
							pageContent(formatter.FormatLine(line))
						}
						afterBuffer = make([]*gitc.DiffLine, 0)
						afterComment = false
					}
				} else {
					beforeBuffer = append(beforeBuffer, line)
					if len(beforeBuffer) > 5 {
						beforeBuffer = append(beforeBuffer[:0], beforeBuffer[1:]...)
					}
				}
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
