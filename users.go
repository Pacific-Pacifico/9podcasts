package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/apnishiksha/9podcasts/auth"
	"github.com/labstack/echo"
)

type user struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func createUsers(c echo.Context) (err error) {

	us := &user{}

	if err = c.Bind(us); err != nil {
		return
	}

	hashedPassword, err := auth.GetHashedPassword(us.Password)
	log.Println("hashedPassword= ", hashedPassword)
	if err != nil {
		log.Println("Error generating hashedPassword :", err)
		return
	}

	us.Created_at = time.Now()
	us.Updated_at = time.Now()

	q := `INSERT INTO users(name,email,password,created_at,updated_at) VALUES($1,$2,$3,$4,$5) RETURNING id`
	err = db.QueryRow(q, us.Name, us.Email, hashedPassword, us.Created_at, us.Updated_at).Scan(&us.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	return c.JSON(http.StatusOK, us)
}

func getUsers(c echo.Context) (err error) {

	var usr []user

	q := `SELECT id,name,email,password,created_at,updated_at FROM users`

	rows, err := db.Query(q)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var us user

		err = rows.Scan(&us.ID, &us.Name, &us.Email, &us.Password, &us.Created_at, &us.Updated_at)
		if err != nil {
			fmt.Println(err)
		}

		usr = append(usr, us)
	}

	return c.JSON(http.StatusOK, usr)
}
