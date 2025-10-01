package config

import (
	"log"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Type string
	Path string
}

type JWTConfig struct {
	Secret     string
	ExpireTime int // 过期时间（小时）
}

type UploadConfig struct {
	MaxSize      int64 // 最大文件大小（字节）
	AllowedTypes []string
	StoragePath  string
}

var AppConfig *Config

func InitConfig() {
	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "release"),
		},
		Database: DatabaseConfig{
			Type: "sqlite",
			Path: "./gotux.db",
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpireTime: 24 * 7, // 7天
		},
		Upload: UploadConfig{
			MaxSize:      10 * 1024 * 1024, // 10MB
			AllowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
			StoragePath:  "./uploads",
		},
	}

	// 确保上传目录存在
	if err := os.MkdirAll(AppConfig.Upload.StoragePath, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
