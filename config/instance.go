package config

var config = &Config{}

func GetConfig() *Config {
	return config
}

func SetInstance(cfg *Config) {
	config = cfg
}
