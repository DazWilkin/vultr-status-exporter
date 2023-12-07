package api

import (
	"log/slog"
	"testing"
)

func TestAlerts(t *testing.T) {
	client := NewClient()

	alerts, err := client.Alerts()
	if err != nil {
		t.Fatal("expected success")
	}

	slog.Info("Result",
		"alerts", alerts,
	)
}
func TestStatus(t *testing.T) {
	client := NewClient()

	status, err := client.Status()
	if err != nil {
		t.Fatal("expected success")
	}

	slog.Info("Result",
		"status", status,
	)
}
