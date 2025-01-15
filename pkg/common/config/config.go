package config

type ServerConfig struct{
	ServerIp string `mapstructure:"server_ip"`
	Version string `mapstructure:"version"`
}
