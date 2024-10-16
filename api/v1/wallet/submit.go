package v1

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

// BroadcastTransaction 广播交易到 Stacks 节点
func BroadcastTransaction(serializedTx string) error {
	// Stacks 节点 API 的 URL（主网）
	apiURLs := []string{
		"https://stacks-node-api.mainnet.stacks.co/v2/transactions",
		"https://api.hiro.so/v2/transactions",
	}
	// 将十六进制字符串转换为字节数组
	txBytes, err := hex.DecodeString(serializedTx)
	if err != nil {
		return fmt.Errorf("无法解码交易: %v", err)
	}
	// 尝试使用每个 API URL 进行请求
	for _, apiURL := range apiURLs {
		// 创建一个新的请求
		req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(txBytes))
		if err != nil {
			return fmt.Errorf("创建请求失败: %v", err)
		}
		// 设置请求头
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-API-Key", "001f2edce38260813b718a1350125c94")
		// 发送 HTTP 请求，将交易广播到 Stacks 节点
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("无法广播交易到 %s: %v\n", apiURL, err)
			continue // 如果请求失败，尝试下一个 API URL
		}
		defer resp.Body.Close()
		// 读取响应体
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("无法读取响应体: %v", err)
		}
		// 检查响应状态码
		if resp.StatusCode == http.StatusOK {
			fmt.Println("交易成功广播到区块链网络")
			return nil // 成功广播，返回
		}
		fmt.Printf("广播交易失败到 %s, 状态码: %d, 响应: %s\n", apiURL, resp.StatusCode, body)
	}
	return fmt.Errorf("所有 API 请求均失败，无法广播交易")
}

// broadcastTransaction 广播交易到Stacks节点
func BroadcastTransaction1(serializedTx string) error {
	// Stacks节点API的URL（此处为主网URL）
	apiURL := "https://stacks-node-api.mainnet.stacks.co/v2/transactions"
	//apiURL:=

	// 将十六进制字符串转换为字节数组
	txBytes, err := hex.DecodeString(serializedTx)
	if err != nil {
		return fmt.Errorf("无法解码交易: %v", err)
	}

	// 发送HTTP请求，将交易广播到Stacks节点
	resp, err := http.Post(apiURL, "application/octet-stream", bytes.NewBuffer(txBytes))
	if err != nil {
		return fmt.Errorf("无法广播交易: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("无法读取响应体: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("广播交易失败, 状态码: %d, 响应: %s", resp.StatusCode, body)
	}

	fmt.Println("交易成功广播到区块链网络")
	return nil
}

// 定义广播交易的方法
func broadcastTransaction(serializedTx string) error {
	// Stacks节点API的URL（此处为主网URL）
	apiURL := "https://stacks-node-api.mainnet.stacks.co/v2/transactions"

	// 将十六进制字符串转换为字节数组
	txBytes, err := hex.DecodeString(serializedTx)
	if err != nil {
		return fmt.Errorf("无法解码交易: %v", err)
	}

	// 发送HTTP请求，将交易广播到Stacks节点
	resp, err := http.Post(apiURL, "application/octet-stream", bytes.NewBuffer(txBytes))
	if err != nil {
		return fmt.Errorf("无法广播交易: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("广播交易失败, 状态码: %d", resp.StatusCode)
	}

	fmt.Println("交易成功广播到区块链网络")
	return nil
}

// 定义广播交易的方法
func BroadcastTransactionFromb(serializedTx string) error {
	// Stacks节点API的URL（此处为主网URL）
	apiURL := "https://stacks-node-api.mainnet.stacks.co/v2/transactions"

	// 将十六进制字符串转换为字节数组
	txBytes, err := hex.DecodeString(serializedTx)
	if err != nil {
		return fmt.Errorf("无法解码交易: %v", err)
	}

	// 发送HTTP请求，将交易广播到Stacks节点
	resp, err := http.Post(apiURL, "application/octet-stream", bytes.NewBuffer(txBytes))
	if err != nil {
		return fmt.Errorf("无法广播交易: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("广播交易失败, 状态码: %d", resp.StatusCode)
	}

	fmt.Println("交易成功广播到区块链网络")
	return nil
}
func Transfer1(privateKey string, to string, memo string, amount *big.Int, nonce *big.Int, fee *big.Int) (string, error) {
	// 现有的代码...
	publicKeyHexStr, err := GetPublicKey(privateKey)
	if err != nil {
		return "", err
	}
	signedtokentransferoptions := createSignedTokenTransferOptions(to, *amount, *fee, *nonce, memo, privateKey)
	stacksTransaction, err := makeUnsignedSTXTokenTransfer(signedtokentransferoptions, publicKeyHexStr)
	if err != nil {
		return "", err
	}
	stacksPrivateKey, err := createStacksPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	stacks := StandardAuthorization{
		stacksTransaction.Auth.AuthType,
		&SingleSigSpendingCondition{
			HashMode:    stacksTransaction.Auth.SpendingCondition.HashMode,
			Signer:      stacksTransaction.Auth.SpendingCondition.Signer,
			Nonce:       *big.NewInt(int64(0)),
			Fee:         *big.NewInt(int64(0)),
			KeyEncoding: stacksTransaction.Auth.SpendingCondition.KeyEncoding,
			Signature:   stacksTransaction.Auth.SpendingCondition.Signature,
		},
		stacksTransaction.Auth.SponsorSpendingCondition,
	}
	tmpAuth := stacksTransaction.Auth
	stacksTransaction.Auth = stacks
	tx := Txid(*stacksTransaction)
	stacksTransaction.Auth = tmpAuth

	signer := &TransactionSigner{
		transaction:   stacksTransaction,
		sigHash:       tx,
		originDone:    false,
		checkOversign: true,
		checkOverlap:  true,
	}

	signer.signOrigin(stacksPrivateKey)
	stacksTransaction.Auth.SpendingCondition.Signature = signer.transaction.Auth.SpendingCondition.Signature

	buf := bytes.NewBuffer(make([]byte, 0))
	buf.Write(getBytes(int64(stacksTransaction.Version), 0))
	chainIdBuffer := bytes.NewBuffer(make([]byte, 0, 4))

	chainIdBuffer.Write(getBytesByLength(stacksTransaction.ChainId, 8))
	buf.Write(sliceByteBuffer(chainIdBuffer))
	buf.Write(serializeAuth(&stacksTransaction.Auth))
	buf.Write(getBytes(int64(stacksTransaction.AnchorMode), 0))
	buf.Write(getBytes(int64(stacksTransaction.PostConditionMode), 0))
	buffer2 := serializeLPList(stacksTransaction.PostConditions)
	buf.Write(buffer2)

	buffer3 := SerializePayload(stacksTransaction.Payload)
	buf.Write(buffer3)

	txSerialize := hex.EncodeToString(sliceByteBuffer(buf))
	txId := Txid(*stacksTransaction)
	transactionRes := TransactionRes{
		TxId:        txId,
		TxSerialize: txSerialize,
	}
	res, err := json.Marshal(transactionRes)
	if err != nil {
		return "", err
	}
	return string(res), nil
	// 广播交易
	err = broadcastTransaction(txSerialize)
	if err != nil {
		return "", fmt.Errorf("广播交易失败: %v", err)
	}
	return string(res), nil
}

// 访问API返回结构体
func ReqAPIwhithStruct(apiURL string, method string, result interface{}) error {
	client := &http.Client{Timeout: 30 * time.Second}
	var req *http.Request
	var err error
	if method == http.MethodPost {
		req, err = http.NewRequest(method, apiURL, bytes.NewBuffer(nil))
	} else {
		req, err = http.NewRequest(method, apiURL, nil)
	}
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("意外的状态码: %v", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	return nil
}

// ReqAPIwhithStructAndBro 函数用于发送 API 请求并解析结果
func ReqAPIwhithStructAndBro(apiURL string, method string, result interface{}) error {
	client := &http.Client{Timeout: 30 * time.Second}
	var req *http.Request
	var err error

	// 根据请求方法创建请求
	if method == http.MethodPost {
		req, err = http.NewRequest(method, apiURL, bytes.NewBuffer(nil))
	} else {
		req, err = http.NewRequest(method, apiURL, nil)
	}

	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头以模拟浏览器
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("意外的状态码: %v, 响应体: %s", resp.StatusCode, body)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %v", err)
	}

	// 解析响应
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	return nil
}
