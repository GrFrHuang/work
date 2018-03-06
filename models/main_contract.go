package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"github.com/astaxie/beego/orm"
	"time"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/utils"
)

type MainContract struct {
	Id             int    `json:"id" orm:"column(id);auto"`
	CompanyType    int    `json:"company_type" orm:"column(company_type);null"` //1标识cp商，2标识渠道
	CompanyId      int    `json:"company_id" orm:"column(company_id);null"`
	BodyMy         int    `json:"body_my" orm:"column(body_my);null"`
	BeginTime      string `json:"begin_time" orm:"column(begin_time);size(10);null"`
	EndTime        string `json:"end_time" orm:"column(end_time);size(10);null"`
	SignTime       string `json:"sign_time" orm:"column(sign_time);size(10);null"`
	State          int    `json:"state" orm:"column(state);null"`
	FileId         string `json:"file_id" orm:"column(file_id);size(255);null"`
	SignPerson     int    `json:"sign_person" orm:"column(sign_person);null"`
	UpdatePerson   int    `json:"update_person" orm:"column(update_person);null"`
	UpdateTime     int64  `json:"update_time" orm:"column(update_time);null"`
	GameIds        string `json:"game_ids" orm:"column(game_ids);size(255);null"`
	Desc           string `json:"desc" orm:"column(desc);null"`
	EffectiveState int    `json:"effective_state" orm:"column(effective_state);null"`
}

func (t *MainContract) TableName() string {
	return "main_contract"
}

func init() {
	orm.RegisterModel(new(MainContract))
}

func AddGameList(gameid string, companyId int, companyType int) (err error) {
	gameids := strings.Split(gameid, ",")
	if len(gameids) == 0 {
		return nil
	}

	o := orm.NewOrm()
	sql := ""

	if companyType == 1 {
		// cp
		sql = "SELECT a.* FROM contract a LEFT JOIN game b ON a.game_id = b.game_id WHERE a.company_type = 0 AND " +
			"b.issue = ? AND a.game_id=? AND effective_state=1 LIMIT 1"
	} else {
		// 渠道
		sql = "SELECT a.* FROM contract a LEFT JOIN channel_company b ON a.channel_code = b.channel_code WHERE " +
			"a.company_type = 1 AND b.company_id=? AND a.game_id=? AND a.effective_state=1 LIMIT 1"
	}

	for _, id := range gameids {
		contract := Contract{}
		o.Raw(sql, companyId, id).QueryRow(&contract)
		if contract.IsMain == 2 {
			continue
		}
		contract.IsMain = 2
		_, err = o.Update(&contract)
		if err != nil {
			return
		}
	}
	return
}

// AddMainContract insert a new MainContract into database and returns
// last inserted Id on success.
func AddMainContract(m *MainContract) (id int64, err error) {
	o := orm.NewOrm()

	if o.QueryTable(new(MainContract)).Filter("CompanyType", m.CompanyType).Filter("CompanyId", m.CompanyId).Exist() {
		return 0, errors.New("已存在该主合同，请勿重复添加!")
	}
	m.UpdateTime = time.Now().Unix()
	id, err = o.Insert(m)

	var page int
	if m.CompanyType == 1 { //cp合同
		page = bean.OPP_CP_MAIN_CONTRACT
	} else if m.CompanyType == 2 { //渠道合同
		page = bean.OPP_CHANNEL_MAIN_CONTRACT
	}

	err = CompareAndAddOperateLog(nil, m, m.UpdatePerson, page, int(id), bean.OPA_INSERT)

	return
}

// GetMainContractById retrieves MainContract by Id. Returns error if
// Id doesn't exist
func GetMainContractById(id int) (v *MainContract, err error) {
	o := orm.NewOrm()
	v = &MainContract{Id: id, EffectiveState: 1}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, errors.New("请求错误")
}

// 根据channel_code获取该渠道所有的游戏
func GetGameName(company_id int, typ int) (ss []Game, err error) {
	o := orm.NewOrm()
	sql := ""
	if typ == 2 {
		sql = "SELECT b.game_id,c.game_name FROM channel_company a LEFT JOIN channel_access b ON a.channel_code=b.channel_code " +
			"LEFT JOIN game c ON b.game_id=c.game_id WHERE a.company_id=?"
	} else {
		sql = "SELECT a.game_id,a.game_name FROM game a WHERE a.issue=? and a.game_id>0"
	}

	_, err = o.Raw(sql, company_id).QueryRows(&ss)
	return
}

// GetAllMainContract retrieves all MainContract matches certain condition. Returns empty list if
// no records exist
func GetAllMainContract(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(MainContract))
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

	var l []MainContract
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

// UpdateMainContract updates MainContract by Id and returns error if
// the record to be updated doesn't exist
func UpdateMainContractById(m *MainContract) (err error) {
	o := orm.NewOrm()
	v := MainContract{Id: m.Id}
	// ascertain id exists in the database
	fields := utils.GetNotEmptyFields(m, "BodyMy", "SignTime", "BeginTime", "EndTime", "State", "FileId",
		"UpdatePerson", "UpdateTime", "GameIds")
	//当Desc之前有值，但是前端传空想删除此字段的值时，上面的函数并不能完成修改，所以需要加下一行代码修改Desc
	fields = append(fields, "Desc")

	if err = o.Read(&v); err != nil {
		return
	}
	if _, err = o.Update(m); err != nil {
		return
	}

	fields = utils.RemoveFields(fields, "UpdatePerson", "UpdateTime")

	var page int
	if m.CompanyType == 1 { //cp合同
		page = bean.OPP_CP_MAIN_CONTRACT
	} else if m.CompanyType == 2 { //渠道合同
		page = bean.OPP_CHANNEL_MAIN_CONTRACT
	}

	if err = CompareAndAddOperateLog(&v, m, m.UpdatePerson, page, m.Id, bean.OPA_UPDATE, fields...); err != nil {
		return
	}
	return
}

func GetMainContractByCompanyId(companyId int, companyType int) (mainContract *MainContract) {
	o := orm.NewOrm()
	main := MainContract{}
	o.QueryTable(new(MainContract)).Filter("CompanyType", companyType).Filter("CompanyId", companyId).
		Filter("EffectiveState", 1).One(&main)
	return &main
}

// DeleteMainContract deletes MainContract by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMainContract(id int) (err error) {
	o := orm.NewOrm()
	v := MainContract{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&MainContract{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
