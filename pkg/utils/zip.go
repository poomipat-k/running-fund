package utils

import (
	"archive/zip"
	"io"
)

// Write to zip then seek start of the file
func WriteToZip(zipWriter *zip.Writer, file io.ReadSeeker, filePath string) error {
	w, err := zipWriter.Create(filePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}
	// Seek start
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}
