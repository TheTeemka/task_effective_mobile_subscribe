package config

import (
	"fmt"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel   string `mapstructure:"LOG_LEVEL" validate:"required,oneof=DEBUG INFO WARN ERROR"`
	PSQLSource string `mapstructure:"PSQL_SOURCE"`
	// optional components to compose PSQL source if PSQL_SOURCE is not provided
	PSQLUser     string `mapstructure:"PSQL_USER"`
	PSQLPassword string `mapstructure:"PSQL_PASSWORD"`
	PSQLHost     string `mapstructure:"PSQL_HOST"`
	PSQLPort     string `mapstructure:"PSQL_PORT"`
	PSQLDB       string `mapstructure:"PSQL_DB"`
	PSQLSSLMode  string `mapstructure:"PSQL_SSLMODE"`
	Port       string `mapstructure:"PORT" validate:"required"`
}

func LoadConfig() *Config {
	godotenv.Load(".env")

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic(err)
	}
	// if PSQLSource was not provided directly, try to compose it from components
	if cfg.PSQLSource == "" {
		// prefer explicit viper values (env or .env)
		user := viper.GetString("PSQL_USER")
		pass := viper.GetString("PSQL_PASSWORD")
		host := viper.GetString("PSQL_HOST")
		if host == "" {
			host = "localhost"
		}
		port := viper.GetString("PSQL_PORT")
		if port == "" {
			port = "5432"
		}
		db := viper.GetString("PSQL_DB")
		ssl := viper.GetString("PSQL_SSLMODE")
		if ssl == "" {
			ssl = "disable"
		}

		if db != "" {
			u := &url.URL{
				Scheme: "postgres",
				Host:   fmt.Sprintf("%s:%s", host, port),
				Path:   db,
			}
			if user != "" {
				if pass != "" {
					u.User = url.UserPassword(user, pass)
				} else {
					u.User = url.User(user)
				}
			}
			q := u.Query()
			q.Set("sslmode", ssl)
			u.RawQuery = q.Encode()
			cfg.PSQLSource = u.String()
		}
	}

	if cfg.PSQLSource == "" {
		panic("PSQL_SOURCE is required: set PSQL_SOURCE or PSQL_USER/PSQL_PASSWORD/PSQL_HOST/PSQL_PORT/PSQL_DB")
	}

	return &cfg
}
