package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	AppEnv     string
	AppPort    string
	AppName    string
	AppDebug  bool

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	JWTSecret               string
	JWTAccessTokenExpire   time.Duration
	JWTRefreshTokenExpire  time.Duration

	AuthCenterURL        string
	AuthCenterAPIKey      string
	AuthCenterRedirectURI string

	FrontendURL string
}

func Load() *Config {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using defaults: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Failed to unmarshal config:", err)
	}

	return &cfg
}

func setDefaults() {
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_NAME", "PR Business")
	viper.SetDefault("APP_DEBUG", "true")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "pr_business")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "pr_business")
	viper.SetDefault("DB_SSLMODE", "disable")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("JWT_SECRET", "change-me-in-production")
	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRE", "24h")
	viper.SetDefault("JWT_REFRESH_TOKEN_EXPIRE", "168h")

	viper.SetDefault("AUTH_CENTER_URL", "http://os.crazyaigc.com")
	viper.SetDefault("AUTH_CENTER_API_KEY", "")
	viper.SetDefault("AUTH_CENTER_REDIRECT_URI", "http://localhost:8080/api/v1/auth/callback")

	viper.SetDefault("FRONTEND_URL", "http://localhost:5173")
}

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	log.Println("Database connected successfully")
	return db, nil
}

func InitRedis(cfg *Config) (interface{}, error) {
	// TODO: 实现Redis连接
	return nil, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.SSLMode)
}
