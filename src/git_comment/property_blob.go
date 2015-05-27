package git_comment

import (
	"bytes"
	"strings"
	"github.com/cevaris/ordered_map"
)

type PropertyBlob struct {
	Properties *ordered_map.OrderedMap
	Message    string
}

const lineSeparator string = "\n"
const itemSeparator string = " "

func CreatePropertyBlob(content string) *PropertyBlob {
	props := ordered_map.NewOrderedMap()
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
				props.Set(name, value)
			}
		}
	}
	return &PropertyBlob{props, message.String()}
}

func (p *PropertyBlob) RawContent() string {
	var content bytes.Buffer
	for kv := range p.Properties.Iter() {
		name, nOk := kv.Key.(string)
		value, vOk := kv.Value.(string)
		if nOk && vOk {
			content.WriteString(name)
			content.WriteString(itemSeparator)
			content.WriteString(value)
			content.WriteString(lineSeparator)
		}
	}
	if len(p.Message) > 0 {
		content.WriteString(lineSeparator)
		content.WriteString(p.Message)
	}
	return content.String()
}
