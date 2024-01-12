package handler

import (
	response "aTalkBackEnd/pkg"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

// var endpoint = os.Getenv("GPT_ENDPOINT")
var bearerToken = os.Getenv("GPT_BEARER_TOKEN")
var endpoint = "https://aigptx.top/v1/chat/completions"

func GPTHandler(c *gin.Context) {
	// 从请求中解析content字段
	var requestData map[string]string
	if err := c.BindJSON(&requestData); err != nil {
		response.SendBadRequestError(c, err.Error())
		return
	}
	content, ok := requestData["content"]
	if !ok {
		response.SendInternalServerError(c, "content field is required")
		return
	}

	data := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": content,
			},
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		response.SendInternalServerError(c, "Error marshalling data")
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		response.SendInternalServerError(c, "Error creating request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		response.SendBadRequestError(c, "Error sending request")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.SendBadRequestError(c, "Error reading response body")
		return
	}

	// 直接返回GPT API的响应
	response.SendSuccess(c, string(body), "GPT API response")
}
