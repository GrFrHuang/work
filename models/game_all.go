package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type GameAll struct {
	Id       int    `json:"id,omitempty" orm:"column(id);auto"`
	GameId   int    `json:"game_id,omitempty" orm:"column(game_id)"`
	GameName string `json:"game_name,omitempty" orm:"column(game_name);size(255);null"`
	Icon     string `json:"icon,omitempty" orm:"column(icon);size(255);null"`
	ShowIcon string `json:"show_icon,omitempty" orm:"column(show_icon);size(255);null"`
}

func (t *GameAll) TableName() string {
	return "game_all"
}

func init() {
	orm.RegisterModel(new(GameAll))
}

// AddGameAll insert a new GameAll into database and returns
// last inserted Id on success.
func AddGameAll(m *GameAll) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetGameAllById retrieves GameAll by Id. Returns error if
// Id doesn't exist
func GetGameAllById(id int) (v *GameAll, err error) {
	o := orm.NewOrm()
	v = &GameAll{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetGameAllById retrieves GameAll by Id. Returns error if
// Id doesn't exist
func GetGameAllNameById(id int) (name string) {
	o := orm.NewOrm()
	v := &GameAll{Id: id}
	o.Read(v)

	return v.GameName
}

// GetAllGameAll retrieves all GameAll matches certain condition. Returns empty list if
// no records exist
func GetAllGameAll(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64, where map[string][]interface{}) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(GameAll))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
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

	var l []GameAll
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

// UpdateGameAll updates GameAll by Id and returns error if
// the record to be updated doesn't exist
func UpdateGameAllById(m *GameAll) (err error) {
	o := orm.NewOrm()
	v := GameAll{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteGameAll deletes GameAll by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGameAll(id int) (err error) {
	o := orm.NewOrm()
	v := GameAll{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&GameAll{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetGameIdByGameName(gameName string) (id int) {
	o := orm.NewOrm()
	v := &GameAll{GameName: gameName}
	o.Read(v, "GameName")
	return v.GameId
}

func GetGameAll(gameName string) (gameAll []GameAll) {
	o := orm.NewOrm()
	o.QueryTable(new(GameAll)).Filter("game_name__istartswith", gameName).All(&gameAll)
	return gameAll
}
