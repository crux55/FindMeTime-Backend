package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DatabaseConfig struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		DatabaseName string `yaml:"databaseName"`
	} `yaml:"database"`
}

func LoadConfig(file string) (*Config, error) {
	var config Config
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
