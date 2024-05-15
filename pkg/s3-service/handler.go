package s3Service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const PUT_PRESIGNED_DURATION_SECOND = 60

type S3Handler struct {
	// awsS3Service S3Service
	presigner Presigner
}

func NewS3Handler(presigner Presigner) *S3Handler {
	return &S3Handler{
		presigner: presigner,
	}
}

func (h *S3Handler) GeneratePresignedUrl(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdFromRequestHeader(r)
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userId", http.StatusForbidden)
		return
	}
	userRole := utils.GetUserRoleFromRequestHeader(r)
	if userRole == "" {
		msg := "userRole is empty"
		err = errors.New(msg)
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userRole", http.StatusForbidden)
		return
	}

	// payload.path start with project_code without prefix /
	var payload GetPresignedPayload
	err = utils.ReadJSON(w, r, &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "payload", http.StatusBadRequest)
	}

	var objectKey string
	bucketName := os.Getenv("AWS_S3_STORE_BUCKET_NAME")
	if userRole == "applicant" {
		objectKey = fmt.Sprintf("applicant/user_%d/%s", userId, payload.Path)
	} else {
		objectKey = fmt.Sprintf("applicant/user_%d/%s", payload.ProjectCreatedByUserId, payload.Path)
	}
	presignedResult, err := h.presigner.GetObject(bucketName, objectKey, 3600)
	if err != nil {
		utils.ErrorJSON(w, err, "presign")
	}
	utils.WriteJSON(w, http.StatusOK, presignedResult)
}

func (h *S3Handler) GetPresignedPutObject(w http.ResponseWriter, r *http.Request) {
	bucketName := os.Getenv("AWS_S3_STORE_BUCKET_NAME")
	objectKey := "presigned/test-presigned.png"
	lifetimeSecs := 300
	request, err := h.generatePresignedPutObject(bucketName, objectKey, int64(lifetimeSecs))
	if err != nil {
		utils.ErrorJSON(w, err, "presigned", http.StatusBadRequest)
		return
	}
	log.Println("===request:", request)
	utils.WriteJSON(w, http.StatusOK, request)
}

func (h *S3Handler) generatePresignedPutObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := h.presigner.PresignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}
