package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"

	v2config "go_sleep_admin/internal/platform/config"
)

type OSSUploader struct {
	client       *oss.Client
	bucket       string
	endpoint     string
	objectPrefix string
	publicBase   string
	allowedExts  []string
	maxSize      int64
	disableSSL   bool
	useCName     bool
}

func NewOSSUploader(cfg v2config.UploadOSSConfig, allowedExts []string, maxSize int64) *OSSUploader {
	endpoint := normalizeEndpoint(cfg.Endpoint)
	ossCfg := oss.LoadDefaultConfig().
		WithRegion(cfg.Region).
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret, cfg.SessionToken))

	if endpoint != "" {
		ossCfg = ossCfg.WithEndpoint(endpoint)
	}
	if cfg.DisableSSL {
		ossCfg = ossCfg.WithDisableSSL(true)
	}
	if cfg.UseInternal {
		ossCfg = ossCfg.WithUseInternalEndpoint(true)
	}
	if cfg.UseCName {
		ossCfg = ossCfg.WithUseCName(true)
	}

	return &OSSUploader{
		client:       oss.NewClient(ossCfg),
		bucket:       cfg.Bucket,
		endpoint:     endpoint,
		objectPrefix: strings.Trim(strings.TrimSpace(cfg.ObjectPrefix), "/"),
		publicBase:   strings.TrimRight(strings.TrimSpace(cfg.PublicBaseURL), "/"),
		allowedExts:  allowedExts,
		maxSize:      maxSize,
		disableSSL:   cfg.DisableSSL,
		useCName:     cfg.UseCName,
	}
}

func (u *OSSUploader) Save(file *multipart.FileHeader, scene string) (*UploadedFile, error) {
	prepared, err := prepareUpload(file, u.allowedExts, u.maxSize)
	if err != nil {
		return nil, err
	}
	defer os.Remove(prepared.tmpPath)

	allowedTypes, subdir := allowedTypeByScene(scene)
	if _, ok := allowedTypes[prepared.fileType]; !ok {
		return nil, fmt.Errorf("file type %d is not allowed for scene %s", prepared.fileType, scene)
	}

	objectKey := u.buildObjectKey(subdir, prepared.random, prepared.ext)
	body, err := os.Open(prepared.tmpPath)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	_, err = u.client.PutObject(context.Background(), &oss.PutObjectRequest{
		Bucket: oss.Ptr(u.bucket),
		Key:    oss.Ptr(objectKey),
		Body:   body,
	})
	if err != nil {
		return nil, fmt.Errorf("put object to oss: %w", err)
	}

	return &UploadedFile{
		Random:     prepared.random,
		OriginName: prepared.originName,
		Path:       u.publicURL(objectKey),
		MD5:        prepared.md5,
		Type:       prepared.fileType,
		Scene:      scene,
		Provider:   "aliyun-oss",
	}, nil
}

func (u *OSSUploader) buildObjectKey(subdir, random, ext string) string {
	parts := []string{}
	if u.objectPrefix != "" {
		parts = append(parts, u.objectPrefix)
	}
	if subdir != "" {
		parts = append(parts, subdir)
	}
	parts = append(parts, time.Now().Format("2006/01/02"), random+ext)
	return path.Join(parts...)
}

func (u *OSSUploader) publicURL(objectKey string) string {
	if u.publicBase != "" {
		return u.publicBase + "/" + objectKey
	}
	scheme := "https"
	if u.disableSSL {
		scheme = "http"
	}
	if u.useCName {
		return scheme + "://" + strings.TrimRight(u.endpoint, "/") + "/" + objectKey
	}
	return scheme + "://" + u.bucket + "." + strings.TrimRight(u.endpoint, "/") + "/" + objectKey
}

func normalizeEndpoint(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		if parsed, err := url.Parse(raw); err == nil {
			return parsed.Host
		}
	}
	return strings.TrimRight(raw, "/")
}
