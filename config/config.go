package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// Config is type config is unmarshelled to
type Config struct {
	ServerPort   string `json:"serverPort"`
	LogDirectory string `json:"logDirectory"`
	DbConnString string `json:"dbConnString"`
}

//ReadConfig unmarshalles config file to Config ̰
func ReadConfig(configFile string) (*Config, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("Config: error when opening file %s : %s ", configFile, err)
	}
	defer f.Close()

	var bytes []byte
	bytes, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Config: error when reading file %s : %s ", configFile, err)
	}

	var config *Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("Config: error when reading file %s : %s ", configFile, err)
	}
	return config, nil
}

//GetServerPort returns a server port with format :{port}
func (config *Config) GetServerPort() string {
	return fmt.Sprintf(":%s", config.ServerPort)
}

//Validate ensures valid configuration values
func (config *Config) Validate() error {
	if len(config.ServerPort) == 0 {
		return fmt.Errorf("Config: ServerPort must not be empty")
	}

	serverPort, err := strconv.Atoi(config.ServerPort)
	if err != nil {
		return fmt.Errorf("Config: ServerPort must be a number")
	}

	if serverPort <= 0 {
		return fmt.Errorf("Config: Invalid ServerPort")
	}

	if len(config.DbConnString) == 0 {
		return fmt.Errorf("Config: DbConnString must not be empty")
	}
	return nil
}
