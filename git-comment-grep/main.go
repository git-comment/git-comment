package main

import (
	gx "exec"
	"fmt"
	gg "git"
	"github.com/kylef/result.go/src/result"
	kp "gopkg.in/alecthomas/kingpin.v2"
	gc "libgitcomment"
	"os"
)

var (
	buildVersion string
	app          = kp.New("git-comment-grep", "Index and look for comments")
	findCmd      = app.Command("find", "Look for comments containing text")
	indexCmd     = app.Command("index", "Index and cache comment content")
	noPager      = app.Flag("nopager", "Disable pager").Bool()
	noColor      = app.Flag("nocolor", "Disable color").Bool()
	text         = findCmd.Arg("text", "Search text").Required().String()
)

func main() {
	app.Version(buildVersion)
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	fatalIfError(app, gc.VersionCheck(pwd, buildVersion), "version")
	switch kp.MustParse(app.Parse(os.Args[1:])) {
	case "find":
		findText(pwd, *text)
	case "index":
		indexComments(pwd)
	}
}

func findText(wd, text string) {
	termHeight, _ := gx.CalculateDimensions()
	var useColor bool
	if !*noColor {
		useColor = gg.ConfiguredBool(wd, "color.pager", false)
	}
	pager := gx.NewPager(app, wd, gg.ConfiguredPager(wd), termHeight, *noPager)
	printer := NewPrinter(useColor, pager)
	fatalIfError(app, printer.PrintCommentsMatching(wd, text), "find")
}

func indexComments(wd string) {
	fmt.Printf("Indexing...")
	fatalIfError(app, IndexComments(wd), "index")
	fmt.Printf("done\n")
}

// Return the success value, otherwise kill the app with
// the error code specified
func fatalIfError(app *kp.Application, r result.Result, code string) interface{} {
	app.FatalIfError(r.Failure, code)
	return r.Success
}
