package config

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"time"
)

type JWTConfig struct {
	SigningMethod *jwt.SigningMethodHMAC
	ExpiresIn     time.Duration
	Token         string
}

type Config struct {
	ServiceHost string
	ServicePort int

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOUseSSL    bool
	MinIOBucket    string
	JWT            JWTConfig
	Redis          RedisConfig
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

const (
	envRedisHost = "REDIS_HOST"
	envRedisPort = "REDIS_PORT"
	envRedisUser = "REDIS_USER"
	envRedisPass = "REDIS_PASSWORD"
)

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServiceHost:    viper.GetString("ServiceHost"),
		ServicePort:    viper.GetInt("ServicePort"),
		MinIOEndpoint:  viper.GetString("endpoint"),
		MinIOAccessKey: viper.GetString("access_key"),
		MinIOSecretKey: viper.GetString("secret_key"),
		MinIOUseSSL:    viper.GetBool("use_ssl"),
		MinIOBucket:    viper.GetString("bucket"),
		JWT: JWTConfig{
			SigningMethod: jwt.SigningMethodHS256,
			ExpiresIn:     time.Hour * 24,
			Token:         "your-secret-key",
		},
	} // создаем объект конфига

	cfg.Redis.Host = os.Getenv(envRedisHost)
	cfg.Redis.Port, err = strconv.Atoi(os.Getenv(envRedisPort))

	if err != nil {
		return nil, fmt.Errorf("redis port must be int value: %w", err)
	}

	cfg.Redis.Password = os.Getenv(envRedisPass)
	cfg.Redis.User = os.Getenv(envRedisUser)

	err = viper.Unmarshal(cfg) // читаем информацию из файла,
	// конвертируем и затем кладем в нашу переменную cfg
	if err != nil {
		return nil, err
	}
	log.Infof("Config loaded: Host=%s, Port=%d", cfg.ServiceHost, cfg.ServicePort)
	log.Infof("MinIO Config: endpoint=%s, access_key=%s, bucket=%s",
		cfg.MinIOEndpoint, cfg.MinIOAccessKey, cfg.MinIOBucket)

	log.Info("config parsed")

	return cfg, nil
}

func (c *Config) InitMinIO() (*minio.Client, error) {
	minioClient, err := minio.New(c.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.MinIOAccessKey, c.MinIOSecretKey, ""),
		Secure: c.MinIOUseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Проверяем существование бакета
	exists, err := minioClient.BucketExists(context.Background(), c.MinIOBucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", c.MinIOBucket)
	}

	return minioClient, nil
}
