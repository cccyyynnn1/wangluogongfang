package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"yara-security-service/internal/models"
	"yara-security-service/pkg/security"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScanner 模拟扫描器
type MockScanner struct {
	mock.Mock
}

func (m *MockScanner) ScanFile(ctx context.Context, filePath string) (*models.ScanResult, error) {
	args := m.Called(ctx, filePath)
	return args.Get(0).(*models.ScanResult), args.Error(1)
}

func (m *MockScanner) ScanDirectory(ctx context.Context, dirPath string, recursive bool, maxDepth int) ([]*models.ScanResult, error) {
	args := m.Called(ctx, dirPath, recursive, maxDepth)
	return args.Get(0).([]*models.ScanResult), args.Error(1)
}

func (m *MockScanner) ScanBuffer(ctx context.Context, data []byte, identifier string) (*models.ScanResult, error) {
	args := m.Called(ctx, data, identifier)
	return args.Get(0).(*models.ScanResult), args.Error(1)
}

// 移除这个方法，因为现在getFileInfo在FileHandler中实现

func (m *MockScanner) ReloadRules() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockScanner) GetRulesInfo() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

func setupTestRouter(t *testing.T) (*gin.Engine, *MockScanner) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockScanner := new(MockScanner)
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	// 创建一个真实的Scanner实例用于测试
	realScanner, err := security.NewScanner("./rules", logger)
	if err != nil {
		// 如果创建失败，跳过测试
		t.Skip("无法创建Scanner实例，跳过测试")
	}

	fileHandler := NewFileHandler(nil, realScanner, logger)

	api := router.Group("/api/v1")
	{
		fileGroup := api.Group("/file")
		{
			fileGroup.POST("/scan", fileHandler.ScanFile)
			fileGroup.GET("/info/:path", fileHandler.GetFileInfo)
			fileGroup.POST("/scan-directory", fileHandler.ScanDirectory)
			fileGroup.POST("/scan-buffer", fileHandler.ScanBuffer)
			fileGroup.GET("/list", fileHandler.GetFileList)
		}
	}

	return router, mockScanner
}

func TestScanFile(t *testing.T) {
	router, _ := setupTestRouter(t)

	// 创建测试文件
	testFile := filepath.Join(t.TempDir(), "test.exe")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	assert.NoError(t, err)

	// 创建请求
	req, err := http.NewRequest("POST", "/api/v1/file/scan", bytes.NewBufferString("file_path="+testFile))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "扫描完成", response.Message)
}

func TestScanFile_FileNotExists(t *testing.T) {
	router, _ := setupTestRouter(t)

	// 创建请求
	req, err := http.NewRequest("POST", "/api/v1/file/scan", bytes.NewBufferString("file_path=/nonexistent/file.exe"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, "文件不存在", response.Message)
}

func TestScanFile_EmptyPath(t *testing.T) {
	router, _ := setupTestRouter(t)

	// 创建请求
	req, err := http.NewRequest("POST", "/api/v1/file/scan", bytes.NewBufferString("file_path="))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "文件路径不能为空", response.Message)
}

func TestGetFileInfo(t *testing.T) {
	router, _ := setupTestRouter(t)

	testFile := filepath.Join(t.TempDir(), "test.exe")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	assert.NoError(t, err)

	// 创建请求
	req, err := http.NewRequest("GET", "/api/v1/file/info/"+testFile, nil)
	assert.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "获取文件信息成功", response.Message)
}

func TestScanDirectory(t *testing.T) {
	router, _ := setupTestRouter(t)

	// 创建测试目录
	testDir := t.TempDir()
	err := os.WriteFile(filepath.Join(testDir, "test1.exe"), []byte("test1"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(testDir, "test2.dll"), []byte("test2"), 0644)
	assert.NoError(t, err)

	// 创建请求
	requestBody := models.DirectoryScanRequest{
		Path:      testDir,
		Recursive: true,
		MaxDepth:  5,
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/file/scan-directory", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "目录扫描完成", response.Message)
}

func TestScanDirectory_DirectoryNotExists(t *testing.T) {
	router, _ := setupTestRouter(t)

	// 创建请求
	requestBody := models.DirectoryScanRequest{
		Path:      "/nonexistent/directory",
		Recursive: true,
		MaxDepth:  5,
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/file/scan-directory", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, "目录不存在", response.Message)
}

func TestGetFileList(t *testing.T) {
	router, _ := setupTestRouter(t)

	// 创建测试目录和文件
	testDir := t.TempDir()
	err := os.WriteFile(filepath.Join(testDir, "test1.txt"), []byte("test1"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(testDir, "test2.txt"), []byte("test2"), 0644)
	assert.NoError(t, err)

	// 创建请求
	req, err := http.NewRequest("GET", "/api/v1/file/list?path="+testDir, nil)
	assert.NoError(t, err)

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "获取文件列表成功", response.Message)

	// 验证数据
	data := response.Data.(map[string]interface{})
	assert.Equal(t, testDir, data["directory"])
	assert.Equal(t, float64(2), data["count"])
}

func TestFileExists(t *testing.T) {
	handler := &FileHandler{}

	// 创建测试文件
	testFile := filepath.Join(t.TempDir(), "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	assert.NoError(t, err)

	// 测试文件存在
	assert.True(t, handler.fileExists(testFile))

	// 测试文件不存在
	assert.False(t, handler.fileExists("/nonexistent/file.txt"))
}

func TestDirectoryExists(t *testing.T) {
	handler := &FileHandler{}

	// 创建测试目录
	testDir := t.TempDir()

	// 测试目录存在
	assert.True(t, handler.directoryExists(testDir))

	// 测试目录不存在
	assert.False(t, handler.directoryExists("/nonexistent/directory"))
}
