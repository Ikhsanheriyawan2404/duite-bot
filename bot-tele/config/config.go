package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken  string `mapstructure:"TELEGRAM_TOKEN"`
	CoreApiUrl     string `mapstructure:"CORE_API_URL"`
	DashboardUrl   string `mapstructure:"DASHBOARD_URL"`
	AppEnv		   string `mapstructure:"APP_ENV"`
}

var AppConfig Config

func LoadConfig(path string) error {
	viper.SetConfigName(".env") // name of config file (without extension)
	viper.SetConfigType("env")  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)   // path to look for the config file in
	viper.AutomaticEnv()        // read in environment variables that match

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
		{"TELEGRAM_TOKEN", AppConfig.TelegramToken},
		{"DASHBOARD_URL", AppConfig.DashboardUrl},
		{"CORE_API_URL", AppConfig.CoreApiUrl},
		{"APP_ENV", AppConfig.AppEnv},
	}

	for _, req := range required {
		if req.value == "" {
			return fmt.Errorf("%s is required", req.field)
		}
	}

	return nil
}