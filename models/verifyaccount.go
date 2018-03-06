package models

type NoCpVerify struct {
	Date    string `json:"date,omitempty"`
	Company *CompanyType `json:"company,omitempty"`
	BodyMy  int `json:"body_my,omitempty"`
	Amount  float64 `json:"amount"`
}

type NoChannelVerify struct {
	Date    string `json:"date,omitempty"`
	Channel *Channel `json:"channel,omitempty"`
	BodyMy  int `json:"body_my,omitempty"`
	Amount  float64 `json:"amount"`
}
