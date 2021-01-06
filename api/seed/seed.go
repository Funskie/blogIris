package seed

import (
	"github.com/Funskie/blogIris/api/models"

	"log"

	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Nickname: "Funskie",
		Email:    "tusty9292@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Wuskie",
		Email:    "chiii57@gmail.com",
		Password: "password",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	models.Post{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		u, err := users[i].SaveUser(db)
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		users[i] = *u
		posts[i].AuthorID = users[i].ID

		_, err = posts[i].SavePost(db)
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
