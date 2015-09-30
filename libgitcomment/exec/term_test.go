package exec

import (
	"testing"

	"github.com/stvp/assert"
)

func TestColorizeActive(t *testing.T) {
	text := Colorize(Green, "hello", true)
	assert.Equal(t, text, "\x1b[32mhello\x1b[0m")
}

func TestColorizeInactive(t *testing.T) {
	text := Colorize(Green, "hello", false)
	assert.Equal(t, text, "hello")
}
