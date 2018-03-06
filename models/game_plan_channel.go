package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/utils"
	"strconv"
)

type GamePlanChannel struct {
	Id            int `json:"id,omitempty" orm:"column(id);auto"`
	GameId        int `json:"game_id,omitempty" orm:"column(game_id);null"`
	//ChannelId     int `json:"channel_id,omitempty" orm:"column(channel_id);null"`
	ChannelCode   string  `json:"channel_code,omitempty" orm:"column(channel_code);null"`
	GiftStatus    int `json:"gift_status,omitempty" orm:"column(gift_status);null"`
	Material      int `json:"material,omitempty" orm:"column(material);null"`
	PackageStatus int `json:"package_status,omitempty" orm:"column(package_status);null"`
	Test          int `json:"test,omitempty" orm:"column(test);null"`
}

type GamePlanChannel_bak struct {
	GameId      int
	Channels    int
	Gifts       int
	Materials   int
	Packages    int
}

func (t *GamePlanChannel) TableName() string {
	return "game_plan_channel"
}

func init() {
	orm.RegisterModel(new(GamePlanChannel))
}

// AddGamePlanChannel insert a new GamePlanChannel into database and returns
// last inserted Id on success.
func AddGamePlanChannel(m *GamePlanChannel, gameId int, channelCode string) (id int64, err error) {
	o := orm.NewOrm()
	m.GameId = gameId
	m.ChannelCode = channelCode
	id, err = o.Insert(m)
	return
}

// 根据游戏id以及渠道id，获取该条数据
func GetGamePlanChannel(gameid int, channelCode string) (v GamePlanChannel, err error) {
	o := orm.NewOrm()
	//var c GamePlanChannel
	err = o.QueryTable(new(GamePlanChannel)).Filter("game_id__exact", gameid).Filter("channel_code__exact", channelCode).One(&v)
	//id = c.Id
	return
}

// GetGamePlanChannelById retrieves GamePlanChannel by Id. Returns error if
// Id doesn't exist
func GetExactGamePlanChannelById(id int) (ma Details) {
	o := orm.NewOrm()
	//var maps []orm.Params

	var maps Details

	_ = o.Raw("select count(channel_code) as channels, coalesce(sum(gift_status),0) as gifts, " +
		"coalesce(sum(material),0) as materials, coalesce(sum(package_status),0) as packages, " +
		"coalesce(sum(test),0) as test from game_plan_channel where game_id = ?", id).QueryRow(&maps)

	ma = maps
	return
}

// GetAllGamePlanChannel retrieves all GamePlanChannel matches certain condition. Returns empty list if
// no records exist
func GetAllGamePlanChannel(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(GamePlanChannel))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
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

	var l []GamePlanChannel
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

//根据id，获取具体数据
func GetOneById(id int)(v *GamePlanChannel, err error) {
	o := orm.NewOrm()
	v = &GamePlanChannel{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateGamePlanChannel updates GamePlanChannel by Id and returns error if
// the record to be updated doesn't exist
func UpdateGamePlanChannel(m *GamePlanChannel) (err error) {
	o := orm.NewOrm()
	//v := GamePlanChannel{GameId: m.GameId, ChannelId: m.ChannelId}

	if m.GiftStatus == 0{
		o.Raw("UPDATE game_plan_channel SET gift_status = 0 WHERE id = ?", m.Id).Exec()
	}

	if m.Material == 0{
		o.Raw("UPDATE game_plan_channel SET material = 0 WHERE id = ?", m.Id).Exec()
	}

	if m.PackageStatus == 0{
		o.Raw("UPDATE game_plan_channel SET package_status = 0 WHERE id = ?", m.Id).Exec()
	}

	if m.Test == 0{
		o.Raw("UPDATE game_plan_channel SET test = 0 WHERE id = ?", m.Id).Exec()
	}

	if m.GiftStatus == 0 && m.Material == 0 && m.PackageStatus == 0{
		return
	}

	fields := utils.GetNotEmptyFields(m, "GiftStatus", "Material", "PackageStatus", "Test")

	//fmt.Printf("fields:%v\n", fields)
	_, err = o.Update(m, fields...)
	if err != nil {
		return
	}

	return
}

// DeleteGamePlanChannel deletes GamePlanChannel by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGamePlanChannel(id int) (err error) {
	o := orm.NewOrm()
	v := GamePlanChannel{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&GamePlanChannel{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//根据游戏ID和渠道ID获取该游戏上线计划的ID
func GetGamePlanChannelId(gameId int, channelCode string)(id int){
	o := orm.NewOrm()
	var maps []orm.Params
	o.Raw("select id from game_plan_channel where game_id = ? and channel_Code = ?", gameId, channelCode).Values(&maps)
	id,_ = strconv.Atoi(maps[0]["id"].(string))
	return
}

//根据游戏ID获取该游戏所有下发渠道列表
func GetAllChannelsByGameId(gameId int)(maps []orm.Params, err error)  {
	o := orm.NewOrm()
	_,err =o.Raw("SELECT cp,name FROM channel WHERE cp IN (SELECT channel_code FROM game_plan_channel" +
		" WHERE game_id = ?)", gameId).Values(&maps)
	return
}

//根据游戏ID获取该游戏上线准备详情
func GetDetailsByGameId(gameId int)(maps []orm.Params, err error)  {
	o := orm.NewOrm()
	_,err =o.Raw("SELECT a.name,b.gift_status,b.material,b.package_status,b.test FROM channel a " +
		"LEFT JOIN game_plan_channel b ON a.cp = b.channel_code WHERE b.game_id = ? " +
		"GROUP BY NAME ORDER BY SUM(b.gift_status+b.material+b.package_status+b.test) asc", gameId).Values(&maps)
	return
}

//根据游戏ID获取该游戏所有下发渠道
func GetChannelsByGameId(gameId int)(maps []orm.Params, err error)  {
	o := orm.NewOrm()
	_,err =o.Raw("SELECT a.cp,a.name,c.cooperate_state FROM channel a LEFT JOIN game_plan_channel b ON a.cp = b.channel_code " +
		" LEFT JOIN channel_company c ON a.cp=c.channel_code WHERE b.game_id = ? ", gameId).Values(&maps)
	return
}

//根据游戏ID获取该游戏所有下发渠道列表中已发礼包渠道列表
func GetGiftChannelsByGameId(gameId int)(maps []orm.Params, err error)  {
	o := orm.NewOrm()
	_,err =o.Raw("SELECT cp,name FROM channel WHERE cp IN (SELECT channel_code FROM game_plan_channel" +
		" WHERE game_id = ? and gift_status = 1)", gameId).Values(&maps)
	return
}

//根据游戏ID获取该游戏所有下发渠道列表中已发包渠道列表
func GetPackageChannelsByGameId(gameId int)(maps []orm.Params, err error)  {
	o := orm.NewOrm()
	_,err =o.Raw("SELECT cp,name FROM channel WHERE cp IN (SELECT channel_code FROM game_plan_channel" +
		" WHERE game_id = ? and package_status = 1)", gameId).Values(&maps)
	return
}

//根据游戏ID获取该游戏所有下发渠道列表中已测包渠道列表
func GetTestChannelsByGameId(gameId int)(maps []orm.Params, err error)  {
	o := orm.NewOrm()
	_,err =o.Raw("SELECT cp,name FROM channel WHERE cp IN (SELECT channel_code FROM game_plan_channel" +
		" WHERE game_id = ? and test = 1)", gameId).Values(&maps)
	return
}

//根据游戏ID获取该游戏所有下发渠道列表中已发素材渠道列表
func GetMaterialChannelsByGameId(gameId int)(maps []orm.Params, err error)  {
	o := orm.NewOrm()
	_,err =o.Raw("SELECT cp,name FROM channel WHERE cp IN (SELECT channel_code FROM game_plan_channel" +
		" WHERE game_id = ? and material = 1)", gameId).Values(&maps)
	return
}