// unused
// 已经用verify_cp代替
package controllers

import (
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
)

// cp对账
type CpVerifyAccountController struct {
	BaseController
}

// URLMapping ...
func (c *CpVerifyAccountController) URLMapping() {

}

// unused
// 获取游戏选项
// @Title 获取游戏
// @Description 获取游戏
// @Success 200 {object} models.Game
// @Failure 403
// @router /games [get]
func (c *CpVerifyAccountController) GetGames() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, []string{"game_id"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

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
func (c *CpVerifyAccountController) GetCompanies() {
	ss, err := models.GetAllDistributionCompanies()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err)
		return
	}

	c.RespJSONData(ss)
}

//// GetAll ...
//// @Title 获取未对账的数据
//// @Description 获取未对账的数据
//// @Param	company_ids		query	string	false	"发行商ID，多个用','隔开"
//// @Param	start		query	string	false	"开始时间 2006-01-02 格式"
//// @Param	end		query	string	false	"结束时间 2006-01-02 格式"
//// @Success 200 {object} models.bean.NoCpVerify
//// @Failure 403
//// @router /pre [get]
//func (c *CpVerifyAccountController) GetNoVerifyAccount() {
//	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CP_VERIFY_ACCOUNT, nil)
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//
//	idsStr := c.GetString("company_ids")
//	start := c.GetString("start", "2006-01-02")
//	end := c.GetString("end", "2200-01-01")
//	var ids []interface{}
//	for _, value := range strings.Split(idsStr, ",") {
//		if value != "" {
//			ids = append(ids, value)
//		}
//	}
//	res, err := models.SumCpNotVerifyByCompany(start, end, ids)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//
//	c.RespJSONData(res)
//}
//
//// @Title 获取未对账的公司列表
//// @Description 无
//// @Success 200 {object} models.Company
//// @router /notcompanies [get]
//func (c *CpVerifyAccountController) GetNotVerifyCompanies() {
//	res, err := models.GetNotVerifyCompanies()
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	c.RespJSONData(res)
//}
//
//// @Title 根据开发商获取未对账游戏的时间区间
//// @Description 根据开发商获取未对账游戏的时间区间
//// @Param	company_id		query	string	false	"发行商ID"
//// @Success 200 {object} models.bean.CpVerifyAccount
//// @router /notdate [get]
//func (c *CpVerifyAccountController) GetNotVerifyDate() {
//	companyId, _ := c.GetInt("company_id")
//	if companyId == 0 {
//		c.RespJSON(bean.CODE_Params_Err, "company_id can't be empty")
//		return
//	}
//	res, err := models.GetNotVerifyDateByComp(companyId)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	c.RespJSONData(res)
//}
//
//// GetGames ...
//// @Title 获取某发行商下时间范围内未对账的游戏
//// @Description 获取某发行商下时间范围内未对账的游戏
//// @Param	my_body  		query	string	false	"公司主体"
//// @Param	company_id		query	string	false	"发行商ID"
//// @Param	start_time		query	string	false	"开始时间"
//// @Param	end_time		query	string	false	"结束时间"
//// @Success 200 {object} models.CpVerifyAccount
//// @Failure 403
//// @router /notgames [get]
//func (c *CpVerifyAccountController) GetNotGames() {
//	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CP_VERIFY_ACCOUNT, []string{"game_id"})
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	myBody, _ := c.GetInt("my_body")
//	companyId, _ := c.GetInt("company_id")
//	if companyId == 0 {
//		c.RespJSON(bean.CODE_Params_Err, "company_id can't be empty")
//		return
//	}
//
//	start := c.GetString("start_time", "2000-01-01")
//	end := c.GetString("end_time", "2200-01-01")
//	res, err := models.GetCpNotVerifyGames(start, end, companyId, myBody, where, )
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	c.RespJSONData(res)
//}

// GetGames ...
// @Title 获取结算部用户
// @Description 获取结算部用户
// @Success 200
// @Failure 403
// @router /getcpverifyuser [get]
func (c *CpVerifyAccountController) GetVerifyUsers() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CP_VERIFY_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	res, err := models.GetCpVerifyUsers()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(res)
}

//// GetGames ...
//// @Title 根据游戏获取对账的起止时间
//// @Description 根据游戏获取对账的起止时间
//// @Param	gameid		query	string	true	"游戏ID"
//// @Success 200
//// @Failure 403
//// @router /novdate [get]
//func (c *CpVerifyAccountController) GetVerifyDateByGameId() {
//	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_CP_VERIFY_ACCOUNT, nil)
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	gameid, err := c.GetInt("gameid")
//	if err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//	res, err := models.GetVerifyDateByGameId(gameid);
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	log.Printf("res; %v", res)
//	c.RespJSONData(res)
//}

//// @Title 添加对账单
//// @Description 添加对账单
//// @Param	body	body 	models.bean.InputCpVerifyAccount  true		"body for CpVerifyAccount content"
//// @Success 201 {int} models.CpVerifyAccount
//// @Failure 403 body is empty
//// @router / [post]
//func (c *CpVerifyAccountController) Post() {
//	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_CP_VERIFY_ACCOUNT, nil)
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	var v models.CpVerifyAccount
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//
//	// 检查状态
//	if v.Status != models.CP_VERIFY_S_VERIFYING && v.Status != models.CP_VERIFY_S_VERIFYIED && v.Status != models.CP_VERIFY_S_RECEIPT {
//		c.RespJSON(bean.CODE_Params_Err, "status must in [10,20,30]")
//		return
//	}
//
//	// 游戏格式是否合法
//	games := []models.GameAmount{}
//	err = json.Unmarshal([]byte(v.Games), &games)
//	if err != nil {
//		c.RespJSON(bean.CODE_Params_Err, "games is not a json string: "+err.Error())
//		return
//	}
//	if len(games) == 0 {
//		c.RespJSON(bean.CODE_Params_Err, "games can't be empty")
//		return
//	}
//
//	// 检查重复
//	if models.CheckCpOrderVerified(&v) {
//		c.RespJSON(bean.CODE_Bad_Request, "have exsit")
//		return
//	}
//	v.CreateUserId = c.Uid()
//	v.UpdateUserId = c.Uid()
//
//	// 添加
//	if _, err := models.AddCpVerifyAccount(&v); err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//
//	gameIds := []interface{}{}
//	for _, v := range games {
//		gameIds = append(gameIds, v.GameId)
//	}
//
//	// 标记为已读
//	affect, err := models.SetOrderCpVerified(gameIds, v.StartTime, v.EndTime, 1)
//	if err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//	if affect == 0 {
//		log.Error("CpVerify-gameIds", "add verifty success but affected order is zero", gameIds)
//	}
//
//	err = DoByCpVerifyStatus(v)
//	if err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//
//	c.RespJSONData(v)
//	return
//}
//
//// GetOne ...
//// @Title Get One
//// @Description get CpVerifyAccount by id
//// @Param	id		path 	string	true		"The key for staticblock"
//// @Success 200 {object} models.CpVerifyAccount
//// @Failure 403 :id is empty
//// @router /:id [get]
//func (c *CpVerifyAccountController) GetOne() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	v, err := models.GetCpVerifyAccountById(id)
//	if err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//	ss := []models.CpVerifyAccount{*v}
//	models.AddCompanyInfo4CpVerify(&ss)
//	models.AddUserInfo4CpVerify(&ss)
//	c.RespJSONData(ss[0])
//}
//
//// GetAll ...
//// @Title Get All
//// @Description 获取自己的对账单
//// @Param	company_ids		query	string	false	"发行商ID，多个用','隔开"
//// @Param	status	query	string	false	"对账单的状态"
//// @Param	start		query	string	false	"开始时间 2006-01-02 格式"
//// @Param	end		query	string	false	"结束时间 2006-01-02 格式"
//// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
//// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
//// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
//// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
//// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
//// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
//// @Success 200 {object} models.CpVerifyAccount
//// @Failure 403
//// @router / [get]
//func (c *CpVerifyAccountController) GetAll() {
//	total, ss, errCode, err := c.getAll()
//	if err != nil {
//		c.RespJSON(errCode, err.Error())
//		return
//	}
//	c.RespJSONDataWithTotal(ss, total)
//}
//
//var lockPutCpVerify sync.Mutex
//
//
//// Put ...
//// @Title Put
//// @Description 更新对账单
//// @Param	id		path 	string	true		"对账id"
//// @Param	body		body 	models.bean.InputCpVerifyAccount	true		"body for CpVerifyAccount content"
//// @Success 200 {object} models.CpVerifyAccount
//// @Failure 403 :id is not int
//// @router /:id [put]
//func (c *CpVerifyAccountController) Put() {
//	lockPutCpVerify.Lock()
//	defer lockPutCpVerify.Unlock()
//
//	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_CP_VERIFY_ACCOUNT, nil)
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	v := models.CpVerifyAccount{Id: id}
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//
//	// 游戏格式是否合法
//	games := []models.GameAmount{}
//	err = json.Unmarshal([]byte(v.Games), &games)
//	if err != nil {
//		c.RespJSON(bean.CODE_Params_Err, "games is not a json string: "+err.Error())
//		return
//	}
//	if len(games) == 0 {
//		c.RespJSON(bean.CODE_Params_Err, "games can't be empty")
//		return
//	}
//	newGameIds := []interface{}{}
//	for _, v := range games {
//		newGameIds = append(newGameIds, v.GameId)
//	}
//
//	v.UpdateTime = int(time.Now().Unix())
//	v.UpdateUserId = c.Uid()
//
//	// start
//	// 先取出老的对账单,回滚order表中标记的已对账流水
//	oldV, err := models.GetCpVerifyAccountById(id)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	// 回滚回款金额
//	if oldV.AmountSettle != 0 {
//		err = models.AddSettlePreAmount(oldV.CompanyId, oldV.AmountSettle)
//		if err != nil {
//			c.RespJSON(bean.CODE_Bad_Request, "回滚回款金额"+err.Error())
//			return
//		}
//	}
//
//	// 回滚order标记
//	oldGames := []models.GameAmount{}
//	err = json.Unmarshal([]byte(oldV.Games), &oldGames)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, "old games is not a json string: "+err.Error())
//		return
//	}
//	oldGameIds := []interface{}{}
//	for _, v := range oldGames {
//		oldGameIds = append(oldGameIds, v.GameId)
//	}
//
//	_, err = models.SetOrderCpVerified(oldGameIds, oldV.StartTime, oldV.EndTime, 2)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//
//	// 标记为已录入
//	_, err = models.SetOrderCpVerified(newGameIds, v.StartTime, v.EndTime, 1)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	// end
//
//
//	if err := models.UpdateCpVerifyAccount(&v); err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//
//	}
//	err = DoByCpVerifyStatus(v)
//	if err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//	c.RespJSONData(v)
//}
//
////根据状态TODO
//func DoByCpVerifyStatus(v models.CpVerifyAccount) (err error) {
//	switch v.Status {
//	case models.CP_VERIFY_S_VERIFYING:
//		return nil
//	case models.CP_VERIFY_S_VERIFYIED:
//		return nil
//	case models.CP_VERIFY_S_RECEIPT:
//		err = models.DoPushSettleToVerify(v.CompanyId)
//		return
//	default:
//		return errors.New("no status")
//	}
//}
//
//// Delete ...
//// @Title Delete
//// @Description delete the CpVerifyAccount
//// @Param	id		path 	string	true		"The id you want to delete"
//// @Success 200 {string} delete success!
//// @Failure 403 id is empty
//// @router /:id [delete]
//func (c *CpVerifyAccountController) Delete() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//
//	if err := models.DeleteCpVerifyAccount(id); err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//	}
//	c.RespJSONData("OK")
//
//}
//
//// @router /download [get]
//func (c *CpVerifyAccountController) DownLoad() {
//	c.Ctx.Input.SetParam("limit", "0")
//	c.Ctx.Input.SetParam("offset", "0")
//
//	_, ss, errCode, err := c.getAll()
//	if err != nil {
//		c.RespJSON(errCode, err.Error())
//		return
//	}
//
//	rs := make([]map[string]interface{}, len(ss))
//	for i, v := range ss {
//		r := map[string]interface{}{
//			"amount_my":       fmt.Sprintf("%.2f", v.AmountMy),
//			"amount_opposite": fmt.Sprintf("%.2f", v.AmountOpposite),
//			"amount_payable":  fmt.Sprintf("%.2f", v.AmountPayable),
//			"amount_settle":   fmt.Sprintf("%.2f", v.AmountSettle),
//		}
//		tStart := v.StartTime
//		tEnd := v.EndTime
//		tCreated := time.Unix(int64(v.CreateTime), 0).Format("2006-01-02 15:04:05")
//		r["time"] = tStart + " -- " + tEnd
//		r["create_time"] = tCreated
//		r["verify_time"] = time.Unix(int64(v.VerifyTime), 0).Format("2006-01-02")
//		r["update_time"] = time.Unix(int64(v.UpdateTime), 0).Format("2006-01-02 15:04:05")
//		if v.CreateUser != nil {
//			r["create_username"] = v.CreateUser.Nickname
//		}
//		if v.UpdateUser != nil {
//			r["update_username"] = v.UpdateUser.Nickname
//		}
//		if v.VerifyUser != nil {
//			r["verify_username"] = v.VerifyUser.Nickname
//		}
//		if v.Company != nil {
//			r["company"] = v.Company.Name
//		}
//		r["status"] = models.CpStatusCode2String[v.Status]
//		rs[i] = r
//	}
//
//	cols := []string{"time", "company", "status", "amount_my", "amount_opposite",
//	                 "amount_payable", "amount_settle", "verify_time", "verify_username", "create_username", "create_time",
//	                 "update_username", "update_time"}
//	maps := map[string]string{
//		cols[0]:  "账单日期",
//		cols[1]:  "发行商",
//		cols[2]:  "状态",
//		cols[3]:  "流水金额",
//		cols[4]:  "对方流水金额",
//		cols[5]:  "应付金额",
//		cols[6]:  "结算金额",
//		cols[7]:  "对账时间",
//		cols[8]:  "对账人",
//		cols[9]:  "创建人",
//		cols[10]: "创建时间",
//		cols[11]: "更新人",
//		cols[12]: "更新时间",
//	}
//	c.RespExcel(rs, "CP对账单", cols, maps)
//}
//
//func (c *CpVerifyAccountController) getAll() (total int64, ss []models.CpVerifyAccount, errCode int, err error) {
//	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CP_VERIFY_ACCOUNT, []string{"game_id"})
//	if err != nil {
//		errCode = bean.CODE_Forbidden
//		return
//	}
//
//	filter, err := tool.BuildFilter(c.Controller, 20)
//	if err != nil {
//		errCode = bean.CODE_Params_Err
//		return
//	}
//	idsStr := c.GetString("company_ids")
//	start := c.GetString("start", "2006-01-02")
//	end := c.GetString("end", time.Now().AddDate(1, 0, 0).Format("2006-01-02"))
//
//	status, err := c.GetInt("status", 0)
//	if err != nil {
//		c.RespJSON(bean.CODE_Not_Found, err.Error())
//		return
//	}
//	compIds := []interface{}{}
//	for _, value := range strings.Split(idsStr, ",") {
//		if value != "" {
//			compIds = append(compIds, value)
//		}
//	}
//	filter.Where = map[string][]interface{}{}
//	if len(compIds) != 0 {
//		filter.Where["company_id__in"] = compIds
//	}
//	if start != "" {
//		filter.Where["start_time__gte"] = []interface{}{start}
//	}
//	if end != "" {
//		filter.Where["end_time__lte"] = []interface{}{end}
//	}
//	if status != 0 {
//		filter.Where["status"] = []interface{}{status}
//	}
//
//	filter.Sortby = []string{"id"}
//	filter.Order = []string{"desc"}
//
//	tool.InjectPermissionWhere(where, &filter.Where)
//
//	ss = []models.CpVerifyAccount{}
//	total, err = tool.GetAllByFilterWithTotal(new(models.CpVerifyAccount), &ss, filter)
//	if err != nil {
//		errCode = bean.CODE_Bad_Request
//		return
//	}
//
//	models.AddCompanyInfo4CpVerify(&ss)
//	models.AddUserInfo4CpVerify(&ss)
//	return
//}
