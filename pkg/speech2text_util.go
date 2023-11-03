package pkg

import (
	"context"
	"fmt"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
)

func CallSpeech2Text(data []byte, languageCode string) (string, error) {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create speech client: %v", err)
	}
	defer client.Close()
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    languageCode, // "en-US", "zh-CN"
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	}

	resp, err := client.Recognize(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to recognize speech: %v", err)
	}

	if len(resp.Results) == 0 {
		return "", nil // No results.
	}

	return resp.Results[0].Alternatives[0].Transcript, nil
}
