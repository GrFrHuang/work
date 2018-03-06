package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type GameOutage struct {
	Id           int    `json:"id,omiempty" orm:"column(id);auto"`
	GameId       int    `json:"game_id,omiempty" orm:"column(game_id);null"`
	IncrTime     int64  `json:"incr_time,omiempty" orm:"column(incr_time);null"`         //关闭新增时间
	RechargeTime int64  `json:"recharge_time,omiempty" orm:"column(recharge_time);null"` //关闭充值时间
	ServerTime   int64  `json:"server_time,omiempty" orm:"column(server_time);null"`     //关闭服务器时间
	CreateTime   int64  `json:"create_time,omiempty" orm:"column(create_time);null"`     //创建时间
	CreatePerson int    `json:"create_person,omiempty" orm:"column(create_person);null"` //创建人
	Desc         string `json:"desc,omiempty" orm:"column(desc);null"`

	Game       *Game `orm:"-" json:"game,omitempty"`
	CreateUser *User `orm:"-" json:"create_user,omitempty"`
}

func (t *GameOutage) TableName() string {
	return "game_outage"
}

func init() {
	orm.RegisterModel(new(GameOutage))
}

// AddGameOutage insert a new GameOutage into database and returns
// last inserted Id on success.
func AddGameOutage(m *GameOutage) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GameOutAgeAddGameInfo(ss *[]GameOutage) {
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

func GetOutAgeGameName(typ int) (games []Game, err error) {
	o := orm.NewOrm()

	var sql string
	if typ == 1 {
		// 停运页所有停运的游戏
		sql = "SELECT a.game_id,b.game_name FROM game_outage a LEFT JOIN game b ON a.game_id=b.game_id "
	} else {
		// 添加游戏停运页，所有能够添加游戏停运的游戏并且还没添加停运的游戏
		sql = "SELECT game_id,game_name FROM game WHERE game_id>0 AND game_id NOT IN(SELECT game_id FROM game_outage) "
	}

	_, err = o.Raw(sql).QueryRows(&games)

	return
}

func GameOutAgeAddUpdateUserInfo(ss *[]GameOutage) {
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

// GetGameOutageById retrieves GameOutage by Id. Returns error if
// Id doesn't exist
func GetGameOutageById(id int) (v *GameOutage, err error) {
	o := orm.NewOrm()
	v = &GameOutage{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllGameOutage retrieves all GameOutage matches certain condition. Returns empty list if
// no records exist
func GetAllGameOutage(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(GameOutage))
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

	var l []GameOutage
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

// UpdateGameOutage updates GameOutage by Id and returns error if
// the record to be updated doesn't exist
func UpdateGameOutageById(m *GameOutage) (err error) {
	o := orm.NewOrm()
	v := GameOutage{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteGameOutage deletes GameOutage by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGameOutage(id int) (err error) {
	o := orm.NewOrm()
	v := GameOutage{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&GameOutage{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
