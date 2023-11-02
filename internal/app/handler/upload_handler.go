package handler

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func UploadToGCSHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	bucketName := "client_post"         // 替换为你的bucket名字
	objectName := "desired-object-name" // 你希望在GCS中存储的对象名称

	// 使用上述的streamFileUpload函数
	if err := streamFileUpload(file, bucketName, objectName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully."})
}

// streamFileUpload uploads an object via a stream.
func streamFileUpload(w io.Reader, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	b := []byte("Hello world.")
	buf := bytes.NewBuffer(b)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

	if _, err = io.Copy(wc, buf); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	//fmt.Fprintf(w, "%v uploaded to %v.\n", object, bucket)

	return nil
}
