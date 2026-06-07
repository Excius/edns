package config

type Config struct {
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
	App    App
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	URL      string
	MAXConns int32
	MINConns int32
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Stream   string
}

type App struct {
	Env string
}
