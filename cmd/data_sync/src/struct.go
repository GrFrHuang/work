package src

type Response struct {
	Status        string
	Error         string
	Error_message string
	//dat interface{} // []map[string]interface{} , []interface, map[string]interface{}
}

type Game struct {
	Game_id   int
	Name      string
	Icon      string
	Show_icon string
}

type GameResponse struct {
	//Status        string
	//Error         string
	//Error_message string
	Response
	Data          []Game
}

type Channel struct {
	Id           int
	Cp           string
	Channel_name string
	Platform     string
}

type ChannelResponse struct {
	Response
	Data []Channel
}

type Order struct {
	Id           int
	Game_id      int
	Cp           string
	Date         string
	Amount       string
	Order_status int // 0:通知游戏成功 2:失败; 3:金额不符 4:通知游戏失败 5:退款 6:作废
	Update_time  int64
}

type OrderResponse struct {
	Response
	Data []Order
}
