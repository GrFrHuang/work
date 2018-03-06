package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool/old_sys"
)

// 阶梯分层
type LadderController struct {
	BaseController
}

// URLMapping ...
func (c *LadderController) URLMapping() {

}

// GetOne ...
// @Title 获取对上游的阶梯规则
// @Description 获取对上游的阶梯规则
// @Param	game_id		query 	string	true		"The key for staticblock"
// @Success 200 {object} tool.old_sys.Ladder4Get
// @router /cp [get]
func (c *LadderController) Get4Cp() {
	gameId, _ := c.GetInt("game_id", 0)

	ls, err := old_sys.GetLadder(gameId, "")
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(ls)
}

// GetOne ...
// @Title 获取对渠道的阶梯规则
// @Description 获取对渠道的阶梯规则
// @Param	game_id		query 	string	true		"游戏id"
// @Param	channel_code		query 	string	true		"渠道code"
// @Success 200 {object} models.AlarmLog
// @Failure 403 :id is empty
// @router /channel [get]
func (c *LadderController) Get4Channel() {
	gameId, _ := c.GetInt("game_id", 0)
	channelCode := c.GetString("channel_code", "")

	if channelCode == "" {
		c.RespJSON(bean.CODE_Params_Err, "channel_code can't be empty")
		return
	}

	ls, err := old_sys.GetLadder(gameId, channelCode)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(ls)
}

// @Title 添加对上游的阶梯规则
// @Description 添加对上游的阶梯规则
// @Param	game_id		formData 	string	true		"游戏id"
// @Param	data		formData 	tool.old_sys.Ladder4Post	true		""
// @Success 200 {object} models.AlarmLog
// @Failure 403 :id is not int
// @router /cp [put]
func (c *LadderController) Put4Cp() {
	gameId, _ := c.GetInt("game_id", 0)
	data := c.GetString("data", "")

	ls := []old_sys.Ladder4Post{}
	if err := json.Unmarshal([]byte(data), &ls); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	err := old_sys.UpdateOrAddLadderList(gameId, "", ls)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("OK")
}

// @Title 添加对渠道的阶梯规则
// @Description 添加对渠道的阶梯规则
// @Param	game_id		formData 	string	true		"游戏id"
// @Param	channel_code		formData 	string	true		"渠道code"
// @Param	data		formData 	tool.old_sys.Ladder4Post	true		""
// @Success 200 {object} models.AlarmLog
// @Failure 403 :id is not int
// @router /channel [put]
func (c *LadderController) Put4Channel() {
	gameId, _ := c.GetInt("game_id", 0)
	data := c.GetString("data", "")
	channelCode := c.GetString("channel_code", "")

	if channelCode == "" {
		c.RespJSON(bean.CODE_Params_Err, "channel_code can't be empty")
		return
	}

	ls := []old_sys.Ladder4Post{}
	if err := json.Unmarshal([]byte(data), &ls); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	err := old_sys.UpdateOrAddLadderList(gameId, channelCode, ls)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("OK")
}

// @Title 获取渠道阶梯计算出的金额
// @Description 获取渠道阶梯计算出的金额
// @Param	game_id		query 	string	true		"游戏id"
// @Param	channel_code		query 	string	true		"渠道code"
// @Param	month		query 	string	true		"Y-m"
// @Success 200 {object} tool.old_sys.Clearing
// @router /clearchannel [get]
func (c *LadderController) GetClearing4Channel() {
	gameId, _ := c.GetInt("game_id", 0)
	channelCode := c.GetString("channel_code", "")
	month := c.GetString("month", "")

	if channelCode == "" || month == "" {
		c.RespJSON(bean.CODE_Params_Err, "channel_code and month can't be empty")
		return
	}

	clearing, err := old_sys.GetClearing(gameId, channelCode, month)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(clearing)
}

// @Title 获取上游阶梯计算出的金额
// @Description 获取上游阶梯计算出的金额
// @Param	game_id		query 	string	true		"游戏id"
// @Param	month		query 	string	true		"Y-m"
// @Success 200 {object} tool.old_sys.Clearing
// @router /clearcp [get]
func (c *LadderController) GetClearing4Cp() {
	gameId, _ := c.GetInt("game_id", 0)
	month := c.GetString("month", "")

	if month == "" {
		c.RespJSON(bean.CODE_Params_Err, "channel_code and month can't be empty")
		return
	}

	clearing, err := old_sys.GetClearing(gameId, "", month)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(clearing)
}
