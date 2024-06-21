package s3Service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

// params userId, projectCode,
func (client *S3Service) UpdateAttachmentZipContent(userId int, projectCode string) error {
	folderPath := filepath.Join("../home/tmp/zip")
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	zipFileTargetPath := filepath.Join(folderPath, fmt.Sprintf("%s.zip", projectCode))
	err = client.downloadZipToLocal(zipFileTargetPath, userId, projectCode)
	if err != nil {
		return err
	}
	return nil
}

func (client *S3Service) downloadZipToLocal(zipFileTargetPath string, userId int, projectCode string) error {
	downloader := manager.NewDownloader(client.S3Client)

	zipFile, err := os.Create(zipFileTargetPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	basePrefix := getBasePrefix(userId, projectCode)
	zipObjectKey := fmt.Sprintf("%s/zip/%s_เอกสารแนบ.zip", basePrefix, projectCode)

	_, err = downloader.Download(context.TODO(), zipFile, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_STORE_BUCKET_NAME")),
		Key:    aws.String(zipObjectKey),
	})
	if err != nil {
		return err
	}
	return nil
}

func getBasePrefix(userId int, projectCode string) string {
	return fmt.Sprintf("applicant/user_%d/%s", userId, projectCode)
}
