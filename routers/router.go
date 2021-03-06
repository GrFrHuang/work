// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"kuaifa.com/kuaifa/work-together/controllers"
	"kuaifa.com/kuaifa/work-together/controllers/workflow"
	"kuaifa.com/kuaifa/work-together/filters"
)

func init() {
	api := beego.NewNamespace("/v1",
		beego.NSRouter("/*", &controllers.BaseController{}, "options:Options"),
		beego.NSNamespace("/session",
			beego.NSInclude(
				&controllers.SessionController{},
			),
		),
		beego.NSNamespace("/workflow_name",
			beego.NSInclude(
				&workflow.WorkflowNameControllers{},
			),
		),
		beego.NSNamespace("/desktop",
			beego.NSInclude(
				&controllers.UserDesktopClientController{},
			),
		),
		beego.NSNamespace("/desktop_screen_log",
			beego.NSInclude(
				&controllers.UserDesktopScreenLogControllers{},
			),
		),
		beego.NSNamespace("/workflow_task",
			beego.NSInclude(
				&workflow.WorkflowTaskController{},
			),
		),
		beego.NSNamespace("/workflow_history",
			beego.NSInclude(
				&workflow.WorkflowHistoryControllers{},
			),
		),
		beego.NSNamespace("/workflow_send_package",
			beego.NSInclude(
				&workflow.WorkflowSendPackage{},
			),
		),
		beego.NSNamespace("/workflow_node",
			beego.NSInclude(
				&workflow.WorkflowNodeControllers{},
			),
		),
		beego.NSNamespace("/face",
			beego.NSInclude(
				&controllers.UserFaceControllers{},
			),
		),
		beego.NSNamespace("/userfaceverify",
			beego.NSInclude(
				&controllers.UserFaceVerifyController{},
			),
		),
		beego.NSNamespace("/sign",
			beego.NSInclude(
				&controllers.UserSignController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/userown",
			beego.NSInclude(
				&controllers.UserOwnController{},
			),
		),
		beego.NSNamespace("/department",
			beego.NSInclude(
				&controllers.DepartmentController{},
			),
		),
		beego.NSNamespace("/role",
			beego.NSInclude(
				&controllers.RoleController{},
			),
		),
		beego.NSNamespace("/permission",
			beego.NSInclude(
				&controllers.PermissionController{},
			),
		),
		beego.NSNamespace("/logs",
			beego.NSInclude(
				&controllers.LogsController{},
			),
		),
		beego.NSNamespace("/settleaccount",
			beego.NSInclude(
				&controllers.SettleDownAccountController{},
			),
		),
		beego.NSNamespace("/remitaccount",
			beego.NSInclude(
				&controllers.RemitDownAccountController{},
			),
		),
		beego.NSNamespace("/alarmrule",
			beego.NSInclude(
				&controllers.AlarmRuleController{},
			),
		),
		beego.NSNamespace("/alarmlog",
			beego.NSInclude(
				&controllers.AlarmLogController{},
			),
		),
		beego.NSNamespace("/ladder",
			beego.NSInclude(
				&controllers.LadderController{},
			),
		),
		beego.NSNamespace("/contract",
			beego.NSInclude(
				&controllers.ContractController{},
			),
		),
		beego.NSNamespace("/contractStatistics",
			beego.NSInclude(
				&controllers.ContractStatisticsController{},
			),
		),
		beego.NSNamespace("/game",
			beego.NSInclude(
				&controllers.GameController{},
			),
		),
		beego.NSNamespace("/gameall",
			beego.NSInclude(
				&controllers.GameAllController{},
			),
		),
		beego.NSNamespace("/gameplan",
			beego.NSInclude(
				&controllers.GamePlanController{},
			),
		),
		beego.NSNamespace("/gameupdate",
			beego.NSInclude(
				&controllers.GameUpdateController{},
			),
		),
		beego.NSNamespace("/expressManage",
			beego.NSInclude(
				&controllers.ExpressManageController{},
			),
		),
		beego.NSNamespace("/operatelog",
			beego.NSInclude(
				&controllers.OperateLogController{},
			),
		),
		beego.NSNamespace("/channelaccess",
			beego.NSInclude(
				&controllers.ChannelAccessController{},
			),
		),

		beego.NSNamespace("/channel",
			beego.NSInclude(
				&controllers.ChannelController{},
			),
		),
		beego.NSNamespace("/sociaty",
			beego.NSInclude(
				&controllers.SociatyPolicyController{},
			),
		),
		beego.NSNamespace("/cpverify",
			beego.NSInclude(
				&controllers.CpVerifyAccountController{},
			),
		),
		beego.NSNamespace("/cpsserver",
			beego.NSInclude(
				&controllers.CpsServerController{},
			),
		),
		beego.NSNamespace("/channel_pre_verify",
			beego.NSInclude(
				&controllers.OrderPreVerifyChannelController{},
			),
		),
		beego.NSNamespace("/cp_pre_verify",
			beego.NSInclude(
				&controllers.OrderPreVerifyCpController{},
			),
		),
		beego.NSNamespace("/verify_channel",
			beego.NSInclude(
				&controllers.VerifyChannelController{},
			),
		),
		beego.NSNamespace("/verify_cp",
			beego.NSInclude(
				&controllers.VerifyCpController{},
			),
		),
		beego.NSNamespace("/channelverify",
			beego.NSInclude(
				&controllers.ChannelVerifyAccountController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&controllers.OrderController{},
			),
		),
		beego.NSNamespace("/types",
			beego.NSInclude(
				&controllers.TypesController{},
			),
		),
		beego.NSNamespace("/asset",
			beego.NSInclude(
				&controllers.AssetController{},
			),
		),
		beego.NSNamespace("/dashboardInfo",
			beego.NSInclude(
				&controllers.DashboardInfoController{},
			),
		),
		beego.NSNamespace("/channelCompany",
			beego.NSInclude(
				&controllers.ChannelCompanyController{},
			),
		),
		beego.NSNamespace("/developCompany",
			beego.NSInclude(
				&controllers.DevelopCompanyController{},
			),
		),
		beego.NSNamespace("/distributionCompany",
			beego.NSInclude(
				&controllers.DistributionCompanyController{},
			),
		),
		beego.NSNamespace("/repair",
			beego.NSInclude(
				&controllers.RepairDataController{},
			),
		),
		beego.NSNamespace("/warning",
			beego.NSInclude(
				&controllers.WarningController{},
			),
		),
		beego.NSNamespace("/warninglog",
			beego.NSInclude(
				&controllers.WarningLogController{},
			),
		),
		beego.NSNamespace("/warningtype",
			beego.NSInclude(
				&controllers.WarningTypeController{},
			),
		),
		beego.NSNamespace("/companytype",
			beego.NSInclude(
				&controllers.CompanyTypeController{},
			),
		),
		beego.NSNamespace("/statistic",
			beego.NSInclude(
				&controllers.StatisticController{},
			),
		),
		beego.NSNamespace("/accountInformation",
			beego.NSInclude(
				&controllers.AccountInformationController{},
			),
		),
		beego.NSNamespace("/verify_cp_elec",
			beego.NSInclude(
				&controllers.VerifyCpElectricController{},
			),
		),
		beego.NSNamespace("/game_outage",
			beego.NSInclude(
				&controllers.GameOutageController{},
			),
		),
		beego.NSNamespace("/mainContract",
			beego.NSInclude(
				&controllers.MainContractController{},
			),
		),
		beego.NSNamespace("/kpi",
			beego.NSInclude(
				&controllers.KpiTaskController{},
			),
		),

	)
	beego.AddNamespace(api)
	bapi := beego.NewNamespace("/dc",
		beego.NSRouter("/*", &controllers.BaseController{}, "options:Options"),
		beego.NSNamespace("/desktop",
			beego.NSInclude(
				&controllers.DesktopClientController{},
			),
		),
	)
	beego.AddNamespace(bapi)

	beego.InsertFilter("/v1/*", beego.BeforeRouter, filters.AuthLogin, true) // 验证登陆
	beego.InsertFilter("/v1/*", beego.AfterExec, filters.Logger, false)      // 日志

	beego.Router("/login", &controllers.AuthViewController{}, "get:Login")

}
