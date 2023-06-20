package config

import (
	"flag"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	server struct {
		Port           string        `env:"PORT" env-default:":8080"`
		AllowedOrigins []string      `env:"ALLOWED_CORS_ORIGINS" env-default:"*"`
		TimeoutRead    time.Duration `env:"SERVER_READ_TIMEOUT" env-default:"15s"`
		TimeoutWrite   time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"15s"`
	}

	flags struct {
		envFilename string
		DevMode     bool
	}

	Config struct {
		DatabaseURL   string `env:"DATABASE_URL"`
		RedisURL      string `env:"REDIS_URL"`
		RedisPassword string `env:"REDIS_PASSWORD"`
		Server        server
		JeagerURL     string `env:"JAEGER_URL" env-default:"http://localhost:14268/api/traces"`
		Flags         flags
		LogLevel      string `env:"LOG_LEVEL" env-default:"debug"`
	}
)

func LoadConfig() (Config, error) {
	var cfg Config

	cfg.Flags = loadFlags()
	if cfg.Flags.DevMode {
		cfg.LogLevel = "dev"
	}

	if cfg.Flags.envFilename != "" {
		if err := cleanenv.ReadConfig(cfg.Flags.envFilename, &cfg); err != nil {
			return Config{}, err
		}
		cfg.Server.Port = ":" + cfg.Server.Port
		return cfg, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, err
	}

	cfg.Server.Port = ":" + cfg.Server.Port

	return cfg, nil
}

func loadFlags() flags {
	var f flags

	flag.BoolVar(&f.DevMode, "dev", false, "Run in dev mode, some features will be disabled.\nFor example, emails will be printed to stdout instead of being sent.")
	flag.StringVar(&f.envFilename, "env", "", "Path to .env file.")
	flag.Parse()

	return f
}
