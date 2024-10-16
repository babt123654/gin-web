package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	baseURL         = "https://api.hiro.so/extended/v2/addresses/SP2XD7417HGPRTREMKF748VNEQPDRR0RMANB7X1NK.cross-peg-out-endpoint-v2-01/transactions?limit=50&offset="
	txDetailBaseURL = "https://api.hiro.so/extended/v1/tx/"
	outputFile      = "results.txt"
)

type TransactionResponse struct {
	Results []Transaction `json:"results"`
}

type Transaction struct {
	Tx struct {
		TxID string `json:"tx_id"`
	} `json:"tx"`
}

type TxDetailResponse struct {
	TxID          string   `json:"tx_id"`
	SenderAddress string   `json:"sender_address"`
	Events        []Events `json:"events"`
	ContractCall  struct {
		FunctionArgs []struct {
			Repr string `json:"repr"`
		} `json:"function_args"`
	} `json:"contract_call"`
}

type Events struct {
	EventIndex  int         `json:"event_index"`
	EventType   string      `json:"event_type"`
	TxID        string      `json:"tx_id"`
	Asset       Asset       `json:"asset,omitempty"`
	ContractLog ContractLog `json:"contract_log,omitempty"`
}

type Asset struct {
	AssetEventType string `json:"asset_event_type"`
	AssetID        string `json:"asset_id"`
	Sender         string `json:"sender"`
	Recipient      string `json:"recipient"`
	Amount         string `json:"amount"`
}

type ContractLog struct {
	ContractID string `json:"contract_id"`
	Topic      string `json:"topic"`
	Value      Value  `json:"value"`
}

type Value struct {
	Hex  string `json:"hex"`
	Repr string `json:"repr"`
}

func TestExplorer(t *testing.T) {
	// 打开输出文件
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer file.Close()

	offset := 0
	limit := 50
	totalRecords := 2253 // 你可以根据 API 返回的 total 字段动态获取

	for offset < totalRecords {
		// 请求第一个 API
		resp, err := http.Get(fmt.Sprintf("%s%d", baseURL, offset))
		if err != nil {
			fmt.Println("Error fetching transactions:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var txResponse TransactionResponse
		if err := json.Unmarshal(body, &txResponse); err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		// 遍历交易记录
		for _, transaction := range txResponse.Results {
			txID := transaction.Tx.TxID
			// 请求第二个 API
			time.Sleep(2 * time.Second) // 等待 2 秒
			txDetailResp, err := http.Get(txDetailBaseURL + txID + "?event_offset=0&event_limit=100&unanchored=true")
			if err != nil {
				fmt.Println("Error fetching transaction details:", err)
				continue
			}
			defer txDetailResp.Body.Close()

			detailBody, err := ioutil.ReadAll(txDetailResp.Body)
			if err != nil {
				fmt.Println("Error reading transaction detail response:", err)
				continue
			}

			var txDetail TxDetailResponse
			if err := json.Unmarshal(detailBody, &txDetail); err != nil {
				fmt.Println("Error unmarshalling transaction detail JSON:", err)
				continue
			}

			// 检查 repr 字段
			for _, event := range txDetail.Events {
				// 确保 ContractLog 存在
				if event.ContractLog.Value.Repr != "" {
					if strings.Contains(event.ContractLog.Value.Repr, "arb") || strings.Contains(event.ContractLog.Value.Repr, "BSC") || strings.Contains(event.ContractLog.Value.Repr, "Arbitrum") || strings.Contains(event.ContractLog.Value.Repr, "arbitrum") {
						// 记录 tx_id 和 sender_address
						_, err := file.WriteString(fmt.Sprintf("tx_id: %s, sender_address: %s, chain: %s\n", txDetail.TxID, txDetail.SenderAddress, event.ContractLog.Value.Repr))
						if err != nil {
							fmt.Println("Error writing to output file:", err)
						}
					}
				}
			}
		}

		// 更新 offset
		offset += limit
	}

	fmt.Println("Processing completed. Results saved to", outputFile)
}
