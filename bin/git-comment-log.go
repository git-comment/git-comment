package main

import (
	gitc "git_comment"
	gite "git_comment/exec"
	gitl "git_comment/log"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"math"
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
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	showComments(pwd)
}

func showComments(pwd string) {
	termHeight, termWidth := gite.CalculateDimensions()
	pager := gitl.NewPager(app, pwd, termHeight, *noPager)
	contextLines := uint32(math.Max(float64(*linesBefore), float64(*linesAfter)))
	diff := gitc.DiffCommits(pwd, *revision, contextLines)
	app.FatalIfError(diff.Failure, "diff")
	formatter := newFormatter(termWidth)
	printer := gitl.NewDiffPrinter(pager, formatter, *linesBefore, *linesAfter)
	printer.PrintDiff(diff.Success.(*gitc.Diff))
}

func newFormatter(termWidth uint16) *gitl.Formatter {
	var useColor bool
	if !*noColor {
		if wd, err := os.Getwd(); err == nil {
			useColor = gitc.ConfiguredBool(wd, "color.pager", false)
		}
	}
	return gitl.NewFormatter(*pretty, *lineNumbers, useColor, termWidth)
}
