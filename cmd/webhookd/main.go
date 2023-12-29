// webhookd is a command line tool to start a go-webhookd daemon and serve requests over HTTP.
package main

import (
	"context"
	"fmt"
	"github.com/bobertrublik/webhook-router/internal/config"
	"github.com/bobertrublik/webhook-router/internal/daemon"
	"github.com/bobertrublik/webhook-router/internal/logger"
	"github.com/bobertrublik/webhook-router/internal/router"
	"github.com/joho/godotenv"
	"github.com/sfomuseum/go-flags/flagset"
	"net/http"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		logger.Log.Error("Error loading the .env file: %v", err)
		os.Exit(1)
	}

	fs := flagset.NewFlagSet("webhooks")

	configFile := fs.String("config", "/etc/config/config.yaml", "Path to config file")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "webhookd is a command line tool to start a go-webhookd daemon and serve requests over HTTP.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	cfg, err := config.NewConfig(*configFile)

	if err != nil {
		logger.Log.Error("Failed to load config %s, %v", configFile, err)
		os.Exit(1)
	}
	ctx := context.Background()

	webhookDaemon, err := daemon.NewWebhookDaemonFromConfig(ctx, cfg)

	if err != nil {
		logger.Log.Error("Failed to create webhook daemon, %v", err)
		os.Exit(1)
	}

	rtr := router.New(webhookDaemon)

	logger.Log.Info("Server listening on http://localhost:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", rtr); err != nil {
		logger.Log.Error("There was an error with the http server: %v", err)
		os.Exit(1)
	}
}
