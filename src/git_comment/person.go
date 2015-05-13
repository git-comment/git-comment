package git_comment

import (
	"errors"
	"regexp"
)

type Person struct {
	Name  string
	Email string
}

// Creates a person from the email address and name in the
// active git config
func ConfiguredPerson() *Person {
	return &Person{}
}

// Parse a property string and create a person. The expected format is:
//
// Name <email@example.com>
//
// If a valid person cannot be created, an error is returned instead
func CreatePerson(properties string) (*Person, error) {
	const invalidProperties = "Invalid property format for person"
	fullRe := regexp.MustCompile(`(.*)\s<(.*@.*)>$`)
	match := fullRe.FindStringSubmatch(properties)
	if len(match) == 3 {
		return &Person{match[1], match[2]}, nil
	} else {
		emailRe := regexp.MustCompile(`\s?<(.*@.*)>$`)
		match = emailRe.FindStringSubmatch(properties)
		if len(match) == 2 {
			return &Person{"", match[1]}, nil
		}
		return nil, errors.New(invalidProperties)
	}
}
