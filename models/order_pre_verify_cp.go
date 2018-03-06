package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"time"
	"kuaifa.com/kuaifa/work-together/tool/old_sys"
	"encoding/json"
	"github.com/Zeniubius/golang_utils/structUtil"
)

type OrderPreVerifyCp struct {
	Id             int     `orm:"column(id);auto" json:"id,omitempty"`
	GameId         int     `orm:"column(game_id);null" json:"game_id,omitempty"`
	CompanyId      int     `orm:"column(company_id);null" json:"company_id,omitempty"`
	Date           string  `orm:"column(date);size(20);null" json:"date,omitempty"`
	Amount         float64 `orm:"column(amount);null" json:"amount_my"`
	AmountTheory   float64 `orm:"column(amount_theory);null" json:"amount_theory"`
	AmountOpposite float64 `orm:"column(amount_opposite);null" json:"amount_opposite"`
	AmountPayable  float64 `orm:"column(amount_payable);null" json:"amount_payable"`
	VerifyId       int64   `orm:"column(verify_id);null" json:"verify_id,omitempty"`

	GameName string `orm:"-" json:"game_name,omitempty"`
}

func (t *OrderPreVerifyCp) TableName() string {
	return "order_pre_verify_cp"
}

func init() {
	orm.RegisterModel(new(OrderPreVerifyCp))
}

// AddOrderPreVerifyCp insert a new OrderPreVerifyCp into database and returns
// last inserted Id on success.
func AddOrderPreVerifyCp(m *OrderPreVerifyCp) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOrderPreVerifyCpById retrieves OrderPreVerifyCp by Id. Returns error if
// Id doesn't exist
func GetOrderPreVerifyCpById(id int) (v *OrderPreVerifyCp, err error) {
	o := orm.NewOrm()
	v = &OrderPreVerifyCp{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateOrderPreVerifyCp updates OrderPreVerifyCp by Id and returns error if
// the record to be updated doesn't exist
func UpdateOrderPreVerifyCpById(m *OrderPreVerifyCp) (err error) {
	o := orm.NewOrm()
	v := OrderPreVerifyCp{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// 验证预对账单能不能被添加 (已经添加过,或者不存在)
func CheckPreVerifyCpCanAdd(preVerify []OrderPreVerifyCp) (err error) {
	preVerifyIds := make([]interface{}, len(preVerify))
	for i, v := range preVerify {
		preVerifyIds[i] = v.Id
	}

	var ps []OrderPreVerifyCp
	_, err = orm.NewOrm().
		QueryTable("order_pre_verify_cp").
		Filter("Id__in", preVerifyIds).
		All(&ps)
	if err != nil {
		return
	}
	if diff := len(preVerifyIds) - len(ps); diff != 0 {
		err = fmt.Errorf("has %d row not found", diff)
		return
	}
	for _, i := range ps {
		if i.VerifyId != 0 {
			err = fmt.Errorf("preVerify id:%d is verified", i.Id)
			return
		}
	}
	return
}

// 标记预对账单的对账单id
func UpdatePreVerifyCpMarkVerifyId(preVerify []OrderPreVerifyCp, verifyId int) (err error) {
	o := orm.NewOrm()

	for _, p := range preVerify {
		p.VerifyId = int64(verifyId)
		_, err = o.Update(&p, "VerifyId", "AmountTheory", "AmountOpposite", "AmountPayable")
		if err != nil {
			return
		}
	}

	return
}

// 重置预对账单
func UpdatePreVerifyCpReset(verifyId int) (err error) {
	o := orm.NewOrm()
	ps := orm.Params{
		"VerifyId": 0,
	}
	_, err = o.QueryTable("order_pre_verify_cp").Filter("verify_id", verifyId).
		Update(ps)
	return
}

// 同步order表
func UpdatePreVerifyCpFromOrder(companyIds []interface{}, month string) (affCount int64, err error) {
	extendWhere := ""
	var args []interface{}

	if month != "" {
		if extendWhere != "" {
			extendWhere = extendWhere + "AND "
		}
		extendWhere += "LEFT(date, 7) = ? "
		args = append(args, month)
	}
	if extendWhere != "" {
		extendWhere = "WHERE " + extendWhere
	}
	if companyIds != nil {
		extendWhere = extendWhere + "AND game.issue=? "
	}
	// 统计order表
	sql := "SELECT order.game_id,game.issue, LEFT (date, 7) AS months, SUM(amount) AS amount " +
		"FROM `order` " +
		"LEFT JOIN game ON game.game_id = `order`.game_id " +
		extendWhere +
		"GROUP BY game_id, months"
	o := orm.NewOrm()
	var values []orm.Params

	if companyIds != nil {
		_, err = o.Raw(sql, args, companyIds).Values(&values)
	} else {
		_, err = o.Raw(sql, args).Values(&values)
	}

	if err != nil {
		return
	}

	var affected int64 = 0
	// 然后遍历写入pre_verify表
	// game_id,company_id,date(month),amount
	inSql := "INSERT INTO order_pre_verify_cp (game_id,company_id,`date`,`amount`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE amount=VALUES(amount)"
	p, _ := o.Raw(inSql).Prepare()
	for _, value := range values {
		gameId := value["game_id"]
		companyId := value["issue"]
		month := value["months"]
		amount := value["amount"]

		r, e := p.Exec(gameId, companyId, month, amount)
		if e != nil {
			err = e
			return
		}
		aff, _ := r.RowsAffected()
		affected = affected + aff
	}

	affCount = affected
	return
}

// 根据我方主体获取没有对账的发行商
func GetNotVerifyCp(bodyMy int) (companies []CompanyType, err error) {
	// 由于有些游戏没有录入发行商,所以这里可能查出null
	nowMonth := time.Now().Format("2006-01")
	sql := "SELECT `order_pre_verify_cp`.company_id " +
		"FROM `order_pre_verify_cp` " +
		"LEFT JOIN contract on contract.game_id = `order_pre_verify_cp`.game_id AND contract.company_type = 0 " +
		"WHERE verify_id = 0 AND contract.body_my = ? AND date!='" + nowMonth + "'" +
		"GROUP BY company_id"
	var rs []orm.Params
	o := orm.NewOrm()
	_, err = o.Raw(sql, bodyMy).Values(&rs)
	if err != nil {
		return
	}

	// 获取所有发行商
	var chans []CompanyType
	_, err = o.QueryTable(new(CompanyType)).All(&chans, "Id", "Name")
	if err != nil {
		return
	}
	chanMap := make(map[int]CompanyType)
	for _, v := range chans {
		chanMap[v.Id] = v
	}

	companies = []CompanyType{}
	for _, r := range rs {
		companyId, _ := util.Interface2Int(r["company_id"], false)
		if companyId != 0 {
			if x, ok := chanMap[int(companyId)]; ok {
				companies = append(companies, x)
			}
		}
	}

	return
}

// 根据 我方主体 和 发行商 获取没有对账的月份
func GetNotVerifyCpTime(bodyMy int, companyId int) (months []string, err error) {
	o := orm.NewOrm()
	nowMonth := time.Now().Format("2006-01")
	sql := "SELECT `order_pre_verify_cp`.date " +
		"FROM `order_pre_verify_cp` " +
		"LEFT JOIN contract on contract.game_id = `order_pre_verify_cp`.game_id AND contract.company_type = 0 " +
		"WHERE verify_id = 0 AND contract.body_my = ? AND `order_pre_verify_cp`.company_id = ? AND date!='" + nowMonth + "'" +
		"GROUP BY date"
	var rs []orm.Params
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

// 根据 我方主体 和 发行商 和 月份  获取没有对账的游戏账单
func GetNotVerifyCpGame(bodyMy int, companyId int, month string) (gameAmounts []OrderPreVerifyCp, err error) {
	o := orm.NewOrm()
	sql := "SELECT order_pre_verify_cp.id,game_all.game_id,amount,amount_opposite,amount_payable,game_name " +
		"FROM `order_pre_verify_cp` " +
		"LEFT JOIN contract on contract.game_id = `order_pre_verify_cp`.game_id AND contract.company_type = 0 " +
		"LEFT JOIN game_all on game_all.game_id = order_pre_verify_cp.game_id " +
		"WHERE verify_id = 0 AND contract.body_my = ? AND `order_pre_verify_cp`.company_id = ? AND date= ? " +
		"ORDER BY amount DESC"
	var rs []orm.Params
	_, err = o.Raw(sql, bodyMy, companyId, month).Values(&rs)
	if err != nil {
		return
	}

	var allGameIds []string
	gameAmounts = []OrderPreVerifyCp{}
	for _, i := range rs {
		id, _ := util.Interface2Int(i["id"], false)
		gameId, _ := util.Interface2Int(i["game_id"], false)
		gameIdStr, _ := util.Interface2String(i["game_id"], false)
		allGameIds = append(allGameIds, gameIdStr)

		gameName, _ := util.Interface2String(i["game_name"], false)
		Amount, _ := util.Interface2Float(i["amount"], false)
		amountOpposite, _ := util.Interface2Float(i["amount_opposite"], false)
		amountPayable, _ := util.Interface2Float(i["amount_payable"], false)
		gameAmounts = append(gameAmounts, OrderPreVerifyCp{
			Amount:         Amount,
			AmountTheory:   0.01,
			AmountOpposite: amountOpposite,
			AmountPayable:  amountPayable,
			Id:             int(id),
			GameName:       gameName,
			GameId:         int(gameId),
		})
	}

	//  获取理论金额
	/*  2017-12-20 修改，之前的是从老系统获取过来，由于某种问题，获取的数据差异较大，所以直接根据分成比例来计算
	clears, _ := old_sys.GetAllClearing(strings.Join(allGameIds, ","), "", month)
	clearsMap := make(map[string]old_sys.Clearings)
	for _, clear := range *clears {
		clearsMap[clear.GameId] = clear
	}

	for i, v := range gameAmounts {
		gameIdStr := strconv.Itoa(v.GameId)
		gameAmounts[i].AmountTheory = clearsMap[gameIdStr].DivideTotal
	}
	*/

	// ------------------新写的计算理论金额--------------
	for i, v := range gameAmounts {
		gameAmounts[i].AmountTheory = computeCpTheoryMoney(v.GameId, gameAmounts[i].Amount)
	}
	// -----------------------------------------------

	return
}

// 根据分成比例计算理论金额
// @params	gameId	游戏id
// @params	amount	该游戏我方流水
func computeCpTheoryMoney(gameId int, amount float64) float64 {
	o := orm.NewOrm()
	var cpCon Contract
	var cpLadders []old_sys.Ladder4Post
	var cpRatio float64       // cp的我方比例
	var cpSlottingFee float64 // cp的通道费

	o.QueryTable("contract").Filter("game_id__exact", gameId).Filter("company_type__exact", 0).
		Filter("effective_state__exact", 1).One(&cpCon)
	json.Unmarshal([]byte(cpCon.Ladder), &cpLadders)

	// 根据cp分成获取大的我方分成比例
	if len(cpLadders) == 0 {
		return 0
	} else if len(cpLadders) == 1 {
		cpRatio = cpLadders[0].Ratio
		cpSlottingFee = cpLadders[0].SlottingFee
	} else {
		cpRatio = cpLadders[0].Ratio
		cpSlottingFee = cpLadders[0].SlottingFee
		for _, ladder := range cpLadders {
			if cpRatio < ladder.Ratio {
				cpRatio = ladder.Ratio
				cpSlottingFee = ladder.SlottingFee
			}
		}
	}

	return structUtil.Round(amount*(1-cpRatio)*(1-cpSlottingFee), 2)
}

// 获取渠道没有对账的信息
func GetCpNotVerifyInfo(limit, offset int) (count int64, notVerify []NoCpVerify, err error) {
	o := orm.NewOrm()
	// date,company_id,body_my,amount
	sql := "SELECT p.date, g.issue, b.body_my, SUM(p.amount) AS amount FROM `order_pre_verify_cp` AS p LEFT  " +
		"JOIN contract b ON p.game_id = b.game_id AND b.company_type = 0 LEFT JOIN game g ON p.game_id=g.game_id " +
		"WHERE verify_id = 0 AND g.issue IS NOT NULL " +
		"GROUP BY p.company_id, b.body_my, p.date  " +
		"ORDER BY p.date DESC, g.issue LIMIT ?,? "
	var rs []orm.Params
	_, err = o.Raw(sql, offset, limit).Values(&rs)
	if err != nil {
		return
	}

	sqlCount := "SELECT COUNT(*) AS count FROM (SELECT p.date FROM `order_pre_verify_cp` AS p LEFT  " +
		"JOIN contract b ON p.game_id = b.game_id AND b.company_type = 0 LEFT JOIN game g ON p.game_id=g.game_id " +
		"WHERE verify_id = 0  AND g.issue IS NOT NULL " +
		"GROUP BY p.company_id, b.body_my, p.date  " +
		"ORDER BY p.date DESC, g.issue) AS c "
	var countPs []orm.Params
	_, err = o.Raw(sqlCount).Values(&countPs)
	if err != nil {
		return
	}
	count, _ = util.Interface2Int(countPs[0]["count"], false)

	// 获取公司
	company := make([]*CompanyType, len(rs))
	for i, r := range rs {
		companyId, _ := util.Interface2Int(r["issue"], false)
		if companyId != 0 {
			channelName, _ := GetCompanyNameByCompanyId(int(companyId))
			company[i] = &CompanyType{Name: channelName, Id: int(companyId)}
		} else {
			company[i] = &CompanyType{Name: "没选择发行商", Id: 0}
		}
	}

	notVerify = []NoCpVerify{}
	for i, r := range rs {
		bodyMy, _ := util.Interface2Int(r["body_my"], false)
		date, _ := util.Interface2String(r["date"], false)
		Amount, _ := util.Interface2Float(r["amount"], false)

		notVerify = append(notVerify, NoCpVerify{
			Date:    date,
			BodyMy:  int(bodyMy),
			Company: company[i],
			Amount:  Amount,
		})
	}

	return
}
