package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type LocalUploader struct {
	localDir    string
	publicPath  string
	allowedExts []string
	maxSize     int64
}

func NewLocalUploader(localDir, publicPath string, allowedExts []string, maxSize int64) *LocalUploader {
	return &LocalUploader{
		localDir:    localDir,
		publicPath:  publicPath,
		allowedExts: allowedExts,
		maxSize:     maxSize,
	}
}

func (u *LocalUploader) Save(file *multipart.FileHeader, scene string) (*UploadedFile, error) {
	prepared, err := prepareUpload(file, u.allowedExts, u.maxSize)
	if err != nil {
		return nil, err
	}
	defer os.Remove(prepared.tmpPath)

	allowedTypes, subdir := allowedTypeByScene(scene)
	if _, ok := allowedTypes[prepared.fileType]; !ok {
		return nil, fmt.Errorf("file type %d is not allowed for scene %s", prepared.fileType, scene)
	}

	if err := os.MkdirAll(u.localDir, 0o755); err != nil {
		return nil, err
	}

	name := prepared.random + prepared.ext
	targetDir := filepath.Join(u.localDir, filepath.FromSlash(subdir))
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return nil, err
	}
	targetPath := filepath.Join(targetDir, name)
	if err := os.Rename(prepared.tmpPath, targetPath); err == nil {
		return &UploadedFile{
			Random:     prepared.random,
			OriginName: prepared.originName,
			Path:       strings.TrimRight(u.publicPath, "/") + "/" + strings.TrimLeft(pathJoinSlash(subdir, name), "/"),
			MD5:        prepared.md5,
			Type:       prepared.fileType,
			Scene:      scene,
			Provider:   "local",
		}, nil
	}

	src, err := os.Open(prepared.tmpPath)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	dst, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	return &UploadedFile{
		Random:     prepared.random,
		OriginName: prepared.originName,
		Path:       strings.TrimRight(u.publicPath, "/") + "/" + strings.TrimLeft(pathJoinSlash(subdir, name), "/"),
		MD5:        prepared.md5,
		Type:       prepared.fileType,
		Scene:      scene,
		Provider:   "local",
	}, nil
}

func pathJoinSlash(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.Trim(part, "/")
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return strings.Join(filtered, "/")
}
