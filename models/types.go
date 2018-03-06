package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/kataras/go-errors"
	"strings"
)

type Types struct {
	Id          int    `json:"id,omitempty" orm:"column(id);auto"`
	Name        string `json:"name" orm:"column(name);size(255);null"`
	ISDeleted   int    `json:"is_deleted,omitempty" orm:"column(is_deleted)"`
	IsDeletable bool   `json:"is_deletable" orm:"column(is_deletable)"`
	Type        int    `json:"type,omitempty" orm:"column(type)"`

	CompanyIds string `json:"company_ids" orm:"-"`
}

func (t *Types) TableName() string {
	return "types"
}

func init() {
	orm.RegisterModel(new(Types))
}

// 添加该种类型下所有的主合同company_id
// companyType:  1--cp  2--渠道
func AddCompanyIdInfo(ss *[]Types, companyType int) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()

	for i, s := range *ss {
		var company_ids []orm.Params
		var ids []string
		o.Raw("SELECT company_id FROM main_contract WHERE company_type=? AND state=?", companyType, s.Id).
			Values(&company_ids)
		for _, id := range company_ids {
			ids = append(ids, id["company_id"].(string))
		}
		if len(ids) > 0 {
			(*ss)[i].CompanyIds = strings.Join(ids, ",")
		} else {
			(*ss)[i].CompanyIds = "-1"
		}
	}
	return
}

func AddTypes(m *Types) (int64, error) {
	o := orm.NewOrm()
	m_tmp := Types{Name: m.Name, Type: m.Type}
	err := o.Read(&m_tmp, "Name", "Type")
	if err == orm.ErrNoRows {
		id, err := o.Insert(m)
		return id, err
	} else if err == nil {
		if m_tmp.ISDeleted == 1 {
			m_tmp.ISDeleted = 0
			if n, err := o.Update(&m_tmp, "is_deleted"); err != nil {
				return n, err
			} else {
				return 1, nil
			}
		}
	}
	return 0, errors.New("该类型已经存在")
}

func DeleteTypes(id int, where map[string][]interface{}) error {
	o := orm.NewOrm()
	m := Types{Id: id}
	qs := QueryTable(&m, where).Filter("Id", id)
	err1 := qs.One(&m)
	if err1 == nil {
		m.ISDeleted = 1
		if _, err := o.Update(&m, "is_deleted"); err != nil {
			return err
		}
	}
	return err1
}

func GetTypesIdByName(name string) (id int, err error) {
	o := orm.NewOrm()
	id = 0
	m := Types{Name: name}
	if err = o.Read(m); err != nil {
		return
	}
	id = m.Id
	return
}

func GetTypesNameById(id int) (name string) {
	o := orm.NewOrm()
	m := Types{Id: id}
	if err := o.Read(&m); err != nil {
		return ""
	}
	return m.Name
}

func GetAllCourier() (t []Types) {
	o := orm.NewOrm()
	sql := "SELECT * from types where type = 9"
	o.Raw(sql).QueryRows(&t)
	return
}
