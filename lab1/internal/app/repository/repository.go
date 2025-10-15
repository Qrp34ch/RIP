package repository

import (
	"github.com/minio/minio-go/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lab1/internal/app/redis"
)

type Repository struct {
	db          *gorm.DB
	minioClient *minio.Client
	bucketName  string
	RedisClient *redis.Client
}

func New(dsn string, minioClient *minio.Client, bucketName string, redisClient *redis.Client) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db:          db,
		minioClient: minioClient,
		bucketName:  bucketName,
		RedisClient: redisClient,
	}, nil
}
