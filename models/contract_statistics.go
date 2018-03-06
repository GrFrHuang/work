package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"time"
)

type ContractStatistics struct {
	Id       int    `json:"id,omitempty" orm:"column(id);auto"`
	Type	 int	`json:"type,omitempty" orm:"column(type)"`
	Time     string `json:"time,omitempty" orm:"column(time);size(255)"`
	User     int    `json:"user,omitempty" orm:"column(user)"`
	Send     int    `json:"send" orm:"column(send)"`
	Complete int    `json:"complete" orm:"column(complete)"`

	People     *User   `orm:"-" json:"people,omitempty"`
}

func (t *ContractStatistics) TableName() string {
	return "contract_statistics"
}

func init() {
	orm.RegisterModel(new(ContractStatistics))
}

// AddContractStatistics insert a new ContractStatistics into database and returns
// last inserted Id on success.
func AddContractStatistics(m *ContractStatistics) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetContractStatisticsById retrieves ContractStatistics by Id. Returns error if
// Id doesn't exist
func GetContractStatisticsById(id int) (v *ContractStatistics, err error) {
	o := orm.NewOrm()
	v = &ContractStatistics{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllContractStatistics retrieves all ContractStatistics matches certain condition. Returns empty list if
// no records exist
func GetAllContractStatistics(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ContractStatistics))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []ContractStatistics
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateContractStatistics updates ContractStatistics by Id and returns error if
// the record to be updated doesn't exist
func UpdateContractStatisticsById(m *ContractStatistics) (err error) {
	o := orm.NewOrm()
	v := ContractStatistics{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// 当合同状态改为已寄出，未回寄或已签订，则更新合同统计表
func UpdateContractStatistics(userId int, state int, contractType int) (err error) {

	m := ContractStatistics{}
	o := orm.NewOrm()
	timeNow := time.Now().Format("2006-01-02")
	var typ int
	if(contractType == 0){
		typ = 1
	}else{
		typ = 2
	}

	err =  o.QueryTable(new(ContractStatistics)).Filter("time__exact", timeNow).Filter("type__exact", typ).Filter("user__exact", userId).One(&m)

	//如果数据库没有记录，则创建
	if err == orm.ErrNoRows{
		m.Time = timeNow
		m.User = userId
		m.Type = typ

		if(state == 152){//已寄出，未回寄
			m.Send = 1
			m.Complete = 0
		}else if(state == 150){//已签订
			m.Send = 0
			m.Complete = 1
		}

		_, err = o.Insert(&m)
	}else{//有记录，需对原有数据修改
		if(state == 152){//已寄出，未回寄
			o.QueryTable(new(ContractStatistics)).Filter("time__exact", timeNow).Filter("type__exact", typ).Filter("user__exact", userId).Update(orm.Params{
				"send": orm.ColValue(orm.ColAdd, 1),
			})
		}else if(state == 150){//已签订
			o.QueryTable(new(ContractStatistics)).Filter("time__exact", timeNow).Filter("type__exact", typ).Filter("user__exact", userId).Update(orm.Params{
				"complete": orm.ColValue(orm.ColAdd, 1),
			})
		}
	}

	return
}

func GetAllContractStatisticWithTotal(result *[]ContractStatistics, typ int, offset int64, limit int64) (total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ContractStatistics)).Filter("type__exact", typ).GroupBy("time")

	total, err = qs.Count()
	if err != nil{
		return
	}

	//rs := []orm.Params{}
	o.Raw("select time, sum(send) as send, sum(complete) as complete from contract_statistics where type = ? group by time order by time desc limit ?, ?", typ, offset, limit).QueryRows(result)

	return
}

//获取某天的合同统计详情
func GetDetails(result *[]ContractStatistics, typ int, time string) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(ContractStatistics)).Filter("type__exact", typ).Filter("time__exact", time).All(result)

	return
}

func AddPeopleInfo(ss *[]ContractStatistics) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.User
	}
	games := []User{}
	_, err := o.QueryTable(new(User)).
		Filter("Id__in", linkIds).
		All(&games, "Id", "Name", "NickName")
	if err != nil {
		return
	}
	gameMap := map[int]User{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.User]; ok {
			(*ss)[i].People = &g
		}
	}
	return
}

// DeleteContractStatistics deletes ContractStatistics by Id and returns error if
// the record to be deleted doesn't exist
func DeleteContractStatistics(id int) (err error) {
	o := orm.NewOrm()
	v := ContractStatistics{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ContractStatistics{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
