package exec

import (
	"fmt"
)

const (
	Black   = "\x1b[30m"
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan    = "\x1b[36m"
	White   = "\x1b[37m"
	Clear   = "\x1b[0m"
)

func Colorize(code, text string, active bool) string {
	if !active {
		return text
	}
	return fmt.Sprintf("%s%s%s", code, text, Clear)
}
