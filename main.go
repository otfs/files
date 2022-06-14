package main

import (
	"embed"
	"files/config"
	_ "files/config"
	"files/files"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

//go:embed config*.yml
var configFs embed.FS

func main() {
	config.InitConfig(configFs)
	g := gin.Default()
	route(g)
	err := g.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func route(g *gin.Engine) {
	dbConfig := config.GetDbConfig()
	filesConfig := config.GetFilesConfig()
	minioConfig := config.GetMinioConfig()
	db := sqlx.MustOpen("mysql", dbConfig.GetConnection())
	filesRepository := files.NewFilesRepository(db)
	filesService := files.NewFilesService(filesRepository, minioConfig)
	files.NewFilesController(g, filesService, filesConfig)
}
