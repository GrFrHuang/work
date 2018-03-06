// 阶梯分层相关

package old_sys

import (
	"encoding/json"
	"strconv"
	"fmt"
)

// 只关心有注释的部分
type Ladder4Get struct {
	Id            int `json:"id,omitempty"`            //
	GameId        int `json:"game_id,omitempty"`       // game id
	Cp            string `json:"cp,omitempty"`         // channel code
	StartTime     string `json:"start_time,omitempty"` // 阶梯生效开始时间 2017-01-01
	EndTime       string `json:"end_time,omitempty"`   // 阶梯生效截止时间 2017-01-01
	Ratio         float64 `json:"ratio"`               // 我方比例
	SlottingFee   float64 `json:"slotting_fee"`        // 通道费
	Rule          string `json:"rule"`                 // 规则
	UpdateTime    int `json:"update_time,omitempty"`   // 更新时间
	IsTrue        bool `json:"-"`                      //
	Status        int `json:"-"`                       //
	IsDelete      bool `json:"-"`                      //
	CalculateType int `json:"-"`                       //
	Interest      float64 `json:"-"`                   //
}

type Ladder4Post struct {
	Id          int `json:"id"`               // 为0 则为新增
	StartTime   string `json:"start_time"`    //
	EndTime     string `json:"end_time"`      //
	Ratio       float64 `json:"ratio"`        // 我方比例
	SlottingFee float64 `json:"slotting_fee"` // 通道费
	Rule        string `json:"rule"`          // 规则，使用&分割条件，目前只支持user、time、money条件类型，数值与条件类型间必须由"<"分割
}

type Clearing struct {
	Total       float64 `json:"total"`       // 总流水
	DivideTotal float64 `json:"divideTotal"` // 我方所得分成
}

type Clearings struct {
	GameId      string  `json:"game_id"`
	Total       float64 `json:"total"`       // 总流水
	DivideTotal float64 `json:"divideTotal"` // 我方所得分成
}

func GetLadder(gameId int, channelCode string) (ls []Ladder4Get, err error) {
	params := map[string]string{
		"game_id": strconv.Itoa(gameId),
	}
	if channelCode != "" {
		params["cp"] = channelCode
	}
	ls = []Ladder4Get{}
	err = request("/forgoapi/allotrule/getlist", params, &ls)
	if err != nil {
		return
	}

	return
}

func UpdateOrAddLadderList(gameId int, channelCode string, ladders []Ladder4Post) (err error) {

	js, err := json.Marshal(ladders)
	if err != nil {
		return
	}
	s := string(js)


	params := map[string]string{
		"game_id": strconv.Itoa(gameId),
		"data":    s,
	}
	if channelCode != "" {
		params["cp"] = channelCode
	}

	err = request("/forgoapi/allotrule/update", params, nil)
	if err != nil {
		return
	}

	return
}

// 计算
// month : Y-m
func GetClearing(gameId int, channelCode string, month string) (c *Clearing, err error) {
	params := map[string]string{
		"game_id": strconv.Itoa(gameId),
		"month":   month,
	}
	if channelCode != "" {
		params["cp"] = channelCode
	}

	err = request("/forgoapi/order/clearing", params, &c)

	if err != nil {
		return
	}

	return
}

func GetAllClearing(gameIds string, channelCode string, month string) (c *[]Clearings, err error) {
	params := map[string]string{
		"game_ids": gameIds,
		"month":   month,
	}
	if channelCode != "" {
		params["cp"] = channelCode
	}

	err = request("/forgoapi/order/clearingall", params, &c)

	if err != nil {
		return
	}

	return
}

func GetGameIncome(game_channel string) (c *[]Clearings, err error) {
	params := map[string]string{
		"gc": game_channel,
	}

	err = request("/forgoapi/order/clearingall", params, &c)
	fmt.Println(c)

	if err != nil {
		return
	}

	return
}
