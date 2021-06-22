package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/labstack/echo"
)

var (
	db    *sql.DB
	err   error
	cache redis.Conn
	env   string = "debug"
)

func initCache() {
	// Initialize the redis connection to a redis instance running on your local machine
	cache, err = redis.DialURL("redis://localhost")
	if err != nil {
		log.Println("redis connection error")
		log.Fatal(err)
	}
}

func main() {

	// setup upload directory
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	if _, err := os.Stat("feeds"); os.IsNotExist(err) {
		os.Mkdir("feeds", 0755)
	}

	connectionString := os.Getenv("CONN")
	if connectionString == "" {
		log.Fatal("error getting connection string")
	}

	db, err = sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	initCache()

	api := echo.New()

	api.GET("/", rootHandler)

	api.POST("/uploads", uploadMedia)
	api.GET("/media/*", getMediaFile)

	api.POST("/categories", createCategory)
	api.GET("/categories", getCategories)
	api.DELETE("/categories/:id", deleteCategory)
	api.PUT("/categories", updateCategory)

	api.POST("/keywords", createKeywords)
	api.GET("/keywords", getKeywords)

	api.POST("/podcasts", createPodcast)
	api.GET("/podcasts", getPodcast)
	api.GET("/podcasts/:id", getPodcastByID) // path parameters
	api.DELETE("/podcasts/:id", deletePodcast)
	api.PUT("/podcasts", updatePodcast)

	api.POST("/episodes", createEpisodes)
	api.GET("/episodes", getEpisodes) // query parameters
	api.DELETE("/episodes/:id", deleteEpisode)
	api.PUT("/episodes", updateEpisode)

	api.POST("/register", createUser)
	api.GET("/profile", getUsers)
	api.DELETE("/users/:id", deleteUser)

	api.PATCH("/update-name", updateUserName)
	api.PATCH("/update-email", updateUserEmail)
	api.POST("/update-password", updateUserPassword)

	api.POST("/signin", signIn)

	api.POST("/dashboard", Dashboard)
	api.POST("/refresh", refreshToken)

	api.GET("/forgot", ForgotPassword)
	api.GET("/resetpass-request", ResetPasswordRequest)
	api.POST("/resetpass", ResetPassword)

	api.HideBanner = true
	api.Start(":9999")
}

func rootHandler(c echo.Context) (err error) {

	return c.String(http.StatusOK, "Welcome to 9Podcasts Api")
}
