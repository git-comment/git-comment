package git_comment

import (
	gitg "git_comment/git"
	"github.com/blevesearch/bleve"
	"github.com/kylef/result.go/src/result"
	git "github.com/libgit2/git2go"
	"os"
	"path/filepath"
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
// @return result.Result<[]*CommentIndex, error>
func CommentsWithContent(repoPath, content string) result.Result {
	return openIndex(repoPath, func(repo *git.Repository, index bleve.Index) result.Result {
		query := bleve.NewQueryStringQuery(content)
		request := bleve.NewSearchRequest(query)
		return result.NewResult(index.Search(request))
	}).FlatMap(func(match interface{}) result.Result {
		hits := match.(bleve.SearchResult).Hits
		indices := make([]*CommentIndex, len(hits))
		for idx, hit := range hits {
			indices[idx] = hitIndex(hit.Fields)
		}
		return result.NewSuccess(indices)
	})
}

// @return result.Result<bool, error>
func IndexComments(repoPath string) result.Result {
	return openIndex(repoPath, func(repo *git.Repository, index bleve.Index) result.Result {
		results := make([]result.Result, 0)
		return gitg.CommentRefIterator(repo, func(ref *git.Reference) {
			CommentFromRef(repo, ref.Name()).FlatMap(func(c interface{}) result.Result {
				comment := c.(*Comment)
				err := index.Index(*comment.ID, commentIndex(comment))
				results = append(results, gitg.BoolResult(true, err))
				return result.NewSuccess(true)
			})
		}).FlatMap(func(value interface{}) result.Result {
			return result.Combine(func(values ...interface{}) result.Result {
				return result.NewSuccess(true)
			}, results...)
		})
	})
}

// @return result.Result<bool, error>
func IndexComment(repoPath string, comment *Comment) result.Result {
	return openIndex(repoPath, func(repo *git.Repository, index bleve.Index) result.Result {
		return gitg.BoolResult(true, index.Index(*comment.ID, commentIndex(comment)))
	})
}

// Open or create a search index
// @return result.Result<bleve.Index, error>
func openIndex(repoPath string, ifSuccess func(*git.Repository, bleve.Index) result.Result) result.Result {
	storage := filepath.Join(repoPath, CommentStorageDir)
	indexPath := filepath.Join(storage, indexFilePath)
	return gitg.WithRepository(repoPath, func(repo *git.Repository) result.Result {
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

func hitIndex(hit map[string]interface{}) *CommentIndex {
	return &CommentIndex{
		hit["Author"].(string),
		hit["Amender"].(string),
		hit["Commit"].(string),
		hit["Content"].(string),
		hit["FileRef"].(string),
	}
}

func commentIndex(comment *Comment) *CommentIndex {
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
