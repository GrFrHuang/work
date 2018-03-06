package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
)

// WarningTypeController operations for WarningType
type WarningTypeController struct {
	BaseController
}

// URLMapping ...
func (c *WarningTypeController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create WarningType
// @Param	body		body 	models.WarningType	true		"body for WarningType content"
// @Success 201 {int} models.WarningType
// @Failure 403 body is empty
// @router / [post]
func (c *WarningTypeController) Post() {
	var v models.WarningType
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddWarningType(&v); err == nil {
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
// @Description get WarningType by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.WarningType
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WarningTypeController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetWarningTypeById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get WarningType
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.WarningType
// @Failure 403
// @router / [get]
func (c *WarningTypeController) GetAll() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_WARNING_TYPE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	//endTime, _ := c.GetInt("end_time")
	//channels := c.GetString("channels")

	tool.InjectPermissionWhere(where, &filter.Where)

	filter.Sortby = []string{"WarningTypeId"}
	filter.Order = []string{"asc"}
	ss := []models.WarningType{}
	total, err := tool.GetAllByFilterWithTotal(new(models.WarningType), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONDataWithTotal(ss, total)
}

// Put ...
// @Title Put
// @Description update the WarningType
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.WarningType	true		"body for WarningType content"
// @Success 200 {object} models.WarningType
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WarningTypeController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.WarningType{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateWarningTypeById(&v); err == nil {
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
// @Description delete the WarningType
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WarningTypeController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteWarningType(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
