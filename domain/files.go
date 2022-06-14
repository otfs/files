package domain

import (
	"io"
	"mime/multipart"
	"net/url"
	"time"
)

type FileInfo struct {
	Id        string    `db:"id"`         // 文件ID
	Key       string    `db:"key"`        // 文件对象key
	Filename  string    `db:"filename"`   // 原始文件名
	FileSize  int64     `db:"file_size"`  // 文件大小
	CreatedAt time.Time `db:"created_at"` // 创建时间
}

type FilesService interface {
	GenerateUploadUrl(request *GenerateUploadUrlRequest) (*GenerateUploadUrlResponse, error)
	UploadFile(header *multipart.FileHeader) (*UploadFileResponse, error)
	FindById(id string) (*FileInfo, error)
	GenerateFileUrl(id string) (*FileInfo, *url.URL, error)
	DownloadFileReader(id string) (io.Reader, int64, *FileInfo, error)
}

type FilesRepository interface {
	Save(fileInfo FileInfo) error
	FindById(id string) (*FileInfo, error)
}

//-------------------------------------------------------------
// 请求响应参数
//-------------------------------------------------------------

type GenerateUploadUrlRequest struct {
	Key string `json:"key" form:"key" binding:"required"` // 文件对象
}

type GenerateUploadUrlResponse struct {
	Id       string            `json:"id"`       // 文件ID
	Url      string            `json:"url"`      // 文件访问地址(默认时效1小时)
	FormData map[string]string `json:"formData"` // 表单数据
}

type UploadFileResponse struct {
	Id       string `json:"id"`       // 文件ID
	Filename string `json:"filename"` // 文件名
	Url      string `json:"url"`      // 文件访问地址
}
