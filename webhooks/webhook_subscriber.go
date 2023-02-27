package webhooks

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type WebhookSubscriber struct {
	endpoint string
}

func (w WebhookSubscriber) Notify(topic string, object any) error {
	endpoint := w.endpoint

	jsonData, err := json.Marshal(map[string]interface{}{
		"event":  topic,
		"object": object,
	})
	
	if err != nil {
		return err
	}

	// todo: wrap to a new package for sending requests
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// todo: check returned status code and implement retry system
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	_ = response.Body.Close()
	
	return nil
}
