package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var dbConfig *DbConfig

func initDbConfig() {
	dbConfig = new(DbConfig)
	dbConfig.ReadConfig(aviper)
}

func GetDbConfig() *DbConfig {
	return dbConfig
}

type DbConfig struct {
	Host     string
	Port     int64
	Username string
	Password string
	Database string
}

func (t *DbConfig) ReadConfig(v *viper.Viper) {
	t.Host = v.GetString("db.host")
	t.Port = v.GetInt64("db.port")
	t.Username = v.GetString("db.username")
	t.Password = v.GetString("db.password")
	t.Database = v.GetString("db.database")
	if t.Host == "" {
		log.Fatalf("db.host can't be empty")
	}
	if t.Port == 0 {
		log.Fatalf("db.port can't be empty")
	}
	if t.Database == "" {
		log.Fatalf("db.database can't be empty")
	}
}

func (t DbConfig) GetConnection() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", t.Username, t.Password, t.Host, t.Port, t.Database)
}
