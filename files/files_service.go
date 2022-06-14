package files

import (
	"context"
	"errors"
	"files/config"
	"files/domain"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mime/multipart"
	"net/url"
	"path"
	"time"
)

type filesService struct {
	filesRepo   domain.FilesRepository
	minioConfig *config.MinioConfig
}

func NewFilesService(filesRepo domain.FilesRepository, minioConfig *config.MinioConfig) *filesService {
	return &filesService{filesRepo: filesRepo, minioConfig: minioConfig}
}

func (t filesService) GenerateUploadUrl(request *domain.GenerateUploadUrlRequest) (*domain.GenerateUploadUrlResponse, error) {
	minioConfig := t.minioConfig
	client, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	// 生成签名
	postPolicy := minio.NewPostPolicy()
	postPolicy.SetBucket(minioConfig.Bucket)
	postPolicy.SetKey(request.Key)
	postPolicy.SetExpires(time.Now().Add(10 * time.Minute))
	url, formData, err := client.PresignedPostPolicy(context.Background(), postPolicy)
	if err != nil {
		return nil, err
	}

	// 响应
	response := &domain.GenerateUploadUrlResponse{
		Id:       uuid.New().String(),
		Url:      url.String(),
		FormData: formData,
	}
	return response, nil
}

func (t filesService) UploadFile(fileHeader *multipart.FileHeader) (*domain.UploadFileResponse, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	minioConfig := t.minioConfig
	client, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	objectKey := time.Now().Format("2006/01/02/") + uuid.New().String() + path.Ext(fileHeader.Filename)
	uploadInfo, err := client.PutObject(context.Background(), minioConfig.Bucket, objectKey, file, fileHeader.Size, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	// 存储记录
	fileInfo := domain.FileInfo{
		Id:        uuid.New().String(),
		Key:       objectKey,
		Filename:  fileHeader.Filename,
		FileSize:  fileHeader.Size,
		CreatedAt: time.Now(),
	}
	if err = t.filesRepo.Save(fileInfo); err != nil {
		return nil, err
	}

	response := &domain.UploadFileResponse{
		Id:       fileInfo.Id,
		Filename: fileInfo.Filename,
		Url:      uploadInfo.Location,
	}
	return response, nil
}

func (t filesService) FindById(id string) (*domain.FileInfo, error) {
	return t.filesRepo.FindById(id)
}

func (t filesService) GenerateFileUrl(id string) (*domain.FileInfo, *url.URL, error) {
	fileInfo, err := t.FindById(id)
	if err != nil {
		return nil, nil, err
	}
	if fileInfo == nil {
		return nil, nil, nil
	}

	minioConfig := t.minioConfig
	client, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, nil, err
	}

	expires := 1 * time.Minute
	signedUrl, err := client.PresignedGetObject(context.Background(), minioConfig.Bucket, fileInfo.Key, expires, url.Values{})
	if err != nil {
		return fileInfo, nil, err
	}
	return fileInfo, signedUrl, nil
}

func (t filesService) DownloadFileReader(id string) (io.Reader, int64, *domain.FileInfo, error) {
	fileInfo, err := t.FindById(id)
	if err != nil {
		return nil, 0, nil, err
	}
	if fileInfo == nil {
		return nil, 0, nil, errors.New("文件id不存在")
	}
	minioConfig := t.minioConfig
	client, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, 0, fileInfo, err
	}

	object, err := client.GetObject(context.Background(), minioConfig.Bucket, fileInfo.Key, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, fileInfo, err
	}
	stat, err := object.Stat()
	if err != nil {
		return nil, 0, fileInfo, err
	}
	return object, stat.Size, fileInfo, nil
}
