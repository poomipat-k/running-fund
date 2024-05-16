package cms

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	s3Service "github.com/poomipat-k/running-fund/pkg/s3-service"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

type CmsHandler struct {
	awsS3Service s3Service.S3Service
}

func NewCmsHandler(awsS3Service s3Service.S3Service) *CmsHandler {
	return &CmsHandler{
		awsS3Service: awsS3Service,
	}
}

func (h *CmsHandler) AdminUploadContentFiles(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(25 << 20); err != nil {
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}
	banner := r.MultipartForm.File["banner"]
	if len(banner) == 0 {
		utils.ErrorJSON(w, errors.New("banner is empty"), "banner", http.StatusBadRequest)
		return
	}
	bucketName := os.Getenv("AWS_S3_STATIC_BUCKET_NAME")
	if bucketName == "" {
		utils.ErrorJSON(w, errors.New("AWS_S3_STATIC_BUCKET_NAME is empty"), "AWS_S3_STATIC_BUCKET_NAME", http.StatusInternalServerError)
		return
	}
	fileHeader := banner[0]
	objectKey := fmt.Sprintf("%s/%s_%d%s", "banner", strings.Split(fileHeader.Filename, ".")[0], time.Now().Unix(), filepath.Ext(fileHeader.Filename))
	err := h.awsS3Service.UploadFilesToS3WithObjectKey(banner, bucketName, objectKey)
	if err != nil {
		utils.ErrorJSON(w, err, "s3Upload", http.StatusInternalServerError)
		return
	}
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		utils.ErrorJSON(w, errors.New("AWS_REGION is missing"), "AWS_REGION", http.StatusInternalServerError)
		return
	}
	fullPath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, awsRegion, objectKey)
	// Todo: upload another small preview image

	utils.WriteJSON(w, http.StatusOK, S3UploadResponse{
		ObjectKey: objectKey,
		FullPath:  fullPath,
	})
}
