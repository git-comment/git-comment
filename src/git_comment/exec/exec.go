package exec

import (
	gitc "git_comment"
	"io"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

// Start an arbitrary command with arguments and wait got
// it to finish
func ExecCommand(program string, args ...string) error {
	cmd := exec.Command(program, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Open the configured pager and a writer for Stdin.
// When the process is complete, close the writer and
// invoke Wait() on the command.
func ExecPager(pwd string) (*exec.Cmd, io.WriteCloser, error) {
	pager := gitc.ConfiguredPager(pwd)
	cmd := exec.Command(*pager)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return cmd, pipe, nil
}

// Fetch an environment variable. If it has length,
// return a reference.
func GetEnv(name string) *string {
	if env := os.Getenv(name); len(env) > 0 {
		return &env
	}
	return nil
}

// Calculate the number of lines visible in the current
// terminal.
// Windows compatibility is uncertain.
func CalculateDimensions() (height uint16, width uint16) {
	var dimensions [4]uint16
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, 2, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0); err != 0 {
		return 0, 0
	}
	return dimensions[0], dimensions[1]
}
