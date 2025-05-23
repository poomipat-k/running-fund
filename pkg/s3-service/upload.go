package s3Service

import (
	"archive/zip"
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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const MAX_UPLOAD_SIZE = 50 << 20 // 50 mb

type S3Service struct {
	S3Client *s3.Client
}

func NewS3Service(s3Client *s3.Client) *S3Service {
	return &S3Service{
		S3Client: s3Client,
	}
}

func (client *S3Service) ZipAndUploadFileToS3(files []*multipart.FileHeader, zipWriters []*zip.Writer, zipFilePrefix string, s3ObjectPrefix string) error {
	for _, fileHeader := range files {
		file, err := OpenFileFromFileHeader(fileHeader)
		if err != nil {
			return err
		}
		defer file.Close()

		fileName := fmt.Sprintf("%s%s", strings.Split(fileHeader.Filename, ".")[0], filepath.Ext(fileHeader.Filename))

		zipFilePath := fmt.Sprintf("%s/%s", zipFilePrefix, fileName)
		for _, zipWriter := range zipWriters {
			err = utils.WriteToZip(zipWriter, file, zipFilePath)
			if err != nil {
				return err
			}
		}

		s3ObjectKey := fmt.Sprintf("%s/%s", s3ObjectPrefix, fileName)
		bucketName := os.Getenv("AWS_S3_STORE_BUCKET_NAME")
		err = client.DoUploadFileToS3(file, bucketName, s3ObjectKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (client *S3Service) UploadFilesToS3(files []*multipart.FileHeader, bucketName, s3ObjectPrefix string) error {
	for _, fileHeader := range files {
		file, err := OpenFileFromFileHeader(fileHeader)
		if err != nil {
			return err
		}
		defer file.Close()

		fileName := fmt.Sprintf("%s%s", strings.Split(fileHeader.Filename, ".")[0], filepath.Ext(fileHeader.Filename))
		s3ObjectKey := fmt.Sprintf("%s/%s", s3ObjectPrefix, fileName)
		err = client.DoUploadFileToS3(file, bucketName, s3ObjectKey)
		if err != nil {
			return err
		}
	}
	return nil
}

// User supplied bucketName and objectKey
func (client *S3Service) AdminUploadFilesToS3WithObjectKey(files []*multipart.FileHeader, bucketName, s3ObjectKey string) error {
	for _, fileHeader := range files {
		file, err := OpenFileFromFileHeaderForAdmin(fileHeader)
		if err != nil {
			return err
		}
		defer file.Close()

		contentType := "application/octet-stream"
		headerContentType := fileHeader.Header["Content-Type"][0]
		if headerContentType != "" {
			contentType = headerContentType
		}
		err = client.DoUploadFileToS3withContentType(file, bucketName, s3ObjectKey, contentType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (client *S3Service) DoUploadFileToS3(file io.Reader, bucketName, objectKey string) error {
	_, err := client.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			objectKey, bucketName, objectKey, err)
		return err
	}

	return nil
}

func (client *S3Service) DoUploadFileToS3withContentType(file io.Reader, bucketName, objectKey string, contentType string) error {
	_, err := client.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			objectKey, bucketName, objectKey, err)
		return err
	}

	return nil
}

// ListObjects lists the objects in a bucket.
func (client *S3Service) ListObjects(bucketName string, prefix string) ([]types.Object, error) {
	result, err := client.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix), // aws.String("applicant/user_34/NOV67_1201/"),
	})
	var contents []types.Object
	if err != nil {
		log.Printf("Couldn't list objects in bucket %v. Here's why: %v\n", bucketName, err)
	} else {
		contents = result.Contents
	}
	return contents, err
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

// openFile and validate file type
func OpenFileFromFileHeader(fileHeader *multipart.FileHeader) (multipart.File, error) {
	if fileHeader.Size > MAX_UPLOAD_SIZE {
		return nil, fmt.Errorf("the uploaded file is too big: %s. Please use an file less than 25MB in size", fileHeader.Filename)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return nil, err
	}

	filetype := http.DetectContentType(buff)
	headerContentType := fileHeader.Header["Content-Type"][0]

	if !isAllowedContentType(filetype) && !isDocType(filetype, headerContentType) {
		return nil, fmt.Errorf("the provided file format is not allowed. got %s", filetype)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func OpenFileFromFileHeaderForAdmin(fileHeader *multipart.FileHeader) (multipart.File, error) {
	if fileHeader.Size > MAX_UPLOAD_SIZE {
		return nil, fmt.Errorf("the uploaded file is too big: %s. Please use an file less than 25MB in size", fileHeader.Filename)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	headerContentType := fileHeader.Header["Content-Type"][0]
	if !adminAllowContentType(headerContentType) {
		return nil, fmt.Errorf("the provided file format is not allowed. got %s", headerContentType)
	}
	return file, nil
}

func adminAllowContentType(contentType string) bool {
	// allow contentType image/*
	first := strings.Split(contentType, "/")[0]
	return first == "image"
}
