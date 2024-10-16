package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Transaction struct {
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
	Fee       int    `json:"fee"`
	// 其他必要字段
}

func signTransaction(privateKey string, tx Transaction) (string, error) {
	// 使用私钥对交易进行签名的逻辑
	// 返回签名后的交易数据
	return "signed_transaction_data", nil
}

func sendTransaction(signedTx string) error {
	url := "https://stacks-node-api.mainnet.stacks.co/v2/transactions"
	payload := map[string]string{"transaction": signedTx}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send transaction: %s", resp.Status)
	}
	return nil
}

func SendTransaction() {
	privateKey := "你的私钥"
	tx := Transaction{
		Recipient: "接收地址",
		Amount:    1000, // 发送的STX数量
		Fee:       200,  // 交易费用
	}

	signedTx, err := signTransaction(privateKey, tx)
	if err != nil {
		log.Fatalf("Error signing transaction: %v", err)
	}

	err = sendTransaction(signedTx)
	if err != nil {
		log.Fatalf("Error sending transaction: %v", err)
	}

	fmt.Println("Transaction sent successfully!")
}
