/*
                   _ooOoo_
                  o8888888o
                  88" . "88
                  (| -_- |)
                  O\  =  /O
               ____/`---'\____
             .'  \\|     |//  `.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||-  \
           |   | \\\  -  /// |   |
           | \_|  ''\---/''  |   |
           \  .-\__  `-`  ___/-. /
         ___`. .'  /--.--\  `. . __
      ."" '<  `.___\_<|>_/___.'  >'"".
     | | :  `- \`.;`\ _ /`;.`/ - ` : | |
     \  \ `-.   \_ __\ /__ _/   .-` /  /
======`-.____`-.___\_____/___.-`____.-'======
                   `=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
         佛祖保佑       永无BUG
*/
package workflow

import (
	"kuaifa.com/kuaifa/work-together/controllers"
	"kuaifa.com/kuaifa/work-together/models/workflow"
)

type WorkflowHistoryControllers struct {
	controllers.BaseController
}


// GetAll ...
// @Title Get All
// @Description get WorkflowName
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.WorkflowName
// @Failure 403
// @router / [get]
func (c *WorkflowHistoryControllers) GetAll(){
	logs ,err := workflow.GetAll()
	if err!=nil {
		c.RespJSONData(err.Error())
	}

	c.RespJSONData(logs)

}

// GetAll ...
// @Title Get All
// @Description get WorkflowName
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.WorkflowName
// @Failure 403
// @router /get_log/ [get]
func (c *WorkflowHistoryControllers) GetLog(){
	id ,err:= c.GetInt("id",0)
	if err!=nil{
		c.RespJSONData("参数格式错误")
	}
	node_id ,err:= c.GetInt("node_id",0)
	if err!=nil{
		c.RespJSONData("参数格式错误")
	}
	user_id ,err := c.GetInt("user_id",0)
	if err!=nil{
		c.RespJSONData("参数格式错误")
	}
	user_name 	:=c.GetString("user_name","")
	workflow_id ,err :=c.GetInt("workflow_id",0)
	if err!=nil{
		c.RespJSONData("参数格式错误")
	}
	task_id,err := c.GetInt("task_id",0)
	if err!=nil{
		c.RespJSONData("参数格式错误")
	}

	where := make(map[string]interface{})

	if id!=0 {
		where[" and id=?"]=id
	}
	if node_id != 0 {
		where[" and node_id=?"]=node_id
	}
	if user_id != 0{
		where[" and user_id=?"]=user_id
	}
	if user_name != "" {
		where[" and user_name=?"]=user_name
	}
	if workflow_id != 0 {
		where[" and workflow_id=?"]=workflow_id
	}
	if task_id != 0 {
		where[" and task_id=?"]=task_id
	}
	log,err:= workflow.GetLog(where)
	if err!=nil {
		c.RespJSONData(err.Error())
	}
	c.RespJSONData(log)



}