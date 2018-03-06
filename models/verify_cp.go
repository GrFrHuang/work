package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bygo/util"
	"time"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/utils"
)

type VerifyCp struct {
	Id             int     `orm:"column(id);auto" json:"id,omitempty"`
	Date           string  `orm:"column(date);size(20);null" json:"date,omitempty"`
	BodyMy         int     `orm:"column(body_my);null" json:"body_my,omitempty"`
	CompanyId      int  `orm:"column(company_id);null" json:"company_id,omitempty"`
	AmountMy       float64 `orm:"column(amount_my);null;digits(16);decimals(2)" json:"amount_my,omitempty"`
	AmountOpposite float64 `orm:"column(amount_opposite);null;digits(16);decimals(2)" json:"amount_opposite,omitempty"`
	AmountPayable  float64 `orm:"column(amount_payable);null;digits(16);decimals(2)" json:"amount_payable,omitempty"`
	AmountTheory   float64 `orm:"column(amount_theory);null;digits(16);decimals(2)" json:"amount_theory,omitempty"`
	AmountRemit    float64 `orm:"column(amount_remit);null;digits(16);decimals(2)" json:"amount_remit,omitempty"`
	Status         int    `orm:"column(status);null" json:"status,omitempty" json:"status,omitempty"`
	BillingDate    string `orm:"column(billing_date);null" json:"billing_date,omitempty"`
	RemitCompanyId int    `orm:"column(remit_company_id);null" json:"remit_company_id,omitempty"`
	VerifyTime     int    `orm:"column(verify_time);null" json:"verify_time,omitempty"`
	VerifyUserId   int    `orm:"column(verify_user_id);null" json:"verify_user_id,omitempty"`
	FileId         int    `orm:"column(file_id);null" json:"file_id,omitempty"`
	FilePreviewId  int    `orm:"column(file_preview_id);null" json:"file_preview_id,omitempty"`
	Desc           string `orm:"column(desc);null" json:"desc,omitempty"`
	CreatedTime    int    `orm:"column(created_time);null" json:"created_time,omitempty"`
	CreatedUserId  int    `orm:"column(created_user_id);null" json:"created_user_id,omitempty"`
	UpdatedTime    int    `orm:"column(updated_time);null" json:"updated_time,omitempty"`
	UpdatedUserId  int    `orm:"column(updated_user_id);null" json:"updated_user_id,omitempty"`

	VerifyUser     *User `orm:"-" json:"verify_user,omitempty"`
	UpdatedUser    *User `orm:"-" json:"updated_user,omitempty"`
	Company        *CompanyType `orm:"-" json:"company,omitempty"` // 发行商
	PreVerifyGames []OrderPreVerifyCp `orm:"-" json:"pre_verify,omitempty"`
}

func (t *VerifyCp) TableName() string {
	return "verify_cp"
}

func init() {
	orm.RegisterModel(new(VerifyCp))
}

// 添加渠道对账单
func AddVerifyCp(m *VerifyCp) (id int64, err error) {
	if m.PreVerifyGames == nil || len(m.PreVerifyGames) == 0 {
		err = errors.New("游戏不能未空")
		return
	}
	m.CreatedTime = int(time.Now().Unix())
	m.UpdatedTime = int(time.Now().Unix())

	err = CheckPreVerifyCpCanAdd(m.PreVerifyGames)
	if err != nil {
		return
	}

	o := orm.NewOrm()
	id, err = o.Insert(m)
	if err != nil {
		return
	}
	err = UpdatePreVerifyCpMarkVerifyId(m.PreVerifyGames, int(id))
	if err != nil {
		o.Delete(m)
		return
	}

	// 操作日志
	CompareAndAddOperateLog(new(VerifyCp),m, m.CreatedUserId, bean.OPP_VERIFY_CP, m.Id, bean.OPA_INSERT)

	return
}

// 获取一个对账单
func GetVerifyCpById(id int) (v *VerifyCp, err error) {
	o := orm.NewOrm()
	v = &VerifyCp{Id: id}
	if err = o.Read(v); err != nil {
		return
	}
	verifyId := v.Id
	// 添加游戏单
	v.PreVerifyGames, err = GetPreVerifyCpWithGameByVerifyId(verifyId, true)
	if err != nil {
		return
	}
	// 添加渠道
	v.Company, err = GetCompanyById(v.CompanyId)
	if err != nil {
		return
	}
	return
}

// 简单统计昨日今日数据
func GetCpSimpleStatistics() (data map[string]interface{}, err error) {
	o := orm.NewOrm()
	todayTime, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		return
	}
	today, err := o.QueryTable("verify_cp").
		Filter("CreatedTime__gt", todayTime.Unix()).
		Count()
	if err != nil {
		return
	}
	yesterdayTime, err := time.Parse("2006-01-02", time.Now().Add(time.Second * -24 * 3600).Format("2006-01-02"))
	if err != nil {
		return
	}
	yesterday, err := o.QueryTable("verify_cp").
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

func GetPreVerifyCpWithGameByVerifyId(verifyId int, withTheory bool) (gameAmounts []OrderPreVerifyCp, err error) {
	o := orm.NewOrm()
	sql := "SELECT " +
		"order_pre_verify_cp.id, game_all.game_id, date,amount, game_name, amount_payable, amount_opposite " +
		"FROM `order_pre_verify_cp` " +
		"LEFT JOIN game_all on game_all.game_id = order_pre_verify_cp.game_id " +
		"WHERE verify_id = ? " +
		"ORDER BY amount desc"
	var rs []orm.Params
	_, err = o.Raw(sql, verifyId).Values(&rs)
	if err != nil {
		return
	}

	var allGameIds []string
	//month := ""
	gameAmounts = []OrderPreVerifyCp{}
	for _, i := range rs {
		id, _ := util.Interface2Int(i["id"], false)
		gameId, _ := util.Interface2Int(i["game_id"], false)
		gameIdStr, _ := util.Interface2String(i["game_id"], false)
		allGameIds = append(allGameIds, gameIdStr)
		gameName, _ := util.Interface2String(i["game_name"], false)
		Amount, _ := util.Interface2Float(i["amount"], false)
		AmountPayable, _ := util.Interface2Float(i["amount_payable"], false)
		AmountOpposite, _ := util.Interface2Float(i["amount_opposite"], false)
		gameAmounts = append(gameAmounts, OrderPreVerifyCp{
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
	}

	return
}

func GetPreVerifyCp(verifyId int) (result []OrderPreVerifyCp, err error) {
	o := orm.NewOrm()
	result = []OrderPreVerifyCp{}
	_, err = o.QueryTable("order_pre_verify_cp").
		Filter("verify_id", verifyId).
		All(&result)
	return
}

func AddPreVerifyCpForVerify(data *[]VerifyCp) (err error) {
	if len(*data) == 0 {
		return
	}
	for i, v := range *data {
		(*data)[i].PreVerifyGames, err = GetPreVerifyCpWithGameByVerifyId(v.Id, false)
		if err != nil {
			return
		}
	}
	return
}

func AddCompanyForVerifyCp(data *[]VerifyCp) (err error) {
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

func AddVerifyAndUpdateUserForVerifyCp(data *[]VerifyCp) (err error) {
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
func UpdateVerifyCpById(m *VerifyCp) (err error) {
	o := orm.NewOrm()
	v := VerifyCp{Id: m.Id}
	if err = o.Read(&v); err != nil {
		return
	}

	// 回滚
	err = UpdatePreVerifyCpReset(v.Id)
	if err != nil {
		return
	}

	// 重新标记
	err = UpdatePreVerifyCpMarkVerifyId(m.PreVerifyGames, v.Id)
	if err != nil {
		return
	}

	m.UpdatedTime = int(time.Now().Unix())
	fs := util.NewFieldsUtil(m).GetNotEmptyFields().
		Exclude("PreVerifyGames", "VerifyUser", "UpdatedUser", "Channel").
		Must("AmountMy", "AmountOpposite", "AmountPayable", "AmountTheory", "BillingDate").
		Fields()

	if _, err = o.Update(m, fs...); err != nil {
		return
	}

	// 操作日志
	utils.RemoveFields(fs, "UpdatedTime")
	CompareAndAddOperateLog(&v, m, m.UpdatedUserId, bean.OPP_VERIFY_CP, m.Id, bean.OPA_UPDATE, fs...)

	return
}

func DeleteVerifyCp(id, uid int) (err error) {
	o := orm.NewOrm()
	v := VerifyCp{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err != nil {
		return
	}

	// 回滚
	err = UpdatePreVerifyCpReset(v.Id)
	if err != nil {
		return
	}

	if _, err = o.Delete(&VerifyCp{Id: id}); err != nil {
		return
	}

	// 操作日志
	CompareAndAddOperateLog(&v, new(VerifyCp), uid, bean.OPP_VERIFY_CP, id, bean.OPA_DELETE)

	return
}

//func MigrationVerifyCpAll() (err error) {
//	allCpVerifyAccount := []CpVerifyAccount{}
//	orm.NewOrm().
//		QueryTable("cp_verify_account").
//		All(&allCpVerifyAccount, "id")
//	orm.NewOrm().Raw("DELETE from verify_cp").Exec()
//	for _, v := range allCpVerifyAccount {
//		err := MigrationVerifyCp(v.Id)
//		if err != nil {
//			log.Error("err", err)
//		}
//	}
//
//	return
//}
//
//func MigrationVerifyCp(oldVerifyId int) (err error) {
//	v, err := GetCpVerifyAccountById(oldVerifyId)
//	if err != nil {
//		return
//	}
//	oldGames := []GameAmount{}
//	err = json.Unmarshal([]byte(v.Games), &oldGames)
//	if err != nil {
//		return
//	}
//
//	companyId := v.CompanyId
//	date := v.StartTime[:7]
//	status := v.Status
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
//	preVerifys := []OrderPreVerifyCp{}
//	//log.Info("debug", oldGames)
//	for _, r := range oldGames {
//		// 添加amount数据
//		ov := OrderPreVerifyCp{
//			AmountPayable:  r.AmountPayable,
//			Amount:         r.AmountMy,
//			AmountOpposite: r.AmountOpposite,
//			AmountTheory:   r.AmountTheory,
//		}
//		// 根据 游戏 和 cp 和 时间 和 主体 找到预对账单
//		preVerify := OrderPreVerifyCp{}
//		err = orm.NewOrm().QueryTable("order_pre_verify_cp").
//			Filter("company_id", companyId).
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
//	newV := VerifyCp{
//		Date:           date,
//		Desc:           desc,
//		BodyMy:         bodyMy,
//		FilePreviewId:  filePreviewId,
//		FileId:         fileId,
//		CompanyId:      companyId,
//		CreatedTime:    createTime,
//		CreatedUserId:  createUserId,
//		UpdatedUserId:  updateUserId,
//		UpdatedTime:    updateTime,
//		PreVerifyGames: preVerifys,
//		VerifyTime:     verifyTime,
//		VerifyUserId:   verifyUserId,
//		Status:         status,
//
//		AmountPayable:  amountPayable,
//		AmountMy:       amountMy,
//		AmountOpposite: amountOpposite,
//		AmountTheory:   amountTheory,
//	}
//	if _, err = AddVerifyCp(&newV); err != nil {
//		return
//	}
//
//	return
//}

//获取结算部人员
func GetCpVerifyUsers() (maps []orm.Params, err error) {
	o := orm.NewOrm()
	//var maps []orm.Params
	_, err = o.QueryTable(new(User)).Filter("DepartmentId", "6").Values(&maps, "Nickname", "Id", "Name")
	if err != nil {
		return
	}
	return
}
