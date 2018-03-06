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
	"github.com/astaxie/beego/orm"
)

// 渠道回款账单
type RemitDownAccountController struct {
	BaseController
}

// URLMapping ...
func (c *RemitDownAccountController) URLMapping() {

}

// Post ...
// @Title Post
// @Description 新增一个已回款账单
// @Param	body		body 	models.swagger.InputRemitDownAccount	true		"body for RemitDownAccounts content"
// @Success 201 {int} models.RemitDownAccount
// @Failure 403 body is empty
// @router / [post]
func (c *RemitDownAccountController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_REMIT_DOWN_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var v models.RemitDownAccount
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if err := c.Validate(&v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if v.Amount == 0 {
		c.RespJSON(bean.CODE_Params_Err, "amount can't be empty")
		return
	}
	//添加回款人 更新人、时间
	v.UpdateChannelUserID = c.Uid()
	v.UpdateRemitTime = time.Now().Unix()

	if _, err := models.AddRemitDownAccounts(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData("OK")
}

//// 获取未回款的账单信息 ...
//// @Title Get One
//// @Description 获取未回款的账单信息
//// @Param	company_ids		query 	string	false		"回款主体ids:1,2,3,6"
//// @Param   start_time		query 	number	false		"开始时间"
//// @Param	end_time		query 	number	false		"结束时间"
//// @Success 200 {object} models.bean.NotSettled
//// @router /pre [get]
//func (c *RemitDownAccountController) GetPreAccount() {
//	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_REMIT_DOWN_ACCOUNT, []string{"channel_code"})
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	idStr := c.GetString("company_ids")
//	start_time, _ := c.GetInt("start_time")
//	end_time, _ := c.GetInt("end_time")
//
//	ids := []interface{}{}
//	for _, v := range strings.Split(idStr, ",") {
//		if v != "" {
//			ids = append(ids, v)
//		}
//	}
//
//	idsWhere := map[string][]interface{}{
//		"remit_company_id__in": ids,
//	}
//	tool.InjectPermissionWhere(where, &idsWhere)
//	v, err := models.GetNotRemitAccount(idsWhere["remit_company_id__in"], start_time, end_time)
//	if err != nil {
//		c.RespJSON(bean.CODE_Not_Found, err.Error())
//		return
//	}
//	c.RespJSONData(v)
//}

// 根据回款主体公司id,获取未全部回款的账单信息 ...
// @Title Get One
// @Description 获取未全部回款的账单信息
// @Param	company_id		query 	string	false		"回款主体id"
// @Success 200 {object} models.bean.NotSettled
// @router /pre [get]
func (c *RemitDownAccountController) GetPreAccount() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_REMIT_DOWN_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	company_id, _ := c.GetInt("company_id", 0)
	if company_id == 0 {
		c.RespJSON(bean.CODE_Bad_Request, "参数错误")
	}

	v, err := models.GetNotRemitAccount(company_id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}
	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get 查看已回款账单
// @Param   body_my  		query 	number	false		"我方主体"
// @Param	company_ids		query 	string	false		"游戏id"
// @Param   start_time		query 	number	false		"开始时间"
// @Param	end_time		query 	number	false		"结束时间"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.RemitDownAccount
// @Failure 403
// @router / [get]
func (c *RemitDownAccountController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}

	c.RespJSONDataWithTotal(ss, total)
}

// @router /download [get]
func (c *RemitDownAccountController) DownLoad() {
	c.Ctx.Input.SetParam("limit", "0")
	c.Ctx.Input.SetParam("offset", "0")

	_, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}

	rs := make([]map[string]interface{}, len(ss))
	for i, v := range ss {
		r := map[string]interface{}{
			"amount": fmt.Sprintf("%.2f", v.Amount),
		}
		bodyMy := "未知"
		if v.BodyMy == 2 {
			bodyMy = "有量"
		} else if v.BodyMy == 1 {
			bodyMy = "云端"
		}
		remitTime := time.Unix(int64(v.RemitTime), 0).Format("2006-01-02")
		tCreated := time.Unix(int64(v.CreatedTime), 0).Format("2006-01-02 15:04:05")
		updatedTime := ""
		if v.UpdateRemitTime != 0 {
			updatedTime = time.Unix(v.UpdateRemitTime, 0).Format("2006-01-02 15:04:05")
		}
		r["time"] = remitTime
		r["created_time"] = tCreated

		r["body_my"] = bodyMy

		if v.User != nil {
			r["remit_username"] = v.User.Nickname
		}
		if v.RemitCompany != nil {
			r["remit_company"] = v.RemitCompany.Name
		}

		if v.UpdateUser != nil {
			r["updater"] = v.UpdateUser.Nickname + " " + updatedTime
		} else {
			r["updater"] = updatedTime
		}

		rs[i] = r
	}

	cols := []string{"time", "remit_company", "body_my", "amount", "remit_username", "updater"}
	maps := map[string]string{
		cols[0]: "时间",
		cols[1]: "回款主体",
		cols[2]: "我方主体",
		cols[3]: "金额",
		cols[4]: "回款人",
		cols[5]: "更新人",
	}

	c.RespExcel(rs, "回款单", cols, maps)
}

func (c *RemitDownAccountController) getAll() (total int64, ss []models.RemitDownAccount, errCode int, err error) {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_REMIT_DOWN_ACCOUNT, nil)
	if err != nil {
		errCode = bean.CODE_Forbidden
		return
	}

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		errCode = bean.CODE_Params_Err
		return
	}

	idStr := c.GetString("company_ids")
	startTime, _ := c.GetInt("start_time")
	endTime, _ := c.GetInt("end_time")
	bodyMy, _ := c.GetInt("body_my")

	var ids []interface{}
	for _, v := range strings.Split(idStr, ",") {
		if v != "" {
			ids = append(ids, v)
		}
	}
	filter.Where = map[string][]interface{}{}
	if len(ids) != 0 {
		filter.Where["remit_company_id__in"] = ids
	}

	if bodyMy != 0 {
		filter.Where["body_my"] = []interface{}{bodyMy}
	}

	if startTime != 0 {
		filter.Where["remit_time__gte"] = []interface{}{startTime}
	}
	if endTime != 0 {
		filter.Where["remit_time__lte"] = []interface{}{endTime}
	}

	tool.InjectPermissionWhere(where, &filter.Where)
	filter.Sortby = []string{"Id"}
	filter.Order = []string{"desc"}
	ss = []models.RemitDownAccount{}
	total, err = tool.GetAllByFilterWithTotal(new(models.RemitDownAccount), &ss, filter)
	if err != nil {
		errCode = bean.CODE_Bad_Request
		return
	}

	models.AddRemitCompanyInfo4Remit(&ss)
	models.AddUserInfo4Remit(&ss)
	models.AddUpdateUserInfoRemit(&ss)
	//models.AddRemitDownDetailInfo(&ss)
	return
}

// @Title 更新对账单
// @Description update the 更新对账单
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.swagger.InputRemitDownAccount	true		"body for RemitDownAccounts content"
// @Success 200 {object} models.RemitDownAccount
// @Failure 403 :id is not int
// @router /:id [put]
func (c *RemitDownAccountController) Put() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_REMIT_DOWN_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.RemitDownAccount{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	//编辑渠道汇款 添加更新人、时间
	v.UpdateChannelUserID = c.Uid()
	v.UpdateRemitTime = time.Now().Unix()
	if err := models.UpdateRemitDownAccountsById(&v, where); err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(v)
}

// @Title 获取一个对账单
// @Description update the 获取一个对账单
// @Param	id		path 	string	true		"The id you want to get"
// @Success 200 {object} models.RemitDownAccount
// @Failure 403 :id is not int
// @router /:id [get]
func (c *RemitDownAccountController) Get() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_REMIT_DOWN_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetRemitDownAccountsById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(v)
}

// Delete ...
// @Title Delete
// @Description delete the RemitDownAccounts
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *RemitDownAccountController) Delete() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_DELETE, bean.PMSM_REMIT_DOWN_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	m := models.RemitDownAccount{Id: id}
	if err := orm.NewOrm().Read(&m); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	if err := models.DeleteRemitDownAccounts(id, where); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	err = models.CompareAndAddOperateLog(m, nil, c.Uid(), bean.OPP_CHANNEL_REMIT, int(id), bean.OPA_DELETE)

	c.RespJSONData("OK")
}

// unused in version 2.0
// @Title 获取渠道
// @Description 获取渠道
// @Success 200 {object} models.Channel
// @Failure 403
// @router /channel [get]
func (c *RemitDownAccountController) GetChannel() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_REMIT_DOWN_ACCOUNT, []string{"channel_code"})
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
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(ss)
}

// @Title 获取所有回款主体
// @Description 获取所有回款主体
// @Success 200 {object} models.Company
// @Failure 403
// @router /remitcompanies [get]
func (c *RemitDownAccountController) GetRemitCompanies() {
	ss, err := models.GetRemitCompaniesByChannel("*")
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(ss)
}
