package main

import (
	gx "exec"
	"fmt"
	gc "libgitcomment"
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
	Disco         = "disco"
	ShortFormat   = "blue([%h] %c %an <%ae>)%nyellow(%t)"
	FullFormat    = "commit  %H%ncomment %C%nAuthor: %an <%ae>%n%b"
	discoFormat   = "cyan(%an) blue(<%ae>)%n[%h][%c] blue(%ad)%n%nyellow(%b)"
	RawFormat     = "yellow(comment %C)%n%v"
	formatPrefix  = "format:"
	invalidFormat = "Unknown pretty format."
	lineNumberMax = 5
)

type Formatter struct {
	format         string
	useLineNumbers bool
	useColor       bool
	useMargin      bool
	termWidth      uint16
	colorMapping   map[string]string
	indent         string
}

func NewFormatter(format string, useLineNumbers, useColor, useMargin bool, termWidth uint16) *Formatter {
	var indent = "\n  "
	var colorMapping map[string]string
	if useLineNumbers {
		indent = "\n            "
	}
	if useColor {
		colorMapping = map[string]string{
			black:      gx.Black,
			red:        gx.Red,
			green:      gx.Green,
			yellow:     gx.Yellow,
			blue:       gx.Blue,
			magenta:    gx.Magenta,
			cyan:       gx.Cyan,
			white:      gx.White,
			resetColor: gx.Clear,
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
	return &Formatter{format, useLineNumbers, useColor, useMargin, termWidth, colorMapping, indent}
}

func (f *Formatter) FormatLine(line *gc.DiffLine) string {
	prefix := f.formatLinePrefix(line)
	number := f.formatLineNumbers(line.OldLineNumber, line.NewLineNumber)
	content := f.formatLineContent(line)
	return fmt.Sprintf("%v %v %v", prefix, number, content)
}

func (f *Formatter) FormatComment(comment *gc.Comment) string {
	var content string
	switch {
	case f.format == Short || len(f.format) == 0:
		content = f.substituteVariables(ShortFormat, comment)
	case f.format == Full:
		content = f.substituteVariables(FullFormat, comment)
	case f.format == Disco:
		content = f.substituteVariables(discoFormat, comment)
	case f.format == Raw:
		format := string(f.substituteVariables(RawFormat, comment))
		content = fmt.Sprintf(format, comment.Serialize())
	case strings.HasPrefix(f.format, formatPrefix):
		content = f.substituteVariables(f.format[len(formatPrefix):], comment)
	}

	var components []byte
	for _, lineContent := range strings.Split(content, "\n") {
		if f.useMargin {
			components = append(components, []byte(fmt.Sprintf("%s%s│%s%s", f.indent, f.colorMapping["magenta"], f.colorMapping["resetColor"], lineContent))...)
		} else {
			components = append(components, []byte(fmt.Sprintf("%s\n", lineContent))...)
		}
	}
	components = append(components, []byte("\n\n")...)
	return string(components)
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
		newLine = gx.Colorize(gx.Green, f.formatLineNumber(newNum), f.useColor)
	}
	if oldNum != -1 {
		oldLine = gx.Colorize(gx.Red, oldLine, f.useColor)
	}
	return fmt.Sprintf("%v%v", oldLine, newLine)
}

func (f *Formatter) FormatFilePath(file *gc.DiffFile) string {
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

func (f *Formatter) formatLinePrefix(line *gc.DiffLine) string {
	switch line.Type {
	case gc.DiffAdd, gc.DiffAddNewline:
		return gx.Colorize(gx.Green, "+", f.useColor)
	case gc.DiffRemove, gc.DiffRemoveNewline:
		return gx.Colorize(gx.Red, "-", f.useColor)
	default:
		return " "
	}
}

func (f *Formatter) formatLineContent(line *gc.DiffLine) string {
	switch line.Type {
	case gc.DiffAddNewline, gc.DiffRemoveNewline:
		return "↵"
	case gc.DiffAdd:
		return gx.Colorize(gx.Green, line.Content, f.useColor)
	case gc.DiffRemove:
		return gx.Colorize(gx.Red, line.Content, f.useColor)
	default:
		return line.Content
	}
}

func (f *Formatter) substituteVariables(format string, comment *gc.Comment) string {
	return replaceAll(replaceAll(format, f.commentMapping(comment)), f.colorMapping)
}

func (f *Formatter) commentMapping(comment *gc.Comment) map[string]string {
	var path = ""
	var line = ""
	if comment.FileRef != nil {
		path = comment.FileRef.Path
		line = fmt.Sprintf("%v", comment.FileRef.Line)
	}
	return map[string]string{
		authorName:           comment.Author.Name,
		authorEmail:          comment.Author.Email,
		authorDateISO8601:    comment.Author.Date.Format(time.RFC3339),
		authorDateUnix:       fmt.Sprintf("%v", comment.Author.Date.Unix()),
		committerName:        comment.Amender.Name,
		committerEmail:       comment.Amender.Email,
		committerDateISO8601: comment.Amender.Date.Format(time.RFC3339),
		committerDateUnix:    fmt.Sprintf("%v", comment.Amender.Date.Unix()),
		commentFull:          *comment.ID,
		commentShort:         (*comment.ID)[:7],
		commitFull:           *comment.Commit,
		commitShort:          (*comment.Commit)[:7],
		bodyContent:          comment.Content,
		titleLine:            comment.Title(),
		filePath:             path,
		lineNumber:           line,
		newLine:              "\n",
		dividerLine:          strings.Repeat("-", int(f.termWidth)),
	}
}

func replaceAll(format string, substitutions map[string]string) string {
	for key, value := range substitutions {
		format = strings.Replace(format, key, value, -1)
	}
	return format
}
