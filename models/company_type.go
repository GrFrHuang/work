package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/utils"
	"reflect"
	"strconv"
	"strings"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

type CompanyType struct {
	Id   int    `json:"id,omitempty" orm:"column(id);auto"`
	Name string `json:"name,omitempty" orm:"column(name);size(255)"`
}

func (t *CompanyType) TableName() string {
	return "company_type"
}

func init() {
	orm.RegisterModel(new(CompanyType))
}

// AddCompanyType insert a new CompanyType into database and returns
// last inserted Id on success.
func AddCompanyType(m *CompanyType, userid int) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	err = CompareAndAddOperateLog(nil, m, userid, bean.OPP_COMPANY, int(id), bean.OPA_INSERT)
	return
}

// GetCompanyTypeById retrieves CompanyType by Id. Returns error if
// Id doesn't exist
func GetCompanyTypeById(id int) (v *CompanyType, err error) {
	o := orm.NewOrm()
	v = &CompanyType{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCompanyType retrieves all CompanyType matches certain condition. Returns empty list if
// no records exist
func GetAllCompanyType(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CompanyType))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
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

	var l []CompanyType
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

// UpdateCompanyType updates CompanyType by Id and returns error if
// the record to be updated doesn't exist
func UpdateCompanyTypeById(m *CompanyType, where map[string][]interface{}, userid int) (err error) {
	o := orm.NewOrm()
	v := CompanyType{Id: m.Id}
	qs := QueryTable(m, where).Filter("Id", m.Id)
	if err = qs.One(&v); err != nil {
		return
	}

	if _, err = o.Update(m); err != nil{
		return
	}

	if err = CompareAndAddOperateLog(&v, m, userid, bean.OPP_COMPANY, v.Id, bean.OPA_UPDATE); err != nil{
		return
	}

	//// ascertain id exists in the database
	//if err = o.Read(&v); err == nil {
	//	var num int64
	//	if num, err = o.Update(m); err == nil {
	//		fmt.Println("Number of records updated in database:", num)
	//	}
	//}
	return
}

// DeleteCompanyType deletes CompanyType by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCompanyType(id int, userid int) (err error) {
	o := orm.NewOrm()
	v := CompanyType{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&CompanyType{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
		err = CompareAndAddOperateLog(v, nil, userid, bean.OPP_COMPANY, v.Id, bean.OPA_DELETE)
	}
	return
}

// GetCompanyById retrieves Company by Id. Returns error if
// Id doesn't exist
func GetCompanyById(id int) (v *CompanyType, err error) {
	o := orm.NewOrm()
	v = &CompanyType{Id: id}
	if err = o.Read(v); err != nil {
		return
	}
	return
}

func GetCompanyNameByCompanyId(id int) (name string, err error) {
	table := "companyId2Name"
	key := strconv.Itoa(id)
	name, err = utils.Redis.HMGETOne(table, key)
	if err != nil {
		return
	}

	if name != "" {
		return
	}

	company := []CompanyType{}
	_, err = orm.NewOrm().QueryTable("company_type").All(&company, "id", "name")
	if err != nil {
		return
	}

	data := make(map[string]interface{}, len(company))
	for _, v := range company {
		data[strconv.Itoa(v.Id)] = v.Name
	}

	name, _ = util.Interface2String(data[key], false)
	if name == "" {
		err = errors.New("404")
		return
	}

	err = utils.Redis.HMSETALL(table, data, 2*60)
	if err != nil {
		return
	}

	return
}



//获取所有的 公司
func GetCompanyList()(companies []CompanyType, err error)  {
	o := orm.NewOrm()
	if _, err = o.QueryTable(new(CompanyType)).All(&companies); err != nil{
		return
	}
	return companies, nil

}

// 根据 company_id 获得所有的公司信息 返回的是个 []数组
func GetCompaniesByIds(links []int) (wantedCompanies []CompanyType,err error) {
	o := orm.NewOrm()
	if _,err = o.QueryTable(new(CompanyType)).Filter("Id__in",links).All(&wantedCompanies); err != nil{
		return
	}
	return wantedCompanies, nil
}