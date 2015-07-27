package git_comment

import (
	"bytes"
	"github.com/cevaris/ordered_map"
	"strings"
	"time"
)

type PropertyBlob struct {
	properties *ordered_map.OrderedMap
	Message    string
}

const lineSeparator string = "\n"
const itemSeparator string = " "

func NewPropertyBlob() *PropertyBlob {
	props := ordered_map.NewOrderedMap()
	return &PropertyBlob{props, ""}
}

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

func (p *PropertyBlob) Set(name string, value interface{}) {
	p.properties.Set(name, value)
}

func (p *PropertyBlob) Get(property string) *string {
	prop, ok := p.properties.Get(property)
	if !ok {
		return nil
	}
	as, ok := prop.(string)
	if ok {
		return &as
	}
	return nil
}

func (p *PropertyBlob) GetTime(property string) *time.Time {
	prop := p.Get(property)
	if prop == nil {
		return nil
	}
	stamp, err := time.Parse(time.RFC822Z, *prop)
	if err != nil {
		return nil
	}
	return &stamp
}

func (p *PropertyBlob) GetPerson(property string) *Person {
	prop := p.Get(property)
	if prop == nil {
		return nil
	}
	if person, _ := CreatePerson(*prop).Dematerialize(); person != nil {
		return person.(*Person)
	}
	return nil
}

func (p *PropertyBlob) GetFileRef(property string) *FileRef {
	prop := p.Get(property)
	if prop == nil {
		return nil
	}
	return CreateFileRef(*prop, false)
}

func (p *PropertyBlob) Serialize() string {
	var content bytes.Buffer
	for kv := range p.properties.Iter() {
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
