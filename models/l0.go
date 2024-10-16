package models

import (
	"github.com/golang-module/carbon/v2"
	"github.com/piupuer/go-helper/ms"
)

type L0 struct {
	ms.M
	L0Id                         uint            `gorm:"comment:user id(SysUser.Id)" json:"userId"`
	SOURCE_CHAIN                 string          `gorm:"comment:SOURCE_CHAIN" json:"SOURCE_CHAIN"`
	SOURCE_TRANSACTION_HASH      string          `gorm:"comment:SOURCE_TRANSACTION_HASH" json:"SOURCE_TRANSACTION_HASH"`
	DESTINATION_CHAIN            string          `gorm:"comment:DESTINATION_CHAIN" json:"DESTINATION_CHAIN"`
	DESTINATION_TRANSACTION_HASH string          `gorm:"comment:DESTINATION_TRANSACTION_HASH" json:"DESTINATION_TRANSACTION_HASH"`
	SENDER_WALLET                string          `gorm:"comment:SENDER_WALLET" json:"SENDER_WALLET"`
	SOURCE_TIMESTAMP_UTC         carbon.DateTime `gorm:"comment:SOURCE_TIMESTAMP_UTC" json:"SOURCE_TIMESTAMP_UTC"`
	PROJECT                      string          `gorm:"PROJECT" json:"PROJECT"`
	NATIVE_DROP_USD              string          `gorm:"NATIVE_DROP_USD" json:"NATIVE_DROP_USD"`
	STARGATE_SWAP_USD            string          `gorm:"comment:STARGATE_SWAP_USD" json:"STARGATE_SWAP_USD"`
}
