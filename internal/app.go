package app

import (
	"os"
	"time"

	"github.com/Tairascii/google-docs-organization/internal/app/usecase"
	"gopkg.in/yaml.v3"
)

const (
	configFilePath = "CONFIG_FILE_PATH"
)

type UseCase struct {
	Org usecase.OrgUseCase
}

type DI struct {
	UseCase UseCase
}

type Config struct {
	Repo struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		DBName       string `yaml:"dbname"`
		Schema       string `yaml:"schema"`
		AppName      string `yaml:"app_name"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
		MaxOpenConns int    `yaml:"max_open_conns"`
	} `yaml:"repo"`
	Server struct {
		Port    string `yaml:"port"`
		Timeout struct {
			Read  time.Duration `yaml:"read"`
			Write time.Duration `yaml:"write"`
			Idle  time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
	GrpcServer struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"grpc_server"`
}

func LoadConfigs() (*Config, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg := &Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
