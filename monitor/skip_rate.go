package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Chainflow/solana-mission-control/alerter"
	"github.com/Chainflow/solana-mission-control/config"
	"github.com/Chainflow/solana-mission-control/types"
)

var (
	solanaBinaryPath = os.Getenv("SOLANA_BINARY_PATH")
	cluster          = os.Getenv("CLUSTER")
)

func SkipRate(cfg *config.Config) (float64, float64, error) {
	var valSkipped, netSkipped, totalSkipped float64
	var validatorCount int64

	if solanaBinaryPath == "" {
		solanaBinaryPath = "solana"
	}

	clusterFlag := "-um"
	if cluster == "testnet" {
		clusterFlag = "-ut"
	}

	log.Printf("Solana binary path (skip rate) : %s (%s)", solanaBinaryPath, cluster)

	cmd := exec.Command(solanaBinaryPath, "validators", "--output", "json", clusterFlag)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while running solana validators cli command %v", err)
		return valSkipped, netSkipped, err
	}

	var result types.SkipRate
	err = json.Unmarshal(out, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return valSkipped, netSkipped, err
	}

	for _, val := range result.Validators {
		if val.IdentityPubkey == cfg.ValDetails.PubKey {
			valSkipped = val.SkipRate
		}
		totalSkipped = totalSkipped + val.SkipRate
		validatorCount++
	}

	// Calculate network skip rate directly from CLI data instead of making RPC call
	if validatorCount > 0 {
		netSkipped = totalSkipped / float64(validatorCount)
	}

	log.Printf("VAL skip rate : %f, Network skip rate : %f", valSkipped, netSkipped)

	return valSkipped, netSkipped, nil
}

func SkipRateAlerts(cfg *config.Config) error {
	valSkipped, netSkipped, err := SkipRate(cfg)
	if err != nil {
		log.Printf("Error while getting SkipRate: %v", err)
	}

	if valSkipped > netSkipped && (valSkipped > float64(cfg.AlertingThresholds.SkipRateThreshold)) {
		if strings.EqualFold(cfg.AlerterPreferences.SkipRateAlerts, "yes") {
			err = alerter.SendTelegramAlert(fmt.Sprintf("SKIP RATE ALERT ::  Your validator SKIP RATE : %f has exceeded network SKIP RATE : %f", valSkipped, netSkipped), cfg)
			if err != nil {
				log.Printf("Error while sending skip rate alert to telegram: %v", err)
			}
			err = alerter.SendEmailAlert(fmt.Sprintf("Your validator SKIP RATE : %f has exceeded network SKIP RATE : %f", valSkipped, netSkipped), cfg)
			if err != nil {
				log.Printf("Error while sending skip rate alert to email: %v", err)
			}
			err = alerter.SendSlackAlert(fmt.Sprintf("SKIP RATE ALERT ::  Your validator SKIP RATE : %f has exceeded network SKIP RATE : %f", valSkipped, netSkipped), cfg)
			if err != nil {
				log.Printf("Error while sending skip rate alert to slack: %v", err)
			}
		}
	}
	return nil
}
