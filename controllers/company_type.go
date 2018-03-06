package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
)

// CompanyTypeController operations for CompanyType
type CompanyTypeController struct {
	BaseController
}

// URLMapping ...
func (c *CompanyTypeController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create CompanyType
// @Param	body		body 	models.CompanyType	true		"body for CompanyType content"
// @Success 201 {int} models.CompanyType
// @Failure 403 body is empty
// @router / [post]
func (c *CompanyTypeController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_COMPANY_TYPE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var v models.CompanyType
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if v.Name == "" {
		c.RespJSON(bean.CODE_Params_Err, "The Company Name field is required.")
		return
	}

	if _, err := models.AddCompanyType(&v, c.Uid()); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(v)
}

// GetOne ...
// @Title Get One
// @Description get CompanyType by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.CompanyType
// @Failure 403 :id is empty
// @router /:id [get]
func (c *CompanyTypeController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetCompanyTypeById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get CompanyType
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.CompanyType
// @Failure 403
// @router / [get]
func (c *CompanyTypeController) GetAll() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_COMPANY_TYPE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	tool.InjectPermissionWhere(where, &filter.Where)
	filter.Sortby = []string{"id"}
	filter.Order = []string{"desc"}

	var m []*models.CompanyType
	total, err := tool.GetAllByFilterWithTotal(new(models.CompanyType), &m, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONDataWithTotal(m, total)
	return
}

// Put ...
// @Title Put
// @Description update the CompanyType
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.CompanyType	true		"body for CompanyType content"
// @Success 200 {object} models.CompanyType
// @Failure 403 :id is not int
// @router /:id [put]
func (c *CompanyTypeController) Put() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_CHANNEL_COMPANY, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.CompanyType{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if v.Name == "" {
		c.RespJSON(bean.CODE_Params_Err, "The Company Name field is required.")
		return
	}

	if err := models.UpdateCompanyTypeById(&v, where, c.Uid()); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData("ok")
}

// Delete ...
// @Title Delete
// @Description delete the CompanyType
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *CompanyTypeController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteCompanyType(id, c.Uid()); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("ok")
}

// @Description 快递管理页面 下拉框 选择收件公司
// @router /list/ [get]
func (c *CompanyTypeController) GetCompanyList(){
	companies, err := models.GetCompanyList()
	if err != nil{
		c.RespJSON(bean.CODE_Params_Err,err.Error())
		return
	}
	c.RespJSONData(companies)
}