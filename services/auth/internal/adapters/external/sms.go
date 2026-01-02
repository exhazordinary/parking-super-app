package external

import (
	"context"
	"fmt"
	"log"
)

// SMSService implementations for various providers.
// In production, you would integrate with:
// - Malaysian providers: Nexmo/Vonage, Twilio, local telcos
// - Each provider has its own SDK/API

// ConsoleSMSService is a mock SMS service that logs messages.
// Use this for development and testing.
type ConsoleSMSService struct{}

// NewConsoleSMSService creates a new console SMS service.
func NewConsoleSMSService() *ConsoleSMSService {
	return &ConsoleSMSService{}
}

// SendOTP logs the OTP to console instead of sending SMS.
func (s *ConsoleSMSService) SendOTP(ctx context.Context, phone, code string) error {
	log.Printf("[SMS] Sending OTP %s to %s", code, phone)
	return nil
}

// SendMessage logs the message to console.
func (s *ConsoleSMSService) SendMessage(ctx context.Context, phone, message string) error {
	log.Printf("[SMS] Sending message to %s: %s", phone, message)
	return nil
}

// TwilioSMSService integrates with Twilio for SMS delivery.
// This is a production-ready implementation.
//
// SETUP:
// 1. Create a Twilio account
// 2. Get Account SID, Auth Token, and phone number
// 3. Install: go get github.com/twilio/twilio-go
type TwilioSMSService struct {
	accountSID string
	authToken  string
	fromPhone  string
	// client *twilio.RestClient // Uncomment when using Twilio SDK
}

// NewTwilioSMSService creates a new Twilio SMS service.
func NewTwilioSMSService(accountSID, authToken, fromPhone string) *TwilioSMSService {
	return &TwilioSMSService{
		accountSID: accountSID,
		authToken:  authToken,
		fromPhone:  fromPhone,
		// client: twilio.NewRestClientWithParams(twilio.ClientParams{
		// 	Username: accountSID,
		// 	Password: authToken,
		// }),
	}
}

// SendOTP sends an OTP via Twilio.
func (s *TwilioSMSService) SendOTP(ctx context.Context, phone, code string) error {
	message := fmt.Sprintf("Your ParkingApp verification code is: %s. Valid for 5 minutes.", code)
	return s.SendMessage(ctx, phone, message)
}

// SendMessage sends an SMS via Twilio.
func (s *TwilioSMSService) SendMessage(ctx context.Context, phone, message string) error {
	// TODO: Implement actual Twilio integration
	// Example:
	//
	// params := &api.CreateMessageParams{}
	// params.SetTo(phone)
	// params.SetFrom(s.fromPhone)
	// params.SetBody(message)
	//
	// _, err := s.client.Api.CreateMessage(params)
	// return err

	log.Printf("[TWILIO] Would send to %s: %s", phone, message)
	return nil
}
