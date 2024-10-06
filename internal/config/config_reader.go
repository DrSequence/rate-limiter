package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadFile() Config {
	var cfg Config
	f, err := os.Open("config/application_properties.yaml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
	}

	return cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
