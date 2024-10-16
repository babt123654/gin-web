package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gin-web/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 创建数 token
func (my MysqlService) CreateVelarToken(r *models.Message) (err error) {
	body, err := my.GetVelarToken("https://mainnet-prod-proxy-service-dedfb0daae85.herokuapp.com/swapapp/swap/tokens")
	fmt.Println(string(body))
	_, value, err := GetFieldValueFromJSON(body, "message")
	var m []models.Message
	json.Unmarshal(value, &m)
	for _, asset := range m {
		fmt.Println("Order:", asset.Symbol+":"+asset.Price)
	}
	my.Q.Tx.Updates(&m)
	return
}
func GetFieldValueFromJSON(data []byte, field string) (string, []byte, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return "", nil, err
	}
	fieldValue, ok := jsonData[field]
	if !ok {
		return "", nil, fmt.Errorf("field '%s' not found in JSON", field)
	}
	fieldValueBytes, err := json.Marshal(fieldValue)
	if err != nil {
		return "", nil, err
	}

	return field, fieldValueBytes, nil
}

// 交易所价格
func (my MysqlService) GetPriceFromOkx() (string, error) {
	url := "https://www.okx.com/api/v5/public/instruments?instType=SPOT&instId=STX-USDT"

	// 创建请求体
	payload := strings.NewReader("")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

//	func GetFieldValueFromJSON(data []byte, field string) (string, interface{}, error) {
//		var jsonData map[string]interface{}
//		err := json.Unmarshal(data, &jsonData)
//		if err != nil {
//			return "", nil, err
//		}
//
//		fieldValue, ok := jsonData[field]
//		if !ok {
//			return "", nil, fmt.Errorf("field '%s' not found in JSON", field)
//		}
//
//		return field, fieldValue, nil
//	}

func GetNodeFromJSON(data []byte, node string) ([]byte, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	node, ok := jsonData[node].(string)
	if !ok {
		return nil, fmt.Errorf("message node not found or not a string")
	}

	nodeData := []byte(node)
	return nodeData, nil
}

// 前端页面监控
func (my MysqlService) GetVelarToken(address string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://mainnet-prod-proxy-service-dedfb0daae85.herokuapp.com/swapapp/swap/tokens"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	return body, nil
}

// 登录 timestamp, method, requestPath, body, secretKey string
func (my MysqlService) SendPrivateRequest(apiKey, secretKey, passphrase, method, requestPath, requestBody string) (string, error) {
	// 生成时间戳
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	// 生成签名
	signature := generateSignature(timestamp, method, requestPath, requestBody, secretKey)
	// 创建请求
	client := &http.Client{}
	url := "https://okx.com" + requestPath
	req, err := http.NewRequest(method, url, strings.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	// 设置请求头
	req.Header.Set("OK-ACCESS-KEY", apiKey)
	req.Header.Set("OK-ACCESS-SIGN", signature)
	req.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("OK-ACCESS-PASSPHRASE", passphrase)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
func generateSignature(timestamp, method, requestPath, body, secretKey string) string {
	data := timestamp + method + requestPath + body
	key := []byte(secretKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature
}
