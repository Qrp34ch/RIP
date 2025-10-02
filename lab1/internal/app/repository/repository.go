package repository

import (
	"github.com/minio/minio-go/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db          *gorm.DB
	minioClient *minio.Client
	bucketName  string
}

func New(dsn string, minioClient *minio.Client, bucketName string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db:          db,
		minioClient: minioClient,
		bucketName:  bucketName,
	}, nil
}
