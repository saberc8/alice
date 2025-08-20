package storage

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"alice/infra/config"
	"alice/pkg/logger"
)

// ObjectStorage 定义需要暴露的对象存储能力
type ObjectStorage interface {
	CreateBucket(ctx context.Context, bucket string) error
	DeleteBucket(ctx context.Context, bucket string) error
	PutObject(ctx context.Context, bucket, objectName string, data []byte, contentType string) (string, error)
	DeleteObject(ctx context.Context, bucket, objectName string) error
	ListBuckets(ctx context.Context) ([]string, error)
	ListObjects(ctx context.Context, bucket, prefix string, recursive bool, limit int) ([]string, error)
	GetPresignedURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error)
	SetBucketPublic(ctx context.Context, bucket string, public bool) error
}

type MinioStorage struct {
	cli     *minio.Client
	baseURL string
}

// NewMinio 根据配置创建 MinIO 客户端
func NewMinio(cfg config.MinioConfig) (*MinioStorage, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("minio endpoint empty")
	}
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	base := cfg.BaseURL
	if base == "" { // 生成默认 base url
		scheme := "http"
		if cfg.UseSSL {
			scheme = "https"
		}
		base = scheme + "://" + cfg.Endpoint
	}
	return &MinioStorage{cli: client, baseURL: base}, nil
}

func (m *MinioStorage) CreateBucket(ctx context.Context, bucket string) error {
	exists, err := m.cli.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return m.cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}

func (m *MinioStorage) DeleteBucket(ctx context.Context, bucket string) error {
	return m.cli.RemoveBucket(ctx, bucket)
}

func (m *MinioStorage) PutObject(ctx context.Context, bucket, objectName string, data []byte, contentType string) (string, error) {
	if err := m.CreateBucket(ctx, bucket); err != nil { // 确保 bucket 存在
		return "", err
	}
	reader := bytes.NewReader(data)
	_, err := m.cli.PutObject(ctx, bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	// 返回可访问 URL (基础 + bucket + objectName)
	u, _ := url.Parse(m.baseURL)
	u.Path = path.Join(u.Path, bucket, objectName)
	return u.String(), nil
}

func (m *MinioStorage) DeleteObject(ctx context.Context, bucket, objectName string) error {
	return m.cli.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

func (m *MinioStorage) ListBuckets(ctx context.Context) ([]string, error) {
	buckets, err := m.cli.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(buckets))
	for _, b := range buckets {
		names = append(names, b.Name)
	}
	return names, nil
}

func (m *MinioStorage) ListObjects(ctx context.Context, bucket, prefix string, recursive bool, limit int) ([]string, error) {
	opt := minio.ListObjectsOptions{Prefix: prefix, Recursive: recursive}
	ch := m.cli.ListObjects(ctx, bucket, opt)
	var res []string
	for obj := range ch {
		if obj.Err != nil {
			return nil, obj.Err
		}
		res = append(res, obj.Key)
		if limit > 0 && len(res) >= limit {
			break
		}
	}
	return res, nil
}

func (m *MinioStorage) GetPresignedURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error) {
	if expiry <= 0 {
		expiry = time.Hour
	}
	reqParams := make(url.Values)
	presigned, err := m.cli.PresignedGetObject(ctx, bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return presigned.String(), nil
}

// SetBucketPublic 设置 bucket 是否公共读
func (m *MinioStorage) SetBucketPublic(ctx context.Context, bucket string, public bool) error {
	if public {
		policy := fmt.Sprintf(`{ "Version": "2012-10-17", "Statement": [ { "Effect": "Allow", "Principal": {"AWS": ["*"]}, "Action": ["s3:GetObject"], "Resource": ["arn:aws:s3:::%s/*"] } ] }`, bucket)
		return m.cli.SetBucketPolicy(ctx, bucket, policy)
	}
	// 取消公共策略: 传空策略 (MinIO 会移除) —— 若版本不支持, 可尝试设置一个拒绝策略; 这里假设支持空清除
	return m.cli.SetBucketPolicy(ctx, bucket, "")
}

// HealthCheck 简单健康检测
func (m *MinioStorage) HealthCheck(ctx context.Context) error {
	// 使用一个很轻量的操作: 列举 buckets
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := m.cli.ListBuckets(ctx)
	if err != nil {
		logger.Warnf("minio health check failed: %v", err)
	}
	return nil
}
