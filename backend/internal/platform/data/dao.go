package dao

import (
	"fmt"
	"slices"
	"sync"

	"github.com/gtkit/logger"
	"github.com/gtkit/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	daoDB Dao
	_     Dao = (*dao)(nil)
)

type (
	MyDB = *gorm.DB
	RDB  = *redis.Redisclient
)

type Dao interface {
	Mdb() MyDB
	Rdb(db int) RDB
	MdbClose() error
	RdbClose() error
	RegisterProvider(name string, provider any) error
	Provider(name string) (any, bool)
	d()
}

type dao struct {
	rdb       map[int]RDB
	mdb       MyDB
	providers map[string]any
	mu        sync.RWMutex
}

func (d *dao) d() {}

func Mdb() MyDB {
	return DB().Mdb()
}

func (d *dao) Mdb() MyDB {
	return d.mdb
}

func (d *dao) MdbClose() error {
	db, err := d.mdb.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func Rdb(db int) RDB {
	return DB().Rdb(db)
}

func RdbExec(db int, fn func(RDB) error) error {
	return fn(Rdb(db))
}

func (d *dao) Rdb(db int) RDB {
	if rdb, ok := d.rdb[db]; ok {
		return rdb
	}

	return redis.Select(db)
}

func (d *dao) RdbClose() error {
	for _, db := range redisDBIndices(d) {
		rdb, ok := d.rdb[db]
		if !ok || rdb == nil || rdb.Client() == nil {
			continue
		}

		if err := rdb.Client().Close(); err != nil {
			return err
		}
	}
	return nil
}

func redisDBIndices(d *dao) []int {
	if d == nil || len(d.rdb) == 0 {
		return nil
	}

	dbs := make([]int, 0, len(d.rdb))
	for db := range d.rdb {
		dbs = append(dbs, db)
	}
	slices.Sort(dbs)
	return dbs
}

func DB() Dao {
	return daoDB
}

func DBClose() {
	if DB() == nil {
		return
	}

	if err := DB().MdbClose(); err != nil {
		logger.Error("[*]Mysql close error", zap.Error(err))
	}
	logger.Info("[*]Mysql close success")
	if err := DB().RdbClose(); err != nil {
		logger.Error("[*]Redis close error", zap.Error(err))
	}
	logger.Info("[*]Redis close success")
}

func RegisterProvider(name string, provider any) error {
	if DB() == nil {
		return fmt.Errorf("data layer is not initialized")
	}

	return DB().RegisterProvider(name, provider)
}

func (d *dao) RegisterProvider(name string, provider any) error {
	if name == "" {
		return fmt.Errorf("provider name is required")
	}
	if provider == nil {
		return fmt.Errorf("provider %q is nil", name)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.providers[name] = provider

	return nil
}

func Provider(name string) (any, bool) {
	if DB() == nil {
		return nil, false
	}

	return DB().Provider(name)
}

func (d *dao) Provider(name string) (any, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	provider, ok := d.providers[name]

	return provider, ok
}
