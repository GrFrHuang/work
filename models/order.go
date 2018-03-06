package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
	"encoding/json"
	"math"
)

type Order struct {
	Id              int64
	GameID          int       `json:"game_id,omitempty" orm:"column(game_id);size(11)" `
	Cp              string    `orm:"size(100)" json:"cp,omitempty"`
	Date            time.Time `orm:"type(date)" json:"date,omitempty"`
	Status          int       `json:"status,omitempty"`
	Amount          float64   `json:"amount,omitempty"`
	CpVerified      int       `json:"cp_verified,omitempty" orm:"column(cp_verified);null" `
	ChannelVerified int       `json:"channel_verified,omitempty" orm:"column(channel_verified);null" `
	Settled         int       `json:"settled,omitempty" orm:"column(settled);null" `
	Remited         int       `json:"remited,omitempty" orm:"column(remited);null" `
	UpdateTime      int64
}

type BodyMy struct {
	Id     int    `json:"id"`
	BodyMy string `json:"body_my"`
	GameId string `json:"game_id"`
}

func init() {
	orm.RegisterModel(new(Order))
}

func (t *Order) TableName() string {
	return "order"
}

func GetDownloadData(fields []string, games []string, channels []string, start string, end string) []Order {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Order))

	var tmp_games []interface{}
	for _, v := range games {
		game_id, _ := strconv.Atoi(v)
		tmp_games = append(tmp_games, game_id)
	}
	var tmp_channels []interface{}
	for _, v := range channels {
		channel_id, _ := strconv.Atoi(v)
		tmp_channels = append(tmp_channels, channel_id)
	}
	if len(tmp_games) > 0 {
		qs = qs.Filter("Game__Id__in", tmp_games...)
	}
	if len(tmp_channels) > 0 {
		qs = qs.Filter("Channel__Id__in", tmp_channels...)
	}
	if start != "" {
		qs = qs.Filter("Date__gte", start)
	}
	if end != "" {
		qs = qs.Filter("Date__lte", end)
	}
	qs.GroupBy("GameId", "Cp", )
	var ords []Order
	qs.All(&ords, fields...)
	return ords
}

//order和sortBy暂时只支持按照总流水排序
func GetGroupData(offset int64, limit int64, where map[string][]interface{}, order []string, sortBy []string) (int64, float64, float64, []orm.Params, error) {

	games := where["game_id__in"]
	channels := where["cp__in"]
	body := where["body__in"]
	group_by := where["groupby__in"]
	date_gt := where["date__gte"]
	date_lt := where["date__lte"]

	fmt.Printf("where:%v, games: %v, channels: %v, group_by: %v, date_gt: %v, date_lt: %v ,body:%v\n",
		where, games, channels, group_by, date_gt, date_lt, body)

	if len(order) == 0 {
		order = append(order, "desc")
	}

	if len(sortBy) == 0 {
		sortBy = append(sortBy, "total")
	}

	condition := ""
	var cond []interface{}

	if games != nil && len(games) > 0 {
		holder := strings.Repeat(",?", len(games))
		condition = fmt.Sprintf("%sAND `order_profit_view`.game_id in(%s) ", condition, holder[1:])
		cond = append(cond, games)
	}

	if channels != nil && len(channels) > 0 {
		chl_holder := strings.Repeat(",?", len(channels))
		condition = fmt.Sprintf("%sAND `order_profit_view`.cp in(%s) ", condition, chl_holder[1:])
		for _, v := range channels {
			cond = append(cond, fmt.Sprintf("%s", v))
		}
	}

	if body != nil && len(body) > 0 {
		condition = fmt.Sprintf("%sAND `order_profit_view`.body_my = %s ", condition, body[0])

		//cond = append(cond, fmt.Sprintf("%s", body[0]))
	}

	if date_gt != nil && len(date_gt) > 0 {
		condition = fmt.Sprintf("%sAND `date`>=? ", condition)
		cond = append(cond, fmt.Sprintf("%s", date_gt[0]))
	} else {
		date_gt = []interface{}{" "}
	}

	if date_lt != nil && len(date_lt) > 0 {
		condition = fmt.Sprintf("%sAND `date`<=? ", condition)
		cond = append(cond, fmt.Sprintf("%s", date_lt[0]))
	} else {
		date_lt = []interface{}{" "}
	}

	if len(condition) > 3 {
		condition = " WHERE " + condition[3:]
	}
	group_fileds := ""
	if len(group_by) != 1 {
		group_fileds = "`order_profit_view`.game_id, cp"
	} else {
		group_fileds = "`order_profit_view`." + group_by[0].(string)
	}

	// total
	o := orm.NewOrm()
	var total_maps []orm.Params
	sql_str_total := fmt.Sprintf("SELECT COUNT(*) AS total, SUM(o.money) AS money,SUM(o.profit) AS profit  FROM (SELECT "+
		"SUM(`order_profit_view`.amount) AS money,SUM(`order_profit_view`.profit) AS profit "+
		"FROM `order_profit_view` %s "+
		"GROUP BY %s) AS o", condition, group_fileds)
	_, err1 := o.Raw(sql_str_total, cond...).Values(&total_maps)
	if err1 != nil {
		return 0, 0.0, 0.0, nil, err1
	}

	cond = append(cond, offset)
	cond = append(cond, limit)

	var maps []orm.Params
	sql_str := fmt.Sprintf("SELECT '" + fmt.Sprintf("%v", date_gt[0]) + "' as start_time, '"+
		fmt.Sprintf("%v", date_lt[0])+ "' as end_time, %s, sum(amount) as total,body_my,sum(profit) as profit"+
		" FROM `order_profit_view` "+
		" %s"+
		" GROUP BY %s "+
		" ORDER BY %s %s "+ " LIMIT ?, ?", group_fileds, condition, group_fileds, sortBy[0], order[0])

	_, err := o.Raw(sql_str, cond...).Values(&maps)
	if err != nil {
		return 0, 0.0, 0.0, nil, err
	}

	money := 0.0
	profit := 0.0
	if total_maps[0]["money"] != nil {
		money, _ = strconv.ParseFloat(total_maps[0]["money"].(string), 64)
	}
	if total_maps[0]["profit"] != nil {
		profit, _ = strconv.ParseFloat(total_maps[0]["profit"].(string), 64)
		profit = math.Trunc((profit+0.5/math.Pow10(2))*math.Pow10(2)) / math.Pow10(2)
	}

	i, _ := strconv.ParseInt(total_maps[0]["total"].(string), 10, 64)
	err = HandleName(&maps)
	return i, money, profit, maps, err
}

func HandleName(result *[]orm.Params) error {
	o := orm.NewOrm()
	game_sql := "SELECT game_id, game_name as name FROM `game_all`"
	var games []orm.Params
	_, err := o.Raw(game_sql).Values(&games)
	if err != nil {
		return err
	}
	game_map := make(map[string]string)
	for _, row := range games {
		if row["game_id"] != nil && row["name"] != nil {
			game_map[row["game_id"].(string)] = row["name"].(string)
		}
	}

	// get the {cp: name}
	channel_sql := "SELECT cp, name FROM `channel`"
	var channels []orm.Params
	_, err = o.Raw(channel_sql).Values(&channels)
	if err != nil {
		return err
	}
	channel_map := make(map[string]string)
	for _, row := range channels {
		if row["cp"] != nil && row["name"] != nil {
			channel_map[row["cp"].(string)] = row["name"].(string)
		}
	}

	// set the game_name and channel_name
	for _, row := range *result {
		if row["game_id"] != nil {
			game_name := game_map[row["game_id"].(string)]
			if game_name != "" {
				row["game_name"] = game_name
			} else {
				row["game_name"] = row["game_id"]
			}
		}
		if row["cp"] != nil {
			channel_name := channel_map[row["cp"].(string)]
			if channel_name != "" {
				row["cp_name"] = channel_name
			} else {
				row["cp_name"] = row["cp"]
			}
		}
	}
	return nil
}

func AddProfitInfo(ss *[]orm.Params) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	for _, row := range *ss {
		var profit Profit
		if row["game_id"] != nil && row["cp"] != nil {
			// 流水明细 有游戏和渠道
			sql := "SELECT body_my,ROUND(SUM(profit),2) as profit FROM profit WHERE game_id = ? AND channel_code = ? AND DATE >= ? AND DATE <= ? "
			o.Raw(sql, row["game_id"], row["cp"], row["start_time"], row["end_time"]).QueryRow(&profit)
		} else if row["game_id"] != nil && row["cp"] == nil {
			// 游戏汇总
			sql := "SELECT body_my,ROUND(SUM(profit),2) as profit FROM profit WHERE game_id = ? AND DATE >= ? AND DATE <= ? "
			o.Raw(sql, row["game_id"], row["start_time"], row["end_time"]).QueryRow(&profit)
		} else {
			// 渠道汇总
			sql := "SELECT body_my,ROUND(SUM(profit),2) as profit FROM profit WHERE channel_code = ? AND DATE >= ? AND DATE <= ? "
			o.Raw(sql, row["cp"], row["start_time"], row["end_time"]).QueryRow(&profit)
		}
		row["profit"] = profit.Profit
		row["bodyMy"] = profit.BodyMy
	}
	return
}

func GetSumProfit(where map[string][]interface{}) (sum float64) {

	games := where["game_id__in"]
	channels := where["cp__in"]
	date_gt := where["date__gte"]
	date_lt := where["date__lte"]

	sql := "SELECT ROUND(SUM(profit),2) AS profit FROM profit "
	if games != nil && len(games) > 0 {
		var gameName []string
		for _, game := range games {
			gameName = append(gameName, fmt.Sprintf("'%s'", game.(string)))
		}
		sql += fmt.Sprintf("WHERE game_id IN (%s) ", strings.Join(gameName, ","))
	}

	if channels != nil && len(channels) > 0 {
		var channelName []string
		for _, channel := range channels {
			channelName = append(channelName, fmt.Sprintf("'%s'", channel.(string)))
		}
		if strings.Contains(sql, "WHERE") {
			sql += fmt.Sprintf("AND channel_code IN(%s) ", strings.Join(channelName, ","))
		} else {
			sql += fmt.Sprintf("WHERE channel_code IN(%s) ", strings.Join(channelName, ","))
		}
	}

	if date_gt != nil && len(date_gt) > 0 && date_lt != nil && len(date_lt) > 0 {
		if strings.Contains(sql, "WHERE") {
			sql += fmt.Sprintf("AND `date` >= '%s' AND `date` <= '%s' ", date_gt[0], date_lt[0])
		} else {
			sql += fmt.Sprintf("WHERE `date` >= '%s' AND `date` <= '%s' ", date_gt[0], date_lt[0])
		}
	}

	o := orm.NewOrm()
	var profit Profit
	o.Raw(sql).QueryRow(&profit)

	return profit.Profit
}

func AddAdditionInfo(body []BodyMy) {

	o := orm.NewOrm()
	for i, v := range body {
		var result []orm.Params
		var gameIds []string
		o.Raw("SELECT DISTINCT(game_id) FROM profit WHERE body_my=?", v.Id).Values(&result)

		for _, res := range result {
			//id, _ := utils.Interface2Int(res["game_id"], false)
			gameIds = append(gameIds, res["game_id"].(string))
		}

		if len(gameIds) == 0 {
			gameIds = append(gameIds, "-1")
		}
		body[i].GameId = strings.Join(gameIds, ",")
	}
}

//按天获取总流水
func TodayTotal(date string) (*DashboardBasicInfo, error) {
	var dateInfo DashboardBasicInfo
	o := orm.NewOrm()
	err := o.Raw("select sum(amount) as `count`,`date` as date_time from `order` where `date`=? AND `status`=0", date).QueryRow(&dateInfo)
	return &dateInfo, err
}

//获取未对账的总金额
func NotReconTotal() (*DashboardBasicInfo, error) {
	var dateInfo DashboardBasicInfo
	nowMonth := time.Now().Format("2006-01")
	o := orm.NewOrm()
	err := o.Raw("SELECT sum(amount) as count FROM order_pre_verify_channel where verify_id = 0 and `date`<?", nowMonth).QueryRow(&dateInfo)
	return &dateInfo, err
}

//获取未回款总金额
func NotRemitTotal() (*DashboardBasicInfo, error) {
	dateInfo := DashboardBasicInfo{}
	o := orm.NewOrm()
	err := o.Raw("SELECT SUM(A.amount) AS count " +
		"FROM ( " +
		"SELECT SUM(amount_payable) AS amount FROM `verify_channel` WHERE `status` > 20 UNION " +
		"SELECT 0 - SUM(amount) AS amount FROM `remit_down_account` " +
		") AS A").QueryRow(&dateInfo)
	if err != nil {
		return nil, err
	}

	return &dateInfo, nil
}

// 获取游戏，在日期区间里某个渠道的总流水
func GetTotalMoney(gameId int, channelCode string, startTime string, endTime string) (money float64) {
	ord := Order{}
	o := orm.NewOrm()
	err := o.Raw("SELECT SUM(amount) AS amount FROM `order` WHERE game_id=? AND cp=? "+
		"AND date>=? AND date <=? AND status=0", gameId, channelCode, startTime, endTime).QueryRow(&ord)
	if err != nil {
		return 0
	}
	return ord.Amount
}

//获取未签合同数
func NotContractTotal() (*DashboardBasicInfo, error) {
	var dateInfo DashboardBasicInfo
	o := orm.NewOrm()
	err := o.Raw("select count(1) as `count` from contract where company_type=1 and state=149 and effective_state=1").QueryRow(&dateInfo)
	return &dateInfo, err
}

// 游戏流水详情查看
type OrderDetail struct {
	GameName                       string  `json:"game_name"`
	ChannelName                    string  `json:"channel_name"`
	ReleaseTime                    string  `json:"release_time"`         // 发行时间
	FirstDayOrder                  float64 `json:"first_day_order"`      // 首发日流水
	FirstDayUserCount              int     `json:"first_day_user_count"` // 首发日新增用户数
	CPResponsiblePersonName        string  `json:"cp_resp_name"`         // CP负责人姓名
	ChannelResponsiblePersonName   string  `json:"channel_resp_name"`    // 渠道商负责人
	OperationResponsiblePersonName string  `json:"operation_resp_name"`  // 运营负责人
}

// GetOrderDetail 获取流水详细查看的内容
func GetOrderDetail(gameId int, channelCode string) (detail *OrderDetail, err error) {
	detail = &OrderDetail{}
	o := orm.NewOrm()

	game := &Game{GameId: gameId}
	if err = o.Read(game, "GameId"); err != nil {
		return
	}
	detail.GameName = game.GameName
	detail.ReleaseTime = time.Unix(game.PublishTime, 0).Format("2006-01-02")
	if game.SubmitPerson != 0 {
		submitUser, err := GetUserById(game.SubmitPerson)
		if err == nil {
			detail.CPResponsiblePersonName = submitUser.Nickname
		}
	}

	channelCompany, err := GetChannelCompanyByCode(channelCode)
	if err != nil {
		return
	}
	var reps []string
	if channelCompany.YunduanResPerName != "" {
		reps = append(reps, channelCompany.YunduanResPerName)
	}
	if channelCompany.YouliangResPerName != "" {
		reps = append(reps, channelCompany.YouliangResPerName)
	}
	detail.ChannelResponsiblePersonName = strings.Join(reps, ",")
	detail.ChannelName = channelCompany.ChannelName

	date, _ := time.Parse("2016-01-02", detail.ReleaseTime)
	order := &Order{GameID: gameId, Cp: channelCode, Date: date}
	if err = o.Read(order, "GameId", "Cp"); err == nil {
		detail.FirstDayOrder = order.Amount
	}

	gamePlan := &GamePlan{GameId: gameId}
	if err = o.Read(gamePlan, "GameId"); err == nil {
		if gamePlan.OperatorPerson != "" {
			var userIds []interface{}
			err = json.Unmarshal([]byte(gamePlan.OperatorPerson), &userIds)
			if err == nil {
				users, err := GetUsersByIds(userIds)
				if err == nil {
					var names []string
					for _, v := range users {
						names = append(names, v)
					}
					detail.OperationResponsiblePersonName = strings.Join(names, ",")
				}
			}
		}
	}
	return
}

type ChannelData struct {
	ChannelCode string  `json:"channel_code"`
	ChannelName string  `json:"channel_name"`
	GameCount   int     `json:"game_count"` // 游戏数量
	Signed      int     `json:"signed"`     // 已签合同数
	NotSign     int     `json:"not_sign"`   // 未签合同数
	NotVerify   int     `json:"not_verify"` // 未对账数
	ShouldPay   float64 `json:"should_pay"` // 应付金额
	Paid        float64 `json:"paid"`       // 回款金额
	NotPay      float64 `json:"not_pay"`    // 未回款金额
	// YunduanRep      string    `json:"yunduan_rep"`          // 云端负责人
	// YouliangRep     string    `json:"youliang_rep"`         // 有量负责人
	Resps string `json:"resps"` // 负责人
}

// GetChannelData 获取渠道汇总数据
func GetChannelData(channelsCode []interface{}, usersId []interface{}, offset int, limit int) (total int64, channelsData []*ChannelData) {
	channelsData = []*ChannelData{}
	o := orm.NewOrm()
	qs := o.QueryTable("channel")
	if len(channelsCode) > 0 {
		qs = qs.Filter("cp__in", channelsCode...)
	}

	if len(usersId) > 0 {
		holder := strings.Repeat(",?", len(usersId))
		var tmpCodes []interface{}
		sql := fmt.Sprintf("SELECT channel_code FROM `channel_company` WHERE yunduan_responsible_person in(%s) OR youliang_responsible_person in(%s)", holder[1:], holder[1:])
		var channelCodes []orm.Params
		usersId2 := append(usersId, usersId...)
		_, err := o.Raw(sql, usersId2...).Values(&channelCodes)
		if err == nil {
			for _, v := range channelCodes {
				tmpCodes = append(tmpCodes, v["channel_code"])
			}
			fmt.Println("channels:", tmpCodes)
			qs = qs.Filter("cp__in", tmpCodes...)
		}
	}

	total, _ = qs.Count()

	result := map[string]*ChannelData{}
	var codes []interface{}

	// Get the channels
	var maps []orm.Params
	_, err := qs.Limit(limit).Offset(offset).Values(&maps, "cp", "name")
	if err != nil || len(maps) == 0 {
		return
	}

	for _, v := range maps {
		cp := v["Cp"].(string)
		name := v["Name"].(string)

		channelData := &ChannelData{}
		channelData.ChannelCode = cp
		channelData.ChannelName = name
		result[cp] = channelData

		codes = append(codes, cp)
	}

	// Add the number of game
	holder := strings.Repeat(",?", len(codes))
	sql := fmt.Sprintf("SELECT `channel_code`, COUNT(DISTINCT `game_id`) AS cnt FROM channel_access WHERE `channel_code` IN(%s) GROUP BY `channel_code`", holder[1:])
	var games []orm.Params
	_, err = o.Raw(sql, codes).Values(&games)
	if err == nil {
		for _, v := range games {
			result[v["channel_code"].(string)].GameCount, _ = strconv.Atoi(v["cnt"].(string))
		}
	}

	// Add the contract data
	sql = fmt.Sprintf("SELECT `channel_code`, COUNT(*) AS cnt FROM contract WHERE `channel_code` in(%s) AND `state`=149 AND `effective_state`=1 AND `company_type`=1 GROUP BY channel_code", holder[1:])
	var contractNotSign []orm.Params
	_, err = o.Raw(sql, codes).Values(&contractNotSign)
	if err == nil {
		for _, v := range contractNotSign {
			result[v["channel_code"].(string)].NotSign, _ = strconv.Atoi(v["cnt"].(string))
		}
	}
	sql = fmt.Sprintf("SELECT `channel_code`, COUNT(*) AS cnt FROM contract WHERE `channel_code` in(%s) AND  `state`=150 AND `company_type`=1 GROUP BY channel_code", holder[1:])
	var contractSigned []orm.Params
	_, err = o.Raw(sql, codes).Values(&contractSigned)
	if err == nil {
		for _, v := range contractSigned {
			result[v["channel_code"].(string)].Signed, _ = strconv.Atoi(v["cnt"].(string))
		}
	}

	// Add not verify
	//sql = fmt.Sprintf("SELECT channel_code, COUNT(*) AS cnt " +
	//	"FROM order_pre_verify_channel as o INNER JOIN contract AS con ON o.`channel_code`=con.`channel_code` AND o.game_id = con.game_id " +
	//	"WHERE channel_code in(%s) AND verify_id = 0 " +
	//	"GROUP BY o.channel_code, con.body_my, o.date", holder[1:])
	sql = fmt.Sprintf("SELECT c.channel_code, COUNT(*) AS cnt "+
		"FROM( "+
		"SELECT o.channel_code,con.body_my, o.`date` "+
		"FROM `order_pre_verify_channel` AS o INNER JOIN contract AS con ON o.`channel_code`=con.`channel_code` AND o.game_id = con.game_id "+
		"WHERE o.channel_code in(%s) AND o.verify_id = 0 "+
		"GROUP BY o.channel_code, con.body_my, o.`date` "+
		") AS c "+
		"GROUP BY c.channel_code ", holder[1:])
	var notVerify []orm.Params
	_, err = o.Raw(sql, codes).Values(&notVerify)
	if err == nil {
		for _, v := range notVerify {
			result[v["channel_code"].(string)].NotVerify, _ = strconv.Atoi(v["cnt"].(string))
		}
	}

	// Add should pay of money
	sql = fmt.Sprintf("SELECT channel_code, SUM(amount_payable) AS cnt FROM `verify_channel` WHERE `status` >20 AND channel_code in(%s) GROUP BY channel_code", holder[1:])
	var shouldPay []orm.Params
	_, err = o.Raw(sql, codes).Values(&shouldPay)
	if err == nil {
		for _, v := range shouldPay {
			result[v["channel_code"].(string)].ShouldPay, _ = strconv.ParseFloat(v["cnt"].(string), 64)
		}
	}

	// Add response person
	sql = fmt.Sprintf("SELECT channel_code, yunduan_responsible_person AS yunduan, youliang_responsible_person AS youliang FROM `channel_company` WHERE channel_code in(%s)", holder[1:])
	var respPerson []orm.Params
	_, err = o.Raw(sql, codes).Values(&respPerson)
	if err == nil {
		for _, v := range respPerson {
			var resps []interface{}
			if v["yunduan"] != nil {
				resps = append(resps, v["yunduan"])
			}
			if v["youliang"] != nil {
				resps = append(resps, v["youliang"])
			}
			userMap, _ := GetUsersByIds(resps)
			var respsName []string
			for _, v := range resps {
				id, _ := strconv.Atoi(v.(string))
				respsName = append(respsName, userMap[id])
			}
			result[v["channel_code"].(string)].Resps = strings.Join(respsName, ",")
		}
	}

	for _, v := range result {
		channelsData = append(channelsData, v)
	}
	return
}

type ChannelDetail struct {
	ChannelCode     string   `json:"channel_code"`
	ChannelName     string   `json:"channel_name"`
	NotSignGames    []string `json:"not_sign_games"`    // 未签订合同的所有游戏
	SignedGames     []string `json:"signed_games"`      // 已签订合同的所有游戏
	NotVerifyMonths []string `json:"not_verify_months"` // 未对账月份
	LatestPaidDates []string `json:"latest_paid_dates"` // 最近回款(显示最近5次)
}

// GetChannelDetail 获取渠道详情
func GetChannelDetail(channelCode string) (channel *ChannelDetail) {
	channel = &ChannelDetail{}
	channel.ChannelCode = channelCode
	// Get channel name
	channel.ChannelName, _ = GetChannelNameByCp(channelCode)

	o := orm.NewOrm()
	// Get the contracts
	sql := fmt.Sprintf("SELECT DISTINCT game_id FROM `contract` WHERE channel_code = '%s' AND `state`=149 AND company_type=1", channelCode)
	var signedGame []orm.Params
	_, err := o.Raw(sql).Values(&signedGame)
	fmt.Println(signedGame)
	if err == nil {
		var gameIds []int
		for _, v := range signedGame {
			id, err := strconv.Atoi(v["game_id"].(string))
			if err == nil {
				gameIds = append(gameIds, id)
			}
		}
		channel.NotSignGames = getGameNamesByIds(gameIds)
		fmt.Println(channel.NotSignGames)
	}
	sql = fmt.Sprintf("SELECT DISTINCT game_id FROM `contract` WHERE channel_code = '%s' AND `state`=150 AND company_type=1", channelCode)
	var notSignGame []orm.Params
	_, err = o.Raw(sql).Values(&notSignGame)
	if err == nil {
		var gameIds []int
		for _, v := range notSignGame {
			id, err := strconv.Atoi(v["game_id"].(string))
			if err == nil {
				gameIds = append(gameIds, id)
			}

		}
		channel.SignedGames = getGameNamesByIds(gameIds)
	}

	// Get Not verify month
	sql = fmt.Sprintf("SELECT o.channel_code, con.body_my, o.`date` "+
		"FROM `order_pre_verify_channel` AS o INNER JOIN contract AS con ON o.`channel_code`=con.`channel_code` AND o.game_id = con.game_id "+
		"WHERE o.verify_id = 0 AND o.`channel_code`='%s' "+
		"GROUP BY o.channel_code, con.body_my, o.`date` "+
		"ORDER BY o.`date` DESC", channelCode)
	var notVerify []orm.Params
	_, err = o.Raw(sql).Values(&notVerify)
	if err == nil {
		channel.NotVerifyMonths = []string{}
		for _, v := range notVerify {
			bodyMyId := 0
			bodyMyStr := v["body_my"]
			if bodyMyStr != nil {
				bodyMyId, _ = strconv.Atoi(v["body_my"].(string))
			}
			bodyMy := parseBodyMy(bodyMyId)
			notStr := bodyMy + "," + v["date"].(string)
			channel.NotVerifyMonths = append(channel.NotVerifyMonths, notStr)
		}
	}

	// TODO: Get Latest paid note
	channel.LatestPaidDates = []string{}
	return
}

// GetAllChannelResponsePerson 获取所有渠道负责人
func GetAllChannelResponsePerson() (users []*UserIdAndName, err error) {
	users = []*UserIdAndName{}
	o := orm.NewOrm()
	sql := "SELECT yunduan_responsible_person AS yunduan, youliang_responsible_person AS youliang FROM `channel_company`"
	var respPerson []orm.Params
	userIds := map[interface{}]bool{}
	_, err = o.Raw(sql).Values(&respPerson)
	if err != nil {
		return
	}
	for _, v := range respPerson {
		if v["yunduan"] != nil {
			userIds[v["yunduan"]] = true
		}
		if v["youliang"] != nil {
			userIds[v["youliang"]] = true
		}
	}
	var userIdsDistinct []interface{}
	for k := range userIds {
		userIdsDistinct = append(userIdsDistinct, k)
	}
	users, err = GetUsersListByIds(userIdsDistinct)
	return
}

func getGameNamesByIds(gameIds []int) (gameNames []string) {
	gameNames = []string{}
	if len(gameIds) == 0 {
		return
	}
	o := orm.NewOrm()
	var ids []interface{}
	for _, v := range gameIds {
		ids = append(ids, v)
	}
	var games []orm.Params
	_, err := o.QueryTable("game").Filter("game_id__in", ids...).Values(&games, "game_name")
	gameNames = []string{}
	if err == nil {
		for _, v := range games {
			gameNames = append(gameNames, v["GameName"].(string))
		}
	}
	return
}

func parseBodyMy(bodyId int) (bodyMy string) {
	if bodyId == 1 {
		bodyMy = "云端"
	} else if bodyId == 2 {
		bodyMy = "有量"
	} else {
		bodyMy = "错误的主体类型"
	}
	return
}
