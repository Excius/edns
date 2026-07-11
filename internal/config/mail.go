package config

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

func NewMailClient(cfg *Config) (*mail.Client, error) {
	var mailClient *mail.Client
	var err error

	if cfg.App.Env == "development" {
		mailClient, err = mail.NewClient(
			cfg.SMTP.Host,
			mail.WithPort(cfg.SMTP.Port),
			mail.WithTLSPolicy(mail.NoTLS),
		)
	} else {
		mailClient, err = mail.NewClient(
			cfg.SMTP.Host,
			mail.WithPort(cfg.SMTP.Port),
			mail.WithTLSPolicy(mail.TLSMandatory),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("new client creation failed with: %w", err)
	}

	return mailClient, nil
}
