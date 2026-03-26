package dao

import (
	"fmt"
	"sync"
	"time"

	v2config "go_sleep_admin/internal/platform/config"

	"github.com/gtkit/logger"
	"github.com/gtkit/orm"
	"github.com/gtkit/orm/zlogger"
	"github.com/gtkit/redis"
	"gorm.io/gorm"
)

// NewWithConfig 使用 typed config 初始化数据层。
func NewWithConfig(cfg *v2config.Config) {
	daoDB = &dao{
		rdb:       initRedisCollectionWithConfig(cfg),
		mdb:       initMysqlWithConfig(cfg),
		providers: make(map[string]any),
		mu:        sync.RWMutex{},
	}
}

func initMysqlWithConfig(cfg *v2config.Config) *gorm.DB {
	orm.MysqlConfig(
		orm.Host(cfg.Database.Host),
		orm.Port(cfg.Database.Port),
		orm.DbType(cfg.Database.DBType),
		orm.Name(cfg.Database.Name),
		orm.User(cfg.Database.Username),
		orm.WithPassword(cfg.Database.Password),
	)

	loggeropt := []zlogger.Option{
		zlogger.WithLogger(logger.Zlog()),
	}
	if cfg.Log.SQL != 0 {
		loggeropt = append(loggeropt, zlogger.WithSqlLog())
	}

	log := zlogger.New(loggeropt...)
	orm.GormConfig(
		orm.PrepareStmt(true),
		orm.SkipDefaultTransaction(true),
		orm.GormLog(log),
	)

	return orm.NewMysql(&dbSetFromConfig{cfg: cfg})
}

type dbSetFromConfig struct {
	cfg *v2config.Config
}

func (d *dbSetFromConfig) Set(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(d.cfg.Database.MaxOpenConn)
	sqlDB.SetMaxIdleConns(d.cfg.Database.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(d.cfg.Database.MaxLifeSeconds) * time.Second)
	logger.Info("Mysql Custom set done!")
}

func initRedisCollectionWithConfig(cfg *v2config.Config) map[int]*redis.Redisclient {
	var dbsconn []redis.ConnConfigOption
	dbsconn = append(dbsconn,
		redis.WithAddr(cfg.Redis.Addr),
		redis.WithPassword(cfg.Redis.Password),
	)

	for _, db := range cfg.Redis.DBs {
		dbsconn = append(dbsconn, redis.WithDB(db))
	}

	return redis.NewCollection(dbsconn...)
}

// InitProviderInitializer 允许 future provider 在 runtime 启动时注册，例如 ES。
type InitProviderInitializer interface {
	Name() string
	Init(cfg *v2config.Config) (any, error)
}

// InitExtraProviders 按顺序初始化额外 provider 并注册进 data 层。
func InitExtraProviders(cfg *v2config.Config, initializers ...InitProviderInitializer) error {
	for _, initializer := range initializers {
		if initializer == nil {
			continue
		}

		provider, err := initializer.Init(cfg)
		if err != nil {
			return fmt.Errorf("init provider %s: %w", initializer.Name(), err)
		}
		if provider == nil {
			continue
		}

		if err := RegisterProvider(initializer.Name(), provider); err != nil {
			return fmt.Errorf("register provider %s: %w", initializer.Name(), err)
		}
	}

	return nil
}
