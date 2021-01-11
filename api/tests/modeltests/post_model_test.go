package modeltests

import (
	"github.com/Funskie/blogIris/api/models"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestFindAllPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error seeding user and post %v\n", err)
	}

	postsFound, err := postInstance.FindAllPosts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting posts %v\n", err)
		return
	}

	assert.Equal(t, len(*postsFound), 2)
}

func TestSavePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Error seeding user %v\n", err)
	}

	post := models.Post{
		ID:       1,
		Title:    "Testing save post title",
		Content:  "tttttttttttttttttttttttttt",
		AuthorID: user.ID,
	}

	savedPost, err := post.SavePost(server.DB)
	if err != nil {
		t.Errorf("this is the error saving post %v\n", err)
		return
	}

	assert.Equal(t, savedPost.ID, post.ID)
	assert.Equal(t, savedPost.Title, post.Title)
	assert.Equal(t, savedPost.Content, post.Content)
	assert.Equal(t, savedPost.AuthorID, post.AuthorID)
}

func TestGetPostByID(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error seeding user and post table %v\n", err)
	}

	postFound, err := postInstance.FindPostByID(server.DB, post.ID)
	if err != nil {
		t.Errorf("this is the error getting post by ID %v\n", err)
		return
	}

	assert.Equal(t, postFound.ID, post.ID)
	assert.Equal(t, postFound.Title, post.Title)
	assert.Equal(t, postFound.Content, post.Content)
	assert.Equal(t, postFound.AuthorID, post.AuthorID)
}

func TestUpdatePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error seeding user and post table %v\n", err)
	}

	postUpdate := models.Post{
		ID:       1,
		Title:    "Test update post title",
		Content:  "new ttttttttttt",
		AuthorID: post.AuthorID,
	}

	updatedPost, err := postUpdate.UpdateAPost(server.DB, postUpdate.ID)
	if err != nil {
		t.Errorf("this is the error updating post %v\n", err)
		return
	}

	assert.Equal(t, postUpdate.ID, updatedPost.ID)
	assert.Equal(t, postUpdate.Title, updatedPost.Title)
	assert.Equal(t, postUpdate.Content, updatedPost.Content)
	assert.Equal(t, postUpdate.AuthorID, updatedPost.AuthorID)
}

func TestDeletePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error seeding user and post table %v\n", err)
	}

	isDelete, err := postInstance.DeleteAPost(server.DB, post.ID, post.AuthorID)
	if err != nil {
		t.Errorf("this is the error delete post %v\n", err)
		return
	}

	assert.Equal(t, isDelete, int64(1))
}
