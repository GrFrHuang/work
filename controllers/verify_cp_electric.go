package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"strings"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"time"
	"github.com/astaxie/beego/orm"
)

// cp电子对账单
type VerifyCpElectricController struct {
	BaseController
}

// URLMapping ...
func (c *VerifyCpElectricController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create VerifyCpElectric
// @Param	body		body 	models.VerifyCpElectric	true		"body for VerifyCpElectric content"
// @Success 201 {int} models.VerifyCpElectric
// @Failure 403 body is empty
// @router / [post]
func (c *VerifyCpElectricController) Post() {
	var v models.VerifyCpElectric
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	v.UpdateUser = c.Uid()
	v.UpdateTime = time.Now().Unix()
	if id, err := models.AddVerifyCpElectric(&v); err == nil {
		for _, detail := range v.Games{
			detail.ElectricId = id
			if _, er := models.AddVerifyCpElectircDetail(&detail); er != nil{
				c.RespJSON(bean.CODE_Bad_Request, err.Error())
				return
			}
		}
	}
	c.RespJSONData("OK")
}

// GetOne ...
// @Title Get One
// @Description get VerifyCpElectric by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.VerifyCpElectric
// @Failure 403 :id is empty
// @router /:id [get]
func (c *VerifyCpElectricController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetVerifyCpElectricById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get VerifyCpElectric
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.VerifyCpElectric
// @Failure 403
// @router / [get]
func (c *VerifyCpElectricController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss,total)
}

func (c *VerifyCpElectricController) getAll() (total int64, ss []models.VerifyCpElectric, errCode int, err error)  {
	_, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CP_VERIFY_ACCOUNT, nil)
	if err != nil {
		errCode = bean.CODE_Forbidden
		return
	}

	bodyMy, _ := c.GetInt("body_my")
	company_ids := c.GetString("company_ids")
	start := c.GetString("start", "")
	end := c.GetString("end", "")

	var limit int = 15
	var offset int = 0
	if li, limitErr := c.GetInt("limit"); limitErr == nil{
		limit = li
	}
	if off, offsetErr := c.GetInt("offset"); offsetErr == nil{
		offset = off
	}

	companies := []string{}
	if company_ids != "" {
		for _, v := range strings.Split(company_ids, ",") {
			companies = append(companies, v)
		}
	}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("verify_cp_electric.*").From("verify_cp_electric")

	if len(companies) != 0 {
		qb = qb.Where("company_id").In(companies...)
	}
	if bodyMy != 0 {
		if strings.Contains(qb.String(), "WHERE"){
			qb = qb.And("body_my = ? ")
		} else {
			qb = qb.Where("body_my = ? ")
		}
	}
	ids := []string{}
	if start != "" && end != ""{
		orm.NewOrm().Raw("SELECT DISTINCT(electric_id) FROM verify_cp_electric_detail WHERE `date`>=? AND `date`<=?",
			start, end).QueryRows(&ids)
		if strings.Contains(qb.String(), "WHERE"){
			qb = qb.And("id").In(ids...)
		} else {
			qb = qb.Where("id").In(ids...)
		}
	}
	qb = qb.OrderBy("verify_cp_electric.id").Desc()
	sql_total := qb.String()

	o := orm.NewOrm()
	if bodyMy == 0 {
		_, _ = o.Raw(sql_total).QueryRows(&ss)
	}else if bodyMy != 0  {
		_, _ = o.Raw(sql_total, bodyMy).QueryRows(&ss)
	}

	qb = qb.Limit(limit).Offset(offset)
	sql := qb.String()
	if bodyMy == 0 {
		_, _ = o.Raw(sql).QueryRows(&ss)
	}else if bodyMy != 0  {
		_, _ = o.Raw(sql, bodyMy).QueryRows(&ss)
	}

	models.AddCompanyForVerifyCpElectric(&ss)
	models.AddUpdateUserForVerifyCpElectric(&ss)
	models.AddVerifyDate(&ss)

	return
}

// @Title 获取已对账,未开票发行商
// @Description
// @Param	body_my		query 	string	true		"我方主体"
// @router /company [get]
func (c *VerifyCpElectricController) GetVerifyCompany() {
	bodyMy, _ := c.GetInt("body_my", 0)
	if bodyMy == 0 {
		c.RespJSON(bean.CODE_Params_Err, "缺少参数，我方主体:body_my")
		return
	}

	channels, err := models.GetVerifyCp(bodyMy)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(channels)
	return
}

// @Title 获取对账的时间
// @Description
// @Param	body_my		query 	string	true		"我方主体"
// @Param	company_id		query 	int	true		"发行商id"
// @router /verify_time [get]
func (c *VerifyCpElectricController) GetVerifyCpTime() {
	bodyMy, _ := c.GetInt("body_my", 0)
	if bodyMy == 0 {
		c.RespJSON(bean.CODE_Params_Err, "缺少参数，我方主体:body_my")
		return
	}
	companyId, _ := c.GetInt("company_id", 0)
	if companyId == 0 {
		c.RespJSON(bean.CODE_Params_Err, "缺少参数，发行商id:company_id")
		return
	}

	date, err := models.GetVerifyCpTime(bodyMy, companyId)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(date)
	return
}

// @Title 获取cp对账的游戏
// @Description delete the VerifyChannel
// @Param	body_my		query 	string	true		"我方主体"
// @Param	company_id		query 	int	true		"发行商id"
// @Param	month		query 	string	true		"月份"
// @router /verify_game [get]
func (c *VerifyCpElectricController) GetVerifyCpGame() {
	companyId, _ := c.GetInt("company_id", 0)
	if companyId == 0 {
		c.RespJSON(bean.CODE_Params_Err, "缺少参数，发行商id:company_id")
		return
	}

	month := c.GetString("month")
	if month == "" {
		c.RespJSON(bean.CODE_Params_Err, "缺少参数，对账日期:month")
		return
	}

	data, err := models.GetVerifyCpGame(companyId, month)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	models.SetGameName(&data)
	models.SetRate(&data)
	c.RespJSONData(data)
	return
}

// Put ...
// @Title Put
// @Description update the VerifyCpElectric
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.VerifyCpElectric	true		"body for VerifyCpElectric content"
// @Success 200 {object} models.VerifyCpElectric
// @Failure 403 :id is not int
// @router /:id [put]
func (c *VerifyCpElectricController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.VerifyCpElectric{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateVerifyCpElectricById(&v); err == nil {
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
// @Description delete the VerifyCpElectric
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *VerifyCpElectricController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteVerifyCpElectric(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
