package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Security SecurityConfig `mapstructure:"security"`
	File     FileConfig     `mapstructure:"file"`
	Process  ProcessConfig  `mapstructure:"process"`
	Registry RegistryConfig `mapstructure:"registry"`
	Network  NetworkConfig  `mapstructure:"network"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	YaraRulesPath  string        `mapstructure:"yara_rules_path"`
	MaxFileSize    string        `mapstructure:"max_file_size"`
	ScanTimeout    time.Duration `mapstructure:"scan_timeout"`
	EnableRealTime bool          `mapstructure:"enable_real_time"`
	APIKeyRequired bool          `mapstructure:"api_key_required"`
	APIKey         string        `mapstructure:"api_key"`
	AllowedIPs     []string      `mapstructure:"allowed_ips"`
}

// FileConfig 文件配置
type FileConfig struct {
	MaxScanDepth      int      `mapstructure:"max_scan_depth"`
	ExcludePatterns   []string `mapstructure:"exclude_patterns"`
	IncludeExtensions []string `mapstructure:"include_extensions"`
}

// ProcessConfig 进程配置
type ProcessConfig struct {
	MaxProcessCount int           `mapstructure:"max_process_count"`
	RefreshInterval time.Duration `mapstructure:"refresh_interval"`
}

// RegistryConfig 注册表配置
type RegistryConfig struct {
	MaxKeyLength int    `mapstructure:"max_key_length"`
	MaxValueSize string `mapstructure:"max_value_size"`
}

// NetworkConfig 网络配置
type NetworkConfig struct {
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout"`
	MaxConnections    int           `mapstructure:"max_connections"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// 如果配置文件不存在，创建默认配置
	if err := viper.ReadInConfig(); err != nil {
		// 创建默认配置文件
		defaultConfig := GetDefaultConfig()
		if err := createDefaultConfigFile(configPath, defaultConfig); err != nil {
			return nil, fmt.Errorf("创建默认配置文件失败: %w", err)
		}

		// 重新读取配置文件
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证和设置默认值
	if err := validateAndSetDefaults(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// createDefaultConfigFile 创建默认配置文件
func createDefaultConfigFile(configPath string, config *Config) error {
	// 将配置写入文件
	viper.SetConfigFile(configPath)
	viper.Set("server", config.Server)
	viper.Set("logging", config.Logging)
	viper.Set("security", config.Security)
	viper.Set("file", config.File)
	viper.Set("process", config.Process)
	viper.Set("registry", config.Registry)
	viper.Set("network", config.Network)
	viper.Set("cors", config.CORS)

	return viper.WriteConfig()
}

// validateAndSetDefaults 验证配置并设置默认值
func validateAndSetDefaults(config *Config) error {
	// 服务器配置验证
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		config.Server.Port = 8081
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Server.ReadTimeout <= 0 {
		config.Server.ReadTimeout = 30 * time.Second
	}
	if config.Server.WriteTimeout <= 0 {
		config.Server.WriteTimeout = 30 * time.Second
	}

	// 安全配置验证
	if config.Security.YaraRulesPath == "" {
		config.Security.YaraRulesPath = "./rules"
	}
	if config.Security.MaxFileSize == "" {
		config.Security.MaxFileSize = "100MB"
	}
	if config.Security.ScanTimeout <= 0 {
		config.Security.ScanTimeout = 60 * time.Second
	}

	// 日志配置验证
	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "stdout"
	}

	// 文件配置验证
	if config.File.MaxScanDepth <= 0 {
		config.File.MaxScanDepth = 10
	}
	if len(config.File.ExcludePatterns) == 0 {
		config.File.ExcludePatterns = []string{"*.tmp", "*.log", "*.cache"}
	}
	if len(config.File.IncludeExtensions) == 0 {
		config.File.IncludeExtensions = []string{".exe", ".dll", ".sys", ".bat", ".cmd", ".ps1", ".vbs", ".js"}
	}

	// 进程配置验证
	if config.Process.MaxProcessCount <= 0 {
		config.Process.MaxProcessCount = 1000
	}
	if config.Process.RefreshInterval <= 0 {
		config.Process.RefreshInterval = 5 * time.Second
	}

	// 注册表配置验证
	if config.Registry.MaxKeyLength <= 0 {
		config.Registry.MaxKeyLength = 256
	}
	if config.Registry.MaxValueSize == "" {
		config.Registry.MaxValueSize = "1MB"
	}

	// 网络配置验证
	if config.Network.ConnectionTimeout <= 0 {
		config.Network.ConnectionTimeout = 30 * time.Second
	}
	if config.Network.MaxConnections <= 0 {
		config.Network.MaxConnections = 1000
	}

	// CORS配置验证
	if len(config.CORS.AllowedOrigins) == 0 {
		config.CORS.AllowedOrigins = []string{"*"}
	}
	if len(config.CORS.AllowedMethods) == 0 {
		config.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE"}
	}
	if len(config.CORS.AllowedHeaders) == 0 {
		config.CORS.AllowedHeaders = []string{"Content-Type", "Authorization"}
	}

	return nil
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8081,
			Host:         "127.0.0.1",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
		Security: SecurityConfig{
			YaraRulesPath:  "./rules",
			MaxFileSize:    "100MB",
			ScanTimeout:    60 * time.Second,
			EnableRealTime: true,
			APIKeyRequired: false,
			APIKey:         "",
			AllowedIPs:     []string{},
		},
		File: FileConfig{
			MaxScanDepth: 10,
			ExcludePatterns: []string{
				"*.tmp", "*.log", "*.cache",
			},
			IncludeExtensions: []string{
				".exe", ".dll", ".sys", ".bat", ".cmd", ".ps1", ".vbs", ".js",
			},
		},
		Process: ProcessConfig{
			MaxProcessCount: 1000,
			RefreshInterval: 5 * time.Second,
		},
		Registry: RegistryConfig{
			MaxKeyLength: 256,
			MaxValueSize: "1MB",
		},
		Network: NetworkConfig{
			ConnectionTimeout: 30 * time.Second,
			MaxConnections:    1000,
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
	}
}
