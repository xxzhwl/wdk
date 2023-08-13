// Package storage minio存储
// Author: wanlizhan
// Date: 2023/8/5
package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/xxzhwl/wdk/dict"
	"github.com/xxzhwl/wdk/uconfig"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient Minio存储客户端
type MinioClient struct {
	*minio.Client
}

// NewMinioClientBySchema 配置获取Minio客户端
func NewMinioClientBySchema(schema string) (MinioClient, error) {
	m, err := uconfig.ME(schema)
	if err != nil {
		return MinioClient{}, err
	}
	minioClient, err := minio.New(dict.S(m, "EndPoint"), &minio.Options{
		Creds:  credentials.NewStaticV4(dict.S(m, "AccessKey"), dict.S(m, "SecretKey"), ""),
		Secure: dict.B(m, "UseSSL"),
	})
	if err != nil {
		return MinioClient{}, err
	}
	return MinioClient{minioClient}, nil
}

// NewMinioClient 默认配置获取Minio客户端
func NewMinioClient() (MinioClient, error) {
	return NewMinioClientBySchema("Minio.Default")
}

// MakeBucket 创建一个Bucket
func (m MinioClient) MakeBucket(bucketName string) error {
	return m.Client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
}

// BucketExists 判断一个Bucket是否存在
func (m MinioClient) BucketExists(bucketName string) (bool, error) {
	return m.Client.BucketExists(context.Background(), bucketName)
}

// UploadFile 上传文件到指定bucket
func (m MinioClient) UploadFile(bucketName, fileName string, content []byte) error {
	reader := bytes.NewReader(content)
	info, err := m.Client.PutObject(context.Background(), bucketName, fileName, reader, int64(reader.Len()), minio.PutObjectOptions{})
	fmt.Println(info.Location)
	return err
}
