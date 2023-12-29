package transformation

import (
	"context"
	"github.com/bobertrublik/webhook-router/internal/webhookd"
)

func init() {

	ctx := context.Background()
	err := RegisterTransformation(ctx, "pass", NewNullTransformation)

	if err != nil {
		panic(err)
	}
}

// NullTransformation implements the `webhookd.WebhookTransformation` interface for a no-op transformation meaning
// the output of the `Transform` method is the same as its input.
type NullTransformation struct {
	webhookd.WebhookTransformation
}

// NewInsecureTransformation returns a new `NullTransformation` instance configured by 'uri' in the form of:
//
//	pass://
func NewNullTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	p := NullTransformation{}
	return &p, nil
}

// Transform returns 'body' unaltered.
func (p *NullTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {
	return body, nil
}
