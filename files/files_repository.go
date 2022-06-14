package files

import (
	"files/domain"
	"github.com/jmoiron/sqlx"
)

type filesRepository struct {
	db *sqlx.DB
}

func NewFilesRepository(db *sqlx.DB) *filesRepository {
	return &filesRepository{
		db: db,
	}
}

func (t filesRepository) Save(fileInfo domain.FileInfo) error {
	sql := "INSERT INTO file_info(id, `key`, filename, file_size, created_at) VALUES(?, ?, ?, ?, ?)"
	_, err := t.db.Exec(sql, fileInfo.Id, fileInfo.Key, fileInfo.Filename, fileInfo.FileSize, fileInfo.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (t filesRepository) FindById(id string) (*domain.FileInfo, error) {
	short := new(domain.FileInfo)
	err := t.db.Get(short, "select id, `key`, filename, file_size from file_info where id = ?", id)
	if err != nil {
		return nil, err
	}
	return short, nil
}
