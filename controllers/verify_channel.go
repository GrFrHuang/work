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
	"github.com/astaxie/beego"
	"net/http"
	"bytes"
	"errors"
)

// 渠道对账单
type VerifyChannelController struct {
	BaseController
}

// URLMapping ...
func (c *VerifyChannelController) URLMapping() {
}

// @Title Post
// @Description create VerifyChannel
// @Param	body		body 	models.VerifyChannel	true		"body for VerifyChannel content"
// @Success 201 {int} models.VerifyChannel
// @Failure 403 body is empty
// @router / [post]
func (c *VerifyChannelController) Post()  {
	var v models.VerifyChannel

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}

	v.CreatedUserId = c.Uid()
	v.UpdatedUserId = c.Uid()
	_, err := models.AddVerifyChannel(&v)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	//if v.Date >= "2018-02" && v.Status == 10 {
	//	sTime, err := time.Parse("2006-01", v.Date)
	//	if err != nil {
	//		c.RespJSON(bean.CODE_Bad_Request, "同步失败")
	//		return
	//	}
	//	fistDay := time.Date(sTime.Year(), sTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	//	lastDay := fistDay.AddDate(0, 1, -1)
	//	var billGames []*models.LyBillGame
	//	var totalAmount float64
	//	var totalAmountPayable float64
	//	for _, item := range v.PreVerifyGames {
	//		temp := &models.LyBillGame{
	//			AnysdkgameId: item.GameId,
	//			GameName: item.GameName,
	//			Date: fistDay.Format("2006-01-02") + " - " + lastDay.Format("2006-01-02"),
	//			Amount: item.Amount,
	//			AmountPayable: item.AmountPayable,
	//		}
	//		totalAmount += item.Amount
	//		totalAmountPayable += item.AmountPayable
	//		billGames = append(billGames, temp)
	//	}
	//	games, err := json.Marshal(billGames)
	//	if err != nil {
	//		c.RespJSON(bean.CODE_Bad_Request, "同步失败")
	//		return
	//	}
	//	bill := &models.LyUserBill{
	//		WtId: int(id),
	//		ChannelPlatform: v.ChannelCode,
	//		BodyMy: v.BodyMy,
	//		Games: string(games),
	//		StartTime:fistDay,
	//		EndTime: lastDay,
	//		TotalAmount:totalAmount,
	//		TotalAmountPayable: totalAmountPayable,
	//		Status: int8(v.Status),
	//		CreateTime: v.CreatedTime,
	//		UpdateTime: v.UpdatedTime,
	//		TicketUrl: strconv.Itoa(v.FileId),
	//		Cmd: 1,
	//	}
	//	if err = sync(bill); err != nil {
	//		fmt.Println(err.Error())
	//		c.RespJSON(bean.CODE_Bad_Request, "同步失败")
	//		return
	//	}
	//}

	c.RespJSONData("OK")
}

func sync(data *models.LyUserBill) error {
	var url = beego.AppConfig.String("cps_url") + "/ly/bill/sync"
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return errors.New("同步失败")
	}
	return nil
}

// @Title Get One
// @Description get VerifyChannel by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.VerifyChannel
// @Failure 403 :id is empty
// @router /:id [get]
func (c *VerifyChannelController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetVerifyChannelById(id)
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
func (c *VerifyChannelController) GetAll() {
	data, total, err := c.getAll()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONDataWithTotal(data, total)
}

// unused 用处不大, 但很费时
// @router /download [get]
func (c *VerifyChannelController) DownLoad() {
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
		r["cp"] = v.Channel.Name
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

func (c *VerifyChannelController) getAll() (data []models.VerifyChannel, total int64, err error) {
	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		return
	}
	channel_str := c.GetString("channels")
	start := c.GetString("start", "2006-01")
	end := c.GetString("end", "2200-01")
	status, _ := c.GetInt("status", 0)

	verifyUserId, _ := c.GetInt("user_id")
	bodyMy, _ := c.GetInt("body_my")

	channelsCodes := []interface{}{}
	if channel_str != "" {
		for _, v := range strings.Split(channel_str, ",") {
			channelsCodes = append(channelsCodes, v)
		}
	}
	if len(channelsCodes) != 0 {
		filter.Where["channel_code__in"] = channelsCodes
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

	vs := []models.VerifyChannel{}
	total, err = tool.GetAllByFilterWithTotal(new(models.VerifyChannel), &vs, filter)
	if err != nil {
		return
	}
	models.AddPreVerifyChannelForVerifyList(&vs)
	models.AddChannelForVerifyList(&vs)
	models.AddVerifyAndUpdateUserForVerifyList(&vs)
	data = vs
	return
}

// @Title 更新对账单
// @Description 更新对账单
// @Param	id		path 	string					true		"The id you want to update"
// @Param	body	body 	models.VerifyChannel	true		"body for VerifyChannel content"
// @Success 200 {object} models.VerifyChannel
// @Failure 403 :id is not int
// @router /:id [put]
func (c *VerifyChannelController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.VerifyChannel{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}

	v.UpdatedUserId = c.Uid()
	if err := models.UpdateVerifyChannelById(&v); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	//if v.Date >= "2018-02" {
	//	sTime, err := time.Parse("2006-01", v.Date)
	//	if err != nil {
	//		c.RespJSON(bean.CODE_Bad_Request, "同步失败")
	//		return
	//	}
	//	fistDay := time.Date(sTime.Year(), sTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	//	lastDay := fistDay.AddDate(0, 1, -1)
	//	var billGames []*models.LyBillGame
	//	var totalAmount float64
	//	var totalAmountPayable float64
	//	for _, item := range v.PreVerifyGames {
	//		temp := &models.LyBillGame{
	//			AnysdkgameId:  item.GameId,
	//			GameName:      item.GameName,
	//			Date:          fistDay.Format("2006-01-02") + " - " + lastDay.Format("2006-01-02"),
	//			Amount:        item.Amount,
	//			AmountPayable: item.AmountPayable,
	//		}
	//		totalAmount += item.Amount
	//		totalAmountPayable += item.AmountPayable
	//		billGames = append(billGames, temp)
	//	}
	//	games, err := json.Marshal(billGames)
	//	if err != nil {
	//		c.RespJSON(bean.CODE_Bad_Request, "同步失败")
	//		return
	//	}
	//	bill := &models.LyUserBill{
	//		WtId:               int(id),
	//		ChannelPlatform:    v.ChannelCode,
	//		BodyMy:             v.BodyMy,
	//		Games:              string(games),
	//		TotalAmount:        totalAmount,
	//		TotalAmountPayable: totalAmountPayable,
	//		StartTime:			fistDay,
	//		EndTime: 			lastDay,
	//		Status:             int8(v.Status),
	//		CreateTime:         v.CreatedTime,
	//		UpdateTime:         v.UpdatedTime,
	//		TicketUrl:          strconv.Itoa(v.FileId),
	//		Cmd:                2,
	//	}
	//	if err = sync(bill); err != nil {
	//		fmt.Println(err)
	//		c.RespJSON(bean.CODE_Bad_Request, "同步失败")
	//		return
	//	}
	//}
	c.RespJSONData("OK")
}

// @Title Delete
// @Description delete the VerifyChannel
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *VerifyChannelController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteVerifyChannel(id,c.Uid()); err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	//if v.Date >= "2018-02" {
	//	bill := &models.LyUserBill{
	//		WtId: int(id),
	//		Cmd:  3,
	//	}
	//	if err := sync(bill); err != nil {
	//		c.RespJSON(bean.CODE_Bad_Request, "同步失败")
	//		return
	//	}
	//}
	c.RespJSONData("OK")
}

// @Title 获取未对账的渠道
// @Description delete the VerifyChannel
// @Param	body_my		query 	string	true		"我方主体"
// @router /not_verify_channel [get]
func (c *VerifyChannelController) GetNotVerifyChannel() {
	bodyMy, _ := c.GetInt("body_my", 0)
	if bodyMy == 0 {
		c.RespJSON(bean.CODE_Params_Err, "body_my can't be empty")
		return
	}

	channels, err := models.GetNotVerifyChannel(bodyMy)
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
// @Param	channel_code		query 	string	true		"渠道code"
// @router /not_verify_channel_time [get]
func (c *VerifyChannelController) GetNotVerifyChannelTime() {
	bodyMy, _ := c.GetInt("body_my", 0)
	if bodyMy == 0 {
		c.RespJSON(bean.CODE_Params_Err, "body_my can't be empty")
		return
	}
	channelCode := c.GetString("channel_code")
	if channelCode == "" {
		c.RespJSON(bean.CODE_Params_Err, "channelCode can't be empty")
		return
	}

	date, err := models.GetNotVerifyChannelTime(bodyMy, channelCode)
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
// @Param	channel_code		query 	string	true		"渠道code"
// @Param	month		query 	string	true		"月份"
// @router /not_verify_channel_game [get]
func (c *VerifyChannelController) GetNotVerifyChannelGame() {
	bodyMy, _ := c.GetInt("body_my", 0)
	if bodyMy == 0 {
		c.RespJSON(bean.CODE_Params_Err, "body_my can't be empty")
		return
	}
	channelCode := c.GetString("channel_code")
	if channelCode == "" {
		c.RespJSON(bean.CODE_Params_Err, "channelCode can't be empty")
		return
	}

	month := c.GetString("month")
	if month == "" {
		c.RespJSON(bean.CODE_Params_Err, "month can't be empty")
		return
	}

	date, err := models.GetNotVerifyChannelGame(bodyMy, channelCode, month)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(date)
	return
}

// @Title 获取渠道未对账的信息
// @Description 获取渠道未对账的信息(日期，主体，渠道，流水)
// @router /not_verify_info [get]
func (c *VerifyChannelController) GetNotVerifyInfo() {
	limit, _ := c.GetInt("limit", 20)
	offset, _ := c.GetInt("offset", 0)
	count,notVerify, err := models.GetChannelNotVerifyInfo(limit,offset)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONDataWithTotal(notVerify,count)
	return
}

// 获取所有渠道
// @Title 获取所有渠道
// @Description 获取所有渠道
// @Success 200 {object} models.Channel
// @Failure 403
// @router /channels [get]
func (c *VerifyChannelController) GetChannel() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, []string{"game_id"})
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

// @Title 获取某渠道回款主体
// @Description 通过渠道码和日期获取所有未对账的游戏
// @Param	cp	query	string	true	"渠道码"
// @Param	start		query	string	true	"开始时间 2006-01-02 格式"
// @Param	end		query	string	true	"结束时间 2006-01-02 格式"
// @Success 200 {object} models.GameAmount
// @Failure 403
// @router /remitcompanies [get]
func (c *VerifyChannelController) GetRemitCompanies() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_CHANNEL_VERIFY_ACCOUNT, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	channelCode := c.GetString("channel_code")

	res, err := models.GetRemitCompaniesByChannel(channelCode)
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(res)
}

//// @Title 从老的对账单生成一个新对账单
//// @Description 从老的对账单生成一个新对账单
//// @Param	id		path 	string	true		"The id you want to delete"
//// @Success 200 {string} delete success!
//// @Failure 403 id is empty
//// @router /migration/:id [get]
//func (c *VerifyChannelController) Migration() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.Atoi(idStr)
//	if err := models.MigrationVerifyChannel(id); err != nil {
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
//func (c *VerifyChannelController) MigrationAll() {
//	if err := models.MigrationVerifyChannelAll(); err != nil {
//		return
//	}
//
//	c.RespJSONData("OK")
//}

// @Title 简单统计信息
// @Description 获取昨天今天的新增账单
// @Success 200 {object} models.VerifyChannel
// @router /simple_statistics [get]
func (c *VerifyChannelController) SimpleStatistics() {
	v, err := models.GetChannelSimpleStatistics()
	if err != nil {
		c.RespJSON(bean.CODE_Bad_Request, err.Error())
		return
	}
	c.RespJSONData(v)
}
