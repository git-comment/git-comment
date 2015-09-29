package search

import (
	gc "github.com/git-comment/git-comment"
	gx "github.com/git-comment/git-comment/exec"
	"github.com/kylef/result.go/src/result"
)

type Printer struct {
	formatter *Formatter
	pager     *gx.Pager
}

func NewPrinter(useColor bool, pager *gx.Pager) *Printer {
	return &Printer{&Formatter{useColor}, pager}
}

func (p *Printer) PrintCommentsMatching(wd, text string) result.Result {
	return CommentsWithContent(wd, text).FlatMap(func(matches interface{}) result.Result {
		for _, comment := range matches.([]*gc.Comment) {
			p.pager.AddContent(p.formatter.FormatComment(comment, text))
		}
		p.pager.Finish()
		return result.NewSuccess(true)
	})
}
