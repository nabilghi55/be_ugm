package notification

import (
	"context"
	"log"
)

type NotificationMessage struct {
	RecipientPhone string
	RecipientEmail string
	CustomerName   string
	PawnNumber     string
	DueDate        string
	AmountDue      int64
	DaysOffset     int // H-7, H-3, H-1, Due Date, H+1, H+3, H+7
}

type NotificationProvider interface {
	Name() string
	SendWhatsApp(ctx context.Context, msg NotificationMessage) error
	SendEmail(ctx context.Context, msg NotificationMessage) error
	SendSMS(ctx context.Context, msg NotificationMessage) error
}

type MockNotificationProvider struct{}

func NewMockNotificationProvider() *MockNotificationProvider {
	return &MockNotificationProvider{}
}

func (n *MockNotificationProvider) Name() string {
	return "Mock Omnichannel Notification Provider"
}

func (n *MockNotificationProvider) SendWhatsApp(ctx context.Context, msg NotificationMessage) error {
	log.Printf("[NOTIFICATION WA] Sent reminder to %s (%s) for Pawn %s (Due: %s, Amount: Rp %d)",
		msg.CustomerName, msg.RecipientPhone, msg.PawnNumber, msg.DueDate, msg.AmountDue)
	return nil
}

func (n *MockNotificationProvider) SendEmail(ctx context.Context, msg NotificationMessage) error {
	log.Printf("[NOTIFICATION EMAIL] Sent email to %s (%s) for Pawn %s",
		msg.CustomerName, msg.RecipientEmail, msg.PawnNumber)
	return nil
}

func (n *MockNotificationProvider) SendSMS(ctx context.Context, msg NotificationMessage) error {
	log.Printf("[NOTIFICATION SMS] Sent SMS to %s (%s)", msg.CustomerName, msg.RecipientPhone)
	return nil
}

type NotificationService struct {
	providers []NotificationProvider
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		providers: []NotificationProvider{NewMockNotificationProvider()},
	}
}

func (s *NotificationService) TriggerDueReminder(ctx context.Context, msg NotificationMessage) {
	for _, p := range s.providers {
		if msg.RecipientPhone != "" {
			_ = p.SendWhatsApp(ctx, msg)
		}
		if msg.RecipientEmail != "" {
			_ = p.SendEmail(ctx, msg)
		}
	}
}
