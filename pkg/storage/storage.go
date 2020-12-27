package storage

import (
	"errors"
	"io"

	"github.com/appleboy/docker-backup-database/pkg/config"
	"github.com/appleboy/docker-backup-database/pkg/storage/core"
	"github.com/appleboy/docker-backup-database/pkg/storage/disk"
	"github.com/appleboy/docker-backup-database/pkg/storage/minio"
)

// Storage for s3 and disk
type Storage interface {
	// CreateBucket for create new folder
	CreateBucket(string, string) error
	// UploadFile for upload single file
	UploadFile(string, string, []byte, io.Reader) error
	// DeleteFile for delete single file
	DeleteFile(string, string) error
	// FilePath for store path + file name
	FilePath(string, string) string
	// GetFile for storage host + bucket + filename
	GetFileURL(string, string) string
	// DownloadFile downloads and saves the object as a file in the local filesystem.
	DownloadFile(string, string, string) error
	// BucketExist check object exist. bucket + filename
	BucketExists(string) (bool, error)
	// FileExist check object exist. bucket + filename
	FileExist(string, string) bool
	// GetContent for storage bucket + filename
	GetContent(string, string) ([]byte, error)
	// Copy Create or replace an object through server-side copying of an existing object.
	CopyFile(string, string, string, string) error
	// Client get storage client
	Client() interface{}
	// SignedURL get signed URL
	SignedURL(string, string, *core.SignedURLOptions) (string, error)
}

// NewEngine return storage interface
func NewEngine(cfg config.Config) (storage Storage, err error) {
	switch cfg.Storage.Driver {
	case "s3":
		return minio.NewEngine(
			cfg.Storage.Endpoint,
			cfg.Storage.AccessID,
			cfg.Storage.SecretKey,
			cfg.Storage.SSL,
			cfg.Storage.InsecureSkipVerify,
			cfg.Storage.Region,
		)
	case "disk":
		return disk.NewEngine(
			cfg.Server.Addr,
			cfg.Storage.Path,
		)
	}

	return nil, errors.New("We don't support Storage Dirver: " + cfg.Storage.Driver)
}
