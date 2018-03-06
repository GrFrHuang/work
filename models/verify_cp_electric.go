package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"strings"
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/tool/old_sys"
)

type VerifyCpElectric struct {
	Id          int    `json:"id,omitempty" orm:"column(id);auto"`
	CompanyId   int    `json:"company_id,omitempty" orm:"column(company_id);null"`
	BodyMy      int    `json:"body_my,omitempty" orm:"column(body_my);null"`
	Desc        string `json:"desc,omitempty" orm:"column(desc);null"`
	ReceiveUser int    `json:"receive_user,omitempty" orm:"column(receive_user);null"`	//账户信息中收件人id
	ContactId   int    `json:"contact_id,omitempty" orm:"column(contact_id);null"`		//合作方收件人id
	UpdateUser  int    `json:"update_user,omitempty" orm:"column(update_user);null"`
	UpdateTime  int64  `json:"update_time,omitempty" orm:"column(update_time);null"`

	Games       []VerifyCpElectricDetail `json:"games,omitempty" orm:"-"`
	Company     *CompanyType `orm:"-" json:"company,omitempty"` // 发行商
	UpdatedUser *User  `orm:"-" json:"updated_user,omitempty"`
	Dates       string `json:"dates,omitempty" orm:"-"`	//对账日期
}

func (t *VerifyCpElectric) TableName() string {
	return "verify_cp_electric"
}

func init() {
	orm.RegisterModel(new(VerifyCpElectric))
}

// AddVerifyCpElectric insert a new VerifyCpElectric into database and returns
// last inserted Id on success.
func AddVerifyCpElectric(m *VerifyCpElectric) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func AddCompanyForVerifyCpElectric(data *[]VerifyCpElectric) (err error) {
	if len(*data) == 0 {
		return
	}
	for i, v := range *data {
		(*data)[i].Company, err = GetCompanyById(v.CompanyId)
		if err != nil {
			return
		}
	}
	return
}

func AddUpdateUserForVerifyCpElectric(data *[]VerifyCpElectric) (err error) {
	if len(*data) == 0 {
		return
	}
	uid := []interface{}{}
	for _, v := range *data {
		uid = append(uid, v.UpdateUser)
	}
	if len(uid) == 0 {
		return
	}
	users := []*User{}
	orm.NewOrm().QueryTable(new(User)).Filter("id__in", uid).All(&users)
	usersMap := map[int]*User{}
	for i, v := range users {
		usersMap[v.Id] = users[i]
	}

	for i, v := range *data {
		(*data)[i].UpdatedUser = usersMap[v.UpdateUser]
	}

	return
}

func AddVerifyDate(data *[]VerifyCpElectric){
	if len(*data) == 0 {
		return
	}

	qb, _ := orm.NewQueryBuilder("mysql")
	for i, v := range *data {
		ss := []VerifyCpElectricDetail{}
		dates := []string{}

		qb.Select("distinct(date)").From("verify_cp_electric_detail").Where("electric_id=?")
		orm.NewOrm().Raw(qb.String(), v.Id).QueryRows(&ss)
		for _, value := range ss{
			dates = append(dates, value.Date)
		}

		(*data)[i].Dates = strings.Join(dates, ",")
	}
}

// GetVerifyCpElectricById retrieves VerifyCpElectric by Id. Returns error if
// Id doesn't exist
func GetVerifyCpElectricById(id int) (v *VerifyCpElectric, err error) {
	o := orm.NewOrm()
	v = &VerifyCpElectric{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// 根据我方主体获取已对账的发行商
func GetVerifyCp(bodyMy int) (companies []CompanyType, err error) {
	// 由于有些游戏没有录入发行商,所以这里可能查出null
	sql := "SELECT DISTINCT(a.company_id) id,b.name FROM verify_cp a LEFT JOIN company_type b ON a.company_id = b.id " +
		"WHERE a.status=20 and a.body_my=?"
	//rs := []orm.Params{}
	o := orm.NewOrm()
	_, err = o.Raw(sql, bodyMy).QueryRows(&companies)
	if err != nil {
		return
	}

	return
}

// 根据 我方主体 和 发行商 获取对账的月份
func GetVerifyCpTime(bodyMy int, companyId int) (months []string, err error) {
	o := orm.NewOrm()
	sql := "SELECT date FROM `verify_cp` WHERE body_my = ? AND status=20 AND company_id = ? "
	rs := []orm.Params{}
	_, err = o.Raw(sql, bodyMy, companyId).Values(&rs)
	if err != nil {
		return
	}

	months = []string{}
	for _, i := range rs {
		months = append(months, i["date"].(string))
	}

	return
}

// 根据 发行商 和 月份  获取没有对账的游戏账单
func GetVerifyCpGame(companyId int, month string) (games []VerifyCpElectricDetail, err error) {
	o := orm.NewOrm()

	months := strings.Split(month, ",")
	for i, mon := range months{
		months[i] = "\"" + mon + "\""
	}

	//sql := fmt.Sprintf("SELECT a.date,a.company_id,a.game_id,b.game_name,a.amount_opposite amount,c.ladders " +
	//	"FROM order_pre_verify_cp a LEFT JOIN game b ON a.game_id=b.game_id " +
	//	"LEFT JOIN contract c ON a.game_id=c.game_id AND c.company_type=0 " +
	//	"WHERE a.company_id=? AND a.date IN(%s) ORDER BY DATE asc", strings.Join(months,","))
	sql := fmt.Sprintf("SELECT a.date,a.company_id,a.game_id,a.amount_opposite amount FROM order_pre_verify_cp a " +
		" WHERE a.company_id=? AND a.date IN(%s) ORDER BY DATE ASC", strings.Join(months,","))
	_, err = o.Raw(sql, companyId).QueryRows(&games)
	if err != nil {
		return
	}

	return
}

//根据阶梯分成，设置对账单的渠道费率和(对方)分成比例
func SetRate(data *[]VerifyCpElectricDetail) (err error) {
	if len(*data) == 0 {
		return
	}
	o := orm.NewOrm()
	ladders := []old_sys.Ladder4Post{}
	contract := Contract{}
	for i, v := range *data {
		o.Raw("select ladders from contract where game_id = ? and company_type=0", v.GameId).QueryRow(&contract)
		json.Unmarshal([]byte(contract.Ladder), &ladders)
		if len(ladders) == 1 {
			(*data)[i].Rate = ladders[0].SlottingFee
			(*data)[i].Ratio = 1 - ladders[0].Ratio
		}
	}
	return
}

//设置对账单的游戏名
func SetGameName(data *[]VerifyCpElectricDetail) (err error) {
	if len(*data) == 0 {
		return
	}
	o := orm.NewOrm()
	game := Game{}
	for i, v := range *data {
		o.Raw("select game_name from game where game_id = ?", v.GameId).QueryRow(&game)
		(*data)[i].GameName = game.GameName
	}
	return
}

// UpdateVerifyCpElectric updates VerifyCpElectric by Id and returns error if
// the record to be updated doesn't exist
func UpdateVerifyCpElectricById(m *VerifyCpElectric) (err error) {
	o := orm.NewOrm()
	v := VerifyCpElectric{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteVerifyCpElectric deletes VerifyCpElectric by Id and returns error if
// the record to be deleted doesn't exist
func DeleteVerifyCpElectric(id int) (err error) {
	o := orm.NewOrm()
	v := VerifyCpElectric{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&VerifyCpElectric{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
