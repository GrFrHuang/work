package models

import (
	"errors"
	"fmt"
	"reflect"
	"github.com/astaxie/beego/orm"
	"strconv"
	"kuaifa.com/kuaifa/work-together/utils"
	"time"
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"strings"
)

type GamePlan struct {
	Id             int    `json:"id,omitempty" orm:"column(id);auto"`
	GameId         int    `json:"game_id,omitempty" orm:"column(game_id)"`
	Group          string `json:"group,omitempty" orm:"column(group);size(255);null"`
	SDKStatus      int    `json:"SDK_status,omitempty" orm:"column(SDK_status);null"`
	Operator       int    `json:"operator,omitempty" orm:"column(operator);null"`
	OperatorPerson string `json:"operator_person,omitempty" orm:"column(operator_person);null"`
	OperatorUpdate int    `json:"operator_update,omitempty" orm:"column(operator_update);null"`
	OperatorTime   int64  `json:"operator_time,omitempty" orm:"column(operator_time);null"`
	CustomerPerson int    `json:"customer_person,omitempty" orm:"column(customer_person);null"`
	CustomerTime   int64  `json:"customer_time,omitempty" orm:"column(customer_time);null"`
	CreateTime     int64  `json:"create_time,omitempty" orm:"column(create_time);null"`

	Channels       string `orm:"-" json:"channels"`                   //用于接收前端运营准备传回的所有渠道列表，用","分割, e.g. 732,945
	Gifts          string `orm:"-" json:"gifts"`                      //用于接收前端运营准备传回的所有礼包列表，用","分割, e.g. 732,945
	Materials      string `orm:"-" json:"materials"`                  //用于接收前端客服准备传回的所有素材列表，用","分割, e.g. 732,945
	Packages       string `orm:"-" json:"packages"`                   //用于接收前端客服准备传回的所有发包列表，用","分割, e.g. 732,945
	Tests          string `orm:"-" json:"tests"`                      //用于接收前端客服准备传回的所有已测包列表，用","分割, e.g. 732,945
	Game *Game   `orm:"-" json:"game,omitempty"`
	SDK  *Types  `orm:"-" json:"sdk,omitempty"`                       //SDK接入状态
	Yunying  *CompanyType  `orm:"-" json:"yunying,omitempty"`               //运营方
	User *User  `orm:"-" json:"user,omitempty"`
	Yunyingpeoples *[]User  `orm:"-" json:"yunyingpeoples,omitempty"` //运营负责人
	User2 *User  `orm:"-" json:"user2,omitempty"`                     //运营更新人
	Ceping *Ceping  `orm:"-" json:"ceping,omitempty"`                 //测评结果
	Details *Details `orm:"-" json:"details,omitempty"`               //礼包情况，素材情况，总渠道数，发包情况
}

type Details struct {
	Channels int	`json:"channels,omitempty" orm:"column(channels);null"`
	Gifts 	int	`json:"gifs,omitempty" orm:"column(gifts);null"`
	Materials	int 	`json:"materials,omitempty" orm:"column(materials);null"`
	Packages	int	`json:"packages,omitempty" orm:"column(packages);null"`
	Test	int	`json:"test,omitempty" orm:"column(test);null"`
}

type Ceping struct {
	Result      string  `json:"result,omitempty" orm:"column(result);null"`
	ReportWord  int     `json:"report_word,omitempty" orm:"column(word);null"`
	ReportExcel int     `json:"report_excel,omitempty" orm:"column(excel);null"`
}

func (t *GamePlan) TableName() string {
	return "game_plan"
}

func init() {
	orm.RegisterModel(new(GamePlan))
}

// AddGamePlan insert a new GamePlan into database and returns
// last inserted Id on success.
func AddGamePlan(m *GamePlan) (id int64, err error) {
	o := orm.NewOrm()
	m.CreateTime = time.Now().Unix()
	m.Group = "否"
	id, err = o.Insert(m)
	return
}

// GetGamePlanById retrieves GamePlan by Id. Returns error if
// Id doesn't exist
func GetGamePlanById(id int) (v *GamePlan, err error) {
	o := orm.NewOrm()
	v = &GamePlan{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllGamePlan retrieves all GamePlan matches certain condition. Returns empty list if
// no records exist
func GetAllGamePlan(query []string, fields []string, sortby []string, order []string,
	offset int64, limit int64, where map[string][]interface{}) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(GamePlan))
	// query k=v
	var tmp_games []interface{}
	if query != nil {
		for _, v := range query {
			// rewrite dot-notation to Object__Attribute
			game_id, _ := strconv.Atoi(v)
			tmp_games = append(tmp_games, game_id)
			//_ = strings.Replace(_, ".", "__", -1)
			//if strings.Contains(_, "isnull") {
			//	qs = qs.Filter(_, (v == "true" || v == "1"))
			//} else {
			//	qs = qs.Filter(_, v)
			//}
		}
		fmt.Printf("tmp_games:%v\n", tmp_games)
		qs = qs.Filter("game_id__in", tmp_games...)
	}

	for k, v := range where {
		qs = qs.Filter(k, v...)
	}

	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []GamePlan
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateGamePlan updates GamePlan by Id and returns error if
// the record to be updated doesn't exist
func UpdateGamePlanById(m *GamePlan, where map[string][]interface{}) (err error) {
	o := orm.NewOrm()

	v := GamePlan{}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}

	fields := utils.GetNotEmptyFields(m, "Group", "SDKStatus", "Operator", "OperatorPerson", "OperatorUpdate",
		"OperatorTime", "Result", "ResultPerson", "ResultTime", "ResultReport", "ResultReportExcel",
		"ResultPackage", "CustomerPerson", "CustomerTime")
	if _, err = o.Update(m, fields...); err != nil{
		return
	}

	return
}

// 客服修改
func CustumerUpdateGamePlanById(m *GamePlan, where map[string][]interface{}) (err error) {
	o := orm.NewOrm()

	v := GamePlan{}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}

	fields := utils.GetNotEmptyFields(m, "Group", "SDKStatus", "CustomerPerson", "CustomerTime")
	if _, err = o.Update(m, fields...); err != nil{
		return
	}

	materials, err := GetMaterialChannelsByGameId(v.GameId)
	if err != nil{
		return
	}
	v.Materials = MakeChannels(materials)

	packages, err := GetPackageChannelsByGameId(v.GameId)
	if err != nil{
		return
	}
	v.Packages = MakeChannels(packages)

	tests, err := GetTestChannelsByGameId(v.GameId)
	if err != nil{
		return
	}
	v.Tests = MakeChannels(tests)

	fields = append(fields, "Materials", "Packages", "Tests")
	fields = utils.RemoveFields(fields, "CustomerPerson", "CustomerTime")//去除CustomerPerson和CustomerTime字段,不计入操作日志

	if err = CompareAndAddOperateLog(&v, m, m.CustomerPerson, bean.OPP_GAME_PLAN_CUSTOMER, v.Id, bean.OPA_UPDATE, fields...); err != nil{
		return
	}

	return
}

//将渠道详情改为只有渠道code组成的字符串
func MakeChannels(channels []orm.Params) (channelCodes string){
	var channel_a []string
	for _, channel := range channels{
		channel_a = append(channel_a, channel["cp"].(string))
	}
	channelCodes = strings.Join(channel_a, ",")
	return
}

// 运营修改
func OperateUpdateGamePlanById(m *GamePlan, where map[string][]interface{}) (err error) {
	o := orm.NewOrm()

	v := GamePlan{}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}

	fields := utils.GetNotEmptyFields(m, "Operator", "OperatorUpdate", "OperatorPerson",  "OperatorTime")
	if _, err = o.Update(m, fields...); err != nil{
		return
	}

	gifts, err := GetGiftChannelsByGameId(v.GameId)
	if err != nil{
		return
	}
	v.Gifts = MakeChannels(gifts)

	fields = append(fields, "Gifts")
	fields = utils.RemoveFields(fields, "OperatorPerson", "OperatorTime")//去除OperatorTime和OperatorPerson字段,不计入操作日志

	if err = CompareAndAddOperateLog(&v, m, m.OperatorUpdate, bean.OPP_GAME_PLAN_OPERATE, v.Id, bean.OPA_UPDATE, fields...); err != nil{
		return
	}

	return
}

// DeleteGamePlan deletes GamePlan by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGamePlan(id int) (err error) {
	o := orm.NewOrm()
	v := GamePlan{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&GamePlan{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//添加礼包情况，素材情况，总渠道数，发包情况，测包情况
func AddExactGamePlanChannel(ss *[]GamePlan){

	if ss == nil || len(*ss) == 0 {
		return
	}
	//o := orm.NewOrm()
	//linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		maps := GetExactGamePlanChannelById(s.GameId)
		(*ss)[i].Details = &maps
	}
	return
}

func AddGameInfo(ss *[]GamePlan){
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.GameId
	}
	games := []Game{}
	_, err := o.QueryTable(new(Game)).
		Filter("game_id__in", linkIds).
		All(&games, "Id", "GameId", "GameName")
	if err != nil {
		return
	}
	gameMap := map[int]Game{}
	for _, g := range games {
		gameMap[g.GameId] = g
	}
	//fmt.Printf("games:%v\n",games)
	//fmt.Printf("gameMap:%v\n",gameMap)
	for i, s := range *ss {
		if g, ok := gameMap[s.GameId]; ok {
			(*ss)[i].Game = &g
		}
	}
	return
}

func AddSDKInfo(ss *[]GamePlan){
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.SDKStatus
	}
	games := []Types{}
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
		if g, ok := gameMap[s.SDKStatus]; ok {
			(*ss)[i].SDK = &g
		}
	}
	return
}

func AddYunyingInfo(ss *[]GamePlan){
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.Operator
	}
	games := []CompanyType{}
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
		if g, ok := gameMap[s.Operator]; ok {
			(*ss)[i].Yunying = &g
		}
	}
	return
}

func AddYunyingPeopleInfo(ss *[]GamePlan) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	for i, s := range *ss {
		ids := []int{}
		json.Unmarshal([]byte(s.OperatorPerson), &ids)
		rrs := []User{}
		for _, id := range ids{
			r := User{}
			_, err := o.QueryTable(new(User)).
				Filter("Id__in", id).
				All(&r, "Id", "Name", "NickName")
			if err != nil {
				return
			}
			rrs = append(rrs, r)
		}
		(*ss)[i].Yunyingpeoples = &rrs
	}
	return
}

func AddKefuPeopleInfo(ss *[]GamePlan) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.CustomerPerson
	}
	games := []User{}
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
		if g, ok := gameMap[s.CustomerPerson]; ok {
			(*ss)[i].User = &g
		}
	}
	return
}

func AddYunyingUpdateInfo(ss *[]GamePlan) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.OperatorUpdate
	}
	games := []User{}
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
		if g, ok := gameMap[s.OperatorUpdate]; ok {
			(*ss)[i].User2 = &g
		}
	}
	return
}

func AddCepingInfo(ss *[]GamePlan){

	if ss == nil || len(*ss) == 0 {
		return
	}
	//o := orm.NewOrm()
	//linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		maps := GetResultByGameId(s.GameId);
		(*ss)[i].Ceping = &maps
		fmt.Printf("maps-----------%v\n", maps)
	}
	return




	//if ss == nil || len(*ss) == 0 {
	//	return
	//}
	//o := orm.NewOrm()
	//linkIds := make([]int, len(*ss))
	//for i, s := range *ss {
	//	linkIds[i] = s.Result
	//}
	//games := []Types{}
	//_, err := o.QueryTable(new(Types)).
	//	Filter("Id__in", linkIds).
	//	All(&games, "Id", "Name")
	//if err != nil {
	//	return
	//}
	//gameMap := map[int]Types{}
	//for _, g := range games {
	//	gameMap[g.Id] = g
	//}
	//for i, s := range *ss {
	//	if g, ok := gameMap[s.Result]; ok {
	//		(*ss)[i].Ceping = &g
	//	}
	//}
	//return


}

func AddCepingPeopleInfo(ss *[]GamePlan) {
	//if ss == nil || len(*ss) == 0 {
	//	return
	//}
	//o := orm.NewOrm()
	//linkIds := make([]int, len(*ss))
	//for i, s := range *ss {
	//	linkIds[i] = s.ResultPerson
	//}
	//games := []User{}
	//_, err := o.QueryTable(new(User)).
	//	Filter("Id__in", linkIds).
	//	All(&games, "Id", "Name", "NickName")
	//if err != nil {
	//	return
	//}
	//gameMap := map[int]User{}
	//for _, g := range games {
	//	gameMap[g.Id] = g
	//}
	//for i, s := range *ss {
	//	if g, ok := gameMap[s.ResultPerson]; ok {
	//		(*ss)[i].User = &g
	//	}
	//}
	return
}

