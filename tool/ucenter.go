package tool

import (
	"kuaifa.com/kuaifa/work-together/uc_client"
	"kuaifa.com/kuaifa/work-together/models"
)

func RegisterUser(name, pwd, mail, phone, nickname string) (user models.User, err error) {
	user.Name, user.Email, err = uc_client.RegisterUser(name, pwd, mail, phone, nickname);
	if err != nil {
		return
	}

	return
}

func GetUserInfoByName(name string) (userInfo *uc_client.UserInfo, err error) {
	return uc_client.GetUserInfoByName(name, "")
}

func ChangeOwnPwd(name, pwd, newpwd, atoken string) error {
	return uc_client.ChangeOwnPwd(name, pwd, newpwd, atoken)
}

func UpdateOwnInfo(name, phone, extra, atoken string) error {
	return uc_client.UpdateUserInfo(name, phone, extra, atoken)
}

func DeleteUser(name string) error {
	return uc_client.DeleteUser(name)
}

func ChangeUserPwd(name, newpwd string) error {
	return uc_client.ChangeUserPwd(name, newpwd)
}
