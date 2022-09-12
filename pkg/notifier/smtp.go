package notifier

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/pkg/config"
)

func sendEmailNotification(subject string, body string, warn bool, config *config.SMTP) error {
	if !warn && config.WarnOnly {
		return nil
	}

	mail := &email.Email{
		To:      config.To,
		From:    config.From,
		Subject: subject,
		Text:    []byte(body),
	}

	if err := sendEmail(config, mail); err != nil {
		return errors.Wrapf(err, "sending email notification failed")
	}

	return nil
}

func sendEmail(config *config.SMTP, e *email.Email) error {
	// auth is set to nil by default
	// workaround for error given if auth is disabled on the smtp server
	// notifier error: "smtp: server doesn't support AUTH"
	var auth smtp.Auth
	if config.Username != "" {
		auth = smtp.PlainAuth("", config.Username, config.Password, config.Server)
	}

	addr := fmt.Sprintf("%v:%v", config.Server, config.Port)
	if config.TlsEnabled {
		config := &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify, ServerName: config.Server}
		return e.SendWithTLS(addr, auth, config)
	} else {
		return e.Send(addr, auth)
	}
}
