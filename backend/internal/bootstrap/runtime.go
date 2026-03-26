package bootstrap

import (
	"fmt"
	"sync/atomic"

	"github.com/gtkit/encry/jwt"

	v2config "go_sleep_admin/internal/platform/config"
	legacydao "go_sleep_admin/internal/platform/data"
)

// Runtime represents process-wide runtime dependencies.
type Runtime struct {
	options     Options
	config      *v2config.Config
	db          legacydao.Dao
	tokenIssuer *jwt.JwtEd25519
	closed      atomic.Bool
}

// NewRuntime builds runtime dependencies that can be reused by commands.
func NewRuntime(options Options) (*Runtime, error) {
	if options.EnvFile == "" {
		return nil, fmt.Errorf("env file is required")
	}

	cfg, err := v2config.Load(options.EnvFile)
	if err != nil {
		return nil, err
	}

	if err := initLegacyRuntime(cfg); err != nil {
		return nil, err
	}

	db := legacydao.DB()
	if db == nil {
		closeLegacyRuntime()
		return nil, fmt.Errorf("legacy dao is not initialized")
	}

	tokenIssuer := currentJWTIssuer()
	if tokenIssuer == nil {
		closeLegacyRuntime()
		return nil, fmt.Errorf("legacy jwt issuer is not initialized")
	}

	return &Runtime{
		options:     options,
		config:      cfg,
		db:          db,
		tokenIssuer: tokenIssuer,
	}, nil
}

// Config returns the typed runtime configuration.
func (r *Runtime) Config() *v2config.Config {
	return r.config
}

// DB returns the legacy DAO while runtime wiring is still migrating.
func (r *Runtime) DB() legacydao.Dao {
	return r.db
}

// TokenIssuer returns the initialized JWT issuer for explicit injection.
func (r *Runtime) TokenIssuer() *jwt.JwtEd25519 {
	if r == nil {
		return nil
	}

	return r.tokenIssuer
}

// Close releases runtime resources exactly once.
func (r *Runtime) Close() {
	if r == nil || !r.closed.CompareAndSwap(false, true) {
		return
	}

	closeLegacyRuntime()
}
