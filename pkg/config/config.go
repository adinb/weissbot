package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

const DefaultEnvironment = "test"

type Weissbot struct {
	Environment string
}

type Discord struct {
	Enabled bool
	Token   string
}

type Line struct {
	Enabled            bool
	Port               uint16
	ChannelSecret      string `toml:"channel_secret"`
	ChannelAccessToken string `toml:"channel_access_token"`
}

type Management struct {
	Enabled bool
	ApiKey  string `toml:"api_key"`
}

type Twitter struct {
	Enabled bool
	Token   string
}

type Config struct {
	Weissbot   Weissbot
	Discord    Discord
	Line       Line
	Management Management
	Twitter    Twitter
}

var defaultConfig = Config{
	Weissbot: Weissbot{
		Environment: DefaultEnvironment,
	},
}

func Load(s string) (*Config, error) {
	cfg := &defaultConfig

	if _, err := toml.Decode(s, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg, err := Load(string(content))
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to parse configuration file %s: %s", filename, err))
		return nil, err
	}

	return cfg, nil
}
