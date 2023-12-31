package dispatcher

import (
	"context"
	"fmt"
	"github.com/bobertrublik/webhook-router/internal/webhookd"
	"github.com/sfomuseum/go-slack/writer"

	"io"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterDispatcher(ctx, "slack", NewSlackDispatcher)

	if err != nil {
		panic(err)
	}
}

// For backwards compatibility

type SlackcatConfig struct {
	WebhookUrl string `json:"webhook_url"`
	Channel    string `json:"channel"`
}

type SlackDispatcher struct {
	webhookd.WebhookDispatcher
	writer io.Writer
}

func NewSlackDispatcher(ctx context.Context, uri string) (webhookd.WebhookDispatcher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %v", err)
	}

	q := u.Query()

	wh_uri := q.Get("webhook")
	wh_channel := q.Get("channel")

	wr, err := writer.NewSlackWriter(wh_uri, wh_channel)

	if err != nil {
		return nil, err
	}

	slack := SlackDispatcher{
		writer: wr,
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

	_, err := sl.writer.Write(body)

	if err != nil {
		code := 999
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return err
	}

	return nil
}
