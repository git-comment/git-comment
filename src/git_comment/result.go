package git_comment

import (
	"github.com/kylef/result.go/src/result"
)

// Converts an error into a failure result,
// or nil into a success result
func BoolResult(value bool, err error) result.Result {
	if err != nil {
		return result.NewFailure(err)
	}
	return result.NewSuccess(value)
}
