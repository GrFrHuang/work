package controllers

import (
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

type UserDesktopScreenLogControllers struct {
	BaseController
}

// GetWaring ...
// @Title GetWaring
// @Description get GetWaring
// @Success 200 {object} workflow.WorkflowTask
// @Failure 403
// @router / [get]
func (c *UserDesktopScreenLogControllers) GetWaring() {
	err, info := models.GetWarnings()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, "未查询到数据")
		return
	}
	c.RespJSONData(info)
}
