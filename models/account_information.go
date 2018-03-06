package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type AccountInformation struct {
	Id             int    `json:"id" orm:"column(id);auto"`
	BodyMy         int    `json:"body_my" orm:"column(body_my);null"`
	AccountName    string `json:"account_name" orm:"column(account_name);size(255);null"`
	Taxpayer       string `json:"taxpayer" orm:"column(taxpayer);size(255);null"`
	BillingAddress string `json:"billing_address" orm:"column(billing_address);size(255);null"`
	Telephone      string `json:"telephone" orm:"column(telephone);size(255);null"`
	Bank           string `json:"bank" orm:"column(bank);size(255);null"`
	AccountNumber  string `json:"account_number" orm:"column(account_number);size(255);null"`
	MailingAddress string `json:"mailing_address" orm:"column(mailing_address);size(255);null"`
}

func (t *AccountInformation) TableName() string {
	return "account_information"
}

func init() {
	orm.RegisterModel(new(AccountInformation))
}

// AddAccountInformation insert a new AccountInformation into database and returns
// last inserted Id on success.
func AddAccountInformation(m *AccountInformation) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAccountInformationById retrieves AccountInformation by Id. Returns error if
// Id doesn't exist
func GetAccountInformationById(id int) (v *AccountInformation, err error) {
	o := orm.NewOrm()
	v = &AccountInformation{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//根据我方主体获取该主体的账户信息
func GetAccountInformation(bodyMy int) (info AccountInformation, err error) {
	o := orm.NewOrm()

	err = o.QueryTable(new(AccountInformation)).
		Filter("body_my__exact", bodyMy).
		One(&info)

	return
}

// UpdateAccountInformation updates AccountInformation by Id and returns error if
// the record to be updated doesn't exist
func UpdateAccountInformationById(m *AccountInformation) (err error) {
	o := orm.NewOrm()
	v := AccountInformation{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}
