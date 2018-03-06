package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/utils"
)

type ExpressManage struct {
	Id             int    `orm:"column(id);auto" json:"id"`
	Number         string `orm:"column(number);size(255)" json:"number"`
	ReceiveCompany int    `orm:"column(receive_company)" json:"receive_company,omitempty"`
	ReceiveAddress string `orm:"column(receive_address)" json:"receive_address,omitempty"`
	BodyMy         int    `orm:"column(body_my)" json:"body_my,omitempty"`
	SendDepartment int    `orm:"column(send_department)" json:"send_department,omitempty"`
	SendPeople     int    `orm:"column(send_people)" json:"send_people,omitempty"`
	ContentId      int    `orm:"column(content_id)" json:"content_id,omitempty"`
	IncludeGames   string `orm:"column(include_games);size(255)" json:"include_games,omitempty"`
	Channels       string `orm:"column(channels);size(255)" json:"channels,omitempty"`
	Details        string `orm:"column(details);size(255)" json:"details,omitempty"`
	SendType       int    `orm:"column(send_type)" json:"send_type,omitempty"`
	ReceivePeople  int    `orm:"column(receive_people)" json:"receive_people,omitempty"`
	CreateTime     int64  `orm:"column(create_time)" json:"create_time,omitempty"`
	CreatePeople   int    `orm:"column(create_people)" json:"create_people,omitempty"`
	TypeCompany    int    `json:"type_company,omitempty" orm:"column(type_company)"`
	Icon           string `json:"icon,omitempty" orm:"column(icon)"`

	SenderNickName  string `json:"sender_nick_name,omitempty" orm:"-"`
	ReceiveUserName string `json:"receive_user_nick_name,omitempty" orm:"-"`
	CompanyName     string `json:"company_name,omitempty" orm:"-"`

	SendUser    *User        `json:"send_user"  orm:"-"`
	Department  *Department  `json:"department,omitempty"  orm:"-"`
	Company     *CompanyType `json:"company,omitempty"  orm:"-"`
	Types       *Types       `json:"type,omitempty"  orm:"-"`
	ReceiveUser *Contact     `json:"receive_user,omitempty"  orm:"-"`
	CreateUser  *User        `json:"create_user,omitempty"  orm:"-"`
}

func (t *ExpressManage) TableName() string {
	return "express_manage"
}

func init() {
	orm.RegisterModel(new(ExpressManage))
}

// AddExpressManage insert a new ExpressManage into database and returns
// last inserted Id on success.
func AddExpressManage(m *ExpressManage) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)

	err = CompareAndAddOperateLog(nil, m, m.CreatePeople, bean.OPP_EXPRESS_MANAGE, int(id), bean.OPA_INSERT)
	return
}

// GetExpressManageById retrieves ExpressManage by Id. Returns error if
// Id doesn't exist
func GetExpressManageById(id int) (v *ExpressManage, err error) {
	o := orm.NewOrm()
	v = &ExpressManage{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//检验是否 已存在 该快递
func CheckExpressManageByNumber(number string) (err error) {
	o := orm.NewOrm()
	express := ExpressManage{}
	qs := o.QueryTable(new(ExpressManage))
	err = qs.Filter("Number__exact", number).One(&express)
	return
}

// 获取 number id 信息 返回[]express对象 仅有number id信息
func GetExpressNumbers() (expresses []ExpressManage, err error) {
	o := orm.NewOrm()
	if _, err = o.QueryTable(new(ExpressManage)).OrderBy("-id").All(&expresses, "Id", "Number"); err != nil {
		return
	}
	return expresses, nil
}

// GetAllExpressManage retrieves all ExpressManage matches certain condition. Returns empty list if
// no records exist
func GetAllExpressManage(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ExpressManage))
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

	var l []ExpressManage
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

// UpdateExpressManage updates ExpressManage by Id and returns error if
// the record to be updated doesn't exist
func UpdateExpressManageById(m *ExpressManage, where map[string][]interface{}, userid int) (err error) {
	o := orm.NewOrm()
	v := ExpressManage{Id: m.Id}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}
	fields := utils.GetNotEmptyFields(m, "Number", "TypeCompany", "ReceiveCompany", "ReceiveAddress", "ReceivePeople",
		"BodyMy", "SendType", "ContentId", "IncludeGames", "Channels", "Details", "SendDepartment", "SendPeople")

	if _, err = o.Update(m, fields...); err != nil {
		return
	}

	if err = CompareAndAddOperateLog(&v, m, userid, bean.OPP_EXPRESS_MANAGE, v.Id, bean.OPA_UPDATE, fields...); err != nil {
		return
	}

	// ascertain id exists in the database
	//if err = o.Read(&v); err == nil {
	//	//var num int64
	//	if _, err = o.Update(m); err == nil {
	//		//fmt.Println("Number of records updated in database:", num)
	//		return nil
	//	}
	//}
	return
}

// DeleteExpressManage deletes ExpressManage by Id and returns error if
// the record to be deleted doesn't exist
func DeleteExpressManage(id int) (err error) {
	o := orm.NewOrm()
	v := ExpressManage{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ExpressManage{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//添加收件公司的信息
func AddReceiverCompanyInfo(ss *[]ExpressManage) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	links := make([]int, len(*ss))
	for i, s := range *ss {
		links[i] = s.ReceiveCompany
	}
	companys := []CompanyType{}
	o := orm.NewOrm()
	qs := o.QueryTable(new(CompanyType))
	if _, err := qs.Filter("Id__in", links).All(&companys, "Id", "Name"); err != nil {
		return
	}
	companyMap := map[int]CompanyType{}
	for _, c := range companys {
		companyMap[c.Id] = c
	}
	for i, s := range *ss {
		if g, ok := companyMap[s.ReceiveCompany]; ok {
			(*ss)[i].Company = &g
		}
	}
	return

}

//添加发件人信息
func AddSendUser(ss *[]ExpressManage) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	links := make([]int, len(*ss))
	for i, s := range *ss {
		links[i] = s.SendPeople
	}
	receivers := []User{}
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	if _, err := qs.Filter("Id__in", links).All(&receivers, "Id", "NickName"); err != nil {
		return
	}
	receiverMap := map[int]User{}
	for _, u := range receivers {
		receiverMap[u.Id] = u
	}
	for i, s := range *ss {
		if g, ok := receiverMap[s.SendPeople]; ok {
			(*ss)[i].SendUser = &g
		}
	}
	return
}

//添加发件人所在的部门信息
func AddSendDepartment(ss *[]ExpressManage) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	links := make([]int, len(*ss))
	for i, s := range *ss {
		links[i] = s.SendDepartment
	}
	departments := []Department{}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Department))
	if _, err := qs.Filter("Id__in", links).All(&departments, "Id", "Name"); err != nil {
		return
	}
	departmentMap := map[int]Department{}
	for _, d := range departments {
		departmentMap[d.Id] = d
	}
	for i, s := range *ss {
		if g, ok := departmentMap[s.SendDepartment]; ok {
			(*ss)[i].Department = &g
		}
	}
	return
}

//添加邮寄内容信息
func AddSendTypeInfo(ss *[]ExpressManage) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	links := make([]int, len(*ss))
	for i, s := range *ss {
		links[i] = s.SendType
	}
	types := []Types{}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Types))
	if _, err := qs.Filter("Id__in", links).All(&types, "Id", "Name"); err != nil {
		return
	}
	typeMap := map[int]Types{}
	for _, t := range types {
		typeMap[t.Id] = t
	}
	for i, s := range *ss {
		if g, ok := typeMap[s.SendType]; ok {
			(*ss)[i].Types = &g
		}
	}
	return
}

//添加收件人信息 仅有id 以及 name信息
func GetReceiverByContactId(ss *[]ExpressManage) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	links := make([]int, len(*ss))
	for i, s := range *ss {
		links[i] = s.ReceivePeople
	}
	types := []Contact{}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Contact))
	if _, err := qs.Filter("Id__in", links).All(&types, "Id", "Name"); err != nil {
		return
	}
	typeMap := map[int]Contact{}
	for _, t := range types {
		typeMap[t.Id] = t
	}
	for i, s := range *ss {
		if g, ok := typeMap[s.ReceivePeople]; ok {
			(*ss)[i].ReceiveUser = &g
		}
	}
	return
}

func GetDataBySendType(send_type int) (int64, []ExpressManage, error) {
	e := []ExpressManage{}
	t := Types{}
	u := User{}
	c := Contact{}
	o := orm.NewOrm()
	sql := "SELECT a.* from express_manage as a, types as b  WHERE " +
		"type = 9 AND content_id=2 and b.id=? order by a.create_time desc"
	total, err := o.Raw(sql, send_type).QueryRows(&e)

	for i := 0; i < len(e); i++ {
		sql = "SELECT * from types where id=?"
		o.Raw(sql, e[i].SendType).QueryRow(&t)
		e[i].Types = &t

		sql = "SELECT * from `user` where id=?"
		o.Raw(sql, e[i].SendPeople).QueryRow(&u)
		e[i].SendUser = &u

		sql = "SELECT * FROM contact WHERE id =?"
		o.Raw(sql, e[i].ReceivePeople).QueryRow(&c)
		e[i].ReceiveUser = &c

		sql = "SELECT * FROM user WHERE id =?"
		o.Raw(sql, e[i].CreatePeople).QueryRow(&u)
		e[i].CreateUser = &u
	}

	return total, e, err

}
