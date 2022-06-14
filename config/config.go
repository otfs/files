package config

import (
	"embed"
	"github.com/spf13/viper"
)

var aviper = viper.GetViper()

func InitConfig(configFs embed.FS) {
	configViper(aviper, configFs)
	initFilesConfig()
	initMinioConfig()
	initDbConfig()
}
