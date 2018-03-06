package workflow

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/controllers"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/models/workflow"
	"strconv"
)

//发包工作流信息
type WorkflowSendPackage struct {
	controllers.BaseController
}

// 通过TaskId获取工作流信息
// @Title Get One
// @Description get WorkflowSendPackage by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} workflow.WorkflowSendPackage
// @Failure 403 :id is empty
// @router /task/:id [get]
func (c *WorkflowSendPackage) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := workflow.GetWorkflowSendPackageByTaskId(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Acceptable, err.Error())
		return
	}
	c.RespJSONData(v)
	return
}

// Post ...
// @Title Post
// @Description create WorkflowTask
// @Param	body		body 	workflow.WorkflowSendPackage	true		"body for WorkflowTask content"
// @Success 201 {int} workflow.WorkflowSendPackage
// @Failure 403 body is empty
// @router / [put]
func (c *WorkflowSendPackage) Put() {
	var v workflow.WorkflowSendPackage
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	if v.Id == 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}

	if v.Status == 2 {

		err := workflow.ChangeUpdateMan(v)
		if err != nil {
			c.RespJSON(bean.CODE_Not_Acceptable, err.Error())
			//return
		}
	}
	err := workflow.ChangeWorkflowInfoById(&v)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Acceptable, err.Error())
		return
	}
	c.RespJSONSuccess()
	return

}

// get ...
// @Title get
// @Description create WorkflowTask
// @Param	body		body 	workflow.WorkflowSendPackage	true		"body for WorkflowTask content"
// @Success 201 {int} workflow.WorkflowSendPackage
// @Failure 403 body is empty
// @router /is_handle [get]
func (c *WorkflowSendPackage) IsHandle() {
	package_id, err := c.GetInt("package_id")
	user := c.Uid()
	//user, err := models.GetUserInfoByToken(token)
	if err != nil {
		c.RespJSONData(err.Error())
	}
	key, err := workflow.IsHandle(package_id, user)
	if err != nil {
		c.RespJSONData(err.Error())
	}
	c.RespJSONData(key)
}
