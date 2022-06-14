package config

import (
	"github.com/spf13/viper"
	"log"
)

var minioConfig *MinioConfig

func initMinioConfig() {
	minioConfig = new(MinioConfig)
	minioConfig.ReadConfig(aviper)
}

func GetMinioConfig() *MinioConfig {
	return minioConfig
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
}

func (t *MinioConfig) ReadConfig(v *viper.Viper) {
	t.Endpoint = v.GetString("minio.endpoint")
	t.AccessKey = v.GetString("minio.accessKey")
	t.SecretKey = v.GetString("minio.secretKey")
	t.Bucket = v.GetString("minio.bucket")

	if t.Endpoint == "" {
		log.Fatalf("minio.endpoint must not be empty")
	}
	if t.AccessKey == "" {
		log.Fatalf("minio.accessKey must not be empty")
	}
	if t.SecretKey == "" {
		log.Fatalf("minio.secretKey must not be empty")
	}
	if t.Bucket == "" {
		log.Fatalf("minio.bucket must not be empty")
	}
}
