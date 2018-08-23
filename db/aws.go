package db

import (
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func s3ImageUpload(img ReqImage) (interface{}, error) {
	if img.Error != nil {
		return nil, img.Error
	}
	newFile, err := os.Create(img.Header.Filename)
	if err != nil {
		return nil, err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, img.File)
	if err != nil {
		return nil, err
	}

	err = newFile.Sync()
	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("ap-southeast-2")}))
	url, err := uploader.Upload(&s3manager.UploadInput{
		Body:        newFile,
		Bucket:      aws.String("brew-site"),
		Key:         aws.String(newFile.Name()),
		ContentType: aws.String(mime.TypeByExtension(filepath.Ext(newFile.Name()))),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		return nil, err
	}

	return url.Location, nil
}
