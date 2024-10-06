package config

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Cache struct {
		Server string `yaml:"server"`
		Prefix string `yaml:"prefix"`
		Limit  int    `yaml:"limit"`
		Ttl    int    `yaml:"ttl"`
	} `yaml:"cache"`
}
