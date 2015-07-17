package main

import (
	"fmt"
	gitc "git_comment"
	ex "git_comment/exec"
	"github.com/kylef/result.go/src/result"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"os/exec"
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
	committerName        = "%cn"
	committerEmail       = "%ce"
	committerDateISO8601 = "%cd"
	committerDateUnix    = "%cU"
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
)

var (
	buildVersion string
	termWidth    uint16
	termHeight   uint16
	useColor     bool
	colorMapping map[string]string
	app          = kp.New("git-comment-log", "List git commit comments")
	pretty       = app.Flag("pretty", "Pretty-print the comments in a format such as short, full, raw, or custom placeholders.").String()
	noPager      = app.Flag("nopager", "Disable pager").Bool()
	noColor      = app.Flag("nocolor", "Disable color").Bool()
	lineNumbers  = app.Flag("line-numbers", "Show line numbers").Bool()
	revision     = app.Arg("revision range", "Filter comments to comments on commits from the specified range").String()
	indent       = "\n  "
)

func main() {
	app.Version(buildVersion)
	kp.MustParse(app.Parse(os.Args[1:]))
	pwd, err := os.Getwd()
	app.FatalIfError(err, "pwd")
	configureEnv()
	showComments(pwd)
}

func configureEnv() {
	if !*noColor {
		if wd, err := os.Getwd(); err == nil {
			useColor = gitc.ConfiguredBool(wd, "color.pager", false)
		}
	}
	termHeight, termWidth = ex.CalculateDimensions()
	if *lineNumbers {
		indent = "\n            "
	}
}

func fatalIfError(r result.Result, code string) interface{} {
	app.FatalIfError(r.Failure, code)
	return r.Success
}

func showComments(pwd string) {
	var usePager bool = termHeight == 0 && !*noPager
	var content []byte
	var writer io.WriteCloser
	var cmd *exec.Cmd
	var err error
	diff := fatalIfError(gitc.DiffCommits(pwd, *revision), "diff")
	pageContent := func(data string) {
		content = append(content, []byte(data)...)
		if !usePager {
			lines := strings.Split(string(content), "\n")
			usePager = !*noPager && uint16(len(lines)) > termHeight-1
		}
		if usePager {
			if writer == nil {
				cmd, writer, err = ex.ExecPager(pwd)
				app.FatalIfError(err, "pager")
			}
			if len(content) > 0 {
				_, err = writer.Write(content)
				content = []byte{}
				app.FatalIfError(err, "writer")
			}
		}
	}
	for _, file := range diff.(*gitc.Diff).Files {
		var printedFileHeader = false
		var afterComment = false
		beforeBuffer := make([]*gitc.DiffLine, 0)
		afterBuffer := make([]*gitc.DiffLine, 0)
		for _, line := range file.Lines {
			if len(line.Comments) > 0 {
				for _, line := range afterBuffer {
					pageContent(formattedLine(line))
				}
				afterBuffer = make([]*gitc.DiffLine, 0)
				if !printedFileHeader {
					pageContent(formattedFilePath(file))
					printedFileHeader = true
				}
				for _, line := range beforeBuffer {
					pageContent(formattedLine(line))
				}
				pageContent(formattedLine(line))
				for _, comment := range line.Comments {
					pageContent(formattedComment(comment))
				}
				beforeBuffer = make([]*gitc.DiffLine, 0)
				afterComment = true
			} else {
				if afterComment {
					afterBuffer = append(afterBuffer, line)
					if len(afterBuffer) == 5 {
						for _, line := range afterBuffer {
							pageContent(formattedLine(line))
						}
						afterBuffer = make([]*gitc.DiffLine, 0)
						afterComment = false
					}
				} else {
					beforeBuffer = append(beforeBuffer, line)
					if len(beforeBuffer) > 5 {
						beforeBuffer = append(beforeBuffer[:0], beforeBuffer[1:]...)
					}
				}
			}
		}
	}
	if !usePager {
		fmt.Println(string(content))
	}

	if writer != nil {
		writer.Close()
		cmd.Wait()
	}
}

func formattedLineNumber(number int) string {
	const lineNumberMax = 5
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

func formattedLineNumbers(oldNum, newNum int) string {
	if !*lineNumbers {
		return ""
	}
	var newLine string
	oldLine := formattedLineNumber(oldNum)
	if oldNum == newNum || newNum == -1 {
		newLine = formattedLineNumber(-1)
	} else {
		newLine = ex.Colorize(ex.Green, formattedLineNumber(newNum), useColor)
	}
	if oldNum != -1 {
		oldLine = ex.Colorize(ex.Red, oldLine, useColor)
	}
	return fmt.Sprintf("%v%v", oldLine, newLine)
}

func formattedFilePath(file *gitc.DiffFile) string {
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
	return fmt.Sprintf("%v\n", path)
}

func formattedLinePrefix(line *gitc.DiffLine) string {
	switch line.Type {
	case gitc.DiffAdd, gitc.DiffAddNewline:
		return ex.Colorize(ex.Green, "+", useColor)
	case gitc.DiffRemove, gitc.DiffRemoveNewline:
		return ex.Colorize(ex.Red, "-", useColor)
	default:
		return " "
	}
}

func formattedLineContent(line *gitc.DiffLine) string {
	switch line.Type {
	case gitc.DiffAddNewline, gitc.DiffRemoveNewline:
		return "↵"
	case gitc.DiffAdd:
		return ex.Colorize(ex.Green, line.Content, useColor)
	case gitc.DiffRemove:
		return ex.Colorize(ex.Red, line.Content, useColor)
	default:
		return line.Content
	}
}

func formattedLine(line *gitc.DiffLine) string {
	prefix := formattedLinePrefix(line)
	number := formattedLineNumbers(line.OldLineNumber, line.NewLineNumber)
	content := formattedLineContent(line)
	return fmt.Sprintf("%v %v %v", prefix, number, content)
}

func formattedComment(comment *gitc.Comment) string {
	if *pretty == Short || len(*pretty) == 0 {
		return substituteVariables(ShortFormat, comment)
	} else if *pretty == Full {
		return substituteVariables(FullFormat, comment)
	} else if *pretty == Raw {
		format := string(substituteVariables(RawFormat, comment))
		return fmt.Sprintf(format, *comment.ID)
	} else if strings.HasPrefix(*pretty, formatPrefix) {
		format := substituteVariables((*pretty)[len(formatPrefix):], comment)
		return substituteColors(format)
	}
	app.FatalUsage(invalidFormat)
	return ""
}

func substituteColors(format string) string {
	if !useColor {
		return format
	}
	if len(colorMapping) == 0 {
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
	}
	return replaceAll(format, colorMapping)
}

func substituteVariables(format string, comment *gitc.Comment) string {
	var path = ""
	var line = ""
	if len(comment.FileRef.Path) > 0 {
		path = comment.FileRef.Path
	}
	if comment.FileRef.Line > 0 {
		line = fmt.Sprintf("%v", comment.FileRef.Line)
	}
	mapping := map[string]string{
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
		newLine:              indent,
		dividerLine:          strings.Repeat("-", int(termWidth)),
	}

	return replaceAll(format, mapping)
}

func replaceAll(format string, substitutions map[string]string) string {
	for key, value := range substitutions {
		format = strings.Replace(format, key, value, -1)
	}
	return format
}
