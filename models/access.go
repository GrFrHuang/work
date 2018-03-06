package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Access struct {
	Id   int    `json:"id,omitempty" orm:"column(id);pk"`
	Name string `json:"name,omitempty" orm:"column(name);size(50);null"`
	Url  string `json:"url,omitempty" orm:"column(url);size(255);null"`
}

func (t *Access) TableName() string {
	return "access"
}

func init() {
	orm.RegisterModel(new(Access))
}

// AddAccess insert a new Access into database and returns
// last inserted Id on success.
func AddAccess(m *Access) (id int64, err error) {
	m.Id = 0
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAccessById retrieves Access by Id. Returns error if
// Id doesn't exist
func GetAccessById(id int) (v *Access, err error) {
	o := orm.NewOrm()
	v = &Access{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateAccess updates Access by Id and returns error if
// the record to be updated doesn't exist
func UpdateAccessById(m *Access) (err error) {
	o := orm.NewOrm()
	v := Access{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAccess deletes Access by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAccess(id int) (err error) {
	o := orm.NewOrm()
	v := Access{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Access{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
