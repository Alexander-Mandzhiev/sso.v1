package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"development" env:"ENV" env-default:"development" env-default:"production"`
	DBUpp         string `yaml:"db_upp" env:"DB_UPP"`
	DBAgroReports string `yaml:"db_agro_reports" env:"DB_AGRO_REPORTS"`
	HTTPServer    `yaml:"http_server"`
	SSO           `yaml:"sso"`
	DBConfig      `yaml:"db_config"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env:"ADDRESS"`
	AddressFrontend string        `yaml:"address_frontend" env:"ADDRESS_FRONTEND"`
	Port            int           `yaml:"port" env:"PORT"`
	PortFrontend    int           `yaml:"port_frontend" env:"PORT"`
	Timeout         time.Duration `yaml:"timeout" env:"TIMEOUT"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT"`
}

type SSO struct {
	SecretKeyAccessToken  string        `yaml:"secret_key_access_token" env:"SECRET_KEY_ACCESS_TOKEN" env-default:"secret access token ttl"`
	SecretKeyRefreshToken string        `yaml:"secret_key_refresh_token" env:"SECRET_KEY_REFRESH_TOKEN" env-default:"secret refresh token ttl"`
	AccessTokenTTL        time.Duration `yaml:"access_token_ttl" env:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL       time.Duration `yaml:"refresh_token_ttl" env:"REFRESH_TOKEN_TTL"`
}

type DBConfig struct {
	MaxOpenConnections    int           `yaml:"max_open_connections" env:"MAX_OPEN_CONNECTIONS" env-default:"10"`
	MaxIdleConnections    int           `json:"max_idle_connections" env:"MAX_IDLE_CONNECTIONS" env-default:"5"`
	ConnectionMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"CONN_MAX_LIFETIME" env-default:"5m"`
}

var Cfg *Config

func Initialize() {
	Cfg = MustLoad()
}

func MustLoad() *Config {
	configPath := fetchConfigFlag()
	if configPath == "" {
		return loadingDataInEnv()
	}
	return MustLoadByPath(configPath)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigFlag() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}

func loadingDataInEnv() *Config {
	loadEnv()
	return &Config{
		Env:           os.Getenv("ENV"),
		DBUpp:         os.Getenv("DB_UPP"),
		DBAgroReports: os.Getenv("DB_AGRO_REPORTS"),
		HTTPServer: HTTPServer{
			Address:     os.Getenv("HTTP_ADDRESS"),
			Timeout:     4 * time.Second,
			IdleTimeout: 60 * time.Second,
		},
		SSO: SSO{
			SecretKeyAccessToken:  os.Getenv("SECRET_KEY_ACCESS_TOKEN"),
			SecretKeyRefreshToken: os.Getenv("SECRET_KEY_REFRESH_TOKEN"),
			AccessTokenTTL:        5 * time.Minute,
			RefreshTokenTTL:       10 * time.Hour,
		},
	}
}

func loadEnv() {
	if err := godotenv.Load(".sso.env"); err != nil {
		log.Println("Warning: .sso.env file not found, using default values.")
	}
}
