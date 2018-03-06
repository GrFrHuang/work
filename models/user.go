package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/bysir-zl/bjson"
	"github.com/bysir-zl/bygo/util"
	"kuaifa.com/kuaifa/work-together/utils"
	"strconv"
	"time"
	"strings"
	"github.com/astaxie/beego"
)

type User struct {
	Id           int    `json:"id,omitempty" orm:"column(id);auto"`
	Name         string `json:"name,omitempty" orm:"column(name);null" `
	Nickname     string `json:"nickname,omitempty" orm:"column(nickname);null" valid:"Required"`
	Email        string `json:"email,omitempty" orm:"column(email);null" valid:"Required"`
	RoleIds      string `json:"role_ids,omitempty" orm:"column(role_ids);null" `
	DepartmentId int    `json:"department_id,omitempty" orm:"column(department_id);null"  valid:"Required"`
	Disabled     int    `json:"disabled,omitempty" orm:"column(disabled);null" `
	PicFileId    int    `json:"pic_file_id,omitempty" orm:"column(pic_file_id);null" `
	OpenId       string `json:"open_id,omitempty" orm:"column(open_id);null" `
	AccessToken  string `json:"access_token,omitempty" orm:"column(access_token);size(255);null" json:"-"`
	RefreshToken string `json:"refresh_token,omitempty" orm:"column(refresh_token);size(255);null" json:"-"`
	CreatedTime  int    `json:"created_time,omitempty" orm:"column(created_time);null" `
	UpdatedTime  int    `json:"updated_time,omitempty" orm:"column(updated_time);null" `
	Phone        string ` json:"phone,omitempty" orm:"column(phone);size(11);null"`

	Roles          *[]Role `orm:"-" json:"roles,omitempty"`
	DepartmentName string  `orm:"-" json:"department_name,omitempty"`

	Pwd string `orm:"-" json:"pwd,omitempty"`

	FacePathId string `orm:"-" json:"face_path_id,omitempty"`

	GameIds  string `orm:"-" json:"gameids,omitempty"`  //仅用于记录商务所负责的游戏id（渠道合同页）
	Channels string `orm:"-" json:"channels,omitempty"` //仅用于记录商务所负责的channel_code（渠道合同页）

	EntryTime int    `json:"entry_time,omitempty" orm:"column(entry_time);null" ` // 入职时间
	Position  string `json:"position,omitempty" orm:"column(position);null" `     // 职位
}
type UserSimple struct {
	Value        string `json:"value"` //user email
	Label        string `json:"label"` //用户名字
	Email        string `json:"email"`
	DepartmentId int    `json:"department_id"`
}
type UserByDepartment struct {
	Value    int          `json:"value"` //部门id
	Label    string       `json:"label"` //部门名字
	Children []UserSimple `json:"children"`
}
type GetParam struct {
	ToUser   string `json:"to_user"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	DepId    int    `json:"dep_id"`
	ToGroup  string `json:"to_group"`
	UserType string `json:"user_type"`
}

func (t *User) TableName() string {
	return "user"
}

func init() {
	orm.RegisterModel(new(User))
}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	m.CreatedTime = int(time.Now().Unix())
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetUserInfoByName(name string) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Name: name}
	if err = o.Read(v, "Name"); err == nil {
		return v, nil
	}
	return nil, err
}

func GetUserInfoByNickName(nickName string) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Nickname: nickName}
	if err = o.Read(v, "Nickname"); err == nil {
		return v, nil
	}
	return nil, err
}

func GetUserInfoByToken(token string) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{AccessToken: token}
	if err = o.Read(v, "AccessToken"); err == nil {
		return v, nil
	}
	return nil, err
}
func IsExitUSerById(uid int) (err error) {
	o := orm.NewOrm()
	v := &User{Id: uid}
	if err = o.Read(v); err != nil {
		return
	}
	return
}

func GetDepartmentIdById(uid int) (err error, departmentId int) {
	o := orm.NewOrm()
	v := &User{Id: uid}
	if err = o.Read(v); err != nil {
		return
	}
	if v.DepartmentId == 0 {
		err = errors.New("没有查到部门归属")
		return
	}
	departmentId = v.DepartmentId
	return
}

func GetNickNameById(uid int) (err error, name string) {
	o := orm.NewOrm()
	v := &User{Id: uid}
	if err = o.Read(v); err != nil {
		return
	}
	if len(v.Nickname) <= 0 {
		err = errors.New("没查到用户名")
		return
	}
	name = v.Nickname
	return
}

func GetUserOpenIdById(uid int) string {
	o := orm.NewOrm()
	v := &User{}
	err := o.QueryTable(v).Filter("Id", uid).One(v, "OpenId")
	if err != nil {
		return ""
	}
	return v.OpenId
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.Read(v); err == nil {
		v.AccessToken = ""
		return v, nil
	}
	return nil, err
}

func GroupAddRoleInfo(us []User) bool {
	if len(us) == 0 {
		return false
	}
	rsMap := GetAllRoleMap()
	if rsMap == nil || len(rsMap) == 0 {
		return false
	}
	for i := range us {
		AddRoleInfo(&us[i], rsMap)
	}
	return true
}

func AddRoleInfo(v *User, rsMap map[int]Role) bool {
	if rsMap == nil || len(rsMap) == 0 {
		return false
	}
	bj, err := bjson.New([]byte(v.RoleIds))
	if err != nil {
		return false
	}
	if l := bj.Len(); l != 0 {
		var rrs []Role
		for i := 0; i < l; i++ {
			if r, ok := rsMap[bj.Index(i).Int()]; ok {
				rrs = append(rrs, r)
			}
		}
		v.Roles = &rrs
	}
	return true
}

func GetUsersByIds(ids []interface{}) (user map[int]string, err error) {
	user = map[int]string{}
	if ids == nil || len(ids) == 0 {
		return
	}

	o := orm.NewOrm()
	var users []User
	_, err = o.QueryTable("user").Filter("id__in", ids...).All(&users, "Id", "Nickname")
	if err != nil {
		return
	}
	for _, v := range users {
		user[v.Id] = v.Nickname
	}
	return
}

type UserIdAndName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetUsersListByIds(ids []interface{}) (userIdAndNames []*UserIdAndName, err error) {
	userIdAndNames = []*UserIdAndName{}
	if ids == nil || len(ids) == 0 {
		return
	}

	o := orm.NewOrm()
	var users []User
	_, err = o.QueryTable("user").Filter("id__in", ids...).All(&users, "Id", "Nickname")
	if err != nil {
		return
	}
	for _, v := range users {
		userIdAndNames = append(userIdAndNames, &UserIdAndName{v.Id, v.Nickname})
	}
	return
}

func GroupAddDepartmentInfo(us []User) bool {
	if len(us) == 0 {
		return false
	}
	dsMap := GetAllDepartmentMap()
	if dsMap == nil || len(dsMap) == 0 {
		return false
	}
	for i, v := range us {
		us[i].DepartmentName = dsMap[v.DepartmentId].Name
	}
	return true
}

func UpdateUserToken(m *User) (err error) {
	o := orm.NewOrm()
	f := utils.GetNotEmptyFields(m, "AccessToken", "RefreshToken")
	num, err := o.Update(m, f...)
	if err != nil {
		return
	}
	if num == 0 {
		err = errors.New("not found")
	}
	return
}

// 更新user的基本信息, 不包含token
func UpdateUserDataById(m *User) (err error) {
	o := orm.NewOrm()

	m.UpdatedTime = int(time.Now().Unix())
	fields := utils.GetNotEmptyFields(m, "DepartmentId", "RoleIds", "UpdatedTime", "Nickname", "PicFileId")
	num, err := o.Update(m, fields...)
	if err != nil {
		return
	}
	if num == 0 {
		err = errors.New("not found")
	}
	return
}

// 更新自己的基本信息
func UpdateUserSelfDataById(m *User) (err error) {
	o := orm.NewOrm()

	m.UpdatedTime = int(time.Now().Unix())
	fields := utils.GetNotEmptyFields(m, "DepartmentId", "UpdatedTime", "Nickname", "PicFileId", "Phone")
	num, err := o.Update(m, fields...)
	if err != nil {
		return
	}
	if num == 0 {
		err = errors.New("not found")
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err != nil {
		return
	}
	// 软删除用户
	if _, err = o.Update(&User{Id: id, Name: ""}, "Name"); err != nil {
		return
	}
	return
}

func GetUsersByDevMent(devMentCode int) (us []User, err error) {
	o := orm.NewOrm()
	us = []User{}
	_, err = o.QueryTable(new(User)).
		Filter("department_id", devMentCode).Filter("disabled", 1).
		All(&us, "NickName", "Id", "Name", "Email", "Phone")
	if err != nil {
		return
	}
	return us, nil
}

//func GetDev(userId int) (de Department, err error) {
//	o := orm.NewOrm()
//	us := User{}
//	de = Department{}
//	_, err = o.QueryTable(new(User)).
//		Filter("id", userId).
//		All(&us, "DepartmentId")
//	_, err = o.QueryTable(new(Department)).
//		Filter("id", us.DepartmentId).
//		All(&de, "Name", "Id")
//	if err != nil {
//		return
//	}
//	return de, nil
//}

//func GetUserSimple()([]*UserSimple){
//	o := orm.NewOrm()
//	userSimple := &UserSimple{}
//	userSimples := []*UserSimple{}
//	users := []User{}
//	o.QueryTable(new(User)).All(&users)
//	for _,user :=range users{
//		userSimple.Value = user.Id
//		userSimple.Label = user.Name
//		userSimple.Email = user.Email
//		userSimple.DepartmentId = user.DepartmentId
//		userSimples = append(userSimples,userSimple)
//	}
//	return userSimples
//}

func GetAllUserByDp() ([]UserByDepartment) {
	deps := GetAllDepartment()
	userByDep := UserByDepartment{}
	var userByDeps []UserByDepartment
	for _, dep := range deps {
		userByDep.Value = dep.Id
		userByDep.Label = dep.Name
		users, _ := GetUsersByDevMent(dep.Id)
		userByDep.Children = ChangeToSimple(users)
		userByDeps = append(userByDeps, userByDep)
	}
	return userByDeps
}

//商务负责人附加该负责人负责游戏id，渠道code
func BusinessAddGameIdAndChannelCodeInfo(ss *[]User) {
	if ss == nil || len(*ss) == 0 {
		return
	}
	o := orm.NewOrm()

	for i, s := range *ss {
		var gameids, channelCodes []orm.Params
		var ids, channels []string

		o.Raw("SELECT DISTINCT(game_id) FROM channel_access WHERE business_person=? ", s.Id).Values(&gameids)
		for _, id := range gameids {
			ids = append(ids, id["game_id"].(string))
		}
		(*ss)[i].GameIds = strings.Join(ids, ",")

		o.Raw("SELECT DISTINCT(channel_code) FROM channel_access WHERE business_person=? ", s.Id).Values(&channelCodes)
		for _, channelCode := range channelCodes {
			channels = append(channels, channelCode["channel_code"].(string))
		}
		(*ss)[i].Channels = strings.Join(channels, ",")

	}

	return
}

func ChangeToSimple(users []User) ([]UserSimple) {
	userSimple := UserSimple{}
	var userSimples []UserSimple
	for _, user := range users {
		userSimple.Value = user.Email
		userSimple.Label = user.Nickname
		userSimple.Email = user.Email
		userSimple.DepartmentId = user.DepartmentId
		userSimples = append(userSimples, userSimple)
	}
	return userSimples

}

func GetUserNameByUserId(id int) (name string, err error) {
	table := "userId2Name"
	key := strconv.Itoa(id)
	name, err = utils.Redis.HMGETOne(table, key)
	if err != nil {
		return
	}

	if name != "" {
		return
	}

	var user []User
	_, err = orm.NewOrm().QueryTable("user").All(&user, "id", "nickname")
	if err != nil {
		return
	}

	data := make(map[string]interface{}, len(user))
	for _, v := range user {
		data[strconv.Itoa(v.Id)] = v.Nickname
	}

	name, _ = util.Interface2String(data[key], false)
	if name == "" {
		err = errors.New("404")
		return
	}

	err = utils.Redis.HMSETALL(table, data, 2*60)
	if err != nil {
		return
	}

	return
}

// 清空token,让token不能使用
func EmptyUserToken(id int) {
	orm.NewOrm().Update(&User{Id: id, AccessToken: ""}, "AccessToken")
}

// 快递管理页面
func GetUserList() (users []User, err error) {
	o := orm.NewOrm()
	if _, err = o.QueryTable(new(User)).All(&users, "Id", "NickName"); err != nil {
		return
	}
	return users, nil
}

func GetUserByDepartment(department string) (user User, err error) {
	o := orm.NewOrm()
	var maps []orm.Params
	sql := "select email from user where department_id = ?"
	_, err = o.Raw(sql, department).Values(&maps)
	var email string
	for i := 0; i < len(maps); i++ {
		if i == 0 {
			email = maps[i]["email"].(string)
			continue
		}
		email = email + ";" + maps[i]["email"].(string)
	}
	user.Email = email

	return user, err
}

// Get user name by user ids array
func GetUsersName(ids []int) (names map[int]string, err error) {
	o := orm.NewOrm()
	users := []User{}
	names = make(map[int]string)
	_, err = o.QueryTable("user").Filter("id__in", ids).Limit(len(ids)).All(&users, "id", "nickname")
	if err != nil {
		beego.Error(err)
		return
	}
	for _, v := range users {
		names[v.Id] = v.Nickname
	}
	return
}
