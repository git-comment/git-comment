package git_comment

import (
	"fmt"
	"regexp"
	"strconv"
)

type FileRef struct {
	Path string
	Line int
}

// Create a ref from a format:
//
// ```
// file_path:line
// ```
//
// or
//
// ```
// file_path
// ```
//
func CreateFileRef(content string) *FileRef {
	lineRe := regexp.MustCompile(`(.*):(\d+)`)
	match := lineRe.FindStringSubmatch(content)
	if len(match) == 3 {
		line, parseErr := strconv.ParseInt(match[2], 0, 0)
		if parseErr == nil {
			return &FileRef{match[1], int(line)}
		}
	}
	return &FileRef{content, 0}
}

// Create a deserializable version of a ref
func (f *FileRef) Serialize() string {
	if f.Line > 0 {
		return fmt.Sprintf("%v:%d", f.Path, f.Line)
	}
	return f.Path
}
