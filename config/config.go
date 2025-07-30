package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type environment string

const (
	// EnvLocal represents the local environment.
	EnvLocal environment = "local"

	// EnvTest represents the test environment.
	EnvTest environment = "test"

	// EnvDevelopment represents the development environment.
	EnvDevelopment environment = "dev"

	// EnvStaging represents the staging environment.
	EnvStaging environment = "staging"

	// EnvQA represents the qa environment.
	EnvQA environment = "qa"

	// EnvProduction represents the production environment.
	EnvProduction environment = "prod"
)

// SwitchEnvironment sets the environment variable used to dictate which environment the application is
// currently running in.
// This must be called prior to loading the configuration in order for it to take effect.
func SwitchEnvironment(env environment) {
	if err := os.Setenv("ZERO_APP_ENVIRONMENT", string(env)); err != nil {
		panic(err)
	}
}

type (
	// Config stores complete configuration.
	Config struct {
		HTTP       HTTPConfig
		App        AppConfig
		Cache      CacheConfig
		Database   DatabaseConfig
		Files      FilesConfig
		Tasks      TasksConfig
		Mail       MailConfig
		Storage    StorageConfig
		WhatsApp   WhatsAppConfig
		AI         AIConfig
		FileUpload FileUploadConfig
		Security   SecurityConfig
		Monitoring MonitoringConfig
	}

	// HTTPConfig stores HTTP configuration.
	HTTPConfig struct {
		Hostname        string
		Port            uint16
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		ShutdownTimeout time.Duration
		TLS             struct {
			Enabled     bool
			Certificate string
			Key         string
		}
	}

	// AppConfig stores application configuration.
	AppConfig struct {
		Name          string
		Host          string
		Environment   environment
		EncryptionKey string
		Timeout       time.Duration
		PasswordToken struct {
			Expiration time.Duration
			Length     int
		}
		EmailVerificationTokenExpiration time.Duration
	}

	// CacheConfig stores the cache configuration.
	CacheConfig struct {
		Capacity   int
		Expiration struct {
			PublicFile  time.Duration
			UserSession time.Duration
			NoteData    time.Duration
			UserProfile time.Duration
			NotesList   time.Duration
			NoteCounts  time.Duration
			PublicNotes time.Duration
		}
	}

	// DatabaseConfig stores the database configuration.
	DatabaseConfig struct {
		Driver         string
		Connection     string
		TestConnection string
	}

	// FilesConfig stores the file system configuration.
	FilesConfig struct {
		Directory string
	}

	// TasksConfig stores the tasks configuration.
	TasksConfig struct {
		Goroutines      int
		ReleaseAfter    time.Duration
		CleanupInterval time.Duration
		ShutdownTimeout time.Duration
	}

	// MailConfig stores the mail configuration.
	MailConfig struct {
		Hostname    string
		Port        uint16
		User        string
		Password    string
		FromAddress string
	}

	// StorageConfig stores the cloud storage configuration.
	StorageConfig struct {
		AWS   AWSStorageConfig
		GCS   GCSStorageConfig
		Azure AzureStorageConfig
		Local LocalStorageConfig
	}

	// AWSStorageConfig stores AWS S3 configuration.
	AWSStorageConfig struct {
		Enabled bool
		Bucket  string
		Region  string
	}

	// GCSStorageConfig stores Google Cloud Storage configuration.
	GCSStorageConfig struct {
		Enabled   bool
		Bucket    string
		ProjectID string `mapstructure:"projectId"`
	}

	// AzureStorageConfig stores Azure Blob Storage configuration.
	AzureStorageConfig struct {
		Enabled   bool
		Account   string
		Key       string
		Container string
	}

	// LocalStorageConfig stores local storage configuration.
	LocalStorageConfig struct {
		Enabled bool
		BaseURL string `mapstructure:"baseUrl"`
	}

	// WhatsAppConfig stores WhatsApp integration configuration.
	WhatsAppConfig struct {
		Enabled            bool
		PhoneNumberID      string `mapstructure:"phoneNumberId"`
		AccessToken        string `mapstructure:"accessToken"`
		WebhookVerifyToken string `mapstructure:"webhookVerifyToken"`
		AppSecret          string `mapstructure:"appSecret"`
		APIVersion         string `mapstructure:"apiVersion"`
		BaseURL            string `mapstructure:"baseUrl"`
	}

	// AIConfig stores AI/LLM configuration.
	AIConfig struct {
		Enabled    bool
		OpenAI     OpenAIConfig
		Anthropic  AnthropicConfig
	}

	// OpenAIConfig stores OpenAI configuration.
	OpenAIConfig struct {
		Enabled   bool
		APIKey    string `mapstructure:"apiKey"`
		Model     string
		MaxTokens int    `mapstructure:"maxTokens"`
	}

	// AnthropicConfig stores Anthropic Claude configuration.
	AnthropicConfig struct {
		Enabled   bool
		APIKey    string `mapstructure:"apiKey"`
		Model     string
		MaxTokens int    `mapstructure:"maxTokens"`
	}

	// FileUploadConfig stores file upload limits and settings.
	FileUploadConfig struct {
		MaxFileSize  string   `mapstructure:"maxFileSize"`
		MaxTotalSize string   `mapstructure:"maxTotalSize"`
		MaxFiles     int      `mapstructure:"maxFiles"`
		AllowedTypes []string `mapstructure:"allowedTypes"`
	}

	// SecurityConfig stores security-related configuration.
	SecurityConfig struct {
		CORS      CORSConfig
		RateLimit RateLimitConfig
		CSP       CSPConfig
	}

	// CORSConfig stores CORS configuration.
	CORSConfig struct {
		Enabled        bool
		AllowedOrigins []string `mapstructure:"allowedOrigins"`
		AllowedMethods []string `mapstructure:"allowedMethods"`
		AllowedHeaders []string `mapstructure:"allowedHeaders"`
	}

	// RateLimitConfig stores rate limiting configuration.
	RateLimitConfig struct {
		Enabled           bool
		RequestsPerMinute int `mapstructure:"requestsPerMinute"`
		BurstSize         int `mapstructure:"burstSize"`
	}

	// CSPConfig stores Content Security Policy configuration.
	CSPConfig struct {
		Enabled    bool
		Directives CSPDirectives
	}

	// CSPDirectives stores CSP directive values.
	CSPDirectives struct {
		DefaultSrc string `mapstructure:"defaultSrc"`
		ScriptSrc  string `mapstructure:"scriptSrc"`
		StyleSrc   string `mapstructure:"styleSrc"`
		ImgSrc     string `mapstructure:"imgSrc"`
	}

	// MonitoringConfig stores monitoring and logging configuration.
	MonitoringConfig struct {
		Metrics        MetricsConfig
		Health         HealthConfig
		RequestLogging RequestLoggingConfig
	}

	// MetricsConfig stores metrics collection configuration.
	MetricsConfig struct {
		Enabled  bool
		Endpoint string
	}

	// HealthConfig stores health check configuration.
	HealthConfig struct {
		Enabled  bool
		Endpoint string
	}

	// RequestLoggingConfig stores request logging configuration.
	RequestLoggingConfig struct {
		Enabled      bool
		ExcludePaths []string `mapstructure:"excludePaths"`
	}
)

// GetConfig loads and returns configuration.
func GetConfig() (Config, error) {
	var c Config

	// Load the config file.
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// Load env variables.
	viper.SetEnvPrefix("zero")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}
