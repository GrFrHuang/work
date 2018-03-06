package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type GameTestNumber struct {
	Id       int    `json:"id,omitempty" orm:"column(id);auto"`
	OutId    int    `json:"out_id,omitempty" orm:"column(out_id);null"`
	Username string `json:"username,omitempty" orm:"column(username);size(255);null"`
	Password string `json:"password,omitempty" orm:"column(password);size(255);null"`
}

func (t *GameTestNumber) TableName() string {
	return "game_test_number"
}

func init() {
	orm.RegisterModel(new(GameTestNumber))
}

// AddGameTestNumber insert a new GameTestNumber into database and returns
// last inserted Id on success.
func AddGameTestNumber(m *GameTestNumber) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetGameTestNumberById retrieves GameTestNumber by Id. Returns error if
// Id doesn't exist
func GetGameTestNumberById(id int) (v *GameTestNumber, err error) {
	o := orm.NewOrm()
	v = &GameTestNumber{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllGameTestNumber retrieves all GameTestNumber matches certain condition. Returns empty list if
// no records exist
func GetAllGameTestNumber(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(GameTestNumber))
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

	var l []GameTestNumber
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

// UpdateGameTestNumber updates GameTestNumber by Id and returns error if
// the record to be updated doesn't exist
func UpdateGameTestNumberById(m *GameTestNumber) (err error) {
	o := orm.NewOrm()
	v := GameTestNumber{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteGameTestNumber deletes GameTestNumber by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGameTestNumber(id int) (err error) {
	o := orm.NewOrm()
	v := GameTestNumber{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&GameTestNumber{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
