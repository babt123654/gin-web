package request

import (
	"github.com/golang-module/carbon/v2"
	"github.com/piupuer/go-helper/pkg/resp"
)

type L0 struct {
	L0Id                         uint            `json:"l0Id" swaggerignore:"true"`
	SOURCE_CHAIN                 string          `json:"SOURCE_CHAIN"`
	SOURCE_TRANSACTION_HASH      string          `json:"SOURCE_TRANSACTION_HASH"`
	DESTINATION_CHAIN            string          `json:"DESTINATION_CHAIN"`
	DESTINATION_TRANSACTION_HASH string          `json:"DESTINATION_TRANSACTION_HASH"`
	SENDER_WALLET                string          `json:"SENDER_WALLET"`
	SOURCE_TIMESTAMP_UTC         carbon.DateTime `json:"SOURCE_TIMESTAMP_UTC"`
	PROJECT                      string          `json:"PROJECT"`
	NATIVE_DROP_USD              string          `json:"NATIVE_DROP_USD"`
	STARGATE_SWAP_USD            string          `json:"STARGATE_SWAP_USD"`
	resp.Page
}

type CreateL0 struct {
	L0Id                         uint            `json:"l0Id" swaggerignore:"true"`
	SOURCE_CHAIN                 string          `json:"SOURCE_CHAIN"`
	SOURCE_TRANSACTION_HASH      string          `json:"SOURCE_TRANSACTION_HASH"`
	DESTINATION_CHAIN            string          `json:"DESTINATION_CHAIN"`
	DESTINATION_TRANSACTION_HASH string          `json:"DESTINATION_TRANSACTION_HASH"`
	SENDER_WALLET                string          `json:"SENDER_WALLET"`
	SOURCE_TIMESTAMP_UTC         carbon.DateTime `json:"SOURCE_TIMESTAMP_UTC"`
	PROJECT                      string          `json:"PROJECT"`
	NATIVE_DROP_USD              string          `json:"NATIVE_DROP_USD"`
	STARGATE_SWAP_USD            string          `json:"STARGATE_SWAP_USD"`
}

func (s CreateL0) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Desc"] = "description"
	return m
}

type UpdateL0 struct {
	L0Id                         uint            `json:"l0Id" swaggerignore:"true"`
	SOURCE_CHAIN                 string          `json:"SOURCE_CHAIN"`
	SOURCE_TRANSACTION_HASH      string          `json:"SOURCE_TRANSACTION_HASH"`
	DESTINATION_CHAIN            string          `json:"DESTINATION_CHAIN"`
	DESTINATION_TRANSACTION_HASH string          `json:"DESTINATION_TRANSACTION_HASH"`
	SENDER_WALLET                string          `json:"SENDER_WALLET"`
	SOURCE_TIMESTAMP_UTC         carbon.DateTime `json:"SOURCE_TIMESTAMP_UTC"`
	PROJECT                      string          `json:"PROJECT"`
	NATIVE_DROP_USD              string          `json:"NATIVE_DROP_USD"`
	STARGATE_SWAP_USD            string          `json:"STARGATE_SWAP_USD"`
}

type ApproveL0 struct {
	L0Id                         uint            `json:"l0Id" swaggerignore:"true"`
	SOURCE_CHAIN                 string          `json:"SOURCE_CHAIN"`
	SOURCE_TRANSACTION_HASH      string          `json:"SOURCE_TRANSACTION_HASH"`
	DESTINATION_CHAIN            string          `json:"DESTINATION_CHAIN"`
	DESTINATION_TRANSACTION_HASH string          `json:"DESTINATION_TRANSACTION_HASH"`
	SENDER_WALLET                string          `json:"SENDER_WALLET"`
	SOURCE_TIMESTAMP_UTC         carbon.DateTime `json:"SOURCE_TIMESTAMP_UTC"`
	PROJECT                      string          `json:"PROJECT"`
	NATIVE_DROP_USD              string          `json:"NATIVE_DROP_USD"`
	STARGATE_SWAP_USD            string          `json:"STARGATE_SWAP_USD"`
}
