package api

type Config struct {
	APIKey string
}

var config *Config

func init() {
	config = &Config{}
}
