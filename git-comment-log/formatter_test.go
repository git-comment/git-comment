package main

import (
	"github.com/stvp/assert"
	gc "libgitcomment"
	"regexp"
	"testing"
	"time"
)

func TestPrettyFormatFilePath(t *testing.T) {
	formatter := NewFormatter("format:%f", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "src/file.c\n\n\n")
}

func TestPrettyFormatFileLine(t *testing.T) {
	formatter := NewFormatter("format:%L", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "12\n\n\n")
}

func TestPrettyFormatAuthorName(t *testing.T) {
	formatter := NewFormatter("format:%an", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "Simon\n\n\n")
}

func TestPrettyFormatExtraText(t *testing.T) {
	formatter := NewFormatter("format:hi, %an!", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "hi, Simon!\n\n\n")
}

func TestPrettyFormatAuthorEmail(t *testing.T) {
	formatter := NewFormatter("format:%ae", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "iceking@example.com\n\n\n")
}

func TestPrettyFormatAuthorDateUnix(t *testing.T) {
	formatter := NewFormatter("format:%aU", false, false, false, 0)
	dateRe := regexp.MustCompile(`^([0-9]{10})\s{3}$`)
	match := dateRe.FindStringSubmatch(formatter.FormatComment(comment()))
	assert.Equal(t, len(match), 2)
}

func TestPrettyFormatAuthorDateISO(t *testing.T) {
	formatter := NewFormatter("format:%ad", false, false, false, 0)
	dateRe := regexp.MustCompile(`^([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2})`)
	text := formatter.FormatComment(comment())
	match := dateRe.FindStringSubmatch(text)
	assert.Equal(t, len(match), 2, "Date not in correct format: ", text)
}

func TestPrettyFormatCommitterName(t *testing.T) {
	formatter := NewFormatter("format:%kn", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "Simon\n\n\n")
}

func TestPrettyFormatCommitterEmail(t *testing.T) {
	formatter := NewFormatter("format:%ke", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "iceking@example.com\n\n\n")
}

func TestPrettyFormatCommitterDateUnix(t *testing.T) {
	formatter := NewFormatter("format:%kU", false, false, false, 0)
	dateRe := regexp.MustCompile(`([0-9]{10})`)
	match := dateRe.FindStringSubmatch(formatter.FormatComment(comment()))
	assert.Equal(t, len(match), 2)
}

func TestPrettyFormatCommitterDateISO(t *testing.T) {
	formatter := NewFormatter("format:%kd", false, false, false, 0)
	dateRe := regexp.MustCompile(`([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2})`)
	text := formatter.FormatComment(comment())
	match := dateRe.FindStringSubmatch(text)
	assert.Equal(t, len(match), 2, "Date not in correct format: ", text)
}

func TestPrettyFormatBodyNoMarginLine(t *testing.T) {
	formatter := NewFormatter("format:%b", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "new comment\nmore context\n\n\n")
}

func TestPrettyFormatBodyWithMarginLine(t *testing.T) {
	formatter := NewFormatter("format:%b", false, false, true, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\n  │new comment\n  │more context\n\n")
}

func TestPrettyFormatTitle(t *testing.T) {
	formatter := NewFormatter("format:%t", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "new comment\n\n\n")
}

func TestPrettyFormatID(t *testing.T) {
	formatter := NewFormatter("format:%C", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "abcabcabcabc\n\n\n")
}

func TestPrettyFormatIDShort(t *testing.T) {
	formatter := NewFormatter("format:%c", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "abcabca\n\n\n")
}

func TestPrettyFormatCommitID(t *testing.T) {
	formatter := NewFormatter("format:%H", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "123444abcabc\n\n\n")
}

func TestPrettyFormatCommitIDShort(t *testing.T) {
	formatter := NewFormatter("format:%h", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "123444a\n\n\n")
}

func TestPrettyFormatBlackColor(t *testing.T) {
	formatter := NewFormatter("format:black(yes)", false, true, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[30myes\x1b[0m\n\n\n")
}

func TestPrettyFormatBlackNoColor(t *testing.T) {
	formatter := NewFormatter("format:black(yes)", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes\n\n\n")
}

func TestPrettyFormatRedColor(t *testing.T) {
	formatter := NewFormatter("format:red(yes)", false, true, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[31myes\x1b[0m\n\n\n")
}

func TestPrettyFormatRedNoColor(t *testing.T) {
	formatter := NewFormatter("format:red(yes)", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes\n\n\n")
}

func TestPrettyFormatGreenColor(t *testing.T) {
	formatter := NewFormatter("format:green(yes)", false, true, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[32myes\x1b[0m\n\n\n")
}

func TestPrettyFormatGreenNoColor(t *testing.T) {
	formatter := NewFormatter("format:green(yes)", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes\n\n\n")
}

func TestPrettyFormatYellowColor(t *testing.T) {
	formatter := NewFormatter("format:yellow(yes)", false, true, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[33myes\x1b[0m\n\n\n")
}

func TestPrettyFormatYellowNoColor(t *testing.T) {
	formatter := NewFormatter("format:yellow(yes)", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes\n\n\n")
}

func TestPrettyFormatBlueColor(t *testing.T) {
	formatter := NewFormatter("format:blue(yes)", false, true, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[34myes\x1b[0m\n\n\n")
}

func TestPrettyFormatBlueNoColor(t *testing.T) {
	formatter := NewFormatter("format:blue(yes)", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes\n\n\n")
}

func TestPrettyFormatMagentaColor(t *testing.T) {
	formatter := NewFormatter("format:magenta(yes)", false, true, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[35myes\x1b[0m\n\n\n")
}

func TestPrettyFormatMagentaNoColor(t *testing.T) {
	formatter := NewFormatter("format:magenta(yes)", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes\n\n\n")
}

func TestPrettyFormatCyanColor(t *testing.T) {
	formatter := NewFormatter("format:cyan(yes)", false, true, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[36myes\x1b[0m\n\n\n")
}

func TestPrettyFormatCyanNoColor(t *testing.T) {
	formatter := NewFormatter("format:cyan(yes)", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes\n\n\n")
}

func TestPrettyFormatWhiteColor(t *testing.T) {
	formatter := NewFormatter("format:white(yes)", false, true, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[37myes\x1b[0m\n\n\n")
}

func TestPrettyFormatWhiteNoColor(t *testing.T) {
	formatter := NewFormatter("format:white(yes)", false, false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes\n\n\n")
}

func comment() *gc.Comment {
	id := "abcabcabcabc"
	ref := &gc.FileRef{"src/file.c", 12, gc.RefLineTypeOld}
	author := &gc.Person{"Simon", "iceking@example.com", time.Now(), "+0200"}
	comment := gc.NewComment("new comment\nmore context", "123444abcabc", ref, author).Success.(*gc.Comment)
	comment.ID = &id
	return comment
}
