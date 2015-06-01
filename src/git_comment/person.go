package git_comment

import (
	"fmt"
	"regexp"
)

type Person struct {
	Name  string
	Email string
}

// Parse a property string and create a person. The expected format is:
// ```
// Name <email@example.com>
// ```
// If a valid person cannot be created, an error is returned instead
func CreatePerson(properties string) *Person {
	fullRe := regexp.MustCompile(`(.*)\s<(.*@.*)>$`)
	match := fullRe.FindStringSubmatch(properties)
	if len(match) == 3 {
		return &Person{match[1], match[2]}
	} else {
		emailRe := regexp.MustCompile(`\s?<(.*@.*)>$`)
		match = emailRe.FindStringSubmatch(properties)
		if len(match) == 2 {
			return &Person{"", match[1]}
		}
		return &Person{properties, ""}
	}
}

func (p *Person) Serialize() string {
	if len(p.Name) > 0 {
		if len(p.Email) > 0 {
			return fmt.Sprintf("%v <%v>", p.Name, p.Email)
		}
		return p.Name
	} else if len(p.Email) > 0 {
		return fmt.Sprintf("<%v>", p.Email)
	}
	return ""
}
