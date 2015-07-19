package exec

import (
	"fmt"
	"syscall"
	"unsafe"
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

// Calculate the number of lines visible in the current
// terminal.
// Windows compatibility is uncertain.
func CalculateDimensions() (height uint16, width uint16) {
	var dimensions [4]uint16
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0); err != 0 {
		return 0, 0
	}
	return dimensions[0], dimensions[1]
}

func Colorize(code, text string, active bool) string {
	if !active {
		return text
	}
	return fmt.Sprintf("%s%s%s", code, text, Clear)
}
