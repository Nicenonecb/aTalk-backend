package pkg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func CallText2Speech(text string) ([]byte, error) {
	var url = "https://aigptx.top/v1/audio/speech"
	data := map[string]interface{}{
		"model":           "tts-1",
		"voice":           "alloy",
		"input":           text,
		"response_format": "mp3",
		"speed":           "1.0",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
