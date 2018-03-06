package controllers

import (
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/tool"
)

type UserDesktopClientController struct {
	BaseController
}

// GetAll ...
// @Title Get All
// @Description get WorkflowName
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} workflow.WorkflowTask
// @Failure 403
// @router / [get]
func (c *UserDesktopClientController) GetAll() {
	var ss []models.UserDesktopClient
	total, err := tool.GetAllWithTotal(c.Controller, new(models.UserDesktopClient), &ss, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, "未查询到数据")
		return
	}
	for i := 0; i < len(ss); i++ {
		err, name := models.GetNickNameById(ss[i].Uid)
		if err == nil {
			ss[i].NickName = name
		}
	}
	c.RespJSONDataWithTotal(ss, total)
}
