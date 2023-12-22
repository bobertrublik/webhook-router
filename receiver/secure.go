package receiver

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/whosonfirst/go-webhookd/v3"
)

func init() {

	ctx := context.Background()
	err := RegisterReceiver(ctx, "secure", NewSecureReceiver)

	if err != nil {
		panic(err)
	}
}

// LogReceiver implements the `webhookd.WebhookReceiver` interface for receiving webhook messages in a secure fashion.
type SecureReceiver struct {
	webhookd.WebhookReceiver
}

// NewSecureReceiver returns a new `SecureReceiver` instance configured by 'uri' in the form of:
//
//	insecure://
func NewSecureReceiver(ctx context.Context, uri string) (webhookd.WebhookReceiver, error) {

	wh := SecureReceiver{}
	return wh, nil
}

// Receive returns the body of the message in 'req'. It does not check its provenance or validate the message body in any way. You should not use this in production.
func (wh SecureReceiver) Receive(ctx context.Context, req *http.Request) ([]byte, *webhookd.WebhookError) {

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

	fmt.Printf("http request: %v", req)
	authHeader := req.Header.Get("Authorization")
	fmt.Printf("authoHeader: %v", authHeader)
	if authHeader == "" {

		code := http.StatusUnauthorized
		message := "Authorization header missing"
		return nil, &webhookd.WebhookError{Code: code, Message: message}
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
