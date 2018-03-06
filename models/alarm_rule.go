package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/utils"
)

type AlarmRule struct {
	Id                 int     `json:"id,omitempty" orm:"column(id);auto" json:"id"`
	Type               string  `json:"type,omitempty" orm:"column(type);size(100);null"`
	IntervalTime       int     `json:"interval_time,omitempty" orm:"column(interval_time);null"`
	Amount             float64 `json:"amount,omitempty" orm:"column(amount);null"`
	CustomIntervalTime int     `json:"custom_interval_time,omitempty" orm:"column(custom_interval_time);null"`
	CustomAmount       float64 `json:"custom_amount,omitempty" orm:"column(custom_amount);null"`
	Readonly           int     `json:"-" orm:"column(readonly);null"`
	GameIds            string     `json:"game_ids,omitempty" orm:"column(game_ids);null"`
	ChannelCodes       string     `json:"channel_codes,omitempty" orm:"column(channel_codes);null"`
	RemitCompanyId     int     `json:"remit_company_id,omitempty" orm:"column(remit_company_id);null"`
}

func (t *AlarmRule) TableName() string {
	return "alarm_rule"
}

// AlarmRule.Type 枚举
const (
	AT_ContractTimeout = "contract_timeout" // 合同过期
	AT_ContractSign    = "contract_sign"    // 未签订合同
	AT_VerifyAccount   = "verify_account"   // 对账
	AT_RemitAccount    = "remit_account"    // 回款
	AT_Order           = "order"            // 流水规则
)

func init() {
	orm.RegisterModel(new(AlarmRule))
}

// AddAlarmRule insert a new AlarmRule into database and returns
// last inserted Id on success.
func AddAlarmRule(m *AlarmRule) (id int64, err error) {
	// todo 应该分几种类型分别add
	if !utils.ItemInArray(m.Type, []string{AT_RemitAccount, AT_VerifyAccount}) {
		err = errors.New("type must in Contract_timeout,Contract_sign,verfiy_account,Remit_account")
		return
	}
	m.Readonly = 2
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAlarmRuleById retrieves AlarmRule by Id. Returns error if
// Id doesn't exist
func GetAlarmRuleById(id int) (v *AlarmRule, err error) {
	o := orm.NewOrm()
	v = &AlarmRule{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

type AlarmRuleWithGroup struct {
	ContractTimeout *AlarmRule `json:"contract_timeout,omitempty"`
	ContractSign    *AlarmRule `json:"contract_sign,omitempty"`
	VerifyAccount   *AlarmRule `json:"verify_account,omitempty"`
	RemitAccount    *AlarmRule `json:"remit_account,omitempty"`
	Order           *AlarmRule `json:"order,omitempty"`
}

func GetAllAlarmRuleWithType() (r *AlarmRuleWithGroup, err error) {
	o := orm.NewOrm()
	as := []AlarmRule{}
	if _, err = o.QueryTable(new(AlarmRule)).All(&as); err != nil {
		return
	}
	r = &AlarmRuleWithGroup{}
	for i := range as {
		a := as[i]
		switch a.Type {
		case AT_ContractTimeout:
			r.ContractTimeout = &a
		case AT_ContractSign:
			r.ContractSign = &a
		case AT_VerifyAccount:
			r.VerifyAccount = &a
		case AT_RemitAccount:
			r.RemitAccount = &a
		case AT_Order:
			r.Order = &a
		}
	}

	return
}

// GetAlarmRuleById retrieves AlarmRule by Id. Returns error if
// Id doesn't exist
func GetAlarmRuleByType(typ string) (as []AlarmRule, err error) {
	o := orm.NewOrm()
	v := &AlarmRule{}
	as = []AlarmRule{}
	_, err = o.QueryTable(v).Filter("Type", typ).All(&as)
	return
}

// UpdateAlarmRule updates AlarmRule by Id and returns error if
// the record to be updated doesn't exist
func UpdateAlarmRuleById(m *AlarmRule) (err error) {
	o := orm.NewOrm()
	v := AlarmRule{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err != nil {
		return
	}
	var f []string
	if v.Readonly == 1 {
		// 系统自带的只能更新自定义字段
		f = utils.GetNotEmptyFields(m, "CustomIntervalTime", "CustomAmount")
	} else {
		// 可更新所有
		f = utils.GetNotEmptyFields(m)
	}
	if _, err = o.Update(m, f...); err != nil {
		return
	}
	return
}

// DeleteAlarmRule deletes AlarmRule by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAlarmRule(id int) (err error) {
	o := orm.NewOrm()
	v := AlarmRule{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err != nil {
		return
	}
	if v.Readonly == 1 {
		err = errors.New("this rule is readonly")
		return
	}

	if _, err = o.Delete(&AlarmRule{Id: id}); err == nil {
		return
	}

	return
}
