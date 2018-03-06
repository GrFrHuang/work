package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/utils"
	"reflect"
	"strings"
	"time"
)

type Contract struct {
	Id          int `json:"id,omitempty" orm:"column(id);auto"`
	GameId      int `json:"game_id,omitempty" orm:"column(game_id);null"`
	CompanyType int `json:"company_type,omitempty" orm:"column(company_type);null"`
	//CompanyId    int     `json:"company_id,omitempty" orm:"column(company_id);null"`
	ChannelCode string `json:"channel_code,omitempty" orm:"column(channel_code);null"`
	BodyMy      int    `json:"body_my,omitempty" orm:"column(body_my);null"` // 我方主体
	//Advance      float64 `json:"advance,omitempty" orm:"column(advance);null"`
	SigningTime    int64  `json:"signing_time,omitempty" orm:"column(signing_time);null"`
	BeginTime      int64  `json:"begin_time,omitempty" orm:"column(begin_time);null"`
	EndTime        int64  `json:"end_time,omitempty" orm:"column(end_time);null"`
	State          int    `json:"state" orm:"column(state);null"`
	Number         int    `json:"number,omitempty" orm:"column(number);size(20);null"`
	CreatePerson   int    `json:"create_person,omitempty" orm:"column(create_person);null"`
	UpdatePerson   int    `json:"update_person,omitempty" orm:"column(update_person);null"`
	Desc           string `json:"desc,omitempty" orm:"column(desc);null"`
	CreateTime     int64  `json:"create_time,omitempty" orm:"column(create_time);null"`
	UpdateTime     int64  `json:"update_time,omitempty" orm:"column(update_time);null"`
	FileId         string `json:"file_id,omitempty" orm:"column(file_id);null" valid:"Required"`
	Ladder         string `json:"ladder_new" orm:"column(ladders);null" valid:"Required"` //阶梯,映射新系统数据库
	EffectiveState int    `json:"effective_state" orm:"column(effective_state);null" valid:"Required"`
	IsMain         int    `json:"is_main" orm:"column(is_main);null"` //是否在主合同中，1:否，2:是
	Accessory		string `json:"accessory" orm:"column(accessory);null"`

	Game               *Game        `orm:"-" json:"game,omitempty"`
	Channel            *Channel     `orm:"-" json:"channel,omitempty"`
	Development        *CompanyType `orm:"-" json:"development,omitempty"`
	UpdateUser         *User        `orm:"-" json:"update_user,omitempty"`
	Status             *Types       `orm:"-" json:"status,omitempty"`
	Ladder_front       interface{}  `orm:"-" json:"ladder_front"`
	Express            *Express     `orm:"-" json:"express,omitempty"`
	Ladders            string       `orm:"-" json:"ladders,omitempty"`              // 阶梯, 在这里提交但是直接提交给老系统
	ChannelCompanyName string       `orm:"-" json:"channel_company_name,omitempty"` // 阶梯, 在这里提交但是直接提交给老系统
	Business           string       `orm:"-" json:"business,omitempty"`             // 渠道合同的商务负责人
}

func (t *Contract) TableName() string {
	return "contract"
}

func init() {
	orm.RegisterModel(new(Contract))
}

// AddContract insert a new Contract into database and returns
// last inserted Id on success.
func AddContract(m *Contract, gameId int, companyType int /*companyId int,*/ , channelCode string, createPerson int,accessory string) (id int64, err error) {
	o := orm.NewOrm()
	m.GameId = gameId
	m.CompanyType = companyType

	if companyType == 1 {
		m.ChannelCode = channelCode
	}
	if m.CreateTime == 0 {
		m.CreateTime = time.Now().Unix()
	}
	m.UpdateTime = time.Now().Unix()
	if m.State == 0 {
		m.State = 149
	}
	m.EffectiveState = 1
	m.CreatePerson = createPerson
	m.FileId = "[]"
	m.Accessory = accessory
	id, err = o.Insert(m)

	var page int
	if companyType == 0 { //cp合同
		page = bean.OPP_CP_CONTRACT
	} else if companyType == 1 { //渠道合同
		page = bean.OPP_CHANNEL_CONTRACT
	}
	err = CompareAndAddOperateLog(nil, m, createPerson, page, int(id), bean.OPA_INSERT)
	return
}

func GetContractByGameid(gameId int, channelId string) (int, error) {
	o := orm.NewOrm()
	var m Contract
	m.GameId = gameId
	m.ChannelCode = channelId
	if err := o.Read(&m, "game_id", "channel_code"); err != nil {
		return 0, err
	}
	return m.Id, nil
}

func DelContractByid(id int) (error) {
	o := orm.NewOrm()
	i, err := o.Delete(&Contract{Id: id})
	if i <= 0 || err != nil {
		return errors.New("删除合同信息失败")
	}
	return nil
}

// GetContractById retrieves Contract by Id. Returns error if
// Id doesn't exist
func GetContractById(id int, flag string) (v *Contract, err error) {
	o := orm.NewOrm()
	v = &Contract{Id: id}
	if err = o.Read(v); err != nil {
		return
	}

	// 附加阶梯信息
	var f interface{}
	json.Unmarshal([]byte(v.Ladder), &f)
	v.Ladder_front = f

	//附加游戏信息
	game := Game{}
	err = o.QueryTable(new(Game)).
		Filter("game_id__in", v.GameId).
		One(&game, "Id", "GameId", "GameName")
	if err != nil {
		return
	} else {
		v.Game = &game
	}

	//附加快递信息
	express := Express{}
	err = o.QueryTable(new(Express)).
		Filter("contract_id__in", v.Id).
		One(&express)
	if err != nil && err != orm.ErrNoRows {
		return
	} else {
		v.Express = &express
	}

	//附加发行商或渠道信息
	if flag == "cp" {
		issue := CompanyType{}
		err = o.Raw("SELECT id,name FROM company_type WHERE id = (SELECT issue FROM game WHERE game_id=( SELECT game_id FROM contract WHERE id = ?))",
			v.Id).QueryRow(&issue)
		//(*ss)[i].ChannelCompanyName = company.Name
		//err = o.QueryTable(new(Company)).
		//	Filter("id__in", v.CompanyId).
		//	One(&issue, "Id", "Name")
		if err != nil {
			return
		} else {
			v.Development = &issue
		}
	} else if flag == "qd" {
		channel := Channel{}
		err = o.QueryTable(new(Channel)).
			Filter("cp__in", v.ChannelCode).
			One(&channel, "ChannelId", "Name", "Cp")
		if err != nil {
			return
		} else {
			v.Channel = &channel
		}
	}

	return
}

func GetAllEditIds(editId int) (ids []orm.Params, err error) {
	o := orm.NewOrm()

	_, err = o.Raw("SELECT b.id,b.state,b.effective_state FROM contract a INNER JOIN contract b ON a.game_id=b.game_id "+
		"AND a.company_type=b.company_type AND a.channel_code=b.channel_code WHERE a.id=? ORDER BY b.id ASC ",
		editId).Values(&ids)

	return
}

// GetAllContract retrieves all Contract matches certain condition. Returns empty list if
// no records exist
func GetAllContract(query map[string][]string, fields []string, sortby []string, order []string,
	offset int64, limit int64, where map[string][]interface{}, companyType int64) (ml []interface{}, count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Contract))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		//k = strings.Replace(k, ".", "__", -1)
		//if strings.Contains(k, "isnull") {
		//	qs = qs.Filter(k, (v == "true" || v == "1"))
		//} else {
		//	qs = qs.Filter(k, v)
		//}
		fmt.Printf("k:%v    v:%v\n", k, v)
	}

	for k, v := range where {
		qs = qs.Filter(k, v...)
	}
	qs = qs.Filter("CompanyType", companyType)
	count, _ = qs.Count()

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
					return nil, count, errors.New("Error: Invalid order. Must be either [asc|desc]")
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
					return nil, count, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, count, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, count, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Contract
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
		return ml, count, nil
	}
	return nil, count, err
}

// UpdateContract updates Contract by Id and returns error if
// the record to be updated doesn't exist
func UpdateContractById(m *Contract) (err error) {
	o := orm.NewOrm()
	m.UpdateTime = time.Now().Unix()
	// 存在直接在原合同上修改到期时间达到续签目的的情况，如原合同状态为已到期，则将原合同状态修改为已签订
	if m.EndTime > time.Now().Unix() && m.State == 151 {
		m.State = 150
	}
	//if m.State == 1 {
	//	contractNumber := new(ContractNumber)
	//	contractNumber.ContractId = m.Id
	//	contractNumber.FileId = m.FileId
	//	number, err := AddContractNumber(contractNumber)
	//	if err == nil {
	//		tmp, _ := strconv.Atoi(strconv.FormatInt(number, 10))
	//		m.Number = tmp
	//	}
	//}
	fields := utils.GetNotEmptyFields(m, "Advance", "SigningTime", "BeginTime", "EndTime", "State", "BodyMy",
		"Number", "FileId", "UpdatePerson", "Ladder", "UpdateTime")
	//当Desc之前有值，但是前端传空想删除此字段的值时，上面的函数并不能完成修改，所以需要加下一行代码修改Desc
	fields = append(fields, "Desc")

	s := Contract{Id: m.Id}
	if err = orm.NewOrm().Read(&s); err != nil {
		return
	}

	_, err = o.Update(m, fields...)
	if err != nil {
		return
	}

	var page int
	if s.CompanyType == 0 { //cp合同
		page = bean.OPP_CP_CONTRACT
	} else if s.CompanyType == 1 { //渠道合同
		page = bean.OPP_CHANNEL_CONTRACT
	}

	fields = utils.RemoveFields(fields, "UpdatePerson", "UpdateTime")

	if err = CompareAndAddOperateLog(&s, m, m.UpdatePerson, page, s.Id, bean.OPA_UPDATE, fields...); err != nil {
		return
	}

	return
}

//根据editId 把该合同的历史合同改为无效，并生成新的合同返回
func Renew(editId int, userId int) (newId int64, err error) {
	o := orm.NewOrm()
	con := Contract{Id: editId}
	if err = o.Read(&con); err != nil {
		return
	}

	//由于每次前端传过来的editId都是最新的有效合同的editId，所以只需要修改此一条记录的有效状态
	oldCon := con
	oldCon.EffectiveState = 2
	oldCon.UpdatePerson = userId
	oldCon.UpdateTime = time.Now().Unix()

	if _, err = o.Update(&oldCon); err != nil {
		return
	}
	var page int
	if oldCon.CompanyType == 0 { //cp合同
		page = bean.OPP_CP_CONTRACT
	} else if oldCon.CompanyType == 1 { //渠道合同
		page = bean.OPP_CHANNEL_CONTRACT
	}
	if err = CompareAndAddOperateLog(&con, &oldCon, userId, page, editId, bean.OPA_UPDATE); err != nil {
		return
	}

	newCon := Contract{}
	newCon.GameId = con.GameId
	newCon.CompanyType = con.CompanyType
	newCon.ChannelCode = con.ChannelCode
	newCon.BodyMy = con.BodyMy
	newCon.State = 149
	newCon.CreatePerson = userId
	newCon.CreateTime = time.Now().Unix()
	newCon.UpdateTime = time.Now().Unix()
	newCon.UpdatePerson = userId
	newCon.EffectiveState = 1
	newCon.IsMain = oldCon.IsMain
	if newId, err = o.Insert(&newCon); err != nil {
		return
	}

	err = CompareAndAddOperateLog(nil, newCon, userId, page, int(newId), bean.OPA_INSERT)
	return
}

func UpdateCpContractBody(body int, gameId int) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("update contract set body_my = ? where game_id = ? and company_type = 0", body, gameId).Exec()

	return
}

func UpdateChannelContractBody(body int, gameId int, channel_code string) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("update contract set body_my = ? where game_id = ? and company_type = 1 and channel_code = ?", body,
		gameId, channel_code).Exec()

	return
}

// 根据gameId，channel_code获取具体有效渠道合同
func GetContractByChannelCodeAndGameId(game_id int, channel_code string) (con *Contract) {
	o := orm.NewOrm()
	con = &Contract{}
	o.QueryTable(new(Contract)).Filter("channel_code", channel_code).Filter("game_id", game_id).Filter("company_type", 1).
		Filter("effective_state", 1).One(con)
	return
}

// 根据游戏id，和合同类型（0:cp,1:渠道）,获取该游戏所有合同
func GetContractByGameAndType(gameId int, companyType int) (cons []Contract) {
	o := orm.NewOrm()

	qs := o.QueryTable(new(Contract)).Filter("game_id__exact", gameId).Filter("company_type__exact", companyType).
		Filter("effective_state__exact", 1)
	qs.All(&cons)

	return cons
}

// DeleteContract deletes Contract by Id and returns error if
// the record to be deleted doesn't exist
func DeleteContract(id int) (err error) {
	o := orm.NewOrm()
	v := Contract{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Contract{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GroupRemitAddGameInfo(ss *[]Contract) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.GameId
	}
	var link []Game
	_, err := o.QueryTable(new(Game)).
		Filter("game_id__in", linkIds).
		All(&link, "GameId", "GameName")
	if err != nil {
		return
	}
	linkMap := map[int]Game{}
	for _, g := range link {
		linkMap[g.GameId] = g
	}
	for i, s := range *ss {
		if g, ok := linkMap[s.GameId]; ok {
			(*ss)[i].Game = &g
		}
	}
	return
}

func GroupRemitAddCpInfo(ss *[]Contract) {

	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()

	for i, s := range *ss {
		var company CompanyType
		o.Raw("SELECT name FROM company_type WHERE id = (SELECT issue FROM game WHERE game_id=( SELECT game_id FROM contract WHERE id = ?))",
			s.Id).QueryRow(&company)
		(*ss)[i].ChannelCompanyName = company.Name
	}

	return
}

func AddCpBodyMyInfo(ss *[]Contract) {

	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()

	for i, s := range *ss {
		var game Game
		o.Raw("SELECT body_my FROM game WHERE game_id = ?", s.GameId).QueryRow(&game)
		(*ss)[i].BodyMy = game.BodyMy
	}

	return
}

func GetAllCpContract(where map[string][]interface{}, limit int64, offset int64) (data []orm.Params, total int64, err error) {

	games := where["game_id__in"]
	state := where["state__exact"]
	bodyMy := where["body_my__exact"]

	fmt.Printf("where:%v, games: %v, state: %v, bodyMy: %v\n", where, games, state, bodyMy)

	condition := ""
	var cond []interface{}

	if games != nil && len(games) > 0 {
		holder := strings.Repeat(",?", len(games))
		condition = fmt.Sprintf("%sAND game_id in(%s) ", condition, holder[1:])
		cond = append(cond, games)
	}

	if state != nil && len(state) > 0 {
		condition = fmt.Sprintf("%sAND `state` = ? ", condition)
		cond = append(cond, fmt.Sprintf("%s", state[0]))
	} else {
		state = []interface{}{" "}

	}

	if bodyMy != nil && len(bodyMy) > 0 {
		condition = fmt.Sprintf("%sAND `b.body_my` = ? ", condition)
		cond = append(cond, fmt.Sprintf("%s", bodyMy[0]))
	} else {
		bodyMy = []interface{}{" "}

	}

	if len(condition) > 3 {
		condition = " WHERE " + condition[3:] + " ORDER BY a.id desc "
	}

	fmt.Printf("condition:%v, cond: %v\n", condition, cond)

	o := orm.NewOrm()
	var total_maps []orm.Params
	sql_str_total := fmt.Sprintf("SELECT COUNT(*) AS total FROM contract AS a LEFT JOIN game AS b ON a.game_id = b.game_id %s", condition)
	_, _ = o.Raw(sql_str_total, cond...).Values(&total_maps)

	//o := orm.NewOrm()
	////sql_str_total := fmt.Sprintf("SELECT COUNT(*) AS total , SUM(o.money) AS money FROM (SELECT SUM(amount) AS money FROM `order` %s GROUP BY %s) AS o", condition, group_fileds)
	//sql := "SELECT a.*, b.body_my FROM contract AS a LEFT JOIN game AS b ON a.game_id = b.game_id " +
	//	"WHERE a.company_type = 0 "
	//
	//
	//
	//rs := []orm.Params{}
	//_, err = o.Raw(sql,offset,limit).Values(&rs)
	//if err != nil {
	//	return
	//}

	return
}

func GroupRemitContractAddChannelInfo(ss *[]Contract) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]string, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.ChannelCode
	}
	var link []Channel
	_, err := o.QueryTable(new(Channel)).
		Filter("cp__in", linkIds).
		All(&link, "Channelid", "Name", "Cp")
	if err != nil {
		return
	}
	linkMap := map[string]Channel{}
	for _, g := range link {
		linkMap[g.Cp] = g
	}
	for i, s := range *ss {
		if g, ok := linkMap[s.ChannelCode]; ok {
			(*ss)[i].Channel = &g
		}
	}
	return
}

// 录入游戏停运后，根据停服时间，修改该游戏的合同状态
func UpdateContractState(game_id int, serverTime int64) (err error) {
	// cp合同
	cpContract := GetContractByGameAndType(game_id, 0)
	if cpContract[0].EndTime < serverTime && cpContract[0].State != 155 && cpContract[0].State != 156 {
		cpContract[0].State = 157 // 修改为即将停运
		UpdateContractById(&cpContract[0])
	}

	// 渠道合同
	channelContracts := GetContractByGameAndType(game_id, 1)
	for _, con := range channelContracts {
		if con.EndTime < serverTime && cpContract[0].State != 155 && cpContract[0].State != 156 {
			con.State = 157 // 修改为即将停运
			UpdateContractById(&con)
		}
	}

	return nil
}

//渠道合同附加商务负责人信息学
func AddBusinessInfo(ss *[]Contract) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()

	for i, s := range *ss {
		var business User
		o.Raw("SELECT nickname FROM user WHERE id = (SELECT business_person FROM channel_access WHERE "+
			"game_id = ? AND channel_code = ?)", s.GameId, s.ChannelCode).QueryRow(&business)
		(*ss)[i].Business = business.Nickname
	}
	//linkIds := make([]int, len(*ss))
	//for i, s := range *ss {
	//	linkIds[i] = s.BusinessPerson
	//}
	//games := []User{}
	//_, err := o.QueryTable(new(User)).
	//	Filter("Id__in", linkIds).
	//	All(&games, "Id", "Name", "NickName")
	//if err != nil {
	//	return
	//}
	//gameMap := map[int]User{}
	//for _, g := range games {
	//	gameMap[g.Id] = g
	//}
	//for i, s := range *ss {
	//	if g, ok := gameMap[s.BusinessPerson]; ok {
	//		(*ss)[i].Business = &g
	//	}
	//}
	return
}

func AddChannelCompanyInfo(ss *[]Contract) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()

	for i, s := range *ss {
		var company CompanyType
		o.Raw("SELECT name FROM company_type WHERE id = (SELECT company_id FROM channel_company WHERE channel_code = ?)",
			s.ChannelCode).QueryRow(&company)
		(*ss)[i].ChannelCompanyName = company.Name
	}

	return
}

func GroupRemitAddContractUserInfo(ss *[]Contract) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	linkIds := make([]int, len(*ss))
	for i, s := range *ss {
		linkIds[i] = s.UpdatePerson
	}
	var games []User
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
		if g, ok := gameMap[s.UpdatePerson]; ok {
			(*ss)[i].UpdateUser = &g
		}
	}
	return
}

func GroupRemitAddContractStatusInfo(ss *[]Contract) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()
	gameIds := make([]int, len(*ss))
	for i, s := range *ss {
		gameIds[i] = s.State
	}
	var games []Types
	_, err := o.QueryTable(new(Types)).Filter("Id__in", gameIds).All(&games, "Id", "Name")
	if err != nil {
		return
	}
	gameMap := map[int]Types{}
	for _, g := range games {
		gameMap[g.Id] = g
	}
	for i, s := range *ss {
		if g, ok := gameMap[s.State]; ok {
			(*ss)[i].Status = &g
		}
	}
	return
}

func ParseLadder2Json(ss *[]Contract) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	for i := range *ss {

		//(*ss)[i].Ladder =
		var f interface{}

		json.Unmarshal([]byte((*ss)[i].Ladder), &f)

		(*ss)[i].Ladder_front = f

		//if s.State == 0 {
		//	(*ss)[i].Status = "未签订"
		//} else if s.State == 1 {
		//	(*ss)[i].Status = "已签订"
		//} else if s.State == 2 {
		//	(*ss)[i].Status = "已到期"
		//} else if s.State == 3 {
		//	(*ss)[i].Status = "进行中"
		//} else {
		//	(*ss)[i].Status = "无合作"
		//}
	}
	return
}

// 通过company_id 以及合同状态为 未签订 已到期 的条件 获得游戏 game_id
func GetGameIdByCompanyIdAndState(id int) (contracts []Contract, err error) {
	o := orm.NewOrm()
	_, err = o.Raw("SELECT c.id,c.game_id FROM contract c LEFT JOIN game a ON c.game_id=a.game_id WHERE a.issue = ?"+
		" AND c.state IN (149,151) AND c.company_type=0", id).QueryRows(&contracts)
	if err != nil {
		return
	}
	return contracts, nil

}

//通过 channelCode 获取合同的 未签订 已签订 的条件 获得游戏 game_id
func GetGameIdByChannelCodeAndState(channel_code string) (contracts []Contract, err error) {
	o := orm.NewOrm()
	_, err = o.Raw("SELECT id,game_id FROM contract c WHERE c.channel_code = ? "+
		"AND c.state IN (149,151)", channel_code).QueryRows(&contracts)
	if err != nil {
		return
	}
	return contracts, nil

}
