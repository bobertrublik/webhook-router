package transformations

import (
	"context"
	"github.com/aaronland/go-chicken"
	"github.com/whosonfirst/go-webhookd/v2"
	"net/url"
	"strconv"
)

type ChickenTransformation struct {
	webhookd.WebhookTransformation
	chicken *chicken.Chicken
}

func NewChickenTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	lang := q.Get("language")
	str_clucking := q.Get("clucking")

	clucking, err := strconv.ParseBool(str_clucking)

	if err != nil {
		return nil, err
	}

	ch, err := chicken.GetChickenForLanguageTag(lang, clucking)

	if err != nil {
		return nil, err
	}

	tr := ChickenTransformation{
		chicken: ch,
	}

	return &tr, nil
}

func (tr *ChickenTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {

	select {
	case <-ctx.Done():
		return nil, nil
	default:
		// pass
	}

	txt := tr.chicken.TextToChicken(string(body))
	return []byte(txt), nil
}
