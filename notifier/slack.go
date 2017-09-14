package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/config"
)

type slackPayload struct {
	Channel     string            `json:"channel"`
	Username    string            `json:"username"`
	IconUrl     string            `json:"icon_url"`
	IconEmoji   string            `json:"icon_emoji"`
	Text        string            `json:"text,omitempty"`
	Attachments []slackAttachment `json:"attachments,omitempty"`
}

type slackAttachment struct {
	Color    string   `json:"color"`
	Title    string   `json:"title"`
	Pretext  string   `json:"pretext"`
	Text     string   `json:"text"`
	MrkdwnIn []string `json:"mrkdwn_in"`
}

func sendSlackNotification(subject string, body string, warn bool, cfg *config.Slack) error {
	if !warn && cfg.WarnOnly {
		return nil
	}

	payload := slackPayload{
		Channel:  cfg.Channel,
		Username: cfg.Username,
	}

	var emoji, color string
	if warn {
		emoji = ":x:"
		color = "danger"
	} else {
		emoji = ":white_check_mark:"
		color = "good"
	}

	title := "backup log"
	pretext := fmt.Sprintf("%s *%s*", emoji, subject)

	a := slackAttachment{
		Color:    color,
		Title:    title,
		Pretext:  pretext,
		Text:     body,
		MrkdwnIn: []string{"text", "pretext"},
	}

	payload.Attachments = []slackAttachment{a}

	data, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrapf(err, "Marshalling slack payload failed")
	}

	b := bytes.NewBuffer(data)

	if res, err := http.Post(cfg.URL, "application/json", b); err != nil {
		return errors.Wrapf(err, "Sending data to slack failed")
	} else {
		defer res.Body.Close()
		statusCode := res.StatusCode
		if statusCode != 200 {
			body, _ := ioutil.ReadAll(res.Body)
			return errors.Errorf("Sending data to slack failed %v", string(body))
		}
	}

	return nil
}
