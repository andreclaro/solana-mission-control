package monitor_test

import (
	"testing"

	"github.com/Chainflow/solana-mission-control/config"
	"github.com/Chainflow/solana-mission-control/monitor"
)

func TestGetIdentity(t *testing.T) {
	cfg, err := config.ReadFromFile()
	if err != nil {
		t.Error("Error while reading config:", err)
	}

	res, err := monitor.GetIdentity(cfg)
	if err != nil {
		t.Error("Error while fetching identity:", err)
	}

	if res.Result.Identity == "" {
		t.Error("Expected non-empty identity, but got empty result")
	}

	if res.Result.Identity != "" {
		t.Log("Got Identity:", res.Result.Identity)
	}
}
