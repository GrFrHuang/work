package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"sync"
	"time"
)

// 存放预付款（有结算但没有对账，则不能往对账单里倒水，则把金额累计到这个表里，当有对账单完成的时候再从这里取）
type RemitPreAccount struct {
	Id             int     `json:"id,omitempty" orm:"column(id);auto"`
	BodyMy         int     `json:"body_my,omitempty" orm:"column(body_my);null"`
	RemitCompanyId int     `json:"remit_company_id,omitempty" orm:"column(remit_company_id);null"`
	UpdatedTime    int64     `json:"updated_time,omitempty" orm:"column(updated_time);null"`
	Amount         float64 `json:"amount,omitempty" orm:"column(amount);null;digits(16);decimals(0)"`
	OffsetAmount   float64 `json:"offset_amount,omitempty" orm:"column(offset_amount);null;digits(16);decimals(0)"`
}

func (t *RemitPreAccount) TableName() string {
	return "remit_pre_account"
}

func init() {
	orm.RegisterModel(new(RemitPreAccount))
}

var lockRemit sync.Mutex

// 当有对账单完成了,或者出来新结账单,就取出余额并注入到未结算的单子
//func DoPushRemitToVerify(remitCompanyId int) (err error) {
//	lockRemit.Lock()
//	defer lockRemit.Unlock()
//
//	o := orm.NewOrm()
//	count := 0
//	for {
//		preAccount := &RemitPreAccount{RemitCompanyId: remitCompanyId}
//		// 没找到渠道,说明没有渠道的已回款的钱,直接返回等到回款
//		e := o.Read(preAccount, "RemitCompanyId")
//		if e != nil {
//			return
//		}
//		if preAccount.Amount == 0 {
//			return
//		}
//
//		maps := []orm.Params{}
//		// 读取对账表中 (回款主体) 和 (应回款和已回款不一样的) 就是未全部回款的
//		_, err = o.Raw("SELECT id, start_time, amount_payable , amount_remit from channel_verify_account WHERE remit_company_id = ? AND amount_payable != amount_remit ORDER BY start_time limit 0,1", remitCompanyId).
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
//		amountSettle := m["amount_remit"].(string)
//		amountPayAbleF, e := strconv.ParseFloat(amountPayAble, 64)
//		if e != nil {
//			err = e
//			return
//		}
//		amountSettleF, e := strconv.ParseFloat(amountSettle, 64)
//		if e != nil {
//			err = e
//			return
//		}
//		not := amountPayAbleF - amountSettleF
//		if not >= preAccount.Amount {
//			// 如果余额装不满对账单，则只装这一个对账单并清空余额
//			nowAccount := preAccount.Amount + amountSettleF
//			if not == preAccount.Amount {
//				nowAccount = amountPayAbleF
//			}
//			cp := &ChannelVerifyAccount{Id: idInt, AmountRemit: nowAccount}
//			_, e := o.Update(cp, "AmountRemit")
//			if e != nil {
//				err = e
//				return
//			}
//			preAccount.Amount = 0
//			_, e = o.Update(preAccount, "Amount")
//			if e != nil {
//				// todo 这里如果错误的话，应该将上面的语句回滚
//				err = e
//				log.Error("SETTLE", e)
//				return
//			}
//			break
//		} else {
//			// 如果余额大于当前账单,则填充满当前账单并重复多次填充剩余账单
//			cp := &ChannelVerifyAccount{Id: idInt, AmountRemit: amountPayAbleF, Status: CHAN_VERIFY_S_REMIT}
//			_, e := o.Update(cp, "AmountRemit", "Status")
//			if e != nil {
//				err = e
//				return
//			}
//
//			preAccount.Amount = preAccount.Amount - not
//			preAccount.UpdatedTime = time.Now().Unix()
//			_, e = o.Update(preAccount, "Amount", "UpdatedTime")
//			if e != nil {
//				// todo 这里如果错误的话，应该将上面的语句回滚
//				err = e
//				log.Error("SETTLE", e)
//				return
//			}
//		}
//		count++
//	}
//
//	return
//}

// 添加channel的余额
var lockRemit2 sync.Mutex

func AddRemitPreAmount(remitCompanyId int, amount float64) (err error) {
	lockRemit2.Lock()
	defer lockRemit2.Unlock()

	o := orm.NewOrm()
	v := &RemitPreAccount{RemitCompanyId: remitCompanyId}
	e := o.Read(v, "RemitCompanyId")
	if e != nil {
		if e == orm.ErrNoRows {
			v.Amount = amount
			v.UpdatedTime = time.Now().Unix()
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
	v.UpdatedTime = time.Now().Unix()
	_, err = o.Update(v, )
	if err != nil {
		return
	}
	return
}

// 更新(amount) 或者 新建一条余额与偏移量
func UpdateOrCreateRemitPre(bodyMy, remitCompanyId int, amount, offset float64) (err error) {
	o := orm.NewOrm()

	rp := RemitPreAccount{}
	o.QueryTable("remit_pre_account").Filter("body_my", bodyMy).Filter("remit_company_id", remitCompanyId).One(&rp)

	if rp.Id == 0 {
		// create
		irp := RemitPreAccount{BodyMy: bodyMy, Amount: amount, OffsetAmount: offset, RemitCompanyId: remitCompanyId, UpdatedTime: time.Now().Unix()}
		_, err = o.Insert(&irp)
		if err != nil {
			return
		}
	} else {
		rp.OffsetAmount = offset
		rp.UpdatedTime = time.Now().Unix()
		rp.Amount = amount
		_, err = o.Update(&rp, "amount", "updated_time", "offset_amount")
		if err != nil {
			return
		}
	}

	return
}

// 累加偏移量
func AddRemitPreAmountOffset(bodyMy, remitCompanyId int, offset float64) (err error) {
	o := orm.NewOrm()

	rp := RemitPreAccount{}
	o.QueryTable("remit_pre_account").
		Filter("body_my", bodyMy).
		Filter("remit_company_id", remitCompanyId).
		One(&rp)

	if rp.Id == 0 {
		// created
		irp := RemitPreAccount{BodyMy: bodyMy, OffsetAmount: offset, RemitCompanyId: remitCompanyId, UpdatedTime: time.Now().Unix()}
		_, err = o.Insert(&irp)
		if err != nil {
			return
		}
	} else {
		rp.OffsetAmount = offset + rp.OffsetAmount
		rp.UpdatedTime = time.Now().Unix()
		_, err = o.Update(&rp, "updated_time", "offset_amount")
		if err != nil {
			return
		}
	}

	return
}

// UpdateRemitPreAccount updates RemitPreAccount by Id and returns error if
// the record to be updated doesn't exist
func UpdateRemitPreAccountById(m *RemitPreAccount) (err error) {
	o := orm.NewOrm()
	v := RemitPreAccount{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRemitPreAccount deletes RemitPreAccount by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRemitPreAccount(id int) (err error) {
	o := orm.NewOrm()

	v := RemitPreAccount{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&RemitPreAccount{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
