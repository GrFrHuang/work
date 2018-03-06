package workflow

import (
	"github.com/astaxie/beego/orm"
)

type WorkflowName struct {
	Id         int    `orm:"column(id);auto" json:"id,omitempty"`
	TbNickname string `orm:"column(tb_nickname);size(100);null" description:"任务别名" json:"tb_nickname,omitempty"`
	TbName     string `orm:"column(tb_name);size(50);null" description:"工作流表名" json:"-"`
}

func (t *WorkflowName) TableName() string {
	return "workflow_name"
}

func init() {
	orm.RegisterModel(new(WorkflowName))
}
