// unused
// 已经用verify_channel 代替
package controllers

// 渠道对账
type ChannelVerifyAccountController struct {
	BaseController
}

// URLMapping ...
func (c *ChannelVerifyAccountController) URLMapping() {

}

//// @Title 获取渠道未对账的数据
//// @Description 获取渠道未对账的数据
//// @Param	games		query	string	false	"游戏ID，多个用","隔开“
//// @Param	channels	query	string	false	"渠道ID，多个用","隔开"
//// @Param	start		query	string	false	"开始时间 2006-01-02 格式"
//// @Param	end		query	string	false	"结束时间 2006-01-02 格式"
//// @Success 200 {object} models.ChannelVerifyAccount
//// @Failure 403
//// @router /pre [get]
//func (c *ChannelVerifyAccountController) GetNoVerifyAccount() {
//	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, nil)
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	gids_str := c.GetString("games")
//	channel_str := c.GetString("channels")
//
//	start := c.GetString("start", "2006-01-02")
//	end := c.GetString("end", time.Now().Format("2006-01-02"))
//	var gids, cps []interface{}
//	for _, value := range strings.Split(gids_str, ",") {
//		if value != "" {
//			gids = append(gids, value)
//		}
//	}
//	for _, value := range strings.Split(channel_str, ",") {
//		if value != "" {
//			cps = append(cps, value)
//		}
//	}
//	res, err := models.GetChanNoVerifyAccount(start, end, gids, cps)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, bean.CodeString(bean.CODE_Bad_Request))
//		return
//	}
//	c.RespJSONData(res)
//}

//// 获取游戏
//// @Title 获取游戏
//// @Description 获取游戏
//// @Success 200 {object} models.Game
//// @Failure 403
//// @router /games [get]
//func (c *ChannelVerifyAccountController) GetGames() {
//	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, []string{"game_id"})
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//
//	filter := &tool.Filter{
//		Where:  where,
//		Fields: []string{"GameId", "GameName"},
//	}
//
//	ss := []models.Game{}
//	err = tool.GetAllByFilter(new(models.Game), &ss, filter)
//
//	if err != nil {
//		c.RespJSON(bean.CODE_Not_Found, err.Error())
//		return
//	}
//
//	c.RespJSONData(ss)
//}

//// 获取所有渠道
//// @Title 获取所有渠道
//// @Description 获取所有渠道
//// @Success 200 {object} models.Channel
//// @Failure 403
//// @router /channels [get]
//func (c *ChannelVerifyAccountController) GetChannel() {
//	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, []string{"game_id"})
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//
//	filter := &tool.Filter{
//		Where:  where,
//		Fields: []string{"Channelid", "Cp", "Name"},
//	}
//
//	ss := []models.Channel{}
//	err = tool.GetAllByFilter(new(models.Channel), &ss, filter)
//
//	if err != nil {
//		c.RespJSON(bean.CODE_Not_Found, err.Error())
//		return
//	}
//
//	c.RespJSONData(ss)
//}

//// @Title 获取未对账的渠道
//// @Description 获取未对账的渠道
//// @Success 200 {object} models.Channel
//// @Failure 403
//// @router /novchans [get]
//func (c *ChannelVerifyAccountController) GetChannels() {
//	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, []string{"game_id"})
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	res, err := models.GetNoVerifyChannels(where)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	c.RespJSONData(res)
//}
//
//// @Title 通过渠道码获取该渠道对账的起止时间
//// @Description 通过渠道码获取该渠道对账的起止时间
//// @Param	cp	query	string	true	"渠道码"
//// @Success 200
//// @Failure 403
//// @router /novdate [get]
//func (c *ChannelVerifyAccountController) GetNotVerifyDate() {
//	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, nil)
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	cp := c.GetString("cp")
//	res, err := models.GetVerifyDateByCp(cp)
//	if err != nil {
//		c.RespJSON(bean.CODE_Not_Found, err.Error())
//		return
//	}
//	c.RespJSONData(res)
//}
//
//// @Title 通过渠道码和日期获取所有未对账的游戏
//// @Description 通过渠道码和日期获取所有未对账的游戏
//// @Param	cp	query	string	true	"渠道码"
//// @Param	start		query	string	true	"开始时间 2006-01-02 格式"
//// @Param	end		query	string	true	"结束时间 2006-01-02 格式"
//// @Success 200 {object} models.GameAmount
//// @Failure 403
//// @router /novgames [get]
//func (c *ChannelVerifyAccountController) GetGamesByCpAndDate() {
//	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, nil)
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	cp := c.GetString("cp")
//	start := c.GetString("start")
//	end := c.GetString("end")
//	bodyMy, _ := c.GetInt("body_my")
//	res, err := models.GetGamesByCpAndDate(cp, start, end, bodyMy, where)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	c.RespJSONData(res)
//}

//// @Title 获取某渠道回款主体
//// @Description 通过渠道码和日期获取所有未对账的游戏
//// @Param	cp	query	string	true	"渠道码"
//// @Param	start		query	string	true	"开始时间 2006-01-02 格式"
//// @Param	end		query	string	true	"结束时间 2006-01-02 格式"
//// @Success 200 {object} models.GameAmount
//// @Failure 403
//// @router /remitcompanies [get]
//func (c *ChannelVerifyAccountController) GetRemitCompanies() {
//	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, nil)
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	channelCode := c.GetString("channel_code")
//
//	res, err := models.GetRemitCompaniesByChannel(channelCode)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	c.RespJSONData(res)
//}

//// @Title Post
//// @Description create ChannelVerifyAccount
//// @Param	body		body 	models.ChannelVerifyAccount	true		"body for ChannelVerifyAccount content"
//// @Success 201 {int} models.ChannelVerifyAccount
//// @Failure 403 body is empty
//// @router / [post]
//func (c *ChannelVerifyAccountController) Post() {
//	var v models.ChannelVerifyAccount
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	if err := c.Validate(&v); err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//	if v.Status != models.CHAN_VERIFY_S_VERIFYING && v.Status != models.CHAN_VERIFY_S_VERIFYIED && v.Status != models.CHAN_VERIFY_S_RECEIPT {
//		c.RespJSON(bean.CODE_Params_Err, "status must in [10,20,30]")
//		return
//	}
//
//	// 已开发票,检查回款主体是否存在
//	if v.Status == models.CHAN_VERIFY_S_RECEIPT {
//		if v.RemitCompanyId == 0 {
//			c.RespJSON(bean.CODE_Params_Err, "remit_company_id can not be empty")
//			return
//		}
//	}
//
//	var games []models.GameAmount
//	// 应该有这几个字段 game_id,game_name,amount_payable,amount_opposite
//	err := json.Unmarshal([]byte(v.GameStr), &games)
//	if err != nil {
//		c.RespJSON(bean.CODE_Params_Err, "games params error")
//		return
//	}
//
//	var gameIds []interface{}
//	for _, value := range games {
//		gameIds = append(gameIds, value.GameId)
//	}
//	// 判断重复
//	if models.CheckChanOrderVerified(&v) {
//		c.RespJSON(bean.CODE_Bad_Request, "have exsit")
//		return
//	}
//
//	v.CreateUserId = c.Uid()
//	v.UpdateUserId = c.Uid()
//	if _, err := models.AddChannelVerifyAccount(&v); err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//
//	affect, err := models.SetOrderChanVerified(gameIds, v.Cp, v.StartTime, v.EndTime, 1)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//
//	if affect == 0 {
//		log.Error("ChannelVerify-gameIds", "add verifty success but affected order is zero", gameIds)
//	}
//
//	err = DoByChanVerifyStatus(v)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	c.RespJSONData(v)
//	return
//}
//
//// @Title Get One
//// @Description get ChannelVerifyAccount by id
//// @Param	id		path 	string	true		"The key for staticblock"
//// @Success 200 {object} models.ChannelVerifyAccount
//// @Failure 403 :id is empty
//// @router /:id [get]
//func (c *ChannelVerifyAccountController) GetOne() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	v, err := models.GetChannelVerifyAccountById(id)
//	if err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//	ss := []models.ChannelVerifyAccount{*v}
//
//	//models.ChanVerifyAddGameInfo(&ss)
//	models.AddUserInfo4ChanVerify(&ss)
//	models.AddCpInfo4ChanVerify(&ss)
//	models.AddRemitCompany4ChanVerify(&ss)
//	c.RespJSONData(ss[0])
//}
//
//// GetAll ...
//// @Title 获取自己所有的订单
//// @Description 获取自己所有的订单
//// @Param	games		query	string	false	"游戏ID，多个用","隔开“
//// @Param	channels	query	string	false	"渠道ID，多个用","隔开"
//// @Param	status	query	string	false	"对账单的状态"
//// @Param	start		query	string	false	"开始时间 2006-01-02 格式"
//// @Param	end		query	string	false	"结束时间 2006-01-02 格式"
//// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
//// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
//// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
//// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
//// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
//// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
//// @Success 200 {object} models.ChannelVerifyAccount
//// @Failure 403
//// @router / [get]
//func (c *ChannelVerifyAccountController) GetAll() {
//	total, ss, errCode, err := c.getAll()
//	if err != nil {
//		c.RespJSON(errCode, err.Error())
//		return
//	}
//	c.RespJSONDataWithTotal(ss, total)
//}
//
//var lockPutChannelVerify sync.Mutex
//
//// @Title 更新对账单
//// @Description 更新对账单
//// @Param	id		path 	string	true		"The id you want to update"
//// @Param	body		body 	models.ChannelVerifyAccount	true		"body for ChannelVerifyAccount content"
//// @Success 200 {object} models.ChannelVerifyAccount
//// @Failure 403 :id is not int
//// @router /:id [put]
//func (c *ChannelVerifyAccountController) Put() {
//	lockPutChannelVerify.Lock()
//	defer lockPutChannelVerify.Unlock()
//
//	_, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_CP_VERIFY_ACCOUNT, nil)
//	if err != nil {
//		c.RespJSON(bean.CODE_Forbidden, err.Error())
//		return
//	}
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	v := models.ChannelVerifyAccount{Id: id}
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//
//	if err := c.Validate(&v); err != nil {
//		c.RespJSON(bean.CODE_Params_Err, err.Error())
//		return
//	}
//
//	// 已开发票,检查回款主体是否存在
//	if v.Status == models.CHAN_VERIFY_S_RECEIPT {
//		if v.RemitCompanyId == 0 {
//			c.RespJSON(bean.CODE_Params_Err, "remit_company_id can not be empty")
//			return
//		}
//	}
//
//	// 游戏格式是否合法
//	games := []models.GameAmount{}
//	err = json.Unmarshal([]byte(v.GameStr), &games)
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
//	v.UpdateUserId = c.Uid()
//
//	// start
//	// 先取出老的对账单,回滚order表中标记的已对账流水
//	oldV, err := models.GetChannelVerifyAccountById(id)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	// 回滚回款金额
//	if oldV.AmountRemit != 0 {
//		err = models.AddRemitPreAmount(oldV.RemitCompanyId, oldV.AmountRemit)
//		if err != nil {
//			c.RespJSON(bean.CODE_Bad_Request, "回滚回款金额"+err.Error())
//			return
//		}
//	}
//
//	// 回滚order标记
//	oldGames := []models.GameAmount{}
//	err = json.Unmarshal([]byte(oldV.GameStr), &oldGames)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, "old games is not a json string: "+err.Error())
//		return
//	}
//	oldGameIds := []interface{}{}
//	for _, v := range oldGames {
//		oldGameIds = append(oldGameIds, v.GameId)
//	}
//	_, err = models.SetOrderChanVerified(oldGameIds, oldV.Cp, oldV.StartTime, oldV.EndTime, 2)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//
//	// 标记为已录入
//	_, err = models.SetOrderChanVerified(newGameIds, v.Cp, v.StartTime, v.EndTime, 1)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	// end
//
//	if err := models.UpdateChanVerifyAccount(&v); err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//
//	err = DoByChanVerifyStatus(v)
//	if err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//		return
//	}
//	c.RespJSONData(v)
//}
//
////根据状态TODO
//func DoByChanVerifyStatus(v models.ChannelVerifyAccount) (err error) {
//	switch v.Status {
//	case models.CHAN_VERIFY_S_VERIFYING:
//		return nil
//
//	case models.CHAN_VERIFY_S_VERIFYIED:
//		return nil
//
//	case models.CHAN_VERIFY_S_RECEIPT:
//		err = models.DoPushRemitToVerify(v.RemitCompanyId)
//		return
//	}
//
//	return
//}

////检查修改的状态，如果为最后一个状态，则不能回滚
//func checkChanVerifyStatus(v models.ChannelVerifyAccount) (err error) {
//	if v.Status != models.CHAN_VERIFY_S_RECEIPT {
//		r, err := models.GetCpVerifyAccountById(v.Id)
//		if err != nil {
//			return err
//		}
//		if r.Status == models.CHAN_VERIFY_S_RECEIPT {
//			return errors.New("can't modify")
//		}
//	} else {
//		err = errors.New("can't change status when it is finished")
//	}
//	return nil
//}

//// Delete ...
//// @Title Delete
//// @Description delete the ChannelVerifyAccount
//// @Param	id		path 	string	true		"The id you want to delete"
//// @Success 200 {string} delete success!
//// @Failure 403 id is empty
//// @router /:id [delete]
//func (c *ChannelVerifyAccountController) Delete() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//
//	if err := models.DeleteChannelVerifyAccount(id); err != nil {
//		c.RespJSON(bean.CODE_Bad_Request, err.Error())
//	}
//	c.RespJSONData("OK")
//}
//
//func (c *ChannelVerifyAccountController) getAll() (total int64, ss []models.ChannelVerifyAccount, errCode int, err error) {
//	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, []string{"game_id"})
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
//	gids_str := c.GetString("games")
//	channel_str := c.GetString("channels")
//
//	start := c.GetString("start", "2006-01-02")
//	end := c.GetString("end", time.Now().AddDate(1, 0, 0).Format("2006-01-02"))
//
//	status, err := c.GetInt("status", 0)
//	if err != nil {
//		c.RespJSON(bean.CODE_Not_Found, err.Error())
//		return
//	}
//	gids := []interface{}{}
//	cps := []interface{}{}
//
//	for _, value := range strings.Split(gids_str, ",") {
//		if value != "" {
//			gids = append(gids, value)
//		}
//	}
//	for _, value := range strings.Split(channel_str, ",") {
//		if value != "" {
//			cps = append(cps, value)
//		}
//	}
//	filter.Where = map[string][]interface{}{}
//	if len(cps) != 0 {
//		filter.Where["cp__in"] = cps
//	}
//	if len(gids) != 0 {
//		filter.Where["game_id__in"] = gids
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
//	filter.Sortby = []string{"id"}
//	filter.Order = []string{"desc"}
//
//	tool.InjectPermissionWhere(where, &filter.Where)
//
//	ss = []models.ChannelVerifyAccount{}
//	total, err = tool.GetAllByFilterWithTotal(new(models.ChannelVerifyAccount), &ss, filter)
//	if err != nil {
//		errCode = bean.CODE_Bad_Request
//		return
//	}
//
//	//models.ChanVerifyAddGameInfo(&ss)
//	models.AddUserInfo4ChanVerify(&ss)
//	models.AddCpInfo4ChanVerify(&ss)
//	models.AddRemitCompany4ChanVerify(&ss)
//	return
//}
//
//// @router /download [get]
//func (c *ChannelVerifyAccountController) DownLoad() {
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
//			"amount_remit":    fmt.Sprintf("%.2f", v.AmountRemit),
//		}
//		tStart := v.StartTime
//		tEnd := v.EndTime
//		tCreated := time.Unix(int64(v.CreateTime), 0).Format("2006-01-02 15:04:05")
//		r["time"] = tStart + " -- " + tEnd
//		r["create_time"] = tCreated
//		r["verify_time"] = time.Unix(int64(v.VerifyTime), 0).Format("2006-01-02 15:04:05")
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
//		r["game_num"] = v.GameNum
//		r["cp"] = v.Channel.Name
//		r["games"] = v.GameStr
//		//if v.GameStr != "" {
//		//	temp := ""
//		//	for _, value := range v.Games {
//		//		temp += value.GameName + ";"
//		//		r["games"] = temp
//		//	}
//		//}
//		r["status"] = models.ChanStatusCode2String[v.Status]
//		rs[i] = r
//	}
//
//	cols := []string{"time", "cp", "games", "game_num", "status", "amount_my", "amount_opposite",
//	                 "amount_payable", "amount_remit", "verify_username", "create_username", "create_time",
//	                 "update_username", "update_time"}
//	maps := map[string]string{
//		cols[0]:  "账单日期",
//		cols[1]:  "渠道",
//		cols[2]:  "游戏",
//		cols[3]:  "游戏数量",
//		cols[4]:  "状态",
//		cols[5]:  "我方流水",
//		cols[6]:  "对方流水",
//		cols[7]:  "应付金额",
//		cols[8]:  "回款金额",
//		cols[9]:  "对账人",
//		cols[10]: "创建人",
//		cols[11]: "创建时间",
//		cols[12]: "更新人",
//		cols[13]: "更新时间",
//	}
//	c.RespExcel(rs, "渠道对账单", cols, maps)
//}
