package graph

import (
	"github.com/nabishec/graphapi/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// ServerResolvers defines methods supported api
type StorageService interface {
	AddPost(post *model.Post) error
	AddComment(comment *model.Comment) error
	GetPosts() ([]*model.Post, error)
	GetPost(id string) (*model.Post, error)
	GetComments(postId string) ([]*model.Comment, error)
}

type Resolver struct {
	DataResolvers StorageService
	subscribers   map[string][]chan *model.Comment
}

func NewResolver(repository StorageService) *Resolver {
	return &Resolver{
		DataResolvers: repository,
		subscribers:   make(map[string][]chan *model.Comment),
	}
}
