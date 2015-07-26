package log

import (
	gitc "git_comment"
	"github.com/stvp/assert"
	"regexp"
	"testing"
	"time"
)

func TestPrettyFormatFilePath(t *testing.T) {
	formatter := NewFormatter("format:%f", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "src/file.c")
}

func TestPrettyFormatFileLine(t *testing.T) {
	formatter := NewFormatter("format:%L", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "12")
}

func TestPrettyFormatAuthorName(t *testing.T) {
	formatter := NewFormatter("format:%an", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "Simon")
}

func TestPrettyFormatExtraText(t *testing.T) {
	formatter := NewFormatter("format:hi, %an!", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "hi, Simon!")
}

func TestPrettyFormatAuthorEmail(t *testing.T) {
	formatter := NewFormatter("format:%ae", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "iceking@example.com")
}

func TestPrettyFormatAuthorDateUnix(t *testing.T) {
	formatter := NewFormatter("format:%aU", false, false, 0)
	dateRe := regexp.MustCompile(`^([0-9]{10})$`)
	match := dateRe.FindStringSubmatch(formatter.FormatComment(comment()))
	assert.Equal(t, len(match), 2)
}

func TestPrettyFormatAuthorDateISO(t *testing.T) {
	formatter := NewFormatter("format:%ad", false, false, 0)
	dateRe := regexp.MustCompile(`^([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}-[0-9]{2}:[0-9]{2})$`)
	match := dateRe.FindStringSubmatch(formatter.FormatComment(comment()))
	assert.Equal(t, len(match), 2)
}

func TestPrettyFormatCommitterName(t *testing.T) {
	formatter := NewFormatter("format:%kn", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "Simon")
}

func TestPrettyFormatCommitterEmail(t *testing.T) {
	formatter := NewFormatter("format:%ke", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "iceking@example.com")
}

func TestPrettyFormatCommitterDateUnix(t *testing.T) {
	formatter := NewFormatter("format:%kU", false, false, 0)
	dateRe := regexp.MustCompile(`^([0-9]{10})$`)
	match := dateRe.FindStringSubmatch(formatter.FormatComment(comment()))
	assert.Equal(t, len(match), 2)
}

func TestPrettyFormatCommitterDateISO(t *testing.T) {
	formatter := NewFormatter("format:%kd", false, false, 0)
	dateRe := regexp.MustCompile(`^([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}-[0-9]{2}:[0-9]{2})$`)
	match := dateRe.FindStringSubmatch(formatter.FormatComment(comment()))
	assert.Equal(t, len(match), 2)
}

func TestPrettyFormatBody(t *testing.T) {
	formatter := NewFormatter("format:%b", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "new comment\nmore context")
}

func TestPrettyFormatTitle(t *testing.T) {
	formatter := NewFormatter("format:%t", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "new comment")
}

func TestPrettyFormatID(t *testing.T) {
	formatter := NewFormatter("format:%C", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "abcabcabcabc")
}

func TestPrettyFormatIDShort(t *testing.T) {
	formatter := NewFormatter("format:%c", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "abcabca")
}

func TestPrettyFormatCommitID(t *testing.T) {
	formatter := NewFormatter("format:%H", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "123444abcabc")
}

func TestPrettyFormatCommitIDShort(t *testing.T) {
	formatter := NewFormatter("format:%h", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "123444a")
}

func TestPrettyFormatBlackColor(t *testing.T) {
	formatter := NewFormatter("format:black(yes)", false, true, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[30myes\x1b[0m")
}

func TestPrettyFormatBlackNoColor(t *testing.T) {
	formatter := NewFormatter("format:black(yes)", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes")
}

func TestPrettyFormatRedColor(t *testing.T) {
	formatter := NewFormatter("format:red(yes)", false, true, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[31myes\x1b[0m")
}

func TestPrettyFormatRedNoColor(t *testing.T) {
	formatter := NewFormatter("format:red(yes)", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes")
}

func TestPrettyFormatGreenColor(t *testing.T) {
	formatter := NewFormatter("format:green(yes)", false, true, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[32myes\x1b[0m")
}

func TestPrettyFormatGreenNoColor(t *testing.T) {
	formatter := NewFormatter("format:green(yes)", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes")
}

func TestPrettyFormatYellowColor(t *testing.T) {
	formatter := NewFormatter("format:yellow(yes)", false, true, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[33myes\x1b[0m")
}

func TestPrettyFormatYellowNoColor(t *testing.T) {
	formatter := NewFormatter("format:yellow(yes)", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes")
}

func TestPrettyFormatBlueColor(t *testing.T) {
	formatter := NewFormatter("format:blue(yes)", false, true, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[34myes\x1b[0m")
}

func TestPrettyFormatBlueNoColor(t *testing.T) {
	formatter := NewFormatter("format:blue(yes)", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes")
}

func TestPrettyFormatMagentaColor(t *testing.T) {
	formatter := NewFormatter("format:magenta(yes)", false, true, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[35myes\x1b[0m")
}

func TestPrettyFormatMagentaNoColor(t *testing.T) {
	formatter := NewFormatter("format:magenta(yes)", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes")
}

func TestPrettyFormatCyanColor(t *testing.T) {
	formatter := NewFormatter("format:cyan(yes)", false, true, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[36myes\x1b[0m")
}

func TestPrettyFormatCyanNoColor(t *testing.T) {
	formatter := NewFormatter("format:cyan(yes)", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes")
}

func TestPrettyFormatWhiteColor(t *testing.T) {
	formatter := NewFormatter("format:white(yes)", false, true, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "\x1b[37myes\x1b[0m")
}

func TestPrettyFormatWhiteNoColor(t *testing.T) {
	formatter := NewFormatter("format:white(yes)", false, false, 0)
	assert.Equal(t, formatter.FormatComment(comment()), "yes")
}

func comment() *gitc.Comment {
	id := "abcabcabcabc"
	ref := &gitc.FileRef{"src/file.c", 12}
	author := &gitc.Person{"Simon", "iceking@example.com", time.Now(), "+0200"}
	comment := gitc.NewComment("new comment\nmore context", "123444abcabc", ref, author).Success.(*gitc.Comment)
	comment.ID = &id
	return comment
}
