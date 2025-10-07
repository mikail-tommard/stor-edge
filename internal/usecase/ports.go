package usecase

import (
	"context"
	"io"
	"stor-edge/internal/domain"
	"time"
)

type IDGet interface {
	New() (domain.FileID, error)
}

type Clock interface {
	Now() time.Time
}

type FileStore interface {
	WriteTemp(ctx context.Context, id domain.FileID) (w io.WriteCloser, tmpRef string, err error)
	Commit(ctx context.Context, tmpRef string) (storageRef string, err error)
	Open(ctx context.Context, storageRef string, rng *[2]int64)(r io.ReadCloser,size int64, contentType string, err error)
	Stat(ctx context.Context, storageRef string) (size int64, err error)
	Delete(ctx context.Context, storageRef string) error
}

type ListQuery struct {
	OwnerID *domain.OwnerID
	Tag *domain.Tag
	Q *string
	Limit int
	Cursor *string
}

type MetaRepo interface {
	Create(ctx context.Context, f domain.File) error
	Get(ctx context.Context, id domain.FileID) (domain.File, error)
	Delete(cxt context.Context, id domain.FileID) error
	List(ctx context.Context, q ListQuery) ([]domain.File, error)
	SaveIdompotency(ctx context.Context, key string, result domain.FileID) (already bool, existdID domain.FileID, err error)
}

type EventPublisher interface {
	FileUploaded(ctx context.Context, f domain.File) error
	FileDeleted(ctx context.Context, f domain.FileID) error
}

type TxManaget interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}