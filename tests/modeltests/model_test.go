package modeltests

import (
	"github.com/Funskie/blogIris/api/controllers"
	"github.com/Funskie/blogIris/api/models"

	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server controllers.Server
var userInstance models.User
var postInstance models.Post

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")
	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect %s database\n", TestDbDriver)
			log.Fatal("This is an error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}

	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect %s database\n", TestDbDriver)
			log.Fatal("This is an error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}

	log.Println("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {
	err := refreshUserTable()
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err = server.DB.Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {
	users := []models.User{
		models.User{
			Nickname: "Funskie 1",
			Email:    "tusty9292@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Funskie 2",
			Email:    "chiii57@gmail.com",
			Password: "password",
		},
	}

	for i := range users {
		err := server.DB.Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndPostTable() error {
	err := server.DB.DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}

	log.Println("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOnePost() (models.User, models.Post, error) {
	err := refreshUserAndPostTable()
	if err != nil {
		return models.User{}, models.Post{}, err
	}

	user := models.User{
		Nickname: "Pet",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err = server.DB.Create(&user).Error
	if err != nil {
		return models.User{}, models.Post{}, err
	}

	post := models.Post{
		Title:    "This is the title pet",
		Content:  "This is the content pet",
		AuthorID: user.ID,
	}

	err = server.DB.Create(&post).Error
	if err != nil {
		return models.User{}, models.Post{}, err
	}
	return user, post, nil
}

func seedUsersAndPosts() ([]models.User, []models.Post, error) {
	users := []models.User{
		models.User{
			Nickname: "Funskie 1",
			Email:    "tusty9292@gmail.com",
			Password: "password",
		},
		models.User{
			Nickname: "Funskie 2",
			Email:    "chiii57@gmail.com",
			Password: "password",
		},
	}
	posts := []models.Post{
		models.Post{
			Title:   "Title 1",
			Content: "Hello world 1",
		},
		models.Post{
			Title:   "Title 2",
			Content: "Hello world 2",
		},
	}

	for i := range users {
		err := server.DB.Create(&users[i]).Error
		if err != nil {
			return []models.User{}, []models.Post{}, err
		}

		posts[i].AuthorID = users[i].ID
		err = server.DB.Create(&posts[i]).Error
		if err != nil {
			return []models.User{}, []models.Post{}, err
		}
	}
	return users, posts, nil
}
