package config

import "github.com/spf13/viper"

var filesConfig *FilesConfig

const defaultMaxSize = 1 * 1024 * 1024

func initFilesConfig() {
	filesConfig = new(FilesConfig)
	filesConfig.ReadConfig(aviper)
}

func GetFilesConfig() *FilesConfig {
	return filesConfig
}

type FilesConfig struct {
	MaxSize int64 // 最大上传文件大小byte
}

func (t *FilesConfig) ReadConfig(v *viper.Viper) {
	viper.SetDefault("files.maxSize", defaultMaxSize)
	t.MaxSize = viper.GetInt64("files.maxSize")
}
