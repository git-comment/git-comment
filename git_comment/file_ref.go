package git_comment

import (
	"fmt"
	"regexp"
	"strconv"
)

type RefLineType int

const (
	RefLineTypeNew RefLineType = 0
	RefLineTypeOld RefLineType = 1
	oldRef                     = ":old"
)

type FileRef struct {
	Path     string
	Line     int
	LineType RefLineType
}

// Create a ref from a format:
//
// ```
// file_path:line:type
// ```
//
// or
//
// ```
// file_path
// ```
//
func CreateFileRef(content string, markDeleted bool) *FileRef {
	ref := DeserializeFileRef(content)
	if markDeleted {
		ref.LineType = RefLineTypeOld
	}
	return ref
}

func DeserializeFileRef(content string) *FileRef {
	lineRe := regexp.MustCompile(`(?U)^(.*)(?::(\d+)(:old)?)?$`)
	match := lineRe.FindStringSubmatch(content)
	var lineType RefLineType = RefLineTypeNew
	var line = 0
	if len(match) > 2 {
		line = deserializeLine(match[2])
	}
	if len(match) > 3 {
		lineType = deserializeLineType(match[3])
	}
	return &FileRef{match[1], line, lineType}
}

// Create a deserializable version of a ref
func (f *FileRef) Serialize() string {
	if f.Line > 0 {
		return fmt.Sprintf("%v:%d%v", f.Path, f.Line, serializeLineType(f.LineType))
	}
	return f.Path
}

func deserializeLine(lineText string) int {
	if line, parseErr := strconv.ParseInt(lineText, 0, 0); parseErr == nil {
		return int(line)
	}
	return 0
}

func serializeLineType(lineType RefLineType) string {
	switch lineType {
	case RefLineTypeOld:
		return oldRef
	default:
		return ""
	}
}

func deserializeLineType(lineType string) RefLineType {
	switch lineType {
	case oldRef:
		return RefLineTypeOld
	default:
		return RefLineTypeNew
	}
}
