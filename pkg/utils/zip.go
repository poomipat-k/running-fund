package utils

import (
	"archive/zip"
	"io"
	"mime/multipart"
)

// Write to zip then seek start of the file
func WriteToZip(zipWriter *zip.Writer, file multipart.File, filePath string) error {
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
