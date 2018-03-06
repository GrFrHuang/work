package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
)

// 工会政策
type SociatyPolicyController struct {
	BaseController
}

// URLMapping ...
func (c *SociatyPolicyController) URLMapping() {
	c.Mapping("Post", c.Post)
	//c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	//c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create SociatyPolicy
// @Param	body		body 	models.SociatyPolicy	true		"body for SociatyPolicy content"
// @Success 201 {int} models.SociatyPolicy
// @Failure 403 body is empty
// @router / [post]
func (c *SociatyPolicyController) Post() {
	var v models.SociatyPolicy
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if data, err := models.AddSociatyPolicy(&v); err == nil {
			c.RespJSONData(data)
			return
		} else {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
}

// GetOne ...
// @Title Get One
// @Description get SociatyPolicy by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.SociatyPolicy
// @Failure 403 :id is empty
// @router /:id [get]
func (c *SociatyPolicyController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetSociatyPolicyById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get SociatyPolicy
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.SociatyPolicy
// @Failure 403
// @router / [get]
func (c *SociatyPolicyController) GetAll() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_SOCIATY_POLICY, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	filter.Where["deleted"] = []interface{}{"0"}
	tool.InjectPermissionWhere(where, &filter.Where)
	ss := []models.SociatyPolicy{}
	err = tool.GetAllByFilter(new(models.SociatyPolicy), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(ss)
}

// Put ...
// @Title Put
// @Description update the SociatyPolicy
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.SociatyPolicy	true		"body for SociatyPolicy content"
// @Success 200 {object} models.SociatyPolicy
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SociatyPolicyController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.SociatyPolicy{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateSociatyPolicyById(&v); err == nil {
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
// @Description delete the SociatyPolicy
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SociatyPolicyController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteSociatyPolicy(id); err == nil {
		c.RespJSONData("delete sociaty_policy, id="+idStr)
		return
	} else {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
}
