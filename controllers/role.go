package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"strconv"
)

// 角色管理
type RoleController struct {
	BaseController
}

// URLMapping ...
func (c *RoleController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description 创建一个角色
// @Param	body		body 	models.Role	true		"body for Role content"
// @Success 201 {int} models.Role
// @Failure 403 body is empty
// @router / [post]
func (c *RoleController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	var v models.Role
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	ids := []int{}
	if err := json.Unmarshal([]byte(v.PermissionIds), &ids); err != nil {
		c.RespJSON(bean.CODE_Params_Err, "PermissionIds is not json string")
		return
	}

	if _, err := models.AddRole(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(v)
}

// GetOne ...
// @Title Get One
// @Description get Role by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Role
// @Failure 403 :id is empty
// @router /:id [get]
func (c *RoleController) GetOne() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetRoleById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get Role
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Role
// @Failure 403
// @router / [get]
func (c *RoleController) GetAll() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	rs := []models.Role{}
	err = tool.GetAll(c.Controller, new(models.Role), &rs, 1000)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	// 添加权限信息
	models.GroupAddPermissionInfo(rs)
	c.RespJSONData(rs)
}

// Put ...
// @Title Put
// @Description 更新一个角色
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Role	true		"body for Role content"
// @Success 200 {object} models.Role
// @Failure 403 :id is not int
// @router /:id [put]
func (c *RoleController) Put() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Role{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	ids := []int{}
	if err := json.Unmarshal([]byte(v.PermissionIds), &ids); err != nil {
		c.RespJSON(bean.CODE_Params_Err, "PermissionIds is not json string")
		return
	}
	if err := models.UpdateRoleById(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(v)
}

// Delete ...
// @Title Delete
// @Description 删除一个角色
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *RoleController) Delete() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_DELETE, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteRole(id); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData("OK")
}
