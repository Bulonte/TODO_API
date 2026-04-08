package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Storage 文件存储接口
type Storage interface {
	// Upload 上传文件
	Upload(ctx context.Context, file io.Reader, filename string) (string, error)
	// Download 下载文件
	Download(ctx context.Context, key string) (io.Reader, error)
	// Delete 删除文件
	Delete(ctx context.Context, key string) error
	// GetURL 获取文件URL
	GetURL(ctx context.Context, key string) (string, error)
}

// LocalStorage 本地文件存储实现
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage 创建本地文件存储实例
func NewLocalStorage(basePath, baseURL string) Storage {
	// 确保基础路径存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		panic(err)
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

// getFilePath 获取文件路径
func (s *LocalStorage) getFilePath(key string) string {
	// 移除可能的前缀
	key = strings.TrimPrefix(key, "/")
	return filepath.Join(s.basePath, key)
}

// Upload 上传文件
func (s *LocalStorage) Upload(ctx context.Context, file io.Reader, filename string) (string, error) {
	// 生成唯一文件名
	uniqueFilename := filepath.Join(time.Now().Format("2006/01/02"), strings.ReplaceAll(filename, " ", "_"))
	filePath := s.getFilePath(uniqueFilename)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", err
	}

	// 创建文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	// 返回相对路径
	return uniqueFilename, nil
}

// Download 下载文件
func (s *LocalStorage) Download(ctx context.Context, key string) (io.Reader, error) {
	filePath := s.getFilePath(key)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Delete 删除文件
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	filePath := s.getFilePath(key)

	// 删除文件
	return os.Remove(filePath)
}

// GetURL 获取文件URL
func (s *LocalStorage) GetURL(ctx context.Context, key string) (string, error) {
	// 生成文件URL
	return s.baseURL + "/" + key, nil
}
