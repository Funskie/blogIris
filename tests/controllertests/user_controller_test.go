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

func TestCreateUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshing user table %v\n", err)
	}

	samples := []struct {
		inputJSON    string
		statusCode   int
		nickname     string
		email        string
		errorMessage string
	}{
		{
			inputJSON:    `{"nickname":"Pet", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   201,
			nickname:     "Pet",
			email:        "pet@gmail.com",
			errorMessage: "",
		},
		{
			inputJSON:    `{"nickname":"Frank", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Email Already Taken",
		},
		{
			inputJSON:    `{"nickname":"Pet", "email": "frank@gmail.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: "Nickname Already Taken",
		},
		{
			inputJSON:    `{"nickname":"Pet", "email": "petgmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"nickname":"", "email": "pet@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Nickname",
		},
		{
			inputJSON:    `{"nickname":"Pet", "email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"nickname":"Pet", "email": "pet@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("this is the error convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["nickname"], v.nickname)
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshing user table %v\n", err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatalf("Error seeding users %v\n", err)
	}

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)

	var users []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		t.Errorf("this is the error convert to json: %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUserByID(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshing user table %v\n", err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Error seeding user %v\n", err)
	}

	samples := []struct {
		id           string
		statusCode   int
		nickname     string
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			nickname:   user.Nickname,
			email:      user.Email,
		},
		{
			id:         "unknow",
			statusCode: 400,
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("this is the error convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["nickname"], v.nickname)
			assert.Equal(t, responseMap["email"], v.email)
		}
	}
}

func TestUpdateUser(t *testing.T) {

	var authEmail, authPassword string
	var authID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshing user table %v\n", err)
	}

	users, err := seedUsers() // Need at least 2 uers to check this func
	if err != nil {
		log.Fatalf("Error seeding users %v\n", err)
	}
	// Get the first user
	for _, v := range users {
		if v.ID == 2 {
			continue
		}
		authID = v.ID
		authEmail = v.Email
		authPassword = "password"
	}

	token, err := server.SignIn(authEmail, authPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id             string
		updateJSON     string
		statusCode     int
		updateNickname string
		updateEmail    string
		tokenGiven     string
		errorMessage   string
	}{
		{
			id:             strconv.Itoa(int(authID)),
			updateJSON:     `{"nickname":"Chi", "email": "chi57@gmail.com", "password": "password"}`,
			statusCode:     200,
			updateNickname: "Chi",
			updateEmail:    "chi57@gmail.com",
			tokenGiven:     tokenString,
			errorMessage:   "",
		},
		{
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"nickname":"Chi", "email": "chi57gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Invalid Email",
		},
		{
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"nickname":"", "email": "chi57@gmail.com", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Nickname",
		},
		{
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"nickname":"Chi", "email": "", "password": "password"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Email",
		},
		{
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"nickname":"Chi", "email": "chi57@gmail.com", "password": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Password",
		},
		{
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"nickname":"Chi", "email": "chi57@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"nickname":"Chi", "email": "chi57@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   "This is error token",
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"nickname":"Funskie 2", "email": "chi57@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Nickname Already Taken",
		},
		{
			id:           strconv.Itoa(int(authID)),
			updateJSON:   `{"nickname":"Chi", "email": "chiii57@gmail.com", "password": "password"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Email Already Taken",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(2)),
			updateJSON:   `{"nickname":"Chi", "email": "chi57@gmail.com", "password": "password"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("PUT", "/users", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUser)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("this is the error convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["nickname"], v.updateNickname)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {

	var authEmail, authPassword string
	var authID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatalf("Error refreshing user table %v\n", err)
	}

	users, err := seedUsers() // Need at least 2 uers to check this func
	if err != nil {
		log.Fatalf("Error seeding users %v\n", err)
	}

	// Get the first user
	for _, v := range users {
		if v.ID == 2 {
			continue
		}
		authID = v.ID
		authEmail = v.Email
		authPassword = "password"
	}

	token, err := server.SignIn(authEmail, authPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(authID)),
			tokenGiven:   "This is incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(authID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(authID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("DELETE", "/users", nil)
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUser)
		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 204 {
			assert.Equal(t, rr.Body.String()[:1], "1")
		}

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
