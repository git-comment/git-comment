package main

import (
	gx "../exec"
	gc "../libgitcomment"
	gg "../libgitcomment/git"
	"fmt"
	kp "gopkg.in/alecthomas/kingpin.v2"
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
	gx.FatalIfError(app, gc.VersionCheck(pwd, buildVersion), "version")
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
	printer := gs.NewPrinter(useColor, pager)
	gx.FatalIfError(app, printer.PrintCommentsMatching(wd, text), "find")
}

func indexComments(wd string) {
	fmt.Printf("Indexing...")
	gx.FatalIfError(app, gs.IndexComments(wd), "index")
	fmt.Printf("done\n")
}
