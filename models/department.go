package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"errors"
)

type Department struct {
	Id   int    `json:"id,omitempty" orm:"column(id);auto"`
	Name string `json:"name,omitempty" orm:"column(name);size(50);null"`
}

func (t *Department) TableName() string {
	return "department"
}

func init() {
	orm.RegisterModel(new(Department))
}

// AddDepartment insert a new Department into database and returns
// last inserted Id on success.
func AddDepartment(m *Department) (id int64, err error) {
	m.Id = 0
	if m.Name == "" {
		err = errors.New("have't param Name")
		return
	}
	o := orm.NewOrm()
	o.Read(m, "Name")
	if m.Id != 0 {
		err = errors.New(m.Name + " department is existed")
		return
	}
	id, err = o.Insert(m)
	return
}

// GetDepartmentById retrieves Department by Id. Returns error if
// Id doesn't exist
func GetDepartmentById(id int) (v *Department, err error) {
	o := orm.NewOrm()
	v = &Department{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateDepartment updates Department by Id and returns error if
// the record to be updated doesn't exist
func UpdateDepartmentById(m *Department) (err error) {
	o := orm.NewOrm()
	v := Department{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteDepartment deletes Department by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDepartment(id int) (err error) {
	o := orm.NewOrm()
	v := Department{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Department{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllDepartmentMap() map[int]Department {
	o := orm.NewOrm()
	qs2 := o.QueryTable(new(Department))
	rs := []Department{}
	rsMap := map[int]Department{}
	_, err := qs2.All(&rs)
	if err != nil {
		return nil
	}
	for _, r := range rs {
		rsMap[r.Id] = r
	}

	return rsMap
}
func GetAllDepartment()([]*Department){
	o := orm.NewOrm()
	qs2 := o.QueryTable(new(Department))
	rs := []*Department{}
	qs2.All(&rs)
	return rs
}
