package controller

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

var db *sql.DB
var err error

func initDB() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB_ACCOUNT := os.Getenv("DB_ACCOUNT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")

	Connection := DB_ACCOUNT + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ")/" + DB_NAME
	sql.Open("mysql", Connection)
	if err != nil {
		panic(err.Error())
	}
}

func SetupRouter() *gin.Engine {
	initDB()
	r := gin.Default()
	r.GET("/join", func(c *gin.Context) {
		var u User
		err := c.Bind(&u)
		if err != nil {
			log.Fatal(err)
		}

		Id, err := u.createUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(Id)
		username := u.Username
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", username),
		})
	})
	return r
}

func (u User) createUser() (Id int, err error) {
	stmt, err := db.Prepare("INSERT INTO  users(username, password) VALUES (?, ?)")
	if err != nil {
		return
	}

	rs, err := stmt.Exec(u.Username, u.Password)
	if err != nil {
		return
	}

	id, err := rs.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}
	Id = int(id)
	defer stmt.Close()
	return
}

func (u User) get() (user User, err error) {
	row := db.QueryRow("SELECT id, username FROM user WHERE username=?", u.Username)
	err = row.Scan(&user.Id, &user.Username, &user.CreatedTime)
	if err != nil {
		return
	}
	return
}

type User struct {
	Id          int
	Username    string    `json:"id" form:"username"`
	Password    string    `json:"Password" form:"password"`
	CreatedTime time.Time `json:"CreatedTime" form:"CreatedTime"`
}
