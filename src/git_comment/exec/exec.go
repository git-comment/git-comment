package exec

import (
	gitc "git_comment"
	"github.com/kylef/result.go/src/result"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Start an arbitrary command with arguments and wait got
// it to finish
// Windows compatibility is uncertain
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
	pager := strings.Split(gitc.ConfiguredPager(pwd), " ")
	cmd := exec.Command(pager[0], pager[1:]...)
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

// Return the success value, otherwise kill the app with
// the error code specified
func FatalIfError(app *kp.Application, r result.Result, code string) interface{} {
	app.FatalIfError(r.Failure, code)
	return r.Success
}
