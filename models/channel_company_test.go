package models

import (
	"testing"
	"github.com/astaxie/beego/orm"
	"encoding/json"
	"fmt"
)

//获取所有错误的回款公司主体
/*
func TestCheckErrorRemit(t *testing.T) {
	companies := []Company{}

	cs := []ChannelCompany{}
	qs := orm.NewOrm().QueryTable(new(ChannelCompany))
	_, err := qs.All(&cs)
	if err != nil {
		return
	}
	if len(cs) == 0 {
		return
	}
	allRemitCompanyIds := []int{}
	for _, c := range cs {
		if c.RemitCompany != "" {
			remitCompanyIds := []int{}
			e := json.Unmarshal([]byte(c.RemitCompany), &remitCompanyIds)
			if e != nil {
				fmt.Println(e.Error())
				continue
			}

			for _, v := range remitCompanyIds {
				allRemitCompanyIds = append(allRemitCompanyIds, v)
			}
		}
	}
	if len(allRemitCompanyIds) == 0 {
		return
	}
	fmt.Println(allRemitCompanyIds)
	_, err = orm.NewOrm().QueryTable(new(Company)).
		Filter("id__in", allRemitCompanyIds).
		Filter("type__lt",3).
		All(&companies, "id", "name")
	if err != nil {
		return
	}
	fmt.Println(companies)
}
*/

//检查所有的对账单里的回款公司
func TestCheckVerifyRemitCompany(t *testing.T) {
	orm.Debug = false
	cs := []VerifyChannel{}
	qs := orm.NewOrm().QueryTable(new(VerifyChannel))
	_, err := qs.Filter("status__gte", 30).All(&cs)
	if err != nil {
		return
	}
	for _, c := range cs {
		if c.RemitCompanyId != 0 {
			css := ChannelCompany{}
			qss := orm.NewOrm().QueryTable(new(ChannelCompany))
			err := qss.Filter("channel_code", c.ChannelCode).One(&css, "remit_company")
			if err != nil {
				return
			}
			remitCompanyIds := []int{}
			e := json.Unmarshal([]byte(css.RemitCompany), &remitCompanyIds)
			if e != nil {
				fmt.Println(e.Error())
				continue
			}
			flag := false
			for _, v := range remitCompanyIds {
				if v == c.RemitCompanyId {
					flag = true
					break
				}
			}
			if flag == false {
				//fmt.Println("err dat")
				fmt.Println(c.Id)
			}
		} else {
			fmt.Println("remitcompanyid err")
		}
	}
}

func TestGetCompanyTypeById(t *testing.T) {
	c, err := GetCompanyTypeById(131)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(c)
}

// TestInsertRemitCompany 根据渠道商表的RemitCompany字段更新RemitCompany表
// 注意：需要保证渠道商中的channel_code字段唯一
func TestInsertRemitCompany(t *testing.T) {
	o := orm.NewOrm()
	var users []*ChannelCompany
	_, err := o.QueryTable("channel_company").All(&users)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range users {
		fmt.Println(*v)
		var remits []int
		err := json.Unmarshal([]byte(v.RemitCompany), &remits)
		if err != nil || len(remits) == 0 || v.ChannelCode == "" {
			continue
		}
		for _, companyId := range remits {
			remitCompany := RemitCompany{}
			remitCompany.ChannelCode = v.ChannelCode
			remitCompany.CompanyId = companyId
			AddRemitCompany(&remitCompany)
		}
	}
}

func TestMap(t *testing.T) {
	//GetContactsByCompanyId()
}
