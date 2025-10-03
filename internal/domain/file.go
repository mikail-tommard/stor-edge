package domain

import (
	"time"
)

type File struct {
	ID FileID
	Name string
	Size int64
	ContentType string
	Checksum Checksum
	CreatedAt time.Time
	OwnerID OwnerID
	Tags []Tag
	StorageRef string
}

type UploadSpec struct {
	Name string
	Size int64
	ContentType string
	IdempotencyKey string
	OwnerID OwnerID
	Tags []Tag
}
