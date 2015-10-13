package main

import (
	gx "exec"
	gg "git"
	"github.com/kylef/result.go/src/result"
	kp "gopkg.in/alecthomas/kingpin.v2"
	gc "libgitcomment"
	"math"
	"os"
)

const (
	defaultContextLines = 3
	linesBeforeConfig   = "comment.logBefore"
	linesAfterConfig    = "comment.logAfter"
)

var (
	buildVersion     string
	app              = kp.New("git-comment-log", "List git commit comments")
	fullDiff         = app.Flag("full-diff", "Show the full diff surrounding the comments").Bool()
	pretty           = app.Flag("pretty", "Pretty-print the comments in a format such as short, full, raw, or custom placeholders.").String()
	enablePager      = app.Flag("pager", "Use pager (Default)").Default("true").Bool()
	enableColor      = app.Flag("color", "Use color (Default)").Default("true").Bool()
	enableMarginLine = app.Flag("margin-line", "Use margin line (Default)").Default("true").Bool()
	lineNumbers      = app.Flag("line-numbers", "Show line numbers").Bool()
	linesBefore      = app.Flag("lines-before", "Number of context lines to show before comments").Short('B').Int64()
	linesAfter       = app.Flag("lines-after", "Number of context lines to show after comments").Short('A').Int64()
	revision         = app.Arg("revision range", "Filter comments to comments on commits from the specified range").String()
	contextLines     uint32
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	fatalIfError(app, gc.VersionCheck(pwd, buildVersion), "version")
	showComments(pwd)
}

func showComments(pwd string) {
	termHeight, termWidth := gx.CalculateDimensions()
	pager := gx.NewPager(app, pwd, gg.ConfiguredPager(pwd), termHeight, !*enablePager)
	computeContextLines(pwd)
	diff := gc.DiffCommits(pwd, *revision, contextLines)
	app.FatalIfError(diff.Failure, "diff")
	formatter := newFormatter(pwd, termWidth)
	printer := newPrinter(pager, formatter)
	printer.PrintDiff(diff.Success.(*gc.Diff))
}

func newFormatter(wd string, termWidth uint16) *Formatter {
	var useColor bool
	if *enableColor {
		useColor = gg.ConfiguredBool(wd, "color.pager", false)
	}
	return NewFormatter(*pretty, *lineNumbers, useColor, *enableMarginLine, termWidth)
}

func newPrinter(pager *gx.Pager, formatter *Formatter) *DiffPrinter {
	printer := NewDiffPrinter(pager, formatter, *linesBefore, *linesAfter)
	printer.PrintFullDiff = *fullDiff
	return printer
}

func computeContextLines(wd string) {
	if *linesBefore == 0 {
		before := int64(gg.ConfiguredInt32(wd, linesBeforeConfig, defaultContextLines))
		linesBefore = &before
	}
	if *linesAfter == 0 {
		after := int64(gg.ConfiguredInt32(wd, linesAfterConfig, defaultContextLines))
		linesAfter = &after
	}
	contextLines = uint32(math.Max(float64(*linesBefore), float64(*linesAfter)))
}

// Return the success value, otherwise kill the app with
// the error code specified
func fatalIfError(app *kp.Application, r result.Result, code string) interface{} {
	app.FatalIfError(r.Failure, code)
	return r.Success
}
