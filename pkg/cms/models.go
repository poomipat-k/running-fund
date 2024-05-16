package cms

type S3UploadResponse struct {
	ObjectKey string `json:"objectKey,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`
}
