package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"time"
)

// DistributionCompanyController operations for DistributionCompany
type DistributionCompanyController struct {
	BaseController
}

// URLMapping ...
func (c *DistributionCompanyController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create DistributionCompany
// @Param	body		body 	models.DistributionCompany	true		"body for DistributionCompany content"
// @Success 201 {int} models.DistributionCompany
// @Failure 403 body is empty
// @router / [post]
func (c *DistributionCompanyController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_DISTRIBUTION_COMPANY, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	var v models.DistributionCompany
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	err1 := c.Validate(&v)
	if err1 != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	v.CreateTime = time.Now().Unix()
	v.UpdateUserID = c.Uid()
	v.UpdateTime = time.Now().Unix()

	if v.YunduanResponsiblePerson == 0 && v.YouliangResponsiblePerson == 0 {
		c.RespJSON(bean.CODE_Params_Err, "商务负责人至少填写一个!")
		return
	}

	if _, err2 := models.AddDistributionCompany(&v); err2 != nil {
		c.RespJSON(bean.CODE_Bad_Request, err2.Error())
		return
	}
	c.RespJSONData(v)
}

// GetOne ...
// @Title Get One
// @Description get DistributionCompany by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DistributionCompany
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DistributionCompanyController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetDistributionCompanyById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	models.GetDistributionCompanyAdditionalInfo(v)

	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get DistributionCompany
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.DistributionCompany
// @Failure 403
// @router / [get]
func (c *DistributionCompanyController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

func (c *DistributionCompanyController) getAll() (total int64, m []*models.DistributionCompany, errCode int, err error) {
	_, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_DISTRIBUTION_COMPANY, nil)
	if err != nil {
		errCode = bean.CODE_Forbidden
		return
	}
	//filter, err := tool.BuildFilter(c.Controller, 20)
	//if err != nil {
	//	errCode = bean.CODE_Params_Err
	//	return
	//}
	//filter.Order = append(filter.Order, "desc")
	//filter.Sortby = append(filter.Sortby, "update_time")
	//
	//tool.InjectPermissionWhere(where, &filter.Where)
	//m = []*models.DistributionCompany{}
	//total, err = tool.GetAllByFilterWithTotal(new(models.DistributionCompany), &m, filter)

	offset, _ := c.GetInt("offset", 0)
	limit, _ := c.GetInt("limit", 15)
	company_id := c.GetString("company_id", "")
	person, _ := c.GetInt("people", 0)

	m, total, err = models.GetAllDistributionCompaniesWithTotal(offset, limit, company_id, person)
	if err != nil {
		errCode = bean.CODE_Bad_Request
		return
	}

	// add the company info
	for _, v := range m {
		models.GetDistributionCompanyAdditionalInfo(v)
	}

	return
}

// GetAllCompanyName ...
// @Title Get All CompanyName
// @Description get ChannelCompany
// @Success 200 {object} models.ChannelCompany
// @Failure 403
// @router /companyname/ [get]
func (c *DistributionCompanyController) GetAllCompanyName() {
	filter, err := tool.BuildFilter(c.Controller, 0)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	filter.Fields = append(filter.Fields, "Id", "CompanyId")

	var ss []models.DistributionCompany
	_, err = tool.GetAllByFilterWithTotal(new(models.DistributionCompany), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	if len(ss) == 0 {
		return
	}
	models.AddDistributionCompanyInfo(&ss)
	models.AddGameIdInfo(&ss)

	c.RespJSONData(ss)
}

// GetContractState ...
// @Title Get All CompanyName
// @Description get ChannelCompany
// @Success 200 {object} models.ChannelCompany
// @Failure 403
// @router /state/ [get]
func (c *DistributionCompanyController) GetContractState() {
	filter, err := tool.BuildFilter(c.Controller, 0)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	// 仅选择主合同状态
	ids := []interface{}{11}
	filter.Where["type__in"] = ids

	var ss []models.Types
	_, err = tool.GetAllByFilterWithTotal(new(models.Types), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	if len(ss) == 0 {
		return
	}
	models.AddCompanyIdInfo(&ss, 1)

	c.RespJSONData(ss)
}

// Put ...
// @Title Put
// @Description update the DistributionCompany
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.DistributionCompany	true		"body for DistributionCompany content"
// @Success 200 {object} models.DistributionCompany
// @Failure 403 :id is not int
// @router /:id [put]
func (c *DistributionCompanyController) Put() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_DISTRIBUTION_COMPANY, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DistributionCompany{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	//编辑发行商 添加更新人、时间
	v.UpdateUserID = c.Uid()
	v.UpdateTime = time.Now().Unix()

	if v.YunduanResponsiblePerson == 0 && v.YouliangResponsiblePerson == 0 {
		c.RespJSON(bean.CODE_Params_Err, "商务负责人至少填写一个!")
		return
	}

	if err := models.UpdateDistributionCompanyById(&v, where); err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(v)
}

// Delete ...
// @Title Delete
// @Description delete the DistributionCompany
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DistributionCompanyController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteDistributionCompany(id); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData("OK")
}

// @Title 获取全部地区
// @Description 获取全部地区
// @Success 200 {string} !
// @router /region [get]
func (c *DistributionCompanyController) GetRegion() {
	result, err := models.GetDistributionCompanyRegion()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(result)
}
