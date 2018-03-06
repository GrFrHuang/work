package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/utils"
)

type Contact struct {
	Id          int    `json:"id" orm:"column(id);auto"`
	CompanyId   int    `json:"company_id" orm:"column(company_id)"`
	Name        string `json:"name" orm:"column(name);size(255);null"`
	Sex         int    `json:"sex" orm:"column(sex);null"`
	Position    string `json:"position" orm:"column(position);size(255);null"`
	Email       string `json:"email" orm:"column(email);size(63);null"`
	Phone       string `json:"phone" orm:"column(phone);size(63);null"`
	Qq          string `json:"qq" orm:"column(qq);size(63);null"`
	Wechart     string `json:"wechart" orm:"column(wechart);size(63);null"`
	CompanyType int    `json:"company_type" orm:"column(company_type);null"`
}

func (t *Contact) TableName() string {
	return "contact"
}

func init() {
	orm.RegisterModel(new(Contact))
}

// AddContact insert a new Contact into database and returns
// last inserted Id on success.
func AddContact(m *Contact) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetContactById retrieves Contact by Id. Returns error if
// Id doesn't exist
func GetContactById(id int) (v *Contact, err error) {
	o := orm.NewOrm()
	v = &Contact{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetContactsByCompanyIds get all contact by company id, Returns
// all contact
func GetContactsByCompanyId(companyId int, companyType int) (contacts []Contact) {
	o := orm.NewOrm()
	_, err := o.QueryTable("contact").
		Filter("company_id", companyId).
		Filter("company_type", companyType).
		OrderBy("id").All(&contacts)
	if err != nil {
		contacts = []Contact{}
	}
	return
}

// DeleteContactsByCompanyId delete the contacts by company id, returns
// error if the errors exist
func DeleteContactsByCompanyId(companyId int) (err error) {
	o := orm.NewOrm()
	var contacts []*Contact
	_, err = o.QueryTable("contact").Filter("company_id", companyId).All(&contacts, "id")
	if err == nil {
		for _, v := range contacts {
			DeleteContact(v.Id)
		}
	}
	return
}

// GetAllContact retrieves all Contact matches certain condition. Returns empty list if
// no records exist
func GetAllContact(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Contact))
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

	var l []Contact
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

// UpdateContact updates Contact by Id and returns error if
// the record to be updated doesn't exist
func UpdateContactById(m *Contact) (err error) {
	//fmt.Println("Debug Info: contacts =", m)
	o := orm.NewOrm()
	v := Contact{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		fields := utils.GetNotEmptyFields(m, "Name", "Sex", "Position", "Email", "Phone", "Qq", "Wechart")
		//fmt.Println("Debug Info: contacts fileds =", fields)
		if num, err = o.Update(m, fields...); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// UpdateOrAddCompanyContacts Update Contact by Contacts , companyId and company type,
// returns error if the record to be updated doesn't exist
func UpdateOrAddCompanyContacts(contactNew []Contact, companyId int, companyType int) (err error) {
	if len(contactNew) == 0 {
		return
	}
	o := orm.NewOrm()

	// 查询该公司当前的所有联系人
	var contactOld []*Contact
	o.QueryTable("contact").Filter("company_id", companyId).Filter("company_type", companyType).All(&contactOld)
	idOld := map[int]bool{}
	for _, v := range contactOld {
		idOld[v.Id] = true
	}

	idNew := map[int]bool{}
	// 增加数据库中不存在的回款主体
	for _, v := range contactNew {
		if idOld[v.Id] == false {
			v.CompanyId = companyId
			v.CompanyType = companyType
			if _, err = AddContact(&v); err != nil {
				return
			}
		} else {
			UpdateContactById(&v)
		}
		idNew[v.Id] = true
	}

	// 删除数据库中多余的回款主体
	for _, v := range contactOld {
		if idNew[v.Id] == false {
			DeleteContact(v.Id)
		}
	}

	return
}

// DeleteContact deletes Contact by Id and returns error if
// the record to be deleted doesn't exist
func DeleteContact(id int) (err error) {
	o := orm.NewOrm()
	v := Contact{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Contact{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// 通过companyId 获得联系人信息
func GetContactByCompanyId(id , flag int) (contacts []Contact, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Contact))
	if _, err = qs.Filter("company_id", id).Filter("company_type",flag).All(&contacts, "Id", "Name"); err != nil {
		return
	}
	return contacts, nil
}
