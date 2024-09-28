package inmemory

import (
	"errors"

	"github.com/nabishec/graphapi/graph/model"
)

// MemoryResolver provides in-memory implementation
type MemoryResolver struct {
	posts    map[string]*model.Post
	comments map[string][]*model.Comment
}

func NewMemoryResolver() *MemoryResolver {
	return &MemoryResolver{
		posts:    make(map[string]*model.Post),
		comments: make(map[string][]*model.Comment),
	}
}

// implementation methods supported api in memory
func (r *MemoryResolver) AddPost(post *model.Post) error {
	r.posts[post.ID] = post
	return nil
}

func (r *MemoryResolver) AddComment(comment *model.Comment) error {
	r.comments[comment.Post.ID] = append(r.comments[comment.Post.ID], comment)
	return nil
}

func (r *MemoryResolver) GetPosts() ([]*model.Post, error) {
	var posts = make([]*model.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *MemoryResolver) GetPost(id string) (*model.Post, error) {
	post, exist := r.posts[id]
	if !exist {
		return nil, errors.New("requested post doesn't exist")
	}
	return post, nil
}

func (r *MemoryResolver) GetComments(postId string) ([]*model.Comment, error) {
	comments, exist := r.comments[postId]
	if !exist {
		return nil, errors.New("requested post doesn't exist")
	}
	return comments, nil
}
