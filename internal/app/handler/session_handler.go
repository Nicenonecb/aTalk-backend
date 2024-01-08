package handler

import (
	"aTalkBackEnd/internal/app/model"
	"aTalkBackEnd/internal/app/service"
	response "aTalkBackEnd/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type SessionHandler struct {
	Service *service.SessionService
}

//func (h *SessionHandler) ListSessions(c *gin.Context) {
//	sessions, err := h.Service.ListAllSessions()
//	if err != nil {
//		response.SendBadRequestError(c, err.Error())
//		return
//	}
//	response.SendSuccess(c, sessions, "Sessions retrieved")
//
//}

func (h *SessionHandler) ListUserSessions(c *gin.Context) {
	userIDInterface, exists := c.Get("userName") // 确保这里的 "userName" 是正确的上下文键
	if !exists {
		response.SendInternalServerError(c, "User ID not found in context")
		return
	}
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		response.SendInternalServerError(c, "User ID has an unexpected type")
		return
	}

	sessions, err := h.Service.ListAllSessionsByUserID(userID)
	if err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}

	response.SendSuccess(c, sessions, "Sessions retrieved")
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
	var session model.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}
	userIDInterface, exists := c.Get("userName")
	if !exists {
		response.SendInternalServerError(c, "user ID not found in context")
		return
	}
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		response.SendInternalServerError(c, "user ID has an unexpected type")
		return
	}
	session.UserID = userID

	if err := h.Service.CreateSession(&session); err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}
	response.SendSuccess(c, session, "Session created")
}

func (h *SessionHandler) DeleteSession(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.SendBadRequestError(c, "Invalid session ID")
		return
	}
	id := uint(id64)
	if err := h.Service.DeleteSession(id); err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}
	response.SendSuccess(c, nil, "Session deleted")
}

func (h *SessionHandler) UpdateSession(c *gin.Context) {
	var session model.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.UpdateSession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

type SessionResponse struct {
	Content       string `json:"content"`
	ContentBinary []byte `json:"content_binary"`
	// ...后面补充
}

func (h *SessionHandler) GetSessionDetails(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.SendBadRequestError(c, "Invalid session ID")
		return
	}
	id := uint(id64)

	// 获取 session 数据
	session, err := h.Service.GetSessionByID(id)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}
	// 拼接 Language, Scene, Detail
	concatenatedDetails := fmt.Sprintf("你现在是%s 大师, 我想和您进行语言文字上的对话，场景是%s,具体细节为%s,请使用%s语言回复,不用回复好的等客套话，直接进入角色，开始第一句话", session.Language, session.Scene, session.Detail, session.Language)

	// 调用 GPT
	gptResponse, err := response.CallGPT(concatenatedDetails)
	if err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}
	fmt.Println("GPT response:", gptResponse)
	content, err := response.ExtractContentFromGPTResponse(gptResponse)
	if err != nil {
		fmt.Println("Error11:", err)
		return
	}
	audioData, err := response.CallText2Speech(content)
	if err != nil {
		fmt.Println("Error22:", err)
		return
	}
	sessionResponse := SessionResponse{
		Content:       content,
		ContentBinary: audioData,
	}
	// 返回GPT的回复
	response.SendSuccess(c, sessionResponse, "GPT response")
}
