package log

import (
	gitc "git_comment"
)

type DiffPrinter struct {
	pager            *Pager
	formatter        *Formatter
	beforeBuffer     []*gitc.DiffLine
	afterBuffer      []*gitc.DiffLine
	beforeBufferSize int
	afterBufferSize  int
	afterComment     bool
}

func NewDiffPrinter(pager *Pager, formatter *Formatter, linesBefore int, linesAfter int) *DiffPrinter {
	printer := &DiffPrinter{}
	printer.pager = pager
	printer.formatter = formatter
	printer.afterBufferSize = linesAfter
	printer.beforeBufferSize = linesBefore
	return printer
}

func (r *DiffPrinter) addLineBeforeComments(line *gitc.DiffLine) {
	r.beforeBuffer = append(r.beforeBuffer, line)
	if len(r.beforeBuffer) > r.beforeBufferSize {
		r.beforeBuffer = append(r.beforeBuffer[:0], r.beforeBuffer[1:]...)
	}
}

func (r *DiffPrinter) addLineAfterComments(line *gitc.DiffLine) {
	r.afterBuffer = append(r.afterBuffer, line)
	if len(r.afterBuffer) == r.afterBufferSize {
		r.printTrailingLines()
		r.afterComment = false
	}
}

func (r *DiffPrinter) PrintDiff(diff *gitc.Diff) {
	for _, file := range diff.Files {
		var printedFileHeader = false
		r.afterComment = false
		r.beforeBuffer = make([]*gitc.DiffLine, 0)
		r.afterBuffer = make([]*gitc.DiffLine, 0)
		for _, line := range file.Lines {
			if len(line.Comments) > 0 {
				r.printTrailingLines()
				if !printedFileHeader {
					r.pager.AddContent(r.formatter.FormatFilePath(file))
					printedFileHeader = true
				}
				r.printLeadingLines()
				r.pager.AddContent(r.formatter.FormatLine(line))
				for _, comment := range line.Comments {
					r.pager.AddContent(r.formatter.FormatComment(comment))
				}
				r.afterComment = true
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

func (r *DiffPrinter) printTrailingLines() {
	if r.pageLines(r.afterBuffer) {
		r.afterBuffer = make([]*gitc.DiffLine, 0)
	}
}

func (r *DiffPrinter) printLeadingLines() {
	if r.pageLines(r.beforeBuffer) {
		r.beforeBuffer = make([]*gitc.DiffLine, 0)
	}
}

// return true if any lines added
func (r *DiffPrinter) pageLines(lines []*gitc.DiffLine) bool {
	for _, line := range lines {
		r.pager.AddContent(r.formatter.FormatLine(line))
	}
	return len(lines) > 0
}
