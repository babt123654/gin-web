package tests

import (
	"encoding/hex"
	"fmt"
	v1 "gin-web/api/v1"
	stx "gin-web/api/v1/wallet"
	"gin-web/models"
	"gin-web/pkg/service"
	"github.com/piupuer/go-helper/pkg/resp"
	"github.com/stretchr/testify/require"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	AlexStx2sUsdt   = 0.005
	ValerStx2aeUsdt = 0.003
)

// STX相关测试
func TestBlockChanwge(t *testing.T) {
	os.Setenv("TEST_CONF", "F:\\Hy\\gin-web-dev\\conf\\")
	Config()
	Mysql()
	q := service.New(ctx)
	//获取区块变动
	q.GetStxBlockChange()
	//s, _ := q.GetArkTokens("https://arkadiko-api.herokuapp.com/api/v1/pages/swap")
	//s, _ := q.GetStackSwapToken("152731")
	//stackSwapWelshPrice := strconv.ParseFloat(s.FromPair.FeeBalanceX, 64) / strconv.ParseFloat(s.ToPair.FeeBalanceY, 64) * s.
	//	fmt.Println()
}

// 垮桥监控 https: //stacks-bridge-api.alexlab.co/v1/token-bridge/orders?address0=SPKXQJP7Z3ZRKRF5QE1BP8YXYX1KMCS4AMGYK1RN&address1=0x5d5316041651d2e52d33f891117db4f8b46c52e6

func TestGetAlexTokens(t *testing.T) {
	stxnuber := 100.00
	alexstxPrice := 0.00
	velarstxPrice := 0.00
	AlexApiurl1 := "https://api.alexgo.io/v2/public/token-prices"
	//AlexApiurl := "https://gql-v1.alexlab.co/v1/graphql"
	//AlexApiQuery := "query FetchLatestPrices {\n  laplace_current_token_price {\n    avg_price_usd\n    token\n  }\n}"
	VelarApiUrl := "https://gateway.velar.network/swapapp/swap/tokens"
	arkSwapApiUrl := "https://arkadiko-api.herokuapp.com/api/v1/pages/swap"
	//stackSwapApiUrl := "https://stacks-bridge-api.alexlab.co/v1/token-bridge/orders?address0=SPKXQJP7Z3ZRKRF5QE1BP8YXYX1KMCS4AMGYK1RN&address1=0x5d5316041651d2e52d33f891117db4f8b46c52e6"

	//c := new(gin.Context)
	//alex tokens
	alexTokens, err := v1.GetAlexTokens1(AlexApiurl1)
	//velar
	velarTokens, err := v1.GetVelarTokens(VelarApiUrl)
	//arkswap
	arkSwapTokens, err := v1.GetArkTokens(arkSwapApiUrl)
	for _, alextoken := range alexTokens.Data {
		if alextoken.ContractID == "trio" || alextoken.ContractID == "TRIO" {
			alexstxPrice = alextoken.LastPriceUsd
			println("trio" + alextoken.ContractID + " : " + strconv.FormatFloat(alextoken.LastPriceUsd, 'f', -1, 64))
		}
		if alextoken.ContractID != "" && alextoken.ContractID == "SP2XD7417HGPRTREMKF748VNEQPDRR0RMANB7X1NK.token-susdt" {
			println("Alex:Stx->sUsdt " + strconv.FormatFloat(alextoken.LastPriceUsd*alexstxPrice*(stxnuber-(stxnuber*AlexStx2sUsdt)), 'f', -1, 64))
			println("Alex:Stx->sUsdt " + strconv.FormatFloat(alextoken.LastPriceUsd*alexstxPrice*stxnuber, 'f', -1, 64))
			println("alex:Stx->sUsdt " + strconv.FormatFloat(alextoken.LastPriceUsd, 'f', -1, 64))
			println("--------------------------------")
		}
	}
	for _, velartoken := range velarTokens.Message {
		if velartoken.Symbol == "stx" || velartoken.Symbol == "STX" {
			velarstxPrice, _ = strconv.ParseFloat(velartoken.Price, 64)
			println("velar" + strconv.FormatFloat(velarstxPrice, 'f', -1, 64))
		}
		if velartoken.ContractAddress != "" && velartoken.ContractAddress == "SP3Y2ZSH8P7D50B0VBTSX11S7XSG24M1VB9YFQA4K.token-aeusdc" {
			aeUsdtPrice, _ := strconv.ParseFloat(velartoken.Price, 64)
			println("valer:Stx->aeUsdt " + strconv.FormatFloat(aeUsdtPrice*velarstxPrice*(stxnuber-(stxnuber*ValerStx2aeUsdt)), 'f', -1, 64))
			println("valer:Stx->aeUsdt " + strconv.FormatFloat(aeUsdtPrice*velarstxPrice*stxnuber, 'f', -1, 64))
			println("aeusdc：" + strconv.FormatFloat(aeUsdtPrice, 'f', -1, 64))
			println("--------------------------------")
		}
	}
	//遍历alex价格
	for _, alextoken := range alexTokens.Data {
		for _, velartoken := range velarTokens.Message {
			if alextoken.ContractID != "" && velartoken.ContractAddress != "" && alextoken.ContractID == velartoken.ContractAddress {
				//println("Alex：" + alextoken.ContractID + " = Velar:" + velartoken.Symbol + " = velar:" + velartoken.ContractAddress)
				parts := strings.SplitN(alextoken.ContractID, ".", 2)
				println("Alex" + parts[1] + " : " + strconv.FormatFloat(alextoken.LastPriceUsd*alexstxPrice*100, 'f', -1, 64))
				velartokenPrice, _ := strconv.ParseFloat(velartoken.Price, 64)
				println("Velar" + velartoken.Symbol + " : " + strconv.FormatFloat(velartokenPrice*velarstxPrice*100, 'f', -1, 64))
				println("--------------------------------")
			}
		}
	}
	//stackswap
	//stackSwapTokens, err := v1.GetStackSwapToken("152583")
	resp.CheckErr(err)
	// 对 response 进行筛选操作，提取所需的数据
	fmt.Println(alexTokens)
	fmt.Println(velarTokens.Message)
	fmt.Println(arkSwapTokens)
	//fmt.Println(stackSwapTokens)
}

// 监控sUsdt价格
func TestGetSsUsdt(t *testing.T) {
	AlexApiurl := "https://gql-v1.alexlab.co/v1/graphql"
	AlexApiQuery := "query FetchLatestPrices {\n  laplace_current_token_price {\n    avg_price_usd\n    token\n  }\n}"
	ticker := time.NewTicker(20 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				v1.FetchTokenPrice(AlexApiurl, AlexApiQuery, 0)
			case <-quit:
				//ticker.Stop()
				return
			}
		}
	}()
	// 等待程序退出
	time.Sleep(5 * time.Minute)
	quit <- struct{}{}

	////alex tokens
	//alexTokens, err := v1.GetAlexTokens(AlexApiurl, AlexApiQuery)
	//
	//alexTokens.Data.LaplaceCurrentTokenPrice[0].AvgPriceUsd
	//resp.CheckErr(err)
	//fmt.Println(s)
}

// 监控TRIO价格
func TestRTIO(t *testing.T) {
	// 每十分钟发起一次请求
	for {
		//v1.MakeRequest()
		alexstxPrice := 0.0
		AlexApiurl1 := "https://api.alexgo.io/v2/public/token-prices"
		//AlexApiurl := "https://gql-v1.alexlab.co/v1/graphql"
		//AlexApiQuery := "query FetchLatestPrices {\n  laplace_current_token_price {\n    avg_price_usd\n    token\n  }\n}"
		alexTokens, _ := v1.GetAlexTokens1(AlexApiurl1)
		for _, alextoken := range alexTokens.Data {
			if alextoken.ContractID == "stx" || alextoken.ContractID == "STX" {
				alexstxPrice = alextoken.LastPriceUsd
				println("Alex" + alextoken.ContractID + " : " + strconv.FormatFloat(alextoken.LastPriceUsd, 'f', -1, 64))
			}
			if alextoken.ContractID == "SP3K8BC0PPEVCV7NZ6QSRWPQ2JE9E5B6N3PA0KBR9.brc20-trio" {
				println("TRIO" + alextoken.ContractID + " : " + strconv.FormatFloat(alextoken.LastPriceUsd, 'f', -1, 64))
				println("100stx = " + strconv.FormatFloat((alextoken.LastPriceUsd*100)/(alexstxPrice*0.98), 'f', -1, 64))
			}
		}
		time.Sleep(40 * time.Second)
	}
}

// 地址监控
func TestMoniterAddr(t *testing.T) {
	ticker := time.NewTicker(20 * time.Second)
	txs := []string{}
	tag := []string{
		"'SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2",
		models.TOKENPRE + "SP2C2YFP12AJZB4MABJBAJ55XECVS7E4PMMZ89YZR.wrapped-stx-token", //aark
		models.TOKENPRE + "SP2C2YFP12AJZB4MABJBAJ55XECVS7E4PMMZ89YZR.usda-token",        //aark
		models.TOKENPRE + "SP2XD7417HGPRTREMKF748VNEQPDRR0RMANB7X1NK.token-abtc",
		models.TOKENPRE + "SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.token-stx-v-1-1", //xyk
		models.TOKENPRE + "SP3Y2ZSH8P7D50B0VBTSX11S7XSG24M1VB9YFQA4K.token-aeusdc",    //xyk
		"'SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc",
		"'SP2XD7417HGPRTREMKF748VNEQPDRR0RMANB7X1NK.token-susdt",
	}
	//SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.amm-pool-v2-01
	//SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.xyk-core-v-1-1

	//api
	addrApiPre := "https://api.hiro.so/extended/v1/tx/mempool?address="
	addrApiAfter := "&limit=30&offset=0&unanchored=true"
	aarkAPI := "https://api.hiro.so/extended/v1/tx/mempool?address=SP2C2YFP12AJZB4MABJBAJ55XECVS7E4PMMZ89YZR.arkadiko-swap-v2-1&limit=30&offset=0&unanchored=true"
	xykAPI := "https://api.hiro.so/extended/v1/tx/mempool?address=SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.xyk-core-v-1-1&limit=30&offset=0&unanchored=true"
	api := "https://api.hiro.so/extended/v1/tx/mempool?address=SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.amm-pool-v2-01&limit=30&offset=0&unanchored=true"
	api1 := addrApiPre + "SP3RQY2KAMTAAZWMAVEMFKC9QBX0MH91010MNBZ9E" + addrApiAfter
	api2 := addrApiPre + "SP2J1K2BX0HNTP8JEZ23MR3Z38JFPQA49X0BJ45AF" + addrApiAfter
	api3 := addrApiPre + "SP3RYGTCPHMMD2MSJ7VHX7AZWQHHXF4JZP7BQW0G0" + addrApiAfter
	api4 := addrApiPre + "SP34VN50ZVG2BG9TKNH2318SYXY2MRKZJXATA43RF" + addrApiAfter
	api5 := addrApiPre + "SP259FA3GEF3QNKBJYD2N03WBE608FK1R6HTGAYVG" + addrApiAfter
	api6 := addrApiPre + "SP2DG03SMAV8Q8JTDHF9F32Y7B3523ZJYM0Q3MK3Y" + addrApiAfter
	api7 := addrApiPre + "SP3ZGH32VKMJCJ2KJTY46E8FCKA4MFYDXCQ8PFX2P" + addrApiAfter
	api8 := addrApiPre + "SP3Z3RAFPPGZ0XHQBZSGN86JJK3C469BXV2ZHTXYX" + addrApiAfter
	apiusda := addrApiPre + "SP1YNFAAGY02F784AAR10TW3AQPEK4Y74S4RADD12" + addrApiAfter
	//contract := []string{"https://explorer.hiro.so/txid/SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.xyk-core-v-1-1?chain=mainnet"}
	contract := ""
	APIXYK := "https://nameless-autumn-scion.stacks-mainnet.quiknode.pro/d01ddfac943c034fad15701a1f74674da5760d1f/extended/v1/address/SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.xyk-core-v-1-1/mempool"
	txId1 := []string{}
	txId2 := []string{}
	txId3 := []string{}
	txId4 := []string{}
	txId5 := []string{}
	txId6 := []string{}
	xyx := []string{}
	XYK := []string{}
	jiaz := []string{} //夹子
	usda := []string{} //夹子
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			txId1 = v1.TxMonitor(api1, contract, txId1, tag, txs)
			txId2 = v1.TxMonitor(api2, contract, txId2, tag, txs)
			txId3 = v1.TxMonitor(api3, contract, txId3, tag, txs)
			txId4 = v1.TxMonitor(api4, contract, txId4, tag, txs)
			txId5 = v1.TxMonitor(api5, contract, txId5, tag, txs)
			txId6 = v1.TxMonitor(api6, contract, txId6, tag, txs)
			jiaz = v1.TxMonitor(api7, contract, jiaz, tag, txs)
			xyx = v1.TxMonitor(api8, contract, xyx, tag, txs)
			usda = v1.TxMonitor(apiusda, contract, usda, tag, txs)
			XYK = v1.TxMonitor(APIXYK, contract, XYK, tag, txs)
			txs = v1.TxMonitor(api, "alex", XYK, tag, txs)
			txs = v1.TxMonitor(aarkAPI, "aark", XYK, tag, txs)
			txs = v1.TxMonitor(xykAPI, "xyk", XYK, tag, txs)
		}
	}
}

// nycbot---测试通过
func TestFoundNycTx(t *testing.T) {
	txs := []string{}
	nowBlock, _ := v1.GetLatestBlockHeight()
	API := "https://api.hiro.so/extended/v1/tx/mempool?address=SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.amm-pool-v2-01&limit=30&offset=0&unanchored=true"
	addr := "SP1KRZZ8H976W67W0XZG5E4DP334FRRWN7JBM39KP"
	account, err := v1.GetCombinedBalance(addr)
	if err != nil {
		fmt.Println("Error fetching account balance:", err)
		return
	}
	req := models.NycBotReqStruct{
		//SenderKey:     "479dff38f63417579564cc18fdf12d20c561d0fc73765dcccf949a052e4d923701",
		SenderKey:     "2ab3352c76241f3e94f960e598e967e4b1970fd4fa3b16d9956a9b80f3ac3f9c01", //2ab3352c76241f3e94f960e598e967e4b1970fd4fa3b16d9956a9b80f3ac3f9c01 -- SP1KRZZ8H976W67W0XZG5E4DP334FRRWN7JBM39KP
		SwapContract:  []string{"https://explorer.hiro.so/txid/SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.xyk-core-v-1-1?chain=mainnet"},
		TokenContract: []string{"'SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc"},
		Txs:           txs,
		Account:       account,
		Api:           API,
		Submiting:     false,
		SubmitBlock:   nowBlock,
	}
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			resp := v1.NycBot(req)
			req = resp
			fmt.Println("等待其他人卖出NYC：", req.Submiting, "nyc", req.NycAmount)
			newHeight, err := v1.GetLatestBlockHeight()
			if err != nil {
				fmt.Println("Error fetching latest block height:", err)
				continue
			}
			if (int64(newHeight) - int64(req.SubmitBlock)) > 0 {
				account, err := v1.GetCombinedBalance(addr)
				if err != nil {
					fmt.Println("Error fetching combined balance:", err)
					continue
				}
				req.Account = account
				req.Submiting = false

			}
		}
	}
}

// 手续费监控，抢单
func TestMoniterFee(t *testing.T) {
	// 创建一个每分钟触发一次的Ticker
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // 确保在程序结束时停止Ticker
	// 使用一个无限循环来监听Ticker
	for {
		select {
		case <-ticker.C:
			v1.OrderSnatch("SP3RQY2KAMTAAZWMAVEMFKC9QBX0MH91010MNBZ9E") //SPN6PP6AWX1QXY86M6XT9PY6AF9SRM6RK9GXZ7F7
		}
	}
}

// nycbot --测试通过
func TestNycBot(t *testing.T) {
	//1.出块监控
	flag := false
	newHeight, err := v1.GetLatestBlockHeight()
	v1.CheckErr(err)
	flag, newHeight = v1.CheckBlockStatus(861694)
	if flag {
		v1.SerchTranstion(newHeight, "", "SP2DG03SMAV8Q8JTDHF9F32Y7B3523ZJYM0Q3MK3Y")
	} //2.交易扫描---扫描卖单
	//3.
}
func TestSeachNYCTx(t *testing.T) {
	txs := []string{}
	addr := "SP21NKRPJACJQAGEGC0DPS6BEPMB1K1X7VJY86N6Z"
	account, _ := v1.GetCombinedBalance(addr)
	v1.SeachNYCTx(account, txs, "479dff38f63417579564cc18fdf12d20c561d0fc73765dcccf949a052e4d923701", addr)
}
func TestGetBlockHeight(t *testing.T) {
	v1.GetLatestBlockHeight1()
	fmt.Println(v1.GetLatestBlockHeight())
}

// 交易测试
func TestsendTranseation(t *testing.T) {
	v1.GetLatestBlockHeight1()
	fmt.Println(v1.GetLatestBlockHeight())
}

// 压单工具
func TestSubmitTxFromAddr(t *testing.T) {
	nowBlock, _ := v1.GetLatestBlockHeight()
	ApiPre := "https://api.hiro.so/extended/v1/tx/mempool?address="
	ApiFix := "&limit=30&offset=0&unanchored=true"
	addr := "SP1H92C38MEGV82R8D5BJ9FVDJ03XBMQG60CMBVNZ"
	account, _ := v1.GetCombinedBalance(addr)
	req := models.SubmitTxFromAddrReq{
		//479dff38f63417579564cc18fdf12d20c561d0fc73765dcccf949a052e4d923701 -- SP1H92C38MEGV82R8D5BJ9FVDJ03XBMQG60CMBVNZ
		//2ab3352c76241f3e94f960e598e967e4b1970fd4fa3b16d9956a9b80f3ac3f9c01 -- SP1KRZZ8H976W67W0XZG5E4DP334FRRWN7JBM39KP
		//c0cee26708dcd2d16ccb6dfde394090a1fafde7cec4fb083d740cbc3021b39bf01 -- SP21NKRPJACJQAGEGC0DPS6BEPMB1K1X7VJY86N6Z
		SenderKey:     "2ab3352c76241f3e94f960e598e967e4b1970fd4fa3b16d9956a9b80f3ac3f9c01",
		TargetAddress: "SP34VN50ZVG2BG9TKNH2318SYXY2MRKZJXATA43RF",
		SubmitModel:   1,
		StartBlock:    int64(nowBlock),
		Api:           "",
		ApiPre:        ApiPre,
		ApiFix:        ApiFix,
		Account:       account,
		MyTx:          "0x74c9f1e2be10700924af9e5516781454761924d3bc6ac276daf47cc6fcc59c4d",
	}
	v1.SubmitTxFromAddr(req)
}

// x-y-k-bot--开发中
func TestVelarBotx(t *testing.T) {
	//api := "https://nameless-autumn-scion.stacks-mainnet.quiknode.pro/d01ddfac943c034fad15701a1f74674da5760d1f/extended/v1/address/SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.xyk-core-v-1-1/mempool"

	//获取pengdingTX ---返回两部分，Need1BlockTx,Need2BlockTx
	//筛选出交易ae的交易确认在卖出还是买入
}
func TestForkTxBot(t *testing.T) {
	addr := "SP1H92C38MEGV82R8D5BJ9FVDJ03XBMQG60CMBVNZ"
	senderKey := "479dff38f63417579564cc18fdf12d20c561d0fc73765dcccf949a052e4d923701"
	account, _ := v1.GetCombinedBalance(addr)
	api := "https://nameless-autumn-scion.stacks-mainnet.quiknode.pro/d01ddfac943c034fad15701a1f74674da5760d1f/extended/v1/address/SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.amm-pool-v2-01/mempool"
	tokenName := "token-y-trait"
	tokenRepr := "'SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc"
	pendingTx := models.PendingTx{}
	//contractAddress := "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM" // 使用大写字母
	//contractName := "amm-pool-v2-01"
	//functionName := "swap-helper"
	//nonce := big.NewInt(int64(1))
	//fee := big.NewInt(int64(0))
	tokenX, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2")
	tokenY, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc")
	factor := stx.NewUintCV(big.NewInt(100000000))
	dx := stx.NewUintCV(big.NewInt(5000000000))
	dy := stx.NewUintCV(big.NewInt(1750000000000))
	txOption1 := stx.SignedContractCallOptions{
		ContractAddress: "",
		ContractName:    "",
		FunctionName:    "",
		FunctionArgs: []stx.ClarityValue{
			&tokenX,
			&tokenY,
			factor,
			dx,
			&stx.SomeCV{stx.OptionalSome, dy},
		},
		SendKey:                 senderKey,
		ValidateWithAbi:         false,
		Fee:                     big.Int{},
		Nonce:                   big.Int{},
		AnchorMode:              0,
		PostConditionMode:       0,
		PostConditions:          nil,
		SerializePostConditions: nil,
	}
	//获取pending交易
	v1.RequestAPIwithStruct(api, http.MethodGet, &pendingTx)
	//筛选出块前交易
	afterBlockTx := v1.GetAfterBlockTX(pendingTx)
	//筛选出手续费最高的交易
	targetTx := v1.GetMaxFeeTXwithTokens(tokenRepr, tokenName, addr, afterBlockTx)
	//复制交易
	txSerialized := v1.SubmitTxFromTx(targetTx, *account, txOption1)
	//txSerialized:=v1.ChangeFee()
	stx.BroadcastTransaction(txSerialized)
}

// 获取当前区块儿高度---测试通过
func TestGetLatestBlockHeightx(t *testing.T) {
	fmt.Println(v1.GetLatestBlockHeight())
}

// redeemNyc---测试
func TestRedeemNYc(t *testing.T) {
	addr := "SP1KRZZ8H976W67W0XZG5E4DP334FRRWN7JBM39KP"
	account, _ := v1.GetCombinedBalance(addr)
	blockHeight, _ := v1.GetLatestBlockHeight()
	req := models.NycBotReqStruct{
		SenderKey: "2ab3352c76241f3e94f960e598e967e4b1970fd4fa3b16d9956a9b80f3ac3f9c01",
		Account:   account,
		NowBlock:  0,
		//SubmitTx:    false,
		SubmitBlock: blockHeight,
		Redeem:      false,
	}
	redeemResp := v1.RedeemNyc(req)
	nowBlock, _ := v1.GetLatestBlockHeight()
	if int64(nowBlock)-redeemResp.NowBlock >= 2 {
		redeemResp.Redeem = true
	}
}

// 测试获取nonce --- 测试通过
func TestRedeemNycForSerch(t *testing.T) {
	api := "https://docs-demo.stacks-mainnet.quiknode.pro/extended/v1/address/" + "SP1KRZZ8H976W67W0XZG5E4DP334FRRWN7JBM39KP" + "/nonces"
	nonceInfo := models.NonceInfo{}
	err := stx.ReqAPIwhithStructAndBro(api, http.MethodGet, &nonceInfo)
	resp.CheckErr(err)
	fmt.Println(nonceInfo)
}
func TestSubmitTx(t *testing.T) {
	addr := "SP1KRZZ8H976W67W0XZG5E4DP334FRRWN7JBM39KP" //479dff38f63417579564cc18fdf12d20c561d0fc73765dcccf949a052e4d923701
	//addr := "SP21NKRPJACJQAGEGC0DPS6BEPMB1K1X7VJY86N6Z"
	account, _ := v1.GetCombinedBalance(addr)
	senderKey := "2ab3352c76241f3e94f960e598e967e4b1970fd4fa3b16d9956a9b80f3ac3f9c01"
	contractAddress := "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM"  // 使用大写字母
	contractName := "amm-pool-v2-01"                                // 合约名称
	functionName := "swap-helper"                                   // 函数名称
	nonce := big.NewInt(int64(account.NonceInfo.PossibleNextNonce)) // nonce
	fee := big.NewInt(int64(13001))                                 // 手续费 5001 = 0.005001
	// 使用指针类型的 tokenX 和 tokenY
	tokenX, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2")
	tokenY, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc")

	fmt.Println(nonce)
	functionArgs := []stx.ClarityValue{
		&tokenX,                              // 使用指针类型的 token-x-trait
		&tokenY,                              // 使用指针类型的 token-y-trait
		stx.NewUintCV(big.NewInt(100000000)), // factor
		//stx.NewUintCV(big.NewInt(5000000000)), // dx 5000000000 == 50
		stx.NewUintCV(big.NewInt(1000000000)), // dx 5000000000 == 50
		//&stx.SomeCV{stx.OptionalSome, stx.NewUintCV(big.NewInt(171005898411))}, // 使用指针类型的 min-dy
		//&stx.SomeCV{stx.OptionalSome, stx.NewUintCV(big.NewInt(1740000000000))}, // 175000000000 ==1750
		&stx.SomeCV{stx.OptionalSome, stx.NewUintCV(big.NewInt(1740000000000))}, // 175000000000 ==1750
	}
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
	}
	// 发起合约调用
	tx, err := stx.MakeContractCall1(txOption)
	v1.CheckErr(err)
	// 序列化交易
	txSerialized := hex.EncodeToString(stx.Serialize(*tx))
	txId := stx.Txid(*tx)
	// 广播交易
	err = stx.BroadcastTransaction(txSerialized)
	fmt.Println(txSerialized)
	fmt.Println(txId)
}

func TestReplicateTransaction2(t *testing.T) {
	// 设置交易参数
	senderKey := "c0cee26708dcd2d16ccb6dfde394090a1fafde7cec4fb083d740cbc3021b39bf01"
	contractAddress := "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM" // 使用大写字母
	contractName := "amm-pool-v2-01"                               // 合约名称
	functionName := "swap-helper"                                  // 函数名称
	nonce := big.NewInt(286)                                       // nonce
	fee := big.NewInt(50002)                                       // 手续费 0.05002

	// 使用指针类型的 tokenX 和 tokenY
	tokenX, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2")
	tokenY, _ := stx.NewContractPrincipalCV1("SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc")

	functionArgs := []stx.ClarityValue{
		&tokenX,                              // 使用指针类型的 token-x-trait
		&tokenY,                              // 使用指针类型的 token-y-trait
		stx.NewUintCV(big.NewInt(100000000)), // factor
		stx.NewUintCV(big.NewInt(10000000)),  // dx
		&stx.SomeCV{stx.OptionalSome, stx.NewUintCV(big.NewInt(1000000000))}, // 使用指针类型的 min-dy
	}

	// 创建交易选项
	txOption := &stx.SignedContractCallOptions{
		ContractAddress: contractAddress,
		ContractName:    contractName,
		FunctionName:    functionName,
		FunctionArgs:    functionArgs,
		SendKey:         senderKey,
		ValidateWithAbi: false,
		Fee:             *fee,
		Nonce:           *nonce,
		AnchorMode:      3,
	}
	// 发起合约调用
	tx, err := stx.MakeContractCall1(txOption)
	require.NoError(t, err)
	// 序列化交易
	txSerialized := hex.EncodeToString(stx.Serialize(*tx))
	txId := stx.Txid(*tx)
	// 输出结果
	fmt.Println("Serialized Transaction:", txSerialized)
	fmt.Println("Transaction ID:", txId)
	// 广播交易
	err = stx.BroadcastTransactionFromb(txSerialized)
	require.NoError(t, err, "Failed to broadcast transaction")
}

// https://nameless-autumn-scion.stacks-mainnet.quiknode.pro/d01ddfac943c034fad15701a1f74674da5760d1f
// https://nameless-autumn-scion.stacks-mainnet.quiknode.pro/d01ddfac943c034fad15701a1f74674da5760d1f
