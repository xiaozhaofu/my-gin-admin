package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	cos "github.com/tencentyun/cos-go-sdk-v5"

	v2config "go_sleep_admin/internal/platform/config"
)

type COSUploader struct {
	client       *cos.Client
	objectPrefix string
	publicBase   string
	allowedExts  []string
	maxSize      int64
}

func NewCOSUploader(cfg v2config.UploadCOSConfig, allowedExts []string, maxSize int64) (*COSUploader, error) {
	u, err := normalizeCOSBucketURL(cfg)
	if err != nil {
		return nil, err
	}
	baseURL := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(baseURL, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})
	return &COSUploader{
		client:       client,
		objectPrefix: strings.Trim(strings.TrimSpace(cfg.ObjectPrefix), "/"),
		publicBase:   strings.TrimRight(strings.TrimSpace(cfg.PublicBaseURL), "/"),
		allowedExts:  allowedExts,
		maxSize:      maxSize,
	}, nil
}

func normalizeCOSBucketURL(cfg v2config.UploadCOSConfig) (*url.URL, error) {
	raw := strings.TrimSpace(cfg.BucketURL)
	if raw == "" {
		return nil, fmt.Errorf("cos bucket_url is required")
	}

	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		u, err := url.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("parse COS bucket url: %w", err)
		}
		if u.Host == "" {
			return nil, fmt.Errorf("invalid COS bucket url: %s", raw)
		}
		u.Path = ""
		u.RawPath = ""
		return u, nil
	}

	host := strings.Trim(raw, "/")
	if !strings.Contains(host, ".cos.") || !strings.Contains(host, ".myqcloud.com") {
		if strings.TrimSpace(cfg.Region) == "" {
			return nil, fmt.Errorf("cos region is required when bucket_url is not a full COS host")
		}
		if !isValidCOSBucketName(host) {
			return nil, fmt.Errorf("invalid COS bucket name %q, expected {name}-{appid}", host)
		}
		host = fmt.Sprintf("%s.cos.%s.myqcloud.com", host, strings.TrimSpace(cfg.Region))
	}

	u, err := url.Parse("https://" + host)
	if err != nil {
		return nil, fmt.Errorf("build COS bucket url: %w", err)
	}
	return u, nil
}

func isValidCOSBucketName(bucket string) bool {
	idx := strings.LastIndex(bucket, "-")
	if idx <= 0 || idx == len(bucket)-1 {
		return false
	}
	for _, r := range bucket[idx+1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func (u *COSUploader) Save(file *multipart.FileHeader, scene string) (*UploadedFile, error) {
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
	if _, err := u.client.Object.PutFromFile(context.Background(), objectKey, prepared.tmpPath, nil); err != nil {
		return nil, fmt.Errorf("put object to cos: %w", err)
	}

	return &UploadedFile{
		Random:     prepared.random,
		OriginName: prepared.originName,
		Path:       u.publicURL(objectKey),
		MD5:        prepared.md5,
		Type:       prepared.fileType,
		Scene:      scene,
		Provider:   "tencent-cos",
	}, nil
}

func (u *COSUploader) buildObjectKey(subdir, random, ext string) string {
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

func (u *COSUploader) publicURL(objectKey string) string {
	if u.publicBase != "" {
		return u.publicBase + "/" + objectKey
	}
	return strings.TrimRight(u.client.BaseURL.BucketURL.String(), "/") + "/" + objectKey
}
