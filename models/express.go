package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Express struct {
	Id         int    `json:"id,omitempty" orm:"column(id);auto"`
	ContractId int    `json:"contract_id,omitempty" orm:"column(contract_id)"`
	//ExpressId1 int    `json:"express_id1,omitempty" orm:"column(express_id1);null"`
	Number1    int    `json:"number1,omitempty" orm:"column(number1);size(255);null"`
	//ExpressId2 int    `json:"express_id2,omitempty" orm:"column(express_id2);null"`
	Number2    int    `json:"number2,omitempty" orm:"column(number2);size(255);null"`
}

func (t *Express) TableName() string {
	return "express"
}

func init() {
	orm.RegisterModel(new(Express))
}

// AddExpress insert a new Express into database and returns
// last inserted Id on success.
func AddExpress(m *Express) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetExpressById retrieves Express by Id. Returns error if
// Id doesn't exist
func GetExpressById(id int) (v *Express, err error) {
	o := orm.NewOrm()
	v = &Express{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllExpress retrieves all Express matches certain condition. Returns empty list if
// no records exist
func GetAllExpress(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Express))
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

	var l []Express
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

// UpdateExpress updates Express by Id and returns error if
// the record to be updated doesn't exist
func UpdateExpressById(m *Express) (err error) {
	o := orm.NewOrm()
	//v := Express{ContractId: m.ContractId}
	// ascertain id exists in the database
	express := Express{}
	if err = o.QueryTable(new(Express)).Filter("contract_id__in", m.ContractId).One(&express); err == nil {
		var num int64
		m.Id = express.Id
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}else if err == orm.ErrNoRows{
		_, err = AddExpress(m)
		if err != nil{
			return
		}
	}
	return
}

// DeleteExpress deletes Express by Id and returns error if
// the record to be deleted doesn't exist
func DeleteExpress(id int) (err error) {
	o := orm.NewOrm()
	v := Express{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Express{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
