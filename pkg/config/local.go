package config

import (
	"fmt"

	configdata "github.com/ijlik/dating-user/pkg/config/data"
	"github.com/joho/godotenv"
)

func getenv(path string) configdata.Config {
	pathEnv := fmt.Sprintf("%s/.env", path)
	if path == "" {
		pathEnv = ".env"
	}

	data, err := godotenv.Read(pathEnv)
	if err != nil {
		return nil
	}

	var mapData = make(map[string]interface{})

	for index, val := range data {
		mapData[index] = val
	}

	return &configdata.ConfigData{
		Data: mapData,
	}
}
