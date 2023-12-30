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
	err := RegisterDispatcher(ctx, "echo", NewEchoDispatcher)

	if err != nil {
		panic(err)
	}
}

// EchoDispatcher implements the `webhookd.WebhookDispatcher` interface for dispatching messages to nowhere.
type EchoDispatcher struct {
	webhookd.WebhookDispatcher
	endpoint string
}

// NewEchoDispatcher returns a new `EchoDispatcher` instance that dispatches messages to nowhere
// configured by 'uri' in the form of:
//
//	echo://
func NewEchoDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {
	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}
	uri = fmt.Sprintf("http://%s", u.Host)
	d := EchoDispatcher{
		endpoint: uri,
	}
	return &d, nil
}

// Dispatch sends 'body' to nowhere.
func (d *EchoDispatcher) Dispatch(ctx context.Context, body []byte) *webhookd.WebhookError {
	responseBody := bytes.NewBuffer(body)
	// Create a new HTTP request
	resp, err := http.Post(d.endpoint, "application/json", responseBody)
	//Handle Error
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
