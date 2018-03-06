package models
//
//import (
//	"errors"
//	"fmt"
//	"github.com/astaxie/beego/orm"
//	"github.com/bysir-zl/bygo/util"
//	"kuaifa.com/kuaifa/work-together/models/bean"
//	"kuaifa.com/kuaifa/work-together/utils"
//	"log"
//	"strconv"
//	"strings"
//	"time"
//	"kuaifa.com/kuaifa/work-together/tool/old_sys"
//	"encoding/json"
//	"github.com/astaxie/beego"
//)
//
//type ChannelVerifyAccount struct {
//	Id             int     `json:"id,omitempty" orm:"column(id);auto"`
//	StartTime      string  `json:"start_time,omitempty" orm:"column(start_time);type(date);null" valid:"Required"`
//	EndTime        string  `json:"end_time,omitempty" orm:"column(end_time);type(date);null" valid:"Required"`
//	Cp             string  `json:"cp,omitempty" orm:"column(cp);size(255);null" valid:"Required"`
//	GameStr        string  `json:"game_str,omitempty" orm:"column(games);null" valid:"Required"`
//	AmountMy       float64 `json:"amount_my,omitempty" orm:"column(amount_my);null;digits(16);decimals(2)" valid:"Required"` // 我方
//	AmountOpposite float64 `json:"amount_opposite,omitempty" orm:"column(amount_opposite);null;digits(16);decimals(2)" `     // 对方
//	AmountPayable  float64 `json:"amount_payable,omitempty" orm:"column(amount_payable);null;digits(16);decimals(2)" `       // 应收
//	AmountTheory   float64 `json:"amount_theory,omitempty" orm:"column(amount_theory);null;digits(16);decimals(2)"`          // 理论金额
//	AmountRemit    float64 `json:"amount_remit,omitempty" orm:"column(amount_remit);null;digits(16);decimals(2)" `           // 已回款
//	Status         int     `json:"status,omitempty" orm:"column(status);null" valid:"Required"`
//	// 回款主体, 在回款时用于关联对账单, 完成对账时在渠道表里的回款主体选择,
//	RemitCompanyId int    `json:"remit_company_id,omitempty" orm:"column(remit_company_id);null"`
//	CreateUserId   int    `json:"create_user_id,omitempty" orm:"column(create_user_id);null" `
//	VerifyUserId   int    `json:"verify_user_id,omitempty" orm:"column(verify_user_id);null" `
//	UpdateUserId   int    `json:"update_user_id,omitempty" orm:"column(update_user_id);null" `
//	VerifyTime     int    `json:"verify_time,omitempty" orm:"column(verify_time);null" valid:"Required"`
//	CreateTime     int    `json:"create_time,omitempty" orm:"column(create_time);null" `
//	UpdateTime     int    `json:"update_time,omitempty" orm:"column(update_time);null" `
//	Desc           string `json:"desc,omitempty" orm:"column(desc);null" `
//	TicketUrl      string `json:"ticket_url,omitempty" orm:"column(ticket_url);size(255);null" `
//	Extra          string `json:"extra,omitempty" orm:"column(extra);size(255);null" `
//	FileId         int    `json:"file_id,omitempty" orm:"column(file_id);null"`
//	FilePreviewId  int    `json:"file_preview_id,omitempty" orm:"column(file_preview_id);"`
//	BodyMy         int     `json:"body_my,omitempty" orm:"column(body_my);"`
//
//	CpName string        `json:"cp_name,omitempty" orm:"-" `
//	//Games        []*GameAmount `orm:"-" json:"games,omitempty"`
//	GameNum      int           `orm:"-" json:"game_num,omitempty"`
//	CreateUser   *User         `orm:"-" json:"create_user,omitempty"`
//	VerifyUser   *User         `orm:"-" json:"verify_user,omitempty"`
//	UpdateUser   *User         `orm:"-" json:"update_user,omitempty"`
//	Channel      *Channel      `orm:"-" json:"channel,omitempty"`
//	RemitCompany *Company      `orm:"-" json:"remit_company,omitempty"`
//}
//
//// 渠道对账和cp对账公用的单个游戏对账的模型,用于存json
//type GameAmount struct {
//	GameAll
//	Date           string  `json:"date,omitempty"`
//	AmountMy       float64 `json:"amount_my"`
//	AmountOpposite float64 `json:"amount_opposite,omitempty"`
//	AmountPayable  float64 `json:"amount_payable,omitempty"`
//
//	AmountTheory float64 `json:"amount_theory"`     // 理论金额
//	Formula      string  `json:"formula,omitempty"` // 算金额的公式
//}
//
////对账单的状态
//const (
//	CHAN_VERIFY_S_VERIFYING = 10  // 对账中
//	CHAN_VERIFY_S_VERIFYIED = 20  // 对账完毕,未开票
//	CHAN_VERIFY_S_RECEIPT   = 30  // 已开票, 此状态必须选择回款主体
//	CHAN_VERIFY_S_REMIT     = 100 // 已回款
//)
//
//var ChanStatusCode2String = map[int]string{
//	CHAN_VERIFY_S_VERIFYING: "对账中",
//	CHAN_VERIFY_S_VERIFYIED: "对账完毕,未开票",
//	CHAN_VERIFY_S_RECEIPT:   "已开票",
//
//	CHAN_VERIFY_S_REMIT: "已回款", // 此状态在实际=已回款金额自动生效
//}
//
//func (t *ChannelVerifyAccount) TableName() string {
//	return "channel_verify_account"
//}
//
//func init() {
//	orm.RegisterModel(new(ChannelVerifyAccount))
//}
//
//func GetChannelUpToDateTime(gameid []interface{}, channel []interface{}) (uptime int, err error) {
//	o := orm.NewOrm()
//	var res []orm.Params
//	var args []interface{}
//	create_str := "SELECT MIN(end_time) updatatime FROM cp_verify_account "
//	if len(gameid) != 0 {
//		create_str += "WHERE game_id = ? "
//		args = append(args, gameid...)
//	}
//	if len(channel) != 0 {
//		create_str += "WHERE channel "
//	}
//	num, err := o.Raw(create_str, args...).Values(&res)
//	log.Printf("res:%v", res)
//	if err == nil && num > 0 {
//		fmt.Println(res[0]["updatatime"])
//		var uptimetm time.Time
//		uptimetm, err = time.ParseInLocation("2006-01-02", res[0]["updatatime"].(string), time.Local)
//		if err != nil {
//			return 0, err
//		}
//		uptime = int(uptimetm.Unix())
//		return
//	}
//	return
//}
//
//// 获取渠道未对账信息
//func GetChanNoVerifyAccount(stime, etime string, games, channels []interface{}) (res *ChannelVerifyAccount, err error) {
//	_, err = time.Parse("2006-01-02", stime)
//	if err != nil {
//		return
//	}
//	_, err = time.Parse("2006-01-02", etime)
//	if err != nil {
//		return
//	}
//	qb, _ := orm.NewQueryBuilder("mysql")
//	gameNum := len(games)
//	chanNum := len(channels)
//	where := "date >= ? AND date <= ? AND channel_verified != 1 "
//	if gameNum > 0 {
//		create := strings.Repeat(",?", gameNum)[1:]
//		where = fmt.Sprintf("%sAND game_id IN (%s) ", where, create)
//	}
//	if chanNum > 0 {
//		create := strings.Repeat(",?", chanNum)[1:]
//		where = fmt.Sprintf("%sAND cp IN (%s)", where, create)
//	}
//	qb.Select("SUM(amount) total").From("`order`").
//		Where(where)
//	o := orm.NewOrm()
//	sql := qb.String()
//	temp := []interface{}{stime, etime}
//	temp = append(temp, games...)
//	temp = append(temp, channels...)
//	var step []orm.Params
//	_, err = o.Raw(sql, temp...).Values(&step)
//	if err != nil {
//		return
//	}
//	var sum float64
//	if step[0]["total"] != nil {
//		sum, err = strconv.ParseFloat(step[0]["total"].(string), 64)
//	} else {
//		sum = float64(0)
//	}
//	log.Printf("sum: %v", sum)
//	res = &ChannelVerifyAccount{
//		StartTime: stime,
//		EndTime:   etime,
//		AmountMy:  sum,
//		Status:    0,
//	}
//	if gameNum == 1 {
//		var id int
//		id, err = strconv.Atoi(games[0].(string))
//		if err != nil {
//			return
//		}
//		res.GameStr = GetGameAllNameById(id)
//	} else {
//		res.GameStr = "*"
//	}
//	if chanNum == 1 {
//		res.CpName, _ = GetChannelNameByCp(channels[0].(string))
//	} else {
//		res.CpName = "*"
//	}
//	return
//}
//
//// 获取未对账的渠道
//func GetNoVerifyChannels(where map[string][]interface{}) (channels []Channel, err error) {
//	orders := []Order{}
//	_, err = QueryTable(new(Order), where).Filter("ChannelVerified__in", 0, 2).GroupBy("Cp").All(&orders, "Id", "Cp")
//	if err != nil {
//		return
//	}
//	cps := make([]string, len(orders))
//	for _, v := range orders {
//		cps = append(cps, v.Cp)
//	}
//	o := orm.NewOrm()
//	chans := []Channel{}
//	_, err = o.QueryTable(new(Channel)).All(&chans, "Id", "Channelid", "Name", "Cp")
//	if err != nil {
//		return
//	}
//	chanMap := make(map[string]Channel)
//	for _, v := range chans {
//		chanMap[v.Cp] = v
//	}
//	for _, v := range orders {
//		val, ok := chanMap[v.Cp]
//		if ok {
//			channels = append(channels, val)
//		}
//	}
//
//	return
//}
//
//// 获取渠道的对账时间
//func GetVerifyDateByCp(cp string) (res []ChannelVerifyAccount, err error) {
//	o := orm.NewOrm()
//	var maps []orm.Params
//	o.Raw("SELECT LEFT (date, 7) AS months,MIN(date) AS mindate FROM `order` "+"WHERE cp = ? AND channel_verified!=1 GROUP BY months", cp).Values(&maps)
//	log.Printf("maps: %v", maps)
//	res = []ChannelVerifyAccount{}
//	if len(maps) == 0 {
//		return
//	}
//	channel := ChannelVerifyAccount{}
//	var tm time.Time
//	for _, value := range maps {
//		start := value["mindate"].(string)
//		tm, err = time.Parse("2006-01", value["months"].(string))
//		if err != nil {
//			return
//		}
//		channel.Cp = cp
//		channel.StartTime = start
//		channel.EndTime = tm.AddDate(0, 1, -1).Format("2006-01-02")
//		res = append(res, channel)
//	}
//	return
//}
//
//// 获取渠道对应的所有未对账游戏
//func GetGamesByCpAndDate(cp, stm, etm string, bodyMy int, where map[string][]interface{}) (gameAmounts []GameAmount, err error) {
//	gameAmounts = []GameAmount{}
//	_, err = time.Parse("2006-01-02", stm)
//	if err != nil {
//		return
//	}
//	_, err = time.Parse("2006-01-02", etm)
//	if err != nil {
//		return
//	}
//	var maps []orm.Params
//	channelid, err := GetChannelIdByCp(cp);
//	if err != nil {
//		return
//	}
//
//	o := orm.NewOrm()
//	sql := "SELECT order.id,order.game_id ,game_all.game_name " +
//		"FROM `order`  LEFT JOIN game_all ON game_all.game_id=`order`.game_id " +
//		"WHERE channel_verified != 1 AND cp = ? AND date >= ? AND date <= ? AND " +
//		"game_all.game_id IN (SELECT DISTINCT game_id FROM `contract` WHERE company_type = 1 AND company_id = ? AND body_my = ?)" +
//		"GROUP BY game_all.game_id"
//	o.Raw(sql, cp, stm, etm, channelid, bodyMy).Values(&maps)
//	if maps == nil || len(maps) == 0 {
//		return
//	}
//
//	var games []string
//	for _, game := range maps {
//		games = append(games, game["game_id"].(string))
//	}
//
//	// 从其他系统获取理论金额
//	clears, _ := old_sys.GetAllClearing(strings.Join(games, ","), cp, stm[0:7])
//	clears_map := make(map[string]old_sys.Clearings)
//	for _, clear := range *clears {
//		clears_map[clear.GameId] = clear
//	}
//
//	for _, value := range maps {
//		id, e := strconv.Atoi(value["game_id"].(string))
//		if e != nil {
//			err = e
//			return
//		}
//		//total, e := strconv.ParseFloat(value["total"].(string), 64)
//		//if e != nil {
//		//	err = e
//		//	return
//		//}
//
//		g := GameAmount{}
//		g.GameName, _ = util.Interface2String(value["game_name"], true)
//		g.GameId = id
//		g.AmountMy = clears_map[value["game_id"].(string)].Total
//		//clearing, _ := old_sys.GetClearing(id, cp, stm[0:7])
//		g.AmountTheory = clears_map[value["game_id"].(string)].DivideTotal
//		g.Date = stm + " - " + etm
//
//		gameAmounts = append(gameAmounts, g)
//
//	}
//
//	return
//}
//
//// 更改order表中的cp对账状态
//func SetOrderChanVerified(games []interface{}, cp, stime, etime string, status int) (affect int, err error) {
//	o := orm.NewOrm()
//	order := make(orm.Params)
//	order["ChannelVerified"] = status
//	affectstep, err := o.QueryTable(new(Order)).
//		Filter("Cp", cp).
//		Filter("Date__gte", stime).
//		Filter("Date__lte", etime).
//		Filter("GameID__in", games...).
//		Update(order)
//	affect = int(affectstep)
//	return
//}
//
////检查对账单是否存在
//func CheckChanOrderVerified(cha *ChannelVerifyAccount) bool {
//	o := orm.NewOrm()
//	return o.QueryTable(new(ChannelVerifyAccount)).
//		Filter("BodyMy", cha.BodyMy).
//		Filter("Cp", cha.Cp).
//		Filter("StartTime__gte", cha.StartTime).
//		Filter("EndTime__lte", cha.EndTime).
//		Filter("GameStr", cha.GameStr).Exist()
//}
//
//// 给对账单添加游戏信息
////func ChanVerifyAddGameInfo(ss *[]ChannelVerifyAccount) (err error) {
////	if ss == nil || len(*ss) == 0 {
////		return
////	}
////	o := orm.NewOrm()
////	gameIds := make([][]GameAmount, len(*ss))
////	var gameSum []interface{}
////	for i, s := range *ss {
////		err = json.Unmarshal([]byte(s.GameStr), &gameIds[i])
////		if err != nil {
////			return
////		}
////		for _, value := range gameIds[i] {
////			gameSum = append(gameSum, value.GameId)
////		}
////	}
////	log.Printf("gameSun: %v", gameIds)
////	//return
////	games := []GameAll{}
////	_, err = o.QueryTable(new(GameAll)).All(&games, "Id", "GameId", "Name")
////	if err != nil {
////		return
////	}
////	gameMap := map[int]GameAll{}
////	for _, g := range games {
////		gameMap[g.GameId] = g
////	}
////
////	for i := range *ss {
////		for _, gid := range gameIds[i] {
////			if g, ok := gameMap[gid.GameId]; ok {
////				temp := GameAmount{GameAll: g,
////					AmountMy:            gid.AmountMy,
////					AmountOpposite:      gid.AmountOpposite,
////					AmountPayable:       gid.AmountPayable,
////				}
////				(*ss)[i].Games = append((*ss)[i].Games, &temp)
////			}
////		}
////		(*ss)[i].GameNum = len((*ss)[i].Games)
////
////	}
////	return
////}
//
//// 给对账单添加用户信息
//func AddUserInfo4ChanVerify(ss *[]ChannelVerifyAccount) {
//	if ss == nil || len(*ss) == 0 {
//		return
//	}
//	o := orm.NewOrm()
//	userids := make([]int, 3*len(*ss))
//	for _, s := range *ss {
//		userids = append(userids, s.CreateUserId, s.UpdateUserId, s.VerifyUserId)
//	}
//	users := []User{}
//	_, err := o.QueryTable(new(User)).Filter("Id__in", userids).All(&users, "Id", "Nickname", "Name", "Email")
//	if err != nil {
//		return
//	}
//	userMap := map[int]User{}
//	for _, user := range users {
//		userMap[user.Id] = user
//	}
//	for i, s := range *ss {
//		if g, ok := userMap[s.CreateUserId]; ok {
//			(*ss)[i].CreateUser = &g
//		}
//		if g, ok := userMap[s.UpdateUserId]; ok {
//			(*ss)[i].UpdateUser = &g
//		}
//		if g, ok := userMap[s.VerifyUserId]; ok {
//			(*ss)[i].VerifyUser = &g
//		}
//	}
//}
//
//// 给渠道对账单添加渠道信息
//func AddCpInfo4ChanVerify(ss *[]ChannelVerifyAccount) {
//	if ss == nil || len(*ss) == 0 {
//		return
//	}
//	o := orm.NewOrm()
//	cps := make([]string, len(*ss))
//	for i, s := range *ss {
//		cps[i] = s.Cp
//	}
//	channs := []Channel{}
//	_, err := o.QueryTable(new(Channel)).Filter("Cp__in", cps).All(&channs, "Id", "Channelid", "Name", "Cp")
//	if err != nil {
//		return
//	}
//	linkMap := map[string]Channel{}
//	for _, g := range channs {
//		linkMap[g.Cp] = g
//	}
//	for i, s := range *ss {
//		if g, ok := linkMap[s.Cp]; ok {
//			(*ss)[i].Channel = &g
//		}
//	}
//	return
//}
//
//// 给渠道对账单添加回款主体信息
//func AddRemitCompany4ChanVerify(ss *[]ChannelVerifyAccount) {
//	if ss == nil || len(*ss) == 0 {
//		return
//	}
//	o := orm.NewOrm()
//	cps := make([]int, len(*ss))
//	for i, s := range *ss {
//		cps[i] = s.RemitCompanyId
//	}
//	linkers := []Company{}
//	_, err := o.QueryTable(new(Company)).Filter("Id__in", cps).All(&linkers, "Id", "Name")
//	if err != nil {
//		return
//	}
//	linkMap := map[int]Company{}
//	for _, g := range linkers {
//		linkMap[g.Id] = g
//	}
//	for i, s := range *ss {
//		if g, ok := linkMap[s.RemitCompanyId]; ok {
//			(*ss)[i].RemitCompany = &g
//		}
//	}
//	return
//}
//
//// 更新对账单
//func UpdateChanVerifyAccount(m *ChannelVerifyAccount) (err error) {
//	m.UpdateTime = int(time.Now().Unix())
//	o := orm.NewOrm()
//	f := utils.GetNotEmptyFields(m, "Status", "VerifyTime", "GameStr",
//		"TicketUrl", "Desc", "VerifyUserId", "UpdateUserId", "UpdateTime", "AmountPayable", "AmountOpposite", "GameStr", "FileId", "FilePreviewId",
//		"RemitCompanyId")
//	old := ChannelVerifyAccount{Id: m.Id}
//	err = o.Read(&old)
//	if err != nil {
//		return
//	}
//
//	num, err := o.Update(m, f...)
//	if err != nil {
//		return
//	}
//	if num == 0 {
//		err = errors.New("not found or update")
//	}
//	return
//}
//
//// AddChannelVerifyAccount insert a new ChannelVerifyAccount into database and returns
//// last inserted Id on success.
//func AddChannelVerifyAccount(m *ChannelVerifyAccount) (id int64, err error) {
//
//	m.CreateTime = int(time.Now().Unix())
//	m.UpdateTime = int(time.Now().Unix())
//
//	o := orm.NewOrm()
//	id, err = o.Insert(m)
//	return
//}
//
//// GetChannelVerifyAccountById retrieves ChannelVerifyAccount by Id. Returns error if
//// Id doesn't exist
//func GetChannelVerifyAccountById(id int) (v *ChannelVerifyAccount, err error) {
//	o := orm.NewOrm()
//	v = &ChannelVerifyAccount{Id: id}
//	if err = o.Read(v); err == nil {
//		return v, nil
//	}
//	return nil, err
//}
//
//// UpdateChannelVerifyAccount updates ChannelVerifyAccount by Id and returns error if
//// the record to be updated doesn't exist
//func UpdateChannelVerifyAccountById(m *ChannelVerifyAccount) (err error) {
//	o := orm.NewOrm()
//	v := ChannelVerifyAccount{Id: m.Id}
//	// ascertain id exists in the database
//	if err = o.Read(&v); err == nil {
//		var num int64
//		if num, err = o.Update(m); err == nil {
//			fmt.Println("Number of records updated in database:", num)
//		}
//	}
//	return
//}
//
//// 删除渠道对账单,
//// 需要回滚order 与 remit_down
//func DeleteChannelVerifyAccount(id int) (err error) {
//	o := orm.NewOrm()
//	oldV := ChannelVerifyAccount{Id: id}
//	if err = o.Read(&oldV); err != nil {
//		return
//	}
//	// 回滚回款金额
//	if oldV.AmountRemit != 0 {
//		err = AddRemitPreAmount(oldV.RemitCompanyId, oldV.AmountRemit)
//		if err != nil {
//			return
//		}
//	}
//
//	// 回滚order标记
//	oldGames := []GameAmount{}
//	err = json.Unmarshal([]byte(oldV.GameStr), &oldGames)
//	if err != nil {
//		return
//	}
//	oldGameIds := []interface{}{}
//	for _, v := range oldGames {
//		oldGameIds = append(oldGameIds, v.GameId)
//	}
//	_, err = SetOrderChanVerified(oldGameIds, oldV.Cp, oldV.StartTime, oldV.EndTime, 2)
//	if err != nil {
//		return
//	}
//	// end
//
//	if _, err = o.Delete(&ChannelVerifyAccount{Id: id}); err != nil {
//		return
//	}
//
//	return
//}
//
//// 获取没有回款的对账单
//func GetNotRemitAccount(companyIds []interface{}, startTime int, endTime int) (s *bean.NotRemit, err error) {
//	where := ""
//	args := []interface{}{}
//	if l := len(companyIds); l != 0 {
//		holder := strings.Repeat(",?", l)
//		where = where + fmt.Sprintf(" `remit_company_id` in (%s) AND ", holder[1:])
//		args = append(args, companyIds...)
//	}
//	if startTime != 0 {
//		sT := time.Unix(int64(startTime), 0).Format("2006-01-02")
//		where = where + " `start_time` >= ? AND"
//		args = append(args, sT)
//	}
//	if endTime != 0 {
//		sT := time.Unix(int64(startTime), 0).Format("2006-01-02")
//		where = where + " `end_time` <= ? AND"
//		args = append(args, sT)
//	}
//	finishStatus := CHAN_VERIFY_S_RECEIPT
//	sql := fmt.Sprintf("SELECT remit_company_id, max(end_time) as end_time,min(start_time) as start_time, SUM(amount_payable - amount_remit) as nots from channel_verify_account WHERE %s amount_payable != amount_remit AND status = %d GROUP By remit_company_id", where, finishStatus)
//	maps := []orm.Params{}
//	o := orm.NewOrm()
//	_, err = o.Raw(sql, args...).Values(&maps)
//	if err != nil {
//		return
//	}
//	if len(maps) == 0 {
//		s = &bean.NotRemit{
//			Amount:      0,
//			EndTime:     int64(endTime),
//			CompanyName: "共0个主体",
//			StartTime:   int64(startTime),
//		}
//		return
//	}
//
//	allCompanyIds := []interface{}{}
//	end_time := ""
//	start_time := ""
//	var allAmount float64 = 0
//	for _, v := range maps {
//		if e := v["end_time"].(string); e != "" && e[0] != '0' && (end_time == "" || end_time < e) {
//			end_time = e
//		}
//		if s := v["start_time"].(string); s != "" && s[0] != '0' && (start_time == "" || start_time > s) {
//			start_time = s
//		}
//
//		nots := v["nots"].(string)
//		companyId := v["remit_company_id"]
//
//		f, e := strconv.ParseFloat(nots, 64)
//		if e != nil {
//			err = e
//			return
//		}
//		allCompanyIds = append(allCompanyIds, companyId)
//		allAmount += f
//	}
//
//	var rStartTime int64 = 0
//	var rEndTime int64 = 0
//	st, er := time.ParseInLocation("2006-01-02", start_time, time.Local)
//	if er == nil {
//		rStartTime = st.Unix()
//	}
//	et, er := time.ParseInLocation("2006-01-02", end_time, time.Local)
//	if er == nil {
//		rEndTime = et.Unix()
//	}
//
//	s = &bean.NotRemit{
//		Amount:      allAmount,
//		StartTime:   rStartTime,
//		EndTime:     rEndTime,
//		CompanyName: fmt.Sprintf("共%d个主体", len(allCompanyIds)),
//	}
//	return
//}
//
//
////get all channel verify account
//func QueryChannelVerifyAccount() ([]orm.Params, error) {
//
//	o := orm.NewOrm()
//	//var id int
//	//accounts := []*ChannelVerifyAccount{}
//	var accounts []orm.Params
//	qs := o.QueryTable("channel_verify_account")
//	_, err := qs.Limit(		0, 0).Values(&accounts, "id", "cp", "start_time", "end_time", "games")
//	if err != nil {
//		beego.Debug(err)
//	}
//	return accounts, err
//}
