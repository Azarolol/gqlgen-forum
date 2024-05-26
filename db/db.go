package db

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/azarolol/gqlen-forum/graph/model"
	"github.com/go-pg/pg"
)

type DB interface {
	insertPost(*model.Post) (*model.Post, error)
	insertComment(*model.Comment) (*model.Comment, error)
	getAllPosts() ([]*model.Post, error)
	selectPostWithID(string) (*model.Post, error)
	selectCommentWithID(string) (*model.Comment, error)
	getCommentsOnPost(string) ([]*model.Comment, error)
	getCommentsOnComment(string) ([]*model.Comment, error)
}

type pgDB struct {
	db *pg.DB
}

type localDB struct {
	posts     map[int]*model.Post
	comments  map[int]*model.Comment
	postID    int
	commentID int
}

func InsertPost(db DB, input *model.Post) (*model.Post, error) {
	return db.insertPost(input)
}

func InsertComment(db DB, input *model.Comment) (*model.Comment, error) {
	return db.insertComment(input)
}

func GetAllPosts(db DB) ([]*model.Post, error) {
	return db.getAllPosts()
}

func SelectPostWithID(db DB, id string) (*model.Post, error) {
	return db.selectPostWithID(id)
}

func SelectCommentWithID(db DB, id string) (*model.Comment, error) {
	return db.selectCommentWithID(id)
}

func GetCommentsOnPost(db DB, id string) ([]*model.Comment, error) {
	return db.getCommentsOnPost(id)
}

func GetCommentsOnComment(db DB, id string) ([]*model.Comment, error) {
	return db.getCommentsOnComment(id)
}

func GetCommentWithComments(db DB, input string) (*model.CommentWithComments, error) {
	comment, err := db.selectCommentWithID(input)
	if err != nil {
		return nil, fmt.Errorf("error getting comment: %v", err)
	}
	comments, err := db.getCommentsOnComment(input)
	if err != nil {
		return nil, fmt.Errorf("error getting comments on comment: %v", err)
	}
	var commentsWithComments []*model.CommentWithComments
	for _, comment := range comments {
		commentWithComment, err := GetCommentWithComments(db, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting comments on comment: %v", err)
		}
		commentsWithComments = append(commentsWithComments, commentWithComment)
	}
	output := &model.CommentWithComments{
		Comment:  comment,
		Comments: commentsWithComments,
	}
	return output, nil
}

func (pg *pgDB) insertPost(post *model.Post) (*model.Post, error) {
	_, err := pg.db.Model(post).Insert()
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (pg *pgDB) insertComment(comment *model.Comment) (*model.Comment, error) {
	_, err := pg.db.Model(comment).Insert()
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (pg *pgDB) getAllPosts() ([]*model.Post, error) {
	var posts []*model.Post

	err := pg.db.Model(&posts).Select()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (pg *pgDB) getCommentsOnPost(input string) ([]*model.Comment, error) {
	var comments []*model.Comment

	err := pg.db.Model(&comments).Where("on_post = true AND comment_on = " + input).Select()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (pg *pgDB) getCommentsOnComment(input string) ([]*model.Comment, error) {
	var comments []*model.Comment

	err := pg.db.Model(&comments).Where("on_post = false AND comment_on = " + input).Select()
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (pg *pgDB) selectPostWithID(id string) (*model.Post, error) {
	var post model.Post
	err := pg.db.Model(&post).Where("id = " + id).Select()
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (pg *pgDB) selectCommentWithID(id string) (*model.Comment, error) {
	var comment model.Comment
	err := pg.db.Model(&comment).Where("id = " + id).Select()
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (ldb *localDB) insertPost(post *model.Post) (*model.Post, error) {
	id := ldb.postID
	post.ID = strconv.Itoa(id)
	ldb.posts[id] = post
	ldb.postID++
	return post, nil
}

func (ldb *localDB) insertComment(comment *model.Comment) (*model.Comment, error) {
	id := ldb.commentID
	comment.ID = strconv.Itoa(id)
	ldb.comments[id] = comment
	ldb.commentID++
	return comment, nil
}

func (ldb *localDB) getAllPosts() ([]*model.Post, error) {
	var (
		keys  []int
		posts []*model.Post
	)
	for key := range ldb.posts {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for key := range keys {
		posts = append(posts, ldb.posts[key])
	}
	return posts, nil
}

func (ldb *localDB) getCommentsOnPost(input string) ([]*model.Comment, error) {
	var (
		keys     []int
		comments []*model.Comment
	)
	for key, comment := range ldb.comments {
		if comment.OnPost && comment.CommentOn == input {
			keys = append(keys, key)
		}
	}
	for key := range keys {
		comments = append(comments, ldb.comments[key])
	}
	return comments, nil
}

func (ldb *localDB) getCommentsOnComment(input string) ([]*model.Comment, error) {
	var (
		keys     []int
		comments []*model.Comment
	)
	for key, comment := range ldb.comments {
		if !comment.OnPost && comment.CommentOn == input {
			keys = append(keys, key)
		}
	}
	for key := range keys {
		comments = append(comments, ldb.comments[key])
	}
	return comments, nil
}

func (ldb *localDB) selectPostWithID(input string) (*model.Post, error) {
	id, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("error converting id: %v", err)
	}
	post, ok := ldb.posts[id]
	if !ok {
		return nil, fmt.Errorf("no post with such id")
	}
	return post, nil
}

func (ldb *localDB) selectCommentWithID(input string) (*model.Comment, error) {
	id, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("error converting id: %v", err)
	}
	comment, ok := ldb.comments[id]
	if !ok {
		return nil, fmt.Errorf("no comment with such id")
	}
	return comment, nil
}

func Connect(opts pg.Options) *pgDB {
	db := pg.Connect(&opts)
	if _, DBStatus := db.Exec("SELECT 1"); DBStatus != nil {
		panic("PostgreSQL is down")
	}
	return &pgDB{db}
}

func CreateLocalDB() *localDB {
	db := localDB{
		posts:     make(map[int]*model.Post),
		comments:  make(map[int]*model.Comment),
		postID:    1,
		commentID: 1,
	}
	return &db
}
