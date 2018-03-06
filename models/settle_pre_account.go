package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"sync"
	"time"
)

// 存放预付款（有结算但没有对账，则不能往对账单里倒水，则把金额累计到这个表里，当有对账单完成的时候再从这里取）
type SettlePreAccount struct {
	Id          int     `json:"id,omitempty" orm:"column(id);auto"`
	CompanyId      int     `json:"company_id,omitempty" orm:"column(company_id);null"`
	UpdatedTime int     `json:"updated_time,omitempty" orm:"column(updated_time);null"`
	Amount      float64 `json:"amount,omitempty" orm:"column(amount);null;digits(16);decimals(0)"`
}

func (t *SettlePreAccount) TableName() string {
	return "settle_pre_account"
}

func init() {
	orm.RegisterModel(new(SettlePreAccount))
}

var lock sync.Mutex

// 当有对账单完成了,或者出来新结账单,就取出余额并注入到未结算的单子
//func DoPushSettleToVerify(companyId int) (err error) {
//	lock.Lock()
//	defer lock.Unlock()
//
//	o := orm.NewOrm()
//	count := 0
//	for {
//		sp := &SettlePreAccount{CompanyId:companyId}
//		// 没找到游戏,说明没有渠道的已结算的钱,直接返回等到结算
//		e := o.Read(sp, "CompanyId")
//		if e != nil {
//			return
//		}
//		if sp.Amount == 0 {
//			return
//		}
//
//		maps := []orm.Params{}
//		_, err = o.Raw("SELECT id , start_time, amount_payable , amount_settle from cp_verify_account WHERE company_id = ? AND amount_payable != amount_settle  ORDER BY start_time limit 0,1", companyId).
//			Values(&maps)
//		if err != nil {
//			return
//		}
//		if len(maps) == 0 {
//			break
//		}
//
//		m := maps[0]
//		id := m["id"].(string)
//		idInt, _ := strconv.Atoi(id)
//		amountPayAble := m["amount_payable"].(string)
//		amountSettle := m["amount_settle"].(string)
//		fP, e := strconv.ParseFloat(amountPayAble, 64)
//		if e != nil {
//			err = e
//			return
//		}
//		fS, e := strconv.ParseFloat(amountSettle, 64)
//		if e != nil {
//			err = e
//			return
//		}
//		not := fP - fS
//		if not >= sp.Amount {
//			// 如果余额装不满对账单，则只装这一个对账单并清空余额
//			x := sp.Amount + fS
//			if not == sp.Amount {
//				x = fP
//			}
//			cp := &CpVerifyAccount{Id:idInt, AmountSettle:x}
//			_, e := o.Update(cp, "AmountSettle")
//			if e != nil {
//				err = e
//				return
//			}
//			sp.Amount = 0
//			_, e = o.Update(sp, "Amount")
//			if e != nil {
//				// todo 这里如果错误的话，应该将上面的语句回滚
//				err = e
//				log.Error("SETTLE", e)
//				return
//			}
//			break
//		} else {
//			// 如果余额大于当前账单,则填充满当前账单并重复多次填充剩余账单
//			cp := &CpVerifyAccount{Id:idInt, AmountSettle:fP, Status:CP_VERIFY_S_SETTLE}
//			_, e := o.Update(cp, "AmountSettle", "Status")
//			if e != nil {
//				err = e
//				return
//			}
//			//hasAmount -= not
//			sp.Amount = sp.Amount - not
//			sp.UpdatedTime = int(time.Now().Unix())
//			_, e = o.Update(sp, "UpdatedTime", "Amount")
//			if e != nil {
//				// todo 这里如果错误的话，应该将上面的语句回滚
//				err = e
//				log.Error("SETTLE", e)
//				return
//			}
//
//		}
//
//		count++
//	}
//
//	return
//}

// 添加发行商的余额
var lock2 sync.Mutex

func AddSettlePreAmount(companyId int, amount float64) (err error) {
	lock2.Lock()
	defer lock2.Unlock()

	o := orm.NewOrm()
	v := &SettlePreAccount{CompanyId:companyId}
	e := o.Read(v, "CompanyId")
	if e != nil {
		if e == orm.ErrNoRows {
			v.Amount = amount
			v.UpdatedTime = int(time.Now().Unix())
			_, e = o.Insert(v)
			if e != nil {
				err = e
			}
			return
		} else {
			err = e
			return
		}
	}
	v.Amount = v.Amount + amount
	v.UpdatedTime = int(time.Now().Unix())
	_, err = o.Update(v, )
	if err != nil {
		return
	}
	return
}

// GetSettlePreAccountById retrieves SettlePreAccount by Id. Returns error if
// Id doesn't exist
func GetSettlePreAccountById(id int) (v *SettlePreAccount, err error) {
	o := orm.NewOrm()
	v = &SettlePreAccount{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateSettlePreAccount updates SettlePreAccount by Id and returns error if
// the record to be updated doesn't exist
func UpdateSettlePreAccountById(m *SettlePreAccount) (err error) {
	o := orm.NewOrm()
	v := SettlePreAccount{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSettlePreAccount deletes SettlePreAccount by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSettlePreAccount(id int) (err error) {
	o := orm.NewOrm()

	v := SettlePreAccount{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SettlePreAccount{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
