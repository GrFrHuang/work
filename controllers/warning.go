package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"fmt"
	"kuaifa.com/kuaifa/work-together/tool"
	"strings"
)

// WarningController operations for Warning
type WarningController struct {
	BaseController
}

// URLMapping ...
func (c *WarningController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Warning
// @Param	body		body 	models.Warning	true		"body for Warning content"
// @Success 201 {int} models.Warning
// @Failure 403 body is empty
// @router / [post]
func (c *WarningController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_WARNING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	var v models.Warning
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddWarning(&v); err == nil {
			c.RespJSONData("OK")
		} else {
			fmt.Println(err.Error())
			c.RespJSONData(err.Error())
		}
	} else {
		fmt.Println(err.Error())
		c.RespJSONData(err.Error())
	}
}

// GetOne ...
// @Title Get One
// @Description get Warning by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Warning
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WarningController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetWarningById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get static warning types
// @Success 200 {object} map[int]string
// @router /types [get]
func (c *WarningController) GetStaticWarningTypes() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_WARNING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Method_Not_Allowed, "没有预警添加权限")
		return
	}
	v := models.GetWarningStaticType()
	fmt.Println(v)
	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get Warning
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Warning
// @Failure 403
// @router / [get]
func (c *WarningController) GetAll() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_WARNING, nil)
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

	filter.Sortby = []string{"CreateTime"}
	filter.Order = []string{"desc"}
	var ss []models.Warning
	total, err := tool.GetAllByFilterWithTotal(new(models.Warning), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONDataWithTotal(ss, total)
}

// GetWarningName ...
// @Title Get All
// @Description get Warning
// @Success 200 {object} models.Warning
// @Failure 403
// @router /name/ [get]
func (c *WarningController) GetWarningName() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_WARNING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	result := models.GetAllWarningName()
	c.RespJSONData(result)
}

// GetWarningDetail ...
// @Title Get All
// @Description get Warning
// @Param	ids		path 	string	true		"The id you want to get"
// @Success 200 {object} models.Warning
// @Failure 403
// @router /detail [get]
func (c *WarningController) GetWarningDetail() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_WARNING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	ids := c.GetString("ids", "")
	if ids == "" {
		c.RespJSON(bean.CODE_Bad_Request, "请求参数错误")
		return
	}

	typeIds := strings.Split(ids, ",")
	if len(typeIds) == 0 {
		c.RespJSON(bean.CODE_Bad_Request, "请求参数错误")
		return
	}

	result := models.GetWarningDetail(typeIds)
	c.RespJSONData(result)
}

// BatchPut ...
// @Title Put
// @Description update the Warning
// @Param	body		body 	models.Warning	true		"body for Warning content"
// @Success 200 {object} models.Warning
// @Failure 403 :id is not int
// @router /batch/ [put]
func (c *WarningController) BatchPut() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_WARNING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	var values []models.Warning
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &values); err == nil {
		for _, v := range values {
			old, _ := models.GetWarningById(v.Id)
			old.UserIds = v.UserIds
			if err := models.UpdateWarningById(old); err == nil {
				c.RespJSONData("OK")
			} else {
				c.RespJSON(bean.CODE_Method_Not_Allowed, err.Error())
				return
			}
		}
	} else {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
}

// Put ...
// @Title Put
// @Description update the Warning
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Warning	true		"body for Warning content"
// @Success 200 {object} models.Warning
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WarningController) Put() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_WARNING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Warning{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateWarningById(&v); err == nil {
			c.RespJSONData("OK")
		} else {
			c.RespJSON(bean.CODE_Method_Not_Allowed, err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
}

// Delete ...
// @Title Delete
// @Description delete the Warning
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WarningController) Delete() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_DELETE, bean.PMSM_WARNING, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteWarning(id); err == nil {
		c.RespJSONData("OK")
	} else {
		c.RespJSONData(err.Error())
	}
}
