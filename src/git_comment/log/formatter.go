package log

import (
	"fmt"
	gitc "git_comment"
	ex "git_comment/exec"
	"strings"
	"time"
)

const (
	commentFull          = "%C"
	commentShort         = "%c"
	commitFull           = "%H"
	commitShort          = "%h"
	filePath             = "%f"
	lineNumber           = "%L"
	authorName           = "%an"
	authorEmail          = "%ae"
	authorDateISO8601    = "%ad"
	authorDateUnix       = "%aU"
	committerName        = "%kn"
	committerEmail       = "%ke"
	committerDateISO8601 = "%kd"
	committerDateUnix    = "%kU"
	bodyContent          = "%b"
	titleLine            = "%t"
	newLine              = "%n"
	dividerLine          = "%d"
	black                = "black("
	red                  = "red("
	green                = "green("
	yellow               = "yellow("
	blue                 = "blue("
	magenta              = "magenta("
	cyan                 = "cyan("
	white                = "white("
	resetColor           = ")"
)

const (
	Short         = "short"
	Full          = "full"
	Raw           = "raw"
	ShortFormat   = "[%h] %c %an <%ae>\n%t\n\n"
	FullFormat    = "commit  %H\ncomment %C\nAuthor: %an <%ae>\n%b\n\n"
	RawFormat     = "comment %C\n%v\n\n"
	formatPrefix  = "format:"
	invalidFormat = "Unknown pretty format."
	lineNumberMax = 5
)

type Formatter struct {
	format         string
	useLineNumbers bool
	useColor       bool
	termWidth      uint16
	colorMapping   map[string]string
	indent         string
}

func NewFormatter(format string, useLineNumbers, useColor bool, termWidth uint16) *Formatter {
	var indent = "\n  "
	var colorMapping map[string]string
	if useLineNumbers {
		indent = "\n            "
	}
	if useColor {
		colorMapping = map[string]string{
			black:      ex.Black,
			red:        ex.Red,
			green:      ex.Green,
			yellow:     ex.Yellow,
			blue:       ex.Blue,
			magenta:    ex.Magenta,
			cyan:       ex.Cyan,
			white:      ex.White,
			resetColor: ex.Clear,
		}
	} else {
		colorMapping = map[string]string{
			black:      "",
			red:        "",
			green:      "",
			yellow:     "",
			blue:       "",
			magenta:    "",
			cyan:       "",
			white:      "",
			resetColor: "",
		}
	}
	return &Formatter{format, useLineNumbers, useColor, termWidth, colorMapping, indent}
}

func (f *Formatter) FormatLine(line *gitc.DiffLine) string {
	prefix := f.formatLinePrefix(line)
	number := f.formatLineNumbers(line.OldLineNumber, line.NewLineNumber)
	content := f.formatLineContent(line)
	return fmt.Sprintf("%v %v %v", prefix, number, content)
}

func (f *Formatter) FormatComment(comment *gitc.Comment) string {
	switch {
	case f.format == Short || len(f.format) == 0:
		return f.substituteVariables(ShortFormat, comment)
	case f.format == Full:
		return f.substituteVariables(FullFormat, comment)
	case f.format == Raw:
		format := string(f.substituteVariables(RawFormat, comment))
		return fmt.Sprintf(format, *comment.ID)
	case strings.HasPrefix(f.format, formatPrefix):
		format := f.substituteVariables(f.format[len(formatPrefix):], comment)
		return fmt.Sprintf("%v\n", f.substituteColors(format))
	}
	return ""
}

func (f *Formatter) formatLineNumber(number int) string {
	var line string
	if number < 0 {
		line = " "
	} else {
		line = fmt.Sprintf("%d", number)
	}
	for len(line) < lineNumberMax {
		line = fmt.Sprintf(" %v", line)
	}
	return line
}

func (f *Formatter) formatLineNumbers(oldNum, newNum int) string {
	if !f.useLineNumbers {
		return ""
	}
	var newLine string
	oldLine := f.formatLineNumber(oldNum)
	if oldNum == newNum || newNum == -1 {
		newLine = f.formatLineNumber(-1)
	} else {
		newLine = ex.Colorize(ex.Green, f.formatLineNumber(newNum), f.useColor)
	}
	if oldNum != -1 {
		oldLine = ex.Colorize(ex.Red, oldLine, f.useColor)
	}
	return fmt.Sprintf("%v%v", oldLine, newLine)
}

func (f *Formatter) FormatFilePath(file *gitc.DiffFile) string {
	var path string
	if file.OldPath == file.NewPath {
		path = file.OldPath
	} else if len(file.OldPath) > 0 && len(file.NewPath) > 0 {
		path = fmt.Sprintf("%v -> %v", file.OldPath, file.NewPath)
	} else if len(file.OldPath) > 0 {
		path = file.OldPath
	} else {
		path = file.NewPath
	}
	return fmt.Sprintf("\n%v\n", path)
}

func (f *Formatter) formatLinePrefix(line *gitc.DiffLine) string {
	switch line.Type {
	case gitc.DiffAdd, gitc.DiffAddNewline:
		return ex.Colorize(ex.Green, "+", f.useColor)
	case gitc.DiffRemove, gitc.DiffRemoveNewline:
		return ex.Colorize(ex.Red, "-", f.useColor)
	default:
		return " "
	}
}

func (f *Formatter) formatLineContent(line *gitc.DiffLine) string {
	switch line.Type {
	case gitc.DiffAddNewline, gitc.DiffRemoveNewline:
		return "↵"
	case gitc.DiffAdd:
		return ex.Colorize(ex.Green, line.Content, f.useColor)
	case gitc.DiffRemove:
		return ex.Colorize(ex.Red, line.Content, f.useColor)
	default:
		return line.Content
	}
}

func (f *Formatter) substituteColors(format string) string {
	return replaceAll(format, f.colorMapping)
}

func (f *Formatter) substituteVariables(format string, comment *gitc.Comment) string {
	return replaceAll(format, f.commentMapping(comment))
}

func (f *Formatter) commentMapping(comment *gitc.Comment) map[string]string {
	var path = ""
	var line = ""
	if comment.FileRef != nil {
		path = comment.FileRef.Path
		line = fmt.Sprintf("%v", comment.FileRef.Line)
	}
	return map[string]string{
		authorName:           comment.Author.Name,
		authorEmail:          comment.Author.Email,
		authorDateISO8601:    comment.CreateTime.Format(time.RFC3339),
		authorDateUnix:       fmt.Sprintf("%v", comment.CreateTime.Unix()),
		committerName:        comment.Amender.Name,
		committerEmail:       comment.Amender.Email,
		committerDateISO8601: comment.AmendTime.Format(time.RFC3339),
		committerDateUnix:    fmt.Sprintf("%v", comment.AmendTime.Unix()),
		commentFull:          *comment.ID,
		commentShort:         (*comment.ID)[:7],
		commitFull:           *comment.Commit,
		commitShort:          (*comment.Commit)[:7],
		bodyContent:          comment.Content,
		titleLine:            strings.Split(comment.Content, "\n")[0],
		filePath:             path,
		lineNumber:           line,
		newLine:              f.indent,
		dividerLine:          strings.Repeat("-", int(f.termWidth)),
	}
}

func replaceAll(format string, substitutions map[string]string) string {
	for key, value := range substitutions {
		format = strings.Replace(format, key, value, -1)
	}
	return format
}
