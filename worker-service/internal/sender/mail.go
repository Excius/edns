package sender

import (
	"context"
	"fmt"
	"html/template"

	"github.com/excius/edns/internal/events"
	"github.com/excius/edns/internal/logger"
	"github.com/wneessen/go-mail"
	"go.uber.org/zap"
)

type MailData struct {
	Message        string
	NotificationID string
}

type EmailSender struct {
	client      *mail.Client
	fromAddress string
}

func NewEmailSender(client *mail.Client, from string) *EmailSender {
	return &EmailSender{
		client:      client,
		fromAddress: from,
	}
}

func (e *EmailSender) Send(ctx context.Context, event events.NotificationEvent) error {

	mailData := MailData{
		Message:        event.Message,
		NotificationID: event.NotificationID,
	}

	tmpl, err := template.ParseFiles("templates/notification.html")
	if err != nil {
		return fmt.Errorf("failed to parse template file: %w", err)
	}

	m := mail.NewMsg()
	if err := m.From(e.fromAddress); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}

	if err := m.To(event.Email); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}

	m.Subject(event.Title)
	m.SetBodyHTMLTemplate(tmpl, mailData)

	if err := e.client.DialAndSendWithContext(ctx, m); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	logger.Log.Info(
		"Email sent",
		zap.String("notification_id", event.NotificationID),
	)

	return nil
}
