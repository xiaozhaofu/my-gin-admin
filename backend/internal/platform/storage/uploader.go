package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type Uploader interface {
	Save(file *multipart.FileHeader, scene string) (*UploadedFile, error)
}

type UploadedFile struct {
	Random     string `json:"random"`
	OriginName string `json:"origin_name"`
	Path       string `json:"path"`
	MD5        string `json:"md5"`
	Type       int8   `json:"type"`
	Scene      string `json:"scene"`
	Provider   string `json:"provider"`
}

type preparedUpload struct {
	random     string
	originName string
	ext        string
	md5        string
	fileType   int8
	tmpPath    string
}

func prepareUpload(file *multipart.FileHeader, allowedExts []string, maxSize int64) (*preparedUpload, error) {
	if file.Size > maxSize {
		return nil, fmt.Errorf("file exceeds %d bytes", maxSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !slices.Contains(allowedExts, ext) {
		return nil, fmt.Errorf("file extension %s not allowed", ext)
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	tmp, err := os.CreateTemp("", "sleep-admin-upload-*"+ext)
	if err != nil {
		return nil, err
	}
	defer tmp.Close()

	hash := md5.New()
	if _, err := io.Copy(io.MultiWriter(tmp, hash), src); err != nil {
		_ = os.Remove(tmp.Name())
		return nil, err
	}

	return &preparedUpload{
		random:     uuid.NewString(),
		originName: file.Filename,
		ext:        ext,
		md5:        hex.EncodeToString(hash.Sum(nil)),
		fileType:   detectType(ext),
		tmpPath:    tmp.Name(),
	}, nil
}

func detectType(ext string) int8 {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return 1
	case ".mp4", ".mov", ".avi", ".mkv":
		return 2
	case ".mp3", ".wav", ".aac", ".m4a":
		return 4
	default:
		return 1
	}
}

func allowedTypeByScene(scene string) (map[int8]struct{}, string) {
	switch strings.TrimSpace(scene) {
	case "cover":
		return map[int8]struct{}{1: {}}, "cover"
	case "article-video":
		return map[int8]struct{}{2: {}}, "article/video"
	case "article-audio":
		return map[int8]struct{}{4: {}}, "article/audio"
	case "article-image":
		return map[int8]struct{}{1: {}}, "article/image"
	default:
		return map[int8]struct{}{1: {}, 2: {}, 4: {}}, "misc"
	}
}
