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

	// 构造请求数据
	data := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": content,
			},
		},
	}

	// 序列化请求数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		response.SendInternalServerError(c, "Error marshalling data")
		return
	}

	// 创建请求
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		response.SendInternalServerError(c, "Error creating request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		response.SendBadRequestError(c, "Error sending request")
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.SendBadRequestError(c, "Error reading response body")
		return
	}
	bodyString := string(body)
	// 提取GPT响应内容
	content1, err := response.ExtractContentFromGPTResponse(bodyString)
	if err != nil {
		response.SendBadRequestError(c, "Error extracting content from GPT response")
		return
	}

	// 构造响应
	sessionResponse := SessionResponse{
		Content: content1,
	}

	// 发送响应
	response.SendSuccess(c, sessionResponse, "GPT API response")
}
