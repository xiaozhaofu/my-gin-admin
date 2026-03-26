package storage

import (
	"fmt"
	"mime/multipart"
	"strings"
)

type UploadGateway interface {
	Save(file *multipart.FileHeader, scene, provider string) (*UploadedFile, error)
}

type gateway struct {
	defaultProvider string
	uploaders       map[string]Uploader
}

func NewUploadGateway(defaultProvider string, uploaders map[string]Uploader) UploadGateway {
	return &gateway{
		defaultProvider: normalizeProvider(defaultProvider),
		uploaders:       uploaders,
	}
}

func (g *gateway) Save(file *multipart.FileHeader, scene, provider string) (*UploadedFile, error) {
	selected := normalizeProvider(provider)
	if selected == "" {
		selected = g.defaultProvider
	}

	uploader, ok := g.uploaders[selected]
	if !ok || uploader == nil {
		return nil, fmt.Errorf("upload provider %s is not configured", selected)
	}

	return uploader.Save(file, scene)
}

func normalizeProvider(provider string) string {
	switch strings.TrimSpace(strings.ToLower(provider)) {
	case "", "default":
		return ""
	case "oss", "aliyun", "aliyun-oss":
		return "aliyun-oss"
	case "cos", "tencent", "tencent-cos":
		return "tencent-cos"
	case "obs", "huawei", "huawei-obs":
		return "huawei-obs"
	case "local":
		return "local"
	default:
		return strings.TrimSpace(strings.ToLower(provider))
	}
}
