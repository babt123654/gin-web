package models

import (
	"math/big"
	"time"
)

const (
	XYKSTX = "'SM1793C4R5PZ4NS4VQ4WMP7SKKYVH8JZEWSZ9HCCR.token-stx-v-1-1"
	AEUSDC = "'SP3Y2ZSH8P7D50B0VBTSX11S7XSG24M1VB9YFQA4K.token-aeusdc"
)
const (
	XYKYFX = "swap-y-for-x"
	XYKXFY = "swap-x-for-y"
	//Dx     = "x-token-trait"
	//Dx     = "y-token-trait"
)

// xyk方法请求体
type XykSwapReq struct {
	Api              string  `json:"api"`
	SenderKey        string  `json:"senderKey"`
	Transactions     []XykTx `json:"transactions"`
	CurrentBlockTime int     `json:"currentBlockTime"`
	DxToken          string  `json:"dxToken"`
	DyToken          string  `json:"dyToken"`
}

// xyk api相关结构体
type AEUSDC_STXreq struct {
	Aeusdc big.Int `json:"aeusdc"`
	Stx    big.Int `json:"stx"`
}

// XykApi 表示API返回的整体结构
type XykApi struct {
	Limit   int     `json:"limit"`
	Offset  int     `json:"offset"`
	Total   int     `json:"total"`
	Results []XykTx `json:"results"`
}

// XykFunctionArgs 表示函数参数的结构
type XykFunctionArgs struct {
	Hex  string `json:"hex"`
	Repr string `json:"repr"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// XykContractCall 表示合约调用的结构
type XykContractCall struct {
	ContractID        string            `json:"contract_id"`
	FunctionName      string            `json:"function_name"`
	FunctionSignature string            `json:"function_signature"`
	FunctionArgs      []XykFunctionArgs `json:"function_args"`
}

// XykTx 表示每个交易结果的结构
type XykTx struct {
	TxID              string          `json:"tx_id"`
	Nonce             int             `json:"nonce"`
	FeeRate           string          `json:"fee_rate"`
	SenderAddress     string          `json:"sender_address"`
	Sponsored         bool            `json:"sponsored"`
	PostConditionMode string          `json:"post_condition_mode"`
	PostConditions    []interface{}   `json:"post_conditions"`
	AnchorMode        string          `json:"anchor_mode"`
	TxStatus          string          `json:"tx_status"`
	ReceiptTime       int             `json:"receipt_time"`
	ReceiptTimeIso    time.Time       `json:"receipt_time_iso"`
	TxType            string          `json:"tx_type"`
	ContractCall      XykContractCall `json:"contract_call"`
}
