package log

import (
	gx "github.com/git-comment/git-comment/exec"

	gc "github.com/git-comment/git-comment"
)

const FullDiffContext = -1

type DiffPrinter struct {
	PrintFullDiff     bool
	pager             *gx.Pager
	formatter         *Formatter
	beforeBuffer      []*gc.DiffLine
	afterBuffer       []*gc.DiffLine
	beforeBufferSize  int64
	afterBufferSize   int64
	afterComment      bool
	printedFileHeader bool
	currentFile       *gc.DiffFile
}

func NewDiffPrinter(pager *gx.Pager, formatter *Formatter, linesBefore int64, linesAfter int64) *DiffPrinter {
	printer := &DiffPrinter{}
	printer.pager = pager
	printer.formatter = formatter
	printer.afterBufferSize = linesAfter
	printer.beforeBufferSize = linesBefore
	return printer
}

func (r *DiffPrinter) PrintDiff(diff *gc.Diff) {
	r.currentFile = nil
	for _, file := range diff.Files {
		r.currentFile = file
		r.printedFileHeader = false
		r.afterComment = false
		r.beforeBuffer = make([]*gc.DiffLine, 0)
		r.afterBuffer = make([]*gc.DiffLine, 0)
		for _, line := range file.Lines {
			if len(line.Comments) > 0 {
				r.printLineWithContext(line)
			} else if r.PrintFullDiff {
				r.printLine(line)
			} else if r.afterComment {
				r.addLineAfterComments(line)
			} else {
				r.addLineBeforeComments(line)
			}
		}
	}
	r.printTrailingLines()
	r.pager.Finish()
}

func (r *DiffPrinter) addLineBeforeComments(line *gc.DiffLine) {
	r.beforeBuffer = append(r.beforeBuffer, line)
	if int64(len(r.beforeBuffer)) > r.beforeBufferSize {
		r.beforeBuffer = append(r.beforeBuffer[:0], r.beforeBuffer[1:]...)
	}
}

func (r *DiffPrinter) addLineAfterComments(line *gc.DiffLine) {
	r.afterBuffer = append(r.afterBuffer, line)
	if int64(len(r.afterBuffer)) == r.afterBufferSize {
		r.printTrailingLines()
		r.afterComment = false
	}
}

func (r *DiffPrinter) printLineWithContext(line *gc.DiffLine) {
	r.printTrailingLines()
	r.printLeadingLines()
	r.printLine(line)
	for _, comment := range line.Comments {
		r.pager.AddContent(r.formatter.FormatComment(comment))
	}
	r.afterComment = true
}

func (r *DiffPrinter) printTrailingLines() {
	if r.printLines(r.afterBuffer) {
		r.afterBuffer = make([]*gc.DiffLine, 0)
	}
}

func (r *DiffPrinter) printLeadingLines() {
	if r.printLines(r.beforeBuffer) {
		r.beforeBuffer = make([]*gc.DiffLine, 0)
	}
}

// return true if any lines added
func (r *DiffPrinter) printLines(lines []*gc.DiffLine) bool {
	for _, line := range lines {
		r.printLine(line)
	}
	return len(lines) > 0
}

func (r *DiffPrinter) printLine(line *gc.DiffLine) {
	if !r.printedFileHeader {
		r.pager.AddContent(r.formatter.FormatFilePath(r.currentFile))
		r.printedFileHeader = true
	}
	r.pager.AddContent(r.formatter.FormatLine(line))
}
