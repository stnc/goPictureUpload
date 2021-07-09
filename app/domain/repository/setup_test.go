package repository

import (
	"log"
	"os"
	"stncCms/app/domain/entity"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func DBConn() (*gorm.DB, error) {
	if _, err := os.Stat("./../../../.env"); !os.IsNotExist(err) {
		var err error
		err = godotenv.Load(os.ExpandEnv("./../../../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
		return LocalDatabase()
	}
	return CIBuild()
}

//Circle CI DB
func CIBuild() (*gorm.DB, error) {
	/*
		var err error
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", "127.0.0.1", "5432", "steven", "food-app-test", "password")
		conn, err := gorm.Open("postgres", DBURL)
		if err != nil {
			log.Fatal("This is the error:", err)
		}
		return conn, nil
	*/

	var err error
	dbdriver := os.Getenv("TEST_DB_DRIVER")
	host := os.Getenv("TEST_DB_HOST")
	password := os.Getenv("TEST_DB_PASSWORD")
	user := os.Getenv("TEST_DB_USER")
	dbname := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")

	DBURL := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	conn, err := gorm.Open(dbdriver, DBURL)
	if err != nil {
		log.Fatal("This is the error:", err)
	}
	return conn, nil
}

//Local DB
func LocalDatabase() (*gorm.DB, error) {
	dbdriver := os.Getenv("TEST_DB_DRIVER")
	host := os.Getenv("TEST_DB_HOST")
	password := os.Getenv("TEST_DB_PASSWORD")
	user := os.Getenv("TEST_DB_USER")
	dbname := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")

	//	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password) //postgresql
	DBURL := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	conn, err := gorm.Open(dbdriver, DBURL)
	if err != nil {
		return nil, err
	} else {
		log.Println("CONNECTED TO: ", dbdriver)
	}

	err = conn.DropTableIfExists(&entity.User{}, &entity.Post{}).Error
	if err != nil {
		return nil, err
	}
	err = conn.Debug().AutoMigrate(
		entity.User{},
		entity.Post{},
	).Error
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func seedUser(db *gorm.DB) (*entity.User, error) {
	user := &entity.User{
		ID:        1,
		FirstName: "vic",
		LastName:  "stev",
		Email:     "steven@example.com",
		Password:  "password",
		DeletedAt: nil,
	}
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func seedUsers(db *gorm.DB) ([]entity.User, error) {
	users := []entity.User{
		{
			ID:        1,
			FirstName: "vic",
			LastName:  "stev",
			Email:     "steven@example.com",
			Password:  "password",
			DeletedAt: nil,
		},
		{
			ID:        2,
			FirstName: "kobe",
			LastName:  "bryant",
			Email:     "kobe@example.com",
			Password:  "password",
			DeletedAt: nil,
		},
	}
	for _, v := range users {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func seedPost(db *gorm.DB) (*entity.Post, error) {
	food := &entity.Post{
		ID:          1,
		PostTitle:   "post title",
		PostContent: "post content",
		UserID:      1, //stncsession.GetUserID2(c)
	}
	err := db.Create(&food).Error
	if err != nil {
		return nil, err
	}
	return food, nil
}

func seedPosts(db *gorm.DB) ([]entity.Post, error) {
	posts := []entity.Post{
		{
			ID:          1,
			PostTitle:   "first post",
			PostContent: "first contn",
			UserID:      1, //stncsession.GetUserID2(c)
		},
		{
			ID:          2,
			PostTitle:   "second post",
			PostContent: "second content",
			UserID:      1, //stncsession.GetUserID2(c)
		},
	}
	for _, v := range posts {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}
