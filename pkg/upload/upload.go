package upload

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

const MAX_UPLOAD_SIZE = 25 << 20 // 25 mb

type S3Service struct {
	S3Client *s3.Client
}

func NewS3Service(s3Client *s3.Client) *S3Service {
	return &S3Service{
		S3Client: s3Client,
	}
}

func (client *S3Service) UploadToS3(files []*multipart.FileHeader, objectPrefix string) error {
	for _, fileHeader := range files {
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			return fmt.Errorf("the uploaded image is too big: %s. Please use an image less than 32MB in size", fileHeader.Filename)
		}

		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			return err
		}

		filetype := http.DetectContentType(buff)
		headerContentType := fileHeader.Header["Content-Type"][0]

		if !isAllowedContentType(filetype) && !isDocType(filetype, headerContentType) {
			return fmt.Errorf("the provided file format is not allowed. got %s", filetype)
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
		log.Println("==bucketName", bucketName)
		fileName := fmt.Sprintf("%s%s", strings.Split(fileHeader.Filename, ".")[0], filepath.Ext(fileHeader.Filename))
		objectKey := fmt.Sprintf("%s/%s", objectPrefix, fileName)
		log.Println("===fileName", fileName)
		log.Println("===objectKey", objectKey)
		// fmt.Sprintf("%s/%s_%d%s", targetDirPath, fileHeader.Filename, time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))

		_, err = client.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				objectKey, bucketName, objectKey, err)
			return err
		}
	}
	return nil
}

func isDocType(detectedType string, contentType string) bool {
	if detectedType == "application/octet-stream" && contentType == "application/msword" {
		return true
	}
	if detectedType == "application/zip" && contentType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
		return true
	}
	return false
}

func isAllowedContentType(mimetype string) bool {
	if mimetype != "image/jpeg" &&
		mimetype != "image/png" &&
		mimetype != "application/msword" &&
		mimetype != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" &&
		mimetype != "application/pdf" {
		return false
	}
	return true
}
