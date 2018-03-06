package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"kuaifa.com/kuaifa/work-together/tool/old_sys"
	"strconv"
	"strings"
	"fmt"
	"time"
)

// 合同
type ContractController struct {
	BaseController
}

// URLMapping ...
func (c *ContractController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Contract
// @Param	body		body 	models.Contract	true		"body for Contract content"
// @Success 201 {int} models.Contract
// @Failure 403 body is empty
// @router / [post]
//func (c *ContractController) Post() {
//	var v models.Contract
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
//		if _, err := models.AddContract(&v); err == nil {
//			c.Ctx.Output.SetStatus(201)
//			c.Data["json"] = v
//		} else {
//			c.Data["json"] = err.Error()
//		}
//	} else {
//		c.Data["json"] = err.Error()
//	}
//	c.ServeJSON()
//}

// GetOne ...
// @Title Get One
// @Description get Contract by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Contract
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ContractController) GetOne() {
	flag := c.GetString("flag")
	var PMSM string
	if flag == "cp" { //cp合同
		PMSM = bean.PMSM_CONTRACT_CP
	}else if flag == "qd" { //渠道合同
		PMSM = bean.PMSM_CONTRACT_CHANNEL
	}else{
		c.RespJSON(bean.CODE_Bad_Request, "参数缺少")
		return
	}
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, PMSM, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetContractById(id, flag)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get Contract
// @Param	gameids	query	string	false	"游戏id"
// @Param	channelids	query	string	false	"渠道id"
// @Param	status	query	string	false	"合同状态"
// @Param	flag	query	string	false	"company type. e.g. cp/channel"
// @Success 200 {object} models.Contract
// @Failure 403
// @router / [get]
func (c *ContractController) GetAll() {

	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss,total)
}

func (c *ContractController) getAll() (total int64, ss []models.Contract, errCode int, err error)  {
	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		errCode = bean.CODE_Params_Err
		return
	}

	idStr := c.GetString("gameids")
	channels := c.GetString("channels")
	status := c.GetString("status")
	body, serr2 := c.GetInt("body")
	flag := c.GetString("flag")

	ids := []interface{}{}
	for _, v := range strings.Split(idStr, ",") {
		if v != "" {
			ids = append(ids, v)
		}
	}

	channel := []interface{}{}
	for _, v := range strings.Split(channels, ",") {
		if v != "" {
			channel = append(channel, v)
		}
	}

	statuss := []interface{}{}
	for _, v := range strings.Split(status, ",") {
		if v != "" {
			statuss = append(statuss, v)
		}
	}

	filter.Order = []string{"desc"}
	filter.Sortby = []string{"id"}
	filter.Where = map[string][]interface{}{}
	if len(ids) != 0 {
		filter.Where["game_id__in"] = ids
	}

	if len(channel) != 0 {
		filter.Where["channel_code__in"] = channel
	}

	if len(statuss) != 0 {
		filter.Where["state__in"] = statuss
	}
	if serr2 == nil {
		filter.Where["body_my__exact"] = []interface{}{body}
	}
	var permission = bean.PMSM_CONTRACT_CP
	if flag != "cp" {
		filter.Where["company_type__exact"] = []interface{}{1}
		permission = bean.PMSM_CONTRACT_CHANNEL
	} else {
		filter.Where["company_type__exact"] = []interface{}{0}
	}
	filter.Where["effective_state__exact"] = []interface{}{1}//只获取当前有效的合同
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, permission, nil)
	if err != nil {
		errCode = bean.CODE_Forbidden
		return
	}

	tool.InjectPermissionWhere(where, &filter.Where)
	ss = []models.Contract{}

	total, err = tool.GetAllByFilterWithTotal(new(models.Contract), &ss, filter)
	if err != nil {
		errCode = bean.CODE_Not_Found
		return
	}

	models.GroupRemitAddGameInfo(&ss)
	if flag == "cp" {
		models.GroupRemitAddCpInfo(&ss)
	} else {
		models.AddChannelCompanyInfo(&ss)
		models.GroupRemitContractAddChannelInfo(&ss)
		models.AddBusinessInfo(&ss)
	}
	models.GroupRemitAddContractUserInfo(&ss)
	models.GroupRemitAddContractStatusInfo(&ss)
	models.ParseLadder2Json(&ss)

	return
}

// @router /cpDownload [get]
func (c *ContractController) CPDownLoad() {
	c.Ctx.Input.SetParam("limit","0")
	c.Ctx.Input.SetParam("offset","0")

	_, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}

	//rs := make([]map[string]interface{}, len(ss))
	rs := []map[string]interface{}{}
	for _, v := range ss {
		ladders := []old_sys.Ladder4Post{}
		json.Unmarshal([]byte(v.Ladder), &ladders)

		r := map[string]interface{}{}

		if v.BodyMy == 1 {
			r["bodyMy"] = "云端"
		}else if v.BodyMy == 2 {
			r["bodyMy"] = "有量"
		}else{
			r["bodyMy"] = "无主体"
		}
		r["issue"] = v.ChannelCompanyName
		r["game"] = v.Game.GameName
		if v.BeginTime != 0 {
			r["signTime"] = time.Unix(v.BeginTime,0).Format("2006-01-02")
		}
		if v.EndTime != 0 {
			r["endTime"] = time.Unix(v.EndTime,0).Format("2006-01-02")
		}
		r["state"] = v.Status.Name
		r["desc"] = v.Desc

		if len(ladders) == 0 {
			rs = append(rs, r)
		} else if len(ladders) == 1 {
			rulesString := strings.Split(ladders[0].Rule, "&")
			if len(rulesString) > 1 {
				for _,rule := range rulesString {
					ru := rule
					result := strings.Split(ru, "<")
					if result[1] == "money" {
						r["money"] = result[0] + "~" + result[2]
					}
					if result[1] == "time" {
						r["time"] = result[0] + "~" + result[2]
					}
					if result[1] == "user" {
						r["user"] = result[0] + "~" + result[2]
					}
				}
			}
			r["ratio"] = ladders[0].Ratio
			r["slotting_fee"] = ladders[0].SlottingFee
			rs = append(rs, r)
		} else {
			for _, ladder := range ladders{
				v := map[string]interface{}{}
				v["bodyMy"] = r["bodyMy"]
				v["issue"] = r["issue"]
				v["game"] = r["game"]
				v["signTime"] = r["signTime"]
				v["endTime"] = r["endTime"]
				v["state"] = r["state"]
				v["desc"] = r["desc"]

				rulesString := strings.Split(ladder.Rule, "&")
				if len(rulesString) > 1 {
					for _, rule := range rulesString {
						result := strings.Split(rule, "<")
						if result[1] == "money" {
							v["money"] = result[0] + "~" + result[2]
						}
						if result[1] == "time" {
							v["time"] = result[0] + "~" + result[2]
						}
						if result[1] == "user" {
							v["user"] = result[0] + "~" + result[2]
						}
					}
				}
				v["ratio"] = ladder.Ratio
				v["slotting_fee"] = ladder.SlottingFee
				rs = append(rs, v)
			}
		}
		//rs[i] = r
	}

	cols := []string{"bodyMy", "issue", "game", "signTime","endTime","state", "money", "time", "user", "ratio", "slotting_fee", "desc"}
	maps := map[string]string{
		cols[0]:"我方主体",
		cols[1]:"发行商",
		cols[2]:"游戏名称",
		cols[3]:"签订时间",
		cols[4]:"终止时间",
		cols[5]:"合同状态",
		cols[6]:"金额",
		cols[7]:"时间",
		cols[8]:"用户",
		cols[9]:"我方比例",
		cols[10]:"通道费",
		cols[11]:"备注",
	}
	tmpFileName := fmt.Sprintf("CP合同-%s", time.Now().Format("20060102150405"))
	c.RespExcel(rs, tmpFileName, cols, maps)
}

// @router /channelDownload [get]
func (c *ContractController) ChannelDownLoad() {
	c.Ctx.Input.SetParam("limit","0")
	c.Ctx.Input.SetParam("offset","0")

	_, ss, errCode, err := c.getAllChannelContract()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}

	//rs := make([]map[string]interface{}, len(ss))
	rs := []map[string]interface{}{}
	for _, v := range ss {
		ladders := []old_sys.Ladder4Post{}
		json.Unmarshal([]byte(v.Ladder), &ladders)

		r := map[string]interface{}{}

		if v.BodyMy == 1 {
			r["bodyMy"] = "云端"
		}else if v.BodyMy == 2 {
			r["bodyMy"] = "有量"
		}else{
			r["bodyMy"] = "无主体"
		}
		if v.Channel != nil {
			r["channelName"] = v.Channel.Name
		}
		r["companyName"] = v.ChannelCompanyName
		r["game"] = v.Game.GameName
		if v.BeginTime != 0 {
			r["signTime"] = time.Unix(v.BeginTime,0).Format("2006-01-02")
		}
		if v.EndTime != 0 {
			r["endTime"] = time.Unix(v.EndTime,0).Format("2006-01-02")
		}
		r["state"] = v.Status.Name
		r["desc"] = v.Desc
		r["business"] = v.Business

		if len(ladders) == 0 {
			rs = append(rs, r)
		} else if len(ladders) == 1 {
			rulesString := strings.Split(ladders[0].Rule, "&")
			if len(rulesString) > 1 {
				for _,rule := range rulesString {
					ru := rule
					result := strings.Split(ru, "<")
					if result[1] == "money" {
						r["money"] = result[0] + "~" + result[2]
					}
					if result[1] == "time" {
						r["time"] = result[0] + "~" + result[2]
					}
					if result[1] == "user" {
						r["user"] = result[0] + "~" + result[2]
					}
				}
			}
			r["ratio"] = ladders[0].Ratio
			r["slotting_fee"] = ladders[0].SlottingFee
			rs = append(rs, r)
		} else {
			for _, ladder := range ladders{
				v := map[string]interface{}{}
				v["bodyMy"] = r["bodyMy"]
				v["channelName"] = r["channelName"]
				v["companyName"] = r["companyName"]
				v["issue"] = r["issue"]
				v["game"] = r["game"]
				v["signTime"] = r["signTime"]
				v["endTime"] = r["endTime"]
				v["state"] = r["state"]
				v["desc"] = r["desc"]

				rulesString := strings.Split(ladder.Rule, "&")
				if len(rulesString) > 1 {
					for _, rule := range rulesString {
						result := strings.Split(rule, "<")
						if result[1] == "money" {
							v["money"] = result[0] + "~" + result[2]
						}
						if result[1] == "time" {
							v["time"] = result[0] + "~" + result[2]
						}
						if result[1] == "user" {
							v["user"] = result[0] + "~" + result[2]
						}
					}
				}
				v["ratio"] = ladder.Ratio
				v["slotting_fee"] = ladder.SlottingFee
				rs = append(rs, v)
			}
		}
		//rs[i] = r
	}

	cols := []string{"bodyMy", "channelName", "companyName", "game", "signTime","endTime","state", "money", "time", "user", "ratio", "slotting_fee", "business", "desc"}
	maps := map[string]string{
		cols[0]:"我方主体",
		cols[1]:"渠道名称",
		cols[2]:"公司名称",
		cols[3]:"游戏名称",
		cols[4]:"签订时间",
		cols[5]:"终止时间",
		cols[6]:"合同状态",
		cols[7]:"金额",
		cols[8]:"时间",
		cols[9]:"用户",
		cols[10]:"我方比例",
		cols[11]:"通道费",
		cols[12]:"商务负责人",
		cols[13]:"备注",
	}
	tmpFileName := fmt.Sprintf("渠道合同-%s", time.Now().Format("20060102150405"))
	c.RespExcel(rs, tmpFileName, cols, maps)
}

// GetAllChannelContract ...
// @Title 获取渠道合同
// @Description get Contract
// @Param	gameids	query	string	false	"游戏id"
// @Param	channelcodes	query	string	false	"渠道code"
// @Param	status	query	string	false	"合同状态"
// @Param	business	query	int	false	"商务负责人id"
// @Success 200 {object} models.Contract
// @Failure 403
// @router /channelContract/ [get]
func (c *ContractController) GetAllChannelContract() {

	total, ss, errCode, err := c.getAllChannelContract()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)

}


func (c *ContractController) getAllChannelContract() (total int64, ss []models.Contract, errCode int, err error)  {
	_, err = models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CONTRACT_CHANNEL, nil)
	if err != nil {
		errCode = bean.CODE_Forbidden
		return
	}

	idStr := c.GetString("gameids")
	channels := c.GetString("channelcodes")
	business, businessErr := c.GetInt("business")
	status := c.GetString("status")
	body, bodyErr := c.GetInt("body")

	var limit int = 15
	var offset int = 0
	if li, limitErr := c.GetInt("limit"); limitErr == nil{
		limit = li
	}
	if off, offsetErr := c.GetInt("offset"); offsetErr == nil{
		offset = off
	}

	qb, _ := orm.NewQueryBuilder("mysql")
	qb = qb.Select("contract.*").From("contract")

	if businessErr == nil{
		qb = qb.LeftJoin("channel_access").On("contract.game_id=channel_access.game_id AND " +
			"contract.channel_code=channel_access.channel_code").Where("channel_access.business_person = ?")
	}

	ids := []string{}
	for _, v := range strings.Split(idStr, ",") {
		if v != "" {
			ids = append(ids, v)
		}
	}
	if len(ids) != 0 {
		if strings.Contains(qb.String(), "WHERE"){
			qb = qb.And("contract.game_id").In(ids...)
		} else {
			qb = qb.Where("contract.game_id").In(ids...)
		}
	}

	channel := []string{}
	for _, v := range strings.Split(channels, ",") {
		if v != "" {
			channel = append(channel, "'" + v + "'")
		}
	}
	if len(channel) != 0 {
		if strings.Contains(qb.String(), "WHERE") {
			qb = qb.And("contract.channel_code").In(channel...)
		} else {
			qb = qb.Where("contract.channel_code").In(channel...)
		}
	}

	statuss := []string{}
	for _, v := range strings.Split(status, ",") {
		if v != "" {
			statuss = append(statuss, v)
		}
	}
	if len(statuss) != 0 {
		if strings.Contains(qb.String(), "WHERE") {
			qb = qb.And("contract.state").In(statuss...)
		} else {
			qb = qb.Where("contract.state").In(statuss...)
		}
	}

	if bodyErr == nil {
		if strings.Contains(qb.String(), "WHERE") {
			qb = qb.And("contract.body_my = ?")
		} else {
			qb = qb.Where("contract.body_my = ?")
		}
	}

	if strings.Contains(qb.String(), "WHERE") {
		qb = qb.And("contract.company_type = 1 and contract.effective_state = 1")
	} else {
		qb = qb.Where("contract.company_type = 1 and contract.effective_state = 1")
	}

	qb = qb.OrderBy("contract.id").Desc()
	ss = []models.Contract{}
	sql_total := qb.String()
	o := orm.NewOrm()

	if businessErr != nil && bodyErr != nil{
		total, _ = o.Raw(sql_total).QueryRows(&ss)
	} else if businessErr == nil && bodyErr != nil{
		total, _ = o.Raw(sql_total, business).QueryRows(&ss)
	} else if businessErr != nil && bodyErr == nil{
		total, _ = o.Raw(sql_total, body).QueryRows(&ss)
	} else if businessErr == nil && bodyErr == nil{
		total, _ = o.Raw(sql_total, business, body).QueryRows(&ss)
	}
	qb = qb.Limit(limit).Offset(offset)
	sql := qb.String()

	if businessErr != nil && bodyErr != nil{
		_, _ = o.Raw(sql).QueryRows(&ss)
	} else if businessErr == nil && bodyErr != nil{
		_, _ = o.Raw(sql, business).QueryRows(&ss)
	} else if businessErr != nil && bodyErr == nil{
		_, _ = o.Raw(sql, body).QueryRows(&ss)
	} else if businessErr == nil && bodyErr == nil{
		_, _ = o.Raw(sql, business, body).QueryRows(&ss)
	}

	models.GroupRemitAddGameInfo(&ss)
	models.AddChannelCompanyInfo(&ss)
	models.GroupRemitContractAddChannelInfo(&ss)
	models.AddBusinessInfo(&ss)
	models.GroupRemitAddContractUserInfo(&ss)
	models.GroupRemitAddContractStatusInfo(&ss)
	models.ParseLadder2Json(&ss)

	return
}

// GetAllEditIds ...
// @Title 获取游戏所有合同id
// @Description get Contract
// @Param	editId		path 	string	true		""
// @router /getAllEditIds/:editId [get]
func (c *ContractController) GetAllEditIds(){
	idStr := c.Ctx.Input.Param(":editId")
	id, err := strconv.Atoi(idStr)

	if err != nil{
		c.RespJSON(bean.CODE_Bad_Request, "参数缺少")
		return
	}

	ids, err := models.GetAllEditIds(id)

	if err != nil{
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(ids)
}

// GetAll ...
// @Title Get All Cp contract
// @Description get Contract
// @Param	gameids	query	string	false	"游戏id"
// @Param	status	query	string	false	"合同状态"
// @Param	bodyMy	query	string	false	"我方主体"
// @Success 200 {object} models.Contract
// @Failure 403
// @router /cp/ [get]
func (c *ContractController) GetCpContract() {
	where, err2 := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CONTRACT_CHANNEL, nil)
	if err2 != nil {
		c.RespJSON(bean.CODE_Forbidden, err2.Error())
		return
	}

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	idStr := c.GetString("gameids")
	status, _ := c.GetInt("status")
	body, _ := c.GetInt("bodyMy")

	ids := []interface{}{}
	for _, v := range strings.Split(idStr, ",") {
		if v != "" {
			ids = append(ids, v)
		}
	}

	if len(ids) != 0 {
		filter.Where["game_id__in"] = ids
	}

	filter.Where["state__exact"] = []interface{}{status}
	filter.Where["body_my__exact"] = []interface{}{body}
	filter.Where["company_type__exact"] = []interface{}{0}

	tool.InjectPermissionWhere(where, &filter.Where)
	//ss := []models.Contract{}

	ss, total, err := models.GetAllCpContract(filter.Where, filter.Limit, filter.Offset)

	//total, err := tool.GetAllByFilterWithTotal(new(models.Contract), &ss, filter)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	//models.GroupRemitAddGameInfo(&ss)
	//models.GroupRemitAddCpInfo(&ss)
	//models.GroupRemitAddContractUserInfo(&ss)
	//models.GroupRemitAddContractStatusInfo(&ss)
	//models.ParseLadder2Json(&ss)
	//models.AddCpBodyMyInfo(&ss)


	c.RespJSONDataWithTotal(ss, total)

}

// Put ...
// @Title Put
// @Description update the Contract
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Contract	true		"body for Contract content"
// @Success 200 {object} models.Contract
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ContractController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Contract{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	//参数判断
	if v.Id == 0 || v.State == 0 {
		c.RespJSON(bean.CODE_Params_Err, "参数缺少")
		return
	}
	if v.State!=149 && v.State!=154 && v.State!=155 {
		if v.BeginTime == 0 || v.EndTime == 0 {
			c.RespJSON(bean.CODE_Params_Err, "参数缺少")
			return
		}
	}

	ls := []old_sys.Ladder4Post{}
	if err := json.Unmarshal([]byte(v.Ladders), &ls); v.State>0 && err != nil {
		c.RespJSON(bean.CODE_Params_Err, "ladder format error, "+err.Error())
		return
	}
	v.Ladder = v.Ladders
	channelCode := ""

	// 判断是否是渠道合同(渠道合同说明是设置 渠道的阶梯)
	s := models.Contract{Id: id}
	if err := orm.NewOrm().Read(&s); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	if s.CompanyType == 1 {
		channel, err := models.GetChannelByChannelCode(s.ChannelCode)
		if err != nil {
			c.RespJSON(bean.CODE_Bad_Request, err.Error())
			return
		}
		channelCode = channel.Cp
	}
	//如果合同状态大于0 才同步老系统
	if v.State>0 {
		if err := old_sys.UpdateOrAddLadderList(s.GameId, channelCode, ls); err != nil {
			c.RespJSON(bean.CODE_Bad_Request, "ladder save error, "+err.Error())
			return
		}
	}

	if v.State == 152 || v.State == 150 { //如果状态修改为已寄出，未回寄或已签订，则更新合同统计表
		if err := models.UpdateContractStatistics(c.Uid(), v.State, v.CompanyType); err != nil{
			c.RespJSON(bean.CODE_Bad_Request, err.Error())
			return
		}
	}
	v.UpdatePerson = c.Uid()
	if err := models.UpdateContractById(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	express := v.Express
	express.ContractId = v.Id
	if err := models.UpdateExpressById(express); err != nil{
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData("OK")
}

// Renew ...
// @Title 合同续签
// @Description 根据给定的老合同id，将该合同的老合同有效状态都改为无效，并新增一条有效合同,返回该合同的id
// @Param	editId		path 	string	true		"The id you want to update"
// @router /renew/:editId [put]
func (c *ContractController) Renew() {
	idStr := c.Ctx.Input.Param(":editId")
	id, err := strconv.Atoi(idStr)

	if err != nil{
		c.RespJSON(bean.CODE_Bad_Request, "参数缺少")
		return
	}


	editId, err := models.Renew(id, c.Uid())
	if err != nil{
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(editId)
}

// Delete ...
// @Title Delete
// @Description delete the Contract
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
//func (c *ContractController) Delete() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	if err := models.DeleteContract(id); err == nil {
//		c.Data["json"] = "OK"
//	} else {
//		c.Data["json"] = err.Error()
//	}
//	c.ServeJSON()
//}

// 获取游戏
// @Title 获取渠道
// @Description 获取渠道
// @Success 200 {object} models.Channel
// @Failure 403
// @router /channel [get]
func (c *ContractController) GetChannel() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CONTRACT_CHANNEL, []string{"game_id"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	filter := &tool.Filter{
		Where:  where,
		Fields: []string{"Channelid", "Cp", "Name"},
	}

	ss := []models.Channel{}
	err = tool.GetAllByFilter(new(models.Channel), &ss, filter)

	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(ss)
}

