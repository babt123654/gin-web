package v1

//
//
//import (
//	"encoding/json"
//	"fmt"
//	v1 "gin-web/api/v1"
//	"gin-web/models"
//	"net/http"
//)
//
//func MakeXYKSwap(req models.XykSwapReq) {
//
//	api := "https://nameless-autumn-scion.stacks-mainnet.quiknode.pro/d01ddfac943c034fad15701a1fxxxxda5760d1f/extended/v1/address/SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.xyk-core-v-1-1/mempool"
//	currentBlockTime, _ := v1.GetLatestBlockHeight() // 当前区块时间，您需要根据实际情况获取
//	senderKey := "您的发送密钥"                            // 替换为实际的发送密钥
//
//	transactions, err := getPendingTransactions(api)
//	if err != nil {
//		fmt.Println("获取交易失败:", err)
//		return
//	}
//	getTargetTransaction(req)
//}
//func getPendingTransactions(api string) ([]models.XykTx, error) {
//	resp, err := http.Get(api)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	var response models.XykApi
//	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
//		return nil, err
//	}
//	return response.Results, nil
//}
//
//func getTargetTransaction(req models.XykSwapReq) {
//	afterBurnTime := make([]models.XykTx, 0)
//	for _, result1 := range req.Transactions {
//		if result1.ReceiptTime > req.CurrentBlockTime {
//			afterBurnTime = append(afterBurnTime, result1)
//		}
//	}
//	for _, tx := range afterBurnTime {
//		if tx.ContractCall.FunctionName == models.XYKXFY {
//			for _, arg := range tx.ContractCall.FunctionArgs {
//				if arg.Name == models.AEUSDT || arg.Name == models.XYKSTX {
//				}
//			}
//		}
//		if tx.ContractCall.FunctionName == models.XYKYFX {
//
//		}
//		== req.DxToken {
//
//				// 检查交易提交时间是否在上一区块后当前区块出块前
//				if tx.ReceiptTime < currentBlockTime {
//				// 这里可以调用广播交易的函数
//				fmt.Printf("找到交易: %s，准备广播...\n", tx.TxID)
//			broadcastTransactionForXyk(tx, senderKey) // 实现广播交易的逻辑
//			}
//		}
//	}
//
//}
//func broadcastTransactionForXyk(tx models.XykTx, senderKey string) {
//	// 创建交易选项
//	txOption := &stx.SignedContractCallOptions{
//		ContractAddress: "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM",
//		ContractName:    "amm-pool-v2-01",
//		FunctionName:    "swap-helper",
//		FunctionArgs:    []stx.ClarityValue{ /* 填入参数 */ },
//		SendKey:         senderKey,
//		ValidateWithAbi: false,
//		//Fee:               big.NewInt(int64(tx.FeeRate)), // 根据实际情况设置手续费
//		//Nonce:             big.NewInt(int64(tx.Nonce)),
//		AnchorMode:        3,
//		PostConditionMode: stx.PostConditionModeAllow,
//	}
//
//	// 发起合约调用
//	_, err := stx.MakeContractCall1(txOption)
//	if err != nil {
//		fmt.Println("交易广播失败:", err)
//	} else {
//		fmt.Println("交易广播成功:", tx.TxID)
//	}
//}
