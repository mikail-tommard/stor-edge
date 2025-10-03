package domain

const (
	DefaultMaxUploadBytes = 10 * 1024 * 1024
	DefaultMaxTagsPerFile = 16
	DefaultMaxNameLen = 255
)

const ulidLen = 26

const tagMaxLen = 32

type FileID string

type OwnerID string

type Checksum string

type Tag string

type Limits struct {
	MaxUploadBytes int64
	MaxTagsPerFile int
	MaxNameLen int
}

func DefaultLimits() Limits {
	return Limits{
		MaxUploadBytes: DefaultMaxUploadBytes,
		MaxTagsPerFile: DefaultMaxTagsPerFile,
		MaxNameLen: DefaultMaxNameLen,
	}
}

// func ValidateFieldID(s string) error
// func NormalizeFileID(s string) (FileId, error)
// func ValidateOwnerID(s string) error
// func ValidateChecksum(s string) error
// func ValidateTag(s string) error
// func SanitazeName(name string, maxLen int) (string, error)
// func ValidContentType(ct string) bool