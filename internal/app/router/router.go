package router

import (
	"aTalkBackEnd/internal/app/handler"
	"aTalkBackEnd/internal/app/middleware"
	"aTalkBackEnd/internal/app/model"
	"aTalkBackEnd/internal/app/repository"
	"aTalkBackEnd/internal/app/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func SetupRoutes() *gin.Engine {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
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
