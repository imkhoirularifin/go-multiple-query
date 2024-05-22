package config

type Config struct {
	Host          string   `env:"HOST" envDefault:"localhost"`
	Port          int      `env:"PORT" envDefault:"8080"`
	ProxyHeader   string   `env:"PROXY_HEADER" envDefault:"X-Forwarded-For"`
	LogFields     []string `env:"LOG_FIELDS" envSeparator:","`
	IsDevelopment bool     `env:"IS_DEVELOPMENT" envDefault:"true"`
	MongoDb       MongoDb
}

type MongoDb struct {
	URI string `env:"MONGODB_URI,notEmpty"`
}
