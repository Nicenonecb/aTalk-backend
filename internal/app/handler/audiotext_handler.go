package handler

import (
	response "aTalkBackEnd/pkg"
	speech2text "aTalkBackEnd/pkg"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SpeechRequest struct {
	ContentType string `json:"contentType"`
	AudioData   string `json:"audioData"`
}

func Speech2TextHandler(c *gin.Context) {
	var req SpeechRequest

	if err := c.ShouldBind(&req); err != nil {
		response.SendBadRequestError(c, "Invalid request parameters")
		return
	}

	audioData, err := base64.StdEncoding.DecodeString(req.AudioData)
	if err != nil {
		fmt.Printf("Error decoding audio data: %v\n", err) // Log the error
		response.SendBadRequestError(c, "Failed to decode audio data")
		return
	}

	fmt.Println("req:", req)
	transcription, err := speech2text.CallSpeech2Text(audioData)
	if err != nil {
		fmt.Printf("Error transcribing audio: %v\n", err) // Log the error
		response.SendBadRequestError(c, "Failed to transcribe audio")
		return
	}

	response.SendSuccess(c, gin.H{"transcription": transcription}, "Audio transcribed")
}

type Text2SpeechRequest struct {
	Content string `json:"content" binding:"required"`
}

func Text2SpeechHandler(c *gin.Context) {
	var req Text2SpeechRequest
	if err := c.BindJSON(&req); err != nil {
		response.SendBadRequestError(c, "Invalid request body")
		return
	}

	fmt.Println("req:", req)
	audioData, err := response.CallText2Speech(req.Content)
	if err != nil {
		response.SendBadRequestError(c, "Failed to generate audio")
		return
	}

	// Return the audio data as binary to the client
	response.SendSuccess(c, gin.H{"audioData": audioData}, "Audio generated")
}
