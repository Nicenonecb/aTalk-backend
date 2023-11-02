package handler

import (
	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DialogueHandler(c *gin.Context) {
	ctx := context.Background()

	// Creates a client.
	client, err := speech.NewClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}
	defer client.Close()

	// The path to the remote audio file to transcribe.
	fileURI := "gs://client_post/desired-object-name"

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: fileURI},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recognize speech"})
		return
	}

	var transcripts []string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			transcripts = append(transcripts, alt.Transcript)
		}
	}
	c.JSON(http.StatusOK, gin.H{"transcripts": transcripts})
}
