package storage

import (
	"io"

	"backup/pkg/config"
	"backup/pkg/storage/core"
	"backup/pkg/storage/disk"
	"backup/pkg/storage/minio"
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

// S3 for storage interface
var S3 Storage

// NewEngine return storage interface
func NewEngine(config config.Config) (err error) {
	switch config.Storage.Driver {
	case "s3":
		S3, err = minio.NewEngine(
			config.Storage.Endpoint,
			config.Storage.AccessID,
			config.Storage.SecretKey,
			config.Storage.SSL,
			config.Storage.Region,
		)
		if err != nil {
			return err
		}
	case "disk":
		S3 = disk.NewEngine(
			config.Server.Addr,
			config.Storage.Path,
		)
	}

	return nil
}

// NewS3Engine return storage interface
func NewS3Engine(endPoint, accessID, secretKey string, ssl bool, region string) (Storage, error) {
	return minio.NewEngine(
		endPoint,
		accessID,
		secretKey,
		ssl,
		region,
	)
}

// NewDiskEngine return storage interface
func NewDiskEngine(host, folder string) (Storage, error) {
	return disk.NewEngine(
		host,
		folder,
	), nil
}
