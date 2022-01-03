package awscmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"
)

type InputConfigure struct {
	Region    string
	Profile   string
	AccessKey string
	SecretKey string
}

type OutputConfigure struct{}

func Configure(ctx context.Context, input *InputConfigure) (*OutputConfigure, error) {
	path, err := homedir.Expand("~/.aws")
	if err != nil {
		return nil, fmt.Errorf("Failed to expand path: %v", err)
	}
	err = os.Mkdir(path, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("Failed to create .aws directory: %v", err)
	}

	configPath := filepath.Join(path, "config")
	credentialsPath := filepath.Join(path, "credentials")

	config, err := iniLoadOrCreate(configPath)
	if err != nil {
		return nil, fmt.Errorf("Fail to read config file: %v", err)
	}
	config.Section(fmt.Sprint("profile", " ", input.Profile)).NewKey("region", input.Region)
	err = config.SaveTo(configPath)
	if err != nil {
		return nil, fmt.Errorf("Fail to save config file: %v", err)
	}

	credentials, err := iniLoadOrCreate(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("Fail to read credentials file: %v", err)
	}
	sec := credentials.Section(input.Profile)
	sec.NewKey("aws_access_key_id", input.AccessKey)
	sec.NewKey("aws_secret_access_key", input.SecretKey)

	err = credentials.SaveTo(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("Fail to save config file: %v", err)
	}

	return &OutputConfigure{}, nil
}

func iniLoadOrCreate(path string) (*ini.File, error) {
	_, err := os.Stat(path)
	if err == nil {
		return ini.Load(path)
	}

	_, err = os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to create file %s", path)
	}

	return ini.Load(path)
}
