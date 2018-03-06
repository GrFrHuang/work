package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"kuaifa.com/kuaifa/work-together/utils/express"
	"strconv"
	"strings"
	"time"
	"fmt"
	"kuaifa.com/kuaifa/work-together/models/codes"
)

// 快递管理
type ExpressManageController struct {
	BaseController
}

// URLMapping ...
func (c *ExpressManageController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Post
// @Description create ExpressManage
// @Param	body		body 	models.ExpressManage	true		"body for ExpressManage content"
// @Success 201 {int} models.ExpressManage
// @Failure 403 body is empty
// @router / [post]
func (c *ExpressManageController) Post() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_INSERT, bean.PMSM_EXPRESS_MANAGE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var express models.ExpressManage
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &express); err == nil {

		//参数校验
		if express.Number == "" || express.ReceiveCompany == 0 || express.ReceiveAddress == "" ||
			express.ReceivePeople == 0 || express.BodyMy == 0 || express.SendType == 0 ||
			express.ContentId == 0 || express.SendDepartment == 0 ||
			express.SendPeople == 0 {
			c.RespJSON(bean.CODE_Forbidden, "参数缺少")
			return
		}
		//三个中必有一个
		if express.IncludeGames == "" && express.Details == "" && express.Channels == "" {
			c.RespJSON(bean.CODE_Params_Err, "请输入合理的数据!")
			return
		}

		if err = models.CheckExpressManageByNumber(express.Number); err == nil {
			c.RespJSON(bean.CODE_Params_Err, "已存在该记录!")
			return
		}

		express.CreateTime = time.Now().Unix()
		express.CreatePeople = c.Uid()

		if _, err := models.AddExpressManage(&express); err == nil {
			c.RespJSONData("OK")
		} else {
			c.RespJSON(bean.CODE_Forbidden, err.Error())
		}
	} else {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
	}
}

// GetOne ...
// @Title Get One
// @Description get ExpressManage by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ExpressManage
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ExpressManageController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetExpressManageById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	} else {
		//contact, err := models.GetContactById(v.ReceivePeople)
		//if err != nil {
		//	c.RespJSON(bean.CODE_Params_Err,err.Error())
		//	return
		//}
		//sendUser, err := models.GetUserById(v.SendPeople)
		//if err != nil{
		//	c.RespJSON(bean.CODE_Params_Err,err.Error())
		//	return
		//}
		//company, err := models.GetCompanyTypeById(v.ReceiveCompany)
		//if err != nil{
		//	c.RespJSON(bean.CODE_Params_Err,err.Error())
		//	return
		//}
		//
		//
		//v.ReceiveUserName = contact.Name
		//v.SenderNickName = sendUser.Nickname
		//v.CompanyName = company.Name

		c.RespJSONData(v)
	}
	//c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get ExpressManage
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ExpressManage
// @Failure 403
// @router / [get]
func (c *ExpressManageController) GetAll() {

	total, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}
	c.RespJSONDataWithTotal(ss, total)

}

func (c *ExpressManageController) getAll() (total int64, ss []models.ExpressManage, errCode int, err error) {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_EXPRESS_MANAGE, nil)
	if err != nil {
		errCode = bean.CODE_Forbidden
		return
	}

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		errCode = bean.CODE_Params_Err
		return
	}

	filter.Fields = append(filter.Fields, "id", "Number", "BodyMy", "IncludeGames", "Channels", "Details",
		"TypeCompany", "SendType", "Icon", "ReceiveCompany", "SendPeople", "ContentId", "CreateTime",
		"CreatePeople", "SendDepartment", "ReceivePeople")
	filter.Order = []string{"desc"}
	filter.Sortby = []string{"create_time"}

	//按收件公司查 id
	idStrCompany := c.GetString("idStrCompany")
	idC := []interface{}{}
	for _, v := range strings.Split(idStrCompany, ",") {
		if v != "" {
			idC = append(idC, v)
		}
	}
	if len(idC) != 0 {
		filter.Where["receive_company__in"] = idC
	}

	//按我方主体搜索
	body, err := c.GetInt("body")
	if err == nil {
		filter.Where["body_my__exact"] = []interface{}{body}
	}

	//按发件人搜索  id
	sendUserStr := c.GetString("sendUsers")
	sendUsers := []interface{}{}
	for _, v := range strings.Split(sendUserStr, ",") {
		if v != "" {
			sendUsers = append(sendUsers, v)
		}
	}
	if len(sendUsers) != 0 {
		filter.Where["send_people__in"] = sendUsers
	}
	//按发件部门搜索 id
	departmentStr := c.GetString("departments")
	departments := []interface{}{}
	for _, v := range strings.Split(departmentStr, ",") {
		if v != "" {
			departments = append(departments, v)
		}
	}
	if len(departments) != 0 {
		filter.Where["send_department__in"] = departments
	}

	//按邮寄内容搜索 id
	contentId, err := c.GetInt("contentid")
	if err == nil {
		filter.Where["content_id__exact"] = []interface{}{contentId}
	}

	//按快递订单号搜索
	number := c.GetString("number")
	if number != "" {
		filter.Where["number__exact"] = []interface{}{number}
	}

	tool.InjectPermissionWhere(where, &filter.Where)
	ss = []models.ExpressManage{}
	total, err = tool.GetAllByFilterWithTotal(new(models.ExpressManage), &ss, filter)

	//添加收件公司信息
	models.AddReceiverCompanyInfo(&ss)
	//添加发件人信息 仅个人信息
	models.AddSendUser(&ss)
	//添加发件人所在的部门信息
	models.AddSendDepartment(&ss)
	//添加邮寄内容信息
	models.AddSendTypeInfo(&ss)
	//添加收件人信息
	models.GetReceiverByContactId(&ss)

	return
}

// @router /download [get]
func (c *ExpressManageController) DownLoad() {
	c.Ctx.Input.SetParam("limit", "0")
	c.Ctx.Input.SetParam("offset", "0")

	_, ss, errCode, err := c.getAll()
	if err != nil {
		c.RespJSON(errCode, err.Error())
		return
	}

	//rs := make([]map[string]interface{}, len(ss))
	rs := []map[string]interface{}{}
	for _, v := range ss {
		r := map[string]interface{}{}

		if (v.BodyMy == 1) {
			r["bodyMy"] = "云端"
		} else if (v.BodyMy == 2) {
			r["bodyMy"] = "有量"
		} else {
			r["bodyMy"] = "无主体"
		}
		if (v.Company != nil) {
			r["company"] = v.Company.Name
		}
		if (v.ContentId == 1) {
			r["content"] = "合同"
		} else if (v.ContentId == 2) {
			r["content"] = "发票"
		} else {
			r["content"] = "其他"
		}

		r["number"] = v.Number
		rs = append(rs, r)
	}

	cols := []string{"bodyMy", "company", "content", "number"}
	maps := map[string]string{
		cols[0]: "我方主体",
		cols[1]: "收件公司",
		cols[2]: "邮寄内容",
		cols[3]: "快递单号",
	}
	tmpFileName := fmt.Sprintf("快递管理-%s", time.Now().Format("20060102150405"))
	c.RespExcel(rs, tmpFileName, cols, maps)
}

// Put ...
// @Title Put
// @Description update the ExpressManage
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ExpressManage	true		"body for ExpressManage content"
// @Success 200 {object} models.ExpressManage
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ExpressManageController) Put() {

	where, err := models.CheckPermission(c.Uid(), bean.PMSA_UPDATE, bean.PMSM_EXPRESS_MANAGE, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.ExpressManage{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	if v.Number == "" || v.ReceiveCompany == 0 || v.ReceivePeople == 0 || v.ReceiveAddress == "" ||
		v.BodyMy == 0 || v.SendType == 0 || v.ContentId == 0 || v.SendPeople == 0 {
		c.RespJSON(bean.CODE_Params_Err, "请输入合理的数据!")
		return
	}
	if v.IncludeGames == "" && v.Details == "" && v.Channels == "" {
		c.RespJSON(bean.CODE_Params_Err, "请输入合理的数据!")
		return
	}

	if err := models.UpdateExpressManageById(&v, where, c.Uid()); err == nil {
		c.RespJSONData("保存成功")
	} else {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
	}
}

// Delete ...
// @Title Delete
// @Description delete the ExpressManage
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
//func (c *ExpressManageController) Delete() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	if err := models.DeleteExpressManage(id); err == nil {
//		c.Data["json"] = "OK"
//	} else {
//		c.Data["json"] = err.Error()
//	}
//	c.ServeJSON()
//}

// Get ...
// @Title 收件人
// @Description 通过companyid  获得收件人下拉
// @Success 200
// @Failure 403
// @router /Express/receiver/:id [get]
func (c *ExpressManageController) GetReceiver() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	flag, _ := c.GetInt("flag")
	contacts, err := models.GetContactByCompanyId(id, flag)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(contacts)

}

// Get ...
// @Title 发部门 一级
// @Description 级联 通过部门 到 具体某个人下
// @Success 200
// @Failure 403
// @router /Express/department [get]
func (c *ExpressManageController) GetSendDepartment() {
	departments := models.GetAllDepartment()
	c.RespJSONData(departments)
}

// Get ...
// @Title 发件人 二级
// @Description 级联 通过部门 到 具体某个人下    通过id 确定部门
// @Success 200
// @Failure 403
// @router /Express/sender/:id [get]
func (c *ExpressManageController) GetSender() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	users, err := models.GetUsersByDevMent(id)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(users)
}

// Get ...
// @Title 邮寄内容 合同 下拉
// @Description 通过 收件公司id(除了渠道类型外) 得出合同为 “未签订”、“已到期”的gameid 在通过gameid 获得相应的游戏
// @Success 200
// @Failure 403
// @router /Express/contract/games/:id [get]
func (c *ExpressManageController) GetContractGames() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	flag := c.GetString("flag")
	var contracts []models.Contract
	var err error
	if flag == "hetong" { //类型为渠道商时
		//根据companyid 获取公司的channel_code
		channel_codes, err := models.GetChannelCodesByCompanyId(id)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
		//根据channel_code获取合同的状态
		contracts, err = models.GetGameIdByChannelCodesAndState(channel_codes)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
	} else {
		//发行商
		contracts, err = models.GetGameIdByCompanyIdAndState(id)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
	}
	length := len(contracts)
	if length == 0 {
		c.RespJSON(bean.CODE_Params_Err, "该公司没有该类型数据!")
		return
	}
	links := make([]int, length)
	for i, s := range contracts {
		links[i] = s.GameId
	}
	games, err := models.GetGamesByGameIds(links)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	c.RespJSONData(games)

}

type ChannelInfo struct {
	ChannelName string `json:"channel_name,omitempty"`
	VerifyTime  int    `json:"verify_time,omitempty"`
}

// Get ...
// @Title 发票 下拉
// @Description 通过 收件公司id 得出该公司对应下的渠道 channel_company(获得channel_code) --> channel -->verify_channel
// @Param	id		path 	string	true		"The id you want
// @Success 200
// @Failure 403
// @router /Express/channel/:id [get]
func (c *ExpressManageController) GetChannelCompany() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	var channelInfo ChannelInfo
	channelCode, err := models.GetChannelCodeByCompanyId(id)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, "没有该记录!")
		return
	}
	if channelCode == "" {
		c.RespJSON(bean.CODE_Params_Err, "没有该记录!")
		return
	}
	channelName, err := models.GetChannelNameByCp(channelCode)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, "没有该记录!")
		return
	}
	if channelName == "" {
		c.RespJSON(bean.CODE_Params_Err, "没有该记录!")
		return
	}
	verifyTime, err := models.GetVerifyTimeByChannelCode(channelCode)
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, "没有该记录!")
		return
	}
	if verifyTime == 0 {
		c.RespJSON(bean.CODE_Params_Err, "没有该记录!")
		return
	}
	channelInfo.ChannelName = channelName
	channelInfo.VerifyTime = verifyTime
	c.RespJSONData(channelInfo)
}

// @Description 下拉框 订单号 仅 有id number信息
// @router /expressList/ [get]
func (c *ExpressManageController) GetExpressNumberList() {
	expresses, err := models.GetExpressNumbers()
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	c.RespJSONData(expresses)
}

// @Description 下拉框 收件公司 由公司类型产生 1--研发商 2--发行商 3--渠道商
// @Param	id		path 	string	true		"The id you want to delete"
// @router /companiesList/:id [get]
func (c *ExpressManageController) GetWantedCompany() {
	idStr := c.Ctx.Input.Param(":id")
	type_id, _ := strconv.Atoi(idStr)
	if type_id == 1 { //研发商
		companies, err := models.GetAllDevelopCompanies()
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
		length := len(companies)
		if length == 0 {
			c.RespJSON(bean.CODE_Params_Err, "该类型的数据为空!")
			return
		}
		links := make([]int, length)
		for i, s := range companies {
			links[i] = s.CompanyId
		}
		wantedCompanies, err := models.GetCompaniesByIds(links)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
		c.RespJSONData(wantedCompanies)

	} else if type_id == 2 { //发行商
		wantedCompanies, err := models.GetAllDistributionCompanies()
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
		c.RespJSONData(wantedCompanies)

	} else if type_id == 3 { //渠道商
		wantedCompanies, err := models.GetAllChannelConmpaniesNew()
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
		c.RespJSONData(wantedCompanies)
	} else {
		c.RespJSON(bean.CODE_Params_Err, "类型有误!")
	}
}

// @Description 获取公司地址
// @router /getAddress/:id [get]
func (c *ExpressManageController) GetAddress() {
	flag := c.GetString("flag")
	idStr := c.Ctx.Input.Param(":id")
	company_id, _ := strconv.Atoi(idStr)
	var address string
	var err error
	if flag == "1" { //研发
		address, err = models.GetDevAddressByCompanyId(company_id)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}

	} else if flag == "2" { //发行
		address, err = models.GetDistAddressByCompanyId(company_id)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}

	} else if flag == "3" { //渠道
		address, err = models.GetChannelAddressByCompanyId(company_id)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Params_Err, "flag 错误")
		return
	}
	if address == "" {
		c.RespJSON(bean.CODE_Params_Err, "暂没有该公司地址记录!")
		return
	}
	c.RespJSONData(address)
}

// @Description 获取 物流 详情
// @router /logistics/ [get]
func (c *ExpressManageController) Logistics() {
	num := c.GetString("num")
	com := c.GetString("com")
	rsp, err := express.GetExpressInfoBy100(com, num)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	c.RespJSONData(rsp)
}

// @Description 获取 物流 详情
// @router /get_express_manage [get]
func (c *ExpressManageController) Get() {

	send_type, err := c.GetInt("send_type", 0)
	if err != nil {
		c.RespDataMsg(codes.Params_Err, "参数错误")
		return
	}
	total, e, err := models.GetDataBySendType(send_type)
	if err != nil {
		c.RespJSON(codes.Params_Err, "查询失败")
		fmt.Println(err.Error())
		return
	}
	c.RespJSONDataWithTotal(e, total)
	return

}
