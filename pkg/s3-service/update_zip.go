package s3Service

// // params userId, projectCode,
// func (client *S3Service) UpdateAttachmentZipContent(userId int, projectCode string, newFiles []*multipart.FileHeader) (*os.File, error) {
// 	folderPath := filepath.Join("../home/tmp/zip")
// 	err := os.MkdirAll(folderPath, os.ModePerm)
// 	if err != nil {
// 		return nil, err
// 	}

// 	zipFileTargetPath := filepath.Join(folderPath, fmt.Sprintf("%s.zip", projectCode))
// 	// Download current zip file from S3
// 	zipFile, err := client.downloadZipToLocal(zipFileTargetPath, userId, projectCode)
// 	if err != nil {
// 		slog.Error("error downloading zip file to local", "error", err)
// 		return nil, err
// 	}

// 	err = client.updateZipFile(zipFileTargetPath, newFiles, projectCode)
// 	if err != nil {
// 		slog.Error("error updating the zip file", "error", err)
// 		return nil, err
// 	}
// 	return zipFile, nil
// }

// func (client *S3Service) updateZipFile(zipFileTargetPath string, newFiles []*multipart.FileHeader, projectCode string) error {
// 	// Open the zip file for reading
// 	zipReader, err := zip.OpenReader(zipFileTargetPath)
// 	if err != nil {
// 		slog.Error("error open a zip file to read", "error", err)
// 		return err
// 	}
// 	defer zipReader.Close()

// 	// Open the zip file for writing
// 	zipFile, err := os.OpenFile(zipFileTargetPath, os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		slog.Error("error open a zip file to write", "error", err)
// 		return err
// 	}
// 	defer zipFile.Close()

// 	// Create a new zip writer
// 	zipWriter := zip.NewWriter(zipFile)
// 	defer zipWriter.Close()

// 	// Copy all files in existing zip
// 	for _, file := range zipReader.File {
// 		writer, err := zipWriter.Create(file.Name)
// 		if err != nil {
// 			slog.Error("error creating a zipWriter", "error", err, "fileName", file.Name)
// 			return err
// 		}
// 		reader, err := file.Open()
// 		if err != nil {
// 			slog.Error("error open a file reader", "error", err, "fileName", file.Name)
// 			return err
// 		}
// 		_, err = io.Copy(writer, reader)
// 		if err != nil {
// 			slog.Error("error at io.Copy", "error", err, "fileName", file.Name)
// 			return err
// 		}
// 		reader.Close()
// 	}
// 	// Write new files to the zip
// 	for _, newFileHeader := range newFiles {
// 		// newFile, err := OpenFileFromFileHeader(newFileHeader)
// 		newFile, err := newFileHeader.Open()
// 		if err != nil {
// 			return err
// 		}
// 		defer newFile.Close()

// 		header := &zip.FileHeader{
// 			UncompressedSize64: uint64(newFileHeader.Size),
// 		}

// 		// Set the file name in the zip header
// 		header.Name = fmt.Sprintf("%s_เอกสารอื่นๆ/%s", projectCode, newFileHeader.Filename) // Todo
// 		header.Method = zip.Deflate
// 		header.Modified = time.Now()

// 		writer, err := zipWriter.CreateHeader(header)
// 		if err != nil {
// 			panic(err)
// 		}
// 		_, err = io.Copy(writer, newFile)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	return nil
// }

// func (client *S3Service) downloadZipToLocal(zipFileTargetPath string, userId int, projectCode string) (*os.File, error) {
// 	downloader := manager.NewDownloader(client.S3Client)

// 	zipFile, err := os.Create(zipFileTargetPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	basePrefix := getBasePrefix(userId, projectCode)
// 	zipObjectKey := fmt.Sprintf("%s/zip/%s_เอกสารแนบ.zip", basePrefix, projectCode)

// 	_, err = downloader.Download(context.TODO(), zipFile, &s3.GetObjectInput{
// 		Bucket: aws.String(os.Getenv("AWS_S3_STORE_BUCKET_NAME")),
// 		Key:    aws.String(zipObjectKey),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return zipFile, nil
// }

// func getBasePrefix(userId int, projectCode string) string {
// 	return fmt.Sprintf("applicant/user_%d/%s", userId, projectCode)
// }
