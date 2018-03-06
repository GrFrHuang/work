package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type RemitDownDetail struct {
	Id              int     `orm:"column(id);auto" json:"id,omitempty"`
	RemitDownId     int     `orm:"column(remit_down_id);null" json:"remit_down_id,omitempty"`
	VerifyChannelId int     `orm:"column(verify_channel_id);null" json:"verify_channel_id,omitempty"`
	RemitMonth      string  `orm:"column(remit_month);size(7);null" json:"remit_month,omitempty"`
	RemitMoney      float64 `orm:"column(remit_money);null;digits(16);decimals(2)" json:"remit_money,omitempty"`
	RemitType       int     `orm:"column(remit_type);null" json:"remit_type,omitempty"`
}

func (t *RemitDownDetail) TableName() string {
	return "remit_down_detail"
}

func init() {
	orm.RegisterModel(new(RemitDownDetail))
}

// AddRemitDownDetail insert a new RemitDownDetail into database and returns
// last inserted Id on success.
func AddRemitDownDetail(m *RemitDownDetail) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRemitDownDetailById retrieves RemitDownDetail by Id. Returns error if
// Id doesn't exist
func GetRemitDownDetailById(id int) (v *RemitDownDetail, err error) {
	o := orm.NewOrm()
	v = &RemitDownDetail{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRemitDownDetail retrieves all RemitDownDetail matches certain condition. Returns empty list if
// no records exist
func GetAllRemitDownDetail(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RemitDownDetail))
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

	var l []RemitDownDetail
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

// UpdateRemitDownDetail updates RemitDownDetail by Id and returns error if
// the record to be updated doesn't exist
func UpdateRemitDownDetailById(m *RemitDownDetail) (err error) {
	o := orm.NewOrm()
	v := RemitDownDetail{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRemitDownDetail deletes RemitDownDetail by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRemitDownDetail(id int) (err error) {
	o := orm.NewOrm()
	v := RemitDownDetail{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&RemitDownDetail{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
