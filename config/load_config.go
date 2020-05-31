package config

import (
	"gcfs/config/data"
	"github.com/jinzhu/configor"
	"os"
)

var Main Configuration

type Configuration struct {
	Database data.Database          `json:"database"`
	Server   data.Server            `json:"server"`
	Routes   data.Routes            `json:"routes"`
	Metadata map[string]interface{} `json:"metadata"`
	Global data.Global `json:"global"`
}

func init() {
	PathConfFile := os.Getenv("GCFS_CONFIG")

	if PathConfFile == "" {
		panic("no env variable found for GCFS_CONFIG.")
	}

	if err := configor.Load(&Main, PathConfFile); err != nil {
		panic(err.Error())
	}
}