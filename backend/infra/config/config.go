package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	// Host 监听的主机地址（如 0.0.0.0 或 127.0.0.1），留空则使用 Port 字段原样作为地址（兼容老配置）
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	EnableSwagger bool   `yaml:"enable_swagger"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey string `yaml:"secret_key"`
	ExpiresIn int    `yaml:"expires_in"` // 小时
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `yaml:"level"`
}

// Load 加载配置
func Load() *Config {
	cfg := &Config{}

	// 尝试从配置文件加载
	if data, err := os.ReadFile("config.yaml"); err == nil {
		if yaml.Unmarshal(data, cfg) == nil {
			applyDefaults(cfg)
			return cfg
		}
	}

	// 如果文件不存在或解析失败，使用环境变量和默认值
	cfg = &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			// 兼容：允许 SERVER_PORT 传入形如 ":8090" 或 "8090" 或 "127.0.0.1:8090"
			Port:          getEnv("SERVER_PORT", ":8090"),
			EnableSwagger: getEnv("ENABLE_SWAGGER", "true") == "true",
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "alice"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET", "alice-secret-key"),
			ExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 24),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
	applyDefaults(cfg)
	return cfg
}

// applyDefaults 填充缺失的默认值 (尤其是 YAML 未显式配置的布尔字段)
func applyDefaults(c *Config) {
	if c.Server.Port == "" {
		c.Server.Port = ":8090"
	}
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	// 若未在 YAML 中声明且未通过 env 设置, 默认开启 swagger
	if !c.Server.EnableSwagger && getEnv("ENABLE_SWAGGER", "") == "" {
		c.Server.EnableSwagger = true
	}
	if c.JWT.ExpiresIn == 0 {
		c.JWT.ExpiresIn = 24
	}
	if c.JWT.SecretKey == "" {
		c.JWT.SecretKey = "alice-secret-key"
	}
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为int，如果不存在或转换失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
