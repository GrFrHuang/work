package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"time"
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/utils"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"strconv"
)

type ChannelAccess struct {
	Id     int `json:"id,omitempty" orm:"column(id);auto"`
	GameId int `json:"game_id,omitempty" orm:"column(game_id)"`
	//ChannelId      int    `json:"channel_id,omitempty" orm:"column(channel_id)"`
	ChannelCode      string `json:"channel_code,omitempty" orm:"column(channel_code)"`
	PublishTime      int64  `json:"publish_time,omitempty" orm:"column(publish_time)"`
	BodyMy           int    `json:"body_my,omitempty" orm:"column(body_my)"`
	Cooperation      int    `json:"cooperation,omitempty" orm:"column(cooperation)"`
	Ladders          string `json:"ladders,omitempty" orm:"column(ladders)"`
	BusinessPerson   int    `json:"business_person,omitempty" orm:"column(business_person)"`
	AccessState      int    `json:"access_state,omitempty" orm:"column(access_state)"`
	AccessUpdateUser int    `json:"access_update_user,omitempty" orm:"column(access_update_user)"`
	AccessUpdateTime int64  `json:"access_update_time,omitempty" orm:"column(access_update_time)"`
	//State          int    `json:"state,omitempty" orm:"column(state)"`
	//Finance        int    `json:"finance,omitempty" orm:"column(finance);null"`
	//ContractState  int    `json:"contract_state,omitempty" orm:"column(contract_state);null"`
	//PayState       string `json:"pay_state,omitempty" orm:"column(pay_state);size(255);null"`
	CreateTime     int64  `json:"create_time,omitempty" orm:"column(create_time);null"`
	UpdateTime     int64  `json:"update_time,omitempty" orm:"column(update_time);null"`
	WorkflowStatus int    `json:"workflow_status" orm:"column(workflow_status);null"`
	Accessory      string `json:"accessory" orm:"column(accessory);null"`
	PactId         string    `json:"pact_id" orm:"column(pactId);null"`

	UpdateChannelTime   int64 `json:"update_channeltime,omitempty" orm:"column(update_channeltime);null"`
	UpdateChannelUserID int   `json:"update_channeluserid,omitempty" orm:"column(update_channeluserid)"`

	UpdateUser      *User    `orm:"-" json:"update_user,omitempty"`
	StateUpdateUser *User    `orm:"-" json:"state_update_user,omitempty"`
	Group           []Group  `json:"group,omitempty" orm:"-"`
	Game            *Game    `orm:"-" json:"game,omitempty"`
	Channel         *Channel `orm:"-" json:"channel,omitempty"`
	Hezuo           *Types   `orm:"-" json:"hezuo,omitempty"`
	Business        *User    `orm:"-" json:"business,omitempty"`
	Caiwu           *User    `orm:"-" json:"caiwu,omitempty"`

	Ladder_front interface{} `orm:"-" json:"ladder_front"`
}

type Group struct {
	Ladder  string   `json:"ladder,omitempty"`
	Channel []string `json:"channel,omitempty"`
}

func (t *ChannelAccess) TableName() string {
	return "channel_access"
}

func init() {
	orm.RegisterModel(new(ChannelAccess))
}

// AddChannelAccess insert a new ChannelAccess into database and returns
// last inserted Id on success.
func AddChannelAccess(m *ChannelAccess) (id int64, err error) {
	m.Id = 0
	//m.State = 1	//默认审核状态为未审核
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = time.Now().Unix()
	o := orm.NewOrm()
	id, err = o.Insert(m)

	err = CompareAndAddOperateLog(nil, m, m.UpdateChannelUserID, bean.OPP_CHANNEL_ACCESS, int(id), bean.OPA_INSERT)
	return
}

func ChangeWorkflowStatus(id, status int) (error) {
	o := orm.NewOrm()
	var m ChannelAccess
	m.Id = id
	m.WorkflowStatus = status
	m.UpdateTime = time.Now().Unix()
	i, err := o.Update(&m, "workflow_status", "update_time")
	fmt.Println(err)
	if err != nil || i <= 0 {
		return errors.New("解锁渠道信息失败")
	}
	return nil
}

func ChangeWorkflowStatusByid(gameId, status int, channelId string) (int, error) {
	o := orm.NewOrm()
	var m ChannelAccess
	m.GameId = gameId
	m.ChannelCode = channelId
	if err := o.Read(&m, "game_id", "channel_code"); err != nil {
		return 0, err
	}
	m.WorkflowStatus = status
	m.UpdateTime = time.Now().Unix()
	i, err := o.Update(&m, "workflow_status", "update_time")
	if err != nil || i <= 0 {
		return 0, errors.New("锁定渠道信息失败")
	}
	return m.Id, nil
}

func ChannelAccessAddGameInfo(ss *[]ChannelAccess) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.GameId
	}
	var games []Game
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

func ChannelAccessAddChannelInfo(ss *[]ChannelAccess) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]string, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.ChannelCode
	}
	var link []Channel
	_, err := o.QueryTable(new(Channel)).
		Filter("cp__in", linkIds).
		All(&link, "Channelid", "Name", "Cp")
	if err != nil {
		return
	}
	linkMap := map[string]Channel{}
	for _, g := range link {
		linkMap[g.Cp] = g
	}
	for i, s := range *ss {
		if g, ok := linkMap[s.ChannelCode]; ok {
			(*ss)[i].Channel = &g
		}
	}
	return
}

func ChannelAccessAddHezuoInfo(ss *[]ChannelAccess) {
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

// 渠道接入统计
func GetChannelAccessStatistics() (statistic []Statistics) {
	statistic = []Statistics{}
	var query []Statistics

	us, err := GetUsersByDevMent(237)
	if err != nil {
		return
	}

	userIds := make([]string, 0)
	statMap := map[int]int{}
	for _, u := range us {
		userIds = append(userIds, strconv.Itoa(u.Id))
	}

	now := time.Now().Format("2006-01-02")
	sql := "SELECT  access_update_user AS user_id,COUNT(FROM_UNIXTIME(create_time,'%Y-%m-%d')) AS `total` FROM channel_access " +
		"WHERE access_update_user IN(" + strings.Join(userIds, ",") + ") AND FROM_UNIXTIME(create_time,'%Y-%m-%d') " +
		"= '" + now + "' GROUP BY access_update_user"

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

func ChannelAccessAddBusinessInfo(ss *[]ChannelAccess) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.BusinessPerson
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
		if g, ok := gameMap[s.BusinessPerson]; ok {
			(*ss)[i].Business = &g
		}
	}
	return
}

//func ChannelAccessAddCaiwuInfo(ss *[]ChannelAccess) {
//	if ss == nil || len(*ss) == 0 {
//		return
//	}
//	o := orm.NewOrm()
//	linkIds := make([]int, len(*ss))
//	for i, s := range *ss {
//		linkIds[i] = s.Finance
//	}
//	games := []User{}
//	_, err := o.QueryTable(new(User)).
//		Filter("Id__in", linkIds).
//		All(&games, "Id", "Name", "NickName")
//	if err != nil {
//		return
//	}
//	gameMap := map[int]User{}
//	for _, g := range games {
//		gameMap[g.Id] = g
//	}
//	for i, s := range *ss {
//		if g, ok := gameMap[s.Finance]; ok {
//			(*ss)[i].Caiwu = &g
//		}
//	}
//	return
//}

func ChannelAccessAddUpdateUserInfo(ss *[]ChannelAccess) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.UpdateChannelUserID
	}
	var users []User
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&users, "Id", "NickName")
	if err != nil {
		return
	}
	userMap := map[int]User{}
	for _, g := range users {
		userMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := userMap[s.UpdateChannelUserID]; ok {
			(*ss)[i].UpdateUser = &g
		}
	}
	return
}

// 附加接入状态修改人信息
func ChannelAccessAddStateInfo(ss *[]ChannelAccess) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.AccessUpdateUser
	}
	var users []User
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&users, "Id", "NickName")
	if err != nil {
		return
	}
	userMap := map[int]User{}
	for _, g := range users {
		userMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := userMap[s.AccessUpdateUser]; ok {
			(*ss)[i].StateUpdateUser = &g
		}
	}
	return
}

func ChannelAccessParseLadder2Json(ss *[]ChannelAccess) {
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

// GetChannelAccessById retrieves ChannelAccess by Id. Returns error if
// Id doesn't exist
func GetChannelAccessById(id int) (v *ChannelAccess, err error) {
	o := orm.NewOrm()
	v = &ChannelAccess{Id: id}
	if err = o.Read(v); err != nil {
		return
	}

	var f interface{}
	json.Unmarshal([]byte(v.Ladders), &f)
	v.Ladder_front = f

	game := Game{}
	err = o.QueryTable(new(Game)).
		Filter("game_id__in", v.GameId).
		One(&game, "Id", "GameId", "GameName")
	if err != nil {
		return
	} else {
		v.Game = &game
	}

	channel := Channel{}
	err = o.QueryTable(new(Channel)).
		Filter("cp__in", v.ChannelCode).
		One(&channel, "Channelid", "Name", "Cp")
	if err != nil {
		return
	} else {
		v.Channel = &channel
	}

	return
}

// GetAllChannelAccess retrieves all ChannelAccess matches certain condition. Returns empty list if
// no records exist
func GetAllChannelAccess(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ChannelAccess))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
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

	var l []ChannelAccess
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

// UpdateChannelAccess updates ChannelAccess by Id and returns error if
// the record to be updated doesn't exist
func UpdateChannelAccessById(m *ChannelAccess, where map[string][]interface{}) (err error) {
	o := orm.NewOrm()
	v := ChannelAccess{}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}
	m.UpdateTime = time.Now().Unix()
	fields := utils.GetNotEmptyFields(m, "BodyMy", "Cooperation", "Ladders", "BusinessPerson", "State", "Finance",
		"ContractState", "PayState", "UpdateTime", "UpdateChannelUserID", "UpdateChannelTime", "AccessState",
		"AccessUpdateUser", "AccessUpdateTime","PactId")

	s := ChannelAccess{Id: m.Id}
	if err = orm.NewOrm().Read(&s); err != nil {
		return
	}
	_, err = o.Update(m, fields...)
	if err != nil {
		return
	}
	fields = utils.RemoveFields(fields, "UpdateChannelUserID", "UpdateTime", "UpdateChannelTime", "AccessUpdateUser", "AccessUpdateTime")
	if err = CompareAndAddOperateLog(&s, m, m.UpdateChannelUserID, bean.OPP_CHANNEL_ACCESS, s.Id, bean.OPA_UPDATE, fields...); err != nil {
		return
	}

	return
}

// DeleteChannelAccess deletes ChannelAccess by Id and returns error if
// the record to be deleted doesn't exist
func DeleteChannelAccess(id int) (err error) {
	o := orm.NewOrm()
	v := ChannelAccess{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ChannelAccess{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
