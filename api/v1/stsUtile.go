package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	v1 "gin-web/api/v1/wallet"
	"gin-web/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getAPI1(address string) (*models.BalanceAPI1, error) {
	url := fmt.Sprintf("https://stacks-node.alexlab.co/v2/accounts/%s?proof=0", address)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var balanceAPI1 models.BalanceAPI1
	err = json.Unmarshal(body, &balanceAPI1)
	if err != nil {
		return nil, err
	}

	return &balanceAPI1, nil
}

func getAPI2(address string) (*models.BalanceAPI2, error) {
	url := fmt.Sprintf("https://api.hiro.so/extended/v1/address/%s/balances", address)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var balanceAPI2 models.BalanceAPI2
	err = json.Unmarshal(body, &balanceAPI2)
	if err != nil {
		return nil, err
	}

	return &balanceAPI2, nil
}

// retryOperation 封装重试逻辑
func retryOperation(maxRetries int, operation func() error) error {
	for i := 0; i < maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}
		time.Sleep(time.Second) // 等待一段时间再重试
	}
	return fmt.Errorf("operation failed after %d attempts", maxRetries)
}

// GetCombinedBalance 获取组合余额
func GetCombinedBalance(address string) (*models.CombinedBalance, error) {
	var balance1 *models.BalanceAPI1
	var balance2 *models.BalanceAPI2
	var err error

	// 封装 API 调用为操作
	err = retryOperation(3, func() error {
		balance1, err = getAPI1(address)
		CheckErr(err)
		return err
	})
	if err != nil {
		fmt.Printf("Failed to get balance1 after retries: %v\n", err)
	}

	err = retryOperation(3, func() error {
		balance2, err = getAPI2(address)
		CheckErr(err)
		return err
	})
	if err != nil {
		fmt.Printf("Failed to get balance2 after retries: %v\n", err)
	}

	api := "https://docs-demo.stacks-mainnet.quiknode.pro/extended/v1/address/" + address + "/nonces"
	nonceInfo := models.NonceInfo{}

	err = retryOperation(3, func() error {
		return v1.ReqAPIwhithStructAndBro(api, http.MethodGet, &nonceInfo)
	})
	if err != nil {
		fmt.Printf("Failed to get nonce info after retries: %v\n", err)
	}

	// 检查 balance1 和 balance2 是否为 nil
	if balance1 == nil || balance2 == nil {
		return nil, fmt.Errorf("one of the balances is nil for address: %s", address)
	}

	return &models.CombinedBalance{
		API1:      *balance1,
		API2:      *balance2,
		NonceInfo: nonceInfo,
	}, nil
}
func RequestAPIwithStruct1(url string, struc interface{}) (*interface{}, error) {
	resp, err := http.Get(url)
	CheckErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	err = json.Unmarshal(body, &struc)
	CheckErr(err)
	return &struc, nil
}

// RequestAPIwithStruct 从指定 API 获取数据并填充到给定的结构体中
func RequestAPIwithStruct(apiURL string, method string, result interface{}) (*interface{}, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	var req *http.Request
	var err error
	if method == http.MethodPost {
		req, err = http.NewRequest(method, apiURL, bytes.NewBuffer(nil))
	} else {
		req, err = http.NewRequest(method, apiURL, nil)
	}
	CheckErr(err)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("意外的状态码: %v", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	return &result, nil
}

// 获取出块儿后TX
func GetAfterBlockTX(pendingTx models.PendingTx) []models.PendingTxResults {
	_, burnTime, _ := GetBurnBlockTime()
	//fmt.Println("出块儿时间：" + fmt.Sprintf("%d%s%s%s%d", burnTime, "=", time.Unix(burnTime, 0).Format("2006-01-02 15:04:05"), ":", height))
	afterBurnTime := make([]models.PendingTxResults, 0)
	// 遍历每条结果
	for _, result1 := range pendingTx.PendingTxResults {
		// 确认当前交易产生在出块后
		if result1.ReceiptTime > burnTime {
			afterBurnTime = append(afterBurnTime, result1)
		}
	}
	return afterBurnTime
}

// 根据交易token，获取手续费最高的TX
func GetMaxFeeTXwithTokens(repr, name, account string, pendinTxs []models.PendingTxResults) models.PendingTxResults {
	for _, pendingtx := range pendinTxs {
		if strings.Contains(pendingtx.ContractCall.FunctionName, "swap") && pendingtx.SenderAddress != account {
			for _, arg := range pendingtx.ContractCall.FunctionArgs {
				if arg.Repr == repr && arg.Name == name {
					return pendingtx
				}
			}
		}
	}
	return models.PendingTxResults{}
}

////使用aip获取在出块儿后为提交的交易
//func GetAfterBlockTXwithAPI(apiURL string,result interface{}) ([]models.PendingTxResults) {
//	RequestAPIwithStruct(apiURL, http.MethodGet, &result)
//	GetAfterBlockTX(&result)
//	return afterBlockTx
//}
