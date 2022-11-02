package cloud

import "time"

type Client interface {
	Add(src string, bucket string) error
	Delete(path string) error
}

type S3Folder struct {
	Name         string
	LastModified time.Time
}
