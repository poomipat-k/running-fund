package projects

import (
	"fmt"
	"log/slog"

	"github.com/jung-kurt/gofpdf"
)

const fontSize = 16
const lineHt = 1.25 * fontSize
const padding = 72
const fontName = "sarabunnew"

func generateApplicantFormPdf(userId int, projectCode, projectName, eventDate, address, subdistrict, district, province string) (string, error) {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, "A4", "")
	w, h := pdf.GetPageSize()
	pdf.AddUTF8Font(fontName, "", "../home/fonts/THSarabunNew.ttf")
	fmt.Printf("width=%v, height=%v\n", w, h)

	// pdf.UnitToPointConvert()

	pdf.SetMargins(padding, padding, padding)

	pdf.SetFont(fontName, "", fontSize)
	pdf.SetAutoPageBreak(true, padding)
	pdf.AddPage()

	pdf.Ln(-1)

	pdf.MultiCell(0, lineHt, "แบบฟอร์มข้อเสนอโครงการสนับสนุนทุนอุปถัมภ์งานเดิน-วิ่งเพื่อสุขภาพ", gofpdf.BorderNone, gofpdf.AlignCenter, false)
	pdf.Ln(8)

	pdf.MultiCell(0, lineHt, " ส่วนที่ 1 ข้อมูลพื้นฐานโครงการ", gofpdf.BorderFull, gofpdf.AlignLeft, false)
	pdf.Ln(8)

	pdf.MultiCell(0, lineHt, fmt.Sprintf("    1.1 ชื่อโครงการ/งานวิ่งเพื่อสุขภาพ: %s", projectName), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.MultiCell(0, lineHt, fmt.Sprintf("    1.2 วันที่จัดกิจกรรม (วันที่/เดือน/ปี):  %s", eventDate), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.MultiCell(0, lineHt, fmt.Sprintf("    1.3 สถานที่จัดกิจกรรม: ณ %s    ตำบล/แขวง %s    อำเภอ/เขต %s    จังหวัด %s",
		address, subdistrict, district, province),
		gofpdf.BorderNone,
		gofpdf.AlignLeft,
		false,
	)

	targetPath := fmt.Sprintf("../home/tmp/pdf/user_%d_%s.pdf", userId, projectCode)
	err := pdf.OutputFileAndClose(targetPath)
	if err != nil {
		slog.Error("error saving a pdf file to a local file", "error", err.Error())
		return "", err
	}
	fmt.Println("== Done ==")
	return targetPath, nil

}
