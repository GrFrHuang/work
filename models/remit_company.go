package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type RemitCompany struct {
	Id          int    `json:"id,omitempty" orm:"column(id);auto"`
	ChannelCode string `json:"channel_code,omitempty" orm:"column(channel_code);size(100)"`
	CompanyId   int    `json:"company_id,omitempty" orm:"column(company_id)"`
}

func (t *RemitCompany) TableName() string {
	return "remit_company"
}

func init() {
	orm.RegisterModel(new(RemitCompany))
}

// AddRemitCompany insert a new RemitCompany into database and returns
// last inserted Id on success.
func AddRemitCompany(m *RemitCompany) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRemitCompanyById retrieves RemitCompany by Id. Returns error if
// Id doesn't exist
func GetRemitCompanyById(id int) (v *RemitCompany, err error) {
	o := orm.NewOrm()
	v = &RemitCompany{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRemitCompany retrieves all RemitCompany matches certain condition. Returns empty list if
// no records exist
func GetAllRemitCompany(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RemitCompany))
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

	var l []RemitCompany
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

// 获取所有或者指定回款主体的公司信息
func GetRemitCompaniesByChannel(channelCode string) (companies []CompanyType, err error) {
	companies = []CompanyType{}

	cs := []RemitCompany{}
	qs := orm.NewOrm().QueryTable("remit_company")
	if channelCode != "" && channelCode != "*" {
		qs = qs.Filter("channel_code", channelCode)
	}
	_, err = qs.All(&cs)
	if err != nil {
		return
	}
	if len(cs) == 0 {
		return
	}
	allRemitCompanyIds := []int{}
	for _, c := range cs {
		allRemitCompanyIds = append(allRemitCompanyIds,c.CompanyId)
	}
	if len(allRemitCompanyIds) == 0 {
		return
	}
	_, err = orm.NewOrm().QueryTable(new(CompanyType)).
		Filter("id__in", allRemitCompanyIds).
		All(&companies, "id", "name")
	if err != nil {
		return
	}

	return
}



// UpdateRemitCompany updates RemitCompany by Id and returns error if
// the record to be updated doesn't exist
func UpdateRemitCompanyById(m *RemitCompany) (err error) {
	o := orm.NewOrm()
	v := RemitCompany{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRemitCompany deletes RemitCompany by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRemitCompany(id int) (err error) {
	o := orm.NewOrm()
	v := RemitCompany{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&RemitCompany{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
