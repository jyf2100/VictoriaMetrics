package notifier

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// AlertManager represents integration provider with Prometheus alert manager
// https://github.com/prometheus/alertmanager
type AlertManager struct {
	alertURL string
	argFunc  AlertURLGenerator
	client   *http.Client
}

// Send an alert or resolve message
func (am *AlertManager) Send(alerts []Alert) error {
	b := &bytes.Buffer{}
	writeamRequest(b, alerts, am.argFunc)
	resp, err := am.client.Post(am.alertURL, "application/json", b)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response from %q: %s", am.alertURL, err)
		}
		return fmt.Errorf("invalid SC %d from %q; response body: %s", resp.StatusCode, am.alertURL, string(body))
	}
	return nil
}

// AlertURLGenerator returns URL to single alert by given name
type AlertURLGenerator func(group, id string) string

const alertManagerPath = "/api/v2/alerts"

// NewAlertManager is a constructor for AlertManager
func NewAlertManager(alertManagerURL string, fn AlertURLGenerator, c *http.Client) *AlertManager {
	return &AlertManager{
		alertURL: strings.TrimSuffix(alertManagerURL, "/") + alertManagerPath,
		argFunc:  fn,
		client:   c,
	}
}
