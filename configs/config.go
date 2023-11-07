package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type DbDsnCfg struct {
	User         string `yaml:"user"`
	DbName       string `yaml:"dbname"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Sslmode      string `yaml:"sslmode"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	Timer        uint32 `yaml:"timer"`
}

func ReadConfig() (*DbDsnCfg, error) {
	dsnConfig := DbDsnCfg{}
	dsnFile, err := os.ReadFile("configs/db_dsn.yaml")
	if err != nil {
		return nil, fmt.Errorf("ReadConfig read file err: %w", err)
	}

	err = yaml.Unmarshal(dsnFile, &dsnConfig)
	if err != nil {
		return nil, fmt.Errorf("ReadConfig unmarshal err: %w", err)
	}

	return &dsnConfig, nil
}
