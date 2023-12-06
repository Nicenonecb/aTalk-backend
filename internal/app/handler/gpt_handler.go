package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content, ok := requestData["content"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content field is required"})
		return
	}

	data := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": content,
			},
			{
				"role":    "system",
				"content": "You are a helpful assistant",
			},
		},
		"temperature": 0.7,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling data"})
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}

	// 直接返回GPT API的响应
	c.String(http.StatusOK, string(body))
}
