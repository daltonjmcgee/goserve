package config

import (
	"encoding/json"
	"goserve/helpers"
)

func ReturnConfig(configFile string) map[string]string {
	var config map[string]string
	configJson, _ := helpers.LoadFile(configFile)
	json.Unmarshal([]byte(configJson), &config)
	return config
}
