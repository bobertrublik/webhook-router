package dispatcher

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bobertrublik/webhook-router/internal/webhookd"
	"io"
	"net/http"

	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterDispatcher(ctx, "slack", NewSlackDispatcher)

	if err != nil {
		panic(err)
	}
}

type SlackDispatcher struct {
	webhookd.WebhookDispatcher
	incomingWebhook string
}

func NewSlackDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %v", err)
	}

	q := u.Query()

	slack := SlackDispatcher{
		incomingWebhook: q.Get("webhook"),
	}

	return &slack, nil
}

func (sl *SlackDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {

	select {
	case <-ctx.Done():
		return nil
	default:
		// pass
	}

	responseBody := bytes.NewBuffer(body)
	// Create a new HTTP request
	resp, err := http.Post(sl.incomingWebhook, "application/json", responseBody)

	if err != nil {
		code := 999
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return err
	}
	defer resp.Body.Close()
	//Read the response body
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		code := 999
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return err
	}

	return nil
}
