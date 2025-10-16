package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Chainflow/solana-mission-control/alerter"
	"github.com/Chainflow/solana-mission-control/config"
	"github.com/Chainflow/solana-mission-control/querier"
	"github.com/Chainflow/solana-mission-control/types"
)

// GetIdentityBalance returns the balance of the identity account
func GetIdentityBalance(cfg *config.Config) (types.Balance, error) {
	log.Println("Getting identity account Balance...")
	return GetAccountBalance(cfg, cfg.ValDetails.PubKey)
}

// GetVoteAccBalance returns the balance of the vote account
func GetVoteAccBalance(cfg *config.Config) (types.Balance, error) {
	log.Println("Getting vote Aaccount Balance...")
	return GetAccountBalance(cfg, cfg.ValDetails.VoteKey)
}

// GetAccountBalance returns the balance for an arbitrary base58 account address
func GetAccountBalance(cfg *config.Config, address string) (types.Balance, error) {
	ops := types.HTTPOptions{
		Endpoint: cfg.Endpoints.RPCEndpoint,
		Method:   http.MethodPost,
		Body: types.Payload{Jsonrpc: "2.0", Method: "getBalance", ID: 1, Params: []interface{}{
			address,
		}},
	}

	var result types.Balance
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	return result, nil
}

// SendBalanceChangeAlert checks balance and DBbalance, If balance dropped to threshold,
// sends Alerts to the validator
func SendBalanceChangeAlert(currentBal int64, cfg *config.Config) error {
	prevBal, err := querier.GetAccountBalFromDB(cfg)
	if err != nil {
		log.Printf("Error while getting bal from db : %v", err)
		return err
	}

	// c := float64(currentBal) / math.Pow(10, 9)
	c := fmt.Sprintf("%.4f", float64(currentBal)/math.Pow(10, 9))
	cBal, _ := strconv.ParseFloat(c, 64)
	current := c + "SOL"
	previous := prevBal + "SOL"

	if strings.EqualFold(cfg.AlerterPreferences.AccountBalanceChangeAlerts, "yes") {
		if cBal < cfg.AlertingThresholds.BalanaceChangeThreshold {
			err = alerter.SendTelegramAlert(fmt.Sprintf("Account Balance Alert: Your account balance has dropped below configured threshold, current balance is : %s", current), cfg)
			if err != nil {
				log.Printf("Error while sending account balance change alert to telegram : %v", err)
				return err
			}

			err = alerter.SendEmailAlert(fmt.Sprintf("Account Balance Alert: Your account balance has dropped below configured threshold, current balance is : %s", current), cfg)
			if err != nil {
				log.Printf("Error while sending account balance change alert to email: %v", err)
				return err
			}

			err = alerter.SendSlackAlert(fmt.Sprintf("Account Balance Alert: Your account balance has dropped below configured threshold, current balance is : %s", current), cfg)
			if err != nil {
				log.Printf("Error while sending account balance change alert to slack: %v", err)
				return err
			}
		}
	}

	// Send delegation alerts
	if prevBal != "" {
		pBal, err := strconv.ParseFloat(prevBal, 64)
		if err != nil {
			log.Printf("Error while converting pBal to float64 : %v ", err)
			return err
		}

		if strings.EqualFold(cfg.AlerterPreferences.DelegationAlerts, "yes") {
			diff := cBal - pBal
			if diff > 50 && diff < 100 { // check and change the condition
				// Alert to telegram
				err = alerter.SendTelegramAlert(fmt.Sprintf("Delegation Alert: Your account balance has changed form %s to %s", previous, current), cfg)
				if err != nil {
					log.Printf("Error while sending delegation alert to telegram : %v", err)
					return err
				}

				// Alert to email
				err = alerter.SendEmailAlert(fmt.Sprintf("Delegation Alert: Your account balance has changed form %s to %s", previous, current), cfg)
				if err != nil {
					log.Printf("Error while sending delegation alert to email : %v", err)
					return err
				}

				// Alert to slack
				err = alerter.SendSlackAlert(fmt.Sprintf("Delegation Alert: Your account balance has changed form %s to %s", previous, current), cfg)
				if err != nil {
					log.Printf("Error while sending delegation alert to slack : %v", err)
					return err
				}
			} else if diff < -50 { // check and change the condition
				// Alert to telegram
				err = alerter.SendTelegramAlert(fmt.Sprintf("Undelegation Alert: Your account balance has changed form %s to %s", previous, current), cfg)
				if err != nil {
					log.Printf("Error while sending undelegation alert to telegram : %v", err)
					return err
				}

				// Alert to email
				err = alerter.SendEmailAlert(fmt.Sprintf("Undelegation Alert: Your account balance has changed form %s to %s", previous, current), cfg)
				if err != nil {
					log.Printf("Error while sending undelegation alert to email : %v", err)
					return err
				}

				// Alert to slack
				err = alerter.SendSlackAlert(fmt.Sprintf("Undelegation Alert: Your account balance has changed form %s to %s", previous, current), cfg)
				if err != nil {
					log.Printf("Error while sending undelegation alert to slack : %v", err)
					return err
				}
			}
		}
	}

	return nil
}
