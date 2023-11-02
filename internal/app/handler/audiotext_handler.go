package handler

import (
	response "aTalkBackEnd/pkg"
	speech2text "aTalkBackEnd/pkg"
	"github.com/gin-gonic/gin"
)

type SpeechRequest struct {
	ContentType string `form:"contentType"`
	AudioData   []byte `form:"audioData" binding:"required"`
}

func Speech2TextHandler(c *gin.Context) {
	var req SpeechRequest

	if err := c.ShouldBind(&req); err != nil {
		response.SendBadRequestError(c, "Invalid request parameters")
		return
	}

	transcription, err := speech2text.CallSpeech2Text(req.AudioData, req.ContentType)
	if err != nil {
		response.SendBadRequestError(c, "Failed to transcribe audio")
		return
	}

	response.SendSuccess(c, gin.H{"transcription": transcription}, "Audio transcribed")
}

type Text2SpeechRequest struct {
	Content     string `json:"content" binding:"required"`
	ContentType string `json:"contentType" binding:"required"`
}

func Text2SpeechHandler(c *gin.Context) {
	var req Text2SpeechRequest
	audioData, err := response.CallText2Speech(req.Content, req.ContentType)
	if err != nil {
		response.SendBadRequestError(c, "Failed to generate audio")
		return
	}

	// Return the audio data as binary to the client
	response.SendSuccess(c, gin.H{"audioData": audioData}, "Audio generated")
}
