package domain

import (
	"fmt"
	"strings"
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

func (s UploadSpec) Validate(l Limits) error {
	if _, err := SanitazeName(s.Name, l.MaxNameLen); err != nil {
		return err
	}

	if s.Size <= 0 {
		return fmt.Errorf("%v: %w", InvalidFieldError{Field: "Size", Reason: "too_small"}, ErrInvalid)
	}

	if s.Size > l.MaxUploadBytes {
		return fmt.Errorf("%v: %w", InvalidFieldError{Field: "Size", Reason: "exceeds_limit"}, ErrToLarge)
	}

	if !ValidContentType(s.ContentType) {
		return fmt.Errorf("%v: %w", InvalidFieldError{Field: "ContentType", Reason: "bad_format"}, ErrInvalid)
	}

	if err := ValidateOwnerID(string(s.OwnerID)); err != nil {
		return err
	}

	if len(s.Tags) > l.MaxTagsPerFile {
		return fmt.Errorf("%v: %w", InvalidFieldError{Field: "Tags", Reason: "exceeds_limit"}, ErrInvalid)
	}
	for _, tag := range s.Tags {
		if err := ValidateTag(string(tag)); err != nil {
			return nil
		}
	}

	k := strings.TrimSpace(s.IdempotencyKey)
	if k == "" || len(k) > 64 {
		return fmt.Errorf("%v: %w", InvalidFieldError{Field: "IdempotencyKey", Reason: "bad_format"}, ErrInvalid)
	}
	for i := 0; i < len(k); i++ {
		if k[i] < 0x21 || k[i] > 0x7e {
			return fmt.Errorf("%v: %w", InvalidFieldError{Field: "IdempotencyKey", Reason: "bad_format"}, ErrInvalid)
		}
	}
	return nil
}

func NewFile(spec UploadSpec, id FileID, checkSum Checksum, now time.Time, storageRef string, l Limits) (File, error) {
	if err := spec.Validate(l); err != nil {
		return File{}, err
	}

	if _, err := NormalizeFileID(string(id)); err != nil {
		return File{}, err
	}

	if err := ValidateChecksum(string(checkSum)); err != nil {
		return  File{}, err
	}

	if storageRef == "" || strings.HasPrefix(storageRef, "/") || strings.Contains(storageRef, "..") || strings.ContainsAny(storageRef, " \t\r\n") {
		return File{}, fmt.Errorf("%v: %w", InvalidFieldError{Field: "StorageRef", Reason: "bad_format"}, ErrInvalid)
	}

	name, err := SanitazeName(spec.Name, l.MaxNameLen)
	if err != nil {
		return File{}, err
	}

	if now.IsZero() {
		return File{}, fmt.Errorf("%v: %w", InvalidFieldError{Field: "CreatedAt", Reason: "empty"}, ErrInvalid)
	}
	now.UTC()

	tags := make([]Tag, len(spec.Tags))
	copy(tags, spec.Tags)

	return File{
		ID: id,
		Name: name,
		Size: spec.Size,
		ContentType: spec.ContentType,
		Checksum: checkSum,
		CreatedAt: now,
		OwnerID: spec.OwnerID,
		Tags: tags,
		StorageRef: storageRef,
	}, nil
}