package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
	"path/filepath"
	"strings"
)

func ParseConfig(cfgPath string) (map[string]interface{}, error) {
	var err error
	var cfgPathAbs string
	if cfgPathAbs, err = filepath.Abs(strings.TrimSpace(cfgPath)); err != nil {
		return nil, err
	}

	file_base, file_name := filepath.Split(cfgPathAbs)
	file_type := strings.TrimPrefix(path.Ext(file_name), ".")
	file_name = strings.TrimSuffix(file_name, path.Ext(file_name))

	switch strings.ToLower(file_type) {
	case "json", "yaml":
	default:
		return nil, fmt.Errorf("unsupported file extension")
	}

	cfg := viper.New()
	cfg.AddConfigPath(file_base)
	cfg.SetConfigType(file_type)
	cfg.SetConfigName(file_name)

	if err := cfg.ReadInConfig(); err != nil {
		return nil, err
	}
	return cfg.AllSettings(), nil
}
