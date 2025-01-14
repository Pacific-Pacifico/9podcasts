package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func signIn(c echo.Context) (err error) {
	enteredCreds := &Credentials{}

	if err = c.Bind(enteredCreds); err != nil {
		return
	}

	result := db.QueryRow("SELECT password FROM users WHERE email=$1", enteredCreds.Email)
	// to store the credentials we get from the database
	storedCreds := &Credentials{}
	// Store the obtained password in storedCreds
	err = result.Scan(&storedCreds.Password)

	// If an entry with the Email does not exist, send an "Unauthorized"(401) status
	if err == sql.ErrNoRows {
		log.Println(err)
		return c.NoContent(http.StatusUnauthorized)
	}

	// If the error is of any other type, send a 500 status
	if err != nil {
		log.Println(err)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(enteredCreds.Password))
	if err != nil {
		// If the two passwords don't match, return a 401 status
		log.Println(err)
		return
	}

	SessionID, err := SetSessionID(enteredCreds.Email)
	if err != nil {
		log.Println(err)
		return
	}

	cookie := new(http.Cookie)
	cookie.Name = "SessionID"
	cookie.Value = SessionID
	cookie.Expires = time.Now().Add(60 * time.Second)
	c.SetCookie(cookie)

	return c.String(http.StatusOK, SessionID)
}
