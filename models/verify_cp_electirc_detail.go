package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type VerifyCpElectricDetail struct {
	Id         int     `json:"id,omitempty" orm:"column(id);auto"`
	ElectricId int64   `json:"electric_id,omitempty" orm:"column(electric_id);null"`
	Date       string  `json:"date,omitempty" orm:"column(date);size(20);null"`	//所属月份
	CompanyId  int	   `json:"company_id,omitempty" orm:"column(company_id);null"`
	GameId     int     `json:"game_id,omitempty" orm:"column(game_id);null"`
	Amount     float64 `json:"amount,omitempty" orm:"column(amount);null"`		//订单金额
	Rate       float64 `json:"rate,omitempty" orm:"column(rate);null"`		//渠道费率
	TaxRate    float64 `json:"tax_rate,omitempty" orm:"column(tax_rate);null"`	//税率
	Ratio      float64 `json:"ratio,omitempty" orm:"column(ratio);null"`		//分成比例

	GameName   string  `json:"game_name,omitempty" orm:"-"`
	Ladders    string  `json:"ladders,omitempty" orm:"-"`
}

func (t *VerifyCpElectricDetail) TableName() string {
	return "verify_cp_electric_detail"
}

func init() {
	orm.RegisterModel(new(VerifyCpElectricDetail))
}

// AddVerifyCpElectircDetail insert a new VerifyCpElectircDetail into database and returns
// last inserted Id on success.
func AddVerifyCpElectircDetail(m *VerifyCpElectricDetail) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetVerifyCpElectircDetailById retrieves VerifyCpElectircDetail by Id. Returns error if
// Id doesn't exist
func GetVerifyCpElectircDetailById(id int) (v *VerifyCpElectricDetail, err error) {
	o := orm.NewOrm()
	v = &VerifyCpElectricDetail{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllVerifyCpElectircDetail retrieves all VerifyCpElectircDetail matches certain condition. Returns empty list if
// no records exist
func GetAllVerifyCpElectircDetail(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(VerifyCpElectricDetail))
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

	var l []VerifyCpElectricDetail
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

// UpdateVerifyCpElectircDetail updates VerifyCpElectircDetail by Id and returns error if
// the record to be updated doesn't exist
func UpdateVerifyCpElectircDetailById(m *VerifyCpElectricDetail) (err error) {
	o := orm.NewOrm()
	v := VerifyCpElectricDetail{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteVerifyCpElectircDetail deletes VerifyCpElectircDetail by Id and returns error if
// the record to be deleted doesn't exist
func DeleteVerifyCpElectircDetail(id int) (err error) {
	o := orm.NewOrm()
	v := VerifyCpElectricDetail{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&VerifyCpElectricDetail{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
