package main

import (
	gitc "git_comment"
	gite "git_comment/exec"
	gitl "git_comment/log"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	buildVersion string
	app          = kp.New("git-comment-log", "List git commit comments")
	pretty       = app.Flag("pretty", "Pretty-print the comments in a format such as short, full, raw, or custom placeholders.").String()
	noPager      = app.Flag("nopager", "Disable pager").Bool()
	noColor      = app.Flag("nocolor", "Disable color").Bool()
	lineNumbers  = app.Flag("line-numbers", "Show line numbers").Bool()
	linesBefore  = app.Flag("lines-before", "Number of context lines to show before comments").Short('B').Int()
	linesAfter   = app.Flag("lines-after", "Number of context lines to show after comments").Short('A').Int()
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
	pager := gitl.NewContentPager(app, pwd, termHeight, *noPager)
	diff := gite.FatalIfError(app, gitc.DiffCommits(pwd, *revision), "diff")
	for _, file := range diff.(*gitc.Diff).Files {
		var printedFileHeader = false
		var afterComment = false
		beforeBuffer := make([]*gitc.DiffLine, 0)
		afterBuffer := make([]*gitc.DiffLine, 0)
		for _, line := range file.Lines {
			if len(line.Comments) > 0 {
				for _, line := range afterBuffer {
					pager.AddContent(formatter.FormatLine(line))
				}
				afterBuffer = make([]*gitc.DiffLine, 0)
				if !printedFileHeader {
					pager.AddContent(formatter.FormatFilePath(file))
					printedFileHeader = true
				}
				for _, line := range beforeBuffer {
					pager.AddContent(formatter.FormatLine(line))
				}
				pager.AddContent(formatter.FormatLine(line))
				for _, comment := range line.Comments {
					pager.AddContent(formatter.FormatComment(comment))
				}
				beforeBuffer = make([]*gitc.DiffLine, 0)
				afterComment = true
			} else {
				if afterComment {
					afterBuffer = append(afterBuffer, line)
					if len(afterBuffer) == *linesAfter {
						for _, line := range afterBuffer {
							pager.AddContent(formatter.FormatLine(line))
						}
						afterBuffer = make([]*gitc.DiffLine, 0)
						afterComment = false
					}
				} else {
					beforeBuffer = append(beforeBuffer, line)
					if len(beforeBuffer) > *linesBefore {
						beforeBuffer = append(beforeBuffer[:0], beforeBuffer[1:]...)
					}
				}
			}
		}
	}
	pager.Finish()
}
