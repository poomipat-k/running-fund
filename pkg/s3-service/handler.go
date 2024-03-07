package s3Service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

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

	// payload.path start with project_code without prefix /
	var payload GetPresignedPayload
	err = utils.ReadJSON(w, r, &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "payload", http.StatusBadRequest)
	}

	userRole := utils.GetUserRoleFromRequestHeader(r)
	if userRole == "" {
		msg := "userRole is empty"
		err = errors.New(msg)
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "userRole", http.StatusForbidden)
		return
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
