package v1

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	stx "gin-web/api/v1/wallet"
	"gin-web/models"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 获取alex tokens
func GetAlexTokens(url, query string) (models.AlexResponseData, error) {
	requestBody, err := json.Marshal(map[string]string{
		"query": query,
	})
	if err != nil {
		return models.AlexResponseData{}, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return models.AlexResponseData{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.AlexResponseData{}, err
	}
	var response models.AlexResponseData
	err = json.Unmarshal(body, &response)
	if err != nil {
		return models.AlexResponseData{}, err
	}
	return response, nil
}

// 获取最新代币价格
func GetAlexTokens1(url string) (models.AlexTokenPrice2, error) {
	var response models.AlexTokenPrice2

	// 发送 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	// 解析 JSON 数据
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// 提取VelarSwap中的token
func GetVelarTokens(url string) (models.VelarResponseData, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(url), nil)
	if err != nil {
		return models.VelarResponseData{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return models.VelarResponseData{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.VelarResponseData{}, err
	}
	var tokensPrice models.VelarResponseData
	err = json.Unmarshal(body, &tokensPrice)
	if err != nil {
		return models.VelarResponseData{}, err
	}
	return tokensPrice, err
}

// 提取ArkSwap中的token url:https://arkadiko-api.herokuapp.com/api/v1/pages/swap
func GetArkTokens(url string) ([]models.ArkaTokens, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var data map[string]json.RawMessage
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	var tokens []models.ArkaTokens
	for _, rawMsg := range data {
		var token models.ArkaTokens
		err = json.Unmarshal(rawMsg, &token)
		if err != nil {
			// 错误发生时的处理代码
			fmt.Printf("解码错误： %v\n", err)
			continue
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// 监控sUsdt价格
type ApiResponse struct {
	Data struct {
		LaplaceCurrentTokenPrice []struct {
			AvgPriceUsd float64 `json:"avg_price_usd"`
			Token       string  `json:"token"`
		} `json:"laplace_current_token_price"`
	} `json:"data"`
}

func FetchTokenPrice(AlexApiurl string, AlexApiQuery string, laterPrice float64) {
	requestBody, err := json.Marshal(map[string]string{
		"query": AlexApiQuery,
	})
	resp, err := http.Post(AlexApiurl, "application/json", bytes.NewBuffer(requestBody))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// 解析API响应
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	// 遍历返回结果，查找token为token-susdc的avg_price_usd值
	for _, tokenPrice := range apiResponse.Data.LaplaceCurrentTokenPrice {
		if tokenPrice.Token == "token-susdt" {
			if laterPrice == 0 && tokenPrice.AvgPriceUsd > 0 {
				laterPrice = tokenPrice.AvgPriceUsd
			}
			if (tokenPrice.AvgPriceUsd - laterPrice) > 0.07 {
				mailBody := "Susdt价格变化:" + strconv.FormatFloat(tokenPrice.AvgPriceUsd, 'f', -1, 64)
				toAddress := "zyj18883502123@gmail.com"
				SendEmail("1334642655@qq.com", "区块变动-差价产生", mailBody, toAddress)
				fmt.Printf("Token: %s, Avg Price USD: %.6f\n", tokenPrice.Token, tokenPrice.AvgPriceUsd)
			}
			break
		}
	}
	CheckErr(err)
}

// 监控垮桥是否到账
func MakeRequest() {
	// Compose the request
	req, err := http.NewRequest("GET", "https://www.oklink.com/api/explorer/v2/btc/inscription/transfer/list?offset=0&limit=20&type=BRC20&address=bc1p7aerd2nfgx7uzj2yld3mwsyyv3lcnpnxflcp629q7zwpp7kg20psmjsn5e&name=TRIO&direction=3&fromOrTo=&t=1721084522041", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add the necessary headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-apikey", "LWIzMWUtNDU0Ny05Mjk5LWI2ZDA3Yjc2MzFhYmEyYzkwM2NjfDI4MzIxOTU2MzMxNDUyMzU=")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("x-utc", "8")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Parse the response
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// Extract the transferTime value
	transferTime := int64(data["data"].(map[string]interface{})["hits"].([]interface{})[0].(map[string]interface{})["transferTime"].(float64))

	// Convert transferTime to UTC+8 time
	utcPlus8Time := time.Unix(transferTime, 0).UTC().Add(8 * time.Hour)

	// Check if the time is after 5:00 AM on July 16, 2024
	targetTime := time.Date(2024, time.July, 16, 5, 0, 0, 0, time.UTC)
	if utcPlus8Time.After(targetTime) {
		fmt.Println("The transferTime is after 5:00 AM on July 16, 2024. Time:", utcPlus8Time)
		// You can also include additional information from the response here
	} else {
		fmt.Println("The transferTime is not after 5:00 AM on July 16, 2024. Time:", utcPlus8Time)
	}
}

// 地址监控
// 定义返回值的结构体
type Principal struct {
	TypeID       string `json:"type_id"`
	Address      string `json:"address"`
	ContractName string `json:"contract_name,omitempty"`
}

type Asset struct {
	ContractName    string `json:"contract_name"`
	AssetName       string `json:"asset_name"`
	ContractAddress string `json:"contract_address"`
}

type PostCondition struct {
	Type          string    `json:"type"`
	ConditionCode string    `json:"condition_code"`
	Amount        string    `json:"amount"`
	Principal     Principal `json:"principal"`
	Asset         Asset     `json:"asset,omitempty"`
}

type Result struct {
	TxID           string          `json:"tx_id"`
	Nonce          int             `json:"nonce"`
	FeeRate        string          `json:"fee_rate"`
	SenderAddress  string          `json:"sender_address"`
	PostConditions []PostCondition `json:"post_conditions"`
	ReceiptTime    int64           `json:"receipt_time"`
}

type Response struct {
	Limit   int      `json:"limit"`
	Offset  int      `json:"offset"`
	Total   int      `json:"total"`
	Results []Result `json:"results"`
}

func FetchMempoolData(TxID, addr string) string {
	url := "https://api.hiro.so/extended/v1/tx/mempool?address=" + addr + "&limit=30&offset=0&unanchored=true"
	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return ""
	}
	defer resp.Body.Close()
	// 解析返回值
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding response:", err)
		return ""
	}
	// 检查results是否为空
	if len(response.Results) == 0 {
		return ""
	}
	// 获取第一个结果的receipt_time
	receiptTime := response.Results[0].ReceiptTime
	// 获取当前时间
	currentTime := time.Now().Unix()
	// 计算时间差
	timeDifference := currentTime - receiptTime
	// 判断时间差是否大于一分钟（60秒）
	if timeDifference > 90 {
		//temp, _ := strconv.ParseFloat(response.Results[0].FeeRate, 64)
		//fmt.Println("无新交易产生" + fmt.Sprintf("%f", temp/1000000))
	} else {
		if TxID == response.Results[0].TxID && TxID != "" {
			return response.Results[0].TxID
		}
		fmt.Println(time.Now().Format("15:04:10") + "一分钟前产生交易，邮件已发送！" + "地址：" + addr + "\n" + "交易ID：" + response.Results[0].TxID + "\n")
		SendEmail("1334642655@qq.com", "监控地址产生交易", time.Now().Format("15:04:00")+"地址："+addr+"\n"+"TX: https://explorer.hiro.so/txid/"+response.Results[0].TxID+"?chain=mainnet\n"+"交易ID：", "zyj18883502123@gmail.com")
		return response.Results[0].TxID
	}
	return ""
}

// 大单监控
func MonitorBigOrder() {
	url := "https://api.hiro.so/extended/v1/tx/mempool?address=SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.amm-pool-v2-01&limit=30&offset=0&unanchored=true"

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	// 解析返回值
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// 获取当前时间
	currentTime := time.Now().Unix()

	// 遍历每条结果
	for _, result := range response.Results {
		for _, postCondition := range result.PostConditions {
			// 检查asset的contract_name
			if postCondition.Asset.ContractName == "token-susdt" || postCondition.Asset.ContractName == "token-abtc" {
				// 检查amount是否大于900000000
				amount := postCondition.Amount
				if amountValue, err := strconv.ParseInt(amount, 10, 64); err == nil && amountValue > 900000000 {
					// 检查时间差是否小于20分钟
					timeDifference := currentTime - result.ReceiptTime
					if timeDifference < 20*60 {
						FeeRate, _ := strconv.ParseFloat(response.Results[0].FeeRate, 64)
						SendEmail("1334642655@qq.com", "套利:", "Addr:"+result.SenderAddress+"Fee:"+fmt.Sprintf("%f", FeeRate/100000)+"TX:"+result.TxID, "zyj18883502123@gmail.com")
						fmt.Println("存在交易套利")
						fmt.Println("Addr:" + result.SenderAddress + "Fee:" + fmt.Sprintf("%f", FeeRate/1000000) + "TX:" + result.TxID)
						return
					}
				}
			}
		}
	}

	fmt.Println("无大单交易")
}

// 抢单-手续费监控
func OrderSnatch(addr string) {
	url := "https://api.hiro.so/extended/v1/tx/mempool?address=SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.amm-pool-v2-01&limit=30&offset=0&unanchored=true"
	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	// 解析返回值
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	flagTime, _, _ := GetBurnBlockTime()
	fmt.Println("出块儿时间：" + fmt.Sprintf("%d", flagTime))

	// 遍历每条结果
	for _, result1 := range response.Results {
		// 检查sender_address是否为传入的addr
		if result1.SenderAddress == addr && result1.ReceiptTime > flagTime {
			target := result1 // 记录为target
			// 遍历post_conditions
			for _, result := range response.Results {
				if result.ReceiptTime > flagTime {
					// 检查第一个asset的contract_address是否相同
					for _, postCondition := range result.PostConditions {
						// 检查receipt_time是否小于target的receipt_time
						if len(postCondition.Asset.ContractAddress) > 0 && postCondition.Asset.ContractAddress == target.PostConditions[0].Asset.ContractAddress {
							// 检查fee_rate是否大于target的fee_rate
							targetFeeRate, _ := strconv.ParseInt(target.FeeRate, 10, 64)
							currentFeeRate, _ := strconv.ParseInt(result.FeeRate, 10, 64)
							if currentFeeRate >= targetFeeRate {
								// 输出结果
								feeRateOutput := currentFeeRate * 100000
								fmt.Printf("Fee Rate: %d, TX ID: %s, Sender Address: %s\n", feeRateOutput, result.TxID, result.SenderAddress)
								break
							}
						}
					}
				}
			}
		}
		//fmt.Println("Addr:" + addr + "没有TX")
	}
}

// 定义响应结构体
type GetStxBlockTimeResp struct {
	Limit   int64 `json:"limit"`
	Offset  int64 `json:"offset"`
	Total   int64 `json:"total"`
	Results []struct {
		Height        int64 `json:"height"`
		BurnBlockTime int64 `json:"burn_block_time"`
	} `json:"results"`
}

// 工具方法： GetStxBlockTime 方法获取区块高度和爆块时间
func GetBurnBlockTime() (int64, int64, error) {
	url := "https://api.hiro.so/extended/v1/block?limit=1&offset=0&unanchored=true"
	resp, err := http.Get(url)
	CheckErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	CheckErr(err)
	var apiResponse GetStxBlockTimeResp
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return 0, 0, err
	}
	if len(apiResponse.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found")
	}
	block := apiResponse.Results[0]
	return block.Height, block.BurnBlockTime, nil
}

// 工具方法：转换数字到字母
func ConvertNumbersToLetters(n int) []string {
	letters := make([]string, 0)
	for i := 1; i <= n; i++ {
		if i <= 26 {
			letters = append(letters, string('a'+i-1))
		} else {
			letters = append(letters, fmt.Sprintf("number %d out of range", i))
		}
	}
	return letters
}

// Transation monitor -- 原始
func TxMonitor(api, defiTag string, contract, tags, txs []string) []string {
	dxAmount := 0.00
	dyAmount := 0.00
	dyToken := ""
	nycTxs := txs
	//url:="https://explorer.hiro.so/txid/"+contract+"?chain=mainnet"
	//url := "https://api.hiro.so/extended/v1/tx/mempool?address=SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.xyk-core-v-1-1&limit=30&offset=0&unanchored=true"
	resp, err := http.Get(api)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return nil
	}
	defer resp.Body.Close()
	// 解析返回值
	var pendingTx models.PendingTx
	if err := json.NewDecoder(resp.Body).Decode(&pendingTx); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil
	}
	//height, burnTime, err := GetBurnBlockTime()
	_, burnTime, err := GetBurnBlockTime()
	//fmt.Println("出块儿时间：" + fmt.Sprintf("%d%s%s%s%d", burnTime, "=", time.Unix(burnTime, 0).Format("2006-01-02 15:04:05"), ":", height))
	afterBurnTime := make([]models.PendingTxResults, 0)
	// 遍历每条结果
	for _, result1 := range pendingTx.PendingTxResults {
		// 确认当前交易产生在出块后
		if result1.ReceiptTime > burnTime {
			afterBurnTime = append(afterBurnTime, result1)
		}
	}
	nbName := ""
	targetName := ""
	feeRate := 0.0
	for _, tx := range afterBurnTime {
		if len(tx.ContractCall.FunctionArgs) > 0 {
			for _, arg := range tx.ContractCall.FunctionArgs {
				for _, tag := range tags {
					tokenTemp := strings.Split(tag, "-")
					//tx包含传入代币
					if arg.Repr == tag {
						parts := strings.Split(arg.Name, "-")
						if len(parts) > 1 {
							nbName = parts[1]
						}
						if len(parts) > 0 && len(parts) < 1 {
							nbName = parts[0]
						} // This will be "x"
						targetName = "d" + nbName
						if defiTag == "xyk" {
							targetName = "x-amount"
						}
						for _, temparg := range tx.ContractCall.FunctionArgs {
							if temparg.Name == targetName {
								//处理dy
								switch tx.ContractCall.FunctionName {
								case "swap-helper":
									if len(tx.ContractCall.FunctionArgs) > 4 {
										reprStr := tx.ContractCall.FunctionArgs[(len(tx.ContractCall.FunctionArgs) - 4)].Repr
										parts := strings.Split(reprStr, "-")
										//提取数字
										re := regexp.MustCompile("[0-9]+")
										dyAmountStr := re.FindAllString(tx.ContractCall.FunctionArgs[(len(tx.ContractCall.FunctionArgs)-1)].Repr, -1)
										combinedStr := strings.Join(dyAmountStr, "")
										dyAmount64, _ := strconv.ParseFloat(combinedStr, 64)
										dyAmount = dyAmount64 / 1e8
										if len(parts) > 0 {
											dyToken = parts[1]
										}
									}
									break
								case "swap-helper-a":
									if len(tx.ContractCall.FunctionArgs) > 5 {
										reprStr := tx.ContractCall.FunctionArgs[(len(tx.ContractCall.FunctionArgs) - 5)].Repr
										parts := strings.Split(reprStr, "-")
										//提取数字
										re := regexp.MustCompile("[0-9]+")
										dyAmountStr := re.FindAllString(tx.ContractCall.FunctionArgs[(len(tx.ContractCall.FunctionArgs)-1)].Repr, -1)
										combinedStr := strings.Join(dyAmountStr, "")
										dyAmount64, _ := strconv.ParseFloat(combinedStr, 64)
										dyAmount = dyAmount64 / 1e8
										if len(parts) > 0 {
											dyToken = parts[1]
										}
									}
									break
								case "swap-helper-b":
									if len(tx.ContractCall.FunctionArgs) > 5 {
										reprStr := tx.ContractCall.FunctionArgs[(len(tx.ContractCall.FunctionArgs) / 5)].Repr
										parts := strings.Split(reprStr, "-")
										//提取数字
										re := regexp.MustCompile("[0-9]+")
										dyAmountStr := re.FindAllString(tx.ContractCall.FunctionArgs[(len(tx.ContractCall.FunctionArgs)-1)].Repr, -1)
										combinedStr := strings.Join(dyAmountStr, "")
										dyAmount64, _ := strconv.ParseFloat(combinedStr, 64)
										dyAmount = dyAmount64 / 1e8
										if len(parts) > 0 {
											dyToken = parts[1]
										}
									}
									break
								//aark部分
								case "swap-x-for-y":
									for _, farg := range tx.ContractCall.FunctionArgs {
										if farg.Name == "dx" || farg.Name == "x-amount" {
											dxn, _ := strconv.Atoi(farg.Repr[1:])
											dxAmount = float64(dxn)
										}
										if farg.Name == "min-dy" {
											re := regexp.MustCompile(`\d+`)
											dyns := re.FindStringSubmatch(farg.Repr)
											dyn, _ := strconv.ParseInt(dyns[0], 10, 64)
											dyAmount = float64(dyn)
										}
									}
									break
								case "swap-y-for-x":
									for _, farg := range tx.ContractCall.FunctionArgs {
										if farg.Name == "dx" || farg.Name == "x-amount" {
											dxn, _ := strconv.Atoi(farg.Repr[1:])
											dxAmount = float64(dxn)
										}
										if farg.Name == "min-dy" {
											re := regexp.MustCompile(`\d+`)
											dyns := re.FindStringSubmatch(farg.Repr)
											dyn, _ := strconv.ParseInt(dyns[0], 10, 64)
											dyAmount = float64(dyn)
										}
									}
									break
								}
								// 使用正则表达式提取最后一段数字
								re := regexp.MustCompile(`(\d+(\.\d+)?)$`)
								amountStrs := re.FindStringSubmatch(temparg.Repr) //数量
								if (defiTag != "aark" && defiTag != "xyk") && len(amountStrs) > 0 {
									amountStr := amountStrs[1] // 获取最后一段数字
									amount16, _ := strconv.ParseFloat(amountStr, 64)
									dxAmount = amount16 / 1e8
									tempfee, _ := strconv.ParseFloat(tx.FeeRate, 64)
									feeRate = tempfee / 1e6
								}
								if defiTag == "aark" || defiTag == "xyk" {
									dxAmount = dxAmount / 1e6
									dyAmount = dyAmount / 1e6
								}
								for i, inStx := range txs {
									if inStx == tx.TxID {
										//fmt.Println("已存在 = ", tx.TxID, "=>", targetName, "=>", amount, tokenTemp[1], "|| Fee:", feeRate)
										if tx.ReceiptTime < burnTime {
											nycTxs = append(nycTxs[:i], nycTxs[i+1:]...)
										}
										return txs
									}
								}
								txs = append(txs, tx.TxID)

								fmt.Println(time.Now().Format("15:04"), defiTag, "|", tx.TxID, "|", tx.SenderAddress[len(tx.SenderAddress)-5:], dxAmount, "<", tokenTemp[1], ">", "=>", dyAmount, "<", dyToken, ">", "| Fee:", feeRate)
								if dxAmount >= 500 {
									//SendEmail("1334642655@qq.com", tx.ContractCall.FunctionName, "时间："+tx.ReceiptTimeIso.String()+"NYC交易："+arg.Name+"数量:"+strconv.FormatFloat(dyAmount, 'f', 2, 64)+"TX: https://explorer.hiro.so/txid/"+tx.TxID+"?chain=mainnet\n", "2821322847@qq.com")
									//SendEmail("1334642655@qq.com", tx.ContractCall.FunctionName, "时间："+tx.ReceiptTimeIso.String()+"NYC交易："+arg.Name+"数量:"+strconv.FormatFloat(dyAmount, 'f', 2, 64)+"TX: https://explorer.hiro.so/txid/"+tx.TxID+"?chain=mainnet\n", "1883502123@gmail.com")
									//SendEmail("1334642655@qq.com", tx.ContractCall.FunctionName, "时间："+tx.ReceiptTimeIso.String()+"NYC交易："+arg.Name+"数量:"+strconv.FormatFloat(amount, 'f', 2, 64)+"TX: https://explorer.hiro.so/txid/"+tx.TxID+"?chain=mainnet\n", "1871437892@qq.com")
									fmt.Println("邮件发送成功！")
								}
							}
						}
					}
				}
			}
		}
	}
	return txs
}
func GetLatestBlockHeight1() (int, error) {
	var height int
	height1, err := RequestAPIwithStruct("https://mempool.space/api/blocks/tip/height", "get", &height)
	CheckErr(err)
	fmt.Println(height1)
	fmt.Println(height)
	return 1, nil
}

// GetLatestBlockHeight 获取最新区块高度
func GetLatestBlockHeight() (int, error) {
	apiURL := "https://mempool.space/api/blocks/tip/height"
	client := &http.Client{Timeout: 30 * time.Second}
	for i := 0; i < 3; i++ { // 最多重试 3 次
		resp, err := client.Get(apiURL)
		if err != nil {
			if i == 2 { // 最后一次尝试
				return 0, fmt.Errorf("获取最新区块高度失败: %v", err)
			}
			time.Sleep(2 * time.Second) // 等待后重试
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("获取区块高度失败: 状态码 %d", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		CheckErr(err)
		var height int
		err = json.Unmarshal(body, &height)
		CheckErr(err)
		return height, nil
	}
	return 0, fmt.Errorf("获取最新区块高度失败: 超过最大重试次数")
}

// CheckBlockStatus 检查出块情况
func CheckBlockStatus(lastHeight int) (bool, int) {
	newHeight := 0
	ticker := time.NewTicker(59 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			newHeight, _ = GetLatestBlockHeight()
			if newHeight > lastHeight {
				lastHeight = newHeight
				return true, newHeight
				ticker.Stop()
			}
		case <-quit:
			return true, newHeight
		}
	}
	return true, newHeight
}

// 检查是否有交易
// 交易1 0.3 交易2 0.2
func SerchTranstion(blockHeight int, TxID, addr string) string {
	txs := []string{}
	//tag := []string{"'SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2", "'SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc", "'SP2XD7417HGPRTREMKF748VNEQPDRR0RMANB7X1NK.token-susdt", "'SP3Y2ZSH8P7D50B0VBTSX11S7XSG24M1VB9YFQA4K.token-aeusdc"}
	//contract := []string{"https://explorer.hiro.so/txid/SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.xyk-core-v-1-1?chain=mainnet"}
	ticker := time.NewTicker(8 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			//扫描交易
			account, _ := GetCombinedBalance("SP2DG03SMAV8Q8JTDHF9F32Y7B3523ZJYM0Q3MK3Y")
			//整合广播交易
			txs, _ = SeachNYCTx(account, txs, "479dff38f63417579564cc18fdf12d20c561d0fc73765dcccf949a052e4d923701", addr)
		case <-quit:
			return ""
		}
	}
	return ""
}

// 扫描NYC交易
func SeachNYCTx(account *models.CombinedBalance, txs []string, key, addr string) ([]string, *stx.SignedContractCallOptions) {
	//amount := 0.00
	//nycTxs := txs
	txOption := &stx.SignedContractCallOptions{}
	url := "https://api.hiro.so/extended/v1/tx/mempool?address=SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.amm-pool-v2-01&limit=30&offset=0&unanchored=true"
	resp, err := http.Get(url)
	CheckErr(err)
	defer resp.Body.Close()
	var pendingTx models.PendingTx
	CheckErr(err)
	if err := json.NewDecoder(resp.Body).Decode(&pendingTx); err != nil {
		fmt.Println("Error decoding response:", err)
		return txs, txOption
	}
	// 检查results是否为空
	if len(pendingTx.PendingTxResults) == 0 {
		//fmt.Println("No results found.")
		return txs, txOption
	}
	height, burnTime, err := GetBurnBlockTime()
	fmt.Println("出块儿时间：" + fmt.Sprintf("%d%s%d", burnTime, ":", height))
	afterBurnTime := make([]models.PendingTxResults, 0)
	// 遍历每条结果
	for _, result1 := range pendingTx.PendingTxResults {
		// 确认当前交易产生在出块后
		if result1.ReceiptTime > burnTime {
			afterBurnTime = append(afterBurnTime, result1)
		}
	}
	for _, tx := range afterBurnTime {
		//找到对应地址的 未提交交易
		if tx.SenderAddress == addr && tx.ContractCall.FunctionName == "swap-helper" {
			//确认交易在出售 NYC
			if tx.ContractCall.FunctionArgs[0].Repr == "'SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc" {
				constractId := strings.Split(tx.ContractCall.ContractID, ".")
				fee, _ := strconv.ParseInt(tx.FeeRate, 10, 64)
				// args
				tokenX, _ := stx.NewContractPrincipalCV1("")
				tokenY, _ := stx.NewContractPrincipalCV1("")
				factor := stx.NewUintCV(big.NewInt(0))
				dx := stx.NewUintCV(big.NewInt(0))
				mindy := stx.NewUintCV(big.NewInt(0))

				for _, arg := range tx.ContractCall.FunctionArgs {
					if arg.Name == "token-x-trait" {
						stokenX := strings.Split(arg.Repr, "'")
						if len(stokenX) > 1 {
							tokenX, _ = stx.NewContractPrincipalCV1(stokenX[1])
						}
					} else if arg.Name == "token-y-trait" {
						stokenY := strings.Split(arg.Repr, "'")
						if len(stokenY) > 1 {
							tokenY, _ = stx.NewContractPrincipalCV1(stokenY[1])
						}
					} else if arg.Name == "factor" {
						sfactors := strings.Split(arg.Repr, "u")
						sfactor, _ := strconv.ParseInt(sfactors[1], 10, 64)
						sbfactior := big.NewInt(sfactor)
						factor = stx.NewUintCV(sbfactior)
					} else if arg.Name == "dx" {
						sdxs := strings.Split(arg.Repr, "u")
						sdx, _ := strconv.ParseInt(sdxs[1], 10, 64)
						submit := sdx / 10 * 7 //提交价格
						sbdx := big.NewInt(submit)
						dx = stx.NewUintCV(sbdx)
					} else if arg.Name == "min-dy" {
						smindys := strings.Split(arg.Repr, "u")
						smindy, _ := strconv.ParseInt(smindys[1], 10, 64)
						smindy /= 10 * 6
						sbminy := big.NewInt(smindy)
						mindy = stx.NewUintCV(sbminy)
					}
				}
				functionArgs := []stx.ClarityValue{
					&tokenX, // 使用指针类型的 token-x-trait
					&tokenY, // 使用指针类型的 token-y-trait
					factor,  // factor
					dx,      // dx
					&stx.SomeCV{stx.OptionalSome, mindy},
				}
				txOption = &stx.SignedContractCallOptions{
					ContractAddress:   constractId[0],
					ContractName:      constractId[1],
					FunctionName:      tx.ContractCall.FunctionName,
					FunctionArgs:      functionArgs,
					SendKey:           key,
					ValidateWithAbi:   false,
					Fee:               *big.NewInt(fee + 100),
					Nonce:             *big.NewInt(int64(account.API1.Nonce + 1)),
					AnchorMode:        3,
					PostConditionMode: stx.PostConditionModeAllow,
				}
			}
		}
		// 发起合约调用
		sumitTx, _ := stx.MakeContractCall1(txOption)
		// 序列化交易
		txSerialized := hex.EncodeToString(stx.Serialize(*sumitTx))
		txId := stx.Txid(*sumitTx)
		// 广播交易
		err = stx.BroadcastTransaction(txSerialized)
		fmt.Println("Serialized Transaction:", txSerialized)
		fmt.Println("Transaction ID:", txId)
		CheckErr(err)
		txs = append(txs, txId)
	}
	return txs, txOption
}

// 根据地址压单 0=压单 1=秒提
func SubmitTxFromAddr(req models.SubmitTxFromAddrReq) {
	//[]models.PendingTxResults
	//req.MySubmitTx
	submitedTxId := ""
	submitedTxStruct := models.PendingTxResults{}
	factor := int64(100000000)
	amount := 0.00
	myTxStruct := models.PendingTxResults{}
	//myTxStructTat := 0
	pendingTxOption := GetSubmitTxByStruct(myTxStruct) //根据地址压单
	if req.TargetAddress != "" && req.Api == "" {
		for {
			api := req.ApiPre + req.TargetAddress + req.ApiFix
			pendingTx := GetPendingTx(api)
			for _, tx := range pendingTx {
				if tx.SenderAddress == req.TargetAddress {
					for _, arg := range tx.ContractCall.FunctionArgs {
						if arg.Name == "token-x-trait" {
							//秒单模式tx需要传入---只执行一次
							//if myTxStructTat == 0 {
							myTxStruct = GetStructByTx(req.MyTx)
							//	myTxStructTat = 1
							//}
							for _, functionArg := range myTxStruct.ContractCall.FunctionArgs {
								if functionArg.Name == "token-x-trait" && functionArg.Repr == arg.Repr { //原始token相同
									for _, arg := range tx.ContractCall.FunctionArgs {
										if arg.Name == "token-y-trait" {
											for _, myarges := range myTxStruct.ContractCall.FunctionArgs {
												if myarges.Name == "token-y-trait" && arg.Repr == myarges.Repr { //目标token相同
													for _, inarg := range tx.ContractCall.FunctionArgs {
														switch inarg.Name {
														case "factor":
															smindys := strings.Split(inarg.Repr, "u")
															factor, _ = strconv.ParseInt(smindys[1], 10, 64)
															break
														case "dx":
															smindys := strings.Split(inarg.Repr, "u")
															dxamount, _ := strconv.ParseFloat(smindys[1], 0)
															n := int(math.Floor(math.Log10(float64(factor))) + 1)
															amount = dxamount / math.Pow(10, float64(n))
															fmt.Println("交易额", amount)
															break
														}
														if amount > 0 {
															//交易整合
															pendingTxOption = GetSubmitTxByStruct(myTxStruct)
															//填充手续费
															if tx.FeeRate >= myTxStruct.FeeRate {
																fee, _ := strconv.ParseInt(tx.FeeRate, 10, 64)
																fee = fee + 1
																if fee > 5010100 {
																	fee = 0
																}
																pendingTxOption.Fee = *big.NewInt(fee)
															}
															//填充none
															nonce := *big.NewInt(int64(0))
															nonce = *big.NewInt(int64(req.Account.API1.Nonce))
															if submitedTxId != "" {
																submitedTxStruct = GetStructByTx(submitedTxId)
																nonce = *big.NewInt(int64(submitedTxStruct.Nonce))
															}
															//压单提交
															if req.SubmitModel == 1 {
																nonce = *big.NewInt(int64(req.Account.API1.Nonce)) //出块儿秒提
															}
															//
															pendingTxOption.Nonce = nonce
															//填充key
															pendingTxOption.SendKey = req.SenderKey
															fmt.Println("Ptx:", pendingTxOption)
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			// 检查出块儿秒提--秒单
			if req.SubmitModel == 1 {
				for {
					//秒单
					newHeight, err := GetLatestBlockHeight()
					fmt.Println("秒单", newHeight)
					CheckErr(err)
					if (int64(newHeight) - req.StartBlock) > 0 {
						submitedTxId = SubmitTxByTxOption(req, pendingTxOption)
						break
					}
					time.Sleep(1 * time.Second) //1秒调用一次
				}
			} else if req.SubmitModel == 0 {
				//压单
				newHeight, err := GetLatestBlockHeight()
				fmt.Println("压单", newHeight)
				CheckErr(err)
				submitedTxId = SubmitTxByTxOption(req, pendingTxOption)
				if (int64(newHeight) - req.StartBlock) > 1 {
					break
				}
			}
			req.SubmitModel = 0
			time.Sleep(1 * time.Second) //1秒调用一次
		}
	} else if req.Api != "" {
		for {
			pendingTx := GetPendingTx(req.Api)
			for _, tx := range pendingTx {
				for _, arg := range tx.ContractCall.FunctionArgs {
					if arg.Name == "token-x-trait" {
						myTxStruct = GetStructByTx(req.MyTx)
						for _, functionArg := range myTxStruct.ContractCall.FunctionArgs {
							if functionArg.Name == "token-x-trait" && functionArg.Repr == arg.Repr {
								if arg.Name == "token-y-trait" && functionArg.Name == "token-y-trait" {
									if functionArg.Repr == arg.Repr {
										for _, inarg := range tx.ContractCall.FunctionArgs {
											switch inarg.Name {
											case "factor":
												smindys := strings.Split(inarg.Repr, "u")
												factor, _ = strconv.ParseInt(smindys[1], 10, 64)
												break
											case "dx":
												smindys := strings.Split(inarg.Repr, "u")
												dxamount, _ := strconv.ParseFloat(smindys[1], 0)
												n := int(math.Floor(math.Log10(float64(factor))) + 1)
												amount = dxamount / math.Pow(10, float64(n))
												fmt.Println("交易额", amount)
												break
											}
											if amount > 10 {
												//交易整合
												pendingTxOption = GetSubmitTxByStruct(myTxStruct)
												//填充手续费
												if tx.FeeRate > myTxStruct.FeeRate {
													fee, _ := strconv.ParseInt(tx.FeeRate, 10, 64)
													fee = fee + 1
													if fee > 5010100 {
														fee = 0
													}
													pendingTxOption.Fee = *big.NewInt(fee)
												}
												//填充none
												nonce := *big.NewInt(int64(0))
												nonce = *big.NewInt(int64(myTxStruct.Nonce)) //压单提交
												//nonce = *big.NewInt(int64(req.Account.API1.Nonce)) //出块儿秒提
												pendingTxOption.Nonce = nonce
												//填充key
												pendingTxOption.SendKey = req.SenderKey
											}
										}
									}
								}
							}
						}
					}
				}
			}
			fmt.Println("1")
			// 检查出块儿秒提--秒单
			if req.SubmitModel == 1 {
				for {
					//秒单
					newHeight, err := GetLatestBlockHeight()
					CheckErr(err)
					if (int64(newHeight) - req.StartBlock) > 0 {
						SubmitTxByTxOption(req, pendingTxOption)
						break
					}
					time.Sleep(1 * time.Second) //1秒调用一次
				}
			} else if req.SubmitModel == 0 {
				//压单
				newHeight, err := GetLatestBlockHeight()
				CheckErr(err)
				SubmitTxByTxOption(req, pendingTxOption)
				if (int64(newHeight) - req.StartBlock) > 1 {
					break
				}
			}
			req.SubmitModel = 0
			time.Sleep(1 * time.Second) //1秒调用一次
		}
	}
	//
	////序列化广播交易
	//SubmitTxByTxOption(req, pendingTxOption)
	//return
}
func SubmitTxByTxOption(req models.SubmitTxFromAddrReq, pendingTxOption *stx.SignedContractCallOptions) string {
	txId := ""
	//提交交易
	if req.SubmitModel == 0 { //压单
		// 发起合约调用
		sumitTx, _ := stx.MakeContractCall1(pendingTxOption)
		// 序列化交易
		txSerialized := hex.EncodeToString(stx.Serialize(*sumitTx))
		txId := stx.Txid(*sumitTx)
		// 广播交易
		err := stx.BroadcastTransaction(txSerialized)
		fmt.Println("Serialized Transaction:", txSerialized)
		fmt.Println("压单：Transaction ID:", txId)
		CheckErr(err)
	} else if req.SubmitModel == 1 { //抢单
		// 发起合约调用
		sumitTx, _ := stx.MakeContractCall1(pendingTxOption)
		// 序列化交易
		txSerialized := hex.EncodeToString(stx.Serialize(*sumitTx))
		txId := stx.Txid(*sumitTx)
		// 广播交易
		err := stx.BroadcastTransaction(txSerialized)
		fmt.Println("Serialized Transaction:", txSerialized)
		fmt.Println("压单：Transaction ID:", txId)
		CheckErr(err)
	}
	return txId
}

// 根据结构体 返回待签名tx --压单 SubmitTxFromAddr
func GetSubmitTxByStruct(txStruct models.PendingTxResults) *stx.SignedContractCallOptions {
	// 使用比较函数判断 txStruct 是否为空
	if IsPendingTxResultsEmpty(txStruct) {
		fmt.Println("txStruct 是空的（刚初始化）")
		return nil
	}
	// 设置交易参数
	contract := strings.Split(txStruct.ContractCall.ContractID, ".")
	//contractAddress := "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM" // 使用大写字母
	//contractName := "amm-pool-v2-01"                               // 合约名称
	//functionName := "swap-helper"                                  // 函数名称
	//nonce := big.NewInt(int64(reqStruct.Account.API1.Nonce))       // nonce
	//fee := big.NewInt(int64(minfee - 1))                           // 手续费 0.001002

	// 使用指针类型的 tokenX 和 tokenY
	tokenX, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2")
	tokenY, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc")
	//tokenZ, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc")
	functionArgs := []stx.ClarityValue{
		&tokenX,                               // 使用指针类型的 token-x-trait
		&tokenY,                               // 使用指针类型的 token-y-trait
		stx.NewUintCV(big.NewInt(100000000)),  // factor
		stx.NewUintCV(big.NewInt(6000000000)), // dx
		&stx.SomeCV{stx.OptionalSome, stx.NewUintCV(big.NewInt(2041000000000))}, // 使用指针类型的 min-dy
	}
	for _, arg := range txStruct.ContractCall.FunctionArgs[:len(txStruct.ContractCall.FunctionArgs)] {
		switch arg.Name {
		case "token-x-trait":
			parts := strings.Split(arg.Repr, "'")
			tokenX, _ = stx.NewContractPrincipalCV1(parts[1])
		case "token-y-trait":
			parts := strings.Split(arg.Repr, "'")
			tokenY, _ = stx.NewContractPrincipalCV1(parts[1])
			//case "token-z-trait": tokenZ, _ = stx.NewContractPrincipalCV1(arg.Repr)
		}
	}
	// 创建交易选项
	txOption := &stx.SignedContractCallOptions{
		ContractAddress:   contract[0],
		ContractName:      contract[1],
		FunctionName:      txStruct.ContractCall.FunctionName,
		FunctionArgs:      functionArgs,
		ValidateWithAbi:   false,
		AnchorMode:        3,
		PostConditionMode: stx.PostConditionModeAllow,
	}
	return txOption
}
func IsPendingTxResultsEmpty(txStruct models.PendingTxResults) bool {
	return txStruct.TxID == "" &&
		txStruct.Nonce == 0 &&
		txStruct.FeeRate == "" &&
		txStruct.SenderAddress == "" &&
		!txStruct.Sponsored &&
		txStruct.PostConditionMode == "" &&
		len(txStruct.PostConditions) == 0 &&
		txStruct.AnchorMode == "" &&
		txStruct.TxStatus == "" &&
		txStruct.ReceiptTime == 0 &&
		txStruct.ReceiptTimeIso.IsZero() &&
		txStruct.TxType == ""
}

// 根据 tx 返回结构体 --压单 SubmitTxFromAddr
func GetStructByTx(tx string) models.PendingTxResults {
	api := fmt.Sprintf("https://api.hiro.so/extended/v1/tx/%s?event_offset=0&event_limit=100&unanchored=true", tx)
	resp, err := http.Get(api)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return models.PendingTxResults{}
	}
	defer resp.Body.Close()
	// 解析返回值
	var pendingTx models.PendingTxResults
	if err := json.NewDecoder(resp.Body).Decode(&pendingTx); err != nil {
		fmt.Println("Error decoding response:", err)
		return models.PendingTxResults{}
	}
	return pendingTx
}

// 返回pending交易 --压单 SubmitTxFromAddr
func GetPendingTx(api string) []models.PendingTxResults {
	resp, err := http.Get(api)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return nil
	}
	defer resp.Body.Close()
	// 解析返回值
	var pendingTx models.PendingTx
	if err := json.NewDecoder(resp.Body).Decode(&pendingTx); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil
	}
	//height, burnTime, err := GetBurnBlockTime()
	_, burnTime, err := GetBurnBlockTime()
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

// 抢单、压单
type PendingTx struct {
	Results []PendingTxResult `json:"results"`
}

type PendingTxResult struct {
	TxID          string       `json:"tx_id"`
	FeeRate       string       `json:"fee_rate"`
	SenderAddress string       `json:"sender_address"`
	ContractCall  ContractCall `json:"contract_call"`
	ReceiptTime   int64        `json:"receipt_time"`
}

type ContractCall struct {
	ContractID   string        `json:"contract_id"`
	FunctionName string        `json:"function_name"`
	FunctionArgs []FunctionArg `json:"function_args"`
}

type FunctionArg struct {
	Repr string `json:"repr"`
	Name string `json:"name"`
}

func SubmitTxFromAddr1(key, txAddr, targetAddr, hitModel string, blockHeight int) {
	// 1. 获取链上交易数据
	//url := "https://api.hiro.so/extended/v1/tx/mempool?address=" + txAddr + "&limit=30&offset=0&unanchored=true"
	url := "https://api.hiro.so/extended/v1/tx/mempool?address=SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.amm-pool-v2-01&limit=30&offset=0&unanchored=true"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching transaction data:", err)
		return
	}
	defer resp.Body.Close()
	var pendingTx PendingTx
	if err := json.NewDecoder(resp.Body).Decode(&pendingTx); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}
	if len(pendingTx.Results) == 0 {
		fmt.Println("No transactions found for the given address.")
		return
	}
	// 2. 遍历交易并处理逻辑
	for _, tx := range pendingTx.Results {
		if tx.SenderAddress == targetAddr {
			fee, _ := strconv.ParseInt(tx.FeeRate, 10, 64)
			tokenX, _ := stx.NewContractPrincipalCV1("")
			tokenY, _ := stx.NewContractPrincipalCV1("")
			factor := stx.NewUintCV(big.NewInt(0))
			dx := stx.NewUintCV(big.NewInt(0))
			mindy := stx.NewUintCV(big.NewInt(0))

			for _, arg := range tx.ContractCall.FunctionArgs {
				if arg.Name == "token-x-trait" {
					stokenX := strings.Split(arg.Repr, "'")
					if len(stokenX) > 1 {
						tokenX, _ = stx.NewContractPrincipalCV1(stokenX[1])
					}
				} else if arg.Name == "token-y-trait" {
					stokenY := strings.Split(arg.Repr, "'")
					if len(stokenY) > 1 {
						tokenY, _ = stx.NewContractPrincipalCV1(stokenY[1])
					}
				} else if arg.Name == "factor" {
					sfactors := strings.Split(arg.Repr, "u")
					sfactor, _ := strconv.ParseInt(sfactors[1], 10, 64)
					sbfactior := big.NewInt(sfactor)
					factor = stx.NewUintCV(sbfactior)
				} else if arg.Name == "dx" {
					sdxs := strings.Split(arg.Repr, "u")
					sdx, _ := strconv.ParseInt(sdxs[1], 10, 64)
					submit := sdx / 10 * 7 // 提交价格
					sbdx := big.NewInt(submit)
					dx = stx.NewUintCV(sbdx)
				} else if arg.Name == "min-dy" {
					smindys := strings.Split(arg.Repr, "u")
					smindy, _ := strconv.ParseInt(smindys[1], 10, 64)
					smindy /= 10 * 6
					sbminy := big.NewInt(smindy)
					mindy = stx.NewUintCV(sbminy)
				}
			}

			functionArgs := []stx.ClarityValue{
				&tokenX,
				&tokenY,
				factor,
				dx,
				&stx.SomeCV{stx.OptionalSome, mindy},
			}

			txOption := &stx.SignedContractCallOptions{
				ContractAddress: txAddr,
				FunctionName:    tx.ContractCall.FunctionName,
				FunctionArgs:    functionArgs,
				SendKey:         key,
				Fee:             *big.NewInt(fee + 100),
			}

			submitTx, _ := stx.MakeContractCall1(txOption)
			txSerialized := stx.Serialize(*submitTx)
			txId := stx.Txid(*submitTx)

			err = stx.BroadcastTransaction(string(txSerialized))
			fmt.Println("Broadcasted Transaction ID:", txId)
			if err != nil {
				fmt.Println("Error broadcasting transaction:", err)
			}

			//// 处理hitModel模式逻辑
			//if hitModel == "1" {
			//	for i := 0; i < 120; i++ { // 1分钟
			//		time.Sleep(500 * time.Millisecond)
			//		CheckForNewTransactions(txAddr, key)
			//	}
			//}
		}
	}
}

func CheckForNewTransactions(txAddr, key string) {
	// 检测新交易
	url := "https://api.hiro.so/extended/v1/tx/mempool?address=" + txAddr
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching new transactions:", err)
		return
	}
	defer resp.Body.Close()
	var pendingTx PendingTx
	if err := json.NewDecoder(resp.Body).Decode(&pendingTx); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}
	if len(pendingTx.Results) == 0 {
		return
	}
	// 遍历并处理新交易
	for _, tx := range pendingTx.Results {
		fmt.Println("New transaction detected:", tx.TxID)
	}
}

// func SerchPendingTx(contract, tags, txs []string) []string {
func NycBot(reqStruct models.NycBotReqStruct) models.NycBotReqStruct {
	balance, _ := strconv.ParseInt(reqStruct.Account.API2.STX.Balance, 10, 64)
	if balance < 100000000 {
		reqStruct = RedeemNyc(reqStruct)
		// 1 = 1000000
	}
	if reqStruct.Submiting {
		fmt.Println("区块未变动", reqStruct.Submiting)
		return reqStruct
	}
	pendingTxs := GetPendingTx(reqStruct.Api)
	targetName := "dx"
	dxAmount := 0.00
	countSell := reqStruct.NycAmount
	minfee := 0.00
	targetTx := make([][]models.PendingTxFunctionArgs, 0)
	for _, tx := range pendingTxs {
		if len(tx.ContractCall.FunctionArgs) > 0 {
			for _, arg := range tx.ContractCall.FunctionArgs {
				for _, tag := range reqStruct.TokenContract {
					//判断是否在做swap
					if strings.Contains(tx.ContractCall.FunctionName, "swap") {
						//判断是否为在交易的代币合约，是否在做卖出动作
						if arg.Repr == tag && strings.Contains(arg.Name, "x") {
							//找到目标交易（卖出目标代币）添加到集合
							targetTx = append(targetTx, tx.ContractCall.FunctionArgs)
							for _, temparg := range tx.ContractCall.FunctionArgs {
								//找到dx
								if temparg.Name == targetName {
									//获取出卖出的值
									//提取数字
									re := regexp.MustCompile("[0-9]+")
									dxAmountStr := re.FindAllString(temparg.Repr, -1)
									combinedStr := strings.Join(dxAmountStr, "")
									dxAmount64, _ := strconv.ParseFloat(combinedStr, 64)
									dxAmount = dxAmount64 / 1e8
									//处理卖出总量---卖出最低手续费
									countSell += dxAmount
									tempfee, _ := strconv.ParseFloat(tx.FeeRate, 64)
									minfee = tempfee
									if minfee > tempfee {
										minfee = tempfee
									}
								}
							}
						}
					}
				}
			}
		}
	}
	if countSell >= 17400 {
		//if countSell >= 25300 {//
		//如果迈出量达到，开始套利
		//	jsonData := []byte(`[
		//   {
		//      "principal": {
		//         "type_id": "principal_standard",
		//         "address": "SPN6PP6AWX1QXY86M6XT9PY6AF9SRM6RK9GXZ7F7"
		//      },
		//      "condition_code": "sent_equal_to",
		//      "amount": "10000000",
		//      "type": "stx"
		//   },
		//   {
		//      "principal": {
		//         "type_id": "principal_contract",
		//         "address": "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM",
		//         "contract_name": "amm-vault-v2-01"
		//      },
		//      "condition_code": "sent_greater_than_or_equal_to",
		//      "amount": "175000000",
		//      "type": "fungible",
		//      "asset": {
		//         "asset_name": "newyorkcitycoin",
		//         "contract_address": "SPSCWDV3RKV5ZRN1FQD84YE1NQFEDJ9R1F4DYQ11",
		//         "contract_name": "newyorkcitycoin-token-v2"
		//      }
		//   }
		//]`)
		//
		//	postConditions, err := convertJSONToPostConditions(jsonData)
		//	fmt.Println(postConditions)

		// 设置交易参数
		senderKey := reqStruct.SenderKey
		contractAddress := "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM"            // 使用大写字母
		contractName := "amm-pool-v2-01"                                          // 合约名称
		functionName := "swap-helper"                                             // 函数名称
		nonce := big.NewInt(int64(reqStruct.Account.NonceInfo.PossibleNextNonce)) // nonce
		//fee := big.NewInt(int64(minfee - 2001))                                   // 手续费 5001 = 0.005001
		fee := big.NewInt(int64(minfee - 1)) // 手续费 5001 = 0.005001
		// 使用指针类型的 tokenX 和 tokenY
		tokenX, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2")
		tokenY, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc")

		functionArgs := []stx.ClarityValue{
			&tokenX,                               // 使用指针类型的 token-x-trait
			&tokenY,                               // 使用指针类型的 token-y-trait
			stx.NewUintCV(big.NewInt(100000000)),  // factor
			stx.NewUintCV(big.NewInt(5000000000)), // dx 5000000000 == 50
			//stx.NewUintCV(big.NewInt(10000000)), // dx 5000000000 == 50
			//&stx.SomeCV{stx.OptionalSome, stx.NewUintCV(big.NewInt(171005898411))}, // 使用指针类型的 min-dy
			&stx.SomeCV{stx.OptionalSome, stx.NewUintCV(big.NewInt(1750000000000))}, // 1750000000000 ==17500
			//&stx.SomeCV{stx.OptionalSome, stx.NewUintCV(big.NewInt(175000000))}, // 1750000000000 ==17500
		}
		if minfee <= 0 {
			fee = big.NewInt(int64(5001))
		}
		// 创建交易选项
		//x := stx.DeserializePostCondition("010316e685b016b3b6cd9ebf35f38e5ae29392e2acd51d0f616c65782d7661756c742d76312d3116e685b016b3b6cd9ebf35f38e5ae29392e2acd51d176167653030302d676f7665726e616e63652d746f6b656e04616c657803000000000078b854")
		//postCondition, _ := convertJSONToPostConditions(jsonData)
		//postConditions := []stx.PostCondition{
		//	{
		//		StacksMessage: stx.StacksMessage{},
		//		ConditionType: 0,
		//		Principal:     nil,
		//		ConditionCode: 0,
		//	},
		//	{
		//		StacksMessage: stx.StacksMessage{
		//			Type: 0,
		//		},
		//		ConditionType: 5,
		//		Principal: stx.PostConditionPrincipalInterface{
		//			TypeID:  "principal_standard",
		//			Address: "SPN6PP6AWX1QXY86M6XT9PY6AF9SRM6RK9GXZ7F7",
		//		},
		//		ConditionCode: "sent_equal_to",
		//		Amount:        "50000000",
		//		Type:          "stx",
		//	},
		//	{
		//		StacksMessage: StacksMessage{
		//			Type: 0,
		//		},
		//		ConditionType: 5,
		//		Principal: Principal{
		//			TypeID:       "principal_contract",
		//			Address:      "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM",
		//			ContractName: new(string),
		//		},
		//		ConditionCode: "sent_greater_than_or_equal_to",
		//		Amount:        "17341626758",
		//		Type:          "fungible",
		//		Asset: &Asset{
		//			AssetName:       "newyorkcitycoin",
		//			ContractAddress: "SPSCWDV3RKV5ZRN1FQD84YE1NQFEDJ9R1F4DYQ11",
		//			ContractName:    "newyorkcitycoin-token-v2",
		//		},
		//	},
		//}
		txOption := &stx.SignedContractCallOptions{
			ContractAddress:   contractAddress,
			ContractName:      contractName,
			FunctionName:      functionName,
			FunctionArgs:      functionArgs,
			SendKey:           senderKey,
			ValidateWithAbi:   false,
			Fee:               *fee,
			Nonce:             *nonce,
			AnchorMode:        3,
			PostConditionMode: stx.PostConditionModeAllow,
			//PostConditions:    postCondition,
		}
		// 发起合约调用
		tx, err := stx.MakeContractCall1(txOption)
		// 序列化交易
		txSerialized := hex.EncodeToString(stx.Serialize(*tx))
		txId := stx.Txid(*tx)
		fmt.Println("fee:", fee)
		// 广播交易
		//for {
		//秒单
		newHeight, err := GetLatestBlockHeight()
		fmt.Println("秒单", newHeight)
		CheckErr(err)
		//if (int64(newHeight) - int64(reqStruct.SubmitBlock)) > 0 {
		err = stx.BroadcastTransaction(txSerialized)
		if err == nil {
			reqStruct.SubmitTxSerialized = txSerialized
			reqStruct.SubmitTx = txId
			reqStruct.SubmitBlock = newHeight
			reqStruct.Submiting = true
		}
		fmt.Println(err)
		//break
		//}
		//time.Sleep(1 * time.Second) //1秒调用 一次
		//}
		CheckErr(err)
		//记录剩余数量
		//sellamount, _ := strconv.ParseFloat("26250", 64)
		//reqStruct.NycAmount = countSell - sellamount
	}
	return reqStruct
}

func convertJSONToPostConditions(jsonData []byte) ([]stx.PostConditionInterface, error) {
	var jsonPostConditions []models.PostConditionJSON
	err := json.Unmarshal(jsonData, &jsonPostConditions)
	if err != nil {
		return nil, err
	}

	var postConditions []stx.PostConditionInterface
	for _, jsonPC := range jsonPostConditions {
		var principal stx.PostConditionPrincipalInterface
		if jsonPC.Principal.TypeID == "principal_standard" {
			principal = stx.PostConditionPrincipal{
				Type:    stx.PRINCIPAL,
				Prefix:  2,
				Address: *stx.CreateAddress(jsonPC.Principal.Address),
			}
		} else if jsonPC.Principal.TypeID == "principal_contract" {
			principal = stx.PostConditionPrincipal{
				Type:    stx.PRINCIPAL,
				Prefix:  2,
				Address: *stx.CreateAddress(jsonPC.Principal.Address),
			}
			// 根据实际情况设置合同名称等其他字段
			// principal.ContractName = stx.createLPString(jsonPC.Principal.ContractName, nil, nil)
		}
		conditionCode, _ := strconv.Atoi(jsonPC.ConditionCode)
		postConditions = append(postConditions, stx.PostCondition{
			StacksMessage: stx.StacksMessage{
				Type: stx.POSTCONDITION,
			},
			ConditionType: 1,
			Principal:     principal,
			ConditionCode: conditionCode,
		})
	}

	return postConditions, nil
}

// random Nyc
func RedeemNyc(req models.NycBotReqStruct) models.NycBotReqStruct {
	// 创建交易选项
	txOption := &stx.SignedContractCallOptions{
		ContractAddress:   "SP8A9HZ3PKST0S42VM9523Z9NV42SZ026V4K39WH",
		ContractName:      "ccd012-redemption-nyc",
		FunctionName:      "redeem-nyc",
		FunctionArgs:      []stx.ClarityValue{},
		SendKey:           req.SenderKey,
		ValidateWithAbi:   false,
		Fee:               *big.NewInt(int64(5101)),
		Nonce:             *big.NewInt(int64(req.Account.NonceInfo.PossibleNextNonce)),
		AnchorMode:        3,
		PostConditionMode: stx.PostConditionModeAllow,
	}
	// 发起合约调用
	tx, err := stx.MakeContractCall1(txOption)
	// 序列化交易
	txSerialized := hex.EncodeToString(stx.Serialize(*tx))
	txId := stx.Txid(*tx)
	// 输出结果
	fmt.Println("Serialized Transaction:", txSerialized)
	fmt.Println("Transaction ID:", txId)
	// 广播交易
	err = stx.BroadcastTransaction(txSerialized)
	if err == nil {
		req.SubmitBlock, _ = GetLatestBlockHeight()
	}
	CheckErr(err)
	req.Redeem = true
	return req
}

// 提交交易
func SubmitNYCTxWithNewFee(fee int64, senderKey string) {

}

// 复制交易
func SubmitTxFromTx(targetTx models.PendingTxResults, account models.CombinedBalance, reqTxoption stx.SignedContractCallOptions) string {
	contractAddress := "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM" // 使用大写字母
	contractName := "amm-pool-v2-01"
	functionName := "swap-helper"
	nonce := big.NewInt(int64(1))
	fee := big.NewInt(int64(0))
	tokenX, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2")
	tokenY, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc")
	factor := stx.NewUintCV(big.NewInt(100000000))
	dx := stx.NewUintCV(big.NewInt(5000000000))
	dy := stx.NewUintCV(big.NewInt(1750000000000))
	//替换值
	parts := strings.Split(targetTx.ContractCall.ContractID, ".")
	targetTxFee, _ := strconv.Atoi(targetTx.FeeRate)
	contractAddress = parts[0]                                     // 使用大写字母
	contractName = parts[1]                                        // 合约名称
	functionName = targetTx.ContractCall.FunctionName              // 函数名称
	nonce = big.NewInt(int64(account.NonceInfo.PossibleNextNonce)) // nonce 	// 手续费 5001 = 0.005001
	fee = big.NewInt(int64(targetTxFee + 1))                       // 手续费 5001 = 0.005001
	switch targetTx.ContractCall.FunctionName {
	case "swap-helper":
		for _, arg := range targetTx.ContractCall.FunctionArgs {
			if arg.Name == "token-x-trait" {
				tokencontract := arg.Repr[1:]
				tokenX, _ = stx.NewContractPrincipalCV1(tokencontract)
			}
			if arg.Name == "token-y-trait" {
				tokencontract := arg.Repr[1:]
				tokenY, _ = stx.NewContractPrincipalCV1(tokencontract)
			}
			if arg.Name == "factor" {
				intValue, _ := strconv.Atoi(arg.Repr[1:])
				factor = stx.NewUintCV(big.NewInt(int64(intValue)))
			}
			if arg.Name == "dx" {
				dxn, _ := strconv.Atoi(arg.Repr[1:])
				dx = stx.NewUintCV(big.NewInt(int64(dxn)))
			}
			if arg.Name == "min-dy" {
				re := regexp.MustCompile(`\d+`)
				dyns := re.FindStringSubmatch(arg.Repr)
				dyn, _ := strconv.ParseInt(dyns[0], 10, 64)
				dy = stx.NewUintCV(big.NewInt(dyn))
			}
		}
		break
	case "swap-helper-a":
		break
	case "swap-helper-b":
		break
	}
	if reqTxoption.FunctionArgs != nil {

	}
	functionArgs := []stx.ClarityValue{
		&tokenX,
		&tokenY,
		factor,
		dx,
		&stx.SomeCV{stx.OptionalSome, dy},
	}
	txOption := &stx.SignedContractCallOptions{
		ContractAddress:   contractAddress,
		ContractName:      contractName,
		FunctionName:      functionName,
		FunctionArgs:      functionArgs,
		SendKey:           reqTxoption.SendKey,
		ValidateWithAbi:   false,
		Fee:               *fee,
		Nonce:             *nonce,
		AnchorMode:        3,
		PostConditionMode: stx.PostConditionModeAllow,
	}
	tx, err := stx.MakeContractCall1(txOption)
	txSerialized := hex.EncodeToString(stx.Serialize(*tx))
	//txId := stx.Txid(*tx)
	//err = stx.BroadcastTransaction(txSerialized)
	CheckErr(err)
	return txSerialized
	//fmt.Println("广播交易成功，fee:", fee, "tx:", txId)
}

// 修改手续费
func ChangeFee(req models.ChangeFee) string {
	contractAddress := "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM" // 使用大写字母
	contractName := "amm-pool-v2-01"
	functionName := "swap-helper"
	nonce := big.NewInt(int64(1))
	fee := big.NewInt(int64(0))
	tokenX, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2")
	tokenY, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc")
	factor := stx.NewUintCV(big.NewInt(100000000))
	dx := stx.NewUintCV(big.NewInt(5000000000))
	dy := stx.NewUintCV(big.NewInt(1750000000000))
	//替换值
	parts := strings.Split(req.ChangeTxStruct.ContractCall.ContractID, ".")

	targetTxFee, _ := strconv.Atoi(req.TargetTxStruct.FeeRate)

	contractAddress = parts[0]                                         // 使用大写字母
	contractName = parts[1]                                            // 合约名称
	functionName = req.ChangeTxStruct.ContractCall.FunctionName        // 函数名称
	nonce = big.NewInt(int64(req.Account.NonceInfo.PossibleNextNonce)) // nonce 	// 手续费 5001 = 0.005001
	fee = big.NewInt(int64(targetTxFee + 1))                           // 手续费 5001 = 0.005001
	switch req.ChangeTxStruct.ContractCall.FunctionName {
	case "swap-helper":
		for _, arg := range req.ChangeTxStruct.ContractCall.FunctionArgs {
			if arg.Name == "token-x-trait" {
				tokencontract := arg.Repr[1:]
				tokenX, _ = stx.NewContractPrincipalCV1(tokencontract)
			}
			if arg.Name == "token-y-trait" {
				tokencontract := arg.Repr[1:]
				tokenY, _ = stx.NewContractPrincipalCV1(tokencontract)
			}
			if arg.Name == "factor" {
				intValue, _ := strconv.Atoi(arg.Repr[1:])
				factor = stx.NewUintCV(big.NewInt(int64(intValue)))
			}
			if arg.Name == "dx" {
				dxn, _ := strconv.Atoi(arg.Repr[1:])
				dx = stx.NewUintCV(big.NewInt(int64(dxn)))
			}
			if arg.Name == "min-dy" {
				dyn, _ := strconv.Atoi(arg.Repr[1:])
				dy = stx.NewUintCV(big.NewInt(int64(dyn)))
			}
		}
		break
	case "swap-helper-a":
		break
	case "swap-helper-b":
		break
	}
	functionArgs := []stx.ClarityValue{
		&tokenX,
		&tokenY,
		factor,
		dx,
		&stx.SomeCV{stx.OptionalSome, dy},
	}
	txOption := &stx.SignedContractCallOptions{
		ContractAddress:   contractAddress,
		ContractName:      contractName,
		FunctionName:      functionName,
		FunctionArgs:      functionArgs,
		SendKey:           req.SenderKey,
		ValidateWithAbi:   false,
		Fee:               *fee,
		Nonce:             *nonce,
		AnchorMode:        3,
		PostConditionMode: stx.PostConditionModeAllow,
	}
	tx, err := stx.MakeContractCall1(txOption)
	txSerialized := hex.EncodeToString(stx.Serialize(*tx))
	//txId := stx.Txid(*tx)
	//err = stx.BroadcastTransaction(txSerialized)
	CheckErr(err)
	//fmt.Println("广播交易成功，fee:", fee,"tx:",txId)
	return txSerialized
}
