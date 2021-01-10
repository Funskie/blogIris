package modeltests

import (
	"github.com/Funskie/blogIris/api/models"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres driver
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestFindAllUsers(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	usersFound, err := userInstance.FindAllUsers(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, len(*usersFound), 2)
}

func TestSaveUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		ID:       1,
		Nickname: "test_save_man",
		Email:    "test@ttt.com",
		Password: "test",
	}
	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("this is the error saving the user: %v\n", err)
		return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Nickname, savedUser.Nickname)
	assert.Equal(t, newUser.Email, savedUser.Email)
}

func TestGetUserByID(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	foundUser, err := userInstance.FindUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error getting the user by ID: %v\n", err)
		return
	}

	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Nickname, user.Nickname)
}

func TestUpdatUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	updateUser := models.User{
		ID:       1,
		Email:    "updated@ttt.com",
		Nickname: "UpdatedUser",
		Password: "new password",
	}
	updatedUser, err := updateUser.UpdateAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	assert.Equal(t, updatedUser.ID, updateUser.ID)
	assert.Equal(t, updatedUser.Email, updateUser.Email)
	assert.Equal(t, updatedUser.Nickname, updateUser.Nickname)
}

func TestDeleteUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	isDelete, err := userInstance.DeleteAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error deleting the user %v\n", err)
		return
	}

	assert.Equal(t, isDelete, int64(1))
}
