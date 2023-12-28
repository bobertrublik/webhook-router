// webhookd is a command line tool to start a go-webhookd daemon and serve requests over HTTP.
package main

import (
	"context"
	"fmt"
	"github.com/bobertrublik/webhook-router/internal/config"
	"github.com/bobertrublik/webhook-router/internal/daemon"
	"github.com/bobertrublik/webhook-router/internal/router"
	"github.com/joho/godotenv"
	"github.com/sfomuseum/go-flags/flagset"
	"log"
	"net/http"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading the .env file: %v", err)
	}

	fs := flagset.NewFlagSet("webhooks")

	config_uri := fs.String("config-uri", "", "A valid Go Cloud runtimevar URI representing your webhookd config.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "webhookd is a command line tool to start a go-webhookd daemon and serve requests over HTTP.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVarsWithFeedback(fs, "WEBHOOKD", false)

	if err != nil {
		log.Fatalf("Failed to set flags from env vars, %v", err)
	}

	ctx := context.Background()

	cfg, err := config.NewConfigFromURI(ctx, *config_uri)

	if err != nil {
		log.Fatalf("Failed to load config %s, %v", *config_uri, err)
	}

	wh_daemon, err := daemon.NewWebhookDaemonFromConfig(ctx, cfg)

	if err != nil {
		log.Fatalf("Failed to create webhook daemon, %v", err)
	}

	rtr := router.New(wh_daemon)

	log.Print("Server listening on http://localhost:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
