package cms

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	s3Service "github.com/poomipat-k/running-fund/pkg/s3-service"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const minDashboardYear = 2023

type CmsHandler struct {
	awsS3Service s3Service.S3Service
	store        cmdStore
}

type cmdStore interface {
	GetReviewPeriod() (ReviewPeriod, error)
	GetAdminWebsiteDashboardDateConfigPreview(fromDate, toDate time.Time, limit, offset int) ([]AdminDateConfigPreviewRow, error)
	AdminUpdateWebsiteConfig(payload AdminUpdateWebsiteConfigRequest) error
	GetLandingPageContent() (LandingConfig, error)
	GetWebsiteConfigData() (AdminUpdateWebsiteConfigRequest, string, error)
	GetFAQ() ([]FAQ, error)
}

func NewCmsHandler(awsS3Service s3Service.S3Service, store cmdStore) *CmsHandler {
	return &CmsHandler{
		awsS3Service: awsS3Service,
		store:        store,
	}
}

func (h *CmsHandler) GetReviewPeriod(w http.ResponseWriter, r *http.Request) {
	period, err := h.store.GetReviewPeriod()
	if err != nil {
		slog.Error(err.Error())
		utils.ErrorJSON(w, err, "")
		return
	}

	utils.WriteJSON(w, http.StatusOK, period)
}

func (h *CmsHandler) AdminUploadContentFiles(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(25 << 20); err != nil {
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}
	formJsonString := r.FormValue("form")
	payload := UploadFileRequest{}
	err := json.Unmarshal([]byte(formJsonString), &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "")
		return
	}

	if payload.Name == "" {
		utils.ErrorJSON(w, errors.New("empty upload formData name"), "name")
		return
	}
	if payload.PathPrefix == "" {
		utils.ErrorJSON(w, errors.New("empty upload formData pathPrefix"), "pathPrefix")
		return
	}

	fileHeaders := r.MultipartForm.File[payload.Name]
	if len(fileHeaders) == 0 {
		utils.ErrorJSON(w, fmt.Errorf("%s is empty", payload.Name), payload.Name, http.StatusBadRequest)
		return
	}
	bucketName := os.Getenv("AWS_S3_STATIC_BUCKET_NAME")
	if bucketName == "" {
		utils.ErrorJSON(w, errors.New("AWS_S3_STATIC_BUCKET_NAME is empty"), "AWS_S3_STATIC_BUCKET_NAME", http.StatusInternalServerError)
		return
	}
	fileHeader := fileHeaders[0]
	objectKey := fmt.Sprintf("%s/%s_%d%s", payload.PathPrefix, strings.Split(fileHeader.Filename, ".")[0], time.Now().Unix(), filepath.Ext(fileHeader.Filename))
	err = h.awsS3Service.AdminUploadFilesToS3WithObjectKey(fileHeaders, bucketName, objectKey)
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

// Admin date setting preview dashboard
func (h *CmsHandler) GetAdminWebsiteDashboardDateConfigPreview(w http.ResponseWriter, r *http.Request) {
	var payload GetAdminDashboardDateConfigPreviewRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "payload", http.StatusBadRequest)
		return
	}

	errField, err := validateAdminWebsiteDashboardDateConfigPreviewRequest(payload)
	if err != nil {
		utils.ErrorJSON(w, err, errField, http.StatusBadRequest)
		return
	}
	loc, err := utils.GetTimeLocation()
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}
	offset := (payload.PageNo - 1) * payload.PageSize
	fromDate := time.Date(payload.FromYear, time.Month(payload.FromMonth), payload.FromDay, 0, 0, 0, 0, loc)
	toDate := time.Date(payload.ToYear, time.Month(payload.ToMonth), payload.ToDay+1, 0, 0, 0, 0, loc)
	records, err := h.store.GetAdminWebsiteDashboardDateConfigPreview(fromDate, toDate, payload.PageSize, offset)
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, records)
}

// Get landing page content
func (h *CmsHandler) GetLandingPageContent(w http.ResponseWriter, r *http.Request) {
	data, err := h.store.GetLandingPageContent()
	if err != nil {
		utils.ErrorJSON(w, err, "store", http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}

// Get FAQ content
func (h *CmsHandler) GetFAQ(w http.ResponseWriter, r *http.Request) {
	faqList, err := h.store.GetFAQ()
	if err != nil {
		utils.ErrorJSON(w, err, "store", http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, http.StatusOK, faqList)
}

func (h *CmsHandler) GetWebsiteConfigData(w http.ResponseWriter, r *http.Request) {
	data, errName, err := h.store.GetWebsiteConfigData()
	if err != nil {
		utils.ErrorJSON(w, err, errName, http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, http.StatusOK, data)
}

func (h *CmsHandler) AdminUpdateWebsiteConfig(w http.ResponseWriter, r *http.Request) {
	var payload AdminUpdateWebsiteConfigRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "payload", http.StatusBadRequest)
		return
	}
	fn, err := validateAdminUpdateWebsiteConfigRequest(payload)
	if err != nil {
		utils.ErrorJSON(w, err, fn, http.StatusBadRequest)
		return
	}
	err = h.store.AdminUpdateWebsiteConfig(payload)
	if err != nil {
		utils.ErrorJSON(w, err, "update website config", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, CommonSuccessResponse{Success: true, Message: "Successfully updated website config"})
}

func validateAdminWebsiteDashboardDateConfigPreviewRequest(payload GetAdminDashboardDateConfigPreviewRequest) (string, error) {
	fn, err := validateFormDateToDate(payload.FromYear, payload.FromMonth, payload.FromDay, payload.ToYear, payload.ToMonth, payload.ToDay)
	if err != nil {
		return fn, err
	}

	if payload.PageNo <= 0 {
		return "pageNo", &PageNoInvalidError{}
	}
	if payload.PageSize < 1 {
		return "pageSize", &PageSizeInvalidError{}
	}

	return "", nil
}

func validateAdminUpdateWebsiteConfigRequest(payload AdminUpdateWebsiteConfigRequest) (string, error) {
	fn, err := validateFormDateToDate(
		payload.Dashboard.FromYear,
		payload.Dashboard.FromMonth,
		payload.Dashboard.FromDay,
		payload.Dashboard.ToYear,
		payload.Dashboard.ToMonth,
		payload.Dashboard.ToDay,
	)
	if err != nil {
		return fn, err
	}

	fn, err = validateFaq(payload.Faq)
	if err != nil {
		return fn, err
	}

	fn, err = validateFooter(payload.Footer)
	if err != nil {
		return fn, err
	}

	return "", nil
}

func validateFormDateToDate(fromYear, fromMonth, fromDay, toYear, toMonth, toDay int) (string, error) {
	if fromYear < minDashboardYear {
		return "fromYear", &FromYearRequiredError{}
	}
	if fromMonth == 0 {
		return "fromMonth", &MonthRequiredError{}
	}
	if fromMonth < 1 || fromMonth > 12 {
		return "fromMonth", &MonthOutOfBoundError{}
	}

	if toYear < minDashboardYear {
		return "toYear", &ToYearRequiredError{}
	}
	if toMonth == 0 {
		return "toMonth", &MonthRequiredError{}
	}
	if toMonth < 1 || toMonth > 12 {
		return "toMonth", &MonthOutOfBoundError{}
	}
	loc, err := utils.GetTimeLocation()
	if err != nil {
		return "timeLocation", nil
	}
	fromDate := time.Date(fromYear, time.Month(fromMonth), fromDay, 0, 0, 0, 0, loc)
	toDate := time.Date(toYear, time.Month(toMonth), toDay+1, 0, 0, 0, 0, loc)
	if fromDate.After(toDate) {
		return "fromDate", &FromDateExceedToDateError{}
	}
	return "", nil
}

func validateFaq(faqList []FAQ) (string, error) {
	for i, faq := range faqList {
		if faq.Question == "" {
			fn := fmt.Sprintf("question[%d]", i)
			return fn, fmt.Errorf("%s is empty", fn)
		}
		if faq.Answer == "" {
			fn := fmt.Sprintf("answer[%d]", i)
			return fn, fmt.Errorf("%s is empty", fn)
		}
	}
	return "", nil
}

func validateFooter(footer FooterConfig) (string, error) {
	if footer.Contact.Email == "" {
		return "email", fmt.Errorf("email is empty")
	}
	if footer.Contact.PhoneNumber == "" {
		return "email", fmt.Errorf("PhoneNumber is empty")
	}
	if footer.Contact.FromHour == "" {
		return "email", fmt.Errorf("FromHour is empty")
	}
	if footer.Contact.FromMinute == "" {
		return "email", fmt.Errorf("FromMinute is empty")
	}
	if footer.Contact.ToHour == "" {
		return "email", fmt.Errorf("ToHour is empty")
	}
	if footer.Contact.ToHour == "" {
		return "email", fmt.Errorf("ToHour is empty")
	}
	return "", nil
}
