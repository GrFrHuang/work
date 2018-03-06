package workflow

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models/workflow"
	"kuaifa.com/kuaifa/work-together/controllers"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"fmt"
	"github.com/astaxie/beego"
	"net/http"
	"io/ioutil"
	"github.com/astaxie/beego/orm"
)

// 新建工作流程
type WorkflowTaskController struct {
	controllers.BaseController
}

// Post ...
// @Title Post
// @Description create WorkflowTask
// @Param	body		body 	workflow.WorkflowTask	true		"body for WorkflowTask content"
// @Success 201 {int} workflow.WorkflowTask
// @Failure 403 body is empty
// @router / [post]
func (c *WorkflowTaskController) Post() {
	var v workflow.WorkflowTask
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	if v.TaskName == "" || v.GameId == 0 || v.ChannelCode == "" || v.WfNameId == 0 || v.DepartmentId == 0 || c.Uid() == 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
		return
	}
	v.UserId = c.Uid()
	err, id := workflow.AddWorkflowTask(&v)

	if err != nil {
		c.RespJSON(bean.CODE_Not_Acceptable, err.Error())
		return
	}
	c.RespJSONData(id)
	return

}

// GetOne ...
// @Title Get One
// @Description get WorkflowTask by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} workflow.WorkflowTask
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WorkflowTaskController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := workflow.GetWorkflowTaskById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Acceptable, err.Error())
		return
	}
	c.RespJSONData(v)
	return
}

// GetOne ...
// @Title Get One
// @Description get WorkflowTask by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} workflow.WorkflowTask
// @Failure 403 :id is empty
// @router /task/:id [get]
func (c *WorkflowTaskController) GetChannelInfoIdByTaskId() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := workflow.GetChannelInfoIdByTaskId(id)
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
// @Param	body		body 	workflow.ChangeProgressInput	true		"body for WorkflowTask content"
// @Success 201 {int} workflow.ChangeProgressInput
// @Failure 403 body is empty
// @router / [put]
func (c *WorkflowTaskController) ChangeProgress() {
	var v workflow.ChangeProgressInput
	fmt.Println(string(c.Ctx.Input.RequestBody))
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	v.UserId = c.Uid()
	if err := workflow.ChangeProgress(&v); err != nil {
		c.RespJSON(bean.CODE_Not_Acceptable, err.Error())
		return
	}

	//token := c.Ctx.Input.Header("Authorization")
	cps := beego.AppConfig.String("cps_url")

	for _, vule := range v.User {
		task_id := strconv.Itoa(v.TaskId)
		task, err := workflow.GetWorkflowTaskById(v.TaskId)
		if v.TaskId == 0 {
			c.RespJSONData("日志记录失败")
		}
		uri := cps + "/v1/message_scrollbar/setmessaget?email=" + vule.Email + "&uname=" + vule.Name + "&task_id=" + task_id + "&task_name=" + task.TaskName + "&game_name=" + task.GameName
		fmt.Println(uri)
		req, err := http.Get(uri)
		defer req.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))

	}

	c.RespJSONSuccess()
	return
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
// @Param	department_id	query	int	false	"Start position of result set. Must be an integer"
// @Param	wf_name_id	query	int	false	"Start position of result set. Must be an integer"
// @Success 200 {object} workflow.WorkflowTask
// @Failure 403
// @router / [get]
func (c *WorkflowTaskController) GetAll() {
	DepartmentId, _ := c.GetInt("department_id", 0)
	WfNameId, _ := c.GetInt("wf_name_id", 0)
	if DepartmentId == 0 || WfNameId == 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数不正确")
	}
	o := orm.NewOrm()
	var node workflow.WorkflowNode
	node.DepartmentId = DepartmentId
	node.WfNameId = WfNameId
	if err := o.Read(&node, "department_id", "wf_name_id"); err != nil {
		c.RespJSONDataWithTotal("", 0)
		return
	}
	var ss []workflow.WorkflowTask
	total, err := tool.GetAllWithTotal(c.Controller, new(workflow.WorkflowTask), &ss, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, "未查询到数据")
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

// GetAll ...
// @Title Get All
// @Description get WorkflowName
// @Success 200 {object} workflow.WorkflowTask
// @Failure 403
// @router /condition [get]
func (c *WorkflowTaskController) GetAllByCondition() {
	data, err := workflow.GetAllByCondition([]string{"task_name", "channel_name", "game_name"})
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(data)
	return
}

// GetAll ...
// @Title Get All
// @Description get WorkflowName
// @Param	id	query	int	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	status	query	int	false	"Fields returned. e.g. col1,col2 ..."
// @Param	remarks	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	channel_id	query	int	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	contract_id	query	int	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Success 200 {object} workflow.WorkflowTask
// @Failure 403
// @router /status [get]
func (c *WorkflowTaskController) ChangeStatucbyId() {
	id, _ := c.GetInt("id", 0)
	status, _ := c.GetInt("status", 0)
	channel_id, _ := c.GetInt("channel_id", 0)
	contract_id, _ := c.GetInt("contract_id", 0)
	remarks := c.GetString("remarks")
	if id == 0 || status == 0 || channel_id == 0 || contract_id == 0 || remarks == "" {
		c.RespJSON(bean.CODE_Bad_Request, "参数不正确")
		return
	}
	if err := workflow.ChangeStatucbyId(id, status, channel_id, contract_id, remarks); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONSuccess()
}
