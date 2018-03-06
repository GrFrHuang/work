package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"sync"
	"time"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/utils"
)

// 已经回款的账单, 后台添加
type RemitDownAccount struct {
	Id             int     `json:"id,omitempty" orm:"column(id);auto"`
	UserId         int     `json:"user_id,omitempty" orm:"column(user_id);null"`
	RemitCompanyId int     `json:"remit_company_id,omitempty" orm:"column(remit_company_id);null" valid:"Required"`
	RemitTime      int     `json:"remit_time,omitempty" orm:"column(remit_time);null" valid:"Required"` // 此账单的开始日期
	Amount         float64 `json:"amount,omitempty" orm:"column(amount);null;digits(16);decimals(0)" valid:"Required"`
	Extend         string  `json:"extend,omitempty" orm:"column(extend);null"`
	FileId         int     `json:"file_id,omitempty" orm:"column(file_id);null"`
	FilePreviewId  int     `json:"file_preview_id,omitempty" orm:"column(file_preview_id);null"`
	CreatedTime    int     `json:"created_time,omitempty" orm:"column(created_time);null;"`
	BodyMy         int     `json:"body_my,omitempty" orm:"column(body_my);null;"`

	UpdateRemitTime     int64 `json:"update_remittime,omitempty" orm:"column(update_remittime);null" `
	UpdateChannelUserID int   `json:"update_remituserid,omitempty" orm:"column(update_remituserid)"`

	Details      []RemitDownDetail `orm:"-" json:"details"`
	RemitCompany *CompanyType      `orm:"-" json:"remit_company,omitempty"`
	User         *User             `orm:"-" json:"user,omitempty"`
	UpdateUser   *User             `orm:"-" json:"update_user,omitempty"`
}

func (t *RemitDownAccount) TableName() string {
	return "remit_down_account"
}

func init() {
	orm.RegisterModel(new(RemitDownAccount))
}

var lockRemit4 sync.Mutex
// 添加
func AddRemitDownAccounts(m *RemitDownAccount) (id int64, err error) {
	lockRemit4.Lock()
	defer lockRemit4.Unlock()
	// 检查重复
	o := orm.NewOrm()
	if o.QueryTable(m).
		Filter("RemitCompanyId", m.RemitCompanyId).
		Filter("Amount", m.Amount).
		Filter("RemitTime", m.RemitTime).
		Exist() {
		err = errors.New("plase don't resubmit")
		return
	}

	m.Id = 0
	m.CreatedTime = int(time.Now().Unix())
	id, err = o.Insert(m)
	if err != nil {
		return
	}

	err = OnRemitDownAccountChanged(RemitDownAccount{}, *m)

	err = CompareAndAddOperateLog(nil, m, m.UpdateChannelUserID, bean.OPP_CHANNEL_REMIT, int(id), bean.OPA_INSERT)

	return
}

func GetNotRemitAccount(company_id int) (verifys []VerifyChannel, err error) {
	o := orm.NewOrm()
	sql := "SELECT a.* FROM verify_channel a LEFT JOIN remit_down_detail b ON a.id=b.verify_channel_id AND " +
		" a.date=b.remit_month WHERE a.amount_payable!=b.remit_money OR b.remit_money IS NULL AND a.remit_company_id=?"

	_, err = o.Raw(sql, company_id).QueryRows(&verifys)
	if err != nil {
		return nil, err
	}

	return verifys, nil
}

// get
func GetRemitDownAccountsById(id int) (v *RemitDownAccount, err error) {
	o := orm.NewOrm()
	v = &RemitDownAccount{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// 更新
func UpdateRemitDownAccountsById(m *RemitDownAccount, where map[string][]interface{}) (err error) {
	o := orm.NewOrm()
	v := RemitDownAccount{}
	// ascertain id exists in the database
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}

	fs := util.NewFieldsUtil(m).
		GetNotEmptyFields().
		Filter("BodyMy", "Amount", "Extend", "UserId").
		Must("Amount", "FileId", "RemitTime", "RemitCompanyId", "UpdateChannelUserID", "UpdateRemitTime").Fields()

	old := RemitDownAccount{Id: m.Id}
	if err = orm.NewOrm().Read(&old); err != nil {
		return
	}

	if _, err = o.Update(m, fs...); err != nil {
		return
	}

	err = OnRemitDownAccountChanged(v, *m)
	fs = utils.RemoveFields(fs, "UpdateChannelUserID", "UpdateRemitTime")

	if err = CompareAndAddOperateLog(&old, m, m.UpdateChannelUserID, bean.OPP_CHANNEL_REMIT, m.Id, bean.OPA_UPDATE, fs...); err != nil {
		return
	}

	return
}

func AddRemitCompanyInfo4Remit(ss *[]RemitDownAccount) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]interface{}, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.RemitCompanyId
	}
	util.UnDuplicatesSlice(&linkIds)

	var link []CompanyType
	_, err := o.QueryTable(new(CompanyType)).
		Filter("Id__in", linkIds).
		All(&link, "Id", "Name")
	if err != nil {
		return
	}
	linkMap := map[int]CompanyType{}
	for _, g := range link {
		linkMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := linkMap[s.RemitCompanyId]; ok {
			(*ss)[i].RemitCompany = &g
		}
	}
	return
}

func AddUserInfo4Remit(ss *[]RemitDownAccount) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]interface{}, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.UserId
	}
	util.UnDuplicatesSlice(&linkIds)

	var link []User
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&link, "Id", "Name", "NickName")
	if err != nil {
		return
	}
	linkMap := map[int]User{}
	for _, g := range link {
		linkMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := linkMap[s.UserId]; ok {
			(*ss)[i].User = &g
		}
	}
	return
}

// 取更新人的信息
func AddUpdateUserInfoRemit(ss *[]RemitDownAccount) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]interface{}, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.UpdateChannelUserID
	}
	util.UnDuplicatesSlice(&linkIds)

	var link []User
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&link, "Id", "NickName")
	if err != nil {
		return
	}
	linkMap := map[int]User{}
	for _, g := range link {
		linkMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := linkMap[s.UpdateChannelUserID]; ok {
			(*ss)[i].UpdateUser = &g
		}
	}
	return
}

// 回款添加对账单信息
func AddRemitDownDetailInfo(ss *[]RemitDownAccount) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	var details []RemitDownDetail
	for i, s := range *ss {
		o.QueryTable(new(RemitDownDetail)).Filter("remit_down_id__exact", s.Id).All(&details)
		(*ss)[i].Details = details
	}
	return
}

// 删除
func DeleteRemitDownAccounts(id int, where map[string][]interface{}) (err error) {
	v := &RemitDownAccount{Id: id}
	qs := QueryTable(v, where).Filter("Id", id)
	if _, err = qs.Delete(); err != nil {
		return
	}

	// bu bu pu pong
	err = OnRemitDownAccountChanged(*v, RemitDownAccount{})
	if err != nil {
		return
	}

	return
}

func OnRemitDownAccountChanged(old, new RemitDownAccount) (err error) {
	needAdd := false

	if old.Id == 0 || new.Id == 0 {
		if old.Id == 0 { // create
			// 4月之前的
			needAdd = new.RemitTime < 1490976000
		} else { // delete
			// 4月之前的
			needAdd = old.RemitTime < 1490976000
		}
	} else {
		needAdd = new.RemitTime < 1490976000
	}

	// change
	if needAdd {
		offSet := new.Amount - old.Amount
		if offSet == 0 {
			return
		}
		err = AddRemitPreAmountOffset(new.BodyMy, new.RemitCompanyId, offSet)
		if err != nil {
			return
		}
	}

	return
}
