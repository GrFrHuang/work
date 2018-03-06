package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"time"
)

type Logs struct {
	Id          int    `json:"id,omitempty" orm:"column(id);auto"`
	UserId      int    `json:"user_id,omitempty" orm:"column(user_id);null"`
	Action      string `json:"action,omitempty" orm:"column(action);null"`
	StatusCode  int    `json:"status_code,omitempty" orm:"column(status_code);null"`
	Input       string `json:"input,omitempty" orm:"column(input);null"`
	CreatedTime int64 `json:"created_time,omitempty" orm:"column(created_time);null"`

	User *User `orm:"-" json:"user,omitempty"`
}

func (t *Logs) TableName() string {
	return "logs"
}

func init() {
	orm.RegisterModel(new(Logs))
}

// AddLogs insert a new Logs into database and returns
// last inserted Id on success.
func AddLogs(m *Logs) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// uid 谁,action 做了什么事,statusCode 做后的状态码(请求的返回码)
func SaveOperation(uid int, action string, statusCode int, input string) {
	log := &Logs{
		UserId:     uid,
		Action:     action,
		StatusCode: statusCode,
		Input:      input,
		CreatedTime:time.Now().Unix(),
	}
	_, err := AddLogs(log)
	if err != nil {
		logs.Error("save log error %v", err)
	}
}

// GetLogsById retrieves Logs by Id. Returns error if
// Id doesn't exist
func GetLogsById(id int) (v *Logs, err error) {
	o := orm.NewOrm()
	v = &Logs{Id: id}
	if err = o.Read(v); err != nil {
		return
	}

	u := &User{}
	o.QueryTable(u).Filter("Id", v.UserId).One(u, "Name", "Nickname")
	v.User = u
	return
}

// DeleteLogs deletes Logs by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLogs(id int) (err error) {
	// can not delete
	return errors.New("forbidden")
	o := orm.NewOrm()
	v := Logs{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Logs{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
