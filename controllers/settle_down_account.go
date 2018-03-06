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

// cp结算账单
type SettleDownAccountController struct {
	BaseController
}

// URLMapping ...
func (c *SettleDownAccountController) URLMapping() {

}

// Post ...
// @Title Post
// @Description 新增一个已结算账单
// @Param	body		body 	models.SettleDownAccount	true		"body for SettleDownAccounts content"
// @Success 201 {int} models.SettleDownAccount
// @Failure 403 body is empty
// @router / [post]
func (c *SettleDownAccountController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_SETTLE_DOWN_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var v models.SettleDownAccount
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	err = c.Validate(&v)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	if v.Amount == 0 {
		c.RespJSON(bean.CODE_Params_Err, "amount can't be 0")
		return
	}
	//结算管理 更新人
	//user,_:=models.GetUserById(c.Uid())
	//v.UpdateSettleUser=user.Nickname
	v.UpdateSettleUserID = c.Uid()
	v.UpdateSettleTime = time.Now().Unix()

	if _, err := models.AddSettleDownAccounts(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(v)
}

//// 获取未结算的账单信息 ...
//// @Title Get One
//// @Description 获取未结算的账单信息
//// @Param	companiy_ids		query 	string	false		"开发商id"
//// @Param   start_time		query 	number	false		"开始时间"
//// @Param	end_time		query 	number	false		"结束时间"
//// @Success 200 {object} models.bean.NotSettled
//// @Failure 403 :id is empty
//// @router /pre [get]
//func (c *SettleDownAccountController) GetPreAccount() {
//	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_SETTLE_DOWN_ACCOUNT, []string{"game_id"})
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err)
//		return
//	}
//
//	idStr := c.GetString("companiy_ids")
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
//		"company_id__in": ids,
//	}
//	tool.InjectPermissionWhere(where, &idsWhere)
//	v, err := models.GetNotSettleAccount(idsWhere["company_id__in"], start_time, end_time)
//	if err != nil {
//		c.RespJSON(bean.CODE_Not_Found, err)
//		return
//	}
//	c.RespJSONData(v)
//}

// GetAll ...
// @Title Get All
// @Description 查看已结算账单
// @Param	company_ids		query 	string	false		"发行商id"
// @Param   start_time		query 	number	false		"开始时间"
// @Param	end_time		query 	number	false		"结束时间"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.SettleDownAccount
// @Failure 403
// @router / [get]
func (c *SettleDownAccountController) GetAll() {
	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)
}

// @router /download [get]
func (c *SettleDownAccountController) DownLoad() {
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
		r["body_my"] = bodyMy

		tStart := time.Unix(int64(v.SettleTime), 0).Format("2006-01-02")
		tCreated := time.Unix(int64(v.CreatedTime), 0).Format("2006-01-02 15:04:05")
		updatedTime := ""
		if v.UpdateSettleTime!=0{
			updatedTime = time.Unix(v.UpdateSettleTime, 0).Format("2006-01-02 15:04:05")
		}

		r["time"] = tStart
		r["created_time"] = tCreated
		if v.User != nil {
			r["username"] = v.User.Nickname
		}
		if v.Company != nil {
			r["company"] = v.Company.Name
		}

		if v.UpdateUser != nil {
			r["updater"] = v.UpdateUser.Nickname + " " + updatedTime
		} else {
			r["updater"] = updatedTime
		}
		rs[i] = r
	}

	cols := []string{"time", "company", "body_my","amount", "username", "updater"}
	maps := map[string]string{
		cols[0]: "时间",
		cols[1]: "发行商",
		cols[2]: "我方主体",
		cols[3]: "金额",
		cols[4]: "结算人",
		cols[5]: "更新人",
	}
	c.RespExcel(rs, "结算单", cols, maps)
}

func (c *SettleDownAccountController) getAll() (total int64, ss []models.SettleDownAccount, errCode int, err error) {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_SETTLE_DOWN_ACCOUNT, []string{"game_id"})
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

	ids := []interface{}{}
	for _, v := range strings.Split(idStr, ",") {
		if v != "" {
			ids = append(ids, v)
		}
	}
	filter.Where = map[string][]interface{}{}
	if len(ids) != 0 {
		filter.Where["company_id__in"] = ids
	}
	if bodyMy != 0 {
		filter.Where["body_my"] = []interface{}{bodyMy}
	}
	if startTime != 0 {
		filter.Where["settle_time__gte"] = []interface{}{startTime}
	}
	if endTime != 0 {
		filter.Where["settle_time__lte"] = []interface{}{endTime}
	}
	tool.InjectPermissionWhere(where, &filter.Where)

	filter.Sortby = []string{"Id"}
	filter.Order = []string{"desc"}
	ss = []models.SettleDownAccount{}
	total, err = tool.GetAllByFilterWithTotal(new(models.SettleDownAccount), &ss, filter)
	if err != nil {
		errCode = bean.CODE_Bad_Request
		return
	}

	models.GroupSettleAddCompanyInfo(&ss)
	models.AddUserInfo4Settle(&ss)
	models.AddUpdateUserInfoSettle(&ss)
	return
}

// 获取某发行商下的所有游戏
// @Title 获取某发行商下的所有游戏
// @Description 获取某发行商下的所有游戏
// @Param	company_id		query 	string	false		"发行商id"
// @Success 200 {object} models.Game
// @Failure 403
// @router /games [get]
func (c *SettleDownAccountController) GetGames() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_SETTLE_DOWN_ACCOUNT, []string{"game_id"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	id := c.GetString("company_id", "")
	if id == "" {
		c.RespJSON(bean.CODE_Params_Err, "company_id can't be empty")
		return
	}

	where["Development"] = []interface{}{id}
	filter := &tool.Filter{
		Where:  where,
		Fields: []string{"GameId", "GameName"},
	}

	ss := []models.Game{}
	err = tool.GetAllByFilter(new(models.Game), &ss, filter)

	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(ss)
}

// 获取发行商选项
// @Title 获取发行商选项
// @Description 获取发行商选项
// @Success 200 {object} models.Company
// @Failure 403
// @router /companies [get]
func (c *SettleDownAccountController) GetCompanies() {
	ss, err := models.GetAllDistributionCompanies()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err)
		return
	}

	c.RespJSONData(ss)
}

// @Title 更新结算单
// @Description 更新结算单
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.SettleDownAccount	true		"body for SettleDownAccounts content"
// @Success 200 {object} models.SettleDownAccount
// @Failure 403 :id is not int
// @router /:id [put]
func (c *SettleDownAccountController) Put() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_SETTLE_DOWN_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.SettleDownAccount{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	v.UpdateSettleUserID = c.Uid()
	v.UpdateSettleTime = time.Now().Unix()

	if err := models.UpdateSettleDownAccountsById(&v); err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(v)
}

// @Title 删除结算单
// @Description 删除结算单
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *SettleDownAccountController) Delete() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_DELETE, bean.PMSM_SETTLE_DOWN_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	m := models.SettleDownAccount{Id: id}
	if err := orm.NewOrm().Read(&m); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	if err := models.DeleteSettleDownAccounts(id); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	err = models.CompareAndAddOperateLog(m, nil, c.Uid(), bean.OPP_CP_REMIT, int(id), bean.OPA_DELETE)

	c.RespJSONData("OK")
}

// @Title 获取结算单
// @Description 获取结算单
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {object} models.SettleDownAccount
// @router /:id [get]
func (c *SettleDownAccountController) Get() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_SETTLE_DOWN_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetSettleDownAccountsById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(v)
}
