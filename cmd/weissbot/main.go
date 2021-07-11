package main

import (
	"fmt"
	"github.com/adinb/weissbot/pkg/discord"
	"os/signal"
	"syscall"

	"github.com/adinb/weissbot/pkg/config"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

const DefaultConfigFilePath = "weissbot.toml"

type flagConfig struct {
	configFilePath string
}

func main() {
	var logger = log.New(os.Stderr, "", log.Ltime|log.Ldate|log.LUTC)

	flagCfg := flagConfig{}
	flag.StringVarP(
		&flagCfg.configFilePath,
		"config-file",
		"f",
		DefaultConfigFilePath,
		"Path to Weissbot configuration file",
	)

	flag.Parse()
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	cfg, err := config.LoadFile(flagCfg.configFilePath)
	if err != nil {
		logger.Println("Could not load the configuration file: %s", flagCfg.configFilePath)
		logger.Fatalf(err.Error())
	}

	logger.Println("Starting Weissbot")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	if cfg.Discord.Enabled {
		logger.Println("Starting Discord bot")
		d, err := discord.NewDiscordService(
			cfg.Weissbot.Environment,
			cfg.Twitter.Token,
			cfg.Discord.Token,
			logger,
		)
		if err == nil {
			err = d.Start()
			if err != nil {
				logger.Println("Failed to start Discord bot, continuing")
			}

			defer func() {
				logger.Println("Stopping Discord bot")
				err := d.Stop()
				if err != nil {
					logger.Println(err)
				}
			}()
		} else {
			logger.Println(err)
		}
	}

	for {
		select {
		case sig := <-sc:
			logger.Printf("Received signal: %s, shutting down weissbot", sig)
			return
		}
	}
}
