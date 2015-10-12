package main

import (
	"errors"
	kp "gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	gg "libgitcomment/git"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	noMessageProvided      = "Aborting comment, no message provided"
	defaultTemplateName    = ".gitcommenttemplate"
	defaultMessageTemplate = "\n# Enter comment content\n# Lines beginning with '#' will be stripped"
)

func getMessageFromEditor(app *kp.Application, repoPath string) string {
	editor := gg.ConfiguredEditor(repoPath)
	file, err := ioutil.TempFile("", "gitc")
	app.FatalIfError(err, "io")
	path := file.Name()
	file.Write(commentTemplateText(app, repoPath))
	file.Close()
	err = ExecCommand(editor, path)
	app.FatalIfError(err, "io")
	content, err := ioutil.ReadFile(path)
	os.Remove(path)
	app.FatalIfError(err, "io")
	return sanitizeMessage(app, string(content))
}

func sanitizeMessage(app *kp.Application, message string) string {
	reg := regexp.MustCompile("(?m)^#.*$")
	stripped := reg.ReplaceAllString(message, "")
	content := strings.TrimSpace(stripped)
	if len(content) == 0 {
		app.FatalIfError(errors.New(noMessageProvided), "io")
	}
	return content
}

func commentTemplateText(app *kp.Application, repoPath string) []byte {
	if config := gg.ConfiguredString(repoPath, "comment.template", ""); len(config) > 0 {
		return contentsOfFile(app, config)
	}
	if defaultConfig := contentsOfFile(nil, defaultTemplatePath()); len(defaultConfig) > 0 {
		return defaultConfig
	}
	return []byte(defaultMessageTemplate)
}

func contentsOfFile(app *kp.Application, filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		if app != nil {
			app.FatalIfError(err, "io")
		} else {
			return []byte{}
		}
	}
	return []byte(content)
}

func defaultTemplatePath() string {
	if current, err := user.Current(); err == nil {
		return filepath.Join(current.HomeDir, defaultTemplateName)
	}
	return ""
}
