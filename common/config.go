package common

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Service ServiceConfig `toml:"service"`
	Client  ClientConfig  `toml:"client"`
}

func DefaultConfig() Config {
	return Config{
		Service: DefaultServiceConfig(),
		Client:  DefaultClientConfig(),
	}
}

func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		ListenAddress:          "unix:///tmp/shellie.sock",
		ChatCompletionEndpoint: "https://api.openai.com/v1",
		APIKey:                 "",
		Organization:           "",
		Project:                "",
		Model:                  "gpt-4o-mini",
	}
}

func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		ServerAddress: "unix:///tmp/shellie.sock",
	}
}

type ClientConfig struct {
	ServerAddress string `toml:"server_address" comment:"Default will use unix socket"`
}

type ServiceConfig struct {
	ListenAddress          string `toml:"listen_address" comment:"Address to listen on for the service, defaults to unix socket"`
	ChatCompletionEndpoint string `toml:"chat_completion_endpoint" comment:"OpenAI chat completion compatible endpoint"`
	APIKey                 string `toml:"api_key" comment:"OpenAI API key"`
	Organization           string `toml:"organization" comment:"OpenAI organization"`
	Project                string `toml:"project" comment:"OpenAI project"`
	Model                  string `toml:"model" comment:"OpenAI model"`
}

func ReadOrCreateConfig() (*Config, error) {
	// Read config from ~/.config/promptsuggestion.json
	HOME := os.Getenv("HOME")
	// create default config
	// check if ~/.config/promptsuggestion.json exists
	config := DefaultConfig()
	if _, err := os.Stat(fmt.Sprintf("%s/.config/shellie.toml", HOME)); os.IsNotExist(err) {
		content, err := toml.Marshal(config)
		if err != nil {
			return nil, err
		}
		os.WriteFile(fmt.Sprintf("%s/.config/shellie.toml", HOME), content, 0644)
	}
	configBytes, err := os.ReadFile(fmt.Sprintf("%s/.config/shellie.toml", HOME))
	if err != nil {
		return nil, err
	}
	if err := toml.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
