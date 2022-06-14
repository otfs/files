package files

import (
	"files/config"
	"files/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"path"
	"strings"
)

type FilesController struct {
	filesService domain.FilesService
	filesConfig  *config.FilesConfig
}

func NewFilesController(g *gin.Engine, filesService domain.FilesService, filesConfig *config.FilesConfig) *FilesController {
	ctl := &FilesController{
		filesService: filesService,
		filesConfig:  filesConfig,
	}
	g.POST("/generate-upload-url", ctl.generateUploadUrl)
	g.POST("/upload", ctl.uploadFile)
	g.POST("/files/:id/generate-url", ctl.generateFileUrl)
	g.POST("/files/generate-url-batch", ctl.generateFileUrlBatch)
	g.GET("/files/:id", ctl.downloadFile)
	return ctl
}

// 生成文件预签名上传地址
func (t *FilesController) generateUploadUrl(c *gin.Context) {
	request := new(domain.GenerateUploadUrlRequest)
	err := c.BindJSON(request)
	if err != nil {
		c.JSON(400, domain.NewResultFail("400", err.Error()))
		return
	}
	response, err := t.filesService.GenerateUploadUrl(request)
	if err != nil {
		c.JSON(500, domain.NewResultFail("500", err.Error()))
		return
	}
	c.JSON(200, domain.NewResultOk(response))
}

// 上传文件
func (t *FilesController) uploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, domain.NewResultFail("500", err.Error()))
		return
	}
	if fileHeader.Size > t.filesConfig.MaxSize {
		msg := fmt.Sprintf("上传文件超过最大大小：%d", t.filesConfig.MaxSize)
		c.JSON(400, domain.NewResultFail("400", msg))
		return
	}

	response, err := t.filesService.UploadFile(fileHeader)
	if err != nil {
		c.JSON(500, domain.NewResultFail("500", err.Error()))
		return
	}
	c.JSON(200, domain.NewResultOk(response))
}

// 生成文件访问地址
func (t *FilesController) generateFileUrl(c *gin.Context) {
	id := c.Param("id")
	_, url, err := t.filesService.GenerateFileUrl(id)
	if err != nil {
		c.JSON(500, domain.NewResultFail("500", err.Error()))
		return
	}

	data := gin.H{"url": url.String()}
	c.JSON(200, domain.NewResultOk(data))
}

// 批量生成文件访问地址
func (t *FilesController) generateFileUrlBatch(c *gin.Context) {
	ids := c.Query("ids")
	strings.Split(ids, ",")
}

// 下载文件
func (t *FilesController) downloadFile(c *gin.Context) {
	id := c.Param("id")
	isDownload := c.DefaultQuery("download", "0")
	fileReader, fileSize, fileInfo, err := t.filesService.DownloadFileReader(id)
	if err != nil {
		c.JSON(500, domain.NewResultFail("500", err.Error()))
		return
	}

	if isDownload == "1" {
		attachmentName := fileInfo.Id + path.Ext(fileInfo.Filename)
		headers := map[string]string{
			"Content-Disposition": "attachment; filename=\"" + url.QueryEscape(attachmentName) + "\"",
		}
		c.DataFromReader(200, fileSize, "application/octet-stream", fileReader, headers)
	} else {
		c.DataFromReader(200, fileSize, "", fileReader, nil)
	}
}
