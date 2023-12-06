package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func CallSpeech2Text(userData []byte) (string, error) {
	var url = "https://aigptx.top/v1/audio/transcriptions"

	// Write the binary data to a .mp3 file
	err := os.WriteFile("output.mp3", userData, 0644)
	if err != nil {
		return "", err
	}

	// Open the .mp3 file
	file, err := os.Open("output.mp3")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the request
	fileWriter, err := writer.CreateFormFile("file", "output.mp3")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", err
	}

	// Add the other fields to the request
	_ = writer.WriteField("model", "whisper-1")

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return "", err
	}

	// Create the request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	type Response struct {
		Text string `json:"text"`
	}
	var respData Response
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return "", err
	}

	fmt.Println("Response body:", string(respBody))
	return respData.Text, nil
}
