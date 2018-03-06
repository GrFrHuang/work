package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"strings"
	"time"
	"kuaifa.com/kuaifa/work-together/tool/old_sys"
	"encoding/json"
	"github.com/Zeniubius/golang_utils/structUtil"
)

type OrderPreVerifyChannel struct {
	Id             int     `orm:"column(id);auto" json:"id,omitempty"`
	GameId         int     `orm:"column(game_id);null" json:"game_id,omitempty"`
	ChannelCode    string  `orm:"column(channel_code);size(50);null" json:"channel_code,omitempty"`
	Date           string  `orm:"column(date);size(20);null" json:"date,omitempty"`
	Amount         float64 `orm:"column(amount);null" json:"amount_my"`
	AmountTheory   float64 `orm:"column(amount_theory);null" json:"amount_theory"`
	AmountOpposite float64 `orm:"column(amount_opposite);null" json:"amount_opposite"`
	AmountPayable  float64 `orm:"column(amount_payable);null" json:"amount_payable"`
	VerifyId       int64   `orm:"column(verify_id);null" json:"verify_id,omitempty"`

	GameName string `orm:"-" json:"game_name,omitempty"`
}

func (t *OrderPreVerifyChannel) TableName() string {
	return "order_pre_verify_channel"
}

func init() {
	orm.RegisterModel(new(OrderPreVerifyChannel))
}

// AddOrderPreVerifyChannel insert a new OrderPreVerifyChannel into database and returns
// last inserted Id on success.
func AddOrderPreVerifyChannel(m *OrderPreVerifyChannel) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetOrderPreVerifyChannelById retrieves OrderPreVerifyChannel by Id. Returns error if
// Id doesn't exist
func GetOrderPreVerifyChannelById(id int) (v *OrderPreVerifyChannel, err error) {
	o := orm.NewOrm()
	v = &OrderPreVerifyChannel{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateOrderPreVerifyChannel updates OrderPreVerifyChannel by Id and returns error if
// the record to be updated doesn't exist
func UpdateOrderPreVerifyChannelById(m *OrderPreVerifyChannel) (err error) {
	o := orm.NewOrm()
	v := OrderPreVerifyChannel{Id: m.Id}
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
func CheckPreVerifyChannelCanAdd(preVerify []OrderPreVerifyChannel) (err error) {
	preVerifyIds := make([]interface{}, len(preVerify))
	for i, v := range preVerify {
		preVerifyIds[i] = v.Id
	}

	var ps []OrderPreVerifyChannel
	_, err = orm.NewOrm().
		QueryTable("order_pre_verify_channel").
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
func UpdatePreVerifyChannelMarkVerifyId(preVerify []OrderPreVerifyChannel, verifyId int) (err error) {
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
func UpdatePreVerifyChannelReset(verifyId int) (err error) {
	o := orm.NewOrm()
	ps := orm.Params{
		"VerifyId": 0,
	}
	_, err = o.QueryTable("order_pre_verify_channel").Filter("verify_id", verifyId).
		Update(ps)
	return
}

// 同步order表
func UpdatePreVerifyChannelFromOrder(channelCodes []interface{}, month string) (affCount int64, err error) {
	extendWhere := ""
	var args []interface{}
	if len(channelCodes) != 0 {
		holder := strings.Repeat(",?", len(channelCodes))
		if extendWhere != "" {
			extendWhere = extendWhere + "AND "
		}
		extendWhere += "cp in (" + holder[1:] + ") "
		args = append(args, channelCodes...)
	}
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

	// 统计order表
	sql := "SELECT order.game_id, cp, LEFT (date, 7) AS months, SUM(amount) AS amount " +
		"FROM `order` " +
		extendWhere +
		"GROUP BY game_id, cp, months"
	o := orm.NewOrm()
	var values []orm.Params
	_, err = o.Raw(sql, args).Values(&values)
	if err != nil {
		return
	}

	var affected int64 = 0
	// 然后遍历写入pre_verify表
	// game_id,channel_code(cp),date(month),amount
	inSql := "INSERT INTO order_pre_verify_channel (game_id,channel_code,`date`,`amount`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE amount=VALUES(amount) "
	p, _ := o.Raw(inSql).Prepare()
	for _, value := range values {
		gameId := value["game_id"]
		cp := value["cp"]
		month := value["months"]
		amount := value["amount"]

		r, e := p.Exec(gameId, cp, month, amount)
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

// 根据我方主体获取没有对账的渠道
func GetNotVerifyChannel(bodyMy int) (channels []Channel, err error) {
	nowMonth := time.Now().Format("2006-01")
	sql := "SELECT order_pre_verify_channel.channel_code " +
		"FROM `order_pre_verify_channel` " +
		"LEFT JOIN contract ON contract.channel_code = order_pre_verify_channel.channel_code AND `order_pre_verify_channel`.game_id = `contract`.game_id AND contract.company_type = 1 " +
		"WHERE verify_id = 0 AND contract.body_my = ? AND date!='" + nowMonth + "' " +
		"GROUP BY order_pre_verify_channel.channel_code"
	var rs []orm.Params
	o := orm.NewOrm()
	_, err = o.Raw(sql, bodyMy).Values(&rs)
	if err != nil {
		return
	}

	// 获取渠道
	channels = make([]Channel, len(rs))
	for i, r := range rs {
		channelCode, _ := r["channel_code"].(string)
		if channelCode != "" {
			channelName, _ := GetChannelNameByChannelCode(channelCode)
			channels[i] = Channel{Name: channelName, Cp: channelCode}
		}
	}

	return
}

// 根据 我方主体 和 渠道 获取没有对账的月份
func GetNotVerifyChannelTime(bodyMy int, channelCode string) (months []string, err error) {
	o := orm.NewOrm()
	nowMonth := time.Now().Format("2006-01")
	sql := "SELECT date " +
		"FROM `order_pre_verify_channel` " +
		"LEFT JOIN contract ON contract.channel_code = `order_pre_verify_channel`.channel_code AND `order_pre_verify_channel`.game_id = `contract`.game_id AND contract.company_type = 1 " +
		"WHERE verify_id = 0 AND contract.body_my = ? AND order_pre_verify_channel.channel_code = ? AND date!='" + nowMonth + "'" +
		"GROUP BY date"
	var rs []orm.Params
	_, err = o.Raw(sql, bodyMy, channelCode).Values(&rs)
	if err != nil {
		return
	}

	months = []string{}
	for _, i := range rs {
		months = append(months, i["date"].(string))
	}

	return
}

// 根据 我方主体 和 渠道 和 月份  获取没有对账的游戏账单
func GetNotVerifyChannelGame(bodyMy int, channelCode string, month string) (gameAmounts []OrderPreVerifyChannel, err error) {
	o := orm.NewOrm()
	sql := "SELECT order_pre_verify_channel.id,game_all.game_id,amount,amount_opposite,amount_payable,game_name " +
		"FROM `order_pre_verify_channel` " +
		"LEFT JOIN contract ON contract.channel_code = `order_pre_verify_channel`.channel_code AND `order_pre_verify_channel`.game_id = `contract`.game_id AND contract.company_type = 1 " +
		"LEFT JOIN game_all on game_all.game_id = order_pre_verify_channel.game_id " +
		"WHERE verify_id = 0 AND contract.body_my = ? AND contract.effective_state=1 AND order_pre_verify_channel.channel_code = ? AND date = ? " +
		"ORDER BY amount DESC"
	var rs []orm.Params
	_, err = o.Raw(sql, bodyMy, channelCode, month).Values(&rs)
	if err != nil {
		return
	}

	var allGameIds []string
	gameAmounts = []OrderPreVerifyChannel{}
	for _, i := range rs {
		id, _ := util.Interface2Int(i["id"], false)
		gameId, _ := util.Interface2Int(i["game_id"], false)
		gameIdStr, _ := util.Interface2String(i["game_id"], false)
		gameName, _ := util.Interface2String(i["game_name"], false)

		allGameIds = append(allGameIds, gameIdStr)

		Amount, _ := util.Interface2Float(i["amount"], false)
		gameAmounts = append(gameAmounts, OrderPreVerifyChannel{
			GameId:       int(gameId),
			Amount:       Amount,
			AmountTheory: 0.01,
			Id:           int(id),
			GameName:     gameName,
		})
	}

	//  获取理论金额
	/*  2017-12-20 修改，之前的是从老系统获取过来，由于某种问题，获取的数据差异较大，所以直接根据分成比例来计算
	clears, _ := old_sys.GetAllClearing(strings.Join(allGameIds, ","), channelCode, month)
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
		gameAmounts[i].AmountTheory = computeChannelTheoryMoney(v.GameId, gameAmounts[i].Amount, channelCode)
	}
	// -----------------------------------------------

	return
}

// 根据分成比例计算理论金额
// @params	gameId	游戏id
// @params	amount	该游戏我方流水
// @params	channelCode	渠道code
func computeChannelTheoryMoney(gameId int, amount float64, channelCode string) float64 {
	o := orm.NewOrm()
	var channelCon Contract
	var channelLadders []old_sys.Ladder4Post
	var channelRatio float64       // 渠道的我方比例
	var channelSlottingFee float64 // 渠道的通道费

	o.QueryTable("contract").Filter("game_id__exact", gameId).Filter("company_type__exact", 1).
		Filter("channel_code__exact", channelCode).Filter("effective_state__exact", 1).One(&channelCon)
	json.Unmarshal([]byte(channelCon.Ladder), &channelLadders)

	// 根据渠道分成比例获取大的我方分成比例
	if len(channelLadders) == 0 {
		return 0
	} else if len(channelLadders) == 1 {
		channelRatio = channelLadders[0].Ratio
		channelSlottingFee = channelLadders[0].SlottingFee
	} else {
		channelRatio = channelLadders[0].Ratio
		channelSlottingFee = channelLadders[0].SlottingFee
		for _, ladder := range channelLadders {
			if channelRatio < ladder.Ratio {
				channelRatio = ladder.Ratio
				channelSlottingFee = ladder.SlottingFee
			}
		}
	}

	return structUtil.Round(amount*channelRatio*(1-channelSlottingFee), 2)
}

// 获取渠道没有对账的信息
// date,channel_code,body_my,amount
func GetChannelNotVerifyInfo(limit, offset int) (count int64, notVerify []NoChannelVerify, err error) {
	o := orm.NewOrm()
	sql := "SELECT p.date, p.channel_code, contract.body_my, sum(p.amount) AS amount FROM `order_pre_verify_channel` AS p " +
		"LEFT JOIN contract ON contract.channel_code = p.channel_code AND p.game_id = `contract`.game_id AND contract.company_type = 1 " +
		"WHERE verify_id = 0 AND contract.id IS NOT NULL GROUP BY p.channel_code, contract.body_my, p.date " +
		"ORDER BY p.date DESC, p.channel_code LIMIT ?,?"
	var rs []orm.Params
	_, err = o.Raw(sql, offset, limit).Values(&rs)
	if err != nil {
		return
	}

	sqlCount := "SELECT count(*) as count FROM (SELECT p.date FROM `order_pre_verify_channel` AS p " +
		"LEFT JOIN contract ON contract.channel_code = p.channel_code AND p.game_id = `contract`.game_id AND contract.company_type = 1 " +
		"WHERE verify_id = 0 AND contract.id IS NOT NULL GROUP BY p.channel_code, contract.body_my, p.date " +
		"ORDER BY p.date DESC, p.channel_code) as c"
	var countPs []orm.Params
	_, err = o.Raw(sqlCount).Values(&countPs)
	if err != nil {
		return
	}
	count, _ = util.Interface2Int(countPs[0]["count"], false)

	// 获取渠道
	channels := make([]*Channel, len(rs))
	for i, r := range rs {
		channelCode, _ := r["channel_code"].(string)
		if channelCode != "" {
			channelName, _ := GetChannelNameByChannelCode(channelCode)
			channels[i] = &Channel{Name: channelName, Cp: channelCode}
		}
	}

	notVerify = []NoChannelVerify{}
	for i, r := range rs {
		bodyMy, _ := util.Interface2Int(r["body_my"], false)
		date, _ := util.Interface2String(r["date"], false)
		Amount, _ := util.Interface2Float(r["amount"], false)

		notVerify = append(notVerify, NoChannelVerify{
			Date:    date,
			BodyMy:  int(bodyMy),
			Channel: channels[i],
			Amount:  Amount,
		})
	}

	return
}
