package graph

import "graphapi/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	posts       map[string]*model.Post
	comments    map[string][]*model.Comment
	subscribers map[string][]chan *model.Comment
}

func NewResolver() *Resolver {
	return &Resolver{
		posts:       make(map[string]*model.Post),
		comments:    make(map[string][]*model.Comment),
		subscribers: make(map[string][]chan *model.Comment),
	}
}
