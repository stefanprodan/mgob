package notifier

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/pkg/errors"

	"github.com/stefanprodan/mgob/pkg/config"
)

func sendEmailNotification(subject string, body string, config *config.SMTP) error {

	msg := "From: \"MGOB\" <" + config.From + ">\r\n" +
		"To: " + strings.Join(config.To, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n"

	addr := fmt.Sprintf("%v:%v", config.Server, config.Port)

	// auth is set to nil by default
	// workaround for error given if auth is disabled on the smtp server
	// notifier error: "smtp: server doesn't support AUTH"
	var auth smtp.Auth
	if config.Username != "" {
		auth = smtp.PlainAuth("", config.Username, config.Password, config.Server)
	}

	if err := smtp.SendMail(addr, auth, config.From, config.To, []byte(msg)); err != nil {
		return errors.Wrapf(err, "sending email notification failed")
	}
	return nil
}
