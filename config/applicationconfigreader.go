package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// ApplicationConfig is type appconfig is unmarshelled to
type ApplicationConfig struct {
	ServerPort   string `json:"serverPort"`
	LogDirectory string `json:"logDirectory"`
}

//ReadApplicationConfig unmarshalles appconfig file to ApplicationConfig ̰
func ReadApplicationConfig(applicationConfigFilePath string) (a *ApplicationConfig, err error) {
	f, err := os.Open(applicationConfigFilePath)
	if err != nil {
		return nil, fmt.Errorf("ApplicationConfig: Error when opening file %s : %v ", applicationConfigFilePath, err)
	}
	defer f.Close()
	var bytes []byte
	bytes, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("ApplicationConfig: Error when reading file %s : %v ", applicationConfigFilePath, err)
	}

	err = json.Unmarshal(bytes, &a)

	if err != nil {
		return nil, fmt.Errorf("ApplicationConfig: Error when reading file %s : %v ", applicationConfigFilePath, err)
	}
	return a, nil
}

//GetServerPort returns a server port with format :{port}
func (appConfig *ApplicationConfig) GetServerPort() string {
	return fmt.Sprintf(":%s", appConfig.ServerPort)
}

//Validate ensures valid configuration values
func (appConfig *ApplicationConfig) Validate() error {
	if len(appConfig.ServerPort) == 0 {
		return errors.New("Config: ServerPort must not be empty")
	}

	serverPort, err := strconv.Atoi(appConfig.ServerPort)
	if err != nil {
		return errors.New("Config: ServerPort must be a number")
	}
	if serverPort <= 0 {
		return errors.New("Config: Invalid ServerPort")
	}
	return nil
}
