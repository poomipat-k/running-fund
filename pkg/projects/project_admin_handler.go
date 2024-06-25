package projects

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const minDashboardYear = 2023

func (h *ProjectHandler) AdminUpdateProject(w http.ResponseWriter, r *http.Request) {
	projectCode := chi.URLParam(r, "projectCode")
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formJsonString := r.FormValue("form")
	payload := AdminUpdateProjectRequest{}
	err := json.Unmarshal([]byte(formJsonString), &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "")
		return
	}
	additionFiles := r.MultipartForm.File["additionFiles"]
	etcFiles := r.MultipartForm.File["etcFiles"]
	field, err := validateAdminUpdateProjectPayload(payload)
	if err != nil {
		utils.ErrorJSON(w, err, field)
		return
	}

	currentProject, err := h.store.GetProjectStatusByProjectCode(projectCode)
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusNotFound)
		return
	}
	err = h.doUpdateProject(currentProject, payload, projectCode, additionFiles, etcFiles)
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, currentProject.ProjectHistoryId)
}

func (h *ProjectHandler) GetAdminRequestDashboard(w http.ResponseWriter, r *http.Request) {
	var payload GetAdminDashboardRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "payload", http.StatusBadRequest)
		return
	}
	errField, err := validateGetAdminDashboardPayload(payload)
	if err != nil {
		utils.ErrorJSON(w, err, errField, http.StatusBadRequest)
		return
	}
	loc, err := getTimeLocation()
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}
	offset := (payload.PageNo - 1) * payload.PageSize
	orderByStmt := strings.Join(payload.SortBy, ", ")
	if payload.IsAsc {
		orderByStmt += " ASC"
	} else {
		orderByStmt += " DESC"
	}
	fromDate := time.Date(payload.FromYear, time.Month(payload.FromMonth), payload.FromDay, 0, 0, 0, 0, loc)
	toDate := time.Date(payload.ToYear, time.Month(payload.ToMonth), payload.ToDay+1, 0, 0, 0, 0, loc)
	records, err := h.store.GetAdminRequestDashboard(fromDate, toDate, orderByStmt, payload.PageSize, offset, payload.ProjectCode, payload.ProjectName, payload.ProjectStatus)
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, records)
}

func (h *ProjectHandler) GetAdminStartedDashboard(w http.ResponseWriter, r *http.Request) {
	var payload GetAdminDashboardRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "payload", http.StatusBadRequest)
		return
	}
	errField, err := validateGetAdminDashboardPayload(payload)
	if err != nil {
		utils.ErrorJSON(w, err, errField, http.StatusBadRequest)
		return
	}
	loc, err := getTimeLocation()
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}
	offset := (payload.PageNo - 1) * payload.PageSize
	orderByStmt := strings.Join(payload.SortBy, ", ")
	if payload.IsAsc {
		orderByStmt += " ASC"
	} else {
		orderByStmt += " DESC"
	}
	fromDate := time.Date(payload.FromYear, time.Month(payload.FromMonth), payload.FromDay, 0, 0, 0, 0, loc)
	toDate := time.Date(payload.ToYear, time.Month(payload.ToMonth), payload.ToDay+1, 0, 0, 0, 0, loc)
	records, err := h.store.GetAdminStartedDashboard(fromDate, toDate, orderByStmt, payload.PageSize, offset, payload.ProjectCode, payload.ProjectName, payload.ProjectStatus)
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, records)
}

func (h *ProjectHandler) GetAdminSummary(w http.ResponseWriter, r *http.Request) {
	var payload GetAdminSummaryRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "payload", http.StatusBadRequest)
		return
	}
	errField, err := validateGetAdminSummaryRequestPayload(payload)
	if err != nil {
		utils.ErrorJSON(w, err, errField, http.StatusBadRequest)
		return
	}
	loc, err := getTimeLocation()
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}
	fromDate := time.Date(payload.FromYear, time.Month(payload.FromMonth), payload.FromDay, 0, 0, 0, 0, loc)
	toDate := time.Date(payload.ToYear, time.Month(payload.ToMonth), payload.ToDay+1, 0, 0, 0, 0, loc)
	records, err := h.store.GetAdminSummary(fromDate, toDate)
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, records)
}

func (h *ProjectHandler) GenerateAdminReport(w http.ResponseWriter, r *http.Request) {
	var payload GenerateAdminReportRequest
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		utils.ErrorJSON(w, err, "payload", http.StatusBadRequest)
		return
	}
	errField, err := validateGenerateAdminReportRequest(payload)
	if err != nil {
		utils.ErrorJSON(w, err, errField, http.StatusBadRequest)
		return
	}
	loc, err := getTimeLocation()
	if err != nil {
		utils.ErrorJSON(w, err, "", http.StatusInternalServerError)
		return
	}
	// fromDate <= project.created_at < toDate
	fromDate := time.Date(payload.FromYear, time.Month(payload.FromMonth), payload.FromDay, 0, 0, 0, 0, loc)
	toDate := time.Date(payload.ToYear, time.Month(payload.ToMonth), payload.ToDay+1, 0, 0, 0, 0, loc)

	buffer, err := h.store.GenerateAdminReport(fromDate, toDate)
	if err != nil {
		utils.ErrorJSON(w, err, "report", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, buffer.String())
}

func (h *ProjectHandler) doUpdateProject(
	currentProject AdminUpdateParam,
	payload AdminUpdateProjectRequest,
	projectCode string,
	additionFiles []*multipart.FileHeader,
	etcFiles []*multipart.FileHeader,
) error {
	currentStatus := currentProject.ProjectStatus
	primaryStatusChanged := hasPrimaryStatusChanged(currentStatus, payload.ProjectStatusPrimary)
	secondaryStatusChanged := payload.ProjectStatusSecondary != currentStatus
	now := time.Now()
	if !primaryStatusChanged && !secondaryStatusChanged {
		// change all other attributes
		err := h.store.UpdateProjectByAdmin(
			AdminUpdateParam{
				ProjectHistoryId:   currentProject.ProjectHistoryId,
				ProjectStatus:      currentStatus,
				AdminScore:         payload.AdminScore,
				FundApprovedAmount: payload.FundApprovedAmount,
				AdminComment:       payload.AdminComment,
				AdminApprovedAt:    currentProject.AdminApprovedAt,
				UpdatedAt:          now,
			},
			currentProject.CreatedBy,
			projectCode,
			additionFiles,
			etcFiles,
		)
		if err != nil {
			return err
		}
		return nil
	}

	var newStatus string
	newStatus = currentStatus
	if primaryStatusChanged {
		newStatus = payload.ProjectStatusPrimary
	} else if secondaryStatusChanged {
		newStatus = payload.ProjectStatusSecondary
	}

	if newStatus == "Approved" {
		// update admin_approved_at to now
		var approvedAt *time.Time
		if payload.AdminApprovedAt != nil {
			approvedAt = payload.AdminApprovedAt
		} else {
			approvedAt = &now
		}
		err := h.store.UpdateProjectByAdmin(
			AdminUpdateParam{
				ProjectHistoryId:   currentProject.ProjectHistoryId,
				ProjectStatus:      newStatus,
				AdminScore:         payload.AdminScore,
				FundApprovedAmount: payload.FundApprovedAmount,
				AdminComment:       payload.AdminComment,
				AdminApprovedAt:    approvedAt,
				UpdatedAt:          now,
			},
			currentProject.CreatedBy,
			projectCode,
			additionFiles,
			etcFiles,
		)
		if err != nil {
			return err
		}
		return nil
	}
	if newStatus == "NotApproved" {
		// update admin_approved_at to nils
		err := h.store.UpdateProjectByAdmin(
			AdminUpdateParam{
				ProjectHistoryId:   currentProject.ProjectHistoryId,
				ProjectStatus:      newStatus,
				AdminScore:         payload.AdminScore,
				FundApprovedAmount: payload.FundApprovedAmount,
				AdminComment:       payload.AdminComment,
				AdminApprovedAt:    nil,
				UpdatedAt:          now,
			},
			currentProject.CreatedBy,
			projectCode,
			additionFiles,
			etcFiles,
		)
		if err != nil {
			return err
		}
		return nil
	}
	err := h.store.UpdateProjectByAdmin(
		AdminUpdateParam{
			ProjectHistoryId:   currentProject.ProjectHistoryId,
			ProjectStatus:      newStatus,
			AdminScore:         payload.AdminScore,
			FundApprovedAmount: payload.FundApprovedAmount,
			AdminComment:       payload.AdminComment,
			AdminApprovedAt:    currentProject.AdminApprovedAt,
			UpdatedAt:          now,
		},
		currentProject.CreatedBy,
		projectCode,
		additionFiles,
		etcFiles,
	)
	if err != nil {
		return err
	}
	return nil
}
