package swagger

import "kuaifa.com/kuaifa/work-together/models"

type OutPutUser struct {
	Id           int `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	RoleIds      string `json:"role_ids,omitempty"`
	DepartmentId int `json:"department_id,omitempty"`
	OpenId       int `json:"open_id,omitempty"`
	CreatedTime  int `json:"created_time,omitempty"`
	UpdatedTime  int `json:"updated_time,omitempty"`

	Roles          []models.Role `json:"roles,omitempty"`
	DepartmentName string `json:"department_name,omitempty"`
}

type CreateUser struct {
	Name         string        `json:"name,omitempty"`
	NickName     string        `json:"nick_name,omitempty"`
	RoleIds      string        `json:"role_ids,omitempty"`
	DepartmentId int        `json:"department_id,omitempty"`

	Pwd   string `json:"pwd,omitempty"`
	Mail  string `json:"mail,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type AdminChangePwd struct {
	NewPwd string `json:"new_pwd,omitempty"`
}

type UpdateUserData struct {
	DepartmentId int `json:"department_id,omitempty"`
	RoleIds      string `json:"role_ids,omitempty"`
}
