package config

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Timezone         string `mapstructure:"timezone"`
	PostgresHost     string `mapstructure:"postgres_host"`
	PostgresPort     string `mapstructure:"postgres_port"`
	PostgresDB       string `mapstructure:"postgres_db"`
	PostgresUser     string `mapstructure:"postgres_user"`
	PostgresPassword string `mapstructure:"postgres_password"`
}

func (c Config) PostgresDSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.PostgresUser, c.PostgresPassword),
		Host:   c.PostgresHost + ":" + c.PostgresPort,
		Path:   "/" + c.PostgresDB,
	}
	q := url.Values{"sslmode": []string{"disable"}}
	u.RawQuery = q.Encode()
	return u.String()
}

func MustLoad() *Config {
	v := viper.New()
	v.SetConfigType("yaml")

	cf := os.Getenv("CONFIG_FILE")
	if cf == "" {
		log.Fatal("CONFIG_FILE is required")
	}
	v.SetConfigFile(cf)

	_ = v.BindEnv("postgres_user", "POSTGRES_USER")
	_ = v.BindEnv("postgres_password", "POSTGRES_PASSWORD")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("config: cannot read %s: %v", cf, err)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		log.Fatalf("config: unmarshal: %v", err)
	}
	validate(&cfg)

	if _, err := time.LoadLocation(cfg.Timezone); err != nil {
		log.Fatalf("config: invalid timezone %q: %v", cfg.Timezone, err)
	}
	return &cfg
}

func validate(cfg *Config) {
	var missing []string

	if cfg.Timezone == "" {
		missing = append(missing, "timezone")
	}
	if cfg.PostgresHost == "" {
		missing = append(missing, "postgres_host")
	}
	if cfg.PostgresPort == "" {
		missing = append(missing, "postgres_port")
	}
	if cfg.PostgresDB == "" {
		missing = append(missing, "postgres_db")
	}
	if cfg.PostgresUser == "" {
		missing = append(missing, "POSTGRES_USER env")
	}
	if cfg.PostgresPassword == "" {
		missing = append(missing, "POSTGRES_PASSWORD env")
	}

	if len(missing) > 0 {
		log.Fatalf("config: missing/invalid keys:\n  - %s", strings.Join(missing, "\n  - "))
	}
}
