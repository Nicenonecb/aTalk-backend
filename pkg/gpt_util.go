package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var endpoint = os.Getenv("GPT_ENDPOINT")
var bearerToken = os.Getenv("GPT_BEARER_TOKEN")

func CallGPT(content string) (string, error) {
	data := map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    []map[string]string{{"role": "user", "content": content}},
		"temperature": 0.7,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ExtractContentFromGPTResponse(response string) (string, error) {
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(response), &data); err != nil {
		return "", err
	}

	choices, ok := data["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("choices field not found or empty")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("error parsing choice")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("error parsing message")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("error parsing content")
	}

	return content, nil
}
