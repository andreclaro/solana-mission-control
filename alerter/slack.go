package alerter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *slackAlert) SendSlackMessage(msg, webhookURL string) error {
	data := map[string]string{
		"text": msg,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send slack message, status code: %d", resp.StatusCode)
	}

	return nil
}
