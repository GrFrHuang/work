package task

import (
	"github.com/astaxie/beego/toolbox"
	"fmt"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models"
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/tool/old_sys"
	"math"
	"time"
	"kuaifa.com/kuaifa/work-together/task/warning"
	"strings"
	"strconv"
	"kuaifa.com/kuaifa/work-together/appcache"
	"github.com/bysir-zl/bygo/log"
)

func Run() {
	// 每天4点计算前一天利润
	profit := toolbox.NewTask("ComputeProfit", "0 0 4 * * * ", func() error {
		ComputeProfit(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
		return nil
	})

	// 每天早上9点，下午6点预警游戏停运
	gameOutage := toolbox.NewTask("GameOutage", "0 0 9,18 * * * ", func() error {
		warning.GameOutageWarning()
		return nil
	})

	// 每天早上9:00，下午6:00根据游戏关服时间下架游戏合同
	downContract := toolbox.NewTask("DownContract", "0 0 9,18 * * * ", func() error {
		warning.DownContract()
		return nil
	})

	// 每天早上9:02，下午6:02预警游戏停运渠道合同待处理
	channelContractOpe := toolbox.NewTask("", "0 2 9,18 * * * ", func() error {
		warning.GameOutageChannelContractOpe()
		return nil
	})

	// 每天早上9:02，下午6:02预警游戏停运cp合同待处理
	cpContractOpe := toolbox.NewTask("", "0 2 9,18 * * * ", func() error {
		warning.GameOutageCpContractOpe()
		return nil
	})
	// 每天早上9:03，下午6:03预警游戏停运渠道合同未下架
	channelContractDown := toolbox.NewTask("", "0 3 9,18 * * * ", func() error {
		warning.GameOutageChannelContractDown()
		return nil
	})
	// 每天早上9:03，下午6:03预警游戏停运CP合同未下架
	cpContractDown := toolbox.NewTask("", "0 3 9,18 * * * ", func() error {
		warning.GameOutageCpContractDown()
		return nil
	})

	// 渠道主合同过期预警
	channelMainContractExpireWarning := toolbox.NewTask("", "0 4 9,18 * * * ", func() error {
		warning.ChannelMainContractExpireWarning()
		return nil
	})

	// CP主合同过期预警
	cpMainContractExpireWarning := toolbox.NewTask("", "0 4 9,18 * * * ", func() error {
		warning.CpMainContractExpireWarning()
		return nil
	})

	// 每月2号凌晨2点，生成cp和渠道的预对账单
	UpdatePreVerifyCpFromOrder := toolbox.NewTask("UpdatePreVerifyCpFromOrder", "0 0 2 2 * *", func() error {
		month := time.Now().AddDate(0, -1, 0).Format("2006-01")
		models.UpdatePreVerifyCpFromOrder(nil, month)
		models.UpdatePreVerifyChannelFromOrder(nil, month)
		return nil
	})

	// 每天凌晨2:30,将通讯录中渠道商的商务负责人同步到渠道接入中商务负责人
	SyncBusiness := toolbox.NewTask("SyncChannelBusiness", "0 0 2 * * *", func() error {
		// 第一次同步，需要全表扫描，手动执行
		//FirstSyncChannelBusiness()
		// 第二次以后同步，定时任务执行
		SyncChannelBusiness()

		// 第一次同步，需要全表扫描，手动执行
		//FirstSyncGameAccessPerson()
		// 第二次以后同步，定时任务执行
		SyncGameAccessPerson()

		return nil
	})

	// 每天凌晨2点判断主合同状态是否到期，并修改主合同状态
	ChangeMainContractState := toolbox.NewTask("ChangeMainContractState", "0 0 2 * * *", func() error {
		ChangeMainContractState()
		return nil
	})

	// 没隔1小时，根据流水判断合同并生成
	AddContract := toolbox.NewTask("AddContract", " 0 0 * * * *", func() error {
		AddContractByOrder()
		return nil
	})

	toolbox.AddTask("profit", profit)
	toolbox.AddTask("gameOutage", gameOutage)
	toolbox.AddTask("downContract", downContract)
	toolbox.AddTask("channelContractOpe", channelContractOpe)
	toolbox.AddTask("cpContractOpe", cpContractOpe)
	toolbox.AddTask("channelContractDown", channelContractDown)
	toolbox.AddTask("cpContractDown", cpContractDown)
	toolbox.AddTask("channelMainContractExpireWarning", channelMainContractExpireWarning)
	toolbox.AddTask("cpMainContractExpireWarning", cpMainContractExpireWarning)
	toolbox.AddTask("UpdatePreVerifyCpFromOrder", UpdatePreVerifyCpFromOrder)
	toolbox.AddTask("SyncChannelBusiness", SyncBusiness)
	toolbox.AddTask("ChangeMainContractState", ChangeMainContractState)
	toolbox.AddTask("AddContract", AddContract)

	toolbox.StartTask()
	fmt.Println("task success")
}

// 每天4点定时计算前一天利润
func ComputeProfit(date string) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")
	var pros []models.Profit
	// 用于存储游戏的cp分成
	cpLadders := make(map[int]string)

	qb = qb.Select("game_id, cp as channel_code, date, amount").From("`order`").Where("`date` = ?")
	// 2018-01-03 加,status表示是否通知cp成功，0：成功
	qb = qb.And("status = 0")

	_, err := o.Raw(qb.String(), date).QueryRows(&pros)
	if err != nil {
		return
	}

	for _, pro := range pros {
		var cpCon models.Contract
		var channelCon models.Contract

		// 如果map中存有此游戏的cp阶梯分成，则直接从map中获取，不再查询数据库
		if cpLadders[pro.GameId] != "" {
			cpCon.Ladder = cpLadders[pro.GameId]
		} else {
			o.QueryTable("contract").Filter("game_id__exact", pro.GameId).Filter("company_type__exact", 0).
				Filter("effective_state__exact", 1).One(&cpCon, "game_id", "body_my", "ladders")
			cpLadders[pro.GameId] = cpCon.Ladder
		}
		o.QueryTable("contract").Filter("game_id__exact", pro.GameId).Filter("company_type__exact", 1).
			Filter("channel_code__exact", pro.ChannelCode).Filter("effective_state__exact", 1).
			One(&channelCon, "game_id", "body_my", "channel_code", "ladders")

		profit := compute(pro.Amount, cpCon.Ladder, channelCon.Ladder, pro.GameId, pro.ChannelCode)
		pro.Profit = profit
		pro.BodyMy = channelCon.BodyMy

		o.Insert(&pro)
	}
}

// 计算利润	amount:流水金额	cpLadder:cp阶梯分成	channelLadder:渠道阶梯分成;
// 如果cp或者渠道任一方分成比例没有，则返回利润为0
func compute(amount float64, cpLadder string, channelLadder string, gameId int, channelCode string) float64 {
	var cpLadders []old_sys.Ladder4Post
	var channelLadders []old_sys.Ladder4Post
	json.Unmarshal([]byte(cpLadder), &cpLadders)
	json.Unmarshal([]byte(channelLadder), &channelLadders)

	var cpRatio float64            // cp的我方比例
	var cpSlottingFee float64      // cp的通道费
	var channelRatio float64       // 渠道的我方比例
	var channelSlottingFee float64 // 渠道的通道费

	// 根据cp分成获取小的我方分成比例
	if len(cpLadders) == 0 {
		return 0
	} else {
		//cpRatio, cpSlottingFee = getLadder(cpLadders)
		cpRatio, cpSlottingFee = getBetterLadder(cpLadders, gameId, channelCode)
	}

	// 根据渠道分成比例获取小的我方分成比例
	if len(channelLadders) == 0 {
		return 0
	} else {
		//channelRatio, channelSlottingFee = getLadder(channelLadders)
		channelRatio, channelSlottingFee = getBetterLadder(channelLadders, gameId, channelCode)
	}

	profit := amount*(1-channelSlottingFee)*channelRatio - amount*(1-cpSlottingFee)*(1-cpRatio)
	pow10 := math.Pow10(2)
	return math.Trunc((profit+0.5/pow10)*pow10) / pow10
}

// 进阶计算，考虑了阶梯分成的金额和时间规则
func getBetterLadder(Ladders []old_sys.Ladder4Post, gameId int, channelCode string) (Ratio float64, SlottingFee float64) {
	if len(Ladders) == 1 {
		Ratio = Ladders[0].Ratio
		SlottingFee = Ladders[0].SlottingFee
		return
	} else {
		var rule string
		Ratio = Ladders[0].Ratio
		SlottingFee = Ladders[0].SlottingFee
		for _, ladder := range Ladders {
			effectiveTimeBegin, _ := time.Parse("2006-01-02", ladder.StartTime)
			effectiveTimeEnd, _ := time.Parse("2006-01-02", ladder.EndTime)
			if time.Now().After(effectiveTimeBegin.AddDate(0, 0, -1)) &&
				time.Now().Before(effectiveTimeEnd.AddDate(0, 0, 1)) {

				rule = ladder.Rule
				if rule != "" {
					rulesString := strings.Split(rule, "&")

					// 同时有多个规则
					if len(rulesString) > 1 {
						money := models.GetTotalMoney(gameId, channelCode, ladder.StartTime, ladder.EndTime)
						flag := false
						for _, rul := range rulesString {
							detail := strings.Split(rul, "<")
							begin := detail[0]
							typ := detail[1]
							end := detail[2]

							if typ == "time" {
								beginTime, _ := time.Parse("2006-01-02", begin)
								endTime, _ := time.Parse("2006-01-02", end)

								if time.Now().After(beginTime.AddDate(0, 0, -1)) &&
									time.Now().Before(endTime.AddDate(0, 0, 1)) {
									flag = true
								} else {
									flag = false
									break
								}
							} else if typ == "money" {
								beginMoney, _ := strconv.ParseFloat(begin, 10)
								endMoney, _ := strconv.ParseFloat(end, 10)

								if money >= beginMoney && money < endMoney {
									flag = true
								} else {
									flag = false
									break
								}
							}
						}
						if flag == true {
							Ratio = ladder.Ratio
							SlottingFee = ladder.SlottingFee
							return
						}
					} else if len(rulesString) == 1 {
						// 只有1个规则
						fmt.Printf("rule:%v\nladder:%v\n", rulesString, Ladders)
						detail := strings.Split(rulesString[0], "<")
						begin := detail[0]
						typ := detail[1]
						end := detail[2]
						if typ == "time" {
							beginTime, _ := time.Parse("2006-01-02", begin)
							endTime, _ := time.Parse("2006-01-02", end)

							if time.Now().After(beginTime.AddDate(0, 0, -1)) &&
								time.Now().Before(endTime.AddDate(0, 0, 1)) {
								Ratio = ladder.Ratio
								SlottingFee = ladder.SlottingFee
								return
							}
						} else if typ == "money" {
							money := models.GetTotalMoney(gameId, channelCode, ladder.StartTime, ladder.EndTime)
							beginMoney, _ := strconv.ParseFloat(begin, 10)
							endMoney, _ := strconv.ParseFloat(end, 10)

							if money >= beginMoney && money < endMoney {
								Ratio = ladder.Ratio
								SlottingFee = ladder.SlottingFee
								return
							}
						}
					}
				} else {
					// 规则为空，按生效时间计算
					Ratio = ladder.Ratio
					SlottingFee = ladder.SlottingFee
					return
				}
			}
		}

		// 其他情况未匹配到规则，则按之前的方式
		Ratio, SlottingFee = getLadder(Ladders)
		return
	}
	return
}

// 目前只粗略计算，cp和渠道方的我方比例都按少的算;
func getLadder(Ladders []old_sys.Ladder4Post) (Ratio float64, SlottingFee float64) {
	if len(Ladders) == 1 {
		Ratio = Ladders[0].Ratio
		SlottingFee = Ladders[0].SlottingFee
	} else {
		Ratio = Ladders[0].Ratio
		SlottingFee = Ladders[0].SlottingFee
		for _, ladder := range Ladders {
			if Ratio > ladder.Ratio {
				Ratio = ladder.Ratio
				SlottingFee = ladder.SlottingFee
			}
		}
	}
	return Ratio, SlottingFee
}

type business struct {
	Yunduan  int
	Youliang int
}

// 第一次同步时需要全表扫描
func
FirstSyncChannelBusiness() {
	o := orm.NewOrm()

	var channelAccess []models.ChannelAccess
	channelCompany := models.ChannelCompany{}

	o.QueryTable(new(models.ChannelAccess)).Limit(10000).All(&channelAccess)

	for _, access := range channelAccess {
		cachekey := fmt.Sprintf("Cache_channel_business_%s", access.ChannelCode)
		cacheinfo := appcache.Bm.Get(cachekey)

		// 先获取到该渠道的商务负责人
		channelBusiness := business{}
		if cacheinfo == nil {
			o.QueryTable(new(models.ChannelCompany)).Filter("channel_code", access.ChannelCode).One(&channelCompany)

			channelBusiness.Yunduan = channelCompany.YunduanResponsiblePerson
			channelBusiness.Youliang = channelCompany.YouliangResponsiblePerson

			appcache.Bm.Put(cachekey, channelBusiness, 1*time.Hour)
		} else {
			channelBusiness = cacheinfo.(business)
		}

		// 判断该渠道接入的商务负责人是否需要修改
		if access.BodyMy == 1 {
			if access.BusinessPerson != channelBusiness.Yunduan && channelBusiness.Yunduan != 0 {
				access.BusinessPerson = channelBusiness.Yunduan
				o.Update(&access)
			}
		} else if access.BodyMy == 2 {
			if access.BusinessPerson != channelBusiness.Youliang && channelBusiness.Youliang != 0 {
				access.BusinessPerson = channelBusiness.Youliang
				o.Update(&access)
			}
		}
	}
}

// 第二次以后同步只需要同步渠道商负责人有更改的
func
SyncChannelBusiness() {
	o := orm.NewOrm()

	var channelAccess []models.ChannelAccess
	var channelCompany []models.ChannelCompany

	// 只需要查询在最近一天有更新的渠道商信息
	o.QueryTable(new(models.ChannelCompany)).Filter("update_time__gte", time.Now().Unix()-25*60*60).All(&channelCompany)

	for _, company := range channelCompany {
		// 先获取到该渠道的商务负责人
		o.QueryTable(new(models.ChannelAccess)).Filter("channel_code", company.ChannelCode).All(&channelAccess)

		for _, access := range channelAccess {
			// 判断该渠道接入的商务负责人是否需要修改
			if access.BodyMy == 1 {
				if access.BusinessPerson != company.YunduanResponsiblePerson && company.YunduanResponsiblePerson != 0 {
					access.BusinessPerson = company.YunduanResponsiblePerson
					o.Update(&access)
				}
			} else if access.BodyMy == 2 {
				if access.BusinessPerson != company.YouliangResponsiblePerson && company.YouliangResponsiblePerson != 0 {
					access.BusinessPerson = company.YouliangResponsiblePerson
					o.Update(&access)
				}
			}
		}
	}
}

// 第一次同步时需要全表扫描
func
FirstSyncGameAccessPerson() {
	o := orm.NewOrm()

	var games []models.Game
	distributionCompany := models.DistributionCompany{}

	o.QueryTable(new(models.Game)).Filter("game_id__gte", 0).Limit(10000).All(&games)

	for _, game := range games {
		cachekey := fmt.Sprintf("Cache_cp_business_%s", game.Issue)
		cacheinfo := appcache.Bm.Get(cachekey)

		// 先获取到该发行商的商务负责人
		cpBusiness := business{}
		if cacheinfo == nil {
			o.QueryTable(new(models.DistributionCompany)).Filter("company_id", game.Issue).One(&distributionCompany)

			cpBusiness.Yunduan = distributionCompany.YunduanResponsiblePerson
			cpBusiness.Youliang = distributionCompany.YouliangResponsiblePerson

			appcache.Bm.Put(cachekey, cpBusiness, 1*time.Hour)
		} else {
			cpBusiness = cacheinfo.(business)
		}

		// 判断该渠道接入的商务负责人是否需要修改
		if game.BodyMy == 1 {
			if game.Issue != cpBusiness.Yunduan && cpBusiness.Yunduan != 0 {
				game.Issue = cpBusiness.Yunduan
				o.Update(&game)
			}
		} else if game.BodyMy == 2 {
			if game.Issue != cpBusiness.Youliang && cpBusiness.Youliang != 0 {
				game.Issue = cpBusiness.Youliang
				o.Update(&game)
			}
		}
	}
}

// 第二次以后同步只需要同步发行商负责人有更改的
func
SyncGameAccessPerson() {
	o := orm.NewOrm()

	var games []models.Game
	var distributionCompanies []models.DistributionCompany

	// 只需要查询在最近一天有更新的渠道商信息
	o.QueryTable(new(models.DistributionCompany)).Filter("update_time__gte", time.Now().Unix()-25*60*60).All(&distributionCompanies)

	for _, company := range distributionCompanies {
		// 先获取到该发行商的游戏接入人
		o.QueryTable(new(models.Game)).Filter("game_id__gte", 0).Filter("issue", company.CompanyId).All(&games)

		for _, game := range games {
			// 判断该游戏接入是否需要修改
			if game.BodyMy == 1 {
				if game.AccessPerson != company.YunduanResponsiblePerson && company.YunduanResponsiblePerson != 0 {
					game.AccessPerson = company.YunduanResponsiblePerson
					o.Update(&game)
				}
			} else if game.BodyMy == 2 {
				if game.AccessPerson != company.YouliangResponsiblePerson && company.YouliangResponsiblePerson != 0 {
					game.AccessPerson = company.YouliangResponsiblePerson
					o.Update(&game)
				}
			}
		}
	}
}

// 修改主合同状态，改为已到期
func ChangeMainContractState() (err error) {
	o := orm.NewOrm()

	var cons []models.MainContract

	timeNow := time.Now().Format("2006-01-02")
	qs := o.QueryTable(new(models.MainContract)).Filter("state__exact", 160).Filter("end_time__lte", timeNow)

	_, err = qs.All(&cons)
	if err != nil {
		return
	}

	for _, con := range cons {
		con.State = 161
		if _, err = o.Update(&con); err != nil {
			return
		} else {
			log.Info("INFO", "Update MainContractState Success,id = ", con.Id)
		}
	}

	return
}
