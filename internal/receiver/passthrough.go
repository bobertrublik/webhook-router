package receiver

import (
	"context"
	"github.com/bobertrublik/webhook-router/internal/webhookd"
	"io"
	"net/http"
)

func init() {

	ctx := context.Background()
	err := RegisterReceiver(ctx, "passthrough", NewPassThroughReceiver)

	if err != nil {
		panic(err)
	}
}

// LogReceiver implements the `webhookd.WebhookReceiver` interface for receiving webhook messages in an insecure fashion.
type PassThroughReceiver struct {
	webhookd.WebhookReceiver
}

// NewPassThroughReceiver returns a new `PassThroughReceiver` instance configured by 'uri' in the form of:
//
//	insecure://
func NewPassThroughReceiver(ctx context.Context, uri string) (webhookd.WebhookReceiver, error) {

	wh := PassThroughReceiver{}
	return wh, nil
}

// Receive returns the body of the message in 'req'. It does not check its provenance or validate the message body in any way. You should not use this in production.
func (wh PassThroughReceiver) Receive(ctx context.Context, req *http.Request) ([]byte, *webhookd.WebhookError) {

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	if req.Method != "POST" {

		code := http.StatusMethodNotAllowed
		message := "Method not allowed"

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {

		code := http.StatusInternalServerError
		message := err.Error()

		err := &webhookd.WebhookError{Code: code, Message: message}
		return nil, err
	}

	return body, nil
}
