package controllers

import (
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
)

// 部门统计
type StatisticController struct {
	BaseController
}

// URLMapping ...
func (c *StatisticController) URLMapping() {

}

// @Title 结算部统计
// @Description 结算部统计
// @Param	start		query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  		query	int64	false	"The timestamp of end. Must be an integer"
// @Param	limit   	query	int64	false	"Limit the size of result set. Must be an integer"
// @Param	offset  	query	int64	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []*StatisticAccounting
// @Failure 403
// @router /accounting [get]
func (c *StatisticController) GetStatisticAccounting() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_ACCOUNTING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	var limit int
	var offset int
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	} else {
		start = time.Now().AddDate(0, -3, 0).Unix()
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	} else {
		end = time.Now().Unix()
	}
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	} else {
		limit = 20
	}
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	} else {
		offset = 0
	}

	result, total := models.StatisticOfAccounting(start, end, limit, offset)
	c.RespJSONDataWithTotal(result, total)
}

// @Title 财务部统计
// @Description 财务部统计
// @Param	start		query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  		query	int64	false	"The timestamp of end. Must be an integer"
// @Param	limit   	query	int64	false	"Limit the size of result set. Must be an integer"
// @Param	offset  	query	int64	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []*StatisticFinancial
// @Failure 403
// @router /financial [get]
func (c *StatisticController) GetStatisticFinancial() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_FINANCIAL, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	var limit int
	var offset int
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	} else {
		start = time.Now().AddDate(0, -3, 0).Unix()
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	} else {
		end = time.Now().Unix()
	}
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	} else {
		limit = 20
	}
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	}

	result, total := models.StatisticOfFinancial(start, end, limit, offset)
	c.RespJSONDataWithTotal(result, total)
}

// @Title 渠道商务部统计
// @Description 渠道商务统计
// @Param	start		query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  		query	int64	false	"The timestamp of end. Must be an integer"
// @Param	limit   	query	int64	false	"Limit the size of result set. Must be an integer"
// @Param	offset  	query	int64	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []*StatisticChannelTrade
// @Failure 403
// @router /channelTrade [get]
func (c *StatisticController) GetStatisticChannelTrade() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CHANNEL_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	var limit int
	var offset int
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	} else {
		start = time.Now().AddDate(0, -3, 0).Unix()
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	} else {
		end = time.Now().Unix()
	}
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	} else {
		limit = 20
	}
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	}

	result, total := models.StatisticOfChannelTrade(start, end, limit, offset)
	c.RespJSONDataWithTotal(result, total)
}

// @Title CP商务部统计
// @Description CP商务统计
// @Param	start		query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  		query	int64	false	"The timestamp of end. Must be an integer"
// @Param	limit   	query	int64	false	"Limit the size of result set. Must be an integer"
// @Param	offset  	query	int64	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []*StatisticCpTrade
// @Failure 403
// @router /cpTrade [get]
func (c *StatisticController) GetStatisticCpTrade() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CHANNEL_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	var limit int
	var offset int
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	} else {
		start = time.Now().AddDate(0, -3, 0).Unix()
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	} else {
		end = time.Now().Unix()
	}
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	} else {
		limit = 20
	}
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	} else {
		offset = 0
	}

	result, total := models.StatisticOfCpTrade(start, end, limit, offset)
	c.RespJSONDataWithTotal(result, total)
}

// @Title 运营部统计
// @Description 运营部统计
// @Param	start		query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  		query	int64	false	"The timestamp of end. Must be an integer"
// @Param	limit   	query	int64	false	"Limit the size of result set. Must be an integer"
// @Param	offset  	query	int64	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /operation [get]
func (c *StatisticController) GetStatisticOperation() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CHANNEL_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	var limit int
	var offset int
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	} else {
		start = time.Now().AddDate(0, -3, 0).Unix()
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	} else {
		end = time.Now().Unix()
	}
	if v, err := c.GetInt("limit"); err == nil {
		limit = v
	} else {
		limit = 20
	}
	if v, err := c.GetInt("offset"); err == nil {
		offset = v
	}

	result, total := models.StatisticOfOperation(start, end, limit, offset)
	c.RespJSONDataWithTotal(result, total)
}

// GetAll ...
// @Title 渠道寄出合同
// @Description 渠道寄出合同
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /channelContractSend [get]
func (c *StatisticController) GetChannelSendContract() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_ACCOUNTING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfChannelSendContract(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title 渠道归档合同
// @Description 渠道归档合同
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /channelContractComplete [get]
func (c *StatisticController) GetChannelCompleteContract() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_ACCOUNTING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfChannelCompleteContract(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title 渠道对账数
// @Description 渠道对账数
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /channelVerify [get]
func (c *StatisticController) GetChannelVerify() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_ACCOUNTING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfChannelVerify(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title CP寄出合同数
// @Description CP寄出合同数
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /cpContractComplete [get]
func (c *StatisticController) GetCpSendContract() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_ACCOUNTING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfCpSendContract(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title CP归档合同
// @Description CP归档合同
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /cpContractComplete [get]
func (c *StatisticController) GetCpCompleteContract() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_ACCOUNTING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfCpCompleteContract(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title CP对账数
// @Description CP对账数
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /cpVerify [get]
func (c *StatisticController) GetCpVerify() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_ACCOUNTING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfCpVerify(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title 渠道回款
// @Description 渠道回款
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /channelPaid [get]
func (c *StatisticController) GetChannelPaid() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_FINANCIAL, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfChannelPaid(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title CP付款
// @Description CP付款
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /cpPaid [get]
func (c *StatisticController) GetCpPaid() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_FINANCIAL, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfCpPaid(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title 渠道接入
// @Description 渠道接入
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /channelAccess [get]
func (c *StatisticController) GetChannelAccess() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CHANNEL_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfChannelAccess(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title 渠道回款(渠道商务部)
// @Description 渠道回款(渠道商务部)
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /channelPaidForTrade [get]
func (c *StatisticController) GetChannelPaid2() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CHANNEL_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfChannelPaid(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title 渠道商信息
// @Description 渠道商信息
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /channelCompany [get]
func (c *StatisticController) GetChannelCompany() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CHANNEL_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfChannelCompany(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title 提测游戏
// @Description 提测游戏
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /gameEvaluate [get]
func (c *StatisticController) GetGameEvaluate() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CP_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfGameEvaluate(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title 接入游戏
// @Description 接入游戏
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /gameAccess [get]
func (c *StatisticController) GetGameAccess() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CP_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfGameAccess(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title CP付款(CP商务部)
// @Description 接入游戏(CP商务部)
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /cpPaidForTrade [get]
func (c *StatisticController) GetCpPaidForTrade() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CP_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfCpPaid(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}

// GetAll ...
// @Title 发行商信息
// @Description 发行商信息
// @Param	start	query	int64	false	"The timestamp of start. Must be an integer"
// @Param	end  	query	int64	false	"The timestamp of end. Must be an integer"
// @Success 200 {object} []*StatisticOperation
// @Failure 403
// @router /distributionCompany [get]
func (c *StatisticController) GetDistributionCompany() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_STATISTIC_CP_TRADE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var start int64
	var end int64
	if v, err := c.GetInt64("start"); err == nil {
		start = v
	}
	if v, err := c.GetInt64("end"); err == nil {
		end = v
	}

	result, err := models.StatisticOfDistributionCompany(start, end)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(result)
}
