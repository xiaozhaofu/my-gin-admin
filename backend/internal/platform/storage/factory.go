package storage

import v2config "go_sleep_admin/internal/platform/config"

func NewUploaders(cfg v2config.UploadConfig) (map[string]Uploader, error) {
	maxSize := cfg.MaxSizeMB << 20
	uploaders := map[string]Uploader{
		"local": NewLocalUploader(cfg.LocalDir, cfg.PublicPath, cfg.AllowedExts, maxSize),
	}

	if cfg.OSS.Bucket != "" && cfg.OSS.AccessKeyID != "" && cfg.OSS.AccessKeySecret != "" {
		uploaders["aliyun-oss"] = NewOSSUploader(cfg.OSS, cfg.AllowedExts, maxSize)
	}
	if cfg.COS.BucketURL != "" && cfg.COS.SecretID != "" && cfg.COS.SecretKey != "" {
		uploader, err := NewCOSUploader(cfg.COS, cfg.AllowedExts, maxSize)
		if err != nil {
			return nil, err
		}
		uploaders["tencent-cos"] = uploader
	}
	if cfg.OBS.Endpoint != "" && cfg.OBS.Bucket != "" && cfg.OBS.AccessKeyID != "" && cfg.OBS.AccessKeySecret != "" {
		uploader, err := NewOBSUploader(cfg.OBS, cfg.AllowedExts, maxSize)
		if err != nil {
			return nil, err
		}
		uploaders["huawei-obs"] = uploader
	}

	return uploaders, nil
}

func NewUploadGatewayWithConfig(cfg v2config.UploadConfig) (UploadGateway, error) {
	uploaders, err := NewUploaders(cfg)
	if err != nil {
		return nil, err
	}
	return NewUploadGateway(cfg.Driver, uploaders), nil
}
