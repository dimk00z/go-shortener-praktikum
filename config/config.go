package config

import (
	"log"
	"net"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App      `yaml:"app"`
		Server   `yaml:"http"`
		GRPC     `yaml:"grpc"`
		Log      `yaml:"logger"`
		Storage  `yaml:"storage"`
		Security `yaml:"security"`
		Workers  `yaml:"workers"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	Server struct {
		Port string `env-required:"true" yaml:"server_address" env:"SERVER_ADDRESS"`
		Host string `env-required:"true" yaml:"base_url" env:"BASE_URL"`
	}
	GRPC struct {
		Port       string `yaml:"port" env:"GRPC_PORT"`
		EnableGRPC bool   `env:"ENABLE_GRPC" yaml:"enable_gprc"`
	}
	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	Storage struct {
		FilePath       string `yaml:"file_storage_path" env:"FILE_STORAGE_PATH"`
		PoolMax        int    `yaml:"pool_max" env:"PG_POOL_MAX"`
		DataSourceName string `yaml:"dsn" env:"DATABASE_DSN"`
		MaxRetries     int    `yaml:"max_retries" env:"MAX_RETRIES"`
	}

	Security struct {
		SecretKey     string `env-required:"true" env:"SECRET_KEY" yaml:"secret_key"`
		EnableHTTPS   bool   `env:"ENABLE_HTTPS" yaml:"enable_https"`
		CertFile      string `env:"CERT_FILE" yaml:"cert_file"`
		KeyFile       string `env:"KEY_FILE" yaml:"key_file"`
		TrustedSubnet string `env:"TRUSTED_SUBNET" yaml:"trusted_subnet"`
	}

	Workers struct {
		WorkersNumber int `env-required:"true" env:"WORKERS_NUMBER" yaml:"workers_number"`
		PoolLength    int `env-required:"true" env:"POOL_LENGTH" yaml:"pool_length"`
	}
)

const defaultConfigPath = "./config/config.yml"

func checkCIDR(s string) error {
	_, _, err := net.ParseCIDR(s)
	if err != nil {
		log.Printf("CIDR parsing error: %v", err)
	}
	return err
}

// NewConfig returns app config.
func LoadConfig() *Config {
	var err error

	parsedFlags := newConfigFlags()
	parsedFlags.parseFlags()

	Cfg := &Config{}
	configPath := parsedFlags.flagConfigFile
	err = cleanenv.ReadConfig(*configPath, Cfg)
	if err != nil {
		log.Fatalf("config error: %v", err)
		return nil
	}

	err = cleanenv.ReadEnv(Cfg)
	if err != nil {
		log.Fatalf("config error: %v", err)
		return nil
	}
	Cfg.checkFlags(parsedFlags)

	if checkCIDR(Cfg.Security.TrustedSubnet) != nil {
		Cfg.Security.TrustedSubnet = ""
	}
	return Cfg
}
