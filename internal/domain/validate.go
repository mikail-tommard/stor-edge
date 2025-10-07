package domain

import (
	"strings"
	"unicode/utf8"

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

func SanitazeName(name string, maxLen int) (string, error) {
	for i := 0; i < len(name); i++ {
		switch name[i] {
		case 0x00, '\r', '\n', '\t':
			return "", &InvalidFieldError{Field: "name", Reason: "invalid_char"}
		}
	}

	for i := 0; i < len(name); i++ {
		if name[i] == '/' || name[i] == '\\' {
			return "", &InvalidFieldError{Field: "name", Reason: "invalid_char"}
		}
	}

	if strings.TrimSpace(name) == "" {
		return "", &InvalidFieldError{Field: "name", Reason: "empty"}
	}

	if utf8.RuneCountInString(name) > maxLen {
		return "", &InvalidFieldError{Field: "name", Reason: "too_long"}
	}

	return name, nil
}

func ValidContentType(ct string) bool {
	for i := 0; i < len(ct); i++ {
		switch ct[i]{
		case ' ', '\t', '\n', '\r':
			return false
		}
	}

	for i := 0; i < len(ct); i++ {
		if ct[i] == ';' {
			return false
		}
	}

	slash := -1
	for i := 0; i < len(ct); i++ {
		if ct[i] == '/' {
			if slash != -1 {
				return false
			}
			slash = i
		}
	}
	if slash <= 0 || slash == len(ct) - 1 {
		return false
	}

	okChar := func(b byte) bool {
		switch {
		case b >= 'A' && b <= 'Z':
			return true
		case b >= 'a' && b <= 'z':
			return true
		case b >= '0' && b <= '9':
			return true
		}
		switch b {
		case '!', '#', '$', '&', '^', '_', '.', '+', '-':
			return true
		}
		return false
	}

	for i := 0; i < slash; i++ {
		if !okChar(ct[i]) {
			return false
		}
	}

	for i := slash + 1; i < len(ct); i++ {
		if !okChar(ct[i]) {
			return false
		}
	}
	return true
}