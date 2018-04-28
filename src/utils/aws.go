package utils

import (
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type chanResult struct {
	Data interface{} `json:"data"`
	OK   bool        `json:"ok"`
}

// UploadToS3 func
func UploadToS3(done chan interface{}, mf multipart.File, mh *multipart.FileHeader) chan Result {
	result := make(chan Result)

	go func() {
		defer fmt.Println("S3 Result channel closed.")
		defer close(result)

		select {
		case <-done:
			return
		case result <- handleS3Upload(mf, mh):
			return
		}
	}()
	return result
}

func handleS3Upload(f multipart.File, h *multipart.FileHeader) Result {
	file, err := os.Create(h.Filename)
	if err != nil {
		result.Error = &Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + ": Error creating file -> " + err.Error(),
		}
		return result
	}
	defer file.Close()

	if _, err = io.Copy(file, f); err != nil {
		result.Error = &Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + ": Error copying file contents -> " + err.Error(),
		}
		return result
	}

	err = file.Sync()
	if err != nil {
		log.Println("Error flushing file to memory -> ", err.Error())
	}

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("ap-southeast-2")}))
	url, err := uploader.Upload(&s3manager.UploadInput{
		Body:        file,
		Bucket:      aws.String("brew-site"),
		Key:         aws.String(file.Name()),
		ContentType: aws.String(mime.TypeByExtension(filepath.Ext(file.Name()))),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		result.Error = &Error{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError) + ": Error uploading file to aws -> " + err.Error(),
		}
		return result
	}
	result.Success = &Success{
		Status: http.StatusOK,
		Data:   url.Location,
	}
	return result
}
