package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var loginAttempts = make(map[string]time.Time)

func CheckBruteForce(c *gin.Context) {
	ip := c.ClientIP()

	if attempt, ok := loginAttempts[ip]; ok && time.Now().Sub(attempt) < 1*time.Minute {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "小伙子你滴操作太频繁了，你是超导？"})
		c.Abort()
		return
	}

	c.Next()
}

func RecordAttempt(c *gin.Context) {
	ip := c.ClientIP()
	loginAttempts[ip] = time.Now()
}
