package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"encoding/json"
	"time"
	"kuaifa.com/kuaifa/work-together/utils"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

type GameUpdate struct {
	Id               int    `json:"id,omitempty" orm:"column(id);auto"`
	GameId           int    `json:"game_id,omitempty" orm:"column(game_id)"`
	GameUpdateTime   int64  `json:"game_update_time,omitempty" orm:"column(game_update_time);null"`
	UpdateType       int    `json:"update_type,omitempty" orm:"column(update_type);null"`
	Version          string `json:"version,omitempty" orm:"column(version);size(255);null"`
	VersionCode      string `json:"version_code,omitempty" orm:"column(version_code);size(255);null"`
	PackageName      string `json:"package_name,omitempty" orm:"column(package_name);size(255);null"`
	Material         string `json:"material,omitempty" orm:"column(material);null"`
	Icon             string `json:"icon,omitempty" orm:"column(icon);null"`
	UpdateChannel    string `json:"update_channel,omitempty" orm:"column(update_channel);size(255);null"`
	NotUpdateChannel string `json:"not_update_channel,omitempty" orm:"column(not_update_channel);size(255);null"`
	StopChannel      string `json:"stop_channel,omitempty" orm:"column(stop_channel);size(255);null"`
	Remark           string `json:"remark,omitempty" orm:"column(remark);size(255);null"`
	CreatePerson     int `json:"create_person,omitempty" orm:"column(create_person);size(255);null"`
	CreateTime       int64  `json:"create_time,omitempty" orm:"column(create_time);null"`
	UpdatePerson     int `json:"update_person,omitempty" orm:"column(update_person);size(255);null"`
	UpdateTime       int64  `json:"update_time,omitempty" orm:"column(update_time);null"`

	Game             *Game   `orm:"-" json:"game,omitempty"`
	CreateUser       *User   `orm:"-" json:"create_user,omitempty"`
	UpdateUser       *User   `orm:"-" json:"update_user,omitempty"`
	UpdateChannels   *[]Channel `orm:"-" json:"update_channels,omitempty"`
	NotUpdateChannels   *[]Channel `orm:"-" json:"not_update_channels,omitempty"`
	StopChannels     *[]Channel `orm:"-" json:"stop_channels,omitempty"`
}

func (t *GameUpdate) TableName() string {
	return "game_update"
}

func init() {
	orm.RegisterModel(new(GameUpdate))
}

// AddGameUpdate insert a new GameUpdate into database and returns
// last inserted Id on success.
func AddGameUpdate(m *GameUpdate) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)

	err = CompareAndAddOperateLog(nil, m, m.CreatePerson, bean.OPP_GAME_UPDATE, int(id), bean.OPA_INSERT)
	return
}

// GetGameUpdateById retrieves GameUpdate by Id. Returns error if
// Id doesn't exist
func GetGameUpdateById(id int) (v *GameUpdate, err error) {
	o := orm.NewOrm()
	v = &GameUpdate{Id: id}
	if err = o.Read(v); err != nil {
		return
	}

	//附加游戏信息
	game := Game{}
	err = o.QueryTable(new(Game)).
		Filter("game_id__in", v.GameId).
		One(&game, "Id", "GameId", "GameName")
	if err != nil {
		return
	}else{
		v.Game = &game
	}

	return
}

func UpdateAddGameInfo(ss *[]GameUpdate){
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

func UpdateAddCreateUserInfo(ss *[]GameUpdate) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.CreatePerson
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
		if g, ok := gameMap[s.CreatePerson]; ok {
			(*ss)[i].CreateUser = &g
		}
	}
	return
}

func UpdateAddUpdateUserInfo(ss *[]GameUpdate) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.UpdatePerson
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
		if g, ok := gameMap[s.UpdatePerson]; ok {
			(*ss)[i].UpdateUser = &g
		}
	}
	return
}

func UpdateAddUpdateChannelInfo(ss *[]GameUpdate) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	for i, s := range *ss {
		channels := []string{}
		json.Unmarshal([]byte(s.UpdateChannel), &channels)
		rrs := []Channel{}
		for _, channel := range channels {
			r := Channel{}
			_, err := o.QueryTable(new(Channel)).
				Filter("cp__in", channel).
				All(&r, "cp", "Name")
			if err != nil {
				return
			}
			rrs = append(rrs, r)
		}
		(*ss)[i].UpdateChannels = &rrs
	}
	return
}

func UpdateAddNotUpdateChannelInfo(ss *[]GameUpdate) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	for i, s := range *ss {
		channels := []string{}
		json.Unmarshal([]byte(s.NotUpdateChannel), &channels)
		rrs := []Channel{}
		for _, channel := range channels{
			r := Channel{}
			_, err := o.QueryTable(new(Channel)).
				Filter("cp__in", channel).
				All(&r, "cp", "Name")
			if err != nil {
				return
			}
			rrs = append(rrs, r)
		}
		(*ss)[i].NotUpdateChannels = &rrs
	}
	return
}

func UpdateAddStopChannelInfo(ss *[]GameUpdate) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	for i, s := range *ss {
		channels := []string{}
		json.Unmarshal([]byte(s.StopChannel), &channels)
		rrs := []Channel{}
		for _, channel := range channels{
			r := Channel{}
			_, err := o.QueryTable(new(Channel)).
				Filter("cp__in", channel).
				All(&r, "cp", "Name")
			if err != nil {
				return
			}
			rrs = append(rrs, r)
		}
		(*ss)[i].StopChannels = &rrs
	}
	return
}

// GetAllGameUpdate retrieves all GameUpdate matches certain condition. Returns empty list if
// no records exist
func GetAllGameUpdate(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(GameUpdate))
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

	var l []GameUpdate
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

// UpdateGameUpdate updates GameUpdate by Id and returns error if
// the record to be updated doesn't exist
func UpdateGameUpdateById(m *GameUpdate, where map[string][]interface{}) (err error) {

	o := orm.NewOrm()
	v := GameUpdate{}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}
	m.UpdateTime = time.Now().Unix()
	fields := utils.GetNotEmptyFields(m, "GameUpdateTime", "UpdateType", "Version", "VersionCode", "PackageName",
		"Material", "Icon", "UpdateChannel", "NotUpdateChannel", "StopChannel", "UpdatePerson", "UpdateTime")
	//Remark，但是前端传空想删除此字段的值时，上面的函数并不能完成修改，所以需要加下一行代码修改Remark
	fields = append(fields, "Remark")
	if _, err = o.Update(m, fields...); err != nil {
		return
	}

	fields = utils.RemoveFields(fields, "UpdatePerson", "UpdateTime")//去除UpdatePerson和UpdateTime字段,不计入操作日志

	if err = CompareAndAddOperateLog(&v, m, m.UpdatePerson, bean.OPP_GAME_UPDATE, v.Id, bean.OPA_UPDATE, fields...); err != nil{
		return
	}
	return
}

// DeleteGameUpdate deletes GameUpdate by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGameUpdate(id int) (err error) {
	o := orm.NewOrm()
	v := GameUpdate{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&GameUpdate{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
