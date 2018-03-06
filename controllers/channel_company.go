package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"strconv"
	"time"
)

// 渠道商
type ChannelCompanyController struct {
	BaseController
}

// URLMapping ...
func (c *ChannelCompanyController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create ChannelCompany
// @Param	body		body 	models.ChannelCompany	true		"body for ChannelCompany content"
// @Success 201 {int} models.ChannelCompany
// @Failure 403 body is empty
// @router / [post]
func (c *ChannelCompanyController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_CHANNEL_COMPANY, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	var v models.ChannelCompany
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if v.ChannelCode == "" {
		c.RespJSON(bean.CODE_Params_Err, "The Channel code field is required.")
		return
	}

	if v.CompanyId == 0 {
		c.RespJSON(bean.CODE_Params_Err, "The Company id field is required.")
		return
	}

	if v.YunduanResponsiblePerson == 0 && v.YouliangResponsiblePerson == 0 {
		c.RespJSON(bean.CODE_Params_Err, "Must chose one Person Responsible at least.")
		return
	}

	if v.CooperateState == 0 {
		c.RespJSON(bean.CODE_Params_Err, "The CooperateState field is required!")
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

	if _, err2 := models.AddChannelCompany(&v); err2 != nil {
		c.RespJSON(bean.CODE_Bad_Request, err2.Error())
		return
	}
	c.RespJSONData(v)
}

// GetOne ...
// @Title Get One
// @Description get ChannelCompany by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ChannelCompany
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ChannelCompanyController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetChannelCompanyById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get ChannelCompany
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ChannelCompany
// @Failure 403
// @router / [get]
func (c *ChannelCompanyController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)

}

func (c *ChannelCompanyController) getAll() (total int64, m []*models.ChannelCompany, errCode int, err error) {
	_, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_COMPANY, nil)
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
	//total, err = tool.GetAllByFilterWithTotal(new(models.ChannelCompany), &m, filter)

	var offset int
	var company_id string
	var channel_code string
	var cooperate_state int
	var person int

	offset, err = c.GetInt("offset", 0)
	company_id = c.GetString("company_id", "")
	channel_code = c.GetString("channel_code", "")
	cooperate_state, err = c.GetInt("cooperate_state", 0)
	person, err = c.GetInt("person", 0)

	m, total, err = models.GetAllChannelCompaniesWithTotal(offset, company_id, channel_code, cooperate_state, person)

	if err != nil {
		errCode = bean.CODE_Bad_Request
		return
	}

	if len(m) == 0 {
		return
	}

	// add the company info
	for _, v := range m {
		models.GetChannelCompanyAdditionalInfo(v)
	}

	return
}

// GetAllCompanyName ...
// @Title Get All CompanyName with channel_code
// @Description get ChannelCompany
// @Success 200 {object} models.ChannelCompany
// @Failure 403
// @router /companyname/ [get]
func (c *ChannelCompanyController) GetAllCompanyName() {
	filter, err := tool.BuildFilter(c.Controller, 0)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	filter.Fields = append(filter.Fields, "Id", "CompanyId", "ChannelCode")

	var ss []models.ChannelCompany
	_, err = tool.GetAllByFilterWithTotal(new(models.ChannelCompany), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	if len(ss) == 0 {
		return
	}
	models.AddCompanyInfo(&ss)

	c.RespJSONData(ss)

}

// GetContractState ...
// @Title Get All CompanyName
// @Description get ChannelCompany
// @Success 200 {object} models.ChannelCompany
// @Failure 403
// @router /state/ [get]
func (c *ChannelCompanyController) GetContractState() {
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
	models.AddCompanyIdInfo(&ss, 2)

	c.RespJSONData(ss)
}

// Put ...
// @Title Put
// @Description update the ChannelCompany
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ChannelCompany	true		"body for ChannelCompany content"
// @Success 200 {object} models.ChannelCompany
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ChannelCompanyController) Put() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_CHANNEL_COMPANY, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.ChannelCompany{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	if v.YunduanResponsiblePerson == 0 && v.YouliangResponsiblePerson == 0 {
		c.RespJSON(bean.CODE_Params_Err, "Must chose one Person Responsible at least.")
		return
	}
	//编辑渠道商信息 添加更新人、时间
	v.UpdateUserID = c.Uid()
	v.UpdateTime = time.Now().Unix()

	if v.YunduanResponsiblePerson == 0 && v.YouliangResponsiblePerson == 0 {
		c.RespJSON(bean.CODE_Params_Err, "商务负责人至少填写一个!")
		return
	}

	if err := models.UpdateChannelCompanyById(&v, where); err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(v)
}

// Delete ...
// @Title Delete
// @Description delete the ChannelCompany
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ChannelCompanyController) Delete() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_DELETE, bean.PMSM_CHANNEL_COMPANY, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteChannelCompany(id); err == nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData("OK")
}

// Get Channels ...
// @Title Get Channels
// @Description get the Channels
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /channels [get]
func (c *ChannelCompanyController) GetChannels() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_COMPANY, []string{"cp"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	filter := &tool.Filter{
		Where:  where,
		Fields: []string{"Cp", "Name"},
	}

	var ss []models.Channel
	err = tool.GetAllByFilter(new(models.Channel), &ss, filter)

	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(ss)
}
