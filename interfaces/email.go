package interfaces

import "context"

type IEmailSender interface {
	SendEmail(ctx context.Context, from string, to string, subject string, content string) error
	SendEmailWithTemplate(
		ctx context.Context,
		from string,
		to string,
		templateUUID string,
		templateVars map[string]interface{},
	) error
}
