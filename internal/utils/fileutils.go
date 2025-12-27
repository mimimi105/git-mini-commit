package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileUtils ファイル操作の共通ユーティリティ
type FileUtils struct {
	basePath string
}

// NewFileUtils 新しいFileUtilsインスタンスを作成
func NewFileUtils(basePath string) *FileUtils {
	return &FileUtils{
		basePath: basePath,
	}
}

// WriteFile ファイルに内容を書き込む（エラーハンドリング付き）
func (fu *FileUtils) WriteFile(filename, content string) error {
	filePath := filepath.Join(fu.basePath, filename)
	
	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	
	return nil
}

// ReadFile ファイルの内容を読み込む（エラーハンドリング付き）
func (fu *FileUtils) ReadFile(filename string) (string, error) {
	filePath := filepath.Join(fu.basePath, filename)
	
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	
	return string(content), nil
}

// MkdirAll ディレクトリを作成する（エラーハンドリング付き）
func (fu *FileUtils) MkdirAll(dirname string) error {
	dirPath := filepath.Join(fu.basePath, dirname)
	
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}
	
	return nil
}

// Stat ファイル/ディレクトリの情報を取得する（エラーハンドリング付き）
func (fu *FileUtils) Stat(filename string) (os.FileInfo, error) {
	filePath := filepath.Join(fu.basePath, filename)
	
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}
	
	return info, nil
}

// Exists ファイル/ディレクトリが存在するかチェック
func (fu *FileUtils) Exists(filename string) bool {
	filePath := filepath.Join(fu.basePath, filename)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// Remove ファイル/ディレクトリを削除する（エラーハンドリング付き）
func (fu *FileUtils) Remove(filename string) error {
	filePath := filepath.Join(fu.basePath, filename)
	
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove file %s: %w", filePath, err)
	}
	
	return nil
}

// RemoveAll ファイル/ディレクトリを再帰的に削除する（エラーハンドリング付き）
func (fu *FileUtils) RemoveAll(filename string) error {
	filePath := filepath.Join(fu.basePath, filename)
	
	if err := os.RemoveAll(filePath); err != nil {
		return fmt.Errorf("failed to remove all %s: %w", filePath, err)
	}
	
	return nil
}

// ReadDir ディレクトリの内容を読み込む（エラーハンドリング付き）
func (fu *FileUtils) ReadDir(dirname string) ([]os.DirEntry, error) {
	dirPath := filepath.Join(fu.basePath, dirname)
	
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}
	
	return entries, nil
}

// JoinPath パスを結合する
func (fu *FileUtils) JoinPath(elem ...string) string {
	return filepath.Join(fu.basePath, filepath.Join(elem...))
}

// GetBasePath ベースパスを取得
func (fu *FileUtils) GetBasePath() string {
	return fu.basePath
}

// SetBasePath ベースパスを設定
func (fu *FileUtils) SetBasePath(basePath string) {
	fu.basePath = basePath
}

