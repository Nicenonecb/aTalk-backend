package router

import (
	"aTalkBackEnd/internal/app/handler"
	"aTalkBackEnd/internal/app/middleware"
	"aTalkBackEnd/internal/app/model"
	"aTalkBackEnd/internal/app/repository"
	"aTalkBackEnd/internal/app/service"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func getDatabaseConnection() (*gorm.DB, error) {
	_, onGoogleCloud := os.LookupEnv("GAE_ENV")

	var dsn string
	if onGoogleCloud {
		// On Google Cloud - connect using Unix socket
		dbUser := os.Getenv("DB_USER")
		dbPwd := os.Getenv("DB_PASS")
		dbName := os.Getenv("DB_NAME")
		unixSocketPath := os.Getenv("INSTANCE_UNIX_SOCKET")

		if dbUser == "" || dbPwd == "" || dbName == "" || unixSocketPath == "" {
			log.Fatal("Environment variables DB_USER, DB_PASS, DB_NAME, and INSTANCE_UNIX_SOCKET must be set.")
		}

		dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true&allowCleartextPasswords=1", dbUser, dbPwd, unixSocketPath, dbName)
	} else {
		dsn = os.Getenv("DB_CONNECTION_STRING")
		if dsn == "" {
			log.Fatal("Environment variable DB_CONNECTION_STRING must be set.")
		}
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupRoutes() *gin.Engine {
	db, err := getDatabaseConnection()
	if err != nil {
		panic(err)
	}

	userHandler := &handler.UserHandler{
		Service: &service.UserService{
			Repo: &repository.UserRepository{
				DB: db,
			},
		},
	}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Session{})

	r := gin.Default()
	r.Use(middleware.CheckBruteForce)
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://127.0.0.1:8081"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	sessionService := &service.SessionService{
		Repo: &repository.SessionRepository{DB: db},
	}

	sessionHandler := &handler.SessionHandler{Service: sessionService}

	r.POST("/gpt-response", handler.GPTHandler)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.GET("/sessions/:id/details", middleware.TokenAuthMiddleware(), sessionHandler.GetSessionDetails)
	r.GET("/sessions", middleware.TokenAuthMiddleware(), sessionHandler.ListSessions)
	r.POST("/sessions", middleware.TokenAuthMiddleware(), sessionHandler.CreateSession)
	r.DELETE("/sessions/:id", middleware.TokenAuthMiddleware(), sessionHandler.DeleteSession)
	r.PUT("/sessions/:id", middleware.TokenAuthMiddleware(), sessionHandler.UpdateSession)
	r.POST("/dialogue", middleware.TokenAuthMiddleware(), handler.DialogueHandler)
	r.POST("/upload", middleware.TokenAuthMiddleware(), handler.UploadToGCSHandler)
	r.POST("/speech2text", middleware.TokenAuthMiddleware(), handler.Speech2TextHandler)
	r.POST("/text2speech", middleware.TokenAuthMiddleware(), handler.Text2SpeechHandler)
	return r
}
