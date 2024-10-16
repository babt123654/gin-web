package request

import (
	"gin-web/models"
	"github.com/golang-module/carbon/v2"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

type CreateVelar struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
	//Results    []Results `json:"results"`
	StatusCode int       `json:"statusCode"`
	Message    []Message `gorm:"foreignKey:MessageId" json:"message"`
}

type Message struct {
	Order            int    `json:"order"`
	Symbol           string `json:"symbol"`
	Name             string `json:"name"`
	ContractAddress  string `json:"contractAddress"`
	ImageURL         string `json:"imageUrl"`
	Price            string `json:"price"`
	Decimal          string `json:"decimal"`
	TokenDecimalNum  int    `json:"tokenDecimalNum"`
	AssetName        string `json:"assetName"`
	PercentChange24H string `json:"percent_change_24h"`
	Vsymbol          string `json:"vsymbol"`
	resp.Page
}

type Velar struct {
	UserId          uint          `json:"-"`
	Status          *req.NullUint `json:"status" form:"status"`
	ApprovalOpinion string        `json:"approvalOpinion" form:"approvalOpinion"`
	Desc            string        `json:"desc" form:"desc"`
	resp.Page
}

func (s CreateVelar) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Desc"] = "description"
	return m
}

type UpdateVelar struct {
	Desc      *string          `json:"desc"`
	StartTime *carbon.DateTime `json:"startTime" swaggertype:"string"`
	EndTime   *carbon.DateTime `json:"endTime" swaggertype:"string"`
}

type ApproveVelar struct {
	Id          uint           `json:"id"`
	AfterStatus uint           `json:"afterStatus"`
	Approved    uint           `json:"approved"`
	User        models.SysUser `json:"user" swaggerignore:"true"`
}
