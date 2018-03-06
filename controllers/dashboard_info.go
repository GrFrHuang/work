package controllers

import (
	"kuaifa.com/kuaifa/work-together/appcache"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

// 用户自己的信息
type DashboardInfoController struct {
	BaseController
}

func (c *DashboardInfoController) URLMapping() {
	c.Mapping("Get", c.Get)
}

// @Title 获取仪表板信息
// @Description 获取仪表板信息
// @Success 200
// @Failure 403 body is empty
// @router / [get]
func (c *DashboardInfoController) Get() {
	uid := c.Uid()
	_, err := models.GetUserById(uid)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	var data = []models.DashboradInfo{
		{
			Title: "今日总流水",
			Tag:   "todayMoney",
			Pmsm:  bean.PMSM_ORDER,
			Href:  "/home/list",
		},
		{
			Title: "未对账金额",
			Tag:   "notReconMoney",
			Pmsm:  bean.PMSM_CHANNEL_VERIFY_ACCOUNT,
			Href:  "/home/channelContract/channelContractB",
		},
		{
			Title: "未回款金额",
			Tag:   "notRemitMoney",
			Pmsm:  bean.PMSM_REMIT_DOWN_ACCOUNT,
			Href:  "/home/channelContract/channelContractC",
		},
		{
			Title: "合同待签数",
			Tag:   "notContractCount",
			Pmsm:  bean.PMSM_CONTRACT_CHANNEL,
			Href:  "/home/channelContract/channelContractA",
		},
	}
	var total int64
	for i, d := range data {
		_, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, d.Pmsm, nil)
		if err == nil {
			data[i].DashboardBasicInfo, err = appcache.GetDashboardBasicInfoCache(d.Tag)
			data[i].Show = true
			total++
		}
	}
	c.RespJSONDataWithTotal(data, total)
}
