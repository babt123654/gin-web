package service

import (
	"encoding/json"
	"fmt"
	"gin-web/models"
	//"github.com/okx/go-wallet-sdk/crypto/vrf/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (my MysqlService) GetStxBlockChange() {
	var node []byte
	var beforeValue, afterValue []byte
	var err error
	// 每过5秒钟调用一次方法
	ticker := time.Tick(5 * time.Second)
	for {
		select {
		case <-ticker:
			//CheckSentUserOnchainStatus()
			node, _ = getUrlResult("https://api.hiro.so/extended/v1/block?limit=1")
			_, afterValue, err = getFieldValueFromJSON(node, "total")
			checkErr(err)
			if jsonEqual(beforeValue, afterValue) {
				fmt.Println("无变动" + string(beforeValue) + "==" + string(afterValue))
			} else {
				fmt.Println(time.Now().Format("01-02 15:04:05") + " ： " + string(beforeValue) + "-->" + string(afterValue))
				body := "区块变动：" + string(beforeValue) + "->" + string(afterValue)
				if string(beforeValue) == "" {
					body = "监控开始：" + time.Now().Format("2006-01-02 15:04:05") + "---" + string(beforeValue) + "-->" + string(afterValue)
				}
				//toaddr := "babt543834198328@gmail.com"
				toaddr := "zyj18883502123@gmail.com"
				sendEmail("1334642655@qq.com", "区块变动", body, toaddr)
				fmt.Println("Email sent successfully!")

				////开始获取其他defi的token价格
				//velarTokens, err := GetVelarToken("https://mainnet-prod-proxy-service-dedfb0daae85.herokuapp.com/swapapp/swap/tokens")
				//checkErr(err)
				//fmt.Println(string(velarTokens))
				//_, value, err := GetFieldValueFromJSON(velarTokens, "message")
				//var m []models.Message
				//json.Unmarshal(value, &m)
				//for _, asset := range m {
				//	if asset.AssetName == "STX" {
				//		//提取指定代币 ：唯一标识 symbol
				//		//fmt.Println("Order:", asset.Symbol+":"+asset.Price)
				//		//对比代币价格
				//
				//	}
				//fmt.Println("Order:", asset.Symbol+":"+asset.Price)
				//}
				//welsh价格
				//Tokens := extractWelshSymbols()
				//for _, token := range Tokens {
				//	if token.Symbol == "welsh" {
				//		fmt.Println(token.Price)
				//	}
				//}
				beforeValue = afterValue
			}
		}
	}
}
func temp() {
	//开始获取welsh价格
	welshTokens, err := GetWelshToken("https://wacky-inflatable-45.s3.us-east-2.amazonaws.com/4e29977e-f8fb-48ce-97dd-2b304c32ff7e/95f46ae7-0cba-4576-a971-041c473e2c63/af114948-07c9-4655-8e71-42905ff40137.json")
	_, tokens, err := GetFieldValueFromJSON(welshTokens, "tokens")
	nane, stxPriceByte, err := GetFieldValueFromJSON(welshTokens, "stxPrice")
	stxPrice := byte2float(stxPriceByte)
	fmt.Println(nane+":", string(stxPriceByte))
	checkErr(err)
	var tokensList []models.Tokens
	var welshsList []models.Tokens
	maxPrice := 0.0
	minPrice := 0.0
	maxPriceApp := ""
	minPriceApp := ""
	body := "无变动"
	json.Unmarshal(tokens, &tokensList)
	//fmt.Println(string(welshTokens))
	for _, token := range tokensList {
		if token.Symbol == "WELSH" {
			welshsList = append(welshsList, token)
			//提取指定代币 ：唯一标识 symbol
			swapApp := applyRegexString(token.SwapContractImage, "\\.(.*?)\\.")
			price := token.Price * stxPrice
			if !(swapApp == "coinmarketcap") {
				if minPrice == 0.0 {
					minPrice = price
					minPriceApp = swapApp
				} else if (price - minPrice) > 0.0001 {
					maxPrice = price
					maxPriceApp = swapApp
				} else {
					minPrice = price
					minPriceApp = swapApp
				}
				if maxPrice-minPrice > 0.0001 {
					body = "差价产生:" + maxPriceApp + "-最大价格:" + strconv.FormatFloat(maxPrice, 'f', -1, 64) + "/n最小价格:" + minPriceApp + strconv.FormatFloat(minPrice, 'f', -1, 64)
					//+ "->" + string(afterValue)
					//sendEmail("1334642655@qq.com", "区块变动-差价产生", body)
					fmt.Println("Email sent successfully!")
				}
			}
			fmt.Println(body)
		}
		//fmt.Println("Order:", asset.Symbol+":"+asset.Price)
	}
}
func GetWelshToken(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(url), nil)
	checkErr(err)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	//fmt.Println(string(body))
	return body, err
}

// 根据url返回json []byte
func getUrlResult(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(url), nil)
	checkErr(err)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	return body, err
}

func extractWelshSymbols() []models.Tokens {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.satscreener.com/stx#"), nil)
	checkErr(err)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	fmt.Println("welsh:" + string(body))
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil
	}
	var Tokens []models.Tokens
	tokens, ok := data["token"].([]interface{})
	if !ok {
		return Tokens
	}
	for _, t := range tokens {
		token, ok := t.(map[string]interface{})
		if ok && token["symbol"] == "welsh" {
			model := models.Tokens{
				ID:                token["id"].(string),
				Name:              token["name"].(string),
				SwapContract:      token["swapContract"].(string),
				TotalSupply:       int(token["totalSupply"].(float64)),
				ContractID:        token["contractId"].(string),
				Symbol:            token["symbol"].(string),
				Price:             token["price"].(float64),
				Txns:              int(token["txns"].(float64)),
				Volume:            token["volume"].(float64),
				LastBlockTxn:      int(token["lastBlockTxn"].(float64)),
				SwapContractImage: token["swapContractImage"].(string),
				HrVolume:          int(token["hrVolume"].(float64)),
				DayVolume:         int(token["dayVolume"].(float64)),
				WeekVolume:        token["weekVolume"].(float64),
				HrChange:          int(token["hrChange"].(float64)),
				DayChange:         int(token["dayChange"].(float64)),
				WeekChange:        token["weekChange"].(float64),
			}
			Tokens = append(Tokens, model)
		}
	}
	return Tokens
}

// 获取welsh各个平台的价格
func GetWelshPrice() []models.Welsh {
	//获取当前区块
	node, _ := getUrlResult("https://api.hiro.so/extended/v1/block?limit=1")
	_, block, _ := getFieldValueFromJSON(node, "total")
	//获取welsh的价格
	var welsh models.Welsh
	var welshs []models.Welsh
	//获取stx的价格
	stxPriceByte, _ := GetStxPrice()
	stxPrice := byte2float(stxPriceByte)
	//提取stackSwap的价格
	stSwapTokens, _ := GetStackSwapToken(string(block))
	xPrice, _ := strconv.ParseFloat(stSwapTokens.FromPair.FeeBalanceX, 64)
	yPrice, _ := strconv.ParseFloat(stSwapTokens.ToPair.FeeBalanceY, 64)
	welshPrice := xPrice / yPrice * stxPrice
	welsh.Price = welshPrice
	welsh.Com = "stackSwap"
	//append(welshs, welsh)
	fmt.Println(welsh.Com + ":" + strconv.FormatFloat(welshPrice, 'f', -1, 64))
	//提取VelarSwap的价格

	return welshs
}

// 提取ArkSwap中的token url:https://arkadiko-api.herokuapp.com/api/v1/pages/swap
func (my MysqlService) GetArkTokens(url string) ([]models.ArkaTokens, error) {
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

// stackswap 代币信息 https://app.stackswap.org/api/v1/swap_v4/SP3NE50GEXFG9SZGTT51P40X2CKYSZ5CC4ZTZ7A2G.welshcorgicoin-token/SP1Z92MPDQEWZXW36VX71Q25HKF5K2EPCJ304F275.wstx-token-v4a/152583
func GetStackSwapToken(c string) (models.StakeSwap, error) {
	url := "https://app.stackswap.org/api/v1/swap_v4/SP3NE50GEXFG9SZGTT51P40X2CKYSZ5CC4ZTZ7A2G.welshcorgicoin-token/SP1Z92MPDQEWZXW36VX71Q25HKF5K2EPCJ304F275.wstx-token-v4a/" + c
	response, err := http.Get(url)
	if err != nil {
		return models.StakeSwap{}, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return models.StakeSwap{}, err
	}
	var data models.StakeSwap
	err = json.Unmarshal(body, &data)
	if err != nil {
		return models.StakeSwap{}, err
	}
	return data, nil
}
func GetStxPrice() ([]byte, error) {
	return reqByUrl("https://explorer.hiro.so/stxPrice?blockBurnTime=current")
}

// 监控是否收到UST

func CheckSentUserOnchainStatus() (string, error) {
	url := "https://stacks-bridge-api.alexlab.co/v1/token-bridge/orders?address0=SPKXQJP7Z3ZRKRF5QE1BP8YXYX1KMCS4AMGYK1RN&address1=0x5d5316041651d2e52d33f891117db4f8b46c52e6"
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var responseData models.SUsdt
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", err
	}
	if len(responseData.Orders) > 0 {
		sentUserOnchainStatus := responseData.Orders[0].SentUserOnchainStatus
		if sentUserOnchainStatus != "" {
			body := "跨链成功"
			toaddr := "1871437892@qq.com"
			sendEmail("1334642655@qq.com", "跨链成功", body, toaddr)
			toaddr = "zyj18883502123@gmail.com"
			sendEmail("1334642655@qq.com", "跨链成功", body, toaddr)
			fmt.Println("Email sent successfully!")
			return sentUserOnchainStatus, nil
		}
	}

	return "", nil
}
