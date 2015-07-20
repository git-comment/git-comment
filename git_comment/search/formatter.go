package search

import (
	"fmt"
	gitc "git_comment"
	ex "git_comment/exec"
	"strings"
)

type Formatter struct {
	useColor bool
}

func NewFormatter(useColor bool) *Formatter {
	return &Formatter{useColor}
}

func (f *Formatter) FormatComment(c *gitc.Comment) string {
	return fmt.Sprintf("%v  %v\n",
		ex.Colorize(ex.Cyan, f.formatHeader(c), f.useColor),
		f.formatTitle(c))
}

func (f *Formatter) formatHeader(c *gitc.Comment) string {
	var path string
	if c.FileRef != nil {
		path = c.FileRef.Serialize()
	}
	name := c.Author.Name
	return fmt.Sprintf("%v:%v:%v:%v\n",
		(*c.ID)[:7],
		(*c.Commit)[:7],
		name,
		path)
}

func (f *Formatter) formatTitle(c *gitc.Comment) string {
	lines := strings.Split(c.Content, "\n")
	title := lines[0]
	if len(lines) > 1 {
		title = fmt.Sprintf("%v...", title)
	}
	return title
}
