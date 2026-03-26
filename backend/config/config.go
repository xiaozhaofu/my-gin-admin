package config

import (
	"bytes"
	"embed"
	"io"

	"github.com/gtkit/json"
)

//go:embed env
var conf embed.FS

func init() {
	json.RegisterFuzzyDecoders()
}

func EmbeddedConfig(env string) io.Reader {
	b, _ := conf.ReadFile("env/" + env + ".yml")
	return bytes.NewReader(b)
}
