package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() *Config {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
		},
		DB: DBConfig{
			URL:      viper.GetString("db.url"),
			MAXConns: viper.GetInt32("db.max_conns"),
			MINConns: viper.GetInt32("db.min_conns"),
		},
		Redis: RedisConfig{
			Addr: viper.GetString("redis.addr"),
		},
		App: App{
			Env: viper.GetString("app.env"),
		},
	}

	return config
}
