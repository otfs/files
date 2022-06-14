package config

import (
	"embed"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	ConfigFileNameMain   = "config.yml"
	ConfigFileNameFormat = "config-%s.yml"
	ConfigFileType       = "yml"
	ConfigEnvKey         = "env"
	ConfigEnvDefault     = "default"
)

//
// configViper 支持多环境配置初始化配置
//
func configViper(v *viper.Viper, cfgFiles embed.FS) {
	v.SetConfigType(ConfigFileType)
	// 主要配置文件
	mainFile, err := cfgFiles.Open(ConfigFileNameMain)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("read config file error: %+v", err)
	}
	if err == nil {
		if err := v.ReadConfig(mainFile); err != nil {
			log.Fatalf("read config file error: %+v", err)
		}
	}

	// 环境
	env := v.GetString(ConfigEnvKey)
	if env == "" {
		env = ConfigEnvDefault
	}
	log.Printf("active env: %s", env)

	// 环境配置文件
	envFileName := fmt.Sprintf(ConfigFileNameFormat, env)
	envFile, err := cfgFiles.Open(envFileName)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("read config file error: %+v", err)
	}
	if err == nil {
		if err := v.MergeConfig(envFile); err != nil {
			log.Fatalf("read config file error: %+v", err)
		}
	}
}
