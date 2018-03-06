package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/utils"
)

type DevelopCompany struct {
	Id           int           `json:"id,omitempty" orm:"column(id);auto"`
	CompanyId    int           `json:"company_id,omitempty" orm:"column(company_id);null"`
	//PlatformName string        `json:"platform_name,omitempty" orm:"column(platform_name);size(255);null"`       // 平台名称
	Region       string        `json:"region,omitempty" orm:"column(region);size(255);null"`                     // 区域
	Website      string        `json:"website,omitempty" orm:"column(website);size(255);null"`                   // 公司网站
	Address      string        `json:"address,omitempty" orm:"column(address);size(255);null"`                   // 公司地址, e.g. ["address1", "address2"]
	ContactIds   string        `json:"contact_ids,omitempty" orm:"column(contact_ids);size(255);null"` // 联系人, e.g. [id1, id2, ...]
	Desc         string        `json:"desc,omitempty" orm:"column(desc)"`                                        // 备注
	UpdateTime   int64         `json:"update_time,omitempty" orm:"column(update_time);null" `
	UpdateUserID int           `json:"update_user_id,omitempty" orm:"column(update_user_id);"` // 更新人ID

	CompanyName    string         `orm:"-" json:"company_name,omitempty"`
	Contacts       []Contact     `orm:"-" json:"contacts,omitempty"`
	ContactsForCompare        string        `orm:"-" json:"contact_for_compare,omitempty"`	//仅用于操作日志比较
	UpdateUserName string         `orm:"-" json:"update_user_name,omitempty"`
}

func (t *DevelopCompany) TableName() string {
	return "develop_company"
}

func init() {
	orm.RegisterModel(new(DevelopCompany))
}

// AddDevelopCompany insert a new DevelopCompany into database and returns
// last inserted Id on success.
func AddDevelopCompany(m *DevelopCompany) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	if err != nil {
		return
	}

	for _, v := range m.Contacts {
		v.CompanyId = int(id)
		v.CompanyType = 1
		AddContact(&v)
	}
	err = CompareAndAddOperateLog(nil, m, m.UpdateUserID, bean.OPP_DEVELOP_COMPANY, int(id), bean.OPA_INSERT)
	return
}

func AddDevelopCompanyInfo(ss *[]DevelopCompany) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.CompanyId
	}
	games := []CompanyType{}
	_, err := o.QueryTable(new(CompanyType)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]CompanyType{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.CompanyId]; ok {
			(*ss)[i].CompanyName = g.Name
		}
	}
	return
}

// GetDevelopCompanyById retrieves DevelopCompany by Id. Returns error if
// Id doesn't exist
func GetDevelopCompanyById(id int) (v *DevelopCompany, err error) {
	o := orm.NewOrm()
	v = &DevelopCompany{Id: id}
	if err = o.Read(v); err != nil {
		return
	}
	GetDevelopCompanyAdditionalInfo(v)
	return
}

func GetDevelopCompanyAdditionalInfo(v *DevelopCompany) {
	// Add the company name
	c, err := GetCompanyTypeById(v.CompanyId)
	if err == nil {
		v.CompanyName = c.Name
	}

	// Add the update user name
	user, err := GetUserById(v.UpdateUserID)
	if err == nil {
		v.UpdateUserName = user.Nickname
	}

	// Add the contacts info
	v.Contacts = GetContactsByCompanyId(v.Id, 1)
	return
}

// UpdateDevelopCompany updates DevelopCompany by Id and returns error if
// the record to be updated doesn't exist
func UpdateDevelopCompanyById(m *DevelopCompany, where map[string][]interface{}) (err error) {
	//fmt.Println("Debug Info: access the UpdateDevelopCompanyById, DevelopCompany=", m.Company)
	//o := orm.NewOrm()
	v := DevelopCompany{Id: m.Id}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}

	oldContacts, _ := json.Marshal(GetContactsByCompanyId(v.Id, 1))
	v.ContactsForCompare = string(oldContacts)

	// Update the contacts info
	UpdateOrAddCompanyContacts(m.Contacts, m.Id, 1)

	o := orm.NewOrm()
	fields := util.NewFieldsUtil(m).GetNotEmptyFields().Must("Desc", "Website", "Region", "Address").
		Exclude("CompanyName", "Contacts", "UpdateUserName").Fields()
	if _, err = o.Update(m, fields...); err != nil {
		return
	}

	newContacts, _ := json.Marshal(m.Contacts)
	m.ContactsForCompare = string(newContacts)

	fields = append(fields, "ContactsForCompare")
	fields = utils.RemoveFields(fields, "UpdateUserID", "UpdateTime")

	if err = CompareAndAddOperateLog(&v, m, m.UpdateUserID, bean.OPP_DEVELOP_COMPANY, v.Id, bean.OPA_UPDATE, fields...); err != nil{
		return
	}
	return
}

// DeleteDevelopCompany deletes DevelopCompany by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDevelopCompany(id int) (err error) {
	o := orm.NewOrm()
	v := DevelopCompany{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		if _, err = o.Delete(&DevelopCompany{Id: id}); err != nil {
			return
		}
	}
	return
}

//获取 研发商 全部 只要 id company_id
func GetAllDevelopCompanies()(companies []DevelopCompany, err error)  {
	o := orm.NewOrm()
	if _, err = o.QueryTable(new(DevelopCompany)).All(&companies,"Id","CompanyId"); err != nil{
		return
	}
	return companies, nil

}
//获取 公司地址
func GetDevAddressByCompanyId(id int)(address string, err error)  {
	developCompany := DevelopCompany{}
	o := orm.NewOrm()
	err = o.QueryTable(new(DevelopCompany)).Filter("CompanyId__exact",id).One(&developCompany,"Address");
	return developCompany.Address, err

}
