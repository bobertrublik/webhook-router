// Package daemon provides methods for implementing a long-running daemon to listen for and process webhooks.
package daemon

import (
	"context"
	"fmt"
	"github.com/bobertrublik/webhook-router/internal/logger"
	"github.com/bobertrublik/webhook-router/internal/webhookd"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bobertrublik/webhook-router/internal/config"
	"github.com/bobertrublik/webhook-router/internal/dispatcher"
	"github.com/bobertrublik/webhook-router/internal/receiver"
	"github.com/bobertrublik/webhook-router/internal/transformation"
	"github.com/bobertrublik/webhook-router/internal/webhook"
)

// type WebhookDaemon is a struct that implements a long-running daemon to listen for	and process webhooks.
type WebhookDaemon struct {
	// webhooks is a dictionary of URIs and their corresponding `webhookd.WebhookHandler` instances.
	webhooks map[string]webhookd.WebhookHandler
}

// NewWebhookDaemonFromConfig() returns a new `WebhookDaemon` derived from configuration data in 'cfg'.
func NewWebhookDaemonFromConfig(ctx context.Context, cfg *config.WebhookConfig) (*WebhookDaemon, error) {

	webhooks := make(map[string]webhookd.WebhookHandler)

	d := WebhookDaemon{
		webhooks: webhooks,
	}

	err := d.AddWebhooksFromConfig(ctx, cfg)

	if err != nil {
		return nil, fmt.Errorf("Failed to add webhooks to daemon, %w", err)
	}

	return &d, nil
}

// AddWebhooksFromConfig() appends the webhooks defined in 'cfg' to 'd'.
func (d *WebhookDaemon) AddWebhooksFromConfig(ctx context.Context, cfg *config.WebhookConfig) error {

	if len(cfg.Webhooks) == 0 {
		return fmt.Errorf("No webhooks defined")
	}

	for i, hook := range cfg.Webhooks {

		if hook.Endpoint == "" {
			return fmt.Errorf("Missing endpoint at offset %d", i+1)
		}

		if hook.Receiver == "" {
			return fmt.Errorf("Missing receiver at offset %d", i+1)
		}

		if len(hook.Dispatchers) == 0 {
			return fmt.Errorf("Missing dispatchers at offset %d", i+1)
		}

		recvUri, err := cfg.GetReceiverConfigByName(hook.Receiver)

		if err != nil {
			return fmt.Errorf("Failed to get receiver config for '%s', %w", hook.Receiver, err)
		}

		recv, err := receiver.NewReceiver(ctx, recvUri)

		if err != nil {
			return fmt.Errorf("Failed to add receiver '%s', %w", recvUri, err)
		}

		var steps []webhookd.WebhookTransformation

		for _, name := range hook.Transformations {

			if strings.HasPrefix(name, "#") {
				continue
			}

			transfUri, err := cfg.GetTransformationConfigByName(name)

			if err != nil {
				return fmt.Errorf("Failed to get transformation configuration for '%s', %w", name, err)
			}

			step, err := transformation.NewTransformation(ctx, transfUri)

			if err != nil {
				return fmt.Errorf("Failed to create new transformation for '%s', %w", transfUri, err)
			}

			steps = append(steps, step)
		}

		var sendto []webhookd.WebhookDispatcher

		for _, name := range hook.Dispatchers {

			if strings.HasPrefix(name, "#") {
				continue
			}

			dispUri, err := cfg.GetDispatcherConfigByName(name)

			if err != nil {
				return fmt.Errorf("Failed to get dispatcher configuration for '%s', %w", name, err)
			}

			disp, err := dispatcher.NewDispatcher(ctx, dispUri)

			if err != nil {
				return fmt.Errorf("Failed to create dispatcher for '%s', %w", dispUri, err)
			}

			sendto = append(sendto, disp)
		}

		wh, err := webhook.NewWebhook(ctx, hook.Endpoint, recv, steps, sendto)

		if err != nil {
			return fmt.Errorf("Failed to create new webhook for '%s', %w", hook.Endpoint, err)
		}

		err = d.AddWebhook(ctx, wh)

		if err != nil {
			return fmt.Errorf("Failed to add new webhook for '%s', %w", hook.Endpoint, err)
		}

	}

	return nil
}

// AddWebhook() adds 'wh' to 'd'.
func (d *WebhookDaemon) AddWebhook(ctx context.Context, wh webhook.Webhook) error {

	endpoint := wh.Endpoint()
	_, ok := d.webhooks[endpoint]

	if ok {
		return fmt.Errorf("endpoint already configured")
	}

	d.webhooks[endpoint] = wh
	return nil
}

// HandlerFuncWithLogger() returns a `http.HandlerFunc` that handles HTTP (webhook) requests and response for 'd'
// logging events to 'logger'.
func (d *WebhookDaemon) ProcessRequest(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	endpoint := r.URL.Path

	wh, ok := d.webhooks[endpoint]

	if !ok {
		http.Error(w, "404 Not found", http.StatusNotFound)
		return fmt.Errorf("Endpoint not found, %s", endpoint)
	}

	t1 := time.Now()

	var ta time.Time
	var tb time.Duration

	var ttr time.Duration // time to receive
	var ttt time.Duration // time to transform
	var ttd time.Duration // time to dispatch

	ta = time.Now()

	rcvr := wh.Receiver()

	body, err := rcvr.Receive(ctx, r)

	// we use -1 to signal that this is an unhandled event but
	// not an error, for example when github sends a ping message
	// (20190212/thisisaaronland)

	if err != nil {

		switch err.Code {
		case webhookd.UnhandledEvent, webhookd.HaltEvent:
			logger.Log.Info("Receiver step (%T)  returned non-fatal error and exiting, %v", rcvr, err)
			return nil
		default:
			http.Error(w, err.Error(), err.Code)
			return fmt.Errorf("Receiver step (%T) failed, %v", rcvr, err)
		}
	}

	tb = time.Since(ta)

	ttr = tb

	ta = time.Now()

	for idx, step := range wh.Transformations() {

		body, err = step.Transform(ctx, body)

		if err != nil {

			switch err.Code {
			case webhookd.UnhandledEvent, webhookd.HaltEvent:
				logger.Log.Info("Transformation step (%T) at offset %d returned non-fatal error and exiting, %v", step, idx, err)
				return nil
			default:
				http.Error(w, err.Error(), err.Code)
				return fmt.Errorf("Transformation step (%T) at offset %d failed, %v", step, idx, err)
			}
		}

		// check to see if there is anything left the transformation
		// https://github.com/whosonfirst/go-webhookd/v3/issues/7
	}

	tb = time.Since(ta)
	ttt = tb

	// check to see if there is anything to dispatch
	// https://github.com/whosonfirst/go-webhookd/v3/issues/7

	ta = time.Now()

	wg := new(sync.WaitGroup)
	ch := make(chan *webhookd.WebhookError)

	for idx, di := range wh.Dispatchers() {

		wg.Add(1)

		go func(idx int, di webhookd.WebhookDispatcher, body []byte) {

			defer wg.Done()

			err = di.Dispatch(ctx, body)

			if err != nil {

				switch err.Code {
				case webhookd.UnhandledEvent, webhookd.HaltEvent:
					logger.Log.Info("Dispatch step (%T) at offset %d returned non-fatal error and exiting, %v", d, idx, err)
					return
				default:
					logger.Log.Error("Dispatch step (%T) at offset %d failed, %v", d, idx, err)
					ch <- err
				}
			}

		}(idx, di, body)
	}

	// https://github.com/whosonfirst/go-webhookd/issues/14
	// this is broken as in len(errors) will always be zero even if
	// there are errors (20190214/thisisaaronland)

	errors := make([]string, 0)

	go func() {

		for e := range ch {
			errors = append(errors, e.Error())
		}
	}()

	wg.Wait()

	if len(errors) > 0 {

		msg := strings.Join(errors, "\n\n")
		http.Error(w, msg, http.StatusInternalServerError)
		return fmt.Errorf("Internal Server Error")
	}

	tb = time.Since(ta)
	ttd = tb

	t2 := time.Since(t1)

	logger.Log.Debug("Time to receive: %v", ttr)
	logger.Log.Debug("Time to transform: %v", ttt)
	logger.Log.Debug("Time to dispatch: %v", ttd)
	logger.Log.Debug("Time to process: %v", t2)

	w.Header().Set("X-Webhookd-Time-To-Receive", fmt.Sprintf("%v", ttr))
	w.Header().Set("X-Webhookd-Time-To-Transform", fmt.Sprintf("%v", ttt))
	w.Header().Set("X-Webhookd-Time-To-Dispatch", fmt.Sprintf("%v", ttd))
	w.Header().Set("X-Webhookd-Time-To-Process", fmt.Sprintf("%v", t2))

	return nil

}

// Start() causes 'd' to listen for, and process, requests.
func (d *WebhookDaemon) Start(w http.ResponseWriter, r *http.Request) error {
	err := d.ProcessRequest(w, r)
	if err != nil {
		return err
	}
	return nil
}
