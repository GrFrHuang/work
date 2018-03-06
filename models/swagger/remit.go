package swagger

type InputRemitDownAccount struct {
	ChannelId int `json:"channel_id"`
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
	Amount    float64 `json:"amount"`
}
