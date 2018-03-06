package controllers

import (
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"strings"
)

// cp预对账
type OrderPreVerifyCpController struct {
	BaseController
}

// URLMapping ...
func (c *OrderPreVerifyCpController) URLMapping() {
}

// @Title 刷新cp预对账单
// @Description 刷新cp预对账单
// @Param	company_ids	query	string	false	"指定要更新的发行商,为空则全部更新"
// @Param	month	query	string	false	"指定要更新的月份,为空则全部更新"
// @Success 200 string "count:1(更新的游戏数量)"
// @Failure 403
// @router / [get]
func (c *OrderPreVerifyCpController) UpdateFromOrder() {
	companyStr := c.GetString("company_ids")
	month := c.GetString("month")
	if month == "" {
		c.RespJSON(bean.CODE_Params_Err, "month can't be empty")
		return
	}
	companyIds := []interface{}{}
	if companyStr != "" {
		for _, v := range strings.Split(companyStr, ",") {
			companyIds = append(companyIds, v)
		}
	}
	aff, err := models.UpdatePreVerifyCpFromOrder(companyIds, month)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(aff)
}
