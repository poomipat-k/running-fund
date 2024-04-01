package projects

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/poomipat-k/running-fund/pkg/utils"
)

const padding = 72
const sr = "sarabunnew"
const srB = "sarabunnewBold"

var expectedParticipantsMap = map[string]string{
	"<=500":     "ต่ำกว่า 500 คน",
	"501-1500":  "501 - 1,500 คน",
	"1501-2500": "1,501 - 2,500 คน",
	"2501-3500": "2,501 - 3,500 คน",
	"3501-4500": "3,501 - 4,500 คน",
	"4501-5500": "4,501 - 5,500 คน",
	">=5501":    "5,501 คนขึ้นไป",
}

var scoreMeaning = []string{"", "ไม่มั่นใจอย่างยิ่ง", "ไม่มั่นใจ", "กลาง ๆ", "มั่นใจ", "มั่นใจอย่างยิ่ง"}

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

	pdf.ImageOptions("../home/images/run_club.png", 120, 10, 40, 0, false, gofpdf.ImageOptions{
		ReadDpi: true,
	}, 0, "")

	pdf.SetFont(srB, "B", 16)
	pdf.Text(200, 30, fmt.Sprintf("ข้อเสนอโครงการ: %s", payload.General.ProjectName))
	pdf.Text(200, 50, fmt.Sprintf("รหัสโครงการ: %s", projectCode))

	// Header end
	// 1. General details
	err := s.generateGeneralDetailsSection(pdf, payload)
	if err != nil {
		return "", err
	}

	// 2. Contact
	err = s.generateContactSection(pdf, payload)
	if err != nil {
		return "", err
	}

	// 3. Details
	s.generateDetailsSection(pdf, payload)

	// 4. Experience
	s.generateExperienceSection(pdf, payload)

	// 5. Fund request
	s.generateFundRequestSection(pdf, payload)

	// save pdf to a file
	tmpPdfFolder := filepath.Join("../home/tmp/pdf")
	err = os.MkdirAll(tmpPdfFolder, os.ModePerm)
	if err != nil {
		return "", err
	}
	targetPath := fmt.Sprintf("%s/user_%d_%s.pdf", tmpPdfFolder, userId, projectCode)
	err = pdf.OutputFileAndClose(targetPath)
	if err != nil {
		slog.Error("error saving a pdf file to a local file", "error", err.Error())
		return "", err
	}
	return targetPath, nil
}

func indent(input string, n int) string {
	return fmt.Sprintf("%s%s", strings.Repeat(" ", n), input)
}

func getDateTimeString(year, month, day, hour, minute int) string {
	monthStr := utils.MonthMapThai[month]
	return fmt.Sprintf("%d %s %d %02d:%02d", day, monthStr, year, hour, minute)
}

func getDateString(year, month, day int) string {
	monthStr := utils.MonthMapThai[month]
	return fmt.Sprintf("%d %s %d", day, monthStr, year)
}

func (s *store) getAddressDetails(addressId int) (AddressDetails, error) {
	var ad AddressDetails
	row := s.db.QueryRow(getAddressDetailsSQL, addressId)
	err := row.Scan(&ad.Postcode, &ad.SubdistrictName, &ad.DistrictName, &ad.ProvinceName)

	switch err {
	case sql.ErrNoRows:
		slog.Error("getAddressDetails(): no row were returned!")
		return AddressDetails{}, err
	case nil:
		return ad, nil
	default:
		slog.Error(err.Error())
		return AddressDetails{}, err
	}
}

func (s *store) generateGeneralDetailsSection(pdf *gofpdf.Fpdf, payload AddProjectRequest) error {
	pdf.SetFont(srB, "B", 18)
	pdf.MultiCell(0, 18, "ข้อมูลทั่วไปโครงการ", gofpdf.BorderNone, gofpdf.AlignCenter, false)
	pdf.Ln(12)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "ส่วนที่ 1 ข้อมูลพื้นฐานโครงการ", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, "1.1 ชื่อโครงการ:", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(0, 16, indent(payload.General.ProjectName, 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "1.2 วันที่จัดงานวิ่ง:", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	fromTimeStr := getDateTimeString(
		payload.General.EventDate.Year,
		payload.General.EventDate.Month,
		payload.General.EventDate.Day,
		*payload.General.EventDate.FromHour,
		*payload.General.EventDate.FromMinute,
	)
	pdf.MultiCell(
		0,
		16,
		indent(fmt.Sprintf("%s - %02d:%02d", fromTimeStr, *payload.General.EventDate.ToHour, *payload.General.EventDate.ToMinute), 6),
		gofpdf.BorderNone,
		gofpdf.AlignLeft,
		false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "1.3 สถานที่จัดกิจกรรม", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	address, err := s.getAddressDetails(payload.General.Address.PostcodeId)
	if err != nil {
		return err
	}
	pdf.MultiCell(
		0,
		16,
		indent(fmt.Sprintf("%s,  %s,  %s,  %s", payload.General.Address.Address, address.SubdistrictName, address.DistrictName, address.ProvinceName), 6),
		gofpdf.BorderNone,
		gofpdf.AlignLeft,
		false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "1.4 เส้นทางที่ใช้สำหรับจัดงานวิ่ง", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(
		0,
		16,
		indent(fmt.Sprintf("จุด Start: %s\n      จุด Finish: %s", payload.General.StartPoint, payload.General.FinishPoint), 6),
		gofpdf.BorderNone,
		gofpdf.AlignLeft,
		false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "1.5 ประเภทระยะทางวิ่งและอัตราค่าสมัคร", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, indent("1.5.1 ประเภทการจัดวิ่ง", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if payload.General.EventDetails.Category.Available.RoadRace {
		pdf.MultiCell(0, 16, indent("- วิ่งถนน (Road Race)", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}

	if payload.General.EventDetails.Category.Available.TrailRunning {
		pdf.MultiCell(0, 16, indent("- Trail Running (การวิ่งตามภูมิประเทศ)", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}

	if payload.General.EventDetails.Category.Available.Other && payload.General.EventDetails.Category.OtherType != "" {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- ประเภทอื่นๆ: %s", payload.General.EventDetails.Category.OtherType), 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("1.5.2 ระยะทางและอัตราค่าสมัครปกติ", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	for _, daf := range payload.General.EventDetails.DistanceAndFee {
		if daf.Checked {
			if daf.Type == "fun" {
				pdf.MultiCell(
					0,
					16,
					indent(fmt.Sprintf("- Fun run (ระยะทางไม่เกิน 10 km)  ค่าสมัคร %.0f บาท", *daf.Fee), 10),
					gofpdf.BorderNone,
					gofpdf.AlignLeft,
					false,
				)
			} else if daf.Type == "mini" {
				pdf.MultiCell(
					0,
					16,
					indent(fmt.Sprintf("- Mini Marathon (ระยะทาง 10 km)  ค่าสมัคร %.0f บาท", *daf.Fee), 10),
					gofpdf.BorderNone,
					gofpdf.AlignLeft,
					false,
				)
			} else if daf.Type == "half" {
				pdf.MultiCell(
					0,
					16,
					indent(fmt.Sprintf("- Half Marathon (ระยะทาง 21.1 km)  ค่าสมัคร %.0f บาท", *daf.Fee), 10),
					gofpdf.BorderNone,
					gofpdf.AlignLeft,
					false,
				)
			} else if daf.Type == "full" {
				pdf.MultiCell(
					0,
					16,
					indent(fmt.Sprintf("- Marathon (ระยะทาง 42.195 km)  ค่าสมัคร %.0f บาท", *daf.Fee), 10),
					gofpdf.BorderNone,
					gofpdf.AlignLeft,
					false,
				)
			} else if *daf.Dynamic {
				pdf.MultiCell(
					0,
					16,
					indent(fmt.Sprintf("- %s  ค่าสมัคร %.0f บาท", daf.Type, *daf.Fee), 10),
					gofpdf.BorderNone,
					gofpdf.AlignLeft,
					false,
				)
			}
		}
	}
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	var vipText string
	if *payload.General.EventDetails.VIP {
		vipText = "- มี"
	} else {
		vipText = "- ไม่มี"
	}
	pdf.MultiCell(0, 16, indent("1.5.2 การเปิดรับสมัครประเภท VIP", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(0, 16, indent(vipText, 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "1.6 จำนวนผู้เข้าร่วมที่ตั้งเป้า", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(0, 16, indent(expectedParticipantsMap[payload.General.ExpectedParticipants], 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "1.7 การใช้บริษัทจัดงาน (Organizer)", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	var organizerText string
	if *payload.General.HasOrganizer {
		organizerText = payload.General.OrganizerName
	} else {
		organizerText = "ไม่ใช้"
	}
	pdf.MultiCell(0, 16, indent(organizerText, 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	return nil
}

func (s *store) generateContactSection(pdf *gofpdf.Fpdf, payload AddProjectRequest) error {
	pdf.Ln(12)
	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "ส่วนที่ 2 ข้อมูลการติดต่อ", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, "2.1 หัวหน้าโครงการ", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	projectHeadStr := indent(fmt.Sprintf("%s %s %s\n      ตำแหน่งในหน่วยงาน/องค์กร: %s\n      ตำแหน่งในการจัดงานครั้งนี้: %s",
		payload.Contact.ProjectHead.Prefix,
		payload.Contact.ProjectHead.FirstName,
		payload.Contact.ProjectHead.LastName,
		payload.Contact.ProjectHead.OrganizationPosition,
		payload.Contact.ProjectHead.EventPosition,
	), 6)
	pdf.MultiCell(0, 16,
		projectHeadStr,
		gofpdf.BorderNone, gofpdf.AlignLeft, false)

	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "2.2 ผู้รับผิดชอบโครงการ", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	projectManagerStr := indent(fmt.Sprintf("%s %s %s\n      ตำแหน่งในหน่วยงาน/องค์กร: %s\n      ตำแหน่งในการจัดงานครั้งนี้: %s",
		payload.Contact.ProjectManager.Prefix,
		payload.Contact.ProjectManager.FirstName,
		payload.Contact.ProjectManager.LastName,
		payload.Contact.ProjectManager.OrganizationPosition,
		payload.Contact.ProjectManager.EventPosition,
	), 6)
	pdf.MultiCell(0, 16,
		projectManagerStr,
		gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "2.3 ผู้ประสานงานโครงการ", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	projectCoordinatorStr := indent(fmt.Sprintf("%s %s %s\n      ตำแหน่งในหน่วยงาน/องค์กร: %s\n      ตำแหน่งในการจัดงานครั้งนี้: %s",
		payload.Contact.ProjectCoordinator.Prefix,
		payload.Contact.ProjectCoordinator.FirstName,
		payload.Contact.ProjectCoordinator.LastName,
		payload.Contact.ProjectCoordinator.OrganizationPosition,
		payload.Contact.ProjectCoordinator.EventPosition,
	), 6)
	pdf.MultiCell(0, 16,
		projectCoordinatorStr,
		gofpdf.BorderNone, gofpdf.AlignLeft, false)

	address, err := s.getAddressDetails(payload.General.Address.PostcodeId)
	if err != nil {
		return err
	}
	pdf.MultiCell(
		0,
		16,
		indent(fmt.Sprintf("ที่อยู่: %s,  %s,  %s,  %s %d", payload.General.Address.Address, address.SubdistrictName, address.DistrictName, address.ProvinceName, address.Postcode), 6),
		gofpdf.BorderNone,
		gofpdf.AlignLeft,
		false)
	pdf.MultiCell(
		0,
		16,
		indent(fmt.Sprintf("อีเมล (E-mail): %s", payload.Contact.ProjectCoordinator.Email), 6),
		gofpdf.BorderNone,
		gofpdf.AlignLeft,
		false)
	pdf.MultiCell(
		0,
		16,
		indent(fmt.Sprintf("ไลน์ไอดี (Line ID): %s", payload.Contact.ProjectCoordinator.LineId), 6),
		gofpdf.BorderNone,
		gofpdf.AlignLeft,
		false)
	pdf.MultiCell(
		0,
		16,
		indent(fmt.Sprintf("หมายเลขโทรศัพท์: %s", payload.Contact.ProjectCoordinator.PhoneNumber), 6),
		gofpdf.BorderNone,
		gofpdf.AlignLeft,
		false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "2.4 ผู้ตัดสินชี้ขาด (Race Director)", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	var raceDirectorStr string
	if payload.Contact.RaceDirector.Who == "projectHead" {
		raceDirectorStr = projectHeadStr
	} else if payload.Contact.RaceDirector.Who == "projectManager" {
		raceDirectorStr = projectManagerStr
	} else if payload.Contact.RaceDirector.Who == "projectCoordinator" {
		raceDirectorStr = projectCoordinatorStr
	} else if payload.Contact.RaceDirector.Who == "other" {
		raceDirectorStr = indent(fmt.Sprintf("%s %s %s",
			payload.Contact.RaceDirector.Alternative.Prefix,
			payload.Contact.RaceDirector.Alternative.FirstName,
			payload.Contact.RaceDirector.Alternative.LastName,
		), 6)
	}
	pdf.MultiCell(0, 16, raceDirectorStr, gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "2.5 ประเภทและชื่อของหน่วยงาน/องค์กรที่เสนอโครงการ", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	var organizationStr string
	if payload.Contact.Organization.Type == "government" {
		organizationStr = fmt.Sprintf("ภาครัฐ: %s", payload.Contact.Organization.Name)
	} else if payload.Contact.Organization.Type == "private_sector" {
		organizationStr = fmt.Sprintf("ภาคเอกชน: %s", payload.Contact.Organization.Name)
	} else if payload.Contact.Organization.Type == "civil_society" {
		organizationStr = fmt.Sprintf("ภาคประชาสังคม: %s", payload.Contact.Organization.Name)
	}

	pdf.MultiCell(0, 16, indent(organizationStr, 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	return nil
}

func (s *store) generateDetailsSection(pdf *gofpdf.Fpdf, payload AddProjectRequest) error {
	pdf.Ln(12)
	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, "ส่วนที่ 3 ข้อมูลข้อเสนอโครงการ และแผนบริหารจัดการงานวิ่งเพื่อสุขภาพ", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, "3.1 รายละเอียดความเป็นมาและวัตถุประสงค์ของการจัดงาน", gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)
	pdf.MultiCell(0, 16, indent("3.1.1 ความเป็นมา", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(0, 16, indent(payload.Details.Background, 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.1.2 วัตถุประสงค์", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(0, 16, indent(payload.Details.Objective, 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.2 ช่องทางการสื่อสารเพื่อแจ้งรายละเอียดของกิจกรรมให้กับนักวิ่งหรือผู้ที่สนใจเข้าร่วม", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.MultiCell(0, 16, indent("3.2.1 ช่องทางสื่อสังคมออนไลน์ (Social Media)", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if payload.Details.Marketing.Online.Available.Facebook {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- Facebook: %s", payload.Details.Marketing.Online.HowTo.Facebook), 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Marketing.Online.Available.Website {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- เว็บไซต์: %s", payload.Details.Marketing.Online.HowTo.Website), 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Marketing.Online.Available.OnlinePage {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- เพจวิ่งออนไลน์: %s", payload.Details.Marketing.Online.HowTo.OnlinePage), 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Marketing.Online.Available.Other {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- ช่องทางออนไลน์อื่นๆ: %s", payload.Details.Marketing.Online.HowTo.Other), 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.2.2 สื่อสารผ่านบุคคล และ/หรือช่องทางอื่น ๆ ที่ไม่ใช้อินเทอร์เน็ต (Offline)", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if payload.Details.Marketing.Offline.Available.PR {
		pdf.MultiCell(0, 16, indent("- ส่งหนังสือให้กับหน่วยงาน/องค์กรอื่นช่วยประชาสัมพันธ์", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Marketing.Offline.Available.LocalOfficial {
		pdf.MultiCell(0, 16, indent("- ประชาสัมพันธ์ผ่านบุคคลในพื้นที่ เช่น กำนัน ผู้ใหญ่บ้าน อสม. ชมรมวิ่ง", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Marketing.Offline.Available.Booth {
		pdf.MultiCell(0, 16, indent("- การตั้งบูธประชาสัมพันธ์/ รับสมัคร", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Marketing.Offline.Available.Billboard {
		pdf.MultiCell(0, 16, indent("- กระจายสื่อในพื้นที่ เช่น ป้าย ไวนิล รถประชาสัมพันธ์", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Marketing.Offline.Available.TV {
		pdf.MultiCell(0, 16, indent("- การลงข่าวหรือโฆษณาทาง TV", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Marketing.Offline.Available.Other {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- ช่องทางออฟไลน์อื่นๆ: %s", payload.Details.Marketing.Offline.Addition), 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.3 ความมั่นใจในการวางแผนการจัดเตรียม อุปกรณ์ สถานที่ และสิ่งอำนวยความสะดวกให้นักวิ่ง", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)

	criteria, err := s.GetApplicantCriteriaForPDF(1)
	if err != nil {
		return err
	}
	for _, cri := range criteria {
		pdf.SetFont(srB, "B", 16)
		key := fmt.Sprintf("q_%d_%d", cri.CriteriaVersion, cri.OrderNumber)
		score := payload.Details.Score[key]
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("3.3.%d %s", cri.OrderNumber, cri.PdfDisplay), 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
		pdf.SetFont(sr, "", 16)
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("-  %s", scoreMeaning[score]), 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
		pdf.Ln(4)
	}

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.4 แผนการดูแลความปลอดภัยทางสุขภาพของนักวิ่ง", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if payload.Details.Safety.Ready.RunnerInformation {
		pdf.MultiCell(0, 16, indent("- ข้อมูลสุขภาพและหมายเลขโทรศัพท์ติดต่อฉุกเฉินของนักวิ่งในแบบฟอร์มลงทะเบียน/ระบบ/BIB", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Safety.Ready.HealthDecider {
		pdf.MultiCell(0, 16, indent("- กำหนดผู้รับผิดชอบ/ผู้ตัดสินใจเรื่องความปลอดภัยด้านสุขภาพ", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Safety.Ready.Ambulance {
		pdf.MultiCell(0, 16, indent("- รถพยาบาลฉุกเฉิน (ambulance) พร้อมแพทย์/พยาบาลเคลื่อนที่", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Safety.Ready.FirstAid {
		pdf.MultiCell(0, 16, indent("- จุดปฐมพยาบาลพร้อมเวชภัณฑ์ เช่น แอมโมเนีย ที่ติดแผล สเปรย์ฉีดคลายกล้ามเนื้อ", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Safety.Ready.AED {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- เครื่อง AED จำนวน %d เครื่อง", payload.Details.Safety.AEDCount), 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Safety.Ready.Insurance {
		pdf.MultiCell(0, 16, indent("- ประกันชีวิตสำหรับนักวิ่ง", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Safety.Ready.Other {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- อื่นๆ: %s", payload.Details.Safety.Addition), 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.5 การวัดระยะทางวิ่งและการจัดการจราจร", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, indent("3.5.1 เส้นทาง (เอกสารแนบ)", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, indent("3.5.2 การวัดระยะทาง", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if payload.Details.Route.Measurement.AthleticsAssociation {
		pdf.MultiCell(0, 16, indent("- ได้รับการวัดและรับรองจากสมาคมกีฬากรีฑาแห่งประเทศไทย", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Route.Measurement.CalibratedBicycle {
		pdf.MultiCell(0, 16, indent("- การวัดระยะทางด้วยจักรยานที่สอบเทียบ (Calibrated Bicycle)", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Route.Measurement.SelfMeasurement {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- ผู้จัดการแข่งขันวัดระยะทางเอง เครื่องมือ: %s", payload.Details.Route.Tool), 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	pdf.Ln(4)

	// checkbox
	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.5.3 การจัดการจราจร", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if payload.Details.Route.TrafficManagement.AskPermission {
		pdf.MultiCell(0, 16, indent("- ก่อนและระหว่างการจัดกิจกรรมมีการติดตั้งป้ายขออภัยในความไม่สะดวกในการใช้เส้นทาง", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Route.TrafficManagement.HasSupporter {
		pdf.MultiCell(0, 16, indent("- มีผู้ช่วยดูแลความปลอดภัย เช่น ตำรวจ อาสาสมัครในพื้นที่", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Route.TrafficManagement.RoadClosure {
		pdf.MultiCell(0, 16, indent("- ขออนุญาตหน่วยงานปิดถนน หรือแบ่งช่องทางการจราจร", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Route.TrafficManagement.Signs {
		pdf.MultiCell(0, 16, indent("- ตั้งป้ายสัญลักษณ์ เช่น ป้ายบอกระยะทาง ป้ายจุดบริการน้ำดื่ม", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Route.TrafficManagement.Lighting {
		pdf.MultiCell(0, 16, indent("- มีการจัดแสงไฟในเส้นทางวิ่ง ในช่วงเส้นทางมืด", 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	pdf.Ln(4)

	// radio
	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.6 ระบบตัดสินของผลของการจัดกิจกรรม", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if payload.Details.Judge.Type == "manual" {
		pdf.MultiCell(0, 16, indent("- ระบบ Manual ใช้กรรมการตัดสิน", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	} else if payload.Details.Judge.Type == "auto" {
		pdf.MultiCell(0, 16, indent("- ระบบ Auto (Chip time) ใช้เครื่องประมวลผลร่วมกับกรรมการตัดสินชี้ขาด", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	} else if payload.Details.Judge.Type == "other" {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- อื่นๆ: %s", payload.Details.Judge.OtherType), 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	pdf.Ln(4)

	// checkbox
	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.7 หน่วยงานระดับพื้นที่ที่ร่วมสนับสนุน", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if payload.Details.Support.Organization.ProvincialAdministration {
		pdf.MultiCell(0, 16, indent("- หน่วยงานด้านการปกครอง เช่น ผู้ว่าราชการ นายอำเภอ นายกเทศบาล ทต. อบต.", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Support.Organization.Safety {
		pdf.MultiCell(0, 16, indent("- หน่วยงานด้านความปลอดภัย เช่น ตำรวจ อปพร. วิทยุกู้ชีพ", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Support.Organization.Health {
		pdf.MultiCell(0, 16, indent("- หน่วยงานด้านการแพทย์ เช่น โรงพยาบาล รพ.สต. อสม.", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Support.Organization.Volunteer {
		pdf.MultiCell(0, 16, indent("- มูลนิธิ อาสาสมัครชุมชน", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Support.Organization.Community {
		pdf.MultiCell(0, 16, indent("- องค์กรระดับชุมชน เช่น โรงเรียน วัด ชุมชน อสม.", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	if payload.Details.Support.Organization.Other {
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- อื่นๆ: %s", payload.Details.Support.Addition), 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	}
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("3.8 วิธีการประเมินผลความสำเร็จ/รับฟังประเด็นท้าทายจากการจัดกิจกรรม", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(0, 16, indent(payload.Details.Feedback, 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)

	return nil
}

func (s *store) generateExperienceSection(pdf *gofpdf.Fpdf, payload AddProjectRequest) {
	pdf.Ln(12)
	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("ส่วนที่ 4 ประสบการณ์ดำเนินการโครงการ", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, indent("4.1 การจัดกิจกรรมวิ่งเพื่อสุขภาพที่ท่านส่งข้อเสนอโครงการในครั้งนี้เป็นการจัดกิจกรรมครั้งแรกหรือไม่", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if *payload.Experience.ThisSeries.FirstTime {
		pdf.MultiCell(0, 16, indent("- ครั้งแรก", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	} else {
		latestThisEventDate := getDateString(
			payload.Experience.ThisSeries.History.Year,
			payload.Experience.ThisSeries.History.Month,
			payload.Experience.ThisSeries.History.Day,
		)
		pdf.MultiCell(0, 16, indent(fmt.Sprintf("- ครั้งนี้ครั้งที่ %d  จัดครั้งล่าสุดเมื่อวันที่ %s", payload.Experience.ThisSeries.History.OrdinalNumber, latestThisEventDate), 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
		pdf.MultiCell(0, 16, indent("การจัดงานครั้งที่ผ่านมา", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
		if payload.Experience.ThisSeries.History.Completed1.Name != "" {
			pdf.MultiCell(
				0,
				16,
				indent(fmt.Sprintf("- ปีที่จัดงาน %d  ชื่องาน %s  จำนวนผู้เข้าร่วม %s คน",
					payload.Experience.ThisSeries.History.Completed1.Year,
					payload.Experience.ThisSeries.History.Completed1.Name,
					utils.FormatInt(int64(payload.Experience.ThisSeries.History.Completed1.Participant)),
				), 6),
				gofpdf.BorderNone,
				gofpdf.AlignLeft,
				false,
			)
		}
		if payload.Experience.ThisSeries.History.Completed2.Name != "" {
			pdf.MultiCell(
				0,
				16,
				indent(fmt.Sprintf("- ปีที่จัดงาน %d  ชื่องาน %s  จำนวนผู้เข้าร่วม %s คน",
					payload.Experience.ThisSeries.History.Completed2.Year,
					payload.Experience.ThisSeries.History.Completed2.Name,
					utils.FormatInt(int64(payload.Experience.ThisSeries.History.Completed2.Participant)),
				), 6),
				gofpdf.BorderNone,
				gofpdf.AlignLeft,
				false,
			)
		}
		if payload.Experience.ThisSeries.History.Completed3.Name != "" {
			pdf.MultiCell(
				0,
				16,
				indent(fmt.Sprintf("- ปีที่จัดงาน %d  ชื่องาน %s  จำนวนผู้เข้าร่วม %s คน",
					payload.Experience.ThisSeries.History.Completed3.Year,
					payload.Experience.ThisSeries.History.Completed3.Name,
					utils.FormatInt(int64(payload.Experience.ThisSeries.History.Completed3.Participant)),
				), 6),
				gofpdf.BorderNone,
				gofpdf.AlignLeft,
				false,
			)
		}
	}
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("4.2 นอกจากการจัดกิจกรรมตามข้อเสนอโครงการในครั้งนี้ ท่านมีประสบการณ์จัดงานวิ่งอื่น ๆ หรือไม่", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if !*payload.Experience.OtherSeries.DoneBefore {
		pdf.MultiCell(0, 16, indent("- ไม่มี", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	} else {
		pdf.MultiCell(0, 16, indent("การจัดงานครั้งที่ผ่านมา", 6), gofpdf.BorderNone, gofpdf.AlignLeft, false)
		if payload.Experience.OtherSeries.History.Completed1.Name != "" {
			pdf.MultiCell(
				0,
				16,
				indent(fmt.Sprintf("- ปีที่จัดงาน %d  ชื่องาน %s  จำนวนผู้เข้าร่วม %s คน",
					payload.Experience.OtherSeries.History.Completed1.Year,
					payload.Experience.OtherSeries.History.Completed1.Name,
					utils.FormatInt(int64(payload.Experience.OtherSeries.History.Completed1.Participant)),
				), 6),
				gofpdf.BorderNone,
				gofpdf.AlignLeft,
				false,
			)
		}
		if payload.Experience.OtherSeries.History.Completed2.Name != "" {
			pdf.MultiCell(
				0,
				16,
				indent(fmt.Sprintf("- ปีที่จัดงาน %d  ชื่องาน %s  จำนวนผู้เข้าร่วม %s คน",
					payload.Experience.OtherSeries.History.Completed2.Year,
					payload.Experience.OtherSeries.History.Completed2.Name,
					utils.FormatInt(int64(payload.Experience.OtherSeries.History.Completed2.Participant)),
				), 6),
				gofpdf.BorderNone,
				gofpdf.AlignLeft,
				false,
			)
		}
		if payload.Experience.OtherSeries.History.Completed3.Name != "" {
			pdf.MultiCell(
				0,
				16,
				indent(fmt.Sprintf("- ปีที่จัดงาน %d  ชื่องาน %s  จำนวนผู้เข้าร่วม %s คน",
					payload.Experience.OtherSeries.History.Completed3.Year,
					payload.Experience.OtherSeries.History.Completed3.Name,
					utils.FormatInt(int64(payload.Experience.OtherSeries.History.Completed3.Participant)),
				), 6),
				gofpdf.BorderNone,
				gofpdf.AlignLeft,
				false,
			)
		}
	}
}

func (s *store) generateFundRequestSection(pdf *gofpdf.Fpdf, payload AddProjectRequest) {
	pdf.Ln(12)
	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("ส่วนที่ 5 ความต้องการขอรับการสนับสนุน", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, indent("5.1 งบประมาณในภาพรวมและการสนับสนุนจากแหล่งอื่น", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.MultiCell(0, 16, indent("5.1.1 การจัดงานครั้งนี้ประมาณการงบประมาณจัดงานทั้งหมด (บาท)", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(0, 16, indent(fmt.Sprintf("-  %s บาท", utils.FormatInt(int64(payload.Fund.Budget.Total))), 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("5.1.2 หน่วยงาน/องค์กรที่ให้การสนับสนุนการจัดงานครั้งนี้", 8), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	pdf.MultiCell(0, 16, indent(payload.Fund.Budget.SupportOrganization, 10), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.Ln(4)

	pdf.SetFont(srB, "B", 16)
	pdf.MultiCell(0, 16, indent("5.2 ความต้องการการสนับสนุนจากสสส. และสมาพันธ์ฯ", 0), gofpdf.BorderNone, gofpdf.AlignLeft, false)
	pdf.SetFont(sr, "", 16)
	if payload.Fund.Request.Type.Fund {
		pdf.MultiCell(
			0,
			16,
			indent(fmt.Sprintf("- งบประมาณ %s บาท", utils.FormatInt(int64(payload.Fund.Request.Details.FundAmount))), 6),
			gofpdf.BorderNone,
			gofpdf.AlignLeft,
			false,
		)
	}
	if payload.Fund.Request.Type.BIB {
		pdf.MultiCell(
			0,
			16,
			indent(fmt.Sprintf("- เบอร์วิ่ง (BIB) จำนวน %s ใบ", utils.FormatInt(int64(payload.Fund.Request.Details.BibAmount))), 6),
			gofpdf.BorderNone,
			gofpdf.AlignLeft,
			false,
		)
	}
	if payload.Fund.Request.Type.Pr {
		pdf.MultiCell(
			0,
			16,
			indent("- การประชาสัมพันธ์กิจกรรม และการรับสมัคร", 6),
			gofpdf.BorderNone,
			gofpdf.AlignLeft,
			false,
		)
	}
	if payload.Fund.Request.Type.Seminar {
		pdf.MultiCell(
			0,
			16,
			indent(fmt.Sprintf("- การอบรม/การพัฒนาศักยภาพ หัวข้อเรื่อง %s", payload.Fund.Request.Details.Seminar), 6),
			gofpdf.BorderNone,
			gofpdf.AlignLeft,
			false,
		)
	}
	if payload.Fund.Request.Type.Other {
		pdf.MultiCell(
			0,
			16,
			indent(fmt.Sprintf("- อื่นๆ %s", payload.Fund.Request.Details.Other), 6),
			gofpdf.BorderNone,
			gofpdf.AlignLeft,
			false,
		)
	}
}
