package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowHistoryControllers"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowHistoryControllers"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowHistoryControllers"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowHistoryControllers"],
		beego.ControllerComments{
			Method: "GetLog",
			Router: `/get_log/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowNameControllers"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowNameControllers"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowNodeControllers"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowNodeControllers"],
		beego.ControllerComments{
			Method: "GetUsersByNextNodeId",
			Router: `/get_next_users/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowNodeControllers"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowNodeControllers"],
		beego.ControllerComments{
			Method: "GetWorkflowNodeByUid",
			Router: `/workflow_node_id/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowSendPackage"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowSendPackage"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/task/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowSendPackage"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowSendPackage"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowSendPackage"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowSendPackage"],
		beego.ControllerComments{
			Method: "IsHandle",
			Router: `/is_handle`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"],
		beego.ControllerComments{
			Method: "GetOne",
			Router: `/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"],
		beego.ControllerComments{
			Method: "GetChannelInfoIdByTaskId",
			Router: `/task/:id`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"],
		beego.ControllerComments{
			Method: "ChangeProgress",
			Router: `/`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"],
		beego.ControllerComments{
			Method: "GetAllByCondition",
			Router: `/condition`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"] = append(beego.GlobalControllerRouter["kuaifa.com/kuaifa/work-together/controllers/workflow:WorkflowTaskController"],
		beego.ControllerComments{
			Method: "ChangeStatucbyId",
			Router: `/status`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
