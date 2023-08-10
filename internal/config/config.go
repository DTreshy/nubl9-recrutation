package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Log LogConfig `koanf:"log"`
	Net NetConfig `koanf:"net"`
}

type LogConfig struct {
	File  string `koanf:"file"`
	Level string `koanf:"level"`
}

type NetConfig struct {
	HTTPBind string `koanf:"http_bind"`
}

func New(path string) (*Config, error) {
	var k = koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, err
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	cfg.Log.Level = strings.ToLower(cfg.Log.Level)

	return &cfg, nil
}
