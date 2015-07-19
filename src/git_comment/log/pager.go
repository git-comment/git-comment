package log

import (
	"bytes"
	"fmt"
	gite "git_comment/exec"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os/exec"
)

type ContentPager struct {
	disablePager bool
	usePager     bool
	content      []byte
	writer       io.WriteCloser
	cmd          *exec.Cmd
	app          *kp.Application
	wd           string
	termHeight   uint16
}

func NewContentPager(app *kp.Application, wd string, termHeight uint16, disablePager bool) *ContentPager {
	pager := &ContentPager{}
	pager.app = app
	pager.wd = wd
	pager.disablePager = disablePager
	pager.termHeight = termHeight
	pager.usePager = termHeight == 0 && !disablePager
	return pager
}

func (p *ContentPager) AddContent(data string) {
	if p.disablePager {
		fmt.Println(data)
	} else {
		p.content = append(p.content, []byte(data)...)
		var err error
		if !p.usePager {
			lines := bytes.Count(p.content, []byte("\n")) + 1
			p.usePager = uint16(lines) > p.termHeight-1
		}
		if p.usePager {
			if p.writer == nil {
				p.cmd, p.writer, err = gite.ExecPager(p.wd)
				p.app.FatalIfError(err, "pager")
			}
			if len(p.content) > 0 {
				_, err = p.writer.Write(p.content)
				p.content = []byte{}
				p.app.FatalIfError(err, "writer")
			}
		}
	}
}

func (p *ContentPager) Finish() {
	if !p.disablePager && !p.usePager {
		fmt.Println(string(p.content))
	}
	if p.writer != nil {
		p.writer.Close()
		p.cmd.Wait()
	}
}
