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
	config := cors.DefaultConfig()

	config.AllowOrigins = []string{"http://127.0.0.1:8081", "http://localhost:8081", "https://aigptx.top"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	r := gin.Default()
	r.Use(middleware.CheckBruteForce)
	r.Use(cors.New(config))
	v1 := r.Group("/v1")

	sessionService := &service.SessionService{
		Repo: &repository.SessionRepository{DB: db},
	}

	sessionHandler := &handler.SessionHandler{Service: sessionService}

	v1.POST("/gpt-response", handler.GPTHandler)

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1.POST("/register", userHandler.Register)
	v1.POST("/login", userHandler.Login)
	v1.GET("/sessions/:id/details", middleware.TokenAuthMiddleware(), sessionHandler.GetSessionDetails)
	v1.GET("/sessions", middleware.TokenAuthMiddleware(), sessionHandler.ListUserSessions)
	v1.POST("/sessions", middleware.TokenAuthMiddleware(), sessionHandler.CreateSession)
	v1.DELETE("/sessions/:id", middleware.TokenAuthMiddleware(), sessionHandler.DeleteSession)
	v1.PUT("/sessions/:id", middleware.TokenAuthMiddleware(), sessionHandler.UpdateSession)
	v1.POST("/dialogue", middleware.TokenAuthMiddleware(), handler.DialogueHandler)
	v1.POST("/upload", middleware.TokenAuthMiddleware(), handler.UploadToGCSHandler)
	v1.POST("/speech2text", middleware.TokenAuthMiddleware(), handler.Speech2TextHandler)
	v1.POST("/text2speech", middleware.TokenAuthMiddleware(), handler.Text2SpeechHandler)

	return r
}
