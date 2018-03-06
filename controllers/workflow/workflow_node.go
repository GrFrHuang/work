package workflow

import (
	"kuaifa.com/kuaifa/work-together/controllers"
	"kuaifa.com/kuaifa/work-together/models/workflow"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

type WorkflowNodeControllers struct {
	controllers.BaseController
}

// @Title Get
// @Description create WorkflowTask
// @Param	WfNameId	query	int	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	status	query	int	false	"1 获取需要能撤销的人,2获取本部门的人"
// @Param	step	query	int	false	"步骤"
// @Success 201 ok
// @Failure 403 body is empty
// @router /get_next_users/ [get]
func (c *WorkflowNodeControllers) GetUsersByNextNodeId() {
	WfNameId, _ := c.GetInt("WfNameId", 0)
	status, _ := c.GetInt("status", 0)
	step, _ := c.GetInt("step", 0)
	userId := c.Uid()
	if userId == 0 || WfNameId == 0 || step == 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}
	info, err := workflow.GetNextUsersByUid(userId, status, step, WfNameId)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(info)
	return
}

// @Title Get
// @Description create WorkflowTask
// @Param	step	query	int	false	"步骤"
// @Success 201 ok
// @Failure 403 body is empty
// @router /workflow_node_id/ [get]
func (c *WorkflowNodeControllers) GetWorkflowNodeByUid() {
	userId := c.Uid()
	step, _ := c.GetInt("step", 0)
	if userId == 0 {
		c.RespJSON(bean.CODE_Params_Err, "未获取到UID")
		return
	}
	if step == 0 {
		c.RespJSON(bean.CODE_Params_Err, "未获取到步骤信息")
		return
	}
	info, err := workflow.GetWorkflowNodeByUid(userId, step)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(info)
	return
}
