package config

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Validate 校验运行所需的关键配置。
func (c *Config) Validate() error {
	var errs []error

	if c.Env == "" {
		errs = append(errs, errors.New("env is required"))
	}

	if c.App.Name == "" {
		errs = append(errs, errors.New("application.name is required"))
	}
	if c.App.Host == "" {
		errs = append(errs, errors.New("application.host is required"))
	}
	if c.App.Port == "" {
		errs = append(errs, errors.New("application.port is required"))
	} else if _, err := strconv.Atoi(c.App.Port); err != nil {
		errs = append(errs, fmt.Errorf("application.port must be numeric: %w", err))
	}
	if c.App.ReadTimeoutSeconds < 0 {
		errs = append(errs, errors.New("application.readtimeout must be >= 0"))
	}
	if c.App.WriteTimeoutSeconds < 0 {
		errs = append(errs, errors.New("application.writertimeout must be >= 0"))
	}
	if c.App.Timezone != "" {
		if _, err := time.LoadLocation(c.App.Timezone); err != nil {
			errs = append(errs, fmt.Errorf("application.timezone is invalid: %w", err))
		}
	}

	if c.Database.DBType == "" {
		errs = append(errs, errors.New("database.dbtype is required"))
	}
	if c.Database.Host == "" {
		errs = append(errs, errors.New("database.host is required"))
	}
	if c.Database.Port == "" {
		errs = append(errs, errors.New("database.port is required"))
	} else if _, err := strconv.Atoi(c.Database.Port); err != nil {
		errs = append(errs, fmt.Errorf("database.port must be numeric: %w", err))
	}
	if c.Database.Name == "" {
		errs = append(errs, errors.New("database.name is required"))
	}
	if c.Database.Username == "" {
		errs = append(errs, errors.New("database.username is required"))
	}

	if c.Redis.Addr == "" {
		errs = append(errs, errors.New("redis.addr is required"))
	}
	if len(c.Redis.DBs) == 0 {
		errs = append(errs, errors.New("redis.dbs must not be empty"))
	}

	if c.Log.Level == "" {
		errs = append(errs, errors.New("log.level is required"))
	}
	if c.Log.Path == "" {
		errs = append(errs, errors.New("log.path is required"))
	}

	if c.JWT.Timeout <= 0 {
		errs = append(errs, errors.New("jwt.timeout must be > 0"))
	}
	if c.JWT.RefreshTimeout <= 0 {
		errs = append(errs, errors.New("jwt.refresh_timeout must be > 0"))
	}

	if c.Upload.Driver == "" {
		errs = append(errs, errors.New("upload.driver is required"))
	}
	if c.Upload.LocalDir == "" {
		errs = append(errs, errors.New("upload.local_dir is required"))
	}
	if c.Upload.PublicPath == "" {
		errs = append(errs, errors.New("upload.public_path is required"))
	}
	if c.Upload.MaxSizeMB <= 0 {
		errs = append(errs, errors.New("upload.max_size_mb must be > 0"))
	}

	if c.SMS.Enabled {
		if c.Aliyun.SMS.AccessID == "" {
			errs = append(errs, errors.New("aliyun.sms_access_id is required when sms.enabled is true"))
		}
		if c.Aliyun.SMS.AccessSecret == "" {
			errs = append(errs, errors.New("aliyun.sms_access_secret is required when sms.enabled is true"))
		}
		if c.Aliyun.SMS.SignName == "" {
			errs = append(errs, errors.New("aliyun.sms_sign_name is required when sms.enabled is true"))
		}
		if c.Aliyun.SMS.TemplateCode == "" {
			errs = append(errs, errors.New("aliyun.sms_template_code is required when sms.enabled is true"))
		}
	}

	return errors.Join(errs...)
}
