package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type WarningLog struct {
	Id          int     `json:"id,omitempty" orm:"column(id);auto"`
	WarningName string  `json:"warning_name,omitempty" orm:"column(warning_name);size(255);null"`
	WarningType string  `json:"warning_type,omitempty" orm:"column(warning_type);size(255);null"`
	Info        string  `json:"info,omitempty" orm:"column(info);null"`
	GameId      int     `json:"game_id,omitempty" orm:"column(game_id);null"`
	ChannelCode string  `json:"channel_code,omitempty" orm:"column(channel_code);size(255);null"`
	Grade       int     `json:"grade,omitempty" orm:"column(grade);null"`
	CreateTime  int64   `json:"create_time,omitempty" orm:"column(create_time);null"`
	Keys        string  `json:"keys,omitempty" orm:"column(keys);size(255);null"`
	Amount      float64 `json:"amount,omitempty" orm:"column(amount);digits(11);decimals(2);null"`
	Date        string  `json:"date,omitempty" orm:"column(date);size(255);null"`

	GameName    string  `orm:"-" json:"game_name,omitempty"`
	ChannelName string  `orm:"-" json:"channel_name,omitempty"`
}

func (t *WarningLog) TableName() string {
	return "warning_log"
}

func init() {
	orm.RegisterModel(new(WarningLog))
}

// AddWarningLog insert a new WarningLog into database and returns
// last inserted Id on success.
func AddWarningLog(m *WarningLog) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetWarningLogById retrieves WarningLog by Id. Returns error if
// Id doesn't exist
func GetWarningLogById(id int) (v *WarningLog, err error) {
	o := orm.NewOrm()
	v = &WarningLog{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllWarningLog retrieves all WarningLog matches certain condition. Returns empty list if
// no records exist
func GetAllWarningLog(query map[string]string, fields []string, sortby []string, order []string,
offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(WarningLog))
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

	var l []WarningLog
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

// UpdateWarningLog updates WarningLog by Id and returns error if
// the record to be updated doesn't exist
func UpdateWarningLogById(m *WarningLog) (err error) {
	o := orm.NewOrm()
	v := WarningLog{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteWarningLog deletes WarningLog by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWarningLog(id int) (err error) {
	o := orm.NewOrm()
	v := WarningLog{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&WarningLog{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
