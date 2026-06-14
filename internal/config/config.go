package config

type Config struct {
	Server   ServerConfig
	WSServer WSServerConfig
	DB       DBConfig
	Redis    RedisConfig
	App      App
}

type ServerConfig struct {
	Port string
}

type WSServerConfig struct {
	Port string
}

type DBConfig struct {
	URL      string
	MAXConns int32
	MINConns int32
}

type RedisConfig struct {
	Addr             string
	Password         string
	DB               int
	Stream           string
	DlqStream        string
	Group            string
	RecoveryInterval int
	RecoveryIdleTime int
	Channel          string
}

type App struct {
	Env string
}
