package jwtauth

import (
	v2config "go_sleep_admin/internal/platform/config"
	"go_sleep_admin/internal/runtime/resource"

	"github.com/gtkit/encry/jwt"
)

// NewWithConfig creates a fresh JWT issuer without mutating package globals.
func NewWithConfig(*v2config.Config) (*jwt.JwtEd25519, error) {
	privateKeyPath := resource.ResolvePath("config/pem/jwtpri.pem")
	publicKeyPath := resource.ResolvePath("config/pem/jwtpub.pem")
	return jwt.NewJwtEd25519(privateKeyPath, publicKeyPath)
}

func JwtEd25519() *jwt.JwtEd25519 {
	issuer, err := NewWithConfig(nil)
	if err != nil {
		panic(err)
	}

	return issuer
}
