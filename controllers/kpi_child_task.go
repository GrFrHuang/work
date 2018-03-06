package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"github.com/astaxie/beego"
)

// KpiChildTaskController operations for KpiChildTask
type KpiChildTaskController struct {
	beego.Controller
}

// URLMapping ...
func (c *KpiChildTaskController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// GetOne ...
// @Title Get One
// @Description get KpiChildTask by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.KpiChildTask
// @Failure 403 :id is empty
// @router /:id [get]
func (c *KpiChildTaskController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetKpiChildTaskById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the KpiChildTask
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.KpiChildTask	true		"body for KpiChildTask content"
// @Success 200 {object} models.KpiChildTask
// @Failure 403 :id is not int
// @router /:id [put]
func (c *KpiChildTaskController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.KpiChildTask{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateKpiChildTaskById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the KpiChildTask
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *KpiChildTaskController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteKpiChildTask(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
