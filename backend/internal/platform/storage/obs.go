package storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"

	v2config "go_sleep_admin/internal/platform/config"
)

type OBSUploader struct {
	client       *obs.ObsClient
	endpoint     string
	bucket       string
	objectPrefix string
	publicBase   string
	allowedExts  []string
	maxSize      int64
}

func NewOBSUploader(cfg v2config.UploadOBSConfig, allowedExts []string, maxSize int64) (*OBSUploader, error) {
	client, err := obs.New(cfg.AccessKeyID, cfg.AccessKeySecret, cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("new obs client: %w", err)
	}
	return &OBSUploader{
		client:       client,
		endpoint:     strings.TrimRight(strings.TrimSpace(cfg.Endpoint), "/"),
		bucket:       cfg.Bucket,
		objectPrefix: strings.Trim(strings.TrimSpace(cfg.ObjectPrefix), "/"),
		publicBase:   strings.TrimRight(strings.TrimSpace(cfg.PublicBaseURL), "/"),
		allowedExts:  allowedExts,
		maxSize:      maxSize,
	}, nil
}

func (u *OBSUploader) Save(file *multipart.FileHeader, scene string) (*UploadedFile, error) {
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
	_, err = u.client.PutFile(&obs.PutFileInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket: u.bucket,
				Key:    objectKey,
			},
		},
		SourceFile: prepared.tmpPath,
	})
	if err != nil {
		return nil, fmt.Errorf("put object to obs: %w", err)
	}

	return &UploadedFile{
		Random:     prepared.random,
		OriginName: prepared.originName,
		Path:       u.publicURL(objectKey),
		MD5:        prepared.md5,
		Type:       prepared.fileType,
		Scene:      scene,
		Provider:   "huawei-obs",
	}, nil
}

func (u *OBSUploader) buildObjectKey(subdir, random, ext string) string {
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

func (u *OBSUploader) publicURL(objectKey string) string {
	if u.publicBase != "" {
		return u.publicBase + "/" + objectKey
	}
	return strings.TrimRight(u.endpoint, "/") + "/" + u.bucket + "/" + objectKey
}
