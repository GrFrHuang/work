package controllers

import (
	"encoding/json"
	"github.com/bysir-zl/bjson"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"strconv"
)

// 权限编辑
type PermissionController struct {
	BaseController
}

// URLMapping ...
func (c *PermissionController) URLMapping() {

}

// Post ...
// @Title 新建
// @Description 新建一个权限
// @Param	body		body 	models.Permission	true		"body for Permission content"
// @Success 200 {object} models.Permission
// @Failure 403 body is empty
// @router / [post]
func (c *PermissionController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	var v models.Permission
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, bean.CodeString(bean.CODE_Params_Err))
		return
	}
	if v.Name == "" {
		c.RespJSON(bean.CODE_Params_Err, "Name is not set")
		return
	}
	switch v.Type {
	case models.Type_can, models.Type_notcan:
		v.Field = ""
		v.Condition = ""
	case models.Type_condition:
		if v.Field == "" || v.Condition == "" {
			c.RespJSON(bean.CODE_Params_Err, "field or condition are not set")
			return
		}
	default:
		c.RespJSON(bean.CODE_Params_Err, "type must in [2,3,4]")
		return
	}
	bj, _ := bjson.New([]byte(v.Methods))
	if bj.Len() == 0 {
		c.RespJSON(bean.CODE_Params_Err, "methods must is a json string")
		return
	}

	if _, err := models.AddPermission(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err)
		return
	}
	c.RespJSONData(v)
	return
}

// GetOne ...
// @Title Get One
// @Description get Permission by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Permission
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PermissionController) GetOne() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetPermissionById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Permission
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Permission
// @Failure 403
// @router / [get]
func (c *PermissionController) GetAll() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	us := []models.Permission{}
	err = tool.GetAll(c.Controller, new(models.Permission), &us, 1000)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(us)
}

// @Title 分组获取权限
// @Description get Permission
// @Success 200 {object} models.Permission
// @Failure 403
// @router /group [get]
func (c *PermissionController) GetAllWithGroup() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	us := []models.Permission{}
	err = tool.GetAll(c.Controller, new(models.Permission), &us, 1000)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	result := map[string][]models.Permission{}

	for _, p := range us {
		_, group := bean.PermissionModelString(p.Model)
		if _, ok := result[group]; !ok {
			result[group] = []models.Permission{}
		}
		result[group] = append(result[group], p)
	}
	type x struct {
		Group string    `json:"group,omitempty"`
		Ps    []models.Permission   `json:"ps,omitempty"`
	}
	orderGroup := []string{
		"超管", "游戏", "游戏流水", "CP合同", "渠道合同", "CP结算", "渠道结算", "CP对账", "渠道对账",
		"用户", "部门", "发行商", "研发商", "渠道商", "日志", "快递管理", "公司列表",
	}
	r := []x{}
	for _, key := range orderGroup {
		r = append(r, x{
			Group: key,
			Ps:    result[key],
		})
	}

	c.RespJSONData(r)
}

// Put ...
// @Title Put
// @Description 更改一个权限
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Permission	true		"body for Permission content"
// @Success 200 {object} models.Permission
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PermissionController) Put() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Permission{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdatePermissionById(&v); err == nil {
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
// @Description delete the Permission
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *PermissionController) Delete() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_DELETE, bean.PMSM_USER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeletePermission(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
