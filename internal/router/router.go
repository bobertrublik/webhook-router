package router

import (
	"fmt"
	"github.com/bobertrublik/webhook-router/internal/daemon"
	"github.com/bobertrublik/webhook-router/internal/logger"
	"net/http"

	"github.com/bobertrublik/webhook-router/internal/middleware"
)

// New sets up our routes and returns a *http.ServeMux.
func New(webhookDaemon *daemon.WebhookDaemon) *http.ServeMux {
	router := http.NewServeMux()

	// This route is always accessible.
	router.Handle("/echo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info("Webhook request received on path /echo")
		w.Header().Set("Content-Type", "application/json")
		err := webhookDaemon.ProcessRequest(w, r)
		if err != nil {
			fmt.Println(err)
		}
	}))

	// This route is only accessible if the user has a valid access_token.
	router.Handle("/api", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Info("Webhook request received on path /api")
			err := webhookDaemon.ProcessRequest(w, r)
			if err != nil {
				fmt.Println(err)
			}
		}),
	))

	return router
}
