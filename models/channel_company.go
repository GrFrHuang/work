package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/utils"
	"strconv"
	"strings"
)

type ChannelCompany struct {
	Id                        int    `json:"id,omitempty" orm:"column(id);auto"`
	ChannelCode               string `json:"channel_code,omitempty" orm:"column(channel_code)"`
	Properties                string `json:"properties,omitempty" orm:"column(properties);size(255);null"`
	Website                   string `json:"website,omitempty" orm:"column(website);size(255);null"`
	Address                   string `json:"address,omitempty" orm:"column(address);size(255);null"`
	Region                    string `json:"region,omitempty" orm:"column(region);size(255);null"`
	PlatformName              string `json:"platform_name,omitempty" orm:"column(platform_name);size(255);null"` // 平台名称
	CompanyId                 int    `json:"company_id,omitempty" orm:"column(company_id); null"`
	YunduanResponsiblePerson  int    `json:"yunduan_responsible_person,omitempty" orm:"column(yunduan_responsible_person);null" `   // 云端负责人ID
	YouliangResponsiblePerson int    `json:"youliang_responsible_person,omitempty" orm:"column(youliang_responsible_person);null" ` // 有量负责人ID
	CooperateState            int    `json:"cooperate_state,omitempty" orm:"column(cooperate_state);null" `                         // 合作状态
	ContactIds                string `json:"contact_ids,omitempty" orm:"column(contact_ids);size(255);null"`                        // 联系人, e.g. [id1, id2, ...]
	Desc                      string `json:"desc,omitempty" orm:"column(desc)"`
	CreateTime                int64  `json:"create_time,omitempty" orm:"column(create_time);null" `
	UpdateTime                int64  `json:"update_time,omitempty" orm:"column(update_time);null" `
	UpdateUserID              int    `json:"update_user_id,omitempty" orm:"column(update_user_id);"` // 更新人ID

	RemitCompany       string         `orm:"-" json:"remit_company,omitempty"` // 回款主体, e.g. [id1, id2, ...]
	CompanyName        string         `orm:"-" json:"company_name,omitempty"`
	ChannelName        string         `orm:"-" json:"channel_name,omitempty"`
	Contacts           []Contact      `orm:"-" json:"contacts,omitempty"`
	ContactsForCompare string         `orm:"-" json:"contact_for_compare,omitempty"` //仅用于操作日志比较
	Remits             []*CompanyType `orm:"-" json:"remits,omitempty"`
	UpdateUserName     string         `orm:"-" json:"update_user_name,omitempty"`
	YunduanResPerName  string         `orm:"-" json:"yunduan_res_per_name,omitempty"`
	YouliangResPerName string         `orm:"-" json:"youliang_res_per_name,omitempty"`
	MainContractId     int            `orm:"-" json:"main_contract_id"`
	MainContractState  string         `orm:"-" json:"main_contract_state"`
}

func (t *ChannelCompany) TableName() string {
	return "channel_company"
}

func init() {
	orm.RegisterModel(new(ChannelCompany))
}

// AddChannelCompany insert a new ChannelCompany into database and returns
// last inserted Id on success.
func AddChannelCompany(m *ChannelCompany) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	if err != nil {
		return
	}

	// Add remit companies
	if m.RemitCompany != "" && m.RemitCompany != "[]" {
		var remits []int
		err = json.Unmarshal([]byte(m.RemitCompany), &remits)
		if err != nil {
			return
		}
		for _, v := range remits {
			var remit RemitCompany
			remit.ChannelCode = m.ChannelCode
			remit.CompanyId = v
			_, err = o.Insert(&remit)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	// Add contacts
	for _, v := range m.Contacts {
		v.CompanyId = int(id)
		v.CompanyType = 3
		AddContact(&v)
	}

	fields := util.NewFieldsUtil(m).GetNotEmptyFields().Fields()
	newContacts, _ := json.Marshal(m.Contacts)
	m.ContactsForCompare = string(newContacts)
	fields = append(fields, "RemitCompany", "ContactsForCompare")

	err = CompareAndAddOperateLog(nil, m, m.UpdateUserID, bean.OPP_CHANNEL_COMPANY, int(id), bean.OPA_INSERT, fields...)
	return
}

// GetChannelCompanyById retrieves ChannelCompany by Id. Returns error if
// Id doesn't exist
func GetChannelCompanyById(id int) (v *ChannelCompany, err error) {
	o := orm.NewOrm()
	v = &ChannelCompany{Id: id}
	if err = o.Read(v); err != nil {
		return
	}
	GetChannelCompanyAdditionalInfo(v)
	return
}

// AddChannelCompanyAdditionalInfo retrieves the additional info of ChannelCompany. Returns
// the ChannelCompany
func GetChannelCompanyAdditionalInfo(v *ChannelCompany) {
	// Add the company name
	c, err := GetCompanyTypeById(v.CompanyId)
	if err == nil {
		v.CompanyName = c.Name
	}

	// Add the channel name
	channel_name, err := GetChannelNameByCp(v.ChannelCode)
	if err == nil {
		v.ChannelName = channel_name
	}

	// Add the update user name
	user, err := GetUserById(v.UpdateUserID)
	if err == nil {
		v.UpdateUserName = user.Nickname
	}

	// Add the response person name
	userYunduan, err := GetUserById(v.YunduanResponsiblePerson)
	if err == nil {
		v.YunduanResPerName = userYunduan.Nickname
	}
	userYouliang, err := GetUserById(v.YouliangResponsiblePerson)
	if err == nil {
		v.YouliangResPerName = userYouliang.Nickname
	}

	// Add the remit company
	//var remits []int
	//var remit_companies []*CompanyType
	//err = json.Unmarshal([]byte(v.RemitCompany), &remits)
	//
	//for _, i := range remits {
	//	company, err := GetCompanyTypeById(i)
	//	if err == nil {
	//		remit_companies = append(remit_companies, company)
	//	} else {
	//		fmt.Println(err)
	//	}
	//}
	//v.Remits = remit_companies

	var remits []*RemitCompany
	var remit_companies []*CompanyType
	var rs []int
	o := orm.NewOrm()
	_, err = o.QueryTable("remit_company").Filter("channel_code", v.ChannelCode).All(&remits)
	if err == nil {
		for _, v := range remits {
			company, err := GetCompanyTypeById(v.CompanyId)
			if err == nil {
				remit_companies = append(remit_companies, company)
			}
			rs = append(rs, v.CompanyId)
		}
	}
	v.Remits = remit_companies
	rsString, _ := json.Marshal(rs)
	v.RemitCompany = string(rsString)

	// Add the contacts info
	v.Contacts = GetContactsByCompanyId(v.CompanyId, 3)

	mainContract := GetMainContractByCompanyId(v.CompanyId, 2)
	v.MainContractId = mainContract.Id
	v.MainContractState = GetTypesNameById(mainContract.State)

	return
}

func AddCompanyInfo(ss *[]ChannelCompany) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.CompanyId
	}
	var companys []CompanyType
	_, err := o.QueryTable(new(CompanyType)).
		Filter("Id__in", linkIds).
		All(&companys, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]CompanyType{}
	for _, g := range companys {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.CompanyId]; ok {
			(*ss)[i].CompanyName = g.Name
		}
	}

	// 由于某些公司主体有多个渠道，所以在返回公司名的时候加以区分
	var channelCompanyName []orm.Params
	sql := "SELECT a.company_id,a.channel_code,b.name FROM channel_company a LEFT JOIN channel b ON a.channel_code=b.cp " +
		"WHERE a.company_id IN(SELECT company_id FROM channel_company GROUP BY company_id HAVING COUNT(id)>1) ORDER BY company_id"
	total, err := o.Raw(sql).Values(&channelCompanyName)
	if total == 0 || err != nil {
		return
	}
	channelMap := map[string]string{}
	for _, name := range channelCompanyName {
		company_id, _ := utils.Interface2String(name["company_id"], false)
		channel_code, _ := utils.Interface2String(name["channel_code"], false)
		channel_name, _ := utils.Interface2String(name["name"], false)

		channelMap[company_id+"_"+channel_code] = channel_name
	}
	for i, s := range *ss {
		if g, ok := channelMap[strconv.Itoa(s.CompanyId)+"_"+s.ChannelCode]; ok {
			(*ss)[i].CompanyName = (*ss)[i].CompanyName + "(" + g + ")"
		}
	}

	return
}

func GetChannelCompanyByCode(code string) (v *ChannelCompany, err error) {
	o := orm.NewOrm()
	v = &ChannelCompany{ChannelCode: code}
	if err = o.Read(v, "channel_code"); err != nil {
		return
	}
	GetChannelCompanyAdditionalInfo(v)
	return
}

// UpdateChannelCompany updates ChannelCompany by Id and returns error if
// the record to be updated doesn't exist
func UpdateChannelCompanyById(m *ChannelCompany, where map[string][]interface{}) (err error) {
	o := orm.NewOrm()
	v := ChannelCompany{Id: m.Id}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}

	// Update remit company
	var remitsOld []*RemitCompany
	o.QueryTable("remit_company").Filter("channel_code", m.ChannelCode).All(&remitsOld)
	codeOld := map[int]bool{}
	for _, v := range remitsOld {
		codeOld[v.CompanyId] = true
	}
	var remitsNew []int
	if err = json.Unmarshal([]byte(m.RemitCompany), &remitsNew); err != nil {
		return
	}
	codeNew := map[int]bool{}
	// 增加数据库中不存在的回款主体
	for _, v := range remitsNew {
		if codeOld[v] == false {
			var newRemit RemitCompany
			newRemit.ChannelCode = m.ChannelCode
			newRemit.CompanyId = v
			_, err = o.Insert(&newRemit)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		codeNew[v] = true
	}
	// 删除数据库中多余的回款主体
	for _, v := range remitsOld {
		if codeNew[v.CompanyId] == false {
			o.Delete(v)
		}
	}

	oldContacts, _ := json.Marshal(GetContactsByCompanyId(v.Id, 3))
	v.ContactsForCompare = string(oldContacts)

	// Update contract
	UpdateOrAddCompanyContacts(m.Contacts, m.CompanyId, 3)

	// fields := utils.GetNotEmptyFields(m, "ChannelCode", "RemitCompanyIds", "Properties", "UpdateUserID", "UpdateTime", "YunduanResponsiblePerson", "YouliangResponsiblePerson")
	fields := util.NewFieldsUtil(m).GetNotEmptyFields().
		Must("Desc", "Website", "Properties", "Address", "PlatformName").
		Exclude("RemitCompany", "GameId", "ChannelCode", "CompanyName", "ChannelName", "UpdateUserName", "Contacts",
		"Remits", "YunduanResPerName", "YouliangResPerName", "MainContractId", "MainContractState").Fields()
	if _, err = o.Update(m, fields...); err != nil {
		return
	}

	newContacts, _ := json.Marshal(m.Contacts)
	m.ContactsForCompare = string(newContacts)

	fields = append(fields, "ContactsForCompare")
	fields = utils.RemoveFields(fields, "UpdateUserID", "UpdateTime")
	if err = CompareAndAddOperateLog(&v, m, m.UpdateUserID, bean.OPP_CHANNEL_COMPANY, v.Id, bean.OPA_UPDATE, fields...); err != nil {
		return
	}
	return
}

//通过 channelCodes 获取合同的 未签订 已到期 的条件 获得游戏 game_id
func GetGameIdByChannelCodesAndState(channel_code []string) (contracts []Contract, err error) {
	o := orm.NewOrm()
	links := [2]int{149, 151}
	_, err = o.QueryTable(new(Contract)).Filter("ChannelCode__in", channel_code).Filter("State__in", links).
		Filter("company_type__exact", 1).All(&contracts)

	//_, err = o.Raw("SELECT id,game_id FROM contract c WHERE c.channel_code in ? " +
	//	"AND c.state IN (149,150)",channel_code).QueryRows(&contracts)
	//if err != nil{
	//	return
	//}
	//return contracts, nil
	return

}

//通过公司company_id 获取渠道codes (channel_code)
func GetChannelCodesByCompanyId(id int) (channelCodes []string, err error) {
	o := orm.NewOrm()
	//annel := ChannelCompany{}
	var channel []*ChannelCompany
	qs := o.QueryTable(new(ChannelCompany))
	_, err = qs.Filter("CompanyId__exact", id).All(&channel)
	//if  err = o.Read(&channel,"Id","ChannelCode"); err != nil{
	//	fmt.Printf("err:%v-----\n", err)
	//	return "", err
	//}
	//fmt.Printf("ChannelName:%v-----\n", channel.ChannelName)
	channelCodes = make([]string, len(channel))
	for _, s := range channel {
		channelCodes = append(channelCodes, s.ChannelCode)
	}
	return channelCodes, err
}

// DeleteChannelCompany deletes ChannelCompany by Id and returns error if
// the record to be deleted doesn't exist
func DeleteChannelCompany(id int) (err error) {
	o := orm.NewOrm()
	v := ChannelCompany{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		DeleteContactsByCompanyId(v.Id)
		var num int64
		if num, err = o.Delete(&ChannelCompany{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//通过公司company_id 获取渠道code (channel_code)
func GetChannelCodeByCompanyId(id int) (channelName string, err error) {
	o := orm.NewOrm()
	channel := ChannelCompany{}
	qs := o.QueryTable(new(ChannelCompany))
	err = qs.Filter("CompanyId__exact", id).One(&channel)
	//if  err = o.Read(&channel,"Id","ChannelCode"); err != nil{
	//	fmt.Printf("err:%v-----\n", err)
	//	return "", err
	//}
	//fmt.Printf("ChannelName:%v-----\n", channel.ChannelName)
	return channel.ChannelCode, err
}

// 获取所有渠道商的公司信息
func GetAllChannelCompanies() (companies []CompanyType, err error) {
	var res []orm.Params
	o := orm.NewOrm()
	_, err = o.Raw("SELECT " +
		"id, company_type.`name`" +
		"FROM `channel_company` " +
		"LEFT JOIN company_type ON channel_company.company_id = company_type.id").
		Values(&res)
	if err != nil {
		return
	}
	companies = []CompanyType{}
	for _, v := range res {
		id, _ := util.Interface2Int(v["id"], false)
		name, _ := util.Interface2String(v["name"], false)
		companies = append(companies, CompanyType{
			Name: name,
			Id:   int(id),
		})
	}

	return
}

func GetAllChannelCompaniesWithTotal(offset int, company_id string, channel_code string,
	cooperate_state int, person int) (m []*ChannelCompany, total int64, err error) {

	sql := "select * from channel_company "
	if company_id != "" {
		sql = sql + fmt.Sprintf("where company_id in (%v)", company_id)
	}
	if channel_code != "" {
		if strings.Contains(sql, "where") {
			sql = sql + fmt.Sprintf(" and channel_code = '%v'", channel_code)
		} else {
			sql = sql + fmt.Sprintf(" where channel_code = '%v'", channel_code)
		}
	}
	if cooperate_state != 0 {
		if strings.Contains(sql, "where") {
			sql = sql + fmt.Sprintf(" and cooperate_state = %v", cooperate_state)
		} else {
			sql = sql + fmt.Sprintf(" where cooperate_state = %v", cooperate_state)
		}
	}
	if person != 0 {
		if strings.Contains(sql, "where") {
			sql = sql + fmt.Sprintf(" and (yunduan_responsible_person = %v or youliang_responsible_person = %v)", person, person)
		} else {
			sql = sql + fmt.Sprintf(" where (yunduan_responsible_person = %v or youliang_responsible_person = %v)", person, person)
		}
	}

	total, err = orm.NewOrm().Raw(sql).QueryRows(&m)
	sql = sql + " order by create_time desc limit 15 " + " offset " + strconv.Itoa(offset)
	_, err = orm.NewOrm().Raw(sql).QueryRows(&m)
	return m, total, err
}

//
func GetAllChannelConmpaniesNew() (companies []CompanyType, err error) {
	o := orm.NewOrm()
	_, err = o.Raw("SELECT channel_company.`company_id` id, company_type.`name` FROM channel_company LEFT JOIN company_type On " +
		"channel_company.company_id = company_type.id").QueryRows(&companies)
	return

}

//获取公司地址
func GetChannelAddressByCompanyId(id int) (address string, err error) {
	developCompany := ChannelCompany{}
	o := orm.NewOrm()
	err = o.QueryTable(new(ChannelCompany)).Filter("CompanyId__exact", id).One(&developCompany, "Address")
	return developCompany.Address, err

}

//获取公司的channel_code
func GetChannelCompanyCodeById(id int) (channel_code string, err error) {
	channelCompany := ChannelCompany{}
	o := orm.NewOrm()
	err = o.QueryTable(new(ChannelCompany)).Filter("CompanyId__exact", id).One(&channelCompany, "Id", "ChannelCode")
	return channelCompany.ChannelCode, err
}
