package config

type Config struct {
	Server       ServerConfig
	WorkerServer WorkerConfig
	WSServer     WSServerConfig
	DB           DBConfig
	Redis        RedisConfig
	SMTP         SMTPConfig
	App          App
}

type ServerConfig struct {
	Port string
}

type WorkerConfig struct {
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
	Addr                 string
	Password             string
	DB                   int
	Stream               string
	DlqStream            string
	Group                string
	RecoveryMessageCount int64
	RecoveryInterval     int
	RecoveryIdleTime     int
	Channel              string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type App struct {
	Env string
}
