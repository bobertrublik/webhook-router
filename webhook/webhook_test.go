package webhook

import (
	"context"
	"testing"

	"github.com/bobertrublik/webhook-router/dispatcher"
	"github.com/bobertrublik/webhook-router/receiver"
	"github.com/bobertrublik/webhook-router/transformation"
)

func TestWebhook(t *testing.T) {

	ctx := context.Background()

	r, err := receiver.NewReceiver(ctx, "insecure://")

	if err != nil {
		t.Fatalf("Failed to create new receiver, %v", err)
	}

	tr, err := transformation.NewTransformation(ctx, "null://")

	if err != nil {
		t.Fatalf("Failed to create new transformation, %v", err)
	}

	d, err := dispatcher.NewDispatcher(ctx, "log://")

	if err != nil {
		t.Fatalf("Failed to create new dispatcher, %v", err)
	}

	endpoint := "/insecure"

	wh, err := NewWebhook(ctx, endpoint, r, []webhookd.WebhookTransformation{tr}, []webhookd.WebhookDispatcher{d})

	if err != nil {
		t.Fatalf("Failed to create new webhook, %v", err)
	}

	if wh.Endpoint() != endpoint {
		t.Fatalf("Invalid endpoint: %s", wh.Endpoint())
	}

	if wh.Receiver() != r {
		t.Fatalf("Unexpected receiver")
	}

}
