package git_comment

import (
	"errors"
	"fmt"
	"github.com/kylef/result.go/src/result"
	"regexp"
	"strconv"
	"time"
)

type Person struct {
	Name       string
	Email      string
	Date       time.Time
	TimeOffset string
}

const (
	invalidPersonError = "Person could not be created from input"
)

// Parse a property string and create a person. The expected format is:
// ```
// Name <email@example.com>
// ```
// If a valid person cannot be created, an error is returned instead
// @return result.Result<*Person, error>
func CreatePerson(properties string) result.Result {
	fullRe := regexp.MustCompile(`(?:(.*)\s)?<(.*@.*)>(?:\s([0-9]+)\s([\-+][0-9]{4}))?`)
	match := fullRe.FindStringSubmatch(properties)
	invalidErr := result.NewFailure(errors.New(invalidPersonError))
	if len(match) == 0 {
		return invalidErr
	}
	name, email := match[1], match[2]
	timestamp := result.NewResult(strconv.ParseInt(match[3], 10, 64))
	return timestamp.Analysis(func(value interface{}) result.Result {
		stamp := time.Unix(value.(int64), 0)
		person := &Person{match[1], match[2], stamp, match[4]}
		return result.NewSuccess(person)
	}, func(err error) result.Result {
		timestamp := time.Now()
		return result.NewSuccess(&Person{name, email, timestamp, timestamp.Format("-0700")})
	})
}

func (p *Person) Serialize() string {
	return fmt.Sprintf("%v <%v> %d %v", p.Name, p.Email, p.Date.Unix(), p.TimeOffset)
}
