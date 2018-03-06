package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Profit struct {
	Id          int     `json:"id,omitempty" orm:"column(id);auto"`
	GameId      int     `json:"game_id,omitempty" orm:"column(game_id);null"`
	ChannelCode string  `json:"channel_code,omitempty" orm:"column(channel_code);size(255);null"`
	Date        string  `json:"date,omitempty" orm:"column(date);size(10);null"`
	Amount      float64 `json:"amount,omitempty" orm:"column(amount);null"`
	Profit      float64 `json:"profit,omitempty" orm:"column(profit);null"`
	BodyMy      int     `json:"body_my,omitempty" orm:"column(body_my);null" description:"我方主体 1:云端 2:有量"`
}

func (t *Profit) TableName() string {
	return "profit"
}

func init() {
	orm.RegisterModel(new(Profit))
}

// AddProfit insert a new Profit into database and returns
// last inserted Id on success.
func AddProfit(m *Profit) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetProfitById retrieves Profit by Id. Returns error if
// Id doesn't exist
func GetProfitById(id int) (v *Profit, err error) {
	o := orm.NewOrm()
	v = &Profit{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllProfit retrieves all Profit matches certain condition. Returns empty list if
// no records exist
func GetAllProfit(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Profit))
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

	var l []Profit
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

// UpdateProfit updates Profit by Id and returns error if
// the record to be updated doesn't exist
func UpdateProfitById(m *Profit) (err error) {
	o := orm.NewOrm()
	v := Profit{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteProfit deletes Profit by Id and returns error if
// the record to be deleted doesn't exist
func DeleteProfit(id int) (err error) {
	o := orm.NewOrm()
	v := Profit{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Profit{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
