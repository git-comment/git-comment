package main

import (
	gitc "git_comment"
	gite "git_comment/exec"
	gitg "git_comment/git"
	gitl "git_comment/log"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"math"
	"os"
)

const defaultContextLines = 3

var (
	buildVersion string
	app          = kp.New("git-comment-log", "List git commit comments")
	fullDiff     = app.Flag("full-diff", "Show the full diff surrounding the comments").Bool()
	pretty       = app.Flag("pretty", "Pretty-print the comments in a format such as short, full, raw, or custom placeholders.").String()
	noPager      = app.Flag("nopager", "Disable pager").Bool()
	noColor      = app.Flag("nocolor", "Disable color").Bool()
	lineNumbers  = app.Flag("line-numbers", "Show line numbers").Bool()
	linesBefore  = app.Flag("lines-before", "Number of context lines to show before comments").Short('B').Int64()
	linesAfter   = app.Flag("lines-after", "Number of context lines to show after comments").Short('A').Int64()
	revision     = app.Arg("revision range", "Filter comments to comments on commits from the specified range").String()
	contextLines uint32
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	gite.FatalIfError(app, gitc.VersionCheck(pwd, buildVersion), "version")
	showComments(pwd)
}

func showComments(pwd string) {
	termHeight, termWidth := gite.CalculateDimensions()
	pager := gite.NewPager(app, pwd, termHeight, *noPager)
	computeContextLines(pwd)
	diff := gitc.DiffCommits(pwd, *revision, contextLines)
	app.FatalIfError(diff.Failure, "diff")
	formatter := newFormatter(pwd, termWidth)
	printer := newPrinter(pager, formatter)
	printer.PrintDiff(diff.Success.(*gitc.Diff))
}

func newFormatter(wd string, termWidth uint16) *gitl.Formatter {
	var useColor bool
	if !*noColor {
		useColor = gitg.ConfiguredBool(wd, "color.pager", false)
	}
	return gitl.NewFormatter(*pretty, *lineNumbers, useColor, termWidth)
}

func newPrinter(pager *gite.Pager, formatter *gitl.Formatter) *gitl.DiffPrinter {
	printer := gitl.NewDiffPrinter(pager, formatter, *linesBefore, *linesAfter)
	printer.PrintFullDiff = *fullDiff
	return printer
}

func computeContextLines(wd string) {
	if *linesBefore == 0 {
		before := int64(gitg.ConfiguredInt32(wd, "comment-log.lines-before", defaultContextLines))
		linesBefore = &before
	}
	if *linesAfter == 0 {
		after := int64(gitg.ConfiguredInt32(wd, "comment-log.lines-after", defaultContextLines))
		linesAfter = &after
	}
	contextLines = uint32(math.Max(float64(*linesBefore), float64(*linesAfter)))
}
