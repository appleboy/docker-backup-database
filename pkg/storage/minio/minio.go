package minio

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"backup/pkg/storage/core"

	"github.com/h2non/filetype"
	"github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/credentials"
)

// Minio client
type Minio struct {
	client *minio.Client
}

// NewEngine struct
func NewEngine(endpoint, accessID, secretKey string, ssl bool, region string) (*Minio, error) {
	var client *minio.Client
	var err error
	if endpoint == "" {
		return nil, errors.New("endpoint can't be empty")
	}

	// Fetching from IAM roles assigned to an EC2 instance.
	if accessID == "" && secretKey == "" {
		iam := credentials.NewIAM("")
		client, err = minio.NewWithCredentials(endpoint, iam, ssl, region)
	} else {
		// Initialize minio client object.
		client, err = minio.NewWithRegion(endpoint, accessID, secretKey, ssl, region)
	}

	if err != nil {
		return nil, err
	}

	client.SetCustomTransport(&http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	})

	return &Minio{
		client: client,
	}, nil
}

// UploadFile to s3 server
func (m *Minio) UploadFile(bucketName, objectName string, content []byte, reader io.Reader) error {
	contentType := ""
	kind, _ := filetype.Match(content)
	if kind != filetype.Unknown {
		contentType = kind.MIME.Value
	}

	if contentType == "" {
		contentType = http.DetectContentType(content)
	}

	userMetaData := map[string]string{"x-amz-acl": "public-read"}
	cacheControl := "max-age=31536000"
	opts := minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: userMetaData,
		CacheControl: cacheControl,
	}
	if reader != nil {
		opts.Progress = reader
	}

	// Upload the zip file with FPutObject
	_, err := m.client.PutObject(
		bucketName,
		objectName,
		bytes.NewReader(content),
		int64(len(content)),
		opts,
	)

	return err
}

// BucketExists verify if bucket exists and you have permission to access it.
func (m *Minio) BucketExists(bucketName string) (bool, error) {
	return m.client.BucketExists(bucketName)
}

// CreateBucket create bucket
func (m *Minio) CreateBucket(bucketName, region string) error {
	exists, err := m.client.BucketExists(bucketName)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return m.client.MakeBucket(bucketName, region)
}

// FilePath for store path + file name
func (m *Minio) FilePath(_, fileName string) string {
	return fmt.Sprintf("%s/%s", os.TempDir(), fileName)
}

// DeleteFile delete file
func (m *Minio) DeleteFile(bucketName, fileName string) error {
	return m.client.RemoveObject(bucketName, fileName)
}

// GetFileURL for storage host + bucket + filename
func (m *Minio) GetFileURL(bucketName, fileName string) string {
	return m.client.EndpointURL().String() + "/" + bucketName + "/" + fileName
}

// DownloadFile downloads and saves the object as a file in the local filesystem.
func (m *Minio) DownloadFile(bucketName, fileName, target string) error {
	return m.client.FGetObject(bucketName, fileName, target, minio.GetObjectOptions{})
}

// GetContent for storage bucket + filename
func (m *Minio) GetContent(bucketName, fileName string) ([]byte, error) {
	object, err := m.client.GetObject(bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(object)

	return buf.Bytes(), nil
}

// CopyFile copy src to dest
func (m *Minio) CopyFile(srcBucket, srcPath, destBucket, destPath string) error {
	src := minio.NewSourceInfo(srcBucket, srcPath, nil)

	// Destination object
	dst, err := minio.NewDestinationInfo(destBucket, destPath, nil, nil)
	if err != nil {
		return err
	}

	// Copy object call
	return m.client.CopyObject(dst, src)
}

// FileExist check object exist. bucket + filename
func (m *Minio) FileExist(bucketName, fileName string) bool {
	_, err := m.client.StatObject(bucketName, fileName, minio.StatObjectOptions{})

	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "AccessDenied" {
			return false
		}
		if errResponse.Code == "NoSuchBucket" {
			return false
		}
		if errResponse.Code == "InvalidBucketName" {
			return false
		}
		if errResponse.Code == "NoSuchKey" {
			return false
		}
		return false
	}

	return true
}

// Client get disk client
func (m *Minio) Client() interface{} {
	return m.client
}

// SignedURL support signed URL
func (m *Minio) SignedURL(bucketName, filename string, opts *core.SignedURLOptions) (string, error) {
	// Check if file exists
	if _, err := m.client.StatObject(bucketName, filename, minio.StatObjectOptions{}); err != nil {
		return "", err
	}

	reqParams := make(url.Values)
	if opts != nil && opts.DefaultFilename != "" {
		reqParams.Set("response-content-disposition", `attachment; filename="`+opts.DefaultFilename+`"`)
	}

	url, err := m.client.PresignedGetObject(bucketName, filename, opts.Expiry, reqParams)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
