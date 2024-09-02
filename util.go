package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func readConfigFile(filename string, v any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(v)
}
