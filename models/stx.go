package models

import (
	"math/big"
	"time"
)

const (
	//代币前缀
	TOKENPRE = "'"
	//代币
	ALEXSTX = "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wstx-v2"
	ABTC    = "SP2XD7417HGPRTREMKF748VNEQPDRR0RMANB7X1NK.token-abtc"
	NYC     = "SP102V8P0F7JX67ARQ77WEA3D3CFB5XW39REDT0AM.token-wnyc"
	SUSDT   = "SP2XD7417HGPRTREMKF748VNEQPDRR0RMANB7X1NK.token-susdt"
	AEUSDT  = "SP3Y2ZSH8P7D50B0VBTSX11S7XSG24M1VB9YFQA4K.token-aeusdc"
)

type BlockResponse struct {
	Nodes map[string]interface{} `json:"nodes"`
}
type Tokens struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	SwapContract      string  `json:"swapContract"`
	TotalSupply       int     `json:"totalSupply"`
	ContractID        string  `json:"contractId"`
	Symbol            string  `json:"symbol"`
	Price             float64 `json:"price"`
	Txns              int     `json:"txns"`
	Volume            float64 `json:"volume"`
	LastBlockTxn      int     `json:"lastBlockTxn"`
	SwapContractImage string  `json:"swapContractImage"`
	HrVolume          int     `json:"hrVolume"`
	DayVolume         int     `json:"dayVolume"`
	WeekVolume        float64 `json:"weekVolume"`
	HrChange          int     `json:"hrChange"`
	DayChange         int     `json:"dayChange"`
	WeekChange        float64 `json:"weekChange"`
}

type ArkaResult struct {
	BlockHeight                        int        `json:"block_height"`
	WrappedStxTokenUsdaToken           ArkaTokens `json:"wrapped-stx-token/usda-token"`
	ArkadikoTokenUsdaToken             ArkaTokens `json:"arkadiko-token/usda-token"`
	WrappedStxTokenArkadikoToken       ArkaTokens `json:"wrapped-stx-token/arkadiko-token"`
	WrappedStxTokenWrappedBitcoin      ArkaTokens `json:"wrapped-stx-token/Wrapped-Bitcoin"`
	WrappedBitcoinUsdaToken            ArkaTokens `json:"Wrapped-Bitcoin/usda-token"`
	WrappedStxTokenWelshcorgicoinToken ArkaTokens `json:"wrapped-stx-token/welshcorgicoin-token"`
	WrappedLydianTokenUsdaToken        ArkaTokens `json:"wrapped-lydian-token/usda-token"`
	LydianTokenUsdaToken               ArkaTokens `json:"lydian-token/usda-token"`
	WxusdWusda                         ArkaTokens `json:"wxusd/wusda"`
}

// arka
type ArkaTokens struct {
	ID               int       `json:"id"`
	TokenXName       string    `json:"token_x_name"`
	TokenYName       string    `json:"token_y_name"`
	TokenXAddress    string    `json:"token_x_address"`
	TokenYAddress    string    `json:"token_y_address"`
	SwapTokenName    string    `json:"swap_token_name"`
	SwapTokenAddress string    `json:"swap_token_address"`
	TvlTokenX        int64     `json:"tvl_token_x"`
	TvlTokenY        int64     `json:"tvl_token_y"`
	TvlUpdatedAt     time.Time `json:"tvl_updated_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	TokenYID         int       `json:"token_y_id"`
	TokenXID         int       `json:"token_x_id"`
	SwapTokenID      int       `json:"swap_token_id"`
	BalanceX         int64     `json:"balance_x"`
	BalanceY         int64     `json:"balance_y"`
	SharesTotal      int64     `json:"shares_total"`
	Enabled          bool      `json:"enabled"`
}

// stake
type StakeSwap struct {
	Type          string            `json:"type"`
	From          string            `json:"from"`
	To            string            `json:"to"`
	FromDirection bool              `json:"from_direction"`
	ToDirection   bool              `json:"to_direction"`
	FromPair      StakeSwapFromPair `json:"from_pair"`
	ToPair        StakeSwapToPair   `json:"to_pair"`
	BridgeToken   string            `json:"bridge_token"`
	BlockHeight   string            `json:"block_height"`
}

type StakeSwapFromPair struct {
	BalanceX       string `json:"balance-x"`
	BalanceY       string `json:"balance-y"`
	FeeBalanceX    string `json:"fee-balance-x"`
	FeeBalanceY    string `json:"fee-balance-y"`
	FeeToAddress   string `json:"fee-to-address"`
	LiquidityToken string `json:"liquidity-token"`
	SharesTotal    string `json:"shares-total"`
	Name           string `json:"name"`
}

type StakeSwapToPair struct {
	BalanceX       string `json:"balance-x"`
	BalanceY       string `json:"balance-y"`
	FeeBalanceX    string `json:"fee-balance-x"`
	FeeBalanceY    string `json:"fee-balance-y"`
	FeeToAddress   string `json:"fee-to-address"`
	LiquidityToken string `json:"liquidity-token"`
	SharesTotal    string `json:"shares-total"`
	Name           string `json:"name"`
}
type Welsh struct {
	StxPrice float64 `json:"stxPrice"`
	Price    float64 `json:"price"`
	Com      string  `json:"com"`
}

// sUsdt 监控
type SUsdt struct {
	Orders                      []Order `json:"orders"`
	StacksBridgeContractID      string  `json:"stacks_bridge_contract_id"`
	StacksCrossPegInContractID  string  `json:"stacks_cross_peg_in_contract_id"`
	StacksCrossPegOutContractID string  `json:"stacks_cross_peg_out_contract_id"`
}
type Order struct {
	SentUserOnchainStatus string `json:"sent_user_onchain_status"`
}

// alex tokens
type AlexTokenPrice struct {
	AvgPriceUSD float64 `json:"avg_price_usd"`
	Token       string  `json:"token"`
}
type AlexResponseData struct {
	Data struct {
		LaplaceCurrentTokenPrice []AlexTokenPrice `json:"laplace_current_token_price"`
	} `json:"data"`
}

// token-alex新
type AlexTokenPrice2 struct {
	Data []Data `json:"data"`
}
type Data struct {
	ContractID   string  `json:"contract_id"`
	LastPriceUsd float64 `json:"last_price_usd"`
}

// velar
type VelarResponseData struct {
	StatusCode int               `json:"statusCode"`
	Message    []VelarTokenPrice `json:"message"`
}
type VelarTokenPrice struct {
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
}

// 定义返回值的结构体
type BlackInfo struct {
	TxID                   string     `json:"tx_id"`
	Nonce                  int        `json:"nonce"`
	FeeRate                string     `json:"fee_rate"`
	SenderAddress          string     `json:"sender_address"`
	Sponsored              bool       `json:"sponsored"`
	PostConditionMode      string     `json:"post_condition_mode"`
	PostConditions         []struct{} `json:"post_conditions"` // 这里可以根据需要定义具体结构
	AnchorMode             string     `json:"anchor_mode"`
	IsUnanchored           bool       `json:"is_unanchored"`
	BlockHash              string     `json:"block_hash"`
	ParentBlockHash        string     `json:"parent_block_hash"`
	BlockHeight            int        `json:"block_height"`
	BlockTime              int64      `json:"block_time"`
	BlockTimeISO           string     `json:"block_time_iso"`
	BurnBlockHeight        int        `json:"burn_block_height"`
	BurnBlockTime          int64      `json:"burn_block_time"`
	BurnBlockTimeISO       string     `json:"burn_block_time_iso"`
	ParentBurnBlockTime    int64      `json:"parent_burn_block_time"`
	ParentBurnBlockTimeISO string     `json:"parent_burn_block_time_iso"`
	Canonical              bool       `json:"canonical"`
	TxIndex                int        `json:"tx_index"`
	TxStatus               string     `json:"tx_status"`
}

// 区块浏览器pending返回值部分
type PendingTx struct {
	Limit            int                `json:"limit"`
	Offset           int                `json:"offset"`
	Total            int                `json:"total"`
	PendingTxResults []PendingTxResults `json:"results"`
}
type PendingTxPrincipal struct {
	TypeID       string `json:"type_id"`
	Address      string `json:"address"`
	ContractName string `json:"contract_name"`
}
type PendingTxAsset struct {
	AssetName       string `json:"asset_name"`
	ContractAddress string `json:"contract_address"`
	ContractName    string `json:"contract_name"`
}
type PendingTxPostConditions struct {
	PendingTxPrincipal PendingTxPrincipal `json:"principal,omitempty"`
	ConditionCode      string             `json:"condition_code"`
	Amount             string             `json:"amount"`
	Type               string             `json:"type"`
	//Principal     Principal `json:"principal,omitempty"`
	PendingTxAsset Asset `json:"asset,omitempty"`
}
type PendingTxFunctionArgs struct {
	Hex  string `json:"hex"`
	Repr string `json:"repr"`
	Name string `json:"name"`
	Type string `json:"type"`
}
type PendingTxContractCall struct {
	ContractID        string                  `json:"contract_id"`
	FunctionName      string                  `json:"function_name"`
	FunctionSignature string                  `json:"function_signature"`
	FunctionArgs      []PendingTxFunctionArgs `json:"function_args"`
}
type PendingTxResults struct {
	TxID              string                    `json:"tx_id"`
	Nonce             int                       `json:"nonce"`
	FeeRate           string                    `json:"fee_rate"`
	SenderAddress     string                    `json:"sender_address"`
	Sponsored         bool                      `json:"sponsored"`
	PostConditionMode string                    `json:"post_condition_mode"`
	PostConditions    []PendingTxPostConditions `json:"post_conditions"`
	AnchorMode        string                    `json:"anchor_mode"`
	TxStatus          string                    `json:"tx_status"`
	ReceiptTime       int64                     `json:"receipt_time"`
	ReceiptTimeIso    time.Time                 `json:"receipt_time_iso"`
	TxType            string                    `json:"tx_type"`
	ContractCall      PendingTxContractCall     `json:"contract_call"`
}

// 交易部分
type StacksTransferSig struct {
	TokenX       string
	TokenY       string
	AmountIn     *big.Int
	MinAmountOut *big.Int
	Fees         *big.Int
	Nonce        uint64
}

// NYC套利部分
type BalanceAPI1 struct {
	Balance      string `json:"balance"`
	Locked       string `json:"locked"`
	UnlockHeight int    `json:"unlock_height"`
	Nonce        int    `json:"nonce"`
}
type Stx struct {
	Balance                   string `json:"balance"`
	TotalSent                 string `json:"total_sent"`
	TotalReceived             string `json:"total_received"`
	TotalFeesSent             string `json:"total_fees_sent"`
	TotalMinerRewardsReceived string `json:"total_miner_rewards_received"`
	Locked                    string `json:"locked"`
}
type FungibleToken struct {
	Balance       string `json:"balance"`
	TotalSent     string `json:"total_sent"`
	TotalReceived string `json:"total_received"`
}
type NonFungibleToken struct {
	Count         string `json:"count"`
	TotalSent     string `json:"total_sent"`
	TotalReceived string `json:"total_received"`
}
type BalanceAPI2 struct {
	STX               Stx                         `json:"stx"`
	FungibleTokens    map[string]FungibleToken    `json:"fungible_tokens"`
	NonFungibleTokens map[string]NonFungibleToken `json:"non_fungible_tokens"`
}

type CombinedBalance struct {
	API1      BalanceAPI1 `json:"api1"`
	API2      BalanceAPI2 `json:"api2"`
	NonceInfo NonceInfo   `json:"NonceInfo"`
}

// nonce --- https://api.hiro.so/extended/v1/address/地址/nonces
type NonceInfo struct {
	LastMempoolTxNonce    *int  `json:"last_mempool_tx_nonce"`
	LastExecutedTxNonce   int   `json:"last_executed_tx_nonce"`
	PossibleNextNonce     int   `json:"possible_next_nonce"`
	DetectedMissingNonces []int `json:"detected_missing_nonces"`
	DetectedMempoolNonces []int `json:"detected_mempool_nonces"`
}

type JSONData struct {
}

// sechaPendinTx
type NycBotReqStruct struct {
	SenderKey          string           `json:"senderKey"`
	ContractAddress    string           `json:"contractAddress"`
	ContractName       string           `json:"contractName"`
	FunctionName       string           `json:"functionName"`
	SubmitBlock        int              `json:"submitBlock"`
	SwapContract       []string         `json:"swapContract"`
	TokenContract      []string         `json:"TokenContract"`
	Txs                []string         `json:"txs"`
	Account            *CombinedBalance `json:"account"`
	NowBlock           int64            `json:"NowBlock"`
	RedeemBlock        int64            `json:"redeemBlock"`
	Redeem             bool             `json:"redeem"`
	Api                string           `json:"api"`
	Submiting          bool             `json:"Submiting"`
	NycAmount          float64          `json:"nycAmount"`
	SubmitTx           string           `json:"submitTx"`
	SubmitTxSerialized string           `json:"submitTxSerialized"`
}
type RedeemNycReq struct {
	SenderKey       string           `json:"senderKey"`
	ContractAddress string           `json:"contractAddress"`
	ContractName    string           `json:"contractName"`
	FunctionName    string           `json:"functionName"`
	SwapContract    []string         `json:"swapContract"`
	TokenContract   []string         `json:"TokenContract"`
	Txs             []string         `json:"txs"`
	Account         *CombinedBalance `json:"account"`
	NowBlock        int64            `json:"NowBlock"`
	SubmitTx        bool             `json:"submitTx"`
	SubmitBlock     int              `json:"submitBlock"`
	Redeem          bool             `json:"redeem"`
}

// 压单程序
type SubmitTxFromAddrReq struct {
	SenderKey     string           `json:"senderKey"`
	TargetAddress string           `json:"TargetAddress"`
	SubmitModel   int              `json:"SubmitModel"`
	StartBlock    int64            `json:"startBlock"`
	Api           string           `json:"api"`
	ApiPre        string           `json:"apiPre"`
	ApiFix        string           `json:"apiFix"`
	Account       *CombinedBalance `json:"account"`
	MyTx          string           `json:"myTx"`
}

// 待签名tx
type GetSubmitTxByStructReq struct {
	TxID              string                    `json:"tx_id"`
	Nonce             int                       `json:"nonce"`
	FeeRate           string                    `json:"fee_rate"`
	SenderAddress     string                    `json:"sender_address"`
	Sponsored         bool                      `json:"sponsored"`
	PostConditionMode string                    `json:"post_condition_mode"`
	PostConditions    []PendingTxPostConditions `json:"post_conditions"`
	AnchorMode        string                    `json:"anchor_mode"`
	TxStatus          string                    `json:"tx_status"`
	ReceiptTime       int64                     `json:"receipt_time"`
	ReceiptTimeIso    time.Time                 `json:"receipt_time_iso"`
	TxType            string                    `json:"tx_type"`
	ContractCall      PendingTxContractCall     `json:"contract_call"`
}

// 解析 JSON 数据到自定义结构体
type PostConditionJSON struct {
	Principal     PostConditionPrincipalJSON `json:"principal"`
	ConditionCode string                     `json:"condition_code"`
	Amount        string                     `json:"amount"`
	Type          string                     `json:"type"`
}

type PostConditionPrincipalJSON struct {
	TypeID       string `json:"type_id"`
	Address      string `json:"address"`
	ContractName string `json:"contract_name,omitempty"`
}

// for --changefee
type ChangeFee struct {
	TargetTxStruct PendingTxResults `json:"targetTx"`
	Account        CombinedBalance  `json:"account"`
	SenderKey      string           `json:"senderKey"`
	ChangeTxStruct PendingTxResults `json:"ChangeTxStruct"`
}
