package workflow

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"kuaifa.com/kuaifa/work-together/models"
	"errors"
)

type WorkflowNode struct {
	Id                    int    `orm:"column(id);auto" json:"id"`
	NodeName              string `orm:"column(node_name);size(100);null" description:"结点名称" json:"node_name"`
	DepartmentId          int    `orm:"column(department_id);null" description:"部门ID" json:"department_id"`
	DepartmentName        string `orm:"column(department_name);size(100);null" description:"角色信息" json:"department_name"`
	Permissions           string `orm:"column(permissions);size(100);null" description:"权限（eg:1,2,3,4  |增,删,改,查）" json:"permissions"`
	Rollback              int    `orm:"column(rollback);null" description:"1 能回退,2不行 " json:"rollback"`
	WfNameId              int    `orm:"column(wf_name_id);null" description:"1 能回退,2不行 " json:"wf_name_id"`
	RollbackId            int    `orm:"column(rollback_id);null" description:"1 能回退,2不行 " json:"rollback_id"`
	CurrentNodeId         string `orm:"column(current_node_id);size(100);null" description:"上一步结点，号分割" json:"current_node_id"`
	NextNodeId            string `orm:"column(next_node_id);size(100);null" description:"下一步结点，号分割" json:"next_node_id"`
	WorkflowShow          int    `orm:"column(workflow_show);null" description:" 1显示渠道信息 2显示合同信息 " json:"workflow_show"`
	ChannelProportionHide int    `orm:"column(channel_proportion_hide);null" description:"显示分成信息" json:"channel_proportion_hide"`
	Step                  int    `orm:"column(step);null" json:"step"`
}

type OutWorkflowNode struct {
	Department string         `json:"department"`
	Data       *[]models.User `json:"data"`
}

func (t *WorkflowNode) TableName() string {

	return "workflow_node"
}

func init() {
	orm.RegisterModel(new(WorkflowNode))
}

func GetWorkflowNodeByDepartmentId(id, step int) (v *WorkflowNode, err error) {
	o := orm.NewOrm()
	v = &WorkflowNode{DepartmentId: id, Step: step}
	if err = o.Read(v, "step"); err == nil {
		return v, nil
	}
	return nil, err
}

func GetWorkflowNodeByIds(ids string) (taskName string, err error) {
	o := orm.NewOrm()
	idArray := strings.Split(ids, ",")
	for _, v := range idArray {
		if len(v) > 0 {
			id, errs := strconv.Atoi(v)
			if errs != nil {
				err = errs
				break
			}
			var info WorkflowNode
			info.Id = id
			if err = o.Read(&info); err != nil {
				return
			}
			taskName += info.NodeName
			return
		} else {
			continue
		}
	}
	return
}

func GetWorkflowNodeByUid(uid, step int) (v *WorkflowNode, err error) {
	var DepartmentId int
	err, DepartmentId = models.GetDepartmentIdById(uid)
	if err != nil {
		return
	}
	return GetWorkflowNodeByDepartmentId(DepartmentId, step)
}

func GetWorlflowBodeDepartmentIdsById(nodeid string, departmentId string) (err error) {
	nodes := strings.Split(nodeid, ",")
	for _, v := range nodes {
		if len(v) > 0 {
			o := orm.NewOrm()
			var index int64
			//index, err = o.QueryTable("workflow_node").Filter("id", nodeid).Filter("department_id", departmentId).Count()
			index, err = o.QueryTable("workflow_node").Filter("id", nodeid).Count()
			if err != nil || index <= 0 {
				err = errors.New("该用户当前不能操作")
				continue
			}
		}
	}
	return
}

func GetNextUsersByUid(id, status, step, WfNameId int) (infos interface{}, err error) {
	//var DepartmentId int
	//err, DepartmentId = models.GetDepartmentIdById(id)
	if err != nil {
		return
	}
	o := orm.NewOrm()
	var node WorkflowNode
	//node.DepartmentId = DepartmentId
	node.WfNameId = WfNameId
	node.Step = step

	if err = o.Read(&node, "wf_name_id", "step"); err != nil {
		err = errors.New("未查询到部门信息或者该用户没此工作流程")
		return
	}
	var nodeIds []string
	if status == 0 {
		nodeIds = strings.Split(node.NextNodeId, ",")
	} else if status == 1 {
		nodeIds = strings.Split(strconv.Itoa(node.RollbackId), ",")
	} else if status == 2 {
		nodeIds = strings.Split(strconv.Itoa(node.DepartmentId), ",")
	}
	if len(nodeIds) <= 0 {
		err = errors.New("未查询到下一步数据")
		return
	}
	var outWork []OutWorkflowNode
	for _, v := range nodeIds {
		id, errs := strconv.Atoi(v)
		if errs != nil {
			err = errs
			break
		}
		if id != 0 {
			var node WorkflowNode
			node.Id = id
			if err = o.Read(&node); err != nil {
				err = errors.New("获取结点信息失败")
				break
			}
			var department models.Department
			department.Id = node.DepartmentId
			if err = o.Read(&department); err != nil {
				err = errors.New("获取部门信息失败")
				break
			}
			var user []models.User
			i, errs := o.QueryTable("user").Filter("department_id", node.DepartmentId).All(&user)
			if i <= 0 || errs != nil {
				err = errors.New("查询数据失败")
				break
			}
			outWork = append(outWork, OutWorkflowNode{department.Name, &user})
		}
	}
	infos = outWork
	return
}

func getNextNodeStep(nextNodeId string) (nextStep int) {

	nextnode, _ := strconv.Atoi(nextNodeId)
	o := orm.NewOrm()
	node := WorkflowNode{Id: nextnode}
	o.Read(&node, "Id")
	nextStep = node.Step
	return
}
