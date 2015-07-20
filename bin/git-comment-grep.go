package main

import (
	"fmt"
	gc "git_comment"
	gx "git_comment/exec"
	gl "git_comment/log"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

var (
	buildVersion string
	app          = kp.New("git-comment-grep", "Index and look for comments")
	findCmd      = app.Command("find", "Look for comments containing text")
	indexCmd     = app.Command("index", "Index and cache comment content")
	noPager      = app.Flag("nopager", "Disable pager").Bool()
	noColor      = app.Flag("nocolor", "Disable color").Bool()
	text         = findCmd.Arg("text", "Search text").String()
)

func main() {
	app.Version(buildVersion)
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	switch kp.MustParse(app.Parse(os.Args[1:])) {
	case "find":
		findText(pwd, *text)
	case "index":
		indexComments(pwd)
	}
}

func findText(wd, text string) {
	termHeight, _ := gx.CalculateDimensions()
	matches := gx.FatalIfError(app, gc.CommentsWithContent(wd, text), "find")
	pager := gl.NewPager(app, wd, termHeight, *noPager)
	for _, index := range matches.([]*gc.Comment) {
		pager.AddContent(formatComment(index))
	}
	pager.Finish()
}

func indexComments(wd string) {
	fmt.Printf("Indexing...")
	gx.FatalIfError(app, gc.IndexComments(wd), "index")
	fmt.Printf("done\n")
}

func formatComment(c *gc.Comment) string {
	var path string
	if c.FileRef != nil {
		path = c.FileRef.Path
	}
	name := c.Author.Name
	lines := strings.Split(c.Content, "\n")
	title := lines[0]
	if len(lines) > 1 {
		title = fmt.Sprintf("%v...", title)
	}

	return fmt.Sprintf("[%v]%v:%v:%v\n  %v\n",
		(*c.ID)[:7],
		(*c.Commit)[:7],
		name,
		path,
		title)
}
