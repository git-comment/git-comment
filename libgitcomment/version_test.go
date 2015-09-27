package libgitcomment

import (
	"github.com/stvp/assert"
	"strings"
	"testing"
)

func TestCompareVersionsSame(t *testing.T) {
	status, err := compareVersion("2.1.0", "2.1.0").Dematerialize()
	assert.Nil(t, err)
	assert.Equal(t, status, VersionStatusEqual)
}

func TestCompareVersionsToolNewer(t *testing.T) {
	status, err := compareVersion("2.1.1", "2.1.0").Dematerialize()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "git-comment --update"))
	assert.Nil(t, status)
}

func TestCompareVersionsRepoNewer(t *testing.T) {
	status, err := compareVersion("2.0.1", "2.1.0").Dematerialize()
	assert.NotNil(t, err)
	assert.Nil(t, status)
}

func TestCompareVersionsToolCorrupted(t *testing.T) {
	status, err := compareVersion("2..1", "2.1.0").Dematerialize()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "2..1"))
	assert.True(t, strings.Contains(err.Error(), "2.1.0"))
	assert.Nil(t, status)
}

func TestCompareVersionsRepoCorrupted(t *testing.T) {
	status, err := compareVersion("1.4.5", "2.0").Dematerialize()
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "1.4.5"))
	assert.True(t, strings.Contains(err.Error(), "2.0"))
	assert.Nil(t, status)
}
