package controllers

import (
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"strings"
)

// 渠道预对账
type OrderPreVerifyChannelController struct {
	BaseController
}

// URLMapping ...
func (c *OrderPreVerifyChannelController) URLMapping() {
}

// @Title 刷新渠道预对账单
// @Description 刷新渠道预对账单
// @Param	channel_ids	query	string	false	"指定要更新的渠道,为空则全部更新"
// @Param	month	query	string	false	"指定要更新的月份,为空则全部更新"
// @Success 200 string "count:1(更新的游戏数量)"
// @Failure 403
// @router / [get]
func (c *OrderPreVerifyChannelController) UpdateFromOrder() {
	channelStr := c.GetString("channel_ids")
	month := c.GetString("month")
	if month == "" {
		c.RespJSON(bean.CODE_Params_Err, "month can't be empty")
		return
	}
	channelCodes := []interface{}{}
	if channelStr != "" {
		for _, v := range strings.Split(channelStr, ",") {
			channelCodes = append(channelCodes, v)
		}
	}
	aff, err := models.UpdatePreVerifyChannelFromOrder(channelCodes, month)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(aff)
}
