package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/log"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"time"
	"kuaifa.com/kuaifa/work-together/utils"
)

type VerifyChannel struct {
	Id             int     `orm:"column(id);auto" json:"id,omitempty"`
	Date           string  `orm:"column(date);size(20);null" json:"date,omitempty"`
	BodyMy         int     `orm:"column(body_my);null" json:"body_my,omitempty"`
	ChannelCode    string  `orm:"column(channel_code);size(50);null" json:"channel_code,omitempty"`
	AmountMy       float64 `orm:"column(amount_my);null;digits(16);decimals(2)" json:"amount_my,omitempty"`
	AmountOpposite float64 `orm:"column(amount_opposite);null;digits(16);decimals(2)" json:"amount_opposite,omitempty"`
	AmountPayable  float64 `orm:"column(amount_payable);null;digits(16);decimals(2)" json:"amount_payable,omitempty"`
	AmountTheory   float64 `orm:"column(amount_theory);null;digits(16);decimals(2)" json:"amount_theory,omitempty"`
	AmountRemit    float64 `orm:"column(amount_remit);null;digits(16);decimals(2)" json:"amount_remit,omitempty"`
	Status         int     `orm:"column(status);null" json:"status,omitempty" json:"status,omitempty"`
	BillingDate    string  `orm:"column(billing_date);null" json:"billing_date,omitempty"`
	RemitCompanyId int     `orm:"column(remit_company_id);null" json:"remit_company_id,omitempty"`
	VerifyTime     int     `orm:"column(verify_time);null" json:"verify_time,omitempty"`
	VerifyUserId   int     `orm:"column(verify_user_id);null" json:"verify_user_id,omitempty"`
	FileId         int     `orm:"column(file_id);null" json:"file_id,omitempty"`
	FilePreviewId  int     `orm:"column(file_preview_id);null" json:"file_preview_id,omitempty"`
	Desc           string  `orm:"column(desc);null" json:"desc,omitempty"`
	CreatedTime    int     `orm:"column(created_time);null" json:"created_time,omitempty"`
	CreatedUserId  int     `orm:"column(created_user_id);null" json:"created_user_id,omitempty"`
	UpdatedTime    int     `orm:"column(updated_time);null" json:"updated_time,omitempty"`
	UpdatedUserId  int     `orm:"column(updated_user_id);null" json:"updated_user_id,omitempty"`

	VerifyUser     *User                   `orm:"-" json:"verify_user,omitempty"`
	UpdatedUser    *User                   `orm:"-" json:"updated_user,omitempty"`
	Channel        *Channel                `orm:"-" json:"channel,omitempty"` // 渠道
	PreVerifyGames []OrderPreVerifyChannel `orm:"-" json:"pre_verify,omitempty"`
}

type LyBillGame struct {
	AnysdkgameId	int 	`json:"anysdkgame_id"`
	GameName		string 	`json:"game_name"`
	Date 			string 	`json:"date"`
	Amount			float64 `json:"amount"`
	AmountPayable	float64 `json:"amount_payable"`
}
type LyUserBill struct {
	Id                 int
	UserId             int
	ChannelPlatform    string
	BodyMy             int
	StartTime          time.Time
	EndTime            time.Time
	Games              string
	TotalAmount        float64
	TotalAmountPayable float64
	TicketUrl          string
	Status             int8
	CreateTime         int
	UpdateTime         int
	WtId			   int
	Cmd			   	   int
}

func (t *VerifyChannel) TableName() string {
	return "verify_channel"
}

func init() {
	orm.RegisterModel(new(VerifyChannel))
}

// 添加渠道对账单
func AddVerifyChannel(m *VerifyChannel) (id int64, err error) {
	if m.PreVerifyGames == nil || len(m.PreVerifyGames) == 0 {
		err = errors.New("游戏不能未空")
		return
	}
	m.CreatedTime = int(time.Now().Unix())
	m.UpdatedTime = int(time.Now().Unix())

	err = CheckPreVerifyChannelCanAdd(m.PreVerifyGames)
	if err != nil {
		return
	}

	o := orm.NewOrm()
	id, err = o.Insert(m)
	if err != nil {
		return
	}
	err = UpdatePreVerifyChannelMarkVerifyId(m.PreVerifyGames, int(id))
	if err != nil {
		o.Delete(m)
		return
	}

	// 操作日志
	CompareAndAddOperateLog(new(VerifyChannel), m, m.CreatedUserId, bean.OPP_VERIFY_CHANNEL, m.Id, bean.OPA_INSERT)

	err = AfterVerifyChannelChanged(VerifyChannel{}, *m)
	if err != nil {
		return
	}

	return
}

// 获取一个对账单
func GetVerifyChannelById(id int) (v *VerifyChannel, err error) {
	o := orm.NewOrm()
	v = &VerifyChannel{Id: id}
	if err = o.Read(v); err != nil {
		return
	}
	verifyId := v.Id
	// 添加游戏单
	v.PreVerifyGames, err = GetPreVerifyChannelWithGameByVerifyId(verifyId, true)
	if err != nil {
		return
	}
	// 添加渠道
	v.Channel, err = GetChannelByChannelCode(v.ChannelCode)
	if err != nil {
		return
	}
	return
}

func GetPreVerifyChannelWithGameByVerifyId(verifyId int, withTheory bool) (gameAmounts []OrderPreVerifyChannel, err error) {
	o := orm.NewOrm()
	sql := "SELECT " +
		"order_pre_verify_channel.id,channel_code, game_all.game_id, amount, game_name, amount_payable, amount_opposite " +
		"FROM `order_pre_verify_channel` " +
		"LEFT JOIN game_all on game_all.game_id = order_pre_verify_channel.game_id " +
		"WHERE verify_id = ? " +
		"ORDER BY amount desc"
	var rs []orm.Params
	_, err = o.Raw(sql, verifyId).Values(&rs)
	if err != nil {
		return
	}

	var allGameIds []string
	gameAmounts = []OrderPreVerifyChannel{}
	//month := ""
	channelCode := ""
	for _, i := range rs {
		id, _ := util.Interface2Int(i["id"], false)
		gameId, _ := util.Interface2Int(i["game_id"], false)
		gameIdStr, _ := util.Interface2String(i["game_id"], false)
		//month, _ = util.Interface2String(i["date"], false)
		channelCode, _ = util.Interface2String(i["channel_code"], false)
		allGameIds = append(allGameIds, gameIdStr)

		gameName, _ := util.Interface2String(i["game_name"], false)
		Amount, _ := util.Interface2Float(i["amount"], false)
		AmountPayable, _ := util.Interface2Float(i["amount_payable"], false)
		AmountOpposite, _ := util.Interface2Float(i["amount_opposite"], false)
		gameAmounts = append(gameAmounts, OrderPreVerifyChannel{
			Amount:         Amount,
			AmountTheory:   0.01,
			AmountOpposite: AmountOpposite,
			AmountPayable:  AmountPayable,
			Id:             int(id),
			GameName:       gameName,
			GameId:         int(gameId),
		})
	}

	if withTheory {
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
	}
	return
}

func GetPreVerifyChannel(verifyId int) (result []OrderPreVerifyChannel, err error) {
	o := orm.NewOrm()
	result = []OrderPreVerifyChannel{}
	_, err = o.QueryTable("order_pre_verify_channel").
		Filter("verify_id", verifyId).
		All(&result)
	return
}

func AddPreVerifyChannelForVerifyList(data *[]VerifyChannel) (err error) {
	if len(*data) == 0 {
		return
	}
	for i, v := range *data {
		(*data)[i].PreVerifyGames, err = GetPreVerifyChannelWithGameByVerifyId(v.Id, false)
		if err != nil {
			return
		}
	}
	return
}

func AddChannelForVerifyList(data *[]VerifyChannel) (err error) {
	if len(*data) == 0 {
		return
	}
	for i, v := range *data {
		(*data)[i].Channel, err = GetChannelByChannelCode(v.ChannelCode)
		if err != nil {
			return
		}
	}
	return
}

func AddVerifyAndUpdateUserForVerifyList(data *[]VerifyChannel) (err error) {
	if len(*data) == 0 {
		return
	}
	var uid []interface{}
	for _, v := range *data {
		uid = append(uid, v.UpdatedUserId, v.VerifyUserId)
	}
	if len(uid) == 0 {
		return
	}
	var users []*User
	orm.NewOrm().QueryTable(new(User)).Filter("id__in", uid).All(&users)
	usersMap := map[int]*User{}
	for i, v := range users {
		usersMap[v.Id] = users[i]
	}

	for i, v := range *data {
		(*data)[i].UpdatedUser = usersMap[v.UpdatedUserId]
		(*data)[i].VerifyUser = usersMap[v.VerifyUserId]
	}

	return
}

// 更新对账单，需要回滚已经添加的对账单
func UpdateVerifyChannelById(m *VerifyChannel) (err error) {
	o := orm.NewOrm()
	v := VerifyChannel{Id: m.Id}
	if err = o.Read(&v); err != nil {
		return
	}

	// 已开票
	//if v.Status == 30 {
	//	if m.Status != 100 {
	//		err = errors.New("已开发票的账单只能修改状态为已回款")
	//		return
	//	}
	//} else if v.Status == 100 {
	//	err = errors.New("已回款的账单不能修改")
	//	return
	//}

	// 2018-01-17修改，去除状态判断，都可以修改游戏清单
	//if v.Status != 30 {
	// 回滚
	err = UpdatePreVerifyChannelReset(v.Id)
	if err != nil {
		return
	}

	// 重新标记
	err = UpdatePreVerifyChannelMarkVerifyId(m.PreVerifyGames, v.Id)
	if err != nil {
		return
	}
	//}

	m.UpdatedTime = int(time.Now().Unix())
	fs := []string{"Status", "UpdatedTime", "UpdatedUserId"}
	if v.Status != 30 {
		fs = util.NewFieldsUtil(m).GetNotEmptyFields().
			Exclude("PreVerifyGames", "VerifyUser", "UpdatedUser", "Channel").
			Must("AmountMy", "AmountOpposite", "AmountPayable", "AmountTheory", "BillingDate").
			Fields()
	}

	if _, err = o.Update(m, fs...); err != nil {
		return
	}

	// 操作日志
	utils.RemoveFields(fs, "UpdatedTime")
	CompareAndAddOperateLog(&v, m, m.UpdatedUserId, bean.OPP_VERIFY_CHANNEL, m.Id, bean.OPA_UPDATE, fs...)

	err = AfterVerifyChannelChanged(v, *m)
	if err != nil {
		return
	}

	return
}

func DeleteVerifyChannel(id, uid int) (err error) {
	o := orm.NewOrm()
	v := VerifyChannel{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err != nil {
		return
	}
	if v.Status >= 30 {
		err = errors.New("已开发票的账单不能再修改")
		return
	}

	// 回滚
	err = UpdatePreVerifyChannelReset(v.Id)
	if err != nil {
		return
	}

	if _, err = o.Delete(&VerifyChannel{Id: id}); err != nil {
		return
	}

	// 操作日志
	CompareAndAddOperateLog(&v, new(VerifyChannel), uid, bean.OPP_VERIFY_CHANNEL, id, bean.OPA_DELETE)

	err = AfterVerifyChannelChanged(v, VerifyChannel{})
	if err != nil {
		return
	}

	return
}

//
//func MigrationVerifyChannelAll() (err error) {
//	// 只在迁移数据时使用 . 危险,请不要执行
//	return
//	allChannelVerifyAccount := []ChannelVerifyAccount{}
//	orm.NewOrm().
//		QueryTable("channel_verify_account").
//		All(&allChannelVerifyAccount, "id")
//
//	orm.NewOrm().Raw("DELETE from verify_channel").Exec()
//	for _, v := range allChannelVerifyAccount {
//		err := MigrationVerifyChannel(v.Id)
//		if err != nil {
//			log.Error("err", err)
//		}
//	}
//
//	return
//}
//
//// 修复 渠道已经对账的预对账单被标记为-1的错误
//func FixVerifyChannelFormOld() {
//	allOldVerify := []ChannelVerifyAccount{}
//	orm.NewOrm().
//		QueryTable("channel_verify_account").
//		All(&allOldVerify, "id")
//
//	allVerify := []VerifyChannel{}
//	orm.NewOrm().
//		QueryTable("verify_channel").
//		All(&allVerify)
//
//	needFix := []VerifyChannel{}
//	for _, v := range allVerify {
//		pres, err := GetPreVerifyChannelWithGameByVerifyId(v.Id, false)
//		if err != nil {
//			//needFix = append(needFix,v)
//			continue
//		}
//		if len(pres) == 0 {
//			needFix = append(needFix, v)
//			continue
//		}
//	}
//
//	log.Info("test-needFix", needFix)
//	//return
//
//	for _, v := range needFix {
//		// 找到老对账单
//		sql := "SELECT * FROM `channel_verify_account` WHERE LEFT(start_time,7)=? AND cp =? AND body_my = ?;"
//		vs := []orm.Params{}
//		l, _ := orm.NewOrm().Raw(sql, v.Date, v.ChannelCode, v.BodyMy).Values(&vs)
//		if l == 0 {
//			log.Error("test", "can not found old verify", v.Date, v.ChannelCode, v.BodyMy)
//			continue
//		}
//
//		oldId, _ := util.Interface2Int(vs[0]["id"], false)
//		if oldId == 0 {
//			log.Error("test", "can not found old verify [id]")
//			continue
//		}
//
//		old := ChannelVerifyAccount{Id: int(oldId)}
//		err := orm.NewOrm().Read(&old)
//		if err != nil {
//			log.Error("test", err)
//			continue
//		}
//
//		// 将已经录入的标记
//		_, err = orm.NewOrm().Update(&ChannelVerifyAccount{Id: int(oldId), Cp: old.Cp + "[1]"}, "Cp")
//		if err != nil {
//			log.Error("test", err)
//			continue
//		}
//
//		// 取出老对账单的游戏
//		oldGames := []GameAmount{}
//		err = json.Unmarshal([]byte(old.GameStr), &oldGames)
//		if err != nil {
//			log.Error("test", err)
//			continue
//		}
//
//		// 更新预对账单
//		gameIds := []string{}
//		for _, g := range oldGames {
//			gameIds = append(gameIds, strconv.Itoa(g.GameId))
//		}
//
//		if len(gameIds) == 0 {
//			log.Error("test", "gameId is null")
//			continue
//		}
//
//		// 根据 游戏 和 cp 和 时间 和 主体 找到预对账单
//		o := orm.NewOrm()
//		sql2 := "SELECT order_pre_verify_channel.id,amount,amount_opposite,amount_payable " +
//			"FROM `order_pre_verify_channel` " +
//			"LEFT JOIN contract ON contract.channel_code = `order_pre_verify_channel`.channel_code AND `order_pre_verify_channel`.game_id = `contract`.game_id AND contract.company_type = 1 " +
//			"WHERE contract.body_my = ? AND order_pre_verify_channel.channel_code = ? AND date = ? AND order_pre_verify_channel.game_id in (" + strings.Join(gameIds, ",") + ")"
//		pres := []orm.Params{}
//		_, err = o.Raw(sql2, v.BodyMy, v.ChannelCode, v.Date).Values(&pres)
//		if err != nil {
//			log.Error("test", err)
//			return
//		}
//		if len(pres) == 0 {
//			continue
//		}
//		presIds := []interface{}{}
//		for _, v := range pres {
//			presIds = append(presIds, v["id"])
//		}
//		updatePs := orm.Params{
//			"verify_id": v.Id,
//		}
//		_, err = o.QueryTable("order_pre_verify_channel").Filter("id__in", presIds).
//			Update(updatePs)
//
//		if err != nil {
//			log.Error("test", err)
//			continue
//		}
//
//	}
//
//}
//
//func MigrationVerifyChannel(oldVerifyId int) (err error) {
//	v, err := GetChannelVerifyAccountById(oldVerifyId)
//	if err != nil {
//		return
//	}
//	oldGames := []GameAmount{}
//	err = json.Unmarshal([]byte(v.GameStr), &oldGames)
//	if err != nil {
//		return
//	}
//
//	cp := v.Cp
//	date := v.StartTime[:7]
//	status := v.Status
//	remitCompanyId := v.RemitCompanyId
//	fileId := v.FileId
//	filePreviewId := v.FilePreviewId
//	createUserId := v.CreateUserId
//	verifyUserId := v.VerifyUserId
//	updateUserId := v.UpdateUserId
//	verifyTime := v.VerifyTime
//	createTime := v.CreateTime
//	updateTime := v.UpdateTime
//	desc := v.Desc
//	bodyMy := v.BodyMy
//	amountTheory := v.AmountTheory
//	amountOpposite := v.AmountOpposite
//	amountMy := v.AmountMy
//	amountPayable := v.AmountPayable
//
//	preVerifys := []OrderPreVerifyChannel{}
//	//log.Info("debug", oldGames)
//	for _, r := range oldGames {
//		// 添加amount数据
//		ov := OrderPreVerifyChannel{
//			AmountPayable:  r.AmountPayable,
//			Amount:         r.AmountMy,
//			AmountOpposite: r.AmountOpposite,
//			AmountTheory:   r.AmountTheory,
//		}
//		// 根据 游戏 和 cp 和 时间 和 主体 找到预对账单
//		preVerify := OrderPreVerifyChannel{}
//		err = orm.NewOrm().QueryTable("order_pre_verify_channel").
//			Filter("ChannelCode", cp).
//			Filter("Date", date).
//			Filter("GameId", r.GameId).
//			Filter("BodyMy", bodyMy).
//			One(&preVerify)
//		if err != nil {
//			return
//		}
//		ov.Id = preVerify.Id
//
//		preVerifys = append(preVerifys, ov)
//	}
//
//	newV := VerifyChannel{
//		Date:           date,
//		Desc:           desc,
//		BodyMy:         bodyMy,
//		FilePreviewId:  filePreviewId,
//		FileId:         fileId,
//		ChannelCode:    cp,
//		CreatedTime:    createTime,
//		CreatedUserId:  createUserId,
//		UpdatedUserId:  updateUserId,
//		UpdatedTime:    updateTime,
//		PreVerifyGames: preVerifys,
//		VerifyTime:     verifyTime,
//		VerifyUserId:   verifyUserId,
//		Status:         status,
//		RemitCompanyId: remitCompanyId,
//
//		AmountPayable:  amountPayable,
//		AmountMy:       amountMy,
//		AmountOpposite: amountOpposite,
//		AmountTheory:   amountTheory,
//	}
//	if _, err = AddVerifyChannel(&newV); err != nil {
//		return
//	}
//
//	return
//}

// 简单统计昨日今日数据
func GetChannelSimpleStatistics() (data map[string]interface{}, err error) {
	o := orm.NewOrm()
	todayTime, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		return
	}
	today, err := o.QueryTable("verify_channel").
		Filter("CreatedTime__gt", todayTime.Unix()).
		Count()
	if err != nil {
		return
	}
	yesterdayTime, err := time.Parse("2006-01-02", time.Now().Add(time.Second * -24 * 3600).Format("2006-01-02"))
	if err != nil {
		return
	}
	yesterday, err := o.QueryTable("verify_channel").
		Filter("CreatedTime__gt", yesterdayTime.Unix()).
		Filter("CreatedTime__lt", todayTime.Unix()).
		Count()
	if err != nil {
		return
	}
	data = map[string]interface{}{
		"today":     today,
		"yesterday": yesterday,
	}
	return
}

func AfterVerifyChannelChanged(old, new VerifyChannel) (err error) {
	// 2017-04-12 不再影响偏移量
	return
	needModifyOffset := false

	if old.Id == 0 || new.Id == 0 {
		if old.Id == 0 { // create
			// 4月之前的
			needModifyOffset = new.Date <= "2017-03"
		} else { // delete
			// 4月之前的
			needModifyOffset = old.Date <= "2017-03"
		}
	} else {
		needModifyOffset = new.Date <= "2017-03"
	}

	// change
	if needModifyOffset {
		var newAmount float64 = 0
		// 判断状态
		// 已开发票才能计入金额
		if new.Status >= 30 {
			newAmount = new.AmountPayable
		}
		var oldAmount float64 = 0
		if old.Status >= 30 {
			oldAmount = old.AmountPayable
		}

		offSet := newAmount - oldAmount
		if offSet == 0 {
			return
		}
		err = AddRemitPreAmountOffset(old.BodyMy, old.RemitCompanyId, 0-offSet)
		log.Info("test", "AddRemitPreAmountOffset", offSet)
		if err != nil {
			return
		}
	}

	return
}

//通过channel_code 获得 对账时间
func GetVerifyTimeByChannelCode(code string) (time int, err error) {
	o := orm.NewOrm()
	verifyChannel := VerifyChannel{}

	qs := o.QueryTable(new(VerifyChannel))
	err = qs.Filter("ChannelCode__exact", code).One(&verifyChannel)

	//if err = o.Read(&verifyChannel,"Id","VerifyTime"); err != nil{
	//	return
	//}
	return verifyChannel.VerifyTime, err

}
