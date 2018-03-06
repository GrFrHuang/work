package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"time"
	"strconv"
)

type ContractNumber struct {
	Id         int    `json:"id,omitempty" orm:"column(id);auto"`
	ContractId int    `json:"contract_id,omitempty" orm:"column(contract_id);null"`
	Number     string `json:"number,omitempty" orm:"column(number);size(20);null"`
	CreateTime int64  `json:"create_time,omitempty" orm:"column(create_time);null"`
	FileId     int    `json:"file_id" orm:"column(file_id);null" valid:"Required"`
}

func (t *ContractNumber) TableName() string {
	return "contract_number"
}

func init() {
	orm.RegisterModel(new(ContractNumber))
}

// AddContractNumber insert a new ContractNumber into database and returns
// last inserted Id on success.
func AddContractNumber(m *ContractNumber) (id int64, err error) {
	o := orm.NewOrm()

	m.CreateTime = time.Now().Unix()
	qs := o.QueryTable(new(ContractNumber))
	var tmp []orm.Params
	qs.OrderBy("-id").Limit(1).Values(&tmp)
	if tmp == nil{
		m.Number = "KF1"
	}else{
		m.Number = "KF" + strconv.FormatInt(tmp[0]["Id"].(int64) + 1, 10)
	}
	id, err = o.Insert(m)

	return
}

// GetContractNumberById retrieves ContractNumber by Id. Returns error if
// Id doesn't exist
func GetContractNumberById(id int) (v *ContractNumber, err error) {
	o := orm.NewOrm()
	v = &ContractNumber{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllContractNumber retrieves all ContractNumber matches certain condition. Returns empty list if
// no records exist
func GetAllContractNumber(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ContractNumber))
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

	var l []ContractNumber
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

// UpdateContractNumber updates ContractNumber by Id and returns error if
// the record to be updated doesn't exist
func UpdateContractNumberById(m *ContractNumber) (err error) {
	o := orm.NewOrm()
	v := ContractNumber{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteContractNumber deletes ContractNumber by Id and returns error if
// the record to be deleted doesn't exist
func DeleteContractNumber(id int) (err error) {
	o := orm.NewOrm()
	v := ContractNumber{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ContractNumber{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
