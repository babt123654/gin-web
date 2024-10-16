package models

import (
	"time"
)

type Principal struct {
	TypeID       string `json:"type_id"`
	ContractName string `json:"contract_name"`
	Address      string `json:"address"`
}

type Asset struct {
	ContractName    string `json:"contract_name"`
	AssetName       string `json:"asset_name"`
	ContractAddress string `json:"contract_address"`
}

type PostConditions struct {
	Type          string    `json:"type"`
	ConditionCode string    `json:"condition_code"`
	Amount        string    `json:"amount"`
	Asset         Asset     `json:"asset,omitempty"`
	Principal     Principal `json:"principal,omitempty"`
}

type FunctionArgs struct {
	Hex  string `json:"hex"`
	Repr string `json:"repr"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ContractCall struct {
	ContractID        string         `json:"contract_id"`
	FunctionName      string         `json:"function_name"`
	FunctionSignature string         `json:"function_signature"`
	FunctionArgs      []FunctionArgs `json:"function_args"`
}

type Results struct {
	TxID              string           `json:"tx_id"`
	Nonce             int              `json:"nonce"`
	FeeRate           string           `json:"fee_rate"`
	SenderAddress     string           `json:"sender_address"`
	Sponsored         bool             `json:"sponsored"`
	PostConditionMode string           `json:"post_condition_mode"`
	PostConditions    []PostConditions `json:"post_conditions"`
	AnchorMode        string           `json:"anchor_mode"`
	TxStatus          string           `json:"tx_status"`
	ReceiptTime       int              `json:"receipt_time"`
	ReceiptTimeIso    time.Time        `json:"receipt_time_iso"`
	TxType            string           `json:"tx_type"`
	ContractCall      ContractCall     `json:"contract_call"`
}
