package entity

import "time"

type ManifestSegment struct {
	S3Storage
	ModifiedTime time.Time
}
