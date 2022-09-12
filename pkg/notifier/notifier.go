package notifier

import "github.com/stefanprodan/mgob/pkg/config"

func SendNotification(subject string, body string, warn bool, plan config.Plan) error {

	var err error
	if plan.SMTP != nil {
		err = sendEmailNotification(subject, body, warn, plan.SMTP)
	}
	if plan.Slack != nil {
		err = sendSlackNotification(subject, body, warn, plan.Slack)
	}
	return err
}
