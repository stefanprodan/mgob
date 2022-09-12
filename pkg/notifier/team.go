package notifier

import (
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/pkg/config"
)

func sendTeamNotification(subject string, body string, warn bool, config *config.Team) error {
	if !warn && config.WarnOnly {
		return nil
	}
	// init the client
	mstClient := goteamsnotify.NewClient()

	// setup webhook url
	webhookUrl := config.WebhookURL

	// setup message card
	msgCard := goteamsnotify.NewMessageCard()
	msgCard.Title = subject
	msgCard.Text = body
	msgCard.ThemeColor = config.ThemeColor
	// send
	if err := mstClient.Send(webhookUrl, msgCard); err != nil {
		return errors.Wrapf(err, "sending Team notification failed")
	}

	return nil
}
