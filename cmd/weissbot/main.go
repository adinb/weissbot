package main

import (
	"context"
	"github.com/adinb/weissbot/pkg/config"
	"github.com/adinb/weissbot/pkg/discord"
	"github.com/adinb/weissbot/pkg/line"
	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const DefaultConfigFilePath = "weissbot.toml"

type flagConfig struct {
	configFilePath string
}

func main() {
	flagCfg := flagConfig{}
	flag.StringVarP(
		&flagCfg.configFilePath,
		"config-file",
		"f",
		DefaultConfigFilePath,
		"Path to Weissbot configuration file",
	)

	flag.Parse()

	var logger = log.New(os.Stderr, "", log.Ltime|log.Ldate|log.LUTC)
	cfg, err := config.LoadFile(flagCfg.configFilePath)
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	logger.Printf("Environment: %s", cfg.Weissbot.Environment)

	// Handle signal
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	errg, gctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}

	signalc := make(chan os.Signal)
	defer close(signalc)
	signal.Notify(signalc, os.Interrupt, syscall.SIGTERM)

	if cfg.Discord.Enabled {
		var d *discord.DiscordController
		errg.Go(func() error {
			if d, err = discord.New(cfg, logger); err != nil {
				return err
			}

			return d.Start()
		})

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-gctx.Done()
			err = d.Stop()
			if err != nil {
				logger.Println(err)
			} else {
				logger.Println("Discord bot stopped")
			}
		}()
	}

	if cfg.Line.Enabled {
		var lineserver *http.Server
		errg.Go(func() error {
			if lineserver, err = line.New(cfg); err != nil {
				return err
			}

			err = line.Start(lineserver, logger)
			if err != nil && err != http.ErrServerClosed {
				return err
			}

			return nil
		})

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-gctx.Done()
			err = lineserver.Shutdown(ctx)
			if err != nil {
				logger.Println(err)
			} else {
				logger.Println("LINE handler stopped")
			}
		}()
	}

	// Block until a SIGINT/SIGTERM signal is received
	// Or until the errgroup is done / canceled
	select {
	case sig := <-signalc:
		logger.Printf("Received %s signal", sig)
		stop()
		wg.Wait()
		err = errg.Wait()
		if err != nil {
			logger.Println(err)
		}

	case <-gctx.Done():
		logger.Printf("Error: %s", errg.Wait())
		wg.Wait()
		os.Exit(1)
	}
}
