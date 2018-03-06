package controllers

import (
	"encoding/json"
	"fmt"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"strconv"
	"strings"
	"time"
)

// cp对账单
type VerifyCpController struct {
	BaseController
}

// URLMapping ...
func (c *VerifyCpController) URLMapping() {
}

// @Title 添加对账单
// @Description 添加对账单
// @Param	body		body 	models.VerifyChannel	true		"body for VerifyChannel content"
// @Success 201 {int} models.VerifyChannel
// @Failure 403 body is empty
// @router / [post]
func (c *VerifyCpController) Post() {
	var v models.VerifyCp

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	v.CreatedUserId = c.Uid()
	v.UpdatedUserId = c.Uid()
	if _, err := models.AddVerifyCp(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData("OK")
}

// @Title 简单统计信息
// @Description 获取昨天今天的新增账单
// @Success 200 {object} models.VerifyChannel
// @router /simple_statistics [get]
func (c *VerifyCpController) SimpleStatistics() {
	v, err := models.GetCpSimpleStatistics()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(v)
}

// @Title Get One
// @Description get VerifyChannel by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.VerifyChannel
// @Failure 403 :id is empty
// @router /:id [get]
func (c *VerifyCpController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetVerifyCpById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(v)
}

// @Title 获取所有对账单
// @Description get VerifyChannel
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.VerifyChannel
// @Failure 403
// @router / [get]
func (c *VerifyCpController) GetAll() {
	data, total, err := c.getAll()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONDataWithTotal(data, total)
}

// @router /download [get]
func (c *VerifyCpController) DownLoad() {
	// todo 下载,2017-03-16
	c.Ctx.Input.SetParam("limit", "0")
	c.Ctx.Input.SetParam("offset", "0")

	ss, _, err := c.getAll()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	rs := make([]map[string]interface{}, len(ss))
	for i, v := range ss {
		r := map[string]interface{}{
			"amount_my":       fmt.Sprintf("%.2f", v.AmountMy),
			"amount_opposite": fmt.Sprintf("%.2f", v.AmountOpposite),
			"amount_payable":  fmt.Sprintf("%.2f", v.AmountPayable),
			"amount_remit":    fmt.Sprintf("%.2f", v.AmountRemit),
		}
		date := v.Date
		tCreated := time.Unix(int64(v.CreatedTime), 0).Format("2006-01-02 15:04:05")
		r["time"] = date
		r["create_time"] = tCreated
		r["verify_time"] = time.Unix(int64(v.VerifyTime), 0).Format("2006-01-02 15:04:05")
		r["update_time"] = time.Unix(int64(v.UpdatedTime), 0).Format("2006-01-02 15:04:05")
		//if v.CreatedUser != nil {
		//	r["create_username"] = v.CreatedUser.Nickname
		//}
		if v.UpdatedUser != nil {
			r["update_username"] = v.UpdatedUser.Nickname
		}
		if v.VerifyUser != nil {
			r["verify_username"] = v.VerifyUser.Nickname
		}
		if v.Company != nil {
			r["company"] = v.Company.Name
		}
		//r["games"] = v.GameStr
		////if v.GameStr != "" {
		////	temp := ""
		////	for _, value := range v.Games {
		////		temp += value.GameName + ";"
		////		r["games"] = temp
		////	}
		////}
		//r["status"] = models.ChanStatusCode2String[v.Status]
		rs[i] = r
	}

	cols := []string{"time", "cp", "games", "game_num", "status", "amount_my", "amount_opposite",
	                 "amount_payable", "amount_remit", "verify_username", "create_username", "create_time",
	                 "update_username", "update_time"}
	maps := map[string]string{
		cols[0]:  "账单日期",
		cols[1]:  "渠道",
		cols[2]:  "游戏",
		cols[3]:  "游戏数量",
		cols[4]:  "状态",
		cols[5]:  "我方流水",
		cols[6]:  "对方流水",
		cols[7]:  "应付金额",
		cols[8]:  "回款金额",
		cols[9]:  "对账人",
		cols[10]: "创建人",
		cols[11]: "创建时间",
		cols[12]: "更新人",
		cols[13]: "更新时间",
	}
	c.RespExcel(rs, "渠道对账单", cols, maps)
}

func (c *VerifyCpController) getAll() (data []models.VerifyCp, total int64, err error) {
	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		return
	}

	verifyUserId, _ := c.GetInt("user_id")
	bodyMy, _ := c.GetInt("body_my")

	company_ids := c.GetString("company_ids")
	start := c.GetString("start", "2006-01")
	end := c.GetString("end", "2200-01")
	status, _ := c.GetInt("status", 0)

	var companies []interface{}
	if company_ids != "" {
		for _, v := range strings.Split(company_ids, ",") {
			companies = append(companies, v)
		}
	}
	if len(companies) != 0 {
		filter.Where["company_id__in"] = companies
	}
	if start != "" {
		filter.Where["date__gte"] = []interface{}{start}
	}
	if end != "" {
		filter.Where["date__lte"] = []interface{}{end}
	}
	if status != 0 {
		filter.Where["status"] = []interface{}{status}
	}
	if bodyMy != 0 {
		filter.Where["body_my"] = []interface{}{bodyMy}
	}
	if verifyUserId != 0 {
		filter.Where["verify_user_id"] = []interface{}{verifyUserId}
	}

	filter.Sortby = []string{"id"}
	filter.Order = []string{"desc"}

	var vs []models.VerifyCp
	total, err = tool.GetAllByFilterWithTotal(new(models.VerifyCp), &vs, filter)
	if err != nil {
		return
	}
	models.AddPreVerifyCpForVerify(&vs)
	models.AddCompanyForVerifyCp(&vs)
	models.AddVerifyAndUpdateUserForVerifyCp(&vs)
	data = vs
	return
}

// @Title 更新对账单
// @Description 更新对账单
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.VerifyChannel	true		"body for VerifyChannel content"
// @Success 200 {object} models.VerifyChannel
// @Failure 403 :id is not int
// @router /:id [put]
func (c *VerifyCpController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.VerifyCp{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	v.UpdatedUserId = c.Uid()
	if err := models.UpdateVerifyCpById(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData("OK")
}

// @Title Delete
// @Description delete the VerifyChannel
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *VerifyCpController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteVerifyCp(id,c.Uid()); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData("OK")
}

// @Title 获取未对账的发行商
// @Description  获取未对账的发行商
// @Param	body_my		query 	string	true		"我方主体"
// @router /not_verify_company [get]
func (c *VerifyCpController) GetNotVerifyCompany() {
	bodyMy, _ := c.GetInt("body_my", 0)
	if bodyMy == 0 {
		c.RespJSON(bean.CODE_Params_Err, "body_my can't be empty")
		return
	}

	channels, err := models.GetNotVerifyCp(bodyMy)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(channels)
	return
}

// @Title 获取未对账渠道的时间
// @Description delete the VerifyChannel
// @Param	body_my		query 	string	true		"我方主体"
// @Param	company_id		query 	int	true		"发行商id"
// @router /not_verify_time [get]
func (c *VerifyCpController) GetNotVerifyCpTime() {
	bodyMy, _ := c.GetInt("body_my", 0)
	if bodyMy == 0 {
		c.RespJSON(bean.CODE_Params_Err, "body_my can't be empty")
		return
	}
	companyId, _ := c.GetInt("company_id", 0)
	if companyId == 0 {
		c.RespJSON(bean.CODE_Params_Err, "company_id can't be empty")
		return
	}

	date, err := models.GetNotVerifyCpTime(bodyMy, companyId)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(date)
	return
}

// @Title 获取未对账渠道的游戏
// @Description delete the VerifyChannel
// @Param	body_my		query 	string	true		"我方主体"
// @Param	company_id		query 	int	true		"发行商id"
// @Param	month		query 	string	true		"月份"
// @router /not_verify_game [get]
func (c *VerifyCpController) GetNotVerifyCpGame() {
	bodyMy, _ := c.GetInt("body_my", 0)
	if bodyMy == 0 {
		c.RespJSON(bean.CODE_Params_Err, "body_my can't be empty")
		return
	}
	companyId, _ := c.GetInt("company_id", 0)
	if companyId == 0 {
		c.RespJSON(bean.CODE_Params_Err, "companyId can't be empty")
		return
	}

	month := c.GetString("month")
	if month == "" {
		c.RespJSON(bean.CODE_Params_Err, "month can't be empty")
		return
	}

	date, err := models.GetNotVerifyCpGame(bodyMy, companyId, month)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(date)
	return
}

// 获取发行商选项
// @Title 获取发行商选项
// @Description 获取发行商选项
// @Success 200 {object} models.Company
// @Failure 403
// @router /companies [get]
func (c *VerifyCpController) GetCompanies() {
	ss, err := models.GetAllDistributionCompanies()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err)
		return
	}

	c.RespJSONData(ss)
}

// @Title 获取Cp未对账的信息
// @Description 获取Cp未对账的信息(日期，主体，渠道，流水)
// @router /not_verify_info [get]
func (c *VerifyCpController) GetNotVerifyInfo() {
	limit, _ := c.GetInt("limit", 20)
	offset, _ := c.GetInt("offset", 0)

	count, notVerify, err := models.GetCpNotVerifyInfo(limit, offset)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONDataWithTotal(notVerify, count)
	return
}

//// @Title 从老的对账单生成一个新对账单
//// @Description 从老的对账单生成一个新对账单
//// @Param	id		path 	string	true		"The id you want to delete"
//// @Success 200 {string} delete success!
//// @Failure 403 id is empty
//// @router /migration/:id [get]
//func (c *VerifyCpController) Migration() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	if err := models.MigrationVerifyCp(id); err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//
//	c.RespJSONData("OK")
//}
//
//// @Title 从老的对账单生成全部对账单
//// @Description 从老的对账单生成全部对账单
//// @Success 200 {string} delete success!
//// @router /migration_all [get]
//func (c *VerifyCpController) MigrationAll() {
//	if err := models.MigrationVerifyCpAll(); err != nil {
//		return
//	}
//
//	c.RespJSONData("OK")
//}
