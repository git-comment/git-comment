package exec

import (
	"syscall"
	"unsafe"
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
