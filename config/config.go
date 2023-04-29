package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RspamdConnectionConfig struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ImapConfig struct {
	Server    string `yaml:"server"`
	Port      uint   `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	UseTLS    bool   `yaml:"useTLS"`
	Directory string `yaml:"directory"`
}

type Config struct {
	Rspamd RspamdConnectionConfig `yaml:"rspamd"`
	Inputs []ImapConfig           `yaml:"inputs"`
}

func LoadFromFile(pathToFile string) (cfg Config, err error) {
	data, err := os.ReadFile(pathToFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &cfg)
	return
}
