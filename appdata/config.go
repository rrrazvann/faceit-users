package appdata

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Config *globalConfig

type globalConfig struct {
	Env      string         `yaml:"env"`
	Api      apiConfig      `yaml:"api"`
	Database databaseConfig `yaml:"database"`
}

type apiConfig struct {
	ListenHost string `yaml:"listen_host"`
}

type databaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
}

func newConfig(configPath string) (*globalConfig, error) {
	config := &globalConfig{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
