package search

import (
	"fmt"
	gc "libgitcomment"
	gx "libgitcomment/exec"
	"path/filepath"
	"strings"
)

type Formatter struct {
	useColor bool
}

func NewFormatter(useColor bool) *Formatter {
	return &Formatter{useColor}
}

func (f *Formatter) FormatComment(c *gc.Comment, highlight string) string {
	return fmt.Sprintf("%v  %v\n",
		gx.Colorize(gx.Cyan, f.formatHeader(c), f.useColor),
		f.formatTitle(c, highlight))
}

func (f *Formatter) formatHeader(c *gc.Comment) string {
	var path string
	if c.FileRef != nil {
		_, path = filepath.Split(c.FileRef.Serialize())
	}
	name := c.Author.Name
	return fmt.Sprintf("%v %v %v:%v\n",
		name,
		c.Author.Date.Format("2006-01-02"),
		(*c.Commit)[:7],
		path)
}

func (f *Formatter) formatTitle(c *gc.Comment, highlight string) string {
	lines := strings.Split(c.Content, "\n")
	title := lines[0]
	if len(lines) > 1 {
		title = fmt.Sprintf("%v...", title)
	}
	if f.useColor {
		title = strings.Replace(title, highlight, gx.Colorize(gx.Red, highlight, true), -1)
	}
	return title
}
