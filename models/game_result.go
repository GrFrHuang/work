package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type GameTest struct {
	Id                int    `json:"id,omitempty" orm:"column(id);auto"`
	OutId             int    `json:"out_id,omitempty" orm:"column(out_id);null"`
	Result            int    `json:"result,omitempty" orm:"column(result);null"`
	ResultPerson      int    `json:"result_person,omitempty" orm:"column(result_person);null"`
	ResultTime        int64  `json:"result_time,omitempty" orm:"column(result_time);null"`
	ResultReportWord  int    `json:"result_report_word,omitempty" orm:"column(result_report_word);null"`
	ResultReportExcel int    `json:"result_report_excel,omitempty" orm:"column(result_report_excel);null"`
	Package           string `json:"package,omitempty" orm:"column(package);null" `	//测试包
	Picture           string `json:"picture,omitempty" orm:"column(picture);null" `	//测试数据截图
	CreateTime        int64  `json:"create_time,omitempty" orm:"column(create_time);null"`
	UpdateTime        int64  `json:"update_time,omitempty" orm:"column(update_time);null"`
}

func (t *GameTest) TableName() string {
	return "game_test"
}

func init() {
	orm.RegisterModel(new(GameTest))
}

// AddGameTest insert a new GameTest into database and returns
// last inserted Id on success.
func AddGameTest(m *GameTest) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetGameTestById retrieves GameTest by Id. Returns error if
// Id doesn't exist
func GetGameTestById(id int) (v *GameTest, err error) {
	o := orm.NewOrm()
	v = &GameTest{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllGameTest retrieves all GameTest matches certain condition. Returns empty list if
// no records exist
func GetAllGameTest(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(GameTest))
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

	var l []GameTest
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

// UpdateGameTest updates GameTest by Id and returns error if
// the record to be updated doesn't exist
func UpdateGameTestById(m *GameTest) (err error) {
	o := orm.NewOrm()
	v := GameTest{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteGameTest deletes GameTest by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGameTest(id int) (err error) {
	o := orm.NewOrm()
	v := GameTest{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&GameTest{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
