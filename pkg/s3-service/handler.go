package s3Service

import (
	"fmt"
	"log"
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
	log.Println("==userId", userId)

	// payload.path start with project_code without prefix /
	var payload GetPresignedPayload
	err = utils.ReadJSON(w, r, &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "payload", http.StatusBadRequest)
	}

	bucketName := os.Getenv("AWS_S3_STORE_BUCKET_NAME")
	objectKey := fmt.Sprintf("applicant/user_%d/%s", userId, payload.Path)
	presignResult, err := h.presigner.GetObject(bucketName, objectKey, 3600)
	if err != nil {
		utils.ErrorJSON(w, err, "presign")
	}
	log.Println("==presignResult", presignResult)
	utils.WriteJSON(w, http.StatusOK, presignResult)
}
