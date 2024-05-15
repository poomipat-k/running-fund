package s3Service

type GetPresignedPayload struct {
	Path                   string `json:"path,omitempty"`
	ProjectCreatedByUserId int    `json:"projectCreatedByUserId,omitempty"`
}

type PutPresignedToStaticBucketRequest struct {
	ObjectKey string `json:"objectKey,omitempty"`
}
