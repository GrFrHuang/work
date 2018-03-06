package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/utils"
	"sync"
)

type UserDesktopScreenLog struct {
	Id         int    `orm:"column(id);auto"`
	Uid        int    `orm:"column(uid);null"`
	Imageurl   string `orm:"column(imageurl);size(255);null"`
	Cycles     int    `orm:"column(cycles);null"`
	CreateTime int    `orm:"column(create_time);null"`
	NickName   string `orm:"-"`
	VoicePath  string `orm:"-"`
}

var syncWait sync.WaitGroup

func (t *UserDesktopScreenLog) TableName() string {
	return "user_desktop_screen_log"
}

func init() {
	orm.RegisterModel(new(UserDesktopScreenLog))
}

func GetWarnings() (error, []UserDesktopScreenLog) {
	o := orm.NewOrm()
	var info []UserDesktopScreenLog
	num, err := o.Raw("SELECT * FROM user_desktop_screen_log  WHERE cycles>0  ORDER BY create_time DESC").QueryRows(&info)
	if num <= 0 || err != nil {
		return errors.New("未查询到数据"), nil
	}
	syncWait.Add(len(info))
	for i := 0; i < len(info); i++ {
		go func(index int) {
			defer syncWait.Done()
			err, name := GetNickNameById(info[index].Uid)
			if err != nil {
				return
			}
			err, path := utils.VoiceByText(name + "存在异常操作")
			if err != nil {
				return
			}
			info[index].NickName = name
			info[index].VoicePath = path
			o.Raw("UPDATE user_desktop_screen_log SET cycles=cycles-1 WHERE cycles>0 AND  id=?", info[index].Id).Exec()
		}(i)

	}
	syncWait.Wait()
	return nil, info
}

// AddUserDesktopScreenLog insert a new UserDesktopScreenLog into database and returns
// last inserted Id on success.
func AddUserDesktopScreenLog(m *UserDesktopScreenLog) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserDesktopScreenLogById retrieves UserDesktopScreenLog by Id. Returns error if
// Id doesn't exist
func GetUserDesktopScreenLogById(id int) (v *UserDesktopScreenLog, err error) {
	o := orm.NewOrm()
	v = &UserDesktopScreenLog{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserDesktopScreenLog retrieves all UserDesktopScreenLog matches certain condition. Returns empty list if
// no records exist
func GetAllUserDesktopScreenLog(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserDesktopScreenLog))
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

	var l []UserDesktopScreenLog
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

// UpdateUserDesktopScreenLog updates UserDesktopScreenLog by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserDesktopScreenLogById(m *UserDesktopScreenLog) (err error) {
	o := orm.NewOrm()
	v := UserDesktopScreenLog{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserDesktopScreenLog deletes UserDesktopScreenLog by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserDesktopScreenLog(id int) (err error) {
	o := orm.NewOrm()
	v := UserDesktopScreenLog{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserDesktopScreenLog{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
