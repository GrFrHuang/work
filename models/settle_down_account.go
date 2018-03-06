package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/utils"
	"sync"
	"time"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

// 已经结算的账单, 后台添加
type SettleDownAccount struct {
	Id                 int     `json:"id,omitempty" orm:"column(id);auto"`
	UserId             int     `json:"user_id,omitempty" orm:"column(user_id);null"`
	CompanyId          int     `json:"company_id,omitempty" orm:"column(company_id);null" valid:"Required"`
	Games              string  `json:"games,omitempty" orm:"column(games);null"` // 游戏列表json {GameAmount}
	SettleTime         int     `json:"settle_time,omitempty" orm:"column(settle_time);null" valid:"Required"`
	Amount             float64 `json:"amount,omitempty" orm:"column(amount);null;digits(16);decimals(0)" valid:"Required"`
	Extend             string  `json:"extend,omitempty" orm:"column(extend);null"`
	FileId             int     `json:"file_id,omitempty" orm:"column(file_id);null" `
	FilePreviewId      int     `json:"file_preview_id,omitempty" orm:"column(file_preview_id);null"`
	CreatedTime        int     `json:"created_time,omitempty" orm:"column(created_time);null;"`
	BodyMy             int     `json:"body_my,omitempty" orm:"column(body_my);null;"`
	UpdateSettleTime   int64   `json:"update_settletime,omitempty" orm:"column(update_settletime);null" `
	UpdateSettleUserID int     `json:"update_settleuserid,omitempty" orm:"column(update_settleuserid)"`

	Company    *CompanyType `orm:"-" json:"company,omitempty"`
	User       *User        `orm:"-" json:"user,omitempty"`
	UpdateUser *User        `orm:"-" json:"update_user,omitempty"`
}

type NotSettled struct {
	Companies []CompanyType `json:"companies"`
	Amount    float64       `json:"amount"`
	StartTime int64         `json:"start_time"`
	EndTime   int64         `json:"end_time"`
}

func (t *SettleDownAccount) TableName() string {
	return "settle_down_account"
}

func init() {
	orm.RegisterModel(new(SettleDownAccount))
}

var lock4 sync.Mutex

// AddSettleDownAccounts insert a new SettleDownAccounts into database and returns
// last inserted Id on success.
func AddSettleDownAccounts(m *SettleDownAccount) (id int64, err error) {
	lock4.Lock()
	defer lock4.Unlock()
	// 检查重复
	o := orm.NewOrm()
	if o.QueryTable(m).
		Filter("CompanyId", m.CompanyId).
		Filter("Amount", m.Amount).
		Filter("SettleTime", m.SettleTime).
		Exist() {
		err = errors.New("plase don't resubmit")
		return
	}

	m.Id = 0
	m.CreatedTime = int(time.Now().Unix())

	id, err = o.Insert(m)
	// todo 余额功能暂时不做 2017-03-23 想的是全量统计
	//if err == nil {
	//	// 添加预付款金额
	//	er := AddSettlePreAmount(m.CompanyId, m.Amount)
	//	if er != nil {
	//		err = er
	//		return
	//	}
	//	// 将预付款注入对账单
	//	er = DoPushSettleToVerify(m.CompanyId)
	//	if er != nil {
	//		err = er
	//		return
	//	}
	//}
	err = CompareAndAddOperateLog(nil, m, m.UpdateSettleUserID, bean.OPP_CP_REMIT, int(id), bean.OPA_INSERT)
	return
}

// GetSettleDownAccountsById retrieves SettleDownAccounts by Id. Returns error if
// Id doesn't exist
func GetSettleDownAccountsById(id int) (v *SettleDownAccount, err error) {
	o := orm.NewOrm()
	v = &SettleDownAccount{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateSettleDownAccounts updates SettleDownAccounts by Id and returns error if
// the record to be updated doesn't exist
func UpdateSettleDownAccountsById(m *SettleDownAccount) (err error) {
	o := orm.NewOrm()
	// ascertain id exists in the database

	fs := util.NewFieldsUtil(m).
		GetNotEmptyFields().
		Filter("BodyMy", "CompanyId", "UpdateSettleUserID", "UpdateSettleTime", "UserId").
		Must("Amount", "FileId", "SettleTime", "Extend").Fields()

	s := SettleDownAccount{Id: m.Id}
	if err = orm.NewOrm().Read(&s); err != nil {
		return
	}

	if _, err = o.Update(m, fs...); err != nil {
		return
	}

	fs = utils.RemoveFields(fs, "UpdateSettleUserID", "UpdateSettleTime")
	if err = CompareAndAddOperateLog(&s, m, m.UpdateSettleUserID, bean.OPP_CP_REMIT, s.Id, bean.OPA_UPDATE, fs...); err != nil {
		return
	}
	return
}

// 添加发行商info
func GroupSettleAddCompanyInfo(ss *[]SettleDownAccount) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	companyIds := make([]interface{}, len(*ss))
	for i, s := range *ss {
		companyIds[i] = s.CompanyId
	}
	utils.UnDuplicatesSlice(&companyIds)

	var companies []CompanyType
	_, err := o.QueryTable(new(CompanyType)).Filter("Id__in", companyIds).All(&companies, "Id", "Name")
	if err != nil {
		return
	}
	companyMap := map[int]CompanyType{}
	for _, g := range companies {
		companyMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := companyMap[s.CompanyId]; ok {
			(*ss)[i].Company = &g
		}
	}
	return
}

func AddUserInfo4Settle(ss *[]SettleDownAccount) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]interface{}, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.UserId
	}
	utils.UnDuplicatesSlice(&linkIds)
	var users []User
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&users, "Id", "Name", "NickName")
	if err != nil {
		return
	}
	link := map[int]User{}
	for _, g := range users {
		link[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := link[s.UserId]; ok {
			(*ss)[i].User = &g
		}
	}
	return
}

//添加更新人信息
func AddUpdateUserInfoSettle(ss *[]SettleDownAccount) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]interface{}, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.UpdateSettleUserID
	}
	utils.UnDuplicatesSlice(&linkIds)
	var users []User
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&users, "Id", "NickName")
	if err != nil {
		return
	}
	link := map[int]User{}
	for _, g := range users {
		link[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := link[s.UpdateSettleUserID]; ok {
			(*ss)[i].UpdateUser = &g
		}
	}
	return
}

func DeleteSettleDownAccounts(id int) (err error) {
	o := orm.NewOrm()
	if _, err = o.Delete(&SettleDownAccount{Id: id}); err != nil {
		return
	}
	return
}
