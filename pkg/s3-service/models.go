package s3Service

type GetPresignedPayload struct {
	Path string `json:"path,omitempty"`
}
