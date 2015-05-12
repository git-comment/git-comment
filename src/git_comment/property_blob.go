package git_comment

import (
	"bytes"
	"strings"
)

type PropertyBlob struct {
	Properties map[string]string
	Message    string
}

const lineSeparator string = "\n"
const itemSeparator string = " "

func CreatePropertyBlob(content string) *PropertyBlob {
	props := map[string]string{}
	parsingMessage := false
	lines := strings.Split(content, lineSeparator)
	count := len(lines)
	var message bytes.Buffer
	for index, line := range lines {
		if len(line) == 0 {
			if parsingMessage {
				message.WriteString(lineSeparator)
			} else {
				parsingMessage = true
			}
		} else if parsingMessage {
			message.WriteString(line)
			if index < count-1 {
				message.WriteString(lineSeparator)
			}
		} else {
			property := strings.SplitN(line, itemSeparator, 2)
			if len(property) == 2 {
				name := property[0]
				value := property[1]
				props[name] = value
			}
		}
	}
	return &PropertyBlob{props, message.String()}
}

func (p *PropertyBlob) RawContent() string {
	var content bytes.Buffer
	for name, value := range p.Properties {
		content.WriteString(name)
		content.WriteString(itemSeparator)
		content.WriteString(value)
		content.WriteString(lineSeparator)
	}
	if len(p.Message) > 0 {
		content.WriteString(lineSeparator)
		content.WriteString(p.Message)
	}
	return content.String()
}
