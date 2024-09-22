package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"

	"github.com/nabishec/graphapi/graph/model"

	"github.com/google/uuid"
)

// MUTTATION
// AddPost is the resolver for the addPost field.
func (r *mutationResolver) AddPost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	post := &model.Post{
		ID:            uuid.New().String(),
		Title:         input.Title,
		Text:          input.Text,
		AllowComments: input.AllowComments,
	}

	r.posts[post.ID] = post
	return post, nil
}

// AddComment is the resolver for the addComment field.
const maxCommentLength int = 2000

func (r *mutationResolver) AddComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	post, exist := r.posts[input.PostID]
	if !exist {
		return nil, errors.New("commented post doesn't exist")
	}
	if !post.AllowComments {
		return nil, errors.New("comments aren't allowed for post")
	}
	if len(input.Text) > maxCommentLength {
		return nil, errors.New("comment is too long")
	}

	//searching parent comment for return it
	var parentComment *model.Comment = nil
	if input.ParentID != nil {
		for _, parentComentExist := range r.comments[input.PostID] {
			if parentComentExist.ID == *input.ParentID {
				parentComment = parentComentExist
				break
			}
		}
	}

	comment := &model.Comment{
		ID:     uuid.New().String(),
		Post:   post,
		Parent: parentComment,
		Text:   input.Text,
	}

	// subscription
	if _, ok := r.subscribers[input.PostID]; ok {
		for _, ch := range r.subscribers[input.PostID] {
			ch <- comment
			// mabe error when channel is overflowing should look in feature
			// select {
			// case i<- comment:
			// default: // chanel overflow
			/*mabe we can we can somehow, in case overflow read data
			from the cahnnel< save it< create new buffered channel
			of lager size and fill it*/
		}
	}

	//delete if dont work
	r.comments[input.PostID] = append(r.comments[input.PostID], comment)

	//r.posts[input.PostID].Comments = r.comments[input.PostID] // added comments in post struct
	return comment, nil

}

//END MUTTATION

// QUERY
// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	var posts = make([]*model.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string, first *int, after *string) (*model.Post, error) {
	post, exist := r.posts[id]
	if !exist {
		return nil, errors.New("requested post doesn't exist")
	}
	// if !post.AllowComments {
	// 	return post, nil
	// }
	var err error
	post.Comments, err = r.PagintionComments(id, first, after)
	if err != nil {
		return nil, err
	}
	return post, nil

}

// when we request post we must use pagination and add comments
// Pagintion
func (r *queryResolver) PagintionComments(id string, first *int, after *string) (*model.CommentConnection, error) {
	var edges []*model.CommentEdge

	comments, exist := r.comments[id]
	if !exist || first == nil {
		return &model.CommentConnection{
			Edges: edges,
			PageInfo: &model.PageInfo{
				EndCursor:   nil,
				HasNextPage: false,
			},
		}, nil
	}

	var startIndex = 0
	if after != nil {
		for i, comment := range comments {
			if comment.ID == *after {
				startIndex = i + 1
			}
		}
	}
	endIndex := startIndex + *first

	var endCursor string

	for i, length := startIndex, len(comments); i < length && i < endIndex; i++ {
		comment := comments[i]
		edges = append(edges, &model.CommentEdge{
			Node:   comment,
			Cursor: comment.ID,
		})
		if (i == length-1) || (i == endIndex-1) {
			endCursor = comment.ID
		}
	}

	var hasNextPage = endIndex < len(comments)

	return &model.CommentConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
	}, nil

}

//END QUERY

// SUBSCRIPTION
// AddedComment is the resolver for the addedComment field.
func (r *subscriptionResolver) AddedComment(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	ch := make(chan *model.Comment)
	r.subscribers[postID] = append(r.subscribers[postID], ch)
	return ch, nil
}

//END SUBSCRIPTION

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
