package config

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App      `yaml:"app"`
		Server   `yaml:"http"`
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
		SecretKey   string `env-required:"true" env:"SECRET_KEY" yaml:"secret_key"`
		EnableHTTPS bool   `env:"ENABLE_HTTPS" yaml:"enable_https"`
		CertFile    string `env:"CERT_FILE" yaml:"cert_file"`
		KeyFile     string `env:"KEY_FILE" yaml:"key_file"`
	}

	Workers struct {
		WorkersNumber int `env-required:"true" env:"WORKERS_NUMBER" yaml:"workers_number"`
		PoolLength    int `env-required:"true" env:"POOL_LENGTH" yaml:"pool_length"`
	}
)

const defaultConfigPath = "./config/config.yml"

type configFlags struct {
	flagPort        *string
	flagHost        *string
	flagFileStorage *string
	flagDBStorage   *string
	flagHTTPS       *bool
	flagConfigFile  *string
}

func newConfigFlags() *configFlags {
	return &configFlags{}
}

func (config *configFlags) parseFlags() {
	config.flagPort = flag.String("a", "", "SERVER_ADDRESS")
	config.flagHost = flag.String("b", "", "BASE_URL")
	config.flagFileStorage = flag.String("f", "", "FILE_STORAGE_PATH")
	config.flagDBStorage = flag.String("d", "", "DATABASE_DSN")
	config.flagHTTPS = flag.Bool("s", false, "ENABLE_HTTPS")
	config.flagConfigFile = flag.String("c", defaultConfigPath, "CONFIG")
	configLongFlag := flag.String("config", defaultConfigPath, "CONFIG")
	flag.Parse()
	if *configLongFlag != defaultConfigPath {
		config.flagConfigFile = configLongFlag
	}
}

func (c *Config) checkFlags(config *configFlags) {

	if *config.flagPort != "" {
		c.Server.Port = *config.flagPort
	}
	if *config.flagHost != "" {
		c.Server.Host = *config.flagHost
	}
	if *config.flagDBStorage != "" {
		c.Storage.DataSourceName = *config.flagDBStorage
	}
	if *config.flagFileStorage != "" {
		c.Storage.FilePath = *config.flagFileStorage
	}
	if *config.flagHTTPS {
		c.Security.EnableHTTPS = *config.flagHTTPS
	}
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

	return Cfg
}
