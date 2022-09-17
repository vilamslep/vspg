package yandex

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vilamslep/vspg/cloud"
)

type Client struct {
	s3client *s3.Client 
	bucketName string
	osSep string
	cloudSep string
	cloudRoot string
}

var ErrLoadingConfiguration = fmt.Errorf("Failed to load cloud configuration") 

func NewClient(root string) (*Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(yandexResolver)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, ErrLoadingConfiguration
	}

	s3client := s3.NewFromConfig(cfg)

	osSep := "/"
	if runtime.GOOS == "windows" {
		osSep = "\\"
	}

	return &Client{
		s3client: s3client, 
		cloudRoot: root, 
		cloudSep: "/",	
		osSep: osSep,
	}, nil
}

func (c Client) Add(src string, bucket string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	if stat, err := f.Stat(); err == nil {
		
		isDir := stat.IsDir()
		f.Close() 
		
		c.bucketName = bucket
		if isDir  {
			return c.uploadDir(src)
		} else {
			return c.uploadFile(src, 
				filepath.Base(filepath.Dir(src)))
		}
	} else {
		return err
	}
}

func (c Client) KeepNecessaryQuantity(bucket string, keepCount int) error {
	
	c.bucketName = bucket
	folders, err := c.getFoldesSlice()
	if err != nil {
		return err
	}

	sort.Slice(folders, func(i int, j int) bool {
		return folders[i].LastModified.Before(folders[j].LastModified)
	})

	if len(folders) > keepCount {
		var ls *s3.ListObjectsV2Output
		var err error
	
		delete := folders[:len(folders)-keepCount]

		for _, f := range delete {
			params := &s3.ListObjectsV2Input{ 
				Bucket: aws.String(c.bucketName),
				Prefix: aws.String(fmt.Sprintf("%s/%s", c.cloudRoot, f.Name)),
			}

			if ls, err = c.s3client.ListObjectsV2(context.TODO(), params); err != nil {
				return nil
			}

			for _, obj := range ls.Contents {
				if err := c.Delete(*obj.Key); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c Client) getFoldesSlice() ([]cloud.S3Folder, error) {
	var ls *s3.ListObjectsV2Output
	var err error

	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(c.cloudRoot),
	}

	if ls, err = c.s3client.ListObjectsV2(context.TODO(), params); err != nil {
		return nil, err
	}

	files := make(map[string]time.Time)
	for _, object := range ls.Contents {
		files[*object.Key] = *object.LastModified
	}

	dirs := make(map[string]time.Time, 0)
	for k, v := range files {
		path := strings.Split(k, "/")
		dirs[path[1]] = v
	}

	folders := make([]cloud.S3Folder, 0, len(dirs))
	for k, v := range dirs {
		folders = append(folders, cloud.S3Folder{Name: k, LastModified: v})
	}

	return folders, nil
}

func (c Client) uploadDir(src string) error {
	dirRoot := filepath.Base(src)
	c.cloudRoot += c.cloudSep + dirRoot
	return filepath.Walk(src,
		func (path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			} 

			if info.IsDir() {
				return nil
			} 
 			return c.uploadFile(path, dirRoot)
		})
}

func (c Client) uploadFile(path string, root string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	slpath := strings.Split(path, c.osSep)
	start := 0
	for i := range slpath {
		if slpath[i] == root {
			start = (i + 1)
			break
		}
	}
	yapath := fmt.Sprintf("%s%s%s", c.cloudRoot, c.cloudSep, strings.Join(slpath[start:], c.cloudRoot)) 
	
	object := &s3.PutObjectInput{
		Bucket:        aws.String(c.bucketName),
		Key:           aws.String(yapath),
		Body:          file,
		ContentLength: info.Size(),
	}

	if _, err = c.s3client.PutObject(context.TODO(), object); err != nil {
		return err
	} else {
		return nil
	}
}

func (c Client) Delete(path string) error {
	deleteParams := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(path),
	}

	if _, err := c.s3client.DeleteObject(context.TODO(), deleteParams); err != nil {
		return err
	}
	return nil
}

func yandexResolver(service string, region string, options ...interface{}) (aws.Endpoint, error) {
	if service == s3.ServiceID && region == "ru-central1" {
		return aws.Endpoint{
			PartitionID:   "yc",
			URL:           "https://storage.yandexcloud.net",
			SigningRegion: "ru-central1",
		}, nil
	}
	return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
}