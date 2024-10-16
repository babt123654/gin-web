package models

type LineaPoints struct {
	RankXP      int    `json:"rank_xp"`
	UserAddress string `json:"user_address"`
	XP          int    `json:"xp"`
	ALP         int    `json:"alp"`
	PLP         int    `json:"plp"`
	EP          int    `json:"ep"`
	RP          int    `json:"rp"`
	VP          int    `json:"vp"`
	BP          int    `json:"bp"`
	EAFlag      int    `json:"ea_flag"`
}
