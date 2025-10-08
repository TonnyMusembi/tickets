package sms

import (
	"context"
	"fmt"
	"os"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSProvider interface {
	SendSMS(ctx context.Context, to, body string) error
}

type TwilioProvider struct {
	client *twilio.RestClient
	from   string
}

func NewTwilioProvider() *TwilioProvider {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	return &TwilioProvider{
		client: client,
		from:   os.Getenv("TWILIO_PHONE_NUMBER"),
	}
}

func (t *TwilioProvider) SendSMS(ctx context.Context, to, body string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(t.from)
	params.SetBody(body)

	resp, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send sms: %w", err)
	}

	fmt.Println("SMS sent, SID:", *resp.Sid)
	return nil
}
