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

func createConfigFileFromStruct(t *testing.T, config *u.Config, configFile string) {
	bytes, err := json.Marshal(config)
	assert.NoError(t, err, "Utils test createConfigFileFromStruct failed")

	err = ioutil.WriteFile(configFile, bytes, 0644)
	assert.NoError(t, err, "Utils test createConfigFileFromStruct failed")
}

func removeConfigFile(t *testing.T, configFile string) {
	err := os.Remove(configFile)
	assert.NoError(t, err, "Failed in Utils test teardown. Unable to remove temp config file")
}

func TestReadConfigReadSuccessfullyFromConfigFile(t *testing.T) {
	config := u.Config{
		ServerPort:   "8080",
		LogDirectory: "logs",
		DbConnString: "actual-db-conn-string",
	}

	configFile := path.Join(os.TempDir(), "/default.conf")

	createConfigFileFromStruct(t, &config, configFile)

	actualConfig, err := u.ReadConfig(configFile)

	assert.NoError(t, err, "error when reading config")
	assert.Equal(t, "8080", actualConfig.ServerPort)
	assert.Equal(t, "logs", actualConfig.LogDirectory)
	assert.Equal(t, "actual-db-conn-string", actualConfig.DbConnString)

	removeConfigFile(t, configFile)
}

func TestReadConfigShouldReturnErrorWhenConfigFilePathIsInvalid(t *testing.T) {

	configFile := "/does_not_exist/default.conf"
	actualConfig, err := u.ReadConfig(configFile)

	assert.Nil(t, actualConfig)
	assert.Error(t, err, "no error when reading an invalid path config file")
	assert.Contains(t, err.Error(), "error when opening file /does_not_exist/default.conf")
}

func TestReadConfigShouldReturnErrorInCaseInvalidConfigJsonFormat(t *testing.T) {

	configFile := path.Join(os.TempDir(), "/default.conf")
	err := ioutil.WriteFile(configFile, []byte("garbage text..."), 0644)
	require.NoError(t, err, "failed when writing to a file")

	actualConfig, err := u.ReadConfig(configFile)

	assert.Nil(t, actualConfig)
	assert.Error(t, err, "no error on reading invalid config file")
	assert.Contains(t, err.Error(), fmt.Sprintf("error when reading file %s", configFile))

	removeConfigFile(t, configFile)
}

func TestEmptyServerPort(t *testing.T) {
	config := u.Config{}

	err := config.Validate()
	require.Error(t, err, "Expected error in config validation")

	assert.Equal(t, "Config: ServerPort must not be empty", err.Error())
}

func TestNonNumericServerPort(t *testing.T) {
	config := u.Config{
		ServerPort: "eight-zero-eight-zero",
	}

	err := config.Validate()
	require.Error(t, err, "Expected error in config validation")
	assert.Equal(t, "Config: ServerPort must be a number", err.Error())
}

func TestInvalidServerPort(t *testing.T) {
	config := u.Config{
		ServerPort: "0",
	}
	err := config.Validate()
	require.Error(t, err, "Expected error in config validation")
	assert.Equal(t, "Config: Invalid ServerPort", err.Error())
}

func TestEmptyDbConnString(t *testing.T) {
	config := u.Config{
		ServerPort: "8080",
	}
	err := config.Validate()
	require.Error(t, err, "Expected error in config validation")
	assert.Equal(t, "Config: DbConnString must not be empty", err.Error())
}

func TestGetServerport(t *testing.T) {
	config := u.Config{
		ServerPort:   "8080",
		DbConnString: "some-conn-string-value",
	}

	err := config.Validate()
	require.NoError(t, err, "Expected no error in config validation")

	serverPort := config.GetServerPort()
	assert.Equal(t, ":8080", serverPort)
}
