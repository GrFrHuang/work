package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"time"
)

// DevelopCompanyController operations for DevelopCompany
type DevelopCompanyController struct {
	BaseController
}

// URLMapping ...
func (c *DevelopCompanyController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create DevelopCompany
// @Param	body		body 	models.DevelopCompany	true		"body for DevelopCompany content"
// @Success 201 {int} models.DevelopCompany
// @Failure 403 body is empty
// @router / [post]
func (c *DevelopCompanyController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_DEVELOP_COMPANY, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	var v models.DevelopCompany
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	err1 := c.Validate(&v)
	if err1 != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	v.UpdateUserID = c.Uid()
	v.UpdateTime =time.Now().Unix()
	if _, err2 := models.AddDevelopCompany(&v); err2 != nil {
		c.RespJSON(bean.CODE_Bad_Request, err2.Error())
		return
	}
	c.RespJSONData(v)
}

// GetOne ...
// @Title Get One
// @Description get DevelopCompany by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DevelopCompany
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DevelopCompanyController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetDevelopCompanyById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get DevelopCompany
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.DevelopCompany
// @Failure 403
// @router / [get]
func (c *DevelopCompanyController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

func (c * DevelopCompanyController) getAll() (total int64, m []*models.DevelopCompany, errCode int, err error)  {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_DEVELOP_COMPANY, nil)
	if err != nil {
		errCode = bean.CODE_Forbidden
		return
	}
	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		errCode = bean.CODE_Params_Err
		return
	}
	filter.Order = append(filter.Order, "desc")
	filter.Sortby = append(filter.Sortby, "update_time")
	tool.InjectPermissionWhere(where, &filter.Where)
	m = []*models.DevelopCompany{}
	total, err = tool.GetAllByFilterWithTotal(new(models.DevelopCompany), &m, filter)
	if err != nil {
		errCode = bean.CODE_Bad_Request
		return
	}
	// add the additional info
	for _, v := range m {
		models.GetDevelopCompanyAdditionalInfo(v)
	}
	return
}

// GetAllCompanyName ...
// @Title Get All CompanyName
// @Description get ChannelCompany
// @Success 200 {object} models.ChannelCompany
// @Failure 403
// @router /companyname/ [get]
func (c *DevelopCompanyController) GetAllCompanyName() {
	filter, err := tool.BuildFilter(c.Controller, 0)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	filter.Fields = append(filter.Fields, "Id", "CompanyId")

	ss := []models.DevelopCompany{}
	_, err = tool.GetAllByFilterWithTotal(new(models.DevelopCompany), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	if len(ss) == 0 {
		return
	}
	models.AddDevelopCompanyInfo(&ss)

	c.RespJSONData(ss)

}

// Put ...
// @Title Put
// @Description update the DevelopCompany
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.DevelopCompany	true		"body for DevelopCompany content"
// @Success 200 {object} models.DevelopCompany
// @Failure 403 :id is not int
// @router /:id [put]
func (c *DevelopCompanyController) Put() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_DEVELOP_COMPANY, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DevelopCompany{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	v.UpdateUserID=c.Uid()
	v.UpdateTime=time.Now().Unix()
	if err := models.UpdateDevelopCompanyById(&v, where); err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(v)
}

// Delete ...
// @Title Delete
// @Description delete the DevelopCompany
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DevelopCompanyController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteDevelopCompany(id); err == nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData("OK")
}

