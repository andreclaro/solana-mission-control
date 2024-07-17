package alerter

// Telegram to send telegram alert interface
type Telegram interface {
	SendTelegramMessage(msgText, botToken string, chatID int64) error
}

type telegramAlert struct{}

// NewTelegramAlerter returns telegram alerter
func NewTelegramAlerter() *telegramAlert {
	return &telegramAlert{}
}

// Email to send mail alert
type Email interface {
	SendEmail(msg, token, toEmail string) error
}

type emailAlert struct{}

// NewEmailAlerter returns email alert
func NewEmailAlerter() *emailAlert {
	return &emailAlert{}
}

// Slack to send slack alert
type Slack interface {
	SendSlackMessage(msg, webhookURL string) error
}

type slackAlert struct{}

func NewSlackAlerter() *slackAlert {
	return &slackAlert{}
}