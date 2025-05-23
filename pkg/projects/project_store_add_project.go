package projects

import (
	"archive/zip"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/poomipat-k/running-fund/pkg/utils"
)

const (
	collaborationStr = "หนังสือนำส่ง"
	attachmentsStr   = "เอกสารแนบ"
)

func (s *store) AddProject(
	payload AddProjectRequest,
	userId int,
	criteria []ApplicantSelfScoreCriteria,
	attachments []Attachments,
) (int, error) {
	projectCode, err := s.generateProjectCode()
	if err != nil {
		return 0, err
	}
	// start transaction
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return failAdd("tx", err)
	}

	defer tx.Rollback()

	now := time.Now()
	// Add address rows
	addressId, err := addGeneralAddress(ctx, tx, payload)
	if err != nil {
		return failAdd("addressId", err)
	}

	projectHeadAddressId, err := addProjectContactPersonAddress(ctx, tx, payload, "projectHead")
	if err != nil {
		return failAdd("projectHeadAddressId", err)
	}

	// Add contact rows
	projectHeadContactId, err := addProjectHeadContactId(ctx, tx, payload, projectHeadAddressId)
	if err != nil {
		return failAdd("projectHeadContactId", err)
	}

	// check projectManager is the same person as projectHead otherwise create a new contact
	var projectManagerContactId int
	if payload.Contact.ProjectManager == payload.Contact.ProjectHead {
		projectManagerContactId = projectHeadContactId
	} else {
		projectManagerAddressId, err := addProjectContactPersonAddress(ctx, tx, payload, "projectManager")
		if err != nil {
			return failAdd("projectManagerAddressId", err)
		}
		projectManagerContactId, err = addProjectManagerContactId(ctx, tx, payload, projectManagerAddressId)
		if err != nil {
			return failAdd("projectManagerContactId", err)
		}
	}

	// check projectManager is the same person as projectHead or projectManager otherwise create a new contact
	var projectCoordinatorContactId int
	if payload.Contact.ProjectCoordinator == payload.Contact.ProjectHead {
		projectCoordinatorContactId = projectHeadContactId
	} else if payload.Contact.ProjectCoordinator == payload.Contact.ProjectManager {
		projectCoordinatorContactId = projectManagerContactId
	} else {
		projectCoordinatorAddressId, err := addProjectContactPersonAddress(ctx, tx, payload, "projectCoordinator")
		if err != nil {
			return failAdd("projectCoordinatorAddressId", err)
		}

		projectCoordinatorContactId, err = addProjectCoordinatorContactId(ctx, tx, payload, projectCoordinatorAddressId)
		if err != nil {
			return failAdd("projectCoordinatorContactId", err)
		}
	}

	var projectRaceDirectorContactId int
	if payload.Contact.RaceDirector.Who == "other" {
		projectRaceDirectorContactId, err = addProjectRaceDirectorContactId(ctx, tx, payload)
		if err != nil {
			return failAdd("projectRaceDirectorContactId", err)
		}
	}
	baseFilePrefix := getBasePrefix(userId, projectCode)
	// Add project_history
	projectHistoryId, err := addProjectHistory(
		ctx,
		tx,
		payload,
		projectCode,
		now,
		addressId,
		projectHeadContactId,
		projectManagerContactId,
		projectCoordinatorContactId,
		projectRaceDirectorContactId,
		baseFilePrefix,
	)
	if err != nil {
		return failAdd("projectHistoryId", err)
	}

	// Add project
	projectId, err := addProjectRow(ctx, tx, projectCode, now, projectHistoryId, userId)
	if err != nil {
		return failAdd("projectId", err)
	}

	// Add distance
	_, err = addDistances(ctx, tx, payload, projectHistoryId)
	if err != nil {
		return failAdd("distanceRowsAffected", err)
	}
	// Add applicant scores
	_, err = addApplicantScores(ctx, tx, payload, projectHistoryId, criteria)
	if err != nil {
		return failAdd("applicantScoreRowsAffected", err)
	}

	// Write zip, upload files and zips
	err = s.handleCreateProjectFiles(baseFilePrefix, userId, projectCode, payload, attachments)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return failAdd("tx.Commit()", err)
	}
	// commit
	slog.Info("success adding a new project", "projectCode", projectCode, "userId", userId)
	return projectId, nil
}

func addProjectHistory(
	ctx context.Context,
	tx *sql.Tx,
	payload AddProjectRequest,
	projectCode string,
	now time.Time,
	addressId int,
	projectHeadContactId int,
	projectManagerContactId int,
	projectCoordinatorContactId int,
	rawProjectRaceDirectorContactId int,
	baseFilePrefix string,
) (int, error) {
	fromDate, toDate, thisSeriesLatestCompletedDate, err := buildTimeFromPayload(payload)
	if err != nil {
		slog.Error("addProjectHistory err get fromDate and toDate", "error", err)
		return 0, err
	}
	projectRaceDirectorContactId := rawProjectRaceDirectorContactId
	if rawProjectRaceDirectorContactId == 0 {
		if payload.Contact.RaceDirector.Who == "projectHead" {
			projectRaceDirectorContactId = projectHeadContactId
		} else if payload.Contact.RaceDirector.Who == "projectManager" {
			projectRaceDirectorContactId = projectManagerContactId
		} else if payload.Contact.RaceDirector.Who == "projectCoordinator" {
			projectRaceDirectorContactId = projectCoordinatorContactId
		}
	}

	var nilAbleThisSeriesLatestCompletedDate *time.Time
	if thisSeriesLatestCompletedDate.IsZero() {
		nilAbleThisSeriesLatestCompletedDate = nil
	} else {
		nilAbleThisSeriesLatestCompletedDate = &thisSeriesLatestCompletedDate
	}

	var id int
	err = tx.QueryRowContext(
		ctx,
		addProjectHistorySQL,
		projectCode,
		1,
		now,         // created_at
		now,         // updated_at
		"Reviewing", // valid status: Reviewing, Reviewed, Revise, NotApproved, Approved, Start, Completed
		payload.Collaborated,
		payload.General.ProjectName,
		fromDate,
		toDate,
		addressId,
		payload.General.StartPoint,
		payload.General.FinishPoint,
		payload.General.EventDetails.Category.Available.RoadRace,
		payload.General.EventDetails.Category.Available.TrailRunning,
		payload.General.EventDetails.Category.Available.Other,
		payload.General.EventDetails.Category.OtherType,
		payload.General.EventDetails.VIP,
		payload.General.EventDetails.VIPFee,
		payload.General.ExpectedParticipants,
		payload.General.HasOrganizer,
		payload.General.OrganizerName,
		projectHeadContactId,
		projectManagerContactId,
		projectCoordinatorContactId,
		projectRaceDirectorContactId,
		payload.Contact.Organization.Type,
		payload.Contact.Organization.Name,
		payload.Details.Background,
		payload.Details.Objective,
		payload.Details.Marketing.Online.Available.Facebook,
		payload.Details.Marketing.Online.HowTo.Facebook,
		payload.Details.Marketing.Online.Available.Website,
		payload.Details.Marketing.Online.HowTo.Website,
		payload.Details.Marketing.Online.Available.OnlinePage,
		payload.Details.Marketing.Online.HowTo.OnlinePage,
		payload.Details.Marketing.Online.Available.Other,
		payload.Details.Marketing.Online.HowTo.Other,
		payload.Details.Marketing.Offline.Available.PR,
		payload.Details.Marketing.Offline.Available.LocalOfficial,
		payload.Details.Marketing.Offline.Available.Booth,
		payload.Details.Marketing.Offline.Available.Billboard,
		payload.Details.Marketing.Offline.Available.TV,
		payload.Details.Marketing.Offline.Available.Other,
		payload.Details.Marketing.Offline.Addition,
		payload.Details.Safety.Ready.RunnerInformation,
		payload.Details.Safety.Ready.HealthDecider,
		payload.Details.Safety.Ready.Ambulance,
		payload.Details.Safety.Ready.FirstAid,
		payload.Details.Safety.Ready.AED,
		payload.Details.Safety.AEDCount,
		payload.Details.Safety.Ready.VolunteerDoctor,
		payload.Details.Safety.Ready.Insurance,
		payload.Details.Safety.Ready.Other,
		payload.Details.Safety.Addition,
		payload.Details.Route.Measurement.AthleticsAssociation,
		payload.Details.Route.Measurement.CalibratedBicycle,
		payload.Details.Route.Measurement.SelfMeasurement,
		payload.Details.Route.Tool,
		payload.Details.Route.TrafficManagement.AskPermission,
		payload.Details.Route.TrafficManagement.HasSupporter,
		payload.Details.Route.TrafficManagement.RoadClosure,
		payload.Details.Route.TrafficManagement.Signs,
		payload.Details.Route.TrafficManagement.Lighting,
		payload.Details.Judge.Type,
		payload.Details.Judge.OtherType,
		payload.Details.Support.Organization.ProvincialAdministration,
		payload.Details.Support.Organization.Safety,
		payload.Details.Support.Organization.Health,
		payload.Details.Support.Organization.Volunteer,
		payload.Details.Support.Organization.Community,
		payload.Details.Support.Organization.Other,
		payload.Details.Support.Addition,
		payload.Details.Feedback,
		*payload.Experience.ThisSeries.FirstTime,
		payload.Experience.ThisSeries.History.OrdinalNumber,
		nilAbleThisSeriesLatestCompletedDate,
		payload.Experience.ThisSeries.History.Completed1.Year,
		payload.Experience.ThisSeries.History.Completed1.Participant,
		payload.Experience.ThisSeries.History.Completed2.Year,
		payload.Experience.ThisSeries.History.Completed2.Participant,
		payload.Experience.ThisSeries.History.Completed3.Year,
		payload.Experience.ThisSeries.History.Completed3.Participant,
		payload.Experience.OtherSeries.DoneBefore,
		payload.Experience.OtherSeries.History.Completed1.Year,
		payload.Experience.OtherSeries.History.Completed1.Name,
		payload.Experience.OtherSeries.History.Completed1.Participant,
		payload.Experience.OtherSeries.History.Completed2.Year,
		payload.Experience.OtherSeries.History.Completed2.Name,
		payload.Experience.OtherSeries.History.Completed2.Participant,
		payload.Experience.OtherSeries.History.Completed3.Year,
		payload.Experience.OtherSeries.History.Completed3.Name,
		payload.Experience.OtherSeries.History.Completed3.Participant,
		payload.Fund.Budget.Total,
		payload.Fund.Budget.SupportOrganization,
		payload.Fund.Request.Type.Fund,
		payload.Fund.Request.Details.FundAmount,
		payload.Fund.Request.Type.BIB,
		payload.Fund.Request.Details.BibAmount,
		payload.Fund.Request.Type.Pr,
		payload.Fund.Request.Type.Seminar,
		payload.Fund.Request.Details.Seminar,
		payload.Fund.Request.Type.Other,
		payload.Fund.Request.Details.Other,
		payload.Fund.Budget.NoAlcoholSponsor,
		baseFilePrefix,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getTimeLocation() (*time.Location, error) {
	loc, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		return &time.Location{}, err
	}
	return loc, nil
}

func addDistances(ctx context.Context, tx *sql.Tx, payload AddProjectRequest, projectHistoryId int) (int64, error) {
	dis := payload.General.EventDetails.DistanceAndFee
	checkedDistances := []DistanceAndFee{}
	for i := 0; i < len(dis); i++ {
		if dis[i].Checked {
			checkedDistances = append(checkedDistances, dis[i])
		}
	}

	valuesStrPlaceholder := []string{}
	values := []any{}

	for i := 0; i < len(checkedDistances); i++ {
		valuesStrPlaceholder = append(valuesStrPlaceholder, fmt.Sprintf("($%d, $%d, $%d, $%d)", 4*i+1, 4*i+2, 4*i+3, 4*i+4))
		values = append(values, checkedDistances[i].Type, checkedDistances[i].Fee, checkedDistances[i].Dynamic, projectHistoryId)
	}
	customSQL := addManyDistanceSQL + strings.Join(valuesStrPlaceholder, ",") + ";"
	stmt, err := tx.Prepare(customSQL)
	if err != nil {
		return 0, err
	}
	result, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func addApplicantScores(
	ctx context.Context,
	tx *sql.Tx,
	payload AddProjectRequest,
	projectHistoryId int,
	criteriaList []ApplicantSelfScoreCriteria,
) (int64, error) {
	valuesStrPlaceholder := []string{}
	values := []any{}
	scores := payload.Details.Score

	for i := 0; i < len(criteriaList); i++ {
		valuesStrPlaceholder = append(valuesStrPlaceholder, fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3))
		scoreName := fmt.Sprintf("q_%d_%d", criteriaList[i].CriteriaVersion, criteriaList[i].OrderNumber)
		score, exist := scores[scoreName]
		if !exist {
			return 0, fmt.Errorf("applicant score %s is not exist", scoreName)
		}
		values = append(values, projectHistoryId, criteriaList[i].Id, score)
	}

	customSQL := addManyApplicantScoreSQL + strings.Join(valuesStrPlaceholder, ",") + ";"
	stmt, err := tx.Prepare(customSQL)
	if err != nil {
		return 0, err
	}
	result, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func addGeneralAddress(ctx context.Context, tx *sql.Tx, payload AddProjectRequest) (int, error) {
	var id int
	err := tx.QueryRowContext(ctx, addAddressSQL, payload.General.Address.Address, payload.General.Address.PostcodeId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func addProjectContactPersonAddress(ctx context.Context, tx *sql.Tx, payload AddProjectRequest, contactType string) (int, error) {
	var id int
	var addressDesc string
	var postcodeId int
	if contactType == "projectHead" {
		addressDesc = payload.Contact.ProjectHead.Address.Address
		postcodeId = payload.Contact.ProjectHead.Address.PostcodeId
	} else if contactType == "projectManager" {
		addressDesc = payload.Contact.ProjectManager.Address.Address
		postcodeId = payload.Contact.ProjectManager.Address.PostcodeId
	} else if contactType == "projectCoordinator" {
		addressDesc = payload.Contact.ProjectCoordinator.Address.Address
		postcodeId = payload.Contact.ProjectCoordinator.Address.PostcodeId
	} else {
		return 0, errors.New("unsupported contactType")
	}
	err := tx.QueryRowContext(
		ctx,
		addAddressSQL,
		addressDesc,
		postcodeId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func addProjectHeadContactId(ctx context.Context, tx *sql.Tx, payload AddProjectRequest, projectHeadAddressId int) (int, error) {
	var id int
	err := tx.QueryRowContext(
		ctx,
		addContactFullSQL,
		payload.Contact.ProjectHead.Prefix,
		payload.Contact.ProjectHead.FirstName,
		payload.Contact.ProjectHead.LastName,
		payload.Contact.ProjectHead.OrganizationPosition,
		payload.Contact.ProjectHead.EventPosition,
		projectHeadAddressId,
		payload.Contact.ProjectHead.Email,
		payload.Contact.ProjectHead.LineId,
		payload.Contact.ProjectHead.PhoneNumber,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func addProjectManagerContactId(ctx context.Context, tx *sql.Tx, payload AddProjectRequest, projectManagerAddressId int) (int, error) {
	var id int
	err := tx.QueryRowContext(
		ctx,
		addContactFullSQL,
		payload.Contact.ProjectManager.Prefix,
		payload.Contact.ProjectManager.FirstName,
		payload.Contact.ProjectManager.LastName,
		payload.Contact.ProjectManager.OrganizationPosition,
		payload.Contact.ProjectManager.EventPosition,
		projectManagerAddressId,
		payload.Contact.ProjectManager.Email,
		payload.Contact.ProjectManager.LineId,
		payload.Contact.ProjectManager.PhoneNumber,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func addProjectCoordinatorContactId(ctx context.Context, tx *sql.Tx, payload AddProjectRequest, projectCoordinatorAddressId int) (int, error) {
	var id int
	err := tx.QueryRowContext(
		ctx,
		addContactFullSQL,
		payload.Contact.ProjectCoordinator.Prefix,
		payload.Contact.ProjectCoordinator.FirstName,
		payload.Contact.ProjectCoordinator.LastName,
		payload.Contact.ProjectCoordinator.OrganizationPosition,
		payload.Contact.ProjectCoordinator.EventPosition,
		projectCoordinatorAddressId,
		payload.Contact.ProjectCoordinator.Email,
		payload.Contact.ProjectCoordinator.LineId,
		payload.Contact.ProjectCoordinator.PhoneNumber,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func addProjectRaceDirectorContactId(ctx context.Context, tx *sql.Tx, payload AddProjectRequest) (int, error) {
	var id int
	err := tx.QueryRowContext(
		ctx,
		addContactOnlyRequiredParamSQL,
		payload.Contact.RaceDirector.Alternative.Prefix,
		payload.Contact.RaceDirector.Alternative.FirstName,
		payload.Contact.RaceDirector.Alternative.LastName,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func addProjectRow(ctx context.Context, tx *sql.Tx, projectCode string, now time.Time, projectHistoryId int, userId int) (int, error) {
	var id int
	err := tx.QueryRowContext(
		ctx,
		addProjectSQL,
		projectCode,
		now,
		projectHistoryId,
		userId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func failAdd(name string, err error) (int, error) {
	return 0, fmt.Errorf("addProject name: %s, error: %w", name, err)
}

func buildTimeFromPayload(payload AddProjectRequest) (time.Time, time.Time, time.Time, error) {
	loc, err := getTimeLocation()
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, err
	}

	fromDate := time.Date(
		payload.General.EventDate.Year,
		time.Month(payload.General.EventDate.Month),
		payload.General.EventDate.Day,
		*payload.General.EventDate.FromHour,
		*payload.General.EventDate.FromMinute,
		0,
		0,
		loc,
	)
	toDate := time.Date(
		payload.General.EventDate.Year,
		time.Month(payload.General.EventDate.Month),
		payload.General.EventDate.Day,
		*payload.General.EventDate.ToHour,
		*payload.General.EventDate.ToMinute,
		0,
		0,
		loc,
	)

	var thisSeriesLatestDate time.Time
	if !(*payload.Experience.ThisSeries.FirstTime) {
		thisSeriesLatestDate = time.Date(
			payload.Experience.ThisSeries.History.Year,
			time.Month(payload.Experience.ThisSeries.History.Month),
			payload.Experience.ThisSeries.History.Day,
			0,
			0,
			0,
			0,
			loc,
		)
	}
	return fromDate, toDate, thisSeriesLatestDate, nil
}

func (s *store) handleCreateProjectFiles(baseFilePrefix string, userId int, projectCode string, payload AddProjectRequest, attachments []Attachments) error {
	// Write users uploaded file to zip files
	zipTmpPath := filepath.Join("../home", fmt.Sprintf("tmp/%s", baseFilePrefix))
	err := os.MkdirAll(zipTmpPath, os.ModePerm)
	if err != nil {
		return err
	}

	attachmentsZipName := fmt.Sprintf("%s_%s.zip", projectCode, attachmentsStr)

	baseZipPath := filepath.Join("../home", fmt.Sprintf("tmp/%s", baseFilePrefix))
	attachmentsZip, err := os.Create(fmt.Sprintf("%s/%s", baseZipPath, attachmentsZipName))
	if err != nil {
		return err
	}
	defer attachmentsZip.Close()
	attachmentsZipWriter := zip.NewWriter(attachmentsZip)
	defer attachmentsZipWriter.Close()

	zipWriterMap := map[string][]*zip.Writer{}
	var collaborationZip *os.File
	var collaborationZipWriter *zip.Writer
	var collaborationZipName string
	if *payload.Collaborated {
		collaborationZipName = fmt.Sprintf("%s_%s.zip", projectCode, collaborationStr)
		collaborationZip, err = os.Create(fmt.Sprintf("%s/%s", baseZipPath, collaborationZipName))
		if err != nil {
			return err
		}
		defer collaborationZip.Close()
		collaborationZipWriter = zip.NewWriter(collaborationZip)
		defer collaborationZipWriter.Close()

		zipWriterMap[collaborationStr] = []*zip.Writer{collaborationZipWriter, attachmentsZipWriter}
	}

	zipWriterMap[attachmentsStr] = []*zip.Writer{attachmentsZipWriter}

	// generate pdf files
	pdfPath, err := s.generateApplicantFormPdf(
		userId,
		projectCode,
		payload,
	)
	if err != nil {
		slog.Error("error generating a pdf for", "projectCode", projectCode)
		return err
	}

	// write pdf file to attachments zip and form zip
	formPdfFile, err := os.Open(pdfPath)
	if err != nil {
		return err
	}
	err = utils.WriteToZip(attachmentsZipWriter, formPdfFile, fmt.Sprintf("%s_แบบฟอร์ม.pdf", projectCode))
	if err != nil {
		return err
	}
	_, err = formPdfFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	bucketName := os.Getenv("AWS_S3_STORE_BUCKET_NAME")
	err = s.awsS3Service.DoUploadFileToS3(formPdfFile, bucketName, fmt.Sprintf("%s/%s_แบบฟอร์ม.pdf", baseFilePrefix, projectCode))
	if err != nil {
		return err
	}

	for _, attachment := range attachments {
		zipWriters := zipWriterMap[attachment.ZipName]
		s3FilePrefix := fmt.Sprintf("%s/%s", baseFilePrefix, attachment.DirName)
		err = s.awsS3Service.ZipAndUploadFileToS3(attachment.Files, zipWriters, fmt.Sprintf("%s_%s", projectCode, attachment.InZipFilePrefix), s3FilePrefix)
		if err != nil {
			return err
		}
	}
	// Write pdf to attachmentZip writer and form formZip writer
	// Then upload non-zipped pdf file to s3

	// close zip writer before upload to s3
	attachmentsZipWriter.Close()
	s3ZipPrefix := fmt.Sprintf("%s/zip", baseFilePrefix)
	if *payload.Collaborated {
		collaborationZipWriter.Close()
		_, err = collaborationZip.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		err = s.awsS3Service.DoUploadFileToS3(collaborationZip, bucketName, fmt.Sprintf("%s/%s", s3ZipPrefix, collaborationZipName))
		if err != nil {
			return err
		}
	}

	_, err = attachmentsZip.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = s.awsS3Service.DoUploadFileToS3(attachmentsZip, bucketName, fmt.Sprintf("%s/%s", s3ZipPrefix, attachmentsZipName))
	if err != nil {
		return err
	}

	// Clean up temp files
	err = os.RemoveAll(filepath.Join(zipTmpPath, ".."))
	if err != nil {
		return err
	}
	err = os.RemoveAll(pdfPath)
	if err != nil {
		return err
	}
	return nil
}
