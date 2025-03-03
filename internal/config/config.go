package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DB_URL          string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const (
	config_file_name = ".gatorconfig.json"
)

func Read() (Config, error) {
	c := Config{}

	//load local json
	path, err := os.UserHomeDir()
	if err != nil {
		return c, err
	}

	file, err := os.Open(path + "/" + config_file_name)
	if err != nil {
		return c, err
	}

	//transfer content of json to struct
	data, err := io.ReadAll(file)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(data, &c)
	return c, err
}

func (conf *Config) SetUser(user_name string) error {
	conf.CurrentUserName = user_name

	data, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	path, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = os.WriteFile(path+"/"+config_file_name, data, 0666) //TODO: FileMode is from documentation example, find out what it means
	if err != nil {
		return err
	}

	return nil
}
