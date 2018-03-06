package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/tool"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

// 报警记录
type AlarmLogController struct {
	BaseController
}

// URLMapping ...
func (c *AlarmLogController) URLMapping() {

}

// Post ...
// @Title Post
// @Description create AlarmLog
// @Param	body		body 	models.AlarmLog	true		"body for AlarmLog content"
// @Success 201 {int} models.AlarmLog
// @Failure 403 body is empty
// @router / [post]
func (c *AlarmLogController) Post() {
	var v models.AlarmLog
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddAlarmLog(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get AlarmLog by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.AlarmLog
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AlarmLogController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetAlarmLogById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get AlarmLog
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.AlarmLog
// @Failure 403
// @router / [get]
func (c *AlarmLogController) GetAll() {
	as := []models.AlarmLog{}

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	filter.Where = map[string][]interface{}{
		"is_hide": {2},
	}
	filter.Order = []string{"desc"}
	filter.Sortby = []string{"id"}

	total, e := tool.GetAllByFilterWithTotal(new(models.AlarmLog), &as, filter)
	if e != nil {
		c.RespJSON(bean.CODE_Bad_Request, e.Error())
		return
	}

	c.RespJSONDataWithTotal(as, total)
}


// Delete ...
// @Title Delete
// @Description delete the AlarmLog
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *AlarmLogController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteAlarmLog(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
