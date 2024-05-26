package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"fmt"

	"github.com/azarolol/gqlen-forum/db"
	"github.com/azarolol/gqlen-forum/graph/model"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	if len(input.Author) > 100 {
		return nil, fmt.Errorf("author can't be more then 100 characters")
	}
	if len(input.Header) > 200 {
		return nil, fmt.Errorf("header can't be more than 2000 characters")
	}
	if len(input.Body) > 2000 {
		return nil, fmt.Errorf("post can't be more than 2000 characters")
	}
	post := &model.Post{
		Author:         input.Author,
		Header:         input.Header,
		Body:           input.Body,
		AllowsComments: input.AllowsComments,
	}

	post, err := db.InsertPost(r.DB, post)
	if err != nil {
		return nil, fmt.Errorf("error inserting new post: %v", err)
	}

	return post, nil
}

// AddCommentOnPost is the resolver for the addCommentOnPost field.
func (r *mutationResolver) AddCommentOnPost(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	post, err := db.SelectPostWithID(r.DB, input.CommentOn)
	if err != nil {
		return nil, fmt.Errorf("error finding post: %v", err)
	}
	if !post.AllowsComments {
		return nil, fmt.Errorf("author restricted leaving comments on this post")
	}
	if len(input.Author) > 100 {
		return nil, fmt.Errorf("author can't be more then 100 characters")
	}
	if len(input.Body) > 2000 {
		return nil, fmt.Errorf("comment can't be more than 2000 characters")
	}
	comment := &model.Comment{
		OnPost:    true,
		CommentOn: input.CommentOn,
		Author:    input.Author,
		Body:      input.Body,
	}

	comment, err = db.InsertComment(r.DB, comment)
	if err != nil {
		return nil, fmt.Errorf("error inserting new comment: %v", err)
	}

	return comment, nil
}

// AddCommentOnComment is the resolver for the addCommentOnComment field.
func (r *mutationResolver) AddCommentOnComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	_, err := db.SelectCommentWithID(r.DB, input.CommentOn)
	if err != nil {
		return nil, fmt.Errorf("error finding comment: %v", err)
	}
	if len(input.Author) > 100 {
		return nil, fmt.Errorf("author can't be more then 100 characters")
	}
	if len(input.Body) > 2000 {
		return nil, fmt.Errorf("comment can't be more than 2000 characters")
	}
	newComment := &model.Comment{
		OnPost:    false,
		CommentOn: input.CommentOn,
		Author:    input.Author,
		Body:      input.Body,
	}

	newComment, err = db.InsertComment(r.DB, newComment)
	if err != nil {
		return nil, fmt.Errorf("error inserting new comment: %v", err)
	}

	return newComment, nil
}

// GetPostsPage is the resolver for the getPostsPage field.
func (r *queryResolver) GetPostsPage(ctx context.Context, limit *int, offset *int) ([]*model.Post, error) {
	posts, err := db.GetAllPosts(r.DB)
	if err != nil {
		return nil, fmt.Errorf("error getting posts: %v", err)
	}

	start := *offset
	end := *limit + *offset

	if end > len(posts) {
		end = len(posts)
	}

	return posts[start:end], nil
}

// GetPostWithComments is the resolver for the getPostWithComments field.
func (r *queryResolver) GetPostWithComments(ctx context.Context, input string, limit *int, offset *int) (*model.PostWithComments, error) {
	post, err := db.SelectPostWithID(r.DB, input)
	if err != nil {
		return nil, fmt.Errorf("error getting post: %v", err)
	}
	comments, err := db.GetCommentsOnPost(r.DB, input)
	if err != nil {
		return nil, fmt.Errorf("error getting comments on post: %v", err)
	}
	var commentsWithComments []*model.CommentWithComments
	for _, comment := range comments {
		commentWithComments, err := db.GetCommentWithComments(r.DB, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting comment with comments: %v", err)
		}
		commentsWithComments = append(commentsWithComments, commentWithComments)
	}
	postWithComments := &model.PostWithComments{
		Post:     post,
		Comments: commentsWithComments,
	}
	return postWithComments, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.