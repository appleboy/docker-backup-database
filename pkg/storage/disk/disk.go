package disk

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"backup/pkg/storage/core"
)

// BUFFERSIZE copy file buffer size
var BUFFERSIZE = 1000

func copy(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// Disk client
type Disk struct {
	Host string
	Path string
}

// NewEngine struct
func NewEngine(host, path string) (*Disk, error) {
	return &Disk{
		Host: host,
		Path: path,
	}, nil
}

// UploadFile to upload file to disk
func (d *Disk) UploadFile(bucketName, fileName string, content []byte, _ io.Reader) error {
	// check folder exists
	// ex: bucket + foo/bar/uuid.tar.gz
	storage := path.Join(d.Path, bucketName, filepath.Dir(fileName))
	if err := os.MkdirAll(storage, os.ModePerm); err != nil {
		return nil
	}
	return ioutil.WriteFile(d.FilePath(bucketName, fileName), content, os.FileMode(0644))
}

// BucketExists verify if bucket exists and you have permission to access it.
func (d *Disk) BucketExists(bucketName string) (bool, error) {
	_, err := os.Stat(bucketName)
	return !os.IsNotExist(err), err
}

// CreateBucket create bucket
func (d *Disk) CreateBucket(bucketName, region string) error {
	storage := path.Join(d.Path, bucketName)
	if err := os.MkdirAll(storage, os.ModePerm); err != nil {
		return nil
	}

	return nil
}

// FilePath for store path + file name
func (d *Disk) FilePath(bucketName, fileName string) string {
	return path.Join(
		d.Path,
		bucketName,
		fileName,
	)
}

// DeleteFile delete file
func (d *Disk) DeleteFile(bucketName, fileName string) error {
	return os.Remove(d.FilePath(bucketName, fileName))
}

// GetFileURL for storage host + bucket + filename
func (d *Disk) GetFileURL(bucketName, fileName string) string {
	if d.Host != "" {
		if u, err := url.Parse(d.Host); err == nil {
			u.Path = path.Join(u.Path, d.Path, bucketName, fileName)
			return u.String()
		}
	}
	return path.Join(d.Path, bucketName, fileName)
}

// DownloadFile downloads and saves the object as a file in the local filesystem.
func (d *Disk) DownloadFile(bucketName, fileName, target string) error {
	return nil
}

// GetContent for storage bucket + filename
func (d *Disk) GetContent(bucketName, fileName string) ([]byte, error) {
	return ioutil.ReadFile(path.Join(d.Path, bucketName, fileName))
}

// CopyFile copy src to dest
func (d *Disk) CopyFile(srcBucketName, srcFile, destBucketName, destFile string) error {
	src := path.Join(d.Path, srcBucketName, srcFile)
	dest := path.Join(d.Path, destBucketName, destFile)
	return copy(src, dest, int64(BUFFERSIZE))
}

// FileExist check object exist. bucket + filename
func (d *Disk) FileExist(bucketName, fileName string) bool {
	src := path.Join(d.Path, bucketName, fileName)
	_, err := os.Stat(src)

	return err == nil
}

// Client get disk client
func (d *Disk) Client() interface{} {
	return nil
}

// SignedURL support signed URL
func (d *Disk) SignedURL(bucketName, filename string, opts *core.SignedURLOptions) (string, error) {
	return d.GetFileURL(bucketName, filename), nil
}
