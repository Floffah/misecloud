package env

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var CliViper = viper.New()
var ControllerViper = viper.New()
var AgentViper = viper.New()

const (
	CliViperKeyPort = "controller.port"
)

const (
	ControllerViperKeyPort    = "api.port"
	ControllerViperKeyIsSetup = "is_setup"
)

func InitCliViper() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	miseDir := filepath.Join(home, ".mise")
	if err := os.MkdirAll(miseDir, 0755); err != nil {
		return err
	}

	miseCliTomlPath := filepath.Join(miseDir, "misecli.toml")

	CliViper.SetConfigName("misecli")
	CliViper.SetConfigType("toml")

	CliViper.SetDefault(CliViperKeyPort, "50051")

	CliViper.AddConfigPath(miseDir)

	CliViper.SetEnvPrefix("misecl")
	CliViper.AutomaticEnv()

	var fileLookupError viper.ConfigFileNotFoundError
	if err := CliViper.ReadInConfig(); err != nil {
		if errors.As(err, &fileLookupError) {
			err := CliViper.WriteConfigAs(miseCliTomlPath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = CliViper.WriteConfigAs(miseCliTomlPath)
	if err != nil {
		return err
	}

	return nil
}

func InitControllerViper() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	miseDir := filepath.Join(home, ".mise")
	if err := os.MkdirAll(miseDir, 0755); err != nil {
		return err
	}

	miseControllerTomlPath := filepath.Join(miseDir, "misecontroller.toml")

	ControllerViper.SetConfigName("misecontroller")
	ControllerViper.SetConfigType("toml")

	ControllerViper.SetDefault(ControllerViperKeyIsSetup, false)
	ControllerViper.SetDefault(ControllerViperKeyPort, "50051")

	ControllerViper.AddConfigPath(miseDir)

	ControllerViper.SetEnvPrefix("miseco")
	ControllerViper.AutomaticEnv()

	var fileLookupError viper.ConfigFileNotFoundError
	if err := ControllerViper.ReadInConfig(); err != nil {
		if errors.As(err, &fileLookupError) {
			err := ControllerViper.WriteConfigAs(miseControllerTomlPath)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = ControllerViper.WriteConfigAs(miseControllerTomlPath)
	if err != nil {
		return err
	}

	return nil
}
