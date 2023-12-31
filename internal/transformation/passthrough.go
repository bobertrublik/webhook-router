package transformation

import (
	"context"
	"github.com/bobertrublik/webhook-router/internal/webhookd"
)

func init() {

	ctx := context.Background()
	err := RegisterTransformation(ctx, "passthrough", NewPassThroughTransformation)

	if err != nil {
		panic(err)
	}
}

// PassThroughTransformation implements the `webhookd.WebhookTransformation` interface for a no-op transformation meaning
// the output of the `Transform` method is the same as its input.
type PassThroughTransformation struct {
	webhookd.WebhookTransformation
}

// NewInsecureTransformation returns a new `PassThroughTransformation` instance configured by 'uri' in the form of:
//
//	pass://
func NewPassThroughTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	p := PassThroughTransformation{}
	return &p, nil
}

// Transform returns 'body' unaltered.
func (p *PassThroughTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {
	return body, nil
}
