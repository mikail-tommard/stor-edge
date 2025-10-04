package domain

import (
	"strings"

	"github.com/oklog/ulid/v2"
)

func ValidateFileID(s string) error {
	if len(s) != ulidLen {
		return &InvalidFieldError{Field: "ID", Reason: "wrong length"}
	}
	if _, err := ulid.Parse(s); err != nil {
		return &InvalidFieldError{Field: "ID", Reason:"invalid ulid"}
	}
	return nil
}

func NormalizeFileID(s string) (FileID, error) {
	normalize := strings.ToLower(strings.TrimSpace(s))
	if err := ValidateFileID(normalize); err != nil {
		return "", err
	}
	return FileID(normalize), nil
}

func ValidateOwnerID(s string) error {
	if len(s) == 0 {
		return &InvalidFieldError{Field: "owner_id", Reason:"empty"}
	}
	if len(s) > 34 {
		return &InvalidFieldError{Field: "owner_id", Reason:"too long"}
	}

	if strings.TrimSpace(s) != s {
		return &InvalidFieldError{Field: "owner_id", Reason: "non_ascii"}
	}

	for i := 0; i < len(s); i++ {
		b := s[i]
		if b < 32 || b > 126 {
			return &InvalidFieldError{Field: "owner_id", Reason:"non_ascii"}
		}
	}

	return nil
}

func ValidateChecksum(s string) error {
	if len(s) != 64 {
		return &InvalidFieldError{Field: "checksum", Reason: "wrong length"}
	}

	for i := 0; i < len(s); i++ {
		c := s[i]

		if c >= '0' && c <= '9' {
			continue
		}
		if c >= 'a' && c <= 'f' {
			continue
		}

		return &InvalidFieldError{Field: "checksum", Reason: "not hex"}
	}

	return nil
}

func ValidateTag(s string) error {
	n := len(s)
	if n == 0 || n > 32 {
		return &InvalidFieldError{Field: "tag", Reason: "bad_format"}
	}

	c := s[0]
	if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
		return &InvalidFieldError{Field: "tag", Reason: "bad_format"}
	}

	for i := 0; i < n; i++ {
		c = s[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			continue
		} 
		return &InvalidFieldError{Field: "tag", Reason: "bad_format"}
	}
	return nil
}	