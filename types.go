package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/twsm000/lenslocked/models/database/postgres"
	"github.com/twsm000/lenslocked/models/services"
)

type EnvConfig struct {
	CSRF       CSRF                `json:"csrf"`
	DBConfig   postgres.Config     `json:"database"`
	Server     Server              `json:"server"`
	Session    Session             `json:"session"`
	SMTPConfig services.SMTPConfig `json:"smtp"`
}

func LoadEnvSettings(fpath, dbDriver string) (*EnvConfig, error) {
	env := EnvConfig{
		DBConfig: postgres.Config{
			Driver: dbDriver,
		},
	}

	fpath, err := filepath.Abs(fpath)
	if err != nil {
		return nil, err
	}
	logInfo.Println("EnvSettingsFilePath:", fpath)

	data, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	envData := bytes.NewBuffer(data)
	decoder := json.NewDecoder(envData)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

type CSRF struct {
	Key    string `json:"key"`
	Secure bool   `json:"secure"`
}

type Server struct {
	Address string `json:"address"`
}

type Session struct {
	TokenSize int `json:"token_size"`
}
