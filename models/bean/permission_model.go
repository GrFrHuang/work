package bean

import "strings"

const (
	PMSA_SELECT = "查"
	PMSA_INSERT = "增"
	PMSA_DELETE = "删"
	PMSA_UPDATE = "改"
)
const (
	PMSM_SUPER                  = "super"
	PMSM_USER                   = "user"
	PMSM_SETTLE_DOWN_ACCOUNT    = "settle_down_account"
	PMSM_REMIT_DOWN_ACCOUNT     = "remit_down_account"
	PMSM_GAME_PLAN              = "game_plan"
	PMSM_GAME_PLAN_OPERATOR     = "game_plan_operator" // 营运
	PMSM_GAME_PLAN_RESULT       = "game_plan_result"   // 测评
	PMSM_GAME_PLAN_PUB          = "game_plan_pub"      // 游戏提测
	PMSM_GAME_PLAN_CUSTOMER     = "game_plan_customer" // 客服
	PMSM_GAME                   = "game"
	PMSM_GAME_ALL               = "game_all"
	PMSM_GAME_CHANNEL_ACCESS    = "game_channel_access" //渠道接入
	PMSM_GAME_UPDATE            = "game_update"         //游戏更新
	PMSM_EXPRESS_MANAGE         = "express_manage"      //快递管理
	PMSM_CONTRACT_CP            = "contract_cp"
	PMSM_CONTRACT_CHANNEL       = "contract_channel"
	PMSM_CP_VERIFY_ACCOUNT      = "cp_verify_account"
	PMSM_CHANNEL_VERIFY_ACCOUNT = "channel_verify_account"
	PMSM_DEPARTMENT             = "department"
	PMSM_ALARM_RULE             = "alarm_rule"
	PMSM_ORDER                  = "order"       // 流水
	PMSM_LOGS                   = "logs"        // logs
	PMSM_ALARM_LOG              = "alarm_log"   // 报警日志
	PMSM_GAME_TYPE              = "game_type"   // 游戏类型
	PMSM_COOPERATION            = "cooperation" // 合作方式
	PMSM_GAME_OUTAGE            = "game_outage" //游戏停运

	PMSM_ISSUE          = "issue"          // 发行商 unused
	PMSM_SOCIATY_POLICY = "sociaty_policy" // 工行政策
	PMSM_RESULT         = "result"         // 评测结果
	PMSM_SDK_STATUS     = "sdk_status"     // sdk接入状态 unused
	PMSM_TYPES          = "types"          // 类型集合

	PMSM_DEVELOPMENT = "development" // 研发商 unused

	PMSM_CHANNEL_COMPANY      = "channel_company"      // 渠道商
	PMSM_DEVELOP_COMPANY      = "develop_company"      // 研发商
	PMSM_DISTRIBUTION_COMPANY = "distribution_company" // 发行商
	PMSM_WARNING              = "warning"              // 预警
	PMSM_WARNING_TYPE         = "warning_type"         // 预警类型
	PMSM_WARNING_LOG          = "warning_log"          // 预警日志
	PMSM_COMPANY_TYPE         = "company_type"         // 公司类型

	PMSM_STATISTIC_FINANCIAL     = "statistic_financial"     // 财务部门统计
	PMSM_STATISTIC_ACCOUNTING    = "statistic_account"       // 结算部统计
	PMSM_STATISTIC_CHANNEL_TRADE = "statistic_channel_trade" // 渠道商务部门统计
	PMSM_STATISTIC_CP_TRADE      = "statistic_cp_trade"      // cp商务部门统计
	PASM_STATISTIC_OPERATION     = "statistic_operation"

	PMSM_USER_FACE_EXAMINE   = "user_face"   // 人脸审核
	PMSM_USER_SIGN_EXAMINE   = "user_sign"   // 签到审核
	PMSM_SMALL_SHEEP_EXAMINE = "small_sheep" // 小绵羊
)

type Router struct {
	Url    string
	Action string // todo 如果设置了action,则有这个action的权限才显示,
}

type P struct {
	Routers []Router
}

// 对应权限能访问的路由
var MenuMap = map[string]P{
	PMSM_GAME_PLAN: {
		Routers: []Router{
			{
				Url: "home/ready", // 概述
			},
		},
	},
	PMSM_GAME_PLAN_RESULT: {
		Routers: []Router{
			{
				Url: "home/GameEvaluation", // 游戏测评
			},
		},
	},
	PMSM_GAME_PLAN_PUB: {
		Routers: []Router{
			{
				Url: "home/reference", // 游戏提测
			},
		},
	},
	PMSM_GAME_CHANNEL_ACCESS: {
		Routers: []Router{
			{
				Url: "home/channelAccess", //渠道接入
			},
		},
	},
	PMSM_GAME_UPDATE: {
		Routers: []Router{
			{
				Url: "home/gameUpdate", //游戏更新
			},
		},
	},
	PMSM_GAME_OUTAGE: {
		Routers: []Router{
			{
				Url: "home/gameOutage", //游戏停运
			},
		},
	},
	PMSM_EXPRESS_MANAGE: {
		Routers: []Router{
			{
				Url: "home/expressManagement", //快递管理
			},
		},
	},
	PMSM_GAME_PLAN_CUSTOMER: {
		Routers: []Router{
			{
				Url: "home/game-manage/pre-online/pre_customer_service", //客服
			},
		},
	},
	PMSM_GAME_PLAN_OPERATOR: {
		Routers: []Router{
			{
				Url: "home/game-manage/pre-online/pre_operation", // 运营
			},
		},
	},

	PMSM_GAME: {
		Routers: []Router{
			{
				Url: "home/gameAccess", // 游戏接入
			},
		},
	},
	PMSM_ORDER: {
		Routers: []Router{
			{
				Url: "home/list", // 流水
			},
		},
	},
	PMSM_CONTRACT_CP: {
		Routers: []Router{
			{
				Url: "home/mgt/mgta", // cp合同
			},
		},
	},
	PMSM_CONTRACT_CHANNEL: {
		Routers: []Router{
			{
				Url: "home/channelContract/channelContractA", // channel合同
			},
		},
	},

	PMSM_CP_VERIFY_ACCOUNT: {
		Routers: []Router{
			{
				Url: "home/mgt/mgtb", // cp 对账
			},
		},
	},
	PMSM_CHANNEL_VERIFY_ACCOUNT: {
		Routers: []Router{
			{
				Url: "home/channelContract/channelContractB", // 渠道对账
			},
		},
	},
	PMSM_SETTLE_DOWN_ACCOUNT: {
		Routers: []Router{
			{
				Url: "home/mgt/mgtc", // 结算
			},
		},
	},
	PMSM_REMIT_DOWN_ACCOUNT: {
		Routers: []Router{
			{
				Url: "home/channelContract/channelContractC", // 回款
			},
		},
	},
	PMSM_USER: {
		Routers: []Router{
			{
				Url: "home/user", // 角色管理
			}, {
				Url: "home/role", // 角色管理
			},
		},
	},
	PMSM_LOGS: {
		Routers: []Router{
			{
				Url: "home/Operation", // 查看日志
			},
		},
	},
	PMSM_ALARM_LOG: {
		Routers: []Router{
			{
				Url: "home/Alarm", // 报警日志
			},
		},
	},
	PMSM_CHANNEL_COMPANY: {
		Routers: []Router{
			{
				Url: "home/addresslist", // 渠道商
			},
		},
	},
	PMSM_DISTRIBUTION_COMPANY: {
		Routers: []Router{
			{
				Url: "home/addresslist", // 发行商
			},
		},
	},
	PMSM_DEVELOP_COMPANY: {
		Routers: []Router{
			{
				Url: "home/addresslist", // 研发商
			},
		},
	},
	PMSM_COMPANY_TYPE: {
		Routers: []Router{
			{
				Url: "home/addresslist", // 公司列表
			},
		},
	},
	PMSM_WARNING: {
		Routers: []Router{
			{
				Url: "home/warning_center", // 
			},
		},
	},
	PMSM_WARNING_LOG: {
		Routers: []Router{
			{
				Url: "home/warning", // 研发商
			},
		},
	},
	PMSM_STATISTIC_FINANCIAL: {
		Routers: []Router{
			{
				Url: "home/departmentStatistics", // 研发商
			},
		},
	},
	PMSM_STATISTIC_ACCOUNTING: {
		Routers: []Router{
			{
				Url: "home/departmentStatistics", // 研发商
			},
		},
	},
	PMSM_STATISTIC_CHANNEL_TRADE: {
		Routers: []Router{
			{
				Url: "home/departmentStatistics", // 研发商
			},
		},
	},
	PMSM_STATISTIC_CP_TRADE: {
		Routers: []Router{
			{
				Url: "home/departmentStatistics", // 研发商
			},
		},
	},
	PASM_STATISTIC_OPERATION: {
		Routers: []Router{
			{
				Url: "home/departmentStatistics", // 研发商
			},
		},
	},
	PMSM_USER_FACE_EXAMINE: {
		Routers: []Router{
			{
				Url: "home/user/examine", // 人脸审核
			},
		},
	},
	PMSM_SMALL_SHEEP_EXAMINE: {
		Routers: []Router{
			{
				Url: "home/user/smallsheep", // 人脸审核
			},
		},
	},
}

// 获取模型对应的权限名称和分组
func PermissionModelString(i string) (name, group string) {
	modelMap := map[string]string{
		PMSM_SUPER: "超管",

		PMSM_ALARM_LOG:               "日志",
		PMSM_ORDER:                   "游戏流水",
		PMSM_GAME_PLAN:               "游戏上线准备|游戏",
		PMSM_GAME_PLAN_RESULT:        "游戏测评|游戏",
		PMSM_GAME_PLAN_PUB:           "游戏提测|游戏",
		PMSM_GAME_PLAN_OPERATOR:      "游戏运营准备|游戏",
		PMSM_GAME_PLAN_CUSTOMER:      "游戏客服准备|游戏",
		PMSM_GAME:                    "游戏",
		PMSM_GAME_CHANNEL_ACCESS:     "渠道接入|游戏",
		PMSM_GAME_UPDATE:             "游戏更新|游戏",
		PMSM_GAME_OUTAGE:             "游戏停运|游戏",
		PMSM_EXPRESS_MANAGE:          "快递管理",
		PMSM_CONTRACT_CP:             "CP合同",
		PMSM_CONTRACT_CHANNEL:        "渠道合同",
		PMSM_CP_VERIFY_ACCOUNT:       "CP对账",
		PMSM_CHANNEL_VERIFY_ACCOUNT:  "渠道对账",
		PMSM_SETTLE_DOWN_ACCOUNT:     "CP结算",
		PMSM_REMIT_DOWN_ACCOUNT:      "渠道结算",
		PMSM_USER:                    "用户",
		PMSM_DEPARTMENT:              "部门",
		PMSM_DISTRIBUTION_COMPANY:    "发行商",
		PMSM_DEVELOP_COMPANY:         "研发商",
		PMSM_CHANNEL_COMPANY:         "渠道商",
		PMSM_COMPANY_TYPE:            "公司列表",
		PMSM_WARNING:                 "预警",
		PMSM_WARNING_LOG:             "预警日志",
		PMSM_STATISTIC_FINANCIAL:     "财务部统计",
		PMSM_STATISTIC_ACCOUNTING:    "结算部统计",
		PMSM_STATISTIC_CHANNEL_TRADE: "渠道商务部统计",
		PMSM_STATISTIC_CP_TRADE:      "cp商务部统计",
		PASM_STATISTIC_OPERATION:     "运营部统计",
		PMSM_USER_FACE_EXAMINE:       "人脸审核|用户",
		PMSM_USER_SIGN_EXAMINE:       "签到审核|用户",
		PMSM_SMALL_SHEEP_EXAMINE:     "小绵羊|用户",
	}
	if v, ok := modelMap[i]; ok {
		ng := strings.Split(v, "|")
		name = ng[0]
		if len(ng) == 2 {
			group = ng[1]
		} else {
			group = name
		}
		return
	}
	return
}
