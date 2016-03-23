package config_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	u "git.mailbox.com/mailbox/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createConfigFileFromStruct(t *testing.T, appConfig *u.ApplicationConfig, configFilePath string) {
	bytes, err := json.Marshal(appConfig)
	assert.NoError(t, err, "Utils test createConfigFileFromStruct failed")
	err = ioutil.WriteFile(configFilePath, bytes, 0644)
	assert.NoError(t, err, "Utils test createConfigFileFromStruct failed")
}

func removeConfigFile(t *testing.T, configFilePath string) {
	err := os.Remove(configFilePath)
	assert.NoError(t, err, "Failed in Utils test teardown. Unable to remove temp config file")
}

func TestReadApplicationConfigReadSuccessfullyFromConfigFile(t *testing.T) {

	appConfig := u.ApplicationConfig{
		ServerPort:   "8080",
		LogDirectory: "logs",
	}
	configFilePath := path.Join(os.TempDir(), "/default.conf")

	createConfigFileFromStruct(t, &appConfig, configFilePath)

	actualApplicationConfig, err := u.ReadApplicationConfig(configFilePath)

	assert.Nil(t, err)
	assert.Equal(t, actualApplicationConfig.ServerPort, "8080")
	assert.Equal(t, actualApplicationConfig.LogDirectory, "logs")

	removeConfigFile(t, configFilePath)
}

func TestReadApplicationConfigShouldReturnErrorInCaseConfigFilePathIsInvalid(t *testing.T) {

	configFilePath := "/does_not_exist/default.conf"
	actualApplicationConfig, err := u.ReadApplicationConfig(configFilePath)

	assert.Nil(t, actualApplicationConfig)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error when opening file /does_not_exist/default.conf")

}

func TestReadApplicationConfigShouldReturnErrorInCaseInvalidConfigJsonFormat(t *testing.T) {

	configFilePath := path.Join(os.TempDir(), "/default.conf")
	err := ioutil.WriteFile(configFilePath, []byte("garbage text..."), 0644)

	if err != nil {
		t.Error("Utils test createConfigFileFromStruct failed")
	}
	actualApplicationConfig, err := u.ReadApplicationConfig(configFilePath)

	assert.Nil(t, actualApplicationConfig)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("Error when reading file %v", configFilePath))

	removeConfigFile(t, configFilePath)
}

func TestEmptyServerPort(t *testing.T) {
	appConfig := u.ApplicationConfig{}
	err := appConfig.Validate()
	require.Error(t, err, "Expected error in config validation")
	assert.Equal(t, err.Error(), "Config: ServerPort must not be empty")
}

func TestNonNumericServerPort(t *testing.T) {
	appConfig := u.ApplicationConfig{
		ServerPort: "eight-zero-eight-zero",
	}
	err := appConfig.Validate()
	require.Error(t, err, "Expected error in config validation")
	assert.Equal(t, err.Error(), "Config: ServerPort must be a number")
}

func TestInvalidServerPort(t *testing.T) {
	appConfig := u.ApplicationConfig{
		ServerPort: "0",
	}
	err := appConfig.Validate()
	require.Error(t, err, "Expected error in config validation")
	assert.Equal(t, err.Error(), "Config: Invalid ServerPort")
}

func TestGetServerport(t *testing.T) {
	appConfig := u.ApplicationConfig{
		ServerPort: "8080",
	}
	err := appConfig.Validate()
	require.NoError(t, err, "Expected no error in config validation")

	serverPort := appConfig.GetServerPort()
	assert.Equal(t, ":8080", serverPort)
}
