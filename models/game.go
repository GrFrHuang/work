package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/utils"
	"strconv"
	"strings"
	"time"
	"errors"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

type Game struct {
	Id                int    `json:"id,omitempty" orm:"column(id);auto"`
	GameId            int    `json:"game_id,omitempty" orm:"column(game_id)"`                                 //游戏ID
	GameName          string `json:"game_name,omitempty" orm:"column(game_name);size(255)"`                   //游戏名
	Source            int    `json:"source" orm:"column(source)"`                                             // 游戏来源
	ImportTime        int64  `json:"import_time,omitempty" orm:"column(import_time);null" `                   //接入时间
	PublishTime       int64  `json:"publish_time,omitempty" orm:"column(publish_time);null" `                 //发行时间
	GameType          int    `json:"game_type,omitempty" orm:"column(game_type);null" `                       //游戏类型
	Cooperation       int    `json:"cooperation,omitempty" orm:"column(cooperation);null" `                   //合作方式
	Development       int    `json:"development,omitempty" orm:"column(development);null" `                   //研发商
	Issue             int    `json:"issue,omitempty" orm:"column(issue);null" `                               //发行商
	Ip                int    `json:"ip,omitempty" orm:"column(ip);null" `                                     //是否有ip
	Star              string `json:"star,omitempty" orm:"column(star);null" `                                 //是否有明星代言
	Budget            string `json:"budget,omitempty" orm:"column(budget);null" `                             //是否有市场预算
	QuantityPolicy    string `json:"quantity_policy,omitempty" orm:"column(quantity_policy);size(255);null" ` //保量政策
	SociatyPolicy     int    `json:"sociaty_policy,omitempty" orm:"column(sociaty_policy);null" `             //公会政策
	AccessPerson      int    `json:"access_person,omitempty" orm:"column(access_person);null" `               //接入人
	Advise            string `json:"advise,omitempty" orm:"column(advise);null" `                             //接入建议:[id:1,id:2],1,2分别表示接入或不接入
	Result            int    `json:"result,omitempty" orm:"column(result);null"`
	ResultPerson      string `json:"result_person,omitempty" orm:"column(result_person);null"`
	ResultTime        int64  `json:"result_time,omitempty" orm:"column(result_time);null"`
	ResultReportWord  int    `json:"result_report_word,omitempty" orm:"column(result_report_word);null"`
	ResultReportExcel int    `json:"result_report_excel,omitempty" orm:"column(result_report_excel);null"`
	Package           string `json:"package,omitempty" orm:"column(package);null" `           //测试包
	Picture           string `json:"picture,omitempty" orm:"column(picture);null" `           //测试数据截图
	Remarks           string `json:"remarks,omitempty" orm:"column(remarks);size(255);null" ` //
	Ladders           string `json:"ladders,omitempty" orm:"column(ladders);null" `           // 阶梯
	Number            string `json:"number,omitempty" orm:"column(number);size(255);null" `   //
	Switch            int8   `json:"switch,omitempty" orm:"column(switch);null" `             //开关，默认为1，表示接入
	CreateTime        int64  `json:"create_time,omitempty" orm:"column(create_time);null" `
	//UpdateTime     int64  `json:"update_time,omitempty" orm:"column(update_time);null" `
	SubmitPerson int `json:"submit_person,omitempty" orm:"column(submit_person);null" `

	BodyMy int `json:"body_my,omitempty" orm:"column(body_my);null"`
	//该字段为 标记 我方主体 1为 云端  2为 有量

	UpdateRefTime  int64 `json:"update_reftime,omitempty" orm:"column(update_reftime);null" `
	UpdateEvalTime int64 `json:"update_evaltime,omitempty" orm:"column(update_evaltime);null" `
	UpdateJrTime   int64 `json:"update_jrtime,omitempty" orm:"column(update_jrtime);null" `
	//以上3个字段分别为 游戏提测、测评、接入更新时间

	UpdateJrUserID   int `json:"update_jruserid,omitempty" orm:"column(update_jruserid)"`
	UpdateRefUserID  int `json:"update_refuserid,omitempty" orm:"column(update_refuserid)"`
	UpdateEvalUserID int `json:"update_evaluserid,omitempty" orm:"column(update_evaluserid)"`
	//以上3个字段分别是 游戏提测、测评、接入更新人

	User          *User        `orm:"-" json:"user,omitempty"`
	UpdateUser    *User        `orm:"-" json:"update_user,omitempty"` //更新人
	Leixing       *Types       `orm:"-" json:"lexing,omitempty"`
	GameSource    *Types       `orm:"-" json:"game_source,omitempty"`
	Hezuo         *Types       `orm:"-" json:"hezuo,omitempty"`
	Yanfa         *CompanyType `orm:"-" json:"yanfa,omitempty"`
	Faxing        *CompanyType `orm:"-" json:"faxing,omitempty"`
	Gonghui       *Types       `orm:"-" json:"gonghui,omitempty"`
	Ceping        *Types       `orm:"-" json:"ceping,omitempty"`        //测评结果
	Cepingpeoples *[]User      `orm:"-" json:"cepingpeoples,omitempty"` //评测人
	Advises       *[]Advise    `orm:"-" json:"advises,omitempty"`       //接入建议
	Ladder_front  interface{}  `orm:"-" json:"ladder_front"`
	//Numbers  *[]GameTestNumber  `orm:"-" json:"numbers,omitempty"`
	//Number  string  `orm:"-" json:"number,omitempty"` //用于接收前台传回的账号密码，e.g. 123456,asdf;123,asdf
	Access *AccessState `orm:"-" json:"access"` //接入状态
}

type Advise struct {
	Id       int    `json:"id,omitempty" orm:"column(id);null"`
	Name     string `json:"name,omitempty" orm:"column(name);null"`
	NickName string `json:"nickname,omitempty" orm:"column(nickname);null"`
	Adv      string `json:"adv,omitempty" orm:"column(adv);null"`
}

//游戏上线状态
type AccessState struct {
	State            int   `json:"state" orm:"column(state);null"` //1:上架，2:下架
	AccessUpdateTime int64 `json:"time"  orm:"column(time);null"`
	Count            int   `json:"-" orm:"column(count);null"` //某游戏所有渠道中上架中的数量
}

// 提测,评测统计
type Statistics struct {
	UserId   int    `json:"user_id" orm:"column(user_id);null"`
	NickName string `json:"nick_name" orm:"column(nickname);null"` // 提测、评测人
	Total    int    `json:"total" orm:"column(total);null"`        // 此人提测、评测游戏数
}

func (t *Game) TableName() string {
	return "game"
}

func init() {
	orm.RegisterModel(new(Game))
}

func GetTCStatistics() (statistic []Statistics) {
	statistic = []Statistics{}
	var query []Statistics

	us, err := GetUsersByDevMent(245)
	if err != nil {
		return
	}

	userIds := make([]string, 0)
	statMap := map[int]int{}
	for _, u := range us {
		userIds = append(userIds, strconv.Itoa(u.Id))
	}

	now := time.Now().Format("2006-01-02")
	sql := "SELECT submit_person AS user_id,COUNT(FROM_UNIXTIME(create_time,'%Y-%m-%d')) AS `total` FROM game " +
		"WHERE submit_person IN(" + strings.Join(userIds, ",") + ") AND FROM_UNIXTIME(create_time,'%Y-%m-%d') = '" + now + "' " +
		"GROUP BY submit_person"

	orm.NewOrm().Raw(sql).QueryRows(&query)

	for _, que := range query {
		statMap[que.UserId] = que.Total
	}

	for _, u := range us {
		result := Statistics{
			UserId:   u.Id,
			NickName: u.Nickname,
			Total:    statMap[u.Id],
		}
		statistic = append(statistic, result)
	}

	return statistic
}

// 接入统计
func GetJRStatistics() (statistic []Statistics) {
	statistic = []Statistics{}
	var query []Statistics

	us, err := GetUsersByDevMent(245)
	if err != nil {
		return
	}

	userIds := make([]string, 0)
	statMap := map[int]int{}
	for _, u := range us {
		userIds = append(userIds, strconv.Itoa(u.Id))
	}

	now := time.Now().Format("2006-01-02")
	sql := "SELECT update_jruserid AS user_id,COUNT(FROM_UNIXTIME(update_jrtime,'%Y-%m-%d')) AS `total` FROM game " +
		"WHERE update_jruserid IN(" + strings.Join(userIds, ",") + ") AND FROM_UNIXTIME(update_jrtime,'%Y-%m-%d') = '" + now + "' " +
		"GROUP BY update_jruserid"

	orm.NewOrm().Raw(sql).QueryRows(&query)

	for _, que := range query {
		statMap[que.UserId] = que.Total
	}

	for _, u := range us {
		result := Statistics{
			UserId:   u.Id,
			NickName: u.Nickname,
			Total:    statMap[u.Id],
		}
		statistic = append(statistic, result)
	}

	return statistic
}

// 评测统计
func GetPCStatistics() (statistic []Statistics) {
	statistic = []Statistics{}
	var query []Statistics

	us, err := GetUsersByDevMent(10)
	if err != nil {
		return
	}

	userIds := make([]string, 0)
	statMap := map[int]int{}
	for _, u := range us {
		userIds = append(userIds, strconv.Itoa(u.Id))
	}

	now := time.Now().Format("2006-01-02")
	sql := "SELECT user_id,COUNT(DISTINCT(`action`)) AS `total`  FROM `logs` WHERE " +
		"`action` LIKE 'put /game/evaluation/%' AND " +
		"user_id IN(" + strings.Join(userIds, ",") + ")  AND FROM_UNIXTIME(created_time,'%Y-%m-%d') = '" + now + "' " +
		"GROUP BY user_id,`action`"

	orm.NewOrm().Raw(sql).QueryRows(&query)

	for _, que := range query {
		statMap[que.UserId] = que.Total
	}

	for _, u := range us {
		result := Statistics{
			UserId:   u.Id,
			NickName: u.Nickname,
			Total:    statMap[u.Id],
		}
		statistic = append(statistic, result)
	}

	return statistic
}

// AddGame insert a new Game into database and returns
// last inserted Id on success.
func AddGame(m *Game) (id int64, err error) {
	o := orm.NewOrm()
	m.CreateTime = time.Now().Unix()
	//m.UpdateTime = time.Now().Unix()
	id, err = o.Insert(m)

	err = CompareAndAddOperateLog(nil, m, m.SubmitPerson, bean.OPP_GAME_REFERENCE, int(id), bean.OPA_INSERT)
	return
}

// GetGameById retrieves Game by Id. Returns error if
// Id doesn't exist
func GetGameById(id int) (v *Game, err error) {
	o := orm.NewOrm()
	v = &Game{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetGameByGameId retrieves Game by game_id. Returns error if
// Id doesn't exist
func GetGameByGameId(id int) (v *Game, err error) {
	fmt.Println(id)
	o := orm.NewOrm()
	v = &Game{GameId: id}
	if err = o.Read(v, "game_id"); err == nil {
		return v, nil
	}
	return nil, err
}

func GetGameUpdateSelect() (games []Game, err error) {
	o := orm.NewOrm()
	_, err = o.Raw("SELECT a.id,a.game_id,a.game_name FROM game a WHERE a.game_id > 0").QueryRows(&games)
	return
}

func GetResultByGameId(id int) (v Ceping) {
	o := orm.NewOrm()
	_ = o.Raw("SELECT b.name result, a.result_report_word word, a.result_report_excel excel"+" FROM game a LEFT JOIN types b ON a.result=b.id WHERE game_id=? ", id).QueryRow(&v)
	return
}

func GetAddGame() (maps []orm.Params, err error) {
	o := orm.NewOrm()
	//var maps []orm.Params

	_, err = o.Raw("SELECT game_id,game_name FROM game_all WHERE game_id NOT IN (SELECT game_id FROM game)").Values(&maps)

	return
}

func GetLatestGame() (maps []orm.Params, err error) {
	o := orm.NewOrm()
	//var maps []orm.Params

	_, err = o.Raw("SELECT a.id, a.game_id, a.game_name, a.publish_time, b.show_icon FROM game a " +
		"LEFT JOIN game_all b ON a.game_id=b.game_id WHERE a.game_id>0 ORDER BY a.create_time DESC  LIMIT 5").Values(&maps)

	return
}

func AddTypeInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	gameIds := make([]int, len(*ss))
	for i, s := range *ss {
		gameIds[i] = s.GameType
	}
	var games []Types
	_, err := o.QueryTable(new(Types)).Filter("Id__in", gameIds).All(&games, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]Types{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.GameType]; ok {
			(*ss)[i].Leixing = &g
		}
	}
	return
}

func AddSourceInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	gameIds := make([]int, len(*ss))
	for i, s := range *ss {
		gameIds[i] = s.Source
	}
	var games []Types
	_, err := o.QueryTable(new(Types)).Filter("Id__in", gameIds).All(&games, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]Types{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.Source]; ok {
			(*ss)[i].GameSource = &g
		}
	}
	return
}

//游戏评测添加评测信息
func AddResultCepingInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.Result
	}
	var games []Types
	_, err := o.QueryTable(new(Types)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]Types{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.Result]; ok {
			(*ss)[i].Ceping = &g
		}
	}
	return
}

//添加测试包信息
//func AddPackageInfo(ss *[]Game){
//	if ss == nil || len(*ss) == 0 {
//		return
//	}
//	o := orm.NewOrm()
//	gameIds := make([]int, len(*ss))
//	for i, s := range *ss {
//		gameIds[i] = s.Id
//	}
//	games := []GameTest{}
//	_, err := o.QueryTable(new(GameTest)).Filter("out_id__in", gameIds).All(&games, "Id", "OutId", "Package")
//	if err != nil {
//		return
//	}
//	gameMap := map[int]GameTest{}
//	for _, g := range games {
//		gameMap[g.OutId] = g
//	}
//	for i, s := range *ss {
//		if g, ok := gameMap[s.Id]; ok {
//			(*ss)[i].Package = g.Package
//		}
//	}
//	return
//}

//添加测试数据截图信息
func AddPictureInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	gameIds := make([]int, len(*ss))
	for i, s := range *ss {
		gameIds[i] = s.Id
	}
	var games []GameTest
	_, err := o.QueryTable(new(GameTest)).Filter("out_id__in", gameIds).All(&games, "Id", "OutId", "Picture")
	if err != nil {
		return
	}
	gameMap := map[int]GameTest{}
	for _, g := range games {
		gameMap[g.OutId] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.Id]; ok {
			(*ss)[i].Picture = g.Picture
		}
	}
	return
}

////添加高级账号信息
//func AddNumberInfo(ss *[]Game){
//	if ss == nil || len(*ss) == 0 {
//		return
//	}
//	o := orm.NewOrm()
//	for i, s := range *ss {
//		rrs := []GameTestNumber{}
//		_, err := o.QueryTable(new(GameTestNumber)).
//			Filter("out_id__in", s.Id).
//			All(&rrs)
//		if err != nil {
//			return
//		}
//		(*ss)[i].Numbers = &rrs
//	}
//	return
//}

func AddHezuoInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.Cooperation
	}
	var games []Types
	_, err := o.QueryTable(new(Types)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]Types{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.Cooperation]; ok {
			(*ss)[i].Hezuo = &g
		}
	}
	return
}

func AddYanfaInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.Development
	}
	var games []CompanyType
	_, err := o.QueryTable(new(CompanyType)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]CompanyType{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.Development]; ok {
			(*ss)[i].Yanfa = &g
		}
	}
	return
}
func AddFaxingInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.Issue
	}
	var games []CompanyType
	_, err := o.QueryTable(new(CompanyType)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]CompanyType{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.Issue]; ok {
			(*ss)[i].Faxing = &g
		}
	}
	return
}
func AddGonghuiInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.SociatyPolicy
	}
	var games []Types
	_, err := o.QueryTable(new(Types)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]Types{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.SociatyPolicy]; ok {
			(*ss)[i].Gonghui = &g
		}
	}
	return
}

func AddUserInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.AccessPerson
	}
	var games []User
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name", "NickName")
	if err != nil {
		return
	}
	gameMap := map[int]User{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.AccessPerson]; ok {
			(*ss)[i].User = &g
		}
	}
	return
}

// 附加游戏接入状态信息
// 游戏如果未接入，则游戏状态为"---"
// 游戏如果已接入，该游戏所有渠道都已下架，则游戏对应下架，如果至少有一个渠道上架中，则该游戏仍然上架中
func AddAccessStateInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	for i, s := range *ss {
		//gameId>0则该游戏已接入
		if s.GameId != 0 {
			state := AccessState{}
			o.Raw("SELECT COUNT(id) count FROM channel_access WHERE game_id = ?", s.GameId).QueryRow(&state)
			// 该游戏已接入但还没有接入渠道,此时该游戏的状态仍然为上架中,时间为游戏接入的时间
			if state.Count == 0 {
				o.Raw("SELECT IFNULL(update_jrtime,0) time FROM game WHERE game_id = ?", s.GameId).QueryRow(&state)
				state.State = 1
				(*ss)[i].Access = &state
			} else {
				o.Raw("SELECT COUNT(access_state) count, IFNULL(MAX(access_update_time),0) time FROM channel_access WHERE game_id = ? AND access_state=1", s.GameId).QueryRow(&state)
				if state.Count >= 1 {
					state.State = 1
					//(*ss)[i].Access.State = 1
				} else { //游戏已下架，查询下架时间
					o.Raw("SELECT IFNULL(MAX(access_update_time),0) time FROM channel_access WHERE game_id = ? AND access_state=2", s.GameId).QueryRow(&state)
					state.State = 2
					//(*ss)[i].Access.State = 2
				}
				(*ss)[i].Access = &state
			}
		}
	}
	return
}

//获取更新人的信息
func AddUpdateInfo(ss *[]Game, flag string) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		if flag == "pc" { //评测页面 该更新人
			linkIds[i] = s.UpdateEvalUserID
		} else if flag == "jr" { //游戏接入
			linkIds[i] = s.UpdateJrUserID
		} else { //游戏提测页面
			linkIds[i] = s.UpdateRefUserID
		}
	}
	var users []User
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&users, "Id", "NickName")
	if err != nil {
		return
	}
	userMap := map[int]User{}
	for _, u := range users {
		userMap[u.Id] = u
	}
	for i, s := range *ss {
		if flag == "pc" { //游戏测评页面
			if u, ok := userMap[s.UpdateEvalUserID]; ok {
				(*ss)[i].UpdateUser = &u
			}

		} else if flag == "jr" { //游戏接入页面
			if u, ok := userMap[s.UpdateJrUserID]; ok {
				(*ss)[i].UpdateUser = &u
			}

		} else { //提测页面
			if u, ok := userMap[s.UpdateRefUserID]; ok {
				(*ss)[i].UpdateUser = &u
			}
		}

	}
	return
}

func GameLadder2Json(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	for i := range *ss {

		//(*ss)[i].Ladder =
		var f interface{}

		json.Unmarshal([]byte((*ss)[i].Ladders), &f)

		(*ss)[i].Ladder_front = f

	}
	return
}

func AddResultCepingPeopleInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	for i, s := range *ss {
		var ids []int
		json.Unmarshal([]byte(s.ResultPerson), &ids)
		var rrs []User
		for _, id := range ids {
			r := User{}
			_, err := o.QueryTable(new(User)).
				Filter("Id__in", id).
				All(&r, "Id", "Name", "NickName")
			if err != nil {
				return
			}
			rrs = append(rrs, r)
		}
		(*ss)[i].Cepingpeoples = &rrs
	}
	return
}

//func GetSeachResultPeoplInfo(ss *[]Game,p []string)  {
//	if ss == nil || len(*ss) == 0 || len(p)==0{
//		return
//	}
//	for i,s := range *ss {
//		s.ResultPerson
//	}
//}

func AddSubmitPeopleInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.SubmitPerson
	}
	var games []User
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name", "NickName")
	if err != nil {
		return
	}
	gameMap := map[int]User{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.SubmitPerson]; ok {
			(*ss)[i].User = &g
		}
	}
	return
}

func AddAdviseInfo(ss *[]Game) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	for i, s := range *ss {
		var advises []string
		json.Unmarshal([]byte(s.Advise), &advises)
		var rrs []Advise
		for _, advise := range advises {
			r := User{}
			ad := strings.Split(advise, ":")
			_, err := o.QueryTable(new(User)).
				Filter("Id__in", ad[0]).
				All(&r, "Id", "Name", "NickName")
			if err != nil {
				return
			}
			adv := Advise{}

			adv.Id = r.Id
			adv.Name = r.Name
			adv.NickName = r.Nickname
			if ad[1] == "1" {
				adv.Adv = "接入"
			} else {
				adv.Adv = "不接入"
			}
			rrs = append(rrs, adv)
		}
		(*ss)[i].Advises = &rrs
	}
	return
}

// 获取权限范围内游戏名称列表
func GetGameNameList(where map[string][]interface{}) (list []interface{}, err error) {
	o := orm.NewOrm()
	game := new(Game)
	qs := o.QueryTable(game)

	for k, v := range where {
		qs = qs.Filter(k, v...)
	}
	var l []orm.Params
	qs.Values(&l, "GameId", "GameName")

	for _, v := range l {
		list = append(list, v)
	}

	return list, err
}

//根据游戏id获取游戏名
func GetGameNameById(gameId int) (string) {
	o := orm.NewOrm()
	game := Game{GameId: gameId}
	o.Read(&game, "GameId")
	return game.GameName
}

// 游戏提测修改
func ReferenceUpdateGameById(m *Game, where map[string][]interface{}) (err error) {
	o := orm.NewOrm()
	v := Game{}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}
	//m.UpdateTime = time.Now().Unix()
	fields := utils.GetNotEmptyFields(m, "GameName", "PublishTime", "GameType", "Issue", "Remarks",
		"Ip", "UpdateRefUserID", "UpdateRefTime", "Star", "Budget", "Number", "Package", "Picture", "Source")

	if _, err = o.Update(m, fields...); err != nil {
		return
	}

	fields = utils.RemoveFields(fields, "UpdateRefUserID", "UpdateRefTime")
	if err = CompareAndAddOperateLog(&v, m, m.UpdateRefUserID, bean.OPP_GAME_REFERENCE, v.Id, bean.OPA_UPDATE, fields...); err != nil {
		return
	}

	return
}

// 游戏评测修改
func GameEvaluationUpdateGameById(m *Game, where map[string][]interface{}) (err error) {
	o := orm.NewOrm()
	v := Game{}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}
	fields := utils.GetNotEmptyFields(m, "Package", "Result", "ResultPerson", "ResultTime", "UpdateEvalUserID", "UpdateEvalTime",
		"ResultReportWord", "ResultReportExcel", "Advise")

	if _, err = o.Update(m, fields...); err != nil {
		return
	}

	fields = utils.RemoveFields(fields, "UpdateEvalUserID", "UpdateEvalTime")
	if err = CompareAndAddOperateLog(&v, m, m.UpdateEvalUserID, bean.OPP_GAME_RESULT, v.Id, bean.OPA_UPDATE, fields...); err != nil {
		return
	}

	return
}

// 游戏接入修改
func GameAccessUpdateGameById(m *Game, where map[string][]interface{}) (err error) {
	o := orm.NewOrm()
	v := Game{}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}
	//m.UpdateTime = time.Now().Unix()
	fields := utils.GetNotEmptyFields(m, "GameId", "ImportTime", "PublishTime", "GameType", "Cooperation", "Development",
		"Issue", "QuantityPolicy", "SociatyPolicy", "AccessPerson", "Remarks", "UpdateTime",
		"Ladders", "UpdateJrUserID", "UpdateJrTime", "BodyMy")

	if _, err = o.Update(m, fields...); err != nil {
		return
	}

	fields = utils.RemoveFields(fields, "UpdateJrUserID", "UpdateJrTime")
	if err = CompareAndAddOperateLog(&v, m, m.UpdateJrUserID, bean.OPP_GAME_ACCESS, v.Id, bean.OPA_UPDATE, fields...); err != nil {
		return
	}

	//若游戏首发时间有修改，同步到渠道接入表中
	if v.PublishTime != m.PublishTime {
		o.Raw("UPDATE channel_access SET publish_time = ? WHERE game_id = ?", m.PublishTime, v.GameId).Exec()
	}

	return
}

// DeleteGame deletes Game by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGame(id int) (err error) {
	o := orm.NewOrm()
	v := Game{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Game{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetGameNameByGameId(id int) (name string, err error) {
	table := "gameId2Name"
	name, err = utils.Redis.HMGETOne(table, strconv.Itoa(id))
	if err != nil {
		return
	}

	if name != "" {
		return
	}

	var games []Game
	_, err = orm.NewOrm().QueryTable("game").All(&games, "game_id", "game_name")
	if err != nil {
		return
	}

	data := make(map[string]interface{}, len(games))
	for _, v := range games {
		data[strconv.Itoa(v.GameId)] = v.GameName
	}

	name, _ = util.Interface2String(data[strconv.Itoa(id)], false)
	if name == "" {
		err = errors.New("404")
		return
	}

	err = utils.Redis.HMSETALL(table, data, 2*60)
	if err != nil {
		return
	}

	return
}

func GetGameByGamename(gamename string) (game Game) {
	game = Game{GameName: gamename}
	o := orm.NewOrm()
	o.Read(&game, "GameName")
	return game
}

//通过 game_id 获得游戏 返回的是一个游戏数组 game_ids
func GetGamesByGameIds(links []int) (games []Game, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Game))
	_, err = qs.Filter("game_id__in", links).All(&games, "GameId", "GameName")
	if err != nil {
		return
	}
	return games, nil

}

//检验该游戏是否存在
func CheckGameIsExistByGameName(name string, source int) (check bool) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Game))
	//var game Game
	check = qs.Filter("game_name__exact", name).Filter("source__exact", source).Exist()
	return check
}
