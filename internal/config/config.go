package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Name    string `yaml:"Name"`
	Host    string `yaml:"Host"`
	Port    string `yaml:"Port"`
	SqlLite struct {
		File string `yaml:"File"`
	} `yaml:"SqlLite"`
	Storage string `yaml:"Storage"`
}

var config Config

func LoadConfig(file string) *Config {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		panic(err)
	}

	return &config
}
