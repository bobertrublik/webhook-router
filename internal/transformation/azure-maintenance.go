package transformation

import (
	"context"
	"encoding/json"
	"github.com/bobertrublik/webhook-router/internal/webhookd"
	"github.com/tidwall/gjson"
)

func init() {

	ctx := context.Background()
	err := RegisterTransformation(ctx, "azure-maintenance", NewAzureMaintenanceTransformation)

	if err != nil {
		panic(err)
	}
}

type Schema struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type     string    `json:"type"`
	Text     *Text     `json:"text,omitempty"`
	Elements []Element `json:"elements,omitempty"`
	Fields   []Field   `json:"fields,omitempty"`
}

type Text struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji,omitempty"`
}

type Element struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Field struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// AzureMaintenanceTransformation implements the `webhookd.WebhookTransformation` interface for a no-op transformation meaning
// the output of the `Transform` method is the same as its input.
type AzureMaintenanceTransformation struct {
	webhookd.WebhookTransformation
}

// NewInsecureTransformation returns a new `AzureMaintenanceTransformation` instance configured by 'uri' in the form of:
//
//	pass://
func NewAzureMaintenanceTransformation(ctx context.Context, uri string) (webhookd.WebhookTransformation, error) {

	p := AzureMaintenanceTransformation{}
	return &p, nil
}

// Transform returns 'body' unaltered.
func (p *AzureMaintenanceTransformation) Transform(ctx context.Context, body []byte) ([]byte, *webhookd.WebhookError) {
	myString := string(body[:])
	alertRule := gjson.Get(myString, "data.essentials.alertRule")
	description := gjson.Get(myString, "data.essentials.description")
	start := gjson.Get(myString, "data.alertContext.properties.impactStartTime")
	end := gjson.Get(myString, "data.alertContext.properties.impactMitigationTime")
	stage := gjson.Get(myString, "data.alertContext.properties.stage")
	communication := gjson.Get(myString, "data.alertContext.properties.communication")
	status := gjson.Get(myString, "data.alertContext.status")

	emoji := true
	schema := Schema{
		Blocks: []Block{
			{
				Type: "header",
				Text: &Text{
					Type:  "plain_text",
					Text:  "Service Health Alert",
					Emoji: emoji,
				},
			},
			{
				Type: "context",
				Elements: []Element{
					{
						Type: "mrkdwn",
						Text: alertRule.String(),
					},
				},
			},
			{
				Type: "section",
				Fields: []Field{
					{
						Type: "mrkdwn",
						Text: start.String(),
					},
					{
						Type: "mrkdwn",
						Text: stage.String(),
					},
					{
						Type: "mrkdwn",
						Text: end.String(),
					},
					{
						Type: "mrkdwn",
						Text: status.String(),
					},
				},
			},
			{
				Type: "section",
				Text: &Text{
					Type: "mrkdwn",
					Text: description.String(),
				},
			},
			{
				Type: "section",
				Text: &Text{
					Type: "mrkdwn",
					Text: communication.String(),
				},
			},
		},
	}
	jsonData, err := json.MarshalIndent(schema, "", "    ")
	if err != nil {
		panic(err)
	}

	return jsonData, nil
}
