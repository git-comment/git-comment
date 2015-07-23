package git_comment

import (
	"github.com/stvp/assert"
	"testing"
)

func TestCompareVersionsSame(t *testing.T) {
	status, err := compareVersion("2.1.0", "2.1.0").Dematerialize()
	assert.Nil(t, err)
	assert.Equal(t, status, VersionStatusEqual)
}

func TestCompareVersionsToolNewer(t *testing.T) {
	status, err := compareVersion("2.1.1", "2.1.0").Dematerialize()
	assert.Nil(t, err)
	assert.Equal(t, status, VersionStatusUpgradeRepo)
}

func TestCompareVersionsRepoNewer(t *testing.T) {
	status, err := compareVersion("2.0.1", "2.1.0").Dematerialize()
	assert.Nil(t, err)
	assert.Equal(t, status, VersionStatusUpgradeTool)
}

func TestCompareVersionsToolCorrupted(t *testing.T) {
	status, err := compareVersion("2..1", "2.1.0").Dematerialize()
	assert.NotNil(t, err)
	assert.Nil(t, status)
}

func TestCompareVersionsRepoCorrupted(t *testing.T) {
	status, err := compareVersion("1.4.5", "2.0").Dematerialize()
	assert.NotNil(t, err)
	assert.Nil(t, status)
}
