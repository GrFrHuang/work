package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

// 报警规则管理
type AlarmRuleController struct {
	BaseController
}

// URLMapping ...
func (c *AlarmRuleController) URLMapping() {
}

// Post ...
// @Title Post
// @Description create AlarmRule
// @Param	body		body 	models.AlarmRule	true		"body for AlarmRule content"
// @Success 201 {int} models.AlarmRule
// @Failure 403 body is empty
// @router / [post]
func (c *AlarmRuleController) Post() {
	var v models.AlarmRule
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if _, err := models.AddAlarmRule(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(v)
}

// GetOne ...
// @Title Get One
// @Description get AlarmRule by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.AlarmRule
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AlarmRuleController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetAlarmRuleById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get AlarmRule
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.AlarmRule
// @Failure 403
// @router / [get]
func (c *AlarmRuleController) GetAll() {
	r, err := models.GetAllAlarmRuleWithType()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(r)
}

// Put ...
// @Title Put
// @Description 更新规则
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.AlarmRule	true		"body for AlarmRule content"
// @Success 200 {object} models.AlarmRule
// @Failure 403 :id is not int
// @router /:id [put]
func (c *AlarmRuleController) Put() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_ALARM_RULE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.AlarmRule{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	//fmt.Printf("v:%v\n", v)
	//return
	if err := models.UpdateAlarmRuleById(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("保存成功")
}

// Delete ...
// @Title Delete
// @Description delete the AlarmRule
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *AlarmRuleController) Delete() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_DELETE, bean.PMSM_ALARM_RULE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteAlarmRule(id); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("OK")
}
