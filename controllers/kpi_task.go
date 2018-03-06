package controllers

import (
	"encoding/json"
	"kuaifa.com/kuaifa/work-together/models"
	"strconv"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"kuaifa.com/kuaifa/work-together/tool"
	"time"
	"runtime"
	"github.com/astaxie/beego/orm"
	"github.com/tealeg/xlsx"
	"github.com/astaxie/beego"
)

var taskPool map[int]chan int

func init() {
	taskPool = make(map[int]chan int)
}

// KpiTaskController operations for KpiTask
type KpiTaskController struct {
	BaseController
}

// URLMapping ...
func (c *KpiTaskController) URLMapping() {

}

// Post ...
// @Title Post
// @Description create KpiTask
// @Param	body		body 	models.KpiTask	true		"body for KpiTask content"
// @Success 201 {int} models.KpiTask
// @Failure 403 body is empty
// @router / [post]
func (c *KpiTaskController) Post() {
	var v models.KpiTask
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if v.AssesseserId == v.PublisherId {
			c.RespJSON(bean.CODE_Params_Err, "考核人与发布人不能相同")
			return
		}
		v.PublisherId = c.Uid()
		if id, err := models.AddKpiTask(&v); err == nil {
			//Add child task
			if len(v.KpiChildTasks) == 0 {
				c.RespJSON(bean.CODE_Bad_Request, "任务不能为空")
				return
			}
			err = models.AddKpiChildTasks(v.KpiChildTasks)
			if err != nil {
				//Rollback task table
				models.DeleteKpiTask(int(id))
				c.RespJSON(bean.CODE_Internal_Server_Error, "添加子任务失败"+err.Error())
				return
			}
			////If timing task
			//if v.PublishType == 1 {
			//	//Add timing task into work pool
			//	ch := make(chan int)
			//	ch <- v.TaskPublishDate
			//	taskPool[int(id)] = ch
			//}
			c.RespJSONData(v)
			return
		} else {
			c.RespJSON(bean.CODE_Forbidden, "添加任务失败"+err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Bad_Request, "请求参数错误"+err.Error())
		return
	}
}

// GetOne ...
// @Title Get One
// @Description get KpiTask by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.KpiTask
// @Failure 403 :id is empty
// @router /:id [get]
func (c *KpiTaskController) GetOne() {
	var uids = []int{}
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetKpiTaskById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, "获取任务详情失败"+err.Error())
		return
	}
	//Add the child task list
	err = models.AddChildTask(v)
	if err != nil {
		//Ignore error with continue execute
		c.RespJSON(bean.CODE_Internal_Server_Error, "获取子任务信息失败"+err.Error())
		return
	}
	//Add user name
	uids = append(uids, v.PublisherId)
	uids = append(uids, v.AssesseserId)
	names, err := models.GetUsersName(uids)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, "获取姓名失败"+err.Error())
		return
	}
	v.AssesseserName = names[v.AssesseserId]
	v.PublisherName = names[v.PublisherId]
	c.RespJSONData(v)
}

// GetAll ...
// @Title Get All
// @Description get KpiTask
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Param	startTime	query	string	false	"开始时间"
// @Param	endTime	    query	string	false	"结束时间"
// @Success 200 {object} models.KpiTask
// @Failure 403
// @router / [get]
func (c *KpiTaskController) GetAll() {
	var kts = []*models.KpiTask{}
	var uids = []int{}
	f, _ := tool.BuildFilter(c.Controller, 20)
	//Filtrate by time
	startTime, _ := c.GetInt("startTime")
	endTime, _ := c.GetInt("endTime")
	if startTime != 0 && endTime != 0 {
		f.Where["task_complete_date__gte"] = []interface{}{startTime}
		f.Where["task_complete_date__lte"] = []interface{}{endTime}
	}
	//The one Just see oneself's department task
	user, _ := models.GetUserById(c.Uid())

	//test set department id = 243
	user.DepartmentId = 243
	f.Where["department_id"] = []interface{}{user.DepartmentId}
	total, err := tool.GetAllByFilterWithTotal(new(models.KpiTask), &kts, f)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, "获取任务列表失败"+err.Error())
		return
	}
	//If kpi tasks array is empty
	if total == 0 {
		c.RespJSONDataWithTotal(kts, total)
		return
	}
	for _, v := range kts {
		uids = append(uids, v.AssesseserId)
		uids = append(uids, v.PublisherId)
	}
	names, err := models.GetUsersName(uids)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, "获取姓名失败"+err.Error())
		return
	}
	for k, v := range kts {
		kts[k].AssesseserName = names[v.AssesseserId]
		kts[k].PublisherName = names[v.PublisherId]
	}
	//Add total score
	kts, err = models.AddTotalScore(kts)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, "获取得分信息失败"+err.Error())
		return
	}
	c.RespJSONDataWithTotal(kts, total)
}

// Put ...
// @Title Put
// @Description update the KpiTask
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.KpiTask	true		"body for KpiTask content"
// @Success 200 {object} models.KpiTask
// @Failure 403 :id is not int
// @router /:id [put]
func (c *KpiTaskController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.KpiTask{Id: id}
	task, _ := models.GetKpiTaskById(id)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateKpiTaskById(&v); err == nil {
			if v.PublishType == 1 {
				//Append time value into channel for kill goroutine
				if task.TaskPublishDate != v.TaskPublishDate {
					ch := make(chan int)
					ch <- v.TaskPublishDate
					taskPool[id] = ch
				}
			}
			c.RespJSONSuccess()
			return
		} else {
			c.RespJSON(bean.CODE_Forbidden, "修改任务失败"+err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Bad_Request, "请求参数错误"+err.Error())
		return
	}
}

// Delete ...
// @Title Delete
// @Description delete the KpiTask
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *KpiTaskController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	task, err := models.GetKpiTaskById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, "未找到该任务"+err.Error())
		return
	}
	if err != nil && err != orm.ErrNoRows {
		c.RespJSON(bean.CODE_Internal_Server_Error, "服务器错误"+err.Error())
		return
	}
	if task.PublishState == 1 {
		c.RespJSON(bean.CODE_Forbidden, "任务发布中，禁止删除")
		return
	}
	if err := models.DeleteKpiTask(id); err == nil {
		//close(taskPool[id])
		c.RespJSONSuccess()
		return
	} else {
		c.RespJSON(bean.CODE_Internal_Server_Error, "删除任务失败"+err.Error())
		return
	}
}

// Put ...
// @Title Put
// @Description put the KpiTask
// @Param	body		body 	[]models.KpiChildTask	true		"body for KpiTask content"
// @Success 200 {string} update success!
// @Failure 403 id is empty
// @router /examining [put]
func (c *KpiTaskController) Examine() {
	tasks := []*models.KpiChildTask{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &tasks); err == nil {
		if len(tasks) == 0 {
			c.RespJSONSuccess()
			return
		}
		err := models.ExamineChildTask(tasks)
		if err != nil {
			c.RespJSON(bean.CODE_Forbidden, "子任务审核失败"+err.Error())
			return
		}
	} else {
		c.RespJSON(bean.CODE_Internal_Server_Error, "请求参数错误"+err.Error())
		return
	}
}

// Put ...
// @Title Put
// @Description put the KpiTask
// @Param	id		path 	string	true		"The id you want to put"
// @Param	body		body 	[]models.KpiChildTask	true		"body for KpiTask content"
// @Success 200 {string} update success!
// @Failure 403 id is empty
// @router /submittingTask/:id [put]
func (c *KpiTaskController) SubmitTask() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	var childTasks models.ChildTasks
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &childTasks); err == nil {
		if &childTasks == nil {
			c.RespJSON(bean.CODE_Bad_Request, "请求参数错误"+err.Error())
			return
		}
	}
	if err := models.WriteTask(&childTasks, id); err == nil {
		c.RespJSONSuccess()
		return
	} else {
		c.RespJSON(bean.CODE_Internal_Server_Error, "任务填写失败"+err.Error())
		return
	}
}

// Get ...
// @Title single person export kpi excel
// @Description single person export kpi excel
// @Param	id		path 	string	true		"The task id you want to export"
// @Success 200 {string} get success!
// @Failure 403 id is empty
// @router /singleExporting/:id [get]
func (c *KpiTaskController) ExportKpiExcel() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	task, err := models.GetKpiTaskById(id)
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, "未找到kpi任务"+err.Error())
		return
	}
	if task.AuditState != 1 {
		c.RespJSON(bean.CODE_Forbidden, "不能导出未审批的kpi文件")
		return
	}
	childTasks, err := models.GetChildTaskByTaskId(id)
	if err != nil {
		c.RespJSON(bean.CODE_Internal_Server_Error, "获取任务详情失败"+err.Error())
		return
	}
	user, err := models.GetUserById(c.Uid())
	if err != nil {
		c.RespJSON(bean.CODE_Forbidden, "kpi考核人信息获取失败"+err.Error())
		return
	}
	department, err := models.GetDepartmentById(user.DepartmentId)
	if err != nil {
		//Ignore current user's department info error
		beego.Error("获取部门信息失败" + err.Error())
		user.DepartmentName = "未知部门"
	}
	user.DepartmentName = department.Name

	file, err := xlsx.OpenFile("../work-together/docs/kpi.xlsx")
	if err != nil {
		c.RespJSON(bean.CODE_Not_Found, "未找到kpi模版文件"+err.Error())
		return
	}
	sheet := file.Sheet["Sheet1"]
	completeDate := time.Unix(int64(task.TaskCompleteDate), 0).Format("2006-01-02")
	entryDate := time.Unix(int64(user.EntryTime), 0).Format("2006-01-02")

	//Edit department|position|name|date
	sheet.Rows[1].Cells[0].Value = "部门: " + user.DepartmentName +
		" 职位: " + user.Position +
		" 姓名: " + user.Name +
		" 入职日期: " + entryDate

	//Add the task complete date
	sheet.Rows[2].Cells[1].Value = completeDate

	for k, v := range childTasks {
		index := k + 4
		sheet.Rows[index].Cells[1].Value = strconv.Itoa(k + 1)
		sheet.Rows[index].Cells[2].Value = v.TaskName
		sheet.Rows[index].Cells[3].Value = v.Period
		sheet.Rows[index].Cells[4].Value = strconv.FormatFloat(v.ProgressRate, 'f', 2, 64)
		sheet.Rows[index].Cells[5].Value = v.Remark
		sheet.Rows[index].Cells[6].Value = strconv.FormatFloat(v.Score, 'f', 2, 64)
		sheet.Rows[index].Cells[7].Value = v.Annotations
	}
	//Set allow cross domain
	c.AllowCross()
	c.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+user.Name+" "+completeDate+" kpi.xls")
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/vnd.open")
	err = file.Write(c.Ctx.ResponseWriter)
	if err != nil {
		beego.Error("kpi文件导出错误" + err.Error())
	}
}

func publishTask(taskPublishDate, id int) {
	subTime := taskPublishDate - int(time.Now().Unix())
	duration, _ := time.ParseDuration(strconv.Itoa(subTime) + "s")
	ticker := time.NewTicker(duration)
	go func(ticker *time.Ticker) {
		for range ticker.C {
			task := models.KpiTask{Id: id, PublishState: 1}
			err := models.UpdateKpiTaskFields(&task, "publish_state")
			if err != nil {
				//todo republish task
			}
			return
		}
	}(ticker)
}

func ExecuteListenTask() {
	for k, v := range taskPool {
		go func(id int) {
			select {
			case <-taskPool[id]:
				if <-taskPool[id] != 0 {
					defer publishTask(<-v, id)
				}
				runtime.Goexit()
			}
			publishTask(<-v, id)
		}(k)
	}
}
