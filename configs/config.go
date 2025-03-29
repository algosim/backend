package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int    `mapstructure:"port" env:"SERVER_PORT"`
		Host string `mapstructure:"host" env:"SERVER_HOST"`
	} `mapstructure:"server"`

	Auth struct {
		JWTSecret string `mapstructure:"jwt_secret" env:"AUTH_JWT_SECRET"`
		TokenTTL  int    `mapstructure:"token_ttl" env:"AUTH_TOKEN_TTL"`
	} `mapstructure:"auth"`

	GoogleOAuth struct {
		ClientID     string   `mapstructure:"client_id" env:"GOOGLE_OAUTH_CLIENT_ID"`
		ClientSecret string   `mapstructure:"client_secret" env:"GOOGLE_OAUTH_CLIENT_SECRET"`
		RedirectURI  string   `mapstructure:"redirect_uri" env:"GOOGLE_OAUTH_REDIRECT_URI"`
		Scopes       []string `mapstructure:"scopes"`
	} `mapstructure:"google_oauth"`
}

var globalConfig *Config

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("auth.token_ttl", 3600)
	viper.SetDefault("google_oauth.scopes", []string{"openid", "email", "profile"})

	// Configure environment variable handling
	viper.AutomaticEnv()
	viper.SetEnvPrefix("") // Don't use a prefix
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	globalConfig = config
	return config, nil
}

// Get returns the global config instance
func Get() *Config {
	return globalConfig
}
