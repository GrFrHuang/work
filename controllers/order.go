package controllers

import (
	"kuaifa.com/kuaifa/work-together/models"
	"github.com/tealeg/xlsx"
	"fmt"
	"time"
	"kuaifa.com/kuaifa/work-together/utils"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"github.com/astaxie/beego/orm"
	"strconv"
	"net/http"
	"strings"
	"encoding/json"
)

// 流水管理
type OrderController struct {
	BaseController
}

func (o *OrderController) URLMapping() {
	//o.Mapping("Get", o.Get)
	//o.Mapping("GetAll", o.GetAll)
}

// 获取游戏
// @Title 获取

// @Description 获取游戏
// @Success 200 {object} models.Order
// @Failure 403
// @router /game [get]
func (c *OrderController) GetGame() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_ORDER, []string{"game_id"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	filter := &tool.Filter{
		Where:  where,
		Fields: []string{"GameId", "GameName"},
	}

	var ss []models.Game
	err = tool.GetAllByFilter(new(models.Game), &ss, filter)

	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(ss)
}

// 获取渠道
// @Title 获取渠道
// @Description 获取渠道
// @Success 200 {object} models.Channel
// @Failure 403
// @router /channel [get]
func (c *OrderController) GetChannel() {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_ORDER, []string{"cp"})
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

// 获取主体
// @Title 获取主体
// @Description 获取主体
// @Success 200 {object} models.Channel
// @Failure 403
// @router /body [get]
func (c *OrderController) GetBody() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_ORDER, nil)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	var body []models.BodyMy

	body = append(body, models.BodyMy{Id: 0, BodyMy: "无主体"})
	body = append(body, models.BodyMy{Id: 1, BodyMy: "云端"})
	body = append(body, models.BodyMy{Id: 2, BodyMy: "有量"})

	models.AddAdditionInfo(body[:])

	c.RespJSONData(body)
}

// GetAll ...
// @Title Get All
// @Description get 查看流水
// @Param	query	query	string	false	"Filter. e.g. col1:in1,2,3;col2:>2|<=3 ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Success 200 {object} models.Order
// @Failure 403
// @router /data [get]
func (c *OrderController) GetAll() {
	total, money, profit, data, err := getOrderData(c)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	//c.RespJSONData(data)
	//models.AddProfitInfo(&data)

	_, err = tool.BuildFilter(c.Controller, 20)
	if err != nil {
		return
	}

	//sumProfit := models.GetSumProfit(filter.Where)
	c.RespJSONDataWithSumAndTotal(data, total, money, profit)
}

// GetUrl ...
// @Title Get Url
// @Description get 获取流水详情下载链接
// @Param	start	query	string	false	"e.g. 2017-01-01"
// @Param	end	query	string	false	"e.g. 2017-12-01"
// @Param	game_id	query	string	false	"e.g. 732"
// @Param	cp	query	string	false	"e.g. anfan"
// @Success 200 {object} map
// @Failure 403
// @router /url [get]
func (c *OrderController) GetUrl() {
	start, err := time.Parse("2006-01-02", c.GetString("start"))
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	end, err := time.Parse("2006-01-02", c.GetString("end"))
	if err != nil {
		c.RespJSON(bean.CODE_Params_Err, err.Error())
		return
	}
	//fmt.Println("test:", end.AddDate(0, -1, 0).Before(start))
	if end.Before(start) || start.AddDate(0, 1, -1).Before(end) {
		c.RespJSON(bean.CODE_Forbidden, "下载订单必须确定时间范围，且不能超过一个月")
		return
	}
	params := map[string]string{
		"start":     c.GetString("start"),
		"end":       c.GetString("end"),
		"game_id":   c.GetString("game_id"),
		"cp":        c.GetString("cp"),
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	}
	fmt.Println("params:", params)
	url, err := utils.GetDetailUrl(utils.ORDER_DETAIL_PATH, params)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	rst := map[string]string{"url": url}
	c.RespJSONData(rst)
}

// DownloadOrder ...
// @Title DownloadOrder
// @Description get 下载流水
// @Param	query	query	string	false	"Filter. e.g. col1:in1,2,3;col2:>2|<=3..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} []map[string]string
// @Failure 403
// @router /download [get]
func (c *OrderController) DownloadOrder() {
	c.Ctx.Input.SetParam("limit", "1000")
	c.Ctx.Input.SetParam("offset", "0")

	group_str := ""
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ";") {
			kv := strings.SplitN(cond, ":", 2)
			println("kv:", kv[0], kv[1])
			if len(kv) == 2 && kv[0] == "groupby" && len(kv[1]) > 2 {
				group_str = kv[1][2:]
				break
			}
		}
	}
	title_values := []string{"起始日期", "结束日期"}
	cols := []string{"start_time", "end_time"}
	groups := strings.Split(group_str, ",")

	if len(groups) == 1 {
		if groups[0] == "game_id" {
			title_values = append(title_values, "游戏名称")
			cols = append(cols, "game_name")
		} else if groups[0] == "cp" {
			title_values = append(title_values, "渠道名称")
			cols = append(cols, "cp_name")
		} else {
		}

	} else {
		title_values = append(title_values, "游戏名称", "渠道名称")
		cols = append(cols, "game_name", "cp_name")
	}
	title_values = append(title_values, "渠道对账", "渠道回款", "总流水", "利润")
	cols = append(cols, "is_cp_verified", "is_channel_verified", "total", "profit")

	_, _, _, data, err := getOrderData(c)
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	//models.AddProfitInfo(&data)

	file := xlsx.NewFile()
	sht, err := file.AddSheet("sheet1")
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	title := sht.AddRow()
	for _, col := range title_values {
		tmp_cell := title.AddCell()
		tmp_cell.Value = col
	}

	for _, row := range data {
		newRow := sht.AddRow()
		for _, value := range cols {
			tmp_value_cell := newRow.AddCell()
			tmp_value_cell.Value = fmt.Sprintf("%v", row[value])
		}
	}

	mAssert := models.NewAsset()
	keyword := c.GetString("query") + "|" + c.GetString("fields") + "|" + c.GetString("offset")
	mAssert.HashStr = tool.EncodeMd5(keyword)
	fileName, err := mAssert.GetPathByKeyword("xlsx")
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}

	e := file.Save(fileName)
	if e != nil {
		c.RespJSON(bean.CODE_Forbidden, e.Error())
		return
	}

	fmt.Println("now:", time.Now().Format("2006-01-02"))
	c.AllowCross()
	c.Ctx.ResponseWriter.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=YUND-%s-%s", strings.Replace(time.Now().Format("2006-01-02"), "-", "", 3), mAssert.HashStr[:8]+".xlsx"))
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/vnd.open")
	aFile := mAssert.GetFilePath()
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, aFile)
}

func getOrderData(c *OrderController) (total int64, money, profit float64, params []orm.Params, err error) {
	where, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_ORDER, nil)
	if err != nil {
		return
	}

	filter, err := tool.BuildFilter(c.Controller, 20)
	if err != nil {
		return
	}

	tool.InjectPermissionWhere(where, &filter.Where)
	total, money, profit, params, err = models.GetGroupData(filter.Offset, filter.Limit, filter.Where, filter.Order, filter.Sortby)
	return
}

// GetAll ...
// @Title Get All
// @Description get GameIncome
// @Param   games  		query 	number	false		    "游戏id"
// @Param	channels		query 	string	false		"渠道id"
// @Param   start_time		query 	number	false		"开始时间"
// @Param	end_time		query 	number	false		"结束时间"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} {}
// @Failure 403
// @router /gameIncome [get]
func (c *OrderController) GetGameIncome() {
	// Todo：添加权限
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_ORDER, []string{"game_id, cp"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	games := c.GetString("games")
	channels := c.GetString("channels")
	startTime := c.GetString("start_time")
	endTime := c.GetString("end_time")
	limit, _ := c.GetInt("limit")
	offset, _ := c.GetInt("offset")

	fmt.Printf("games:%s, channels:%s, start:%s, end:%s limit:%d, offset:%d", games, channels, startTime, endTime, limit, offset)

	var game_ids []interface{}
	for _, v := range strings.Split(games, ",") {
		if v != "" {
			game_ids = append(game_ids, v)
		}
	}
	fmt.Println(fmt.Sprintf("%v", game_ids))

	var channel_codes []interface{}
	for _, v := range strings.Split(channels, ",") {
		if v != "" {
			channel_codes = append(channel_codes, v)
		}
	}

	param, tmp, total, err := parseGameChannels(game_ids, channel_codes, limit, offset)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}

	models.HandleName(&tmp)

	fmt.Printf("param%s, tmp:%v, total:%d \n", param, tmp, total)

	// TODO: get the data from server by url
	//data, err := old_sys.GetGameIncome(param)

	// get the order

	c.RespJSONDataWithTotal(tmp, total)
}

func parseGameChannels(games []interface{}, channels []interface{}, limit int, offset int) (param string, tmp []orm.Params, total_num int64, err error) {
	condition := ""
	var cond []interface{}
	if games != nil && len(games) > 0 {
		holder := strings.Repeat(",?", len(games))
		condition = fmt.Sprintf("%sAND game_id in(%s) ", condition, holder[1:])
		cond = append(cond, games)
	}
	if channels != nil && len(channels) > 0 {
		chl_holder := strings.Repeat(",?", len(channels))
		condition = fmt.Sprintf("%sAND cp in(%s) ", condition, chl_holder[1:])
		for _, v := range channels {
			cond = append(cond, fmt.Sprintf("%s", v))
		}
	}

	if len(condition) > 3 {
		condition = " WHERE " + condition[3:]
	}
	fmt.Println("condition:", condition)

	// get the total of the row
	o := orm.NewOrm()
	var total []orm.Params
	sql_str_total := fmt.Sprintf("SELECT COUNT(*) AS total FROM (SELECT game_id, cp FROM `order` %s GROUP BY game_id, cp) AS o", condition)
	_, err = o.Raw(sql_str_total, cond...).Values(&total)
	if err != nil {
		return
	}
	total_num, _ = strconv.ParseInt(total[0]["total"].(string), 10, 64)

	cond = append(cond, offset)
	cond = append(cond, limit)
	var gc []orm.Params
	sql := fmt.Sprintf("SELECT game_id, cp FROM `order` %s GROUP BY game_id, cp LIMIT ?, ?", condition)
	_, err = o.Raw(sql, cond...).Values(&gc)
	if err != nil {
		return
	}
	param_bytes, _ := json.Marshal(gc)
	param = string(param_bytes)

	// TODO: delete the tmp var
	tmp = gc
	return
}

// 流水详情
// @Title 流水详情
// @Description 流水详情
// @Param   game_id  		query 	number	false		    "游戏id"
// @Param	channel_code    query 	string	false		    "渠道code"
// @Success 200 {object} models.Channel
// @Failure 403
// @router /detail [get]
func (c *OrderController) GetDetailOrder() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_ORDER, []string{"cp"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	game_id, err := c.GetInt("game_id", 0)
	channel_code := c.GetString("channel_code", "")
	if game_id == 0 || channel_code == "" || err != nil {
		c.RespJSON(bean.CODE_Params_Err, "error params")
		return
	}

	detail, err := models.GetOrderDetail(game_id, channel_code)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, err.Error())
		return
	}

	c.RespJSONData(detail)
}

// 渠道汇总
// @Title 流水详情
// @Description 流水详情
// @Param	channel_code    query 	string	false		    "渠道code"
// @Param	resp_person     query 	string	false		    "负责人"
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.[]*ChannelData
// @Failure 403
// @router /channeldata [get]
func (c *OrderController) GetChannelData() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_ORDER, []string{"cp"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	offset, _ := c.GetInt("offset")
	limit, _ := c.GetInt("limit")
	channelCodes := c.GetString("channel_code", "")
	respPerson := c.GetString("resp_person", "")
	var codes []interface{}
	for _, v := range strings.Split(channelCodes, ",") {
		if v != "" {
			codes = append(codes, v)
		}
	}

	var persons []interface{}
	for _, v := range strings.Split(respPerson, ",") {
		if v != "" {
			persons = append(persons, v)
		}
	}

	total, data := models.GetChannelData(codes, persons, offset, limit)
	c.RespJSONDataWithTotal(data, total)
}

// 渠道汇总详情
// @Title 渠道汇总详情
// @Description 渠道汇总详情
// @Param	channel_code    query 	string	false		    "渠道code"
// @Success 200 {object} models.ChannelDetail
// @Failure 403
// @router /channeldetail [get]
func (c *OrderController) GetDetailChannel() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_ORDER, []string{"cp"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}
	channelCode := c.GetString("channel_code", "")
	detail := models.GetChannelDetail(channelCode)
	c.RespJSONData(detail)
}

// 获取所有渠道负责人
// @Title 获取所有渠道负责人
// @Description 获取所有渠道负责人
// @Success 200 {object} models.ChannelDetail
// @Failure 403
// @router /responsepersons [get]
func (c *OrderController) GetAllReponsePerson() {
	_, err := models.CheckPermission(c.Uid(), bean.PMSA_SELECT, bean.PMSM_ORDER, []string{"cp"})
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, err.Error())
		return
	}

	persons, err := models.GetAllChannelResponsePerson()
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, err.Error())
		return
	}
	c.RespJSONData(persons)
}
