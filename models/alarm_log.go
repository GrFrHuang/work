package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type AlarmLog struct {
	Id          int    `json:"id,omitempty" orm:"column(id);auto"`
	CreatedTime int    `json:"created_time,omitempty" orm:"column(created_time);null"`
	IsHide      int    `json:"is_hide,omitempty" orm:"column(is_hide);null"`
	Context     string `json:"context,omitempty" orm:"column(context);null"`
	Type        string `json:"type,omitempty" orm:"column(type);size(255);null"`
	Extend      string `json:"extend,omitempty" orm:"column(extend);null"`
}

func (t *AlarmLog) TableName() string {
	return "alarm_log"
}

func init() {
	orm.RegisterModel(new(AlarmLog))
}

// AddAlarmLog insert a new AlarmLog into database and returns
// last inserted Id on success.
func AddAlarmLog(m *AlarmLog) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAlarmLogById retrieves AlarmLog by Id. Returns error if
// Id doesn't exist
func GetAlarmLogById(id int) (v *AlarmLog, err error) {
	o := orm.NewOrm()
	v = &AlarmLog{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateAlarmLog updates AlarmLog by Id and returns error if
// the record to be updated doesn't exist
func UpdateAlarmLogById(m *AlarmLog) (err error) {
	o := orm.NewOrm()
	v := AlarmLog{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAlarmLog deletes AlarmLog by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAlarmLog(id int) (err error) {
	o := orm.NewOrm()
	v := AlarmLog{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AlarmLog{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
