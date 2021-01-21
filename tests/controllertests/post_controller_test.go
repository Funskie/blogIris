package controllertests

import (
	"github.com/Funskie/blogIris/api/models"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreatePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Error seeding user %v\n", err)
	}

	token, err := server.SignIn(user.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title": "Test create title", "content": "Test create content", "author_id": 1}`,
			statusCode:   201,
			title:        "Test create title",
			content:      "Test create content",
			author_id:    user.ID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			inputJSON:    `{"title": "Test create title", "content": "Test create content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			inputJSON:    `{"title": "", "content": "Test create content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			inputJSON:    `{"title": "Test create title 2", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			inputJSON:    `{"title": "Test create title 2", "content": "Test create content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			inputJSON:    `{"title": "Test create title 2", "content": "Test create content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "Test create title 2", "content": "Test create content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "Test create title 2", "content": "Test create content", "author_id": 5}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreatePost)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("this is the error convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetPosts(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error seeding users and posts %v\n", err)
	}

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetPosts)
	handler.ServeHTTP(rr, req)

	var posts []models.Post
	err = json.Unmarshal([]byte(rr.Body.String()), &posts)
	if err != nil {
		t.Errorf("this is the error convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(posts), 2)
}

func TestGetPostByID(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	_, post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatalf("Error seeding users and posts %v\n", err)
	}

	samples := []struct {
		id           string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(post.ID)),
			statusCode: 200,
			title:      post.Title,
			content:    post.Content,
			author_id:  post.AuthorID,
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("GET", "/posts", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetPost)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("this is the error convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], post.Title)
			assert.Equal(t, responseMap["content"], post.Content)
			assert.Equal(t, responseMap["author_id"], float64(post.AuthorID))
		}
	}
}

func TestUpdatePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error seeding users and posts %v\n", err)
	}

	token, err := server.SignIn(users[0].Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			updateJSON:   `{"title": "Test update title", "content": "Test update content", "author_id": 1}`,
			statusCode:   200,
			title:        "Test update title",
			content:      "Test update content",
			author_id:    users[0].ID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			updateJSON:   `{"title": "Title 2", "content": "Test update content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			updateJSON:   `{"title": "", "content": "Test update content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			updateJSON:   `{"title": "Test update title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			updateJSON:   `{"title": "Test update title", "content": "Test update content"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			updateJSON:   `{"title": "Test update title", "content": "Test update content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			updateJSON:   `{"title": "Test update title", "content": "Test update content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			updateJSON:   `{"title": "Test update title", "content": "Test update content", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("PUT", "/posts", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdatePost)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("this is the error convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeletePost(t *testing.T) {

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}

	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatalf("Error seeding users and posts %v\n", err)
	}

	token, err := server.SignIn(users[0].Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			statusCode:   401,
			author_id:    users[0].ID,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			statusCode:   401,
			author_id:    users[0].ID,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			statusCode: 400,
			author_id:  users[0].ID,
			tokenGiven: tokenString,
		},
		{
			id:           strconv.Itoa(int(posts[1].ID)),
			statusCode:   401,
			author_id:    users[1].ID,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(posts[0].ID)),
			statusCode:   204,
			author_id:    users[0].ID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("DELETE", "/posts", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeletePost)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 204 {
			assert.Equal(t, rr.Body.String()[:1], "1")
		}
		if v.statusCode == 401 || v.statusCode == 400 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
