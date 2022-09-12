package config

import "flag"

type configFlags struct {
	flagPort          *string
	flagHost          *string
	flagFileStorage   *string
	flagDBStorage     *string
	flagHTTPS         *bool
	flagConfigFile    *string
	flagTrustedSubnet *string
	flagGRPCPort      *string
	flagEnableGRPC    *bool
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
	config.flagTrustedSubnet = flag.String("t", "", "TRUSTED_SUBNET")
	config.flagGRPCPort = flag.String("g", "", "GRPC_PORT")
	config.flagEnableGRPC = flag.Bool("enable_gprc", false, "ENABLE_GPRC")

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
	if *config.flagTrustedSubnet != "" {
		c.Security.TrustedSubnet = *config.flagTrustedSubnet
	}
	if *config.flagGRPCPort != "" {
		c.GRPC.Port = *config.flagGRPCPort
	}
	if *config.flagEnableGRPC {
		c.GRPC.EnableGRPC = *config.flagEnableGRPC
	}
}
