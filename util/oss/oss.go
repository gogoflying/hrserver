package oss

import (
	"bytes"
	"io/ioutil"

	ali "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var OssClient *ali.Client

func DitalOSS(endpoint, accessId, accessKey string) {
	var (
		err error
	)
	if OssClient, err = ali.New(endpoint, accessId, accessKey); err != nil {
		panic(err)
	}
}

func WriteFile(bucketName, fileName string, data []byte) error {
	bucket, err := OssClient.Bucket(bucketName)
	if err != nil {
		return err
	}
	return bucket.PutObject(fileName, bytes.NewReader(data))
}

func ReadFile(bucketName, fileName string) ([]byte, error) {
	bucket, err := OssClient.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	reader, err := bucket.GetObject(fileName)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}

func DeleteFile(bucketName, fileName string) error {
	bucket, err := OssClient.Bucket(bucketName)
	if err != nil {
		return err
	}
	return bucket.DeleteObject(fileName)
}
