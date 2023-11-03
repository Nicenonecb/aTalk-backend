package handler

import (
	"aTalkBackEnd/internal/app/model"
	"aTalkBackEnd/internal/app/service"
	response "aTalkBackEnd/pkg"
	utility "aTalkBackEnd/pkg"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func (h *UserHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}
	if !utility.IsValidEmail(user.Email) {
		response.SendBadRequestError(c, "Invalid email address")
		return
	}
	if err := h.Service.Register(&user); err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}
	response.SendSuccess(c, nil, "Registration successful")
}

func (h *UserHandler) Login(c *gin.Context) {
	var loginDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	isAuthenticated, err := h.Service.Authenticate(loginDetails.Username, loginDetails.Password)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	if !isAuthenticated {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
		response.StatusUnauthorized(c, "Invalid login credentials")
		return
	}

	user, err := h.Service.FindByUsername(loginDetails.Username)

	token, err := utility.GenerateToken(user.ID) // Assuming username can be used as userID
	if err != nil {

		response.SendInternalServerError(c, err.Error())
		return
	}
	response.SendSuccess(c, gin.H{"token": token}, "Login successful")
}

func (h *UserHandler) Logout(c *gin.Context) {
	// TODO：后面做toke的白名单
	response.SendSuccess(c, nil, "Logout successful, please delete the token from the client side")
}
