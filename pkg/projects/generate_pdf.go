package projects

import (
	"fmt"
	"log/slog"

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

	pdf.SetFont(srB, "B", 18)
	pdf.MultiCell(0, 18, "ข้อมูลทั่วไปโครงการ", gofpdf.BorderNone, gofpdf.AlignCenter, false)

	targetPath := fmt.Sprintf("../home/tmp/pdf/user_%d_%s.pdf", userId, projectCode)
	err := pdf.OutputFileAndClose(targetPath)
	if err != nil {
		slog.Error("error saving a pdf file to a local file", "error", err.Error())
		return "", err
	}
	fmt.Println("== Done ==")
	return targetPath, nil

}
