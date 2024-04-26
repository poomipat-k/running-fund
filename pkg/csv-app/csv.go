package myCsv

import (
	"bytes"
	"encoding/csv"
	"errors"
)

func GenCsvBuffer(records [][]string) (*bytes.Buffer, error) {
	if len(records) == 0 {
		return nil, errors.New("records cannot be empty")
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.WriteAll(records)
	if err != nil {
		return nil, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return &buf, nil
}
