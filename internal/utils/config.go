package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config 設定管理
type Config struct {
	// ストレージ設定
	Storage StorageConfig `json:"storage"`
	
	// ログ設定
	Logging LoggingConfig `json:"logging"`
	
	// Git設定
	Git GitConfig `json:"git"`
	
	// パフォーマンス設定
	Performance PerformanceConfig `json:"performance"`
}

// StorageConfig ストレージ設定
type StorageConfig struct {
	BasePath     string `json:"base_path"`
	MaxSize      int64  `json:"max_size"`
	Compression  bool   `json:"compression"`
	Encryption   bool   `json:"encryption"`
	BackupCount  int    `json:"backup_count"`
}

// LoggingConfig ログ設定
type LoggingConfig struct {
	Level      string `json:"level"`
	Output     string `json:"output"`
	MaxSize    int64  `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// GitConfig Git設定
type GitConfig struct {
	DefaultUser  string `json:"default_user"`
	DefaultEmail string `json:"default_email"`
	Timeout      int    `json:"timeout"`
	RetryCount   int    `json:"retry_count"`
}

// PerformanceConfig パフォーマンス設定
type PerformanceConfig struct {
	MaxConcurrentOperations int `json:"max_concurrent_operations"`
	CacheSize               int `json:"cache_size"`
	Timeout                 int `json:"timeout"`
}

// DefaultConfig デフォルト設定を取得
func DefaultConfig() *Config {
	return &Config{
		Storage: StorageConfig{
			BasePath:     ".git/mini-commits",
			MaxSize:      100 * 1024 * 1024, // 100MB
			Compression:  true,
			Encryption:   false,
			BackupCount:  5,
		},
		Logging: LoggingConfig{
			Level:      "INFO",
			Output:     "stdout",
			MaxSize:    10 * 1024 * 1024, // 10MB
			MaxBackups: 3,
			MaxAge:     7, // days
			Compress:   true,
		},
		Git: GitConfig{
			DefaultUser:  "git-mini-commit",
			DefaultEmail: "git-mini-commit@localhost",
			Timeout:      30, // seconds
			RetryCount:   3,
		},
		Performance: PerformanceConfig{
			MaxConcurrentOperations: 10,
			CacheSize:               1000,
			Timeout:                 30, // seconds
		},
	}
}

// ConfigManager 設定管理機能
type ConfigManager struct {
	config     *Config
	configPath string
	logger     *Logger
}

// NewConfigManager 新しいConfigManagerインスタンスを作成
func NewConfigManager(configPath string, logger *Logger) *ConfigManager {
	return &ConfigManager{
		config:     DefaultConfig(),
		configPath: configPath,
		logger:     logger,
	}
}

// LoadConfig 設定ファイルから設定を読み込む
func (cm *ConfigManager) LoadConfig() error {
	if cm.configPath == "" {
		cm.logger.Debug("No config path specified, using default config")
		return nil
	}
	
	// 設定ファイルが存在しない場合はデフォルト設定を作成
	if !cm.fileExists(cm.configPath) {
		cm.logger.Info("Config file not found, creating default config: %s", cm.configPath)
		return cm.SaveConfig()
	}
	
	cm.logger.Debug("Loading config from: %s", cm.configPath)
	
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	cm.config = &config
	cm.logger.Info("Config loaded successfully")
	
	return nil
}

// SaveConfig 設定をファイルに保存
func (cm *ConfigManager) SaveConfig() error {
	if cm.configPath == "" {
		cm.logger.Debug("No config path specified, skipping save")
		return nil
	}
	
	cm.logger.Debug("Saving config to: %s", cm.configPath)
	
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	cm.logger.Info("Config saved successfully")
	return nil
}

// GetConfig 設定を取得
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// SetConfig 設定を設定
func (cm *ConfigManager) SetConfig(config *Config) {
	cm.config = config
}

// GetStorageConfig ストレージ設定を取得
func (cm *ConfigManager) GetStorageConfig() StorageConfig {
	return cm.config.Storage
}

// GetLoggingConfig ログ設定を取得
func (cm *ConfigManager) GetLoggingConfig() LoggingConfig {
	return cm.config.Logging
}

// GetGitConfig Git設定を取得
func (cm *ConfigManager) GetGitConfig() GitConfig {
	return cm.config.Git
}

// GetPerformanceConfig パフォーマンス設定を取得
func (cm *ConfigManager) GetPerformanceConfig() PerformanceConfig {
	return cm.config.Performance
}

// LoadFromEnv 環境変数から設定を読み込む
func (cm *ConfigManager) LoadFromEnv() {
	cm.logger.Debug("Loading config from environment variables")
	
	// ストレージ設定
	if basePath := os.Getenv("GIT_MINI_COMMIT_STORAGE_PATH"); basePath != "" {
		cm.config.Storage.BasePath = basePath
	}
	if maxSize := os.Getenv("GIT_MINI_COMMIT_MAX_SIZE"); maxSize != "" {
		if size, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
			cm.config.Storage.MaxSize = size
		}
	}
	if compression := os.Getenv("GIT_MINI_COMMIT_COMPRESSION"); compression != "" {
		cm.config.Storage.Compression = strings.ToLower(compression) == "true"
	}
	
	// ログ設定
	if level := os.Getenv("GIT_MINI_COMMIT_LOG_LEVEL"); level != "" {
		cm.config.Logging.Level = strings.ToUpper(level)
	}
	if output := os.Getenv("GIT_MINI_COMMIT_LOG_OUTPUT"); output != "" {
		cm.config.Logging.Output = output
	}
	
	// Git設定
	if user := os.Getenv("GIT_MINI_COMMIT_DEFAULT_USER"); user != "" {
		cm.config.Git.DefaultUser = user
	}
	if email := os.Getenv("GIT_MINI_COMMIT_DEFAULT_EMAIL"); email != "" {
		cm.config.Git.DefaultEmail = email
	}
	if timeout := os.Getenv("GIT_MINI_COMMIT_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			cm.config.Git.Timeout = t
		}
	}
	
	// パフォーマンス設定
	if maxConcurrent := os.Getenv("GIT_MINI_COMMIT_MAX_CONCURRENT"); maxConcurrent != "" {
		if mc, err := strconv.Atoi(maxConcurrent); err == nil {
			cm.config.Performance.MaxConcurrentOperations = mc
		}
	}
	if cacheSize := os.Getenv("GIT_MINI_COMMIT_CACHE_SIZE"); cacheSize != "" {
		if cs, err := strconv.Atoi(cacheSize); err == nil {
			cm.config.Performance.CacheSize = cs
		}
	}
	
	cm.logger.Info("Config loaded from environment variables")
}

// ValidateConfig 設定の妥当性をチェック
func (cm *ConfigManager) ValidateConfig() error {
	// ストレージ設定の検証
	if cm.config.Storage.BasePath == "" {
		return fmt.Errorf("storage base path cannot be empty")
	}
	if cm.config.Storage.MaxSize <= 0 {
		return fmt.Errorf("storage max size must be positive")
	}
	if cm.config.Storage.BackupCount < 0 {
		return fmt.Errorf("storage backup count cannot be negative")
	}
	
	// ログ設定の検証
	validLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	if !cm.contains(validLevels, cm.config.Logging.Level) {
		return fmt.Errorf("invalid log level: %s", cm.config.Logging.Level)
	}
	if cm.config.Logging.MaxSize <= 0 {
		return fmt.Errorf("log max size must be positive")
	}
	if cm.config.Logging.MaxBackups < 0 {
		return fmt.Errorf("log max backups cannot be negative")
	}
	if cm.config.Logging.MaxAge < 0 {
		return fmt.Errorf("log max age cannot be negative")
	}
	
	// Git設定の検証
	if cm.config.Git.DefaultUser == "" {
		return fmt.Errorf("git default user cannot be empty")
	}
	if cm.config.Git.DefaultEmail == "" {
		return fmt.Errorf("git default email cannot be empty")
	}
	if cm.config.Git.Timeout <= 0 {
		return fmt.Errorf("git timeout must be positive")
	}
	if cm.config.Git.RetryCount < 0 {
		return fmt.Errorf("git retry count cannot be negative")
	}
	
	// パフォーマンス設定の検証
	if cm.config.Performance.MaxConcurrentOperations <= 0 {
		return fmt.Errorf("max concurrent operations must be positive")
	}
	if cm.config.Performance.CacheSize < 0 {
		return fmt.Errorf("cache size cannot be negative")
	}
	if cm.config.Performance.Timeout <= 0 {
		return fmt.Errorf("performance timeout must be positive")
	}
	
	return nil
}

// fileExists ファイルが存在するかチェック
func (cm *ConfigManager) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// contains スライスに要素が含まれているかチェック
func (cm *ConfigManager) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

