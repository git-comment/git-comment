package search

import (
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve"
	gc "github.com/git-comment/git-comment"
	gg "github.com/git-comment/git-comment/git"
	"github.com/kylef/result.go/src/result"
	git "gopkg.in/libgit2/git2go.v23"
)

const (
	indexFilePath = "index"
)

type CommentIndex struct {
	Author  string
	Amender string
	Commit  string
	Content string
	FileRef string
}

// Find all comments matching text
// @return result.Result<[]*Comment, error>
func CommentsWithContent(repoPath, content string) result.Result {
	return openIndex(repoPath, func(repo *git.Repository, index bleve.Index) result.Result {
		query := bleve.NewQueryStringQuery(content)
		request := bleve.NewSearchRequest(query)
		return result.NewResult(index.Search(request)).FlatMap(func(match interface{}) result.Result {
			hits := match.(*bleve.SearchResult).Hits
			comments := make([]*gc.Comment, len(hits))
			for idx, hit := range hits {
				gc.CommentByID(repo, hit.ID).FlatMap(func(comment interface{}) result.Result {
					comments[idx] = comment.(*gc.Comment)
					return result.Result{}
				})
			}
			return result.NewSuccess(comments)
		})
	})
}

// @return result.Result<bool, error>
func IndexComments(repoPath string) result.Result {
	return openIndex(repoPath, func(repo *git.Repository, index bleve.Index) result.Result {
		results := make([]result.Result, 0)
		batch := index.NewBatch()
		return gg.CommentRefIterator(repo, func(ref *git.Reference) {
			gc.CommentFromRef(repo, ref.Name()).FlatMap(func(c interface{}) result.Result {
				comment := c.(*gc.Comment)
				err := batch.Index(*comment.ID, commentIndex(comment))
				results = append(results, gg.BoolResult(true, err))
				return result.NewSuccess(true)
			})
		}).FlatMap(func(value interface{}) result.Result {
			return result.Combine(func(values ...interface{}) result.Result {
				return gg.BoolResult(true, index.Batch(batch))
			}, results...)
		})
	})
}

// @return result.Result<bool, error>
func IndexComment(repoPath string, comment *gc.Comment) result.Result {
	return openIndex(repoPath, func(repo *git.Repository, index bleve.Index) result.Result {
		return gg.BoolResult(true, index.Index(*comment.ID, commentIndex(comment)))
	})
}

// Open or create a search index
// @return result.Result<bleve.Index, error>
func openIndex(repoPath string, ifSuccess func(*git.Repository, bleve.Index) result.Result) result.Result {
	storage := filepath.Join(repoPath, gc.CommentStorageDir)
	indexPath := filepath.Join(storage, indexFilePath)
	return gg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
		os.Mkdir(storage, 0700)
		success := func(index interface{}) result.Result {
			return ifSuccess(repo, index.(bleve.Index))
		}
		return result.NewResult(bleve.Open(indexPath)).Analysis(success, func(err error) result.Result {
			mapping := bleve.NewIndexMapping()
			index := result.NewResult(bleve.New(indexPath, mapping))
			return index.FlatMap(success)
		})
	})
}

func commentIndex(comment *gc.Comment) *CommentIndex {
	var filePath = ""
	if comment.FileRef != nil {
		filePath = comment.FileRef.Serialize()
	}
	return &CommentIndex{
		comment.Author.Serialize(),
		comment.Amender.Serialize(),
		*comment.Commit,
		comment.Content,
		filePath,
	}
}
