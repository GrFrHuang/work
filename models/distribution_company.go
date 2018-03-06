package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/utils"
	"strings"
	"strconv"
)

type DistributionCompany struct {
	Id                        int    `json:"id,omitempty" orm:"column(id);auto"`
	CompanyId                 int    `json:"company_id,omitempty" orm:"column(company_id);null"`
	PlatformName              string `json:"platform_name,omitempty" orm:"column(platform_name);size(255);null"` // 平台名称
	Region                    string `json:"region,omitempty" orm:"column(region);size(255);null"`
	Website                   string `json:"website,omitempty" orm:"column(website);size(255);null"`
	Address                   string `json:"address,omitempty" orm:"column(address);size(255);null"` // 公司地址, e.g. ["address1", "address2"]
	AccountName               string `json:"account_name,omitempty" orm:"column(account_name);size(255);null"`
	Bank                      string `json:"bank,omitempty,omitempty" orm:"olumn(bank);size(255);null"`
	AccountNumber             string `json:"account_number,omitempty" orm:"olumn(account_number);size(255);null"`
	ContactIds                string `json:"contact_ids,omitempty" orm:"column(contact_ids);size(255);null"` // 联系人, e.g. [id1, id2, ...]
	Desc                      string `json:"desc,omitempty" orm:"column(desc)"`                              // 备注
	CreateTime                int64  `json:"create_time,omitempty" orm:"column(create_time);null" `
	UpdateTime                int64  `json:"update_time,omitempty" orm:"column(update_time);null" `
	UpdateUserID              int    `json:"update_user_id,omitempty" orm:"column(update_user_id);"`                           // 更新人ID
	YunduanResponsiblePerson  int    `json:"yunduan_responsible_person,omitempty" orm:"column(yunduan_responsible_person);"`   // 云端负责人id
	YouliangResponsiblePerson int    `json:"youliang_responsible_person,omitempty" orm:"column(youliang_responsible_person);"` // 更新人ID
	YunduanResPerName         string `orm:"-" json:"yunduan_res_per_name,omitempty"`
	YouliangResPerName        string `orm:"-" json:"youliang_res_per_name,omitempty"`

	CompanyName        string    `orm:"-" json:"company_name,omitempty"`
	Contacts           []Contact `orm:"-" json:"contacts,omitempty"`
	ContactsForCompare string    `orm:"-" json:"contact_for_compare,omitempty"` //仅用于操作日志比较
	UpdateUserName     string    `orm:"-" json:"update_user_name,omitempty"`
	GameIds            string    `orm:"-" json:"gameids,omitempty"` //该发行商所有发行的游戏（已接入并上线）
	MainContractId     int       `orm:"-" json:"main_contract_id"`
	MainContractState  string    `orm:"-" json:"main_contract_state"`
}

func (t *DistributionCompany) TableName() string {
	return "distribution_company"
}

func init() {
	orm.RegisterModel(new(DistributionCompany))
}

// AddDistributionCompany insert a new DistributionCompany into database and returns
// last inserted Id on success.
func AddDistributionCompany(m *DistributionCompany) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	if err != nil {
		return
	}

	for _, v := range m.Contacts {
		v.CompanyId = int(id)
		v.CompanyType = 2
		AddContact(&v)
	}
	err = CompareAndAddOperateLog(nil, m, m.UpdateUserID, bean.OPP_DISTRIBUTION_COMPANY, int(id), bean.OPA_INSERT)
	return
}

// GetDistributionCompanyById retrieves DistributionCompany by Id. Returns error if
// Id doesn't exist
func GetDistributionCompanyById(id int) (v *DistributionCompany, err error) {
	o := orm.NewOrm()
	v = &DistributionCompany{Id: id}
	if err = o.Read(v); err != nil {
		return
	}
	GetDistributionCompanyAdditionalInfo(v)
	return
}

// GetDistributionCompanyById retrieves DistributionCompany by Id. Returns error if
// Id doesn't exist
func GetDistributionCompanyByCompanyId(id int) (v *DistributionCompany, err error) {
	o := orm.NewOrm()
	v = &DistributionCompany{CompanyId: id}
	if err = o.Read(v, "CompanyId"); err != nil {
		return
	}
	GetDistributionCompanyAdditionalInfo(v)
	return
}

func GetAllDistributionCompaniesWithTotal(offset int, limit int, company_id string, person int) (m []*DistributionCompany,
	total int64, err error) {
	sql := "select * from distribution_company "
	if company_id != "" {
		sql = sql + fmt.Sprintf("where company_id in (%v)", company_id)
	}

	if person != 0 {
		if strings.Contains(sql, "where") {
			sql = sql + fmt.Sprintf(" and (yunduan_responsible_person = %v or youliang_responsible_person = %v)",
				person, person)
		} else {
			sql = sql + fmt.Sprintf(" where (yunduan_responsible_person = %v or youliang_responsible_person = %v)",
				person, person)
		}
	}

	total, err = orm.NewOrm().Raw(sql).QueryRows(&m)
	sql = sql + " order by create_time desc limit " + strconv.Itoa(limit) + " offset " + strconv.Itoa(offset)
	_, err = orm.NewOrm().Raw(sql).QueryRows(&m)
	return m, total, err
}

func GetDistributionCompanyAdditionalInfo(v *DistributionCompany) {
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
	v.Contacts = GetContactsByCompanyId(v.Id, 2)
	// Add the response person name
	userYunduan, err := GetUserById(v.YunduanResponsiblePerson)
	if err == nil {
		v.YunduanResPerName = userYunduan.Nickname
	}
	userYouliang, err := GetUserById(v.YouliangResponsiblePerson)
	if err == nil {
		v.YouliangResPerName = userYouliang.Nickname
	}

	mainContract := GetMainContractByCompanyId(v.CompanyId, 1)
	v.MainContractId = mainContract.Id
	v.MainContractState = GetTypesNameById(mainContract.State)
	return
}

// UpdateDistributionCompany updates DistributionCompany by Id and returns error if
// the record to be updated doesn't exist
func UpdateDistributionCompanyById(m *DistributionCompany, where map[string][]interface{}) (err error) {
	v := DistributionCompany{Id: m.Id}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}

	oldContacts, _ := json.Marshal(GetContactsByCompanyId(v.Id, 2))
	v.ContactsForCompare = string(oldContacts)

	// Update the contacts info
	UpdateOrAddCompanyContacts(m.Contacts, m.Id, 2)

	o := orm.NewOrm()
	fields := util.NewFieldsUtil(m).GetNotEmptyFields().Must("Desc", "Website", "PlatformName", "Region", "Address",
		"AccountName", "Bank", "AccountNumber").Exclude("CompanyName", "Contacts", "UpdateUserName", "YunduanResPerName",
		"YouliangResPerName", "MainContractId", "MainContractState").Fields()

	if _, err = o.Update(m, fields...); err != nil {
		return
	}

	newContacts, _ := json.Marshal(m.Contacts)
	m.ContactsForCompare = string(newContacts)

	fields = append(fields, "ContactsForCompare")
	fields = utils.RemoveFields(fields, "UpdateUserID", "UpdateTime")
	if err = CompareAndAddOperateLog(&v, m, m.UpdateUserID, bean.OPP_DISTRIBUTION_COMPANY, v.Id, bean.OPA_UPDATE, fields...); err != nil {
		return
	}

	return
}

func AddDistributionCompanyInfo(ss *[]DistributionCompany) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.CompanyId
	}
	var games []CompanyType
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

//添加该公司所有发行的游戏
func AddGameIdInfo(ss *[]DistributionCompany) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()

	for i, s := range *ss {
		var gameids []orm.Params
		var ids []string
		o.Raw("SELECT b.game_id FROM distribution_company a LEFT JOIN game b ON a.company_id=b.issue "+
			" WHERE b.game_id > 0 AND a.company_id = ?", s.CompanyId).Values(&gameids)
		for _, id := range gameids {
			ids = append(ids, id["game_id"].(string))
		}
		(*ss)[i].GameIds = strings.Join(ids, ",")

	}

	return
}

// DeleteDistributionCompany deletes DistributionCompany by Id and returns error if
// the record to be deleted doesn't exist
func DeleteDistributionCompany(id int) (err error) {
	o := orm.NewOrm()
	v := DistributionCompany{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		DeleteContactsByCompanyId(v.Id)
		var num int64
		if num, err = o.Delete(&DistributionCompany{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// 获取所有发行商的公司信息
func GetAllDistributionCompanies() (companies []CompanyType, err error) {
	var res []orm.Params
	o := orm.NewOrm()
	_, err = o.Raw("SELECT " +
		"company_id, company_type.`name` " +
		"FROM `distribution_company` " +
		"LEFT JOIN company_type ON distribution_company.company_id = company_type.id").
		Values(&res)
	if err != nil {
		return
	}
	companies = []CompanyType{}
	for _, v := range res {
		id, _ := util.Interface2Int(v["company_id"], false)
		name, _ := util.Interface2String(v["name"], false)
		companies = append(companies, CompanyType{
			Name: name,
			Id:   int(id),
		})
	}

	return
}

func GetDistributionCompanyRegion() (region []string, err error) {
	var res []orm.Params
	o := orm.NewOrm()
	_, err = o.Raw("SELECT DISTINCT`region` AS region FROM `distribution_company`").Values(&res)
	if err != nil {
		return
	}
	region = []string{}
	for _, v := range res {
		r, _ := util.Interface2String(v["region"], false)
		if r != "" {
			region = append(region, r)
		}
	}

	return
}

//获取 公司地址
func GetDistAddressByCompanyId(id int) (address string, err error) {
	developCompany := DistributionCompany{}
	o := orm.NewOrm()
	err = o.QueryTable(new(DistributionCompany)).Filter("CompanyId__exact", id).One(&developCompany, "Address")
	return developCompany.Address, err

}
