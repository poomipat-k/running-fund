package s3Service

import v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"

type GetPresignedPayload struct {
	Path                   string `json:"path,omitempty"`
	ProjectCreatedByUserId int    `json:"projectCreatedByUserId,omitempty"`
}

type PutPresignedToStaticBucketRequest struct {
	ObjectKey string `json:"objectKey,omitempty"`
}

type PutPresignedToStaticBucketResponse struct {
	Presigned *v4.PresignedHTTPRequest `json:"presigned,omitempty"`
	FullPath  string                   `json:"fullPath,omitempty"`
}
