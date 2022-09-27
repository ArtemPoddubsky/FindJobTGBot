package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Token  string `toml:"token"`
	Host   string `toml:"host"`
	Server struct {
		Port string `toml:"port"`
	} `toml:"server"`
	DB struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Database string `toml:"database"`
	} `toml:"database"`
}

func GetConfig() Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalln(err)
	}

	return cfg
}
