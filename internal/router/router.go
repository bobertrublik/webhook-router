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
	router.Handle("/api/public", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello from a public endpoint! You don't need to be authenticated to see this."}`))
	}))

	// This route is only accessible if the user has a valid access_token.
	router.Handle("/api", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Info("Retrieved new request")
			err := webhookDaemon.ProcessRequest(w, r)
			if err != nil {
				fmt.Println(err)
			}
		}),
	))

	return router
}
