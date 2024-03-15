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

const MAX_UPLOAD_SIZE = 25 << 20 // 25 mb

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
		file, err := openFileFromFileHeader(fileHeader)
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

		err = client.DoUploadFileToS3(file, s3ObjectKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (client *S3Service) UploadFilesToS3(files []*multipart.FileHeader, s3ObjectPrefix string) error {
	for _, fileHeader := range files {
		file, err := openFileFromFileHeader(fileHeader)
		if err != nil {
			return err
		}
		defer file.Close()

		fileName := fmt.Sprintf("%s%s", strings.Split(fileHeader.Filename, ".")[0], filepath.Ext(fileHeader.Filename))
		s3ObjectKey := fmt.Sprintf("%s/%s", s3ObjectPrefix, fileName)

		err = client.DoUploadFileToS3(file, s3ObjectKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (client *S3Service) DoUploadFileToS3(file io.Reader, objectKey string) error {
	bucketName := os.Getenv("AWS_S3_STORE_BUCKET_NAME")
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
		// log.Println(result.Contents)
		contents = result.Contents
	}
	return contents, err
}

// DownloadFile gets an object from a bucket and stores it in a local file.
func (client *S3Service) DownloadFile(bucketName string, objectKey string, fileName string) error {
	result, err := client.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
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
func openFileFromFileHeader(fileHeader *multipart.FileHeader) (multipart.File, error) {
	if fileHeader.Size > MAX_UPLOAD_SIZE {
		return nil, fmt.Errorf("the uploaded image is too big: %s. Please use an image less than 25MB in size", fileHeader.Filename)
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
