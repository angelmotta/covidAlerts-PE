package util

import (
	"github.com/spf13/viper"
)

// Values are read by Viper from a config file or environment variables
type Config struct {
	DBHost				string `mapstructure:"DB_HOST"`
	DBPort				string `mapstructure:"DB_PORT"`
	DBUser				string `mapstructure:"DB_USER"`
	DBPass				string `mapstructure:"DB_PASS"`
	DBName				string `mapstructure:"DB_NAME"`
	TApiKey				string `mapstructure:"TWITTER_API_KEY"`
	TApiSecretKey		string `mapstructure:"TWITTER_API_SECRET_KEY"`
	TAccessToken		string `mapstructure:"TWITTER_ACCESS_TOKEN"`
	TAccessTokenSecret	string `mapstructure:"TWITTER_ACCESS_TOKEN_SECRET"`
	UrlNewCases			string `mapstructure:"URL_POSITIVE_CASES"`
	UrlDeceased			string `mapstructure:"URL_DECEASED_CASES"`
	DirPositiveFiles	string `mapstructure:"dataFiles/"`
	DirDeceasedFiles	string `mapstructure:"dataFiles/"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(pathFile string) (config Config, err error) {
	viper.AddConfigPath(pathFile)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	// Viper will also read env variables and override values if these exist
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Read values into Config object
	err = viper.Unmarshal(&config)
	return
}