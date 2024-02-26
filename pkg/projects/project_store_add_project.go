package projects

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"
)

func (s *store) AddProject(
	payload AddProjectRequest,
	userId int,
	criteria []ApplicantSelfScoreCriteria,
	attachments []DetailsFiles,
) (int, error) {
	projectCode, err := s.generateProjectCode()
	if err != nil {
		return 0, err
	}
	baseFilePrefix := getBasePrefix(userId, projectCode)
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

	projectCoordinatorAddressId, err := addProjectCoordinatorAddress(ctx, tx, payload)
	if err != nil {
		return failAdd("projectCoordinatorAddressId", err)
	}

	// Add contact rows
	projectHeadContactId, err := addProjectHeadContactId(ctx, tx, payload)
	if err != nil {
		return failAdd("projectHeadContactId", err)
	}

	projectManagerContactId, err := addProjectManagerContactId(ctx, tx, payload)
	if err != nil {
		return failAdd("projectManagerContactId", err)
	}

	projectCoordinatorContactId, err := addProjectCoordinatorContactId(ctx, tx, payload, projectCoordinatorAddressId)
	if err != nil {
		return failAdd("projectCoordinatorContactId", err)
	}

	var projectRaceDirectorContactId int
	if payload.Contact.RaceDirector.Who == "other" {
		projectRaceDirectorContactId, err = addProjectRaceDirectorContactId(ctx, tx, payload)
		if err != nil {
			return failAdd("projectRaceDirectorContactId", err)
		}
	}
	// Add project
	projectId, err := addProjectRow(ctx, tx, payload, projectCode, now, userId)
	if err != nil {
		return failAdd("projectId", err)
	}
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
	log.Println("==projectHistoryId", projectHistoryId)

	// Add distance
	distanceRowsAffected, err := addDistances(ctx, tx, payload, projectHistoryId)
	if err != nil {
		return failAdd("distanceRowsAffected", err)
	}
	log.Println("==distanceRowsAffected", distanceRowsAffected)
	// Add applicant scores
	applicantScoreRowsAffected, err := addApplicantScores(ctx, tx, payload, projectHistoryId, criteria)
	if err != nil {
		return failAdd("applicantScoreRowsAffected", err)
	}
	log.Println("==applicantScoreRowsAffected", applicantScoreRowsAffected)
	// // update project.project_history_id
	// err = updateProjectFK(ctx, tx, projectHistoryId, projectId)
	// if err != nil {
	// 	return failAdd("update project.project_history_id", err)
	// }

	// // upload files
	// for _, files := range attachments {
	// 	err = s.awsS3Service.UploadToS3(files.Files, fmt.Sprintf("%s/%s", baseFilePrefix, files.DirName))
	// 	if err != nil {
	// 		slog.Error("Failed to upload files to s3", "dirName", files.DirName, "error", err.Error())
	// 		return 0, err
	// 	}
	// }

	err = tx.Commit()
	if err != nil {
		return failAdd("tx.Commit()", err)
	}
	// commit
	return projectId, nil
}

func updateProjectFK(ctx context.Context, tx *sql.Tx, projectHistoryId int, projectId int) error {
	log.Println("===projectHistoryId", projectHistoryId)
	log.Println("===projectId", projectId)
	// Debug
	rows, err := tx.QueryContext(ctx, "SELECT id, project_code FROM project WHERE id=$1 LIMIT 1;", projectId)
	if err != nil {
		log.Panic("====panic 1", err)
	}
	data := []struct {
		Id               int
		ProjectCode      string
		ProjectHistoryId int
	}{}
	for rows.Next() {
		var row struct {
			Id               int
			ProjectCode      string
			ProjectHistoryId int
		}
		err := rows.Scan(&row.Id, &row.ProjectCode)
		if err != nil {
			log.Panic("====panic 2", err)
		}

		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		log.Panic("===panic 3", err)
	}
	log.Println("==reach here")
	log.Println(data[0])

	// prjHist
	rows, err = tx.QueryContext(ctx, "SELECT id, project_code, project_id FROM project_history WHERE id=$1 LIMIT 1;", projectHistoryId)
	if err != nil {
		log.Panic("====panic 4", err)
	}
	data2 := []struct {
		Id          int
		ProjectCode string
		ProjectId   sql.NullInt64
	}{}
	for rows.Next() {
		var row struct {
			Id          int
			ProjectCode string
			ProjectId   sql.NullInt64
		}
		err := rows.Scan(&row.Id, &row.ProjectCode, &row.ProjectId)
		if err != nil {
			log.Panic("====panic 5", err)
		}

		data2 = append(data2, row)
	}
	err = rows.Err()
	if err != nil {
		log.Panic("===panic 6", err)
	}
	log.Println("==reach here")
	log.Println(data2[0])

	result, err := tx.ExecContext(ctx, "SET CONSTRAINTS ALL DEFERRED;")
	if err != nil {
		log.Panic("===Panic7 ", err)
	}
	log.Println("===result:", result)

	// END debug
	var modifiedId int
	err = tx.QueryRowContext(ctx, updateProjectProjectHistoryIdSQL, projectHistoryId, projectId).Scan(&modifiedId)
	log.Println("===modifiedId", modifiedId)
	if err != nil {
		return err
	}
	log.Println()
	return nil
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

	var id int
	err = tx.QueryRowContext(
		ctx,
		addProjectHistorySQL,
		projectCode,
		1,
		now,
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
		thisSeriesLatestCompletedDate,
		payload.Experience.ThisSeries.History.Completed1.Year,
		payload.Experience.ThisSeries.History.Completed1.Name,
		payload.Experience.ThisSeries.History.Completed1.Participant,
		payload.Experience.ThisSeries.History.Completed2.Year,
		payload.Experience.ThisSeries.History.Completed2.Name,
		payload.Experience.ThisSeries.History.Completed2.Participant,
		payload.Experience.ThisSeries.History.Completed3.Year,
		payload.Experience.ThisSeries.History.Completed3.Name,
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

	valuesStrStatement := []string{}
	values := []any{}

	for i := 0; i < len(checkedDistances); i++ {
		valuesStrStatement = append(valuesStrStatement, fmt.Sprintf("($%d, $%d, $%d, $%d)", 4*i+1, 4*i+2, 4*i+3, 4*i+4))
		values = append(values, checkedDistances[i].Type, checkedDistances[i].Fee, checkedDistances[i].Dynamic, projectHistoryId)
	}
	customSQL := addManyDistanceSQL + strings.Join(valuesStrStatement, ",") + ";"
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
	valuesStrStatement := []string{}
	values := []any{}
	scores := payload.Details.Score

	for i := 0; i < len(criteriaList); i++ {
		valuesStrStatement = append(valuesStrStatement, fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3))
		scoreName := fmt.Sprintf("q_%d_%d", criteriaList[i].CriteriaVersion, criteriaList[i].OrderNumber)
		score, exist := scores[scoreName]
		if !exist {
			return 0, fmt.Errorf("applicant score %s is not exist", scoreName)
		}
		values = append(values, projectHistoryId, criteriaList[i].Id, score)
	}

	customSQL := addManyApplicantScoreSQL + strings.Join(valuesStrStatement, ",") + ";"
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

func addProjectCoordinatorAddress(ctx context.Context, tx *sql.Tx, payload AddProjectRequest) (int, error) {
	var id int
	err := tx.QueryRowContext(
		ctx,
		addAddressSQL,
		payload.Contact.ProjectCoordinator.Address.Address,
		payload.Contact.ProjectCoordinator.Address.PostcodeId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func addProjectHeadContactId(ctx context.Context, tx *sql.Tx, payload AddProjectRequest) (int, error) {
	var id int
	err := tx.QueryRowContext(
		ctx,
		addContactMainSQL,
		payload.Contact.ProjectHead.Prefix,
		payload.Contact.ProjectHead.FirstName,
		payload.Contact.ProjectHead.LastName,
		payload.Contact.ProjectHead.OrganizationPosition,
		payload.Contact.ProjectHead.EventPosition,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func addProjectManagerContactId(ctx context.Context, tx *sql.Tx, payload AddProjectRequest) (int, error) {
	var id int
	err := tx.QueryRowContext(
		ctx,
		addContactMainSQL,
		payload.Contact.ProjectManager.Prefix,
		payload.Contact.ProjectManager.FirstName,
		payload.Contact.ProjectManager.LastName,
		payload.Contact.ProjectManager.OrganizationPosition,
		payload.Contact.ProjectManager.EventPosition,
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

func addProjectRow(ctx context.Context, tx *sql.Tx, payload AddProjectRequest, projectCode string, now time.Time, userId int) (int, error) {
	var id int
	err := tx.QueryRowContext(
		ctx,
		addProjectSQL,
		projectCode,
		now,
		userId,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	log.Println("===AddProjectRow id:", id)
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
