package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"github.com/astaxie/beego/orm"
	"database/sql"
	_"github.com/astaxie/beego/session/mysql"
)

type RepairData struct {
	Id        int       `orm:"column(id);auto"`
	Cp        string    `orm:"column(cp);size(255);null"`
	StartTime string    `orm:"column(start_time);type(date);null"`
	EndTime   string    `orm:"column(end_time);type(date);null"`
	//GameId    []int       `orm:"column(game_id);null"`
	GameId int          `orm:"column(game_id);null"`
}

type GameId int

func (t *RepairData) TableName() string {
	return "repair_data"
}

func init() {
	orm.RegisterModel(new(RepairData))
}

// AddRepairData insert a new RepairData into database and returns
// last inserted Id on success.
func AddRepairData(m *RepairData) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRepairDataById retrieves RepairData by Id. Returns error if
// Id doesn't exist
func GetRepairDataById(id int) (v *RepairData, err error) {
	o := orm.NewOrm()
	v = &RepairData{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRepairData retrieves all RepairData matches certain condition. Returns empty list if
// no records exist
func GetAllRepairData(query map[string]string, fields []string, sortby []string, order []string,
		offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RepairData))
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

	var l []RepairData
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

// UpdateRepairData updates RepairData by Id and returns error if
// the record to be updated doesn't exist
func UpdateRepairDataById(m *RepairData) (err error) {
	o := orm.NewOrm()
	v := RepairData{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRepairData deletes RepairData by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRepairData(id int) (err error) {
	o := orm.NewOrm()
	v := RepairData{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&RepairData{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}


//clear table
func ClearTable() error{
	db, err := sql.Open("mysql", "kftest:123456@tcp(10.8.230.17:3308)/work_together?charset=utf8")
	if err != nil{
		return err
	}
	stmt,err := db.Prepare("Truncate Table repair_data")
	result ,err := stmt.Exec()
	if result == nil {
		return err
	}
	return err
}