package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	PSQLSource string `mapstructure:"PSQL_SOURCE" validate:"required"`
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
		panic("config validation failed: " + err.Error())
	}

	return &cfg
}
