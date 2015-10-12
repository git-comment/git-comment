package libgitcomment

import (
	gg "../git"
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/kylef/result.go/src/result"
	git "gopkg.in/libgit2/git2go.v23"
	"path/filepath"
)

type VersionStatus int

const (
	VersionStatusEqual       VersionStatus = 0
	VersionStatusUpgradeRepo VersionStatus = 1
	VersionStatusUpgradeTool VersionStatus = 2
)

const (
	versionRef       = "version"
	toolInvalidError = "git-comment version corrupted. Please file a bug report.\ntool: %v\nrepo: %v"
	upgradeMessage   = "Updating git-comment version in use"
	upgradeToolError = "The version of git-comment used in this repository is newer than the version installed. Please upgrade."
	upgradeRepoError = "The version of git-comment used in this repository is out of date. Please upgrade by running `git-comment --update`"
)

// Check the version of git-comment in use against
// the version in use in the repository.
// @return result.Result<VersionStatus, error>
func VersionCheck(repoPath, toolVersion string) result.Result {
	return gg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
		return readVersion(repo).Analysis(func(version interface{}) result.Result {
			return compareVersion(toolVersion, version.(string))
		}, func(err error) result.Result {
			if git.IsErrorCode(err, git.ErrNotFound) {
				return writeVersion(repo, toolVersion)
			}
			return result.NewFailure(err)
		})
	})
}

// Migrate the repo version to the installed version of
// the tool
func VersionUpdate(repoPath, toolVersion string) error {
	return nil
}

// @return result.Result<VersionStatus, error>
func compareVersion(toolVersion, repoVersion string) result.Result {
	errorMsg := fmt.Sprintf(toolInvalidError, toolVersion, repoVersion)
	invalidErr := result.NewFailure(errors.New(errorMsg))
	vt := result.NewResult(semver.Make(toolVersion)).RecoverWith(invalidErr)
	vr := result.NewResult(semver.Make(repoVersion)).RecoverWith(invalidErr)
	return result.Combine(func(values ...interface{}) result.Result {
		vt, vr := values[0].(semver.Version), values[1].(semver.Version)
		return comparisonStatus(vt.Compare(vr))
	}, vt, vr)
}

func comparisonStatus(code int) result.Result {
	switch code {
	case -1:
		return result.NewFailure(errors.New(upgradeToolError))
	case 1:
		return result.NewFailure(errors.New(upgradeRepoError))
	default:
		return result.NewSuccess(VersionStatusEqual)
	}
}

// @return result.Result<VersionStatus, error>
func writeVersion(repo *git.Repository, version string) result.Result {
	oid := result.NewResult(repo.CreateBlobFromBuffer([]byte(version)))
	return oid.FlatMap(func(oid interface{}) result.Result {
		path := filepath.Join(gg.CommentRefBase, versionRef)
		return result.NewResult(repo.References.Create(path,
			oid.(*git.Oid), false, upgradeMessage))
	}).FlatMap(func(ref interface{}) result.Result {
		return result.NewSuccess(VersionStatusEqual)
	})
}

// @return result.Result<string, error>
func readVersion(repo *git.Repository) result.Result {
	path := filepath.Join(gg.CommentRefBase, versionRef)
	ref := result.NewResult(repo.References.Lookup(path))
	return ref.FlatMap(func(ref interface{}) result.Result {
		oid := ref.(*git.Reference).Target()
		return result.NewResult(repo.LookupBlob(oid))
	}).FlatMap(func(blob interface{}) result.Result {
		contents := blob.(*git.Blob).Contents()
		return result.NewSuccess(string(contents))
	})
}
