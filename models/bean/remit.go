package bean

type NotRemit struct {
	CompanyName string `json:"company_name,omitempty"`
	Amount      float64 `json:"amount"`
	StartTime   int64 `json:"start_time"`
	EndTime     int64 `json:"end_time"`
}
