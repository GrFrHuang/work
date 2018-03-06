package workflow

import (
	"kuaifa.com/kuaifa/work-together/controllers"
	"kuaifa.com/kuaifa/work-together/models/workflow"
	"kuaifa.com/kuaifa/work-together/tool"
	"kuaifa.com/kuaifa/work-together/models/bean"
)
//工作流
type WorkflowNameControllers struct {
	controllers.BaseController
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
// @Success 200 {object} models.WorkflowName
// @Failure 403
// @router / [get]
func (c *WorkflowNameControllers) GetAll() {
	var ss []workflow.WorkflowName
	total, err := tool.GetAllWithTotal(c.Controller, new(workflow.WorkflowName), &ss, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}
