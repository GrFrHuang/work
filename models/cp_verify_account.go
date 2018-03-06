package models
//
//import (
//	"errors"
//	"fmt"
//	"github.com/astaxie/beego/orm"
//	"github.com/bysir-zl/bygo/util"
//	"kuaifa.com/kuaifa/work-together/utils"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//	"kuaifa.com/kuaifa/work-together/tool/old_sys"
//	"encoding/json"
//)
//
//type CpVerifyAccount struct {
//	Id             int     `json:"id,omitempty" orm:"column(id);auto"`
//	CompanyId      int `json:"company_id,omitempty" orm:"column(company_id);null"` // 发行商id
//	Games          string `json:"games,omitempty" orm:"column(games);null"`        // 游戏列表json {GameAmount}
//	Status         int    `json:"status,omitempty" orm:"column(status);null" `
//	StartTime      string    `json:"start_time,omitempty" orm:"column(start_time);null"`
//	EndTime        string    `json:"end_time,omitempty" orm:"column(end_time);null" `
//	AmountMy       float64 `json:"amount_my,omitempty" orm:"column(amount_my);null;digits(16);decimals(2)"`             // 我方流水
//	AmountOpposite float64 `json:"amount_opposite,omitempty" orm:"column(amount_opposite);null;digits(16);decimals(2)"` // 对方流水
//	AmountTheory   float64 `json:"amount_theory,omitempty" orm:"column(amount_theory);null;digits(16);decimals(2)"`     // 理论金额
//	AmountPayable  float64 `json:"amount_payable,omitempty" orm:"column(amount_payable);null;digits(16);decimals(2)"`   // 应收款
//	AmountSettle   float64 `json:"amount_settle,omitempty" orm:"column(amount_settle);null;digits(16);decimals(2)"`     // 已收款
//	VerifyTime     int    `json:"verify_time,omitempty" orm:"column(verify_time);null"`
//	CreateUserId   int    `json:"create_user_id,omitempty" orm:"column(create_user_id);null"`
//	VerifyUserId   int    `json:"verify_user_id,omitempty" orm:"column(verify_user_id);null"`
//	UpdateUserId   int    `json:"update_user_id,omitempty" orm:"column(update_user_id);null"`
//	CreateTime     int    `json:"create_time,omitempty" orm:"column(create_time);null"`
//	UpdateTime     int    `json:"update_time,omitempty" orm:"column(update_time);null"`
//	TicketUrl      string `json:"ticket_url,omitempty" orm:"column(ticket_url);null"`
//	Extra          string  `json:"extra,omitempty" orm:"column(extra);size(255);null"`
//	Desc           string  `json:"desc,omitempty" orm:"column(desc);null"`
//	FileId         int     `json:"file_id,omitempty" orm:"column(file_id);null"`
//	FilePreviewId  int     `json:"file_preview_id,omitempty" orm:"column(file_preview_id);"`
//	BodyMy         int     `json:"body_my,omitempty" orm:"column(body_my);"`
//
//	//UpdateCheckTime     int64  `json:"update_checktime,omitempty" orm:"column(update_checktime);null" `
//	//UpdateCheckUser       string `json:"update_checkuser,omitempty" orm:"column(update_checkuser);size(255)"`
//
//	Company    *Company `orm:"-" json:"company,omitempty"`
//	CreateUser *User `orm:"-" json:"create_user,omitempty"`
//	VerifyUser *User `orm:"-" json:"verify_user,omitempty"`
//	UpdateUser *User `orm:"-" json:"update_user,omitempty"`
//}
//
////对账单的状态
//const (
//	CP_VERIFY_S_VERIFYING = 10  // 对账中
//	CP_VERIFY_S_VERIFYIED = 20  // 对账完毕,未开票
//	CP_VERIFY_S_RECEIPT   = 30  // 已收票
//	CP_VERIFY_S_SETTLE    = 100 // 已打款
//)
//
//var CpStatusCode2String = map[int]string{
//	CP_VERIFY_S_VERIFYING: "对账中",
//	CP_VERIFY_S_VERIFYIED: "对账完毕,未开票",
//	CP_VERIFY_S_RECEIPT:   "已收票",
//
//	CP_VERIFY_S_SETTLE: "已打款", // 此状态在实际=已结算金额自动生效, 不能由前端修改
//}
//
//func (t *CpVerifyAccount) TableName() string {
//	return "cp_verify_account"
//}
//
//func init() {
//	orm.RegisterModel(new(CpVerifyAccount))
//}
//
//// 通过公司获取未对账的全部流水
//func SumCpNotVerifyByCompany(sTime, eTime string, companyIds []interface{}) (res CpVerifyAccount, err error) {
//	_, err = time.Parse("2006-01-02", sTime)
//	if err != nil {
//		return
//	}
//	_, err = time.Parse("2006-01-02", eTime)
//	if err != nil {
//		return
//	}
//	o := orm.NewOrm()
//
//	var sql string
//	var raw orm.RawSeter
//	var maps []orm.Params
//
//	if len(companyIds) == 0 {
//		sql = "SELECT SUM(amount) total,MAX(date) last_time,MIN(date) first_time FROM `order` WHERE cp_verified != 1 AND date >= ? AND date <= ?"
//		raw = o.Raw(sql, sTime, eTime)
//	} else {
//		whereIn := strings.Repeat(",?", len(companyIds))[1:]
//		sql = fmt.Sprintf("SELECT SUM(amount) total,MAX(date) last_time,MIN(date) first_time,game.game_id,game.game_name FROM `order` LEFT JOIN game ON game.game_id=`order`.game_id WHERE game.issue in (%s) AND cp_verified != 1 AND date >= ? AND date <= ? ", whereIn)
//		arg := append(companyIds, sTime, eTime)
//		raw = o.Raw(sql, arg...)
//	}
//	_, err = raw.Values(&maps)
//	if err != nil || maps == nil || len(maps) == 0 {
//		return
//	}
//	v := maps[0]
//	if v["total"] == nil {
//		res.AmountMy = 0
//		res.StartTime = sTime
//		res.EndTime = eTime
//		res.Company = &Company{Name: "无"}
//		return
//	}
//
//	sum, _ := util.Interface2Float(v["total"], false)
//
//	res.AmountMy = sum
//	res.StartTime = sTime
//	res.EndTime = eTime
//	res.Company = &Company{Name: "*"}
//
//	return
//}
//
//// unused on v2
//// 由于现在cp对账是以发行商来
//// 获取多个游戏未对账的数据
////func SumNoCpVerifyByGame(stime, etime string, gameids []interface{}) (res *CpVerifyAccount, err error) {
////	_, err = time.Parse("2006-01-02", stime)
////	if err != nil {
////		return
////	}
////	_, err = time.Parse("2006-01-02", etime)
////	if err != nil {
////		return
////	}
////
////	//var gamejson []byte
////	sql := "SELECT SUM(amount) as total FROM `order` WHERE date >= ? AND date <= ? AND cp_verified != 1 "
////	if l := len(gameids); l > 0 {
////		create := strings.Repeat(",?", l)[1:]
////		sql = sql + fmt.Sprintf("AND game_id IN (%s) ", create)
////	}
////	temp := []interface{}{
////		stime,
////		etime,
////	}
////	temp = append(temp, gameids...)
////	var step []orm.Params
////	o := orm.NewOrm()
////	l, err := o.Raw(sql, temp...).Values(&step)
////	if err != nil {
////		return
////	}
////	if l == 0 {
////		res = &CpVerifyAccount{
////			AmountMy:   0,
////			StartTime:  stime,
////			EndTime:    etime,
////			Status:     0,
////			Games:      `无`,
////		}
////		return
////	}
////
////	var sum float64
////	if step[0]["total"] != nil {
////		sum, err = strconv.ParseFloat(step[0]["total"].(string), 64)
////	} else {
////		sum = float64(0)
////	}
////
////	res = &CpVerifyAccount{
////		AmountMy:   sum,
////		StartTime:  stime,
////		EndTime:    etime,
////		Status:     0,
////		Games:      `*`,
////	}
////	return
////}
//
////func CpVerifyAddGameInfo(ss *[]CpVerifyAccount) {
////	if ss == nil || len(*ss) == 0 {
////		return
////	}
////	o := orm.NewOrm()
////	gameIds := make([]int, len(*ss))
////	for i, s := range *ss {
////		gameIds[i] = s.GameId
////	}
////	games := []Game{}
////	_, err := o.QueryTable(new(Game)).Filter("GameId__in", gameIds).All(&games, "Id", "GameId", "GameName")
////	if err != nil {
////		return
////	}
////	gameMap := map[int]Game{}
////	for _, g := range games {
////		gameMap[g.GameId] = g
////	}
////	for i, s := range *ss {
////		if g, ok := gameMap[s.GameId]; ok {
////			(*ss)[i].Game = &g
////		}
////	}
////	return
////}
//
//func AddUserInfo4CpVerify(ss *[]CpVerifyAccount) {
//	if ss == nil || len(*ss) == 0 {
//		return
//	}
//	o := orm.NewOrm()
//
//	userids := make([]int, 3*len(*ss))
//	for _, s := range *ss {
//		userids = append(userids, s.CreateUserId, s.UpdateUserId, s.VerifyUserId)
//	}
//	users := []User{}
//	_, err := o.QueryTable(new(User)).Filter("Id__in", userids).All(&users, "Id", "Nickname", "Name");
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
//func AddCompanyInfo4CpVerify(ss *[]CpVerifyAccount) {
//	if ss == nil || len(*ss) == 0 {
//		return
//	}
//
//	o := orm.NewOrm()
//	linkIds := make([]int, len(*ss))
//	for i, s := range *ss {
//		linkIds[i] = s.CompanyId
//	}
//	link := []Company{}
//	_, err := o.QueryTable(new(Company)).
//		Filter("Id__in", linkIds).
//		All(&link, "Id", "Name")
//	if err != nil {
//		return
//	}
//	linkMap := map[int]Company{}
//	for _, g := range link {
//		linkMap[g.Id] = g
//	}
//	for i, s := range *ss {
//		if g, ok := linkMap[s.CompanyId]; ok {
//			(*ss)[i].Company = &g
//		}
//	}
//	return
//}
//
//// 获取未对账的发行商
//func GetNotVerifyCompanies() (companies []Company, err error) {
//	companies = []Company{}
//	sql := "SELECT company.`name` , company.`id` FROM `order` " +
//		"RIGHT JOIN `game` ON order.game_id=game.game_id " +
//		"LEFT JOIN company ON company.id = issue " +
//		"WHERE cp_verified != 1 " +
//		"GROUP BY game.issue"
//	values := []orm.Params{}
//	_, err = orm.NewOrm().Raw(sql).Values(&values)
//	if err != nil {
//		return
//	}
//	for _, v := range values {
//		companyIdStr, _ := util.Interface2String(v["id"], true)
//		companyName, _ := util.Interface2String(v["name"], true)
//		companyId, _ := strconv.Atoi(companyIdStr)
//		comp := Company{Id: companyId, Name: companyName}
//		companies = append(companies, comp)
//	}
//
//	return
//}
//
//// 获取某发行商下的未对账时间
//func GetNotVerifyDateByComp(companyId int) (cpV []CpVerifyAccount, err error) {
//	cpV = []CpVerifyAccount{}
//	sql := "SELECT LEFT (date, 7) AS month,MIN(date) AS mindate FROM `order` LEFT JOIN game ON game.game_id=`order`.game_id WHERE game.issue = ? AND cp_verified!=1 GROUP BY month"
//	values := []orm.Params{}
//	_, err = orm.NewOrm().Raw(sql, companyId).Values(&values)
//	if err != nil {
//		return
//	}
//
//	for _, v := range values {
//		month := v["month"].(string)
//		start := v["mindate"].(string)
//		tm, e := time.Parse("2006-01", month)
//		if e != nil {
//			err = e
//			return
//		}
//
//		comp := CpVerifyAccount{
//			StartTime: start,
//			EndTime:   tm.AddDate(0, 1, -1).Format("2006-01-02"),
//		}
//		cpV = append(cpV, comp)
//	}
//
//	return
//}
//
//// 获取发行商和时间范围内未对账的游戏
//func GetCpNotVerifyGames(stm, etm string, companyId int, myBody int, where map[string][]interface{}) (gameAmounts []GameAmount, err error) {
//	gameAmounts = []GameAmount{}
//	_, err = time.Parse("2006-01-02", stm)
//	if err != nil {
//		return
//	}
//	_, err = time.Parse("2006-01-02", etm)
//	if err != nil {
//		return
//	}
//
//	var maps []orm.Params
//	o := orm.NewOrm()
//	// SELECT SUM(amount) total,MAX(date) last_time,MIN(date) first_time,game.game_id,game.game_name FROM `order` LEFT JOIN game ON game.game_id=`order`.game_id WHERE game.issue=1 AND cp_verified != 1 AND date >= "2016-01-01" AND date <= "2018-01-01" GROUP BY game_id
//	sql := "SELECT MAX(date) last_time,MIN(date) first_time,game.game_id,game.game_name " +
//		"FROM `order` LEFT JOIN game ON game.game_id=`order`.game_id " +
//		"WHERE game.issue=? AND `order`.cp_verified != 1 AND `order`.`date` >= ? AND `order`.`date` <= ? AND " +
//		"game.game_id IN (SELECT DISTINCT game_id FROM `contract` WHERE company_type = 0 AND company_id = ? AND body_my = ?)" +
//		"GROUP BY game_id"
//
//	_, err = o.Raw(sql, companyId, stm, etm, companyId, myBody).
//		Values(&maps)
//	if err != nil || maps == nil || len(maps) == 0 {
//		return
//	}
//
//	//fmt.Println(maps)
//	var games []string
//	for _, game := range maps {
//		games = append(games, game["game_id"].(string))
//	}
//
//	clears, _ := old_sys.GetAllClearing(strings.Join(games, ","), "", stm[0:7])
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
//		g.GameName = value["game_name"].(string)
//		g.GameId = id
//		g.AmountMy = clears_map[value["game_id"].(string)].Total
//		//clearing, _ := old_sys.GetClearing(id, "", stm[0:7])
//		g.AmountTheory = clears_map[value["game_id"].(string)].DivideTotal
//		g.Date = value["first_time"].(string) + " - " + value["last_time"].(string)
//
//		gameAmounts = append(gameAmounts, g)
//	}
//	return
//}
//
////// 根据游戏id获取该游戏需要对账的时间
////func GetVerifyDateByGameId(gameid int) (res []CpVerifyAccount, err error) {
////	o := orm.NewOrm()
////	var maps []orm.Params
////	_, err = o.Raw("SELECT LEFT (date, 7) AS months,MIN(date) AS mindate FROM "+"`order` WHERE game_id = ? AND cp_verified!=1 GROUP BY months", gameid).Values(&maps)
////	res = []CpVerifyAccount{}
////	if len(maps) == 0 {
////		return
////	}
////	if err != nil {
////		return
////	}
////	cp := CpVerifyAccount{}
////	var tm time.Time
////	if len(maps) != 0 {
////		for _, value := range maps {
////			start := value["mindate"].(string)
////			tm, err = time.Parse("2006-01", value["months"].(string));
////			if err != nil {
////				return
////			}
////			cp.GameId = gameid
////			cp.StartTime = start
////			cp.EndTime = tm.AddDate(0, 1, -1).Format("2006-01-02")
////			res = append(res, cp)
////		}
////	}
////	return
////}
//
////获取结算部人员
//func GetCpVerifyUsers() (maps []orm.Params, err error) {
//	o := orm.NewOrm()
//	//var maps []orm.Params
//	_, err = o.QueryTable(new(User)).Filter("DepartmentId", "6").Values(&maps, "Nickname", "Id", "Name")
//	if err != nil {
//		return
//	}
//	return
//}
//
//// 更改order表中的多个游戏的cp对账状态
//func SetOrderCpVerified(gameIds []interface{}, stime, etime string, status int) (affect int, err error) {
//	o := orm.NewOrm()
//	order := make(orm.Params)
//	order["CpVerified"] = status
//	affectstep, err := o.QueryTable(new(Order)).
//		Filter("GameID__in", gameIds).
//		Filter("Date__gte", stime).
//		Filter("Date__lte", etime).
//	//Filter("CpVerified__in", 0, 2).
//		Update(order)
//	affect = int(affectstep)
//	return
//}
//
//var lock_fasdghjklxchjkvbukdefg sync.Mutex
//// 检查对账单是否存在
//func CheckCpOrderVerified(cp *CpVerifyAccount) (bool) {
//	lock_fasdghjklxchjkvbukdefg.Lock()
//	defer lock_fasdghjklxchjkvbukdefg.Unlock()
//	o := orm.NewOrm()
//	return o.QueryTable(new(CpVerifyAccount)).
//		Filter("BodyMy", cp.BodyMy).
//		Filter("CompanyId", cp.CompanyId).
//		Filter("StartTime__gte", cp.StartTime).
//		Filter("EndTime__lte", cp.EndTime).
//		Exist()
//}
//
//// 更新对账单
//func UpdateCpVerifyAccount(m *CpVerifyAccount) (err error) {
//	o := orm.NewOrm()
//	old := CpVerifyAccount{Id: m.Id}
//	err = o.Read(&old)
//	if err != nil {
//		return
//	}
//	//if old.Status == CP_VERIFY_S_RECEIPT || old.Status == CP_VERIFY_S_SETTLE {
//	//	err = errors.New("can't modify when account is finished")
//	//	return
//	//}
//
//	// v1:
//	// 不能修改金额与游戏,因为涉及到连表(order表)操作
//	// 2017-03-11 : 现在需要修改金额与游戏:
//	// 回滚已经标记的游戏(order表中),然后重新标记
//	f := utils.GetNotEmptyFields(m, "Status", "VerifyTime", "Games",
//		"FileId", "FilePreviewId", "Desc",
//		"VerifyUserId", "UpdateUserId", "UpdateTime", )
//
//	num, err := o.Update(m, f...)
//	if err != nil {
//		return
//	}
//	if num == 0 {
//		err = errors.New("not found")
//	}
//	return
//}
//
//// AddCpVerifyAccount insert a new CpVerifyAccount into database and returns
//// last inserted Id on success.
//func AddCpVerifyAccount(m *CpVerifyAccount) (id int64, err error) {
//	m.CreateTime = int(time.Now().Unix())
//	m.UpdateTime = int(time.Now().Unix())
//
//	o := orm.NewOrm()
//	id, err = o.Insert(m)
//	return
//}
//
//// GetCpVerifyAccountById retrieves CpVerifyAccount by Id. Returns error if
//// Id doesn't exist
//func GetCpVerifyAccountById(id int) (v *CpVerifyAccount, err error) {
//	o := orm.NewOrm()
//	v = &CpVerifyAccount{Id: id}
//	if err = o.Read(v); err == nil {
//		return v, nil
//	}
//	return nil, err
//}
//
//// UpdateCpVerifyAccount updates CpVerifyAccount by Id and returns error if
//// the record to be updated doesn't exist
//func UpdateCpVerifyAccountById(m *CpVerifyAccount) (err error) {
//	o := orm.NewOrm()
//	v := CpVerifyAccount{Id: m.Id}
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
//// 删除cp对账单,
//// 需要回滚
//func DeleteCpVerifyAccount(id int) (err error) {
//	o := orm.NewOrm()
//	oldV := CpVerifyAccount{Id: id}
//	// ascertain id exists in the database
//	if err = o.Read(&oldV); err != nil {
//		return
//	}
//
//	// start
//	// 回滚回款金额
//	if oldV.AmountSettle != 0 {
//		err = AddSettlePreAmount(oldV.CompanyId, oldV.AmountSettle)
//		if err != nil {
//			return
//		}
//	}
//
//	// 回滚order标记
//	oldGames := []GameAmount{}
//	err = json.Unmarshal([]byte(oldV.Games), &oldGames)
//	if err != nil {
//		return
//	}
//	oldGameIds := []interface{}{}
//	for _, v := range oldGames {
//		oldGameIds = append(oldGameIds, v.GameId)
//	}
//
//	_, err = SetOrderCpVerified(oldGameIds, oldV.StartTime, oldV.EndTime, 2)
//	if err != nil {
//		return
//	}
//	// end
//
//	if _, err = o.Delete(&CpVerifyAccount{Id: id}); err != nil {
//		return
//	}
//
//	return
//}
//
//// 获取没有结算的对账单
//func GetNotSettleAccount(gameIds []interface{}, startTime int, endTime int) (s *NotSettled, err error) {
//	where := ""
//	args := []interface{}{}
//	if l := len(gameIds); l != 0 {
//		holder := strings.Repeat(",?", l)
//		where = where + fmt.Sprintf(" `company_id` in (%s) AND ", holder[1:])
//		args = append(args, gameIds...)
//	}
//	if startTime != 0 {
//		sT := time.Unix(int64(startTime), 0).Format("2006-01-02")
//		where = where + " `start_time` >= ? AND"
//		args = append(args, sT)
//	}
//	if endTime != 0 {
//		sT := time.Unix(int64(endTime), 0).Format("2006-01-02")
//		where = where + " `end_time` <= ? AND"
//		args = append(args, sT)
//	}
//
//	finishStatus := CP_VERIFY_S_RECEIPT
//	sql := fmt.Sprintf("SELECT company_id, max(end_time) as end_time,min(start_time) as start_time, SUM(amount_payable - amount_settle) as nots from cp_verify_account WHERE %s amount_payable != amount_settle AND status = %d GROUP By company_id ", where, finishStatus)
//	maps := []orm.Params{}
//	o := orm.NewOrm()
//	_, err = o.Raw(sql, args...).Values(&maps)
//	if err != nil {
//		return
//	}
//	if len(maps) == 0 {
//		s = &NotSettled{
//			Amount:    0,
//			EndTime:   int64(endTime),
//			StartTime: int64(startTime),
//		}
//		return
//	}
//
//	allCompanyIds := []interface{}{}
//	end_time := ""
//	start_time := ""
//	var allNotSettle float64 = 0
//	for _, v := range maps {
//		if e := v["end_time"].(string);
//			e != "" && e[0] != '0' && (end_time == "" || end_time < e ) {
//			end_time = e
//		}
//		if s := v["start_time"].(string);
//			s != "" && s[0] != '0' && (start_time == "" || start_time > s ) {
//			start_time = s
//		}
//
//		nots := v["nots"].(string)
//		company_id := v["company_id"]
//
//		f, e := strconv.ParseFloat(nots, 64)
//		if e != nil {
//			err = e
//			return
//		}
//		allCompanyIds = append(allCompanyIds, company_id)
//		allNotSettle += f
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
//	s = &NotSettled{
//		Amount:    allNotSettle,
//		StartTime: rStartTime,
//		EndTime:   rEndTime,
//		Companies: []Company{
//			{
//				Name: fmt.Sprintf("共%d个发行商", len(allCompanyIds)),
//			},
//		},
//	}
//	return
//}
