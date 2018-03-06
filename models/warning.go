package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
)

type Warning struct {
	Id          int     `json:"id,omitempty" orm:"column(id);auto"`
	Name        string  `json:"name,omitempty" orm:"column(name);size(200);null"`
	Type        string  `json:"type,omitempty" orm:"column(type);size(200);null"`
	Intervals   string  `json:"intervals,omitempty" orm:"column(intervals);size(255);null"`
	ChannelCode string  `json:"channel_code,omitempty" orm:"column(channel_code);null"`
	GameId      string  `json:"game_id,omitempty" orm:"column(game_id);null"`
	Amount      float64 `json:"amount,omitempty" orm:"column(amount);null;digits(10);decimals(2)"`
	Grade       int     `json:"grade,omitempty" orm:"column(grade);null"`
	UserIds     string  `json:"user_ids,omitempty" orm:"column(user_ids);null"`
	Emails      string  `json:"emails,omitempty" orm:"column(emails);null"`
	State       int8    `json:"state,omitempty" orm:"column(state);null"`
	IsRepeat    int8    `json:"is_repeat,omitempty" orm:"column(is_repeat);null"`
	IsPushHome  int8    `json:"is_push_home,omitempty" orm:"column(is_push_home);null"`
	IsSendMail  int8    `json:"is_send_mail,omitempty" orm:"column(is_send_mail);null"`
	CreateTime  int64   `json:"create_time,omitempty" orm:"column(create_time);null"`
}

const (
	WARNING_GAME_ACCESS                 = "新游接入"
	WARNING_GAME_RELEASE_UPDATE         = "首发变更"
	WARNING_GAME_UPDATE                 = "游戏更新"
	WARNING_CP_CONTRACT_EXPIRE          = "CP合同到期"
	WARNING_CHANNEL_CONTRACT_EXPIRE     = "渠道合同到期"
	WARNING_CHANNEL_CONTRACT_SIGN       = "渠道合同未签订"
	WARNING_CP_ORDER_VERIFY             = "CP未对账"
	WARNING_CHANNEL_ORDER_VERIFY        = "渠道未对账"
	WARNING_CHANNEL_ORDER_PAY           = "渠道未回款"
	WARNING_CHANNEL_ORDER_THRESHOLD     = "渠道流水&阈值"
	WARNING_GAME_PROFIT                 = "分成比例亏损"
	WARNING_CHANNEL_SECURITY_DEPOSIT    = "保证金预警"
	WARNING_GAME_OUTAGE                 = "游戏停运公告"
	WARNING_GAME_OUTAGE_CP_OPERATE      = "停运游戏CP合同待处理"
	WARNING_GAME_OUTAGE_CHANNEL_OPERATE = "停运游戏渠道合同待处理"
	WARNING_GAME_OUTAGE_CP_DOWN         = "停运游戏CP合同未下架"
	WARNING_GAME_OUTAGE_CHANNEL_DOWN    = "停运游戏渠道合同未下架"

	WARNING_CHANNEL_MAIN_CONTRACT_EXPIRE = "渠道主合同到期"
	WARNING_CP_MAIN_CONTRACT_EXPIRE      = "cp主合同到期"
)

var Warning_types = map[int]string{
	1:  "新游接入",
	2:  "首发变更",
	3:  "游戏更新",
	4:  "渠道合同到期",
	5:  "渠道合同未签订",
	6:  "渠道未对账",
	7:  "渠道未回款",
	8:  "渠道流水&阈值",
	9:  "分成比例亏损",
	10: "保证金预警",
}

func (t *Warning) TableName() string {
	return "warning"
}

func init() {
	orm.RegisterModel(new(Warning))
}

// AddWarning insert a new Warning into database and returns
// last inserted Id on success.
func AddWarning(m *Warning) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetWarningById retrieves Warning by Id. Returns error if
// Id doesn't exist
func GetWarningById(id int) (v *Warning, err error) {
	o := orm.NewOrm()
	if err = o.Read(v); err == nil {
		return v, nil
	}
	v = &Warning{Id: id}
	return nil, err
}

// GetWarningByType retrieves Warning by type. Returns error if
// type doesn't exist
func GetWarningByType(tp string) (as []Warning, err error) {
	o := orm.NewOrm()
	v := &Warning{}
	as = []Warning{}
	_, err = o.QueryTable(v).Filter("Type", tp).All(&as)
	return
}

type Wtype struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func GetWarningStaticType() (types []Wtype) {
	types = []Wtype{}
	for k, v := range Warning_types {
		tmp := Wtype{Id: k, Name: v}
		types = append(types, tmp)
	}
	return
}

func GetAllWarningName() (result []Warning) {
	o := orm.NewOrm()
	o.QueryTable(new(Warning)).All(&result, "id", "type")

	return
}

func GetWarningDetail(ids []string) (result []Warning) {
	o := orm.NewOrm()
	o.QueryTable(new(Warning)).Filter("id__in", ids).All(&result, "id", "type", "user_ids")

	return
}

// GetAllWarning retrieves all Warning matches certain condition. Returns empty list if
// no records exist
func GetAllWarning(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Warning))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, v == "true" || v == "1")
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

	var l []Warning
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

// UpdateWarning updates Warning by Id and returns error if
// the record to be updated doesn't exist
func UpdateWarningById(m *Warning) (err error) {
	fmt.Println("access put", m)
	o := orm.NewOrm()
	v := Warning{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		fields := util.NewFieldsUtil(m).GetNotEmptyFields().
			Filter("Intervals", "Amount", "Position", "Grade", "State", "IsRepeat", "UserIds").
			Must("State", "IsRepeat").
			Fields()
		if num, err = o.Update(m, fields...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}

	return
}

// DeleteWarning deletes Warning by Id and returns error if
// the record to be deleted doesn't exist
func DeleteWarning(id int) (err error) {
	o := orm.NewOrm()
	v := Warning{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Warning{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func ParseWarningGrade(grade int) (color string) {
	switch grade {
	case 1:
		color = "蓝色"
	case 2:
		color = "黄色"
	case 3:
		color = "橙色"
	case 4:
		color = "红色"
	default:
		color = "ERROR"
	}
	return
}
