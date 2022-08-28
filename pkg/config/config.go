package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
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
	Addr               string
	CommandPrefix      string `toml:"command_prefix"`
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

type Root struct {
	Weissbot   Weissbot
	Discord    Discord
	Line       Line
	Management Management
	Twitter    Twitter
}

var defaultConfig = Root{
	Weissbot: Weissbot{
		Environment: DefaultEnvironment,
	},
}

func Load(s string) (*Root, error) {
	cfg := &defaultConfig

	if _, err := toml.Decode(s, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadFile(filename string) (*Root, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("loading config file: %w", err)
	}

	rootCfg, err := Load(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", filename, err)
	}

	return rootCfg, nil
}
