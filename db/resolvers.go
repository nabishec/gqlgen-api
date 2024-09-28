package db

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabishec/graphapi/graph/model"
)

// DatabaseResolver provides database implementation
type DatabaseResolver struct {
	DB *sqlx.DB
}

func NewDatabaseResolver(db *sqlx.DB) *DatabaseResolver {
	return &DatabaseResolver{
		DB: db,
	}
}

// implementation methods supported api in database
func (r *DatabaseResolver) AddPost(post *model.Post) error {
	_, err := r.DB.Exec("INSERT INTO posts (id, title, content, allow_comments) VALUES ($1, $2, $3, $4)", idToUUID(post.ID), post.Title, post.Content, post.AllowComments)
	return err
}

func (r *DatabaseResolver) AddComment(comment *model.Comment) error {
	var err error
	if comment.Parent != nil {
		_, err = r.DB.Exec("INSERT INTO comments (id, post_id, parent_id, content) VALUES ($1, $2, $3, $4)", idToUUID(comment.ID), idToUUID(comment.Post.ID), idToUUID(comment.ID), comment.Content)
	} else {
		_, err = r.DB.Exec("INSERT INTO comments (id, post_id, content) VALUES ($1, $2, $3)", idToUUID(comment.ID), idToUUID(comment.Post.ID), comment.Content)
	}
	return err
}

func (r *DatabaseResolver) GetPosts() ([]*model.Post, error) {
	var postsDB []struct {
		ID            uuid.UUID `db:"id"`
		Title         string    `db:"title"`
		Content       string    `db:"content"`
		AllowComments bool      `db:"allow_comments"`
	}
	err := r.DB.Select(&postsDB, "SELECT * FROM posts")
	if err != nil {
		return nil, err
	}

	var posts []*model.Post
	for _, postDB := range postsDB {
		posts = append(posts, &model.Post{
			ID:            postDB.ID.String(),
			Title:         postDB.Title,
			Content:       postDB.Content,
			AllowComments: postDB.AllowComments,
		})
	}

	return posts, nil
}

func (r *DatabaseResolver) GetPost(id string) (*model.Post, error) {
	var postDB struct {
		ID            uuid.UUID `db:"id"`
		Title         string    `db:"title"`
		Content       string    `db:"content"`
		AllowComments bool      `db:"allow_comments"`
	}
	err := r.DB.Get(&postDB, "SELECT * FROM posts WHERE id=$1", idToUUID(id))
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		ID:            postDB.ID.String(),
		Title:         postDB.Title,
		Content:       postDB.Content,
		AllowComments: postDB.AllowComments,
	}
	return post, nil
}

func (r *DatabaseResolver) GetComments(postId string) ([]*model.Comment, error) {
	var commentsDB []struct {
		ID       uuid.UUID `db:"id"`
		PostID   uuid.UUID `db:"post_id"`
		ParentID uuid.UUID `db:"parent_id"`
		Content  string    `db:"content"`
	}
	err := r.DB.Select(&commentsDB, "SELECT * FROM comments WHERE post_id=$1", idToUUID(postId))
	if err != nil {
		return nil, err
	}

	var comments []*model.Comment
	for _, commentDB := range commentsDB {
		// Fetch the parent comment if it's not nil
		var parentComment *model.Comment
		if commentDB.ParentID != uuid.Nil {
			parentComment, _ = r.GetComment(commentDB.ParentID.String())
		}

		post, _ := r.GetPost(commentDB.PostID.String())

		comments = append(comments, &model.Comment{
			ID:      commentDB.ID.String(),
			Post:    post,
			Parent:  parentComment,
			Content: commentDB.Content,
		})
	}

	return comments, nil
}

func (r *DatabaseResolver) GetComment(id string) (*model.Comment, error) {
	var commentDB struct {
		ID       uuid.UUID `db:"id"`
		PostID   uuid.UUID `db:"post_id"`
		ParentID uuid.UUID `db:"parent_id"`
		Content  string    `db:"content"`
	}
	err := r.DB.Get(&commentDB, "SELECT * FROM comments WHERE id=$1", idToUUID(id))
	if err != nil {
		return nil, err
	}

	// Fetch the parent comment if it exists
	var parentComment *model.Comment
	if commentDB.ParentID != uuid.Nil {
		parentComment, _ = r.GetComment(commentDB.ParentID.String())
	}

	// Fetch the post related to the comment
	post, _ := r.GetPost(commentDB.PostID.String())

	comment := &model.Comment{
		ID:      commentDB.ID.String(),
		Post:    post,
		Parent:  parentComment,
		Content: commentDB.Content,
	}

	return comment, nil
}

func idToUUID(inputID string) uuid.UUID {
	id, err := uuid.Parse(inputID)
	if err != nil {
		log.Fatal("Error converting string to UUID", err, "in id", inputID)
	}
	return id
}
