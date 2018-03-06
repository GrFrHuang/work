package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
)

// 初始化权限表的数据
func main() {
	// v1 目前只有 (能) 权限, type = 2
	type p struct {
		PermissionModel string
		Actions         []string
		Group           string // 别名(组)
	}

	o := orm.NewOrm()

	oldPs := []models.Permission{}
	o.QueryTable(new(models.Permission)).All(&oldPs)
	oldId := map[string]int{}
	for _, p := range oldPs {
		oldId[p.Name] = p.Id
	}

	pNames := []string{}
	ps := []p{}

	pNames = append(pNames, "查看日志")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_ALARM_LOG,
		Actions:         []string{"查"},
		Group:           "sys",
	})
	pNames = append(pNames, "查看游戏流水")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_ORDER,
		Actions:         []string{"查"},
	})

	// game
	pNames = append(pNames, "查看上线准备概况")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN,
		Actions:         []string{"查"},
	})
	// 测评
	pNames = append(pNames, "查看游戏测评")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN_RESULT,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "编辑游戏测评")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN_RESULT,
		Actions:         []string{"改", "删", "增"},
	})

	// 提测
	pNames = append(pNames, "查看游戏提测")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN_PUB,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "编辑游戏提测")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN_PUB,
		Actions:         []string{"改", "删"},
	})
	pNames = append(pNames, "添加游戏提测")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN_PUB,
		Actions:         []string{"增"},
	})

	// 运营准备
	pNames = append(pNames, "查看游戏运营准备")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN_OPERATOR,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "编辑游戏运营准备")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN_OPERATOR,
		Actions:         []string{"改"},
	})
	pNames = append(pNames, "查看游戏客服准备")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN_CUSTOMER,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "编辑游戏客服准备")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME_PLAN_CUSTOMER,
		Actions:         []string{"改"},
	})
	pNames = append(pNames, "查看游戏")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "添加游戏")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME,
		Actions:         []string{"增"},
	})
	pNames = append(pNames, "编辑游戏")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_GAME,
		Actions:         []string{"改", "删"},
	})

	/* contract*/
	pNames = append(pNames, "查看CP合同")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CONTRACT_CP,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "编辑CP合同")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CONTRACT_CP,
		Actions:         []string{"改"},
	})
	pNames = append(pNames, "查看渠道合同")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CONTRACT_CHANNEL,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "编辑渠道合同")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CONTRACT_CHANNEL,
		Actions:         []string{"改"},
	})

	/* 对账*/
	pNames = append(pNames, "查看CP对账")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CP_VERIFY_ACCOUNT,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "添加CP对账")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CP_VERIFY_ACCOUNT,
		Actions:         []string{"增"},
	})
	pNames = append(pNames, "编辑CP对账")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CP_VERIFY_ACCOUNT,
		Actions:         []string{"改"},
	})
	pNames = append(pNames, "查看渠道对账")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CHANNEL_VERIFY_ACCOUNT,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "添加渠道对账")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CHANNEL_VERIFY_ACCOUNT,
		Actions:         []string{"增"},
	})
	pNames = append(pNames, "编辑渠道对账")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CHANNEL_VERIFY_ACCOUNT,
		Actions:         []string{"改"},
	})

	/* cp结算*/
	pNames = append(pNames, "查看CP结算")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_SETTLE_DOWN_ACCOUNT,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "添加CP结算")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_SETTLE_DOWN_ACCOUNT,
		Actions:         []string{"增"},
	})
	pNames = append(pNames, "编辑CP结算")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_SETTLE_DOWN_ACCOUNT,
		Actions:         []string{"改"},
	})

	/* 渠道回款*/
	pNames = append(pNames, "查看渠道结算")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_REMIT_DOWN_ACCOUNT,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "添加渠道结算")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_REMIT_DOWN_ACCOUNT,
		Actions:         []string{"增"},
	})
	pNames = append(pNames, "编辑渠道结算")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_REMIT_DOWN_ACCOUNT,
		Actions:         []string{"改"},
	})

	/* 用户*/
	pNames = append(pNames, "查看用户")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_USER,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "添加用户")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_USER,
		Actions:         []string{"增"},
	})
	pNames = append(pNames, "编辑用户")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_USER,
		Actions:         []string{"改", "删"},
	})

	/* 部门*/
	pNames = append(pNames, "查看部门")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_DEPARTMENT,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "添加部门")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_DEPARTMENT,
		Actions:         []string{"增"},
	})
	pNames = append(pNames, "编辑部门")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_DEPARTMENT,
		Actions:         []string{"改", "删"},
	})

	// 发行商
	pNames = append(pNames, "查看发行商")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_DISTRIBUTION_COMPANY,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "编辑发行商")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_DISTRIBUTION_COMPANY,
		Actions:         []string{bean.PMSA_UPDATE, bean.PMSA_DELETE, "增"},
	})
	// 渠道商
	pNames = append(pNames, "查看渠道商")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CHANNEL_COMPANY,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "编辑渠道商")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_CHANNEL_COMPANY,
		Actions:         []string{bean.PMSA_UPDATE, bean.PMSA_DELETE, "增"},
	})
	// 研发商
	pNames = append(pNames, "查看研发商")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_DEVELOP_COMPANY,
		Actions:         []string{"查"},
	})
	pNames = append(pNames, "编辑研发商")
	ps = append(ps, p{
		PermissionModel: bean.PMSM_DEVELOP_COMPANY,
		Actions:         []string{bean.PMSA_UPDATE, bean.PMSA_DELETE, "增"},
	})

	dele := &models.Permission{
		Readonly: 1,
	}
	_, err := o.Delete(dele, "Readonly")
	if err != nil {
		panic(err)
	}

	_, err = o.Insert(&models.Permission{
		Name:     "超管",
		Readonly: 1,
		Model:    "super",
		Type:     models.Type_supper,
		Methods:  "",
		Id:       1,
	})

	if err != nil {
		panic(err)
		return
	}

	// 全部添加到数据库里面
	for i, name := range pNames {
		b, _ := json.Marshal(&ps[i].Actions)
		pe := &models.Permission{
			Name:     name,
			Id:       oldId[name],
			Readonly: 1,
			Model:    ps[i].PermissionModel,
			Type:     models.Type_can,
			Methods:  string(b),
		}
		o.Insert(pe)
	}
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", beego.AppConfig.String("mysqluser"),
		beego.AppConfig.String("mysqlpass"), beego.AppConfig.String("mysqlurls"),
		beego.AppConfig.String("mysqlport"), beego.AppConfig.String("mysqldb"))
	orm.RegisterDataBase("default", "mysql", link)
	orm.Debug = true
}
