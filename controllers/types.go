package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
)

// 类型集合
type TypesController struct {
	BaseController
}

// URLMapping ...
func (c *TypesController) URLMapping() {
	c.Mapping("Post", c.Post)
	//c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	//c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Types
// @Param	body		body 	models.Types	true		"body for Types content"
// @Success 201 {int} models.Types
// @Failure 403 body is empty
// @router / [post]
func (c *TypesController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_TYPES, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	var v models.Types
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if data, err := models.AddTypes(&v); err == nil {
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

// GetAll ...
// @Title Get All
// @Description get Types
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Types
// @Failure 403
// @router / [get]
func (c *TypesController) GetAll() {
	//where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_TYPES, nil)
	//if err != nil {
	//	c.RespJSON(bean.CODE_Forbidden, err.Error())
	//	return
	//}
	filter, err := tool.BuildFilter(c.Controller, 0)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	filter.Where["is_deleted"] = []interface{}{"0"}
	ss := []models.Types{}
	err = tool.GetAllByFilter(new(models.Types), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(ss)
}

// Delete ...
// @Title Delete
// @Description delete the Types
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TypesController) Delete() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_DELETE, bean.PMSM_TYPES, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteTypes(id, where); err == nil {
		c.RespJSONData("delete types, id=" + idStr)
		return
	} else {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
}

// @router /get_courier [get]
func (c *TypesController) GetAllCourier() {

	types := models.GetAllCourier()
	c.RespJSONData(types)

}
