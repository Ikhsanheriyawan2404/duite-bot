package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LLMApiKey      string `mapstructure:"LLM_API_KEY"`
	ServerPort     string `mapstructure:"SERVER_PORT"`

	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`

	LLMApiUrl 	   string `mapstructure:"LLM_API_URL"`
}

var AppConfig Config

func LoadConfig(path string) error {
	viper.SetConfigName(".env") // name of config file (without extension)
	viper.SetConfigType("env")  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)   // path to look for the config file in
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	viper.AutomaticEnv()        // read in environment variables that match

	// Set defaults
	viper.SetDefault("SERVER_PORT", "80")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No .env file found, using environment variables or defaults")
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config into struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("unable to decode into config struct: %w", err)
	}

	// Validate required config
	required := []struct {
		field string
		value string
	}{
		{"LLM_API_KEY", AppConfig.LLMApiKey},
		{"LLM_API_URL", AppConfig.LLMApiUrl},
		{"DB_HOST", AppConfig.DBHost},
		{"DB_PORT", AppConfig.DBPort},
		{"DB_USER", AppConfig.DBUser},
		{"DB_PASSWORD", AppConfig.DBPassword},
		{"DB_NAME", AppConfig.DBName},
	}

	for _, req := range required {
		if req.value == "" {
			return fmt.Errorf("%s is required", req.field)
		}
	}

	return nil
}