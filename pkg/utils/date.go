package utils

import "time"

const TIMEZONE = "Asia/Bangkok"

var MonthMapThai = map[int]string{
	1:  "มกราคม",
	2:  "กุมภาพันธ์",
	3:  "มีนาคม",
	4:  "เมษายน",
	5:  "พฤษภาคม",
	6:  "มิถุนายน",
	7:  "กรกฎาคม",
	8:  "สิงหาคม",
	9:  "กันยายน",
	10: "ตุลาคม",
	11: "พฤศจิกายน",
	12: "ธันวาคม",
}

func GetTimeLocation() (*time.Location, error) {
	loc, err := time.LoadLocation(TIMEZONE)
	if err != nil {
		return &time.Location{}, err
	}
	return loc, nil
}
