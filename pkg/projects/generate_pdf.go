package projects

import (
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const padding = 72
const sr = "sarabunnew"
const srB = "sarabunnewBold"

func (s *store) generateApplicantFormPdf(userId int, projectCode string, payload AddProjectRequest) (string, error) {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, "A4", "")
	w, h := pdf.GetPageSize()
	pdf.AddUTF8Font(sr, "", "../home/fonts/THSarabunNew.ttf")
	pdf.AddUTF8Font(srB, "B", "../home/fonts/THSarabunNewBold.ttf")
	fmt.Printf("width=%v, height=%v\n", w, h)

	// pdf.UnitToPointConvert()

	pdf.SetMargins(padding, padding, padding)
	pdf.SetAutoPageBreak(true, padding)
	pdf.AddPage()

	// Header start
	pdf.ImageOptions("../home/images/sss.png", 52, 16, 52, 0, false, gofpdf.ImageOptions{
		ReadDpi: true,
	}, 0, "")

	pdf.ImageOptions("../home/images/run_club.png", 120, 10, 50, 0, false, gofpdf.ImageOptions{
		ReadDpi: true,
	}, 0, "")

	pdf.SetFont(srB, "B", 16)
	pdf.Text(200, 30, fmt.Sprintf("ข้อเสนอโครงการ: %s", payload.General.ProjectName))
	pdf.Text(200, 50, fmt.Sprintf("รหัสโครงการ: %s", projectCode))

	// Header end

	pdf.SetFont(srB, "B", 18)
	pdf.MultiCell(0, 18, "ข้อมูลทั่วไปโครงการ", gofpdf.BorderNone, gofpdf.AlignCenter, false)
	pdf.Ln(12)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "ส่วนที่ 1 ข้อมูลพื้นฐานโครงการ", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, "1.1 ชื่อโครงการ:", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(0, 16, indent(payload.General.ProjectName, 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)

	pdf.MultiCell(0, 16, "1.2 วันที่จัดงานวิ่ง:", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)

	loc, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		return "", err
	}
	fromTime := convertEventDateToDateTime(
		payload.General.EventDate.Year,
		payload.General.EventDate.Month,
		payload.General.EventDate.Day,
		*payload.General.EventDate.FromHour,
		*payload.General.EventDate.FromMinute,
		loc,
	)

	log.Println("==1", fromTime.String())
	log.Println("==2", fromTime.Local().String())
	pdf.MultiCell(0, 16, indent(fromTime.String(), 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)

	targetPath := fmt.Sprintf("../home/tmp/pdf/user_%d_%s.pdf", userId, projectCode)
	err = pdf.OutputFileAndClose(targetPath)
	if err != nil {
		slog.Error("error saving a pdf file to a local file", "error", err.Error())
		return "", err
	}
	fmt.Println("== Done ==")
	return targetPath, nil
}

func indent(input string, n int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", n), input)
}

func convertEventDateToDateTime(year, month, day, hour, minute int, loc *time.Location) time.Time {
	fromDate := time.Date(
		year,
		time.Month(month),
		day,
		hour,
		minute,
		0,
		0,
		loc,
	)
	return fromDate
}
