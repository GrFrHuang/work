package controllers

import (
	"net/http"
	"kuaifa.com/kuaifa/work-together/tool"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"time"
	"image"
	"github.com/astaxie/beego"
	"strconv"
	"fmt"
)

// unused
type DesktopClientController struct {
	BaseController
}

type DesktopResponse struct {
	Code int
	Data interface{}
}

// Post ...
// @Title 桌面客户端登录接口
// @Description 桌面客户端登录接口
// @Param	email		formData 	string	true		"account"
// @Param	password		formData 	string	true		"password"
// @Success 200 ok
// @Failure 403 body is empty
// @router /login [post]
func (c *DesktopClientController) Login() {
	var v bean.CreateSession
	v.Email = c.GetString("email", "")
	v.Password = c.GetString("password", "")
	loginRet, err := tool.CreateSession(v)
	if err != nil {
		c.RespJSON(http.StatusForbidden, err.Error())
		return
	}
	u, err := models.GetUserInfoByName(loginRet.UserName)
	if err != nil {
		c.RespJSON(http.StatusBadRequest, err.Error())
		return
	}

	if u.Disabled == 2 {
		c.RespJSON(http.StatusForbidden, "此账号被禁止登陆")
		return
	}
	err = models.UpdateUserInfo(u.Id, loginRet.AccessToken)
	if err != nil {
		c.RespJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.RespJSON(http.StatusOK, bean.OutPutSession{Uid: u.Id, Token: loginRet.AccessToken})
}


// Post ...
// @Title 截屏图片更新接口
// @Description 截屏图片更新接口
// @Param	token		formData 	string	true		"token"
// @Param	imageurl		formData 	string	true		"imageurl"
// @Success 200 ok
// @Failure 403 body is empty
// @router /screeninfo [post]
func (c *DesktopClientController) Screeninfo() {
	token := c.GetString("token", "")
	imageurl := c.GetString("imageurl", "")
	uid,err := userinfo(token)
	if err != nil {
		c.RespJSON(http.StatusForbidden, "token err")
		return
	}
	data := models.UserDesktopScreenLog{Uid:uid,Imageurl:imageurl,CreateTime:int(time.Now().Unix())}
	_,err = models.AddUserDesktopScreenLog(&data)
	if err != nil {
		c.RespJSON(http.StatusForbidden, err.Error())
		return
	}
	c.RespJSON(http.StatusOK,&DesktopResponse{})
}


// Post ...
// @Title Post
// @Description upload file
// @Param	asset_name	form	file	true	"文件"
// @Success 201 {int} models.Asset
// @Failure 403 body is empty
// @router /uploadPic [post]
func (c *DesktopClientController) UploadPic() {
	token := c.GetString("token", "")
	uid,err := userinfo(token)
	if err != nil {
		c.RespJSON(http.StatusForbidden, "token err")
		return
	}
	c.commonUploadPic(uid,image.Point{},
		[]string{
			"png",
			"jpeg",
			"gif",
			"jpg",
		})
}


// @Title 根据id获取资源文件
// @Description get Game by id
// @Param	id		path 	string	true		"file id"
// @Success 200 {object} models.Asset
// @Failure 403 :id is empty
// @router /image/:id [get]
func (c *DesktopClientController) Get() {
	mAssert := models.NewAsset()
	idInt, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	id, _ := c.GetInt64("id", idInt)
	if id != 0 {
		mAssert.Id = id
		mAssert.Read("Id")
	} else {
		c.RespJSON(bean.CODE_Params_Err, "404")
		return
	}

	var name string
	if mAssert.Name != "" {
		name = mAssert.Name
	} else {
		name = strconv.FormatInt(mAssert.Id, 10) + "." + mAssert.ExtType
	}

	c.Ctx.ResponseWriter.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", name))
	c.Ctx.ResponseWriter.Header().Set("Content-Type", mAssert.MimeType)
	aFile := mAssert.GetFilePath()
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, aFile)
}


func (c *DesktopClientController) commonUploadPic(uid int,expectSize image.Point, expFormat []string) {
	f, h, err := c.GetFile("asset_name")
	defer f.Close()
	if err != nil {
		beego.Warn(err)
		c.RespJSON(bean.CODE_Params_Err, err.Error())
	} else {
		mAssert := models.NewAsset()
		m, err := mAssert.SaveAsset(f, h, uid)
		if err != nil {
			beego.Warn(err)
			c.RespJSON(bean.CODE_Bad_Request, err.Error())
		} else {
			// 图片格式
			retB, err := mAssert.CheckImage(image.Point{}, expFormat)
			if retB {
				type Resdata struct{
					Url string `json:"url"`
				}
				resdata := Resdata{Url: "/dc/desktop/image/"+strconv.Itoa(int(m.Id))}
				c.RespJSONData(&resdata)
			} else {
				c.RespJSON(bean.CODE_Params_Err, err.Error())
			}
		}
	}
}


func userinfo(token string) (uid int, err error){
	user,err := models.GetUserDesktopClientByCtoken(token)
	if err == nil {
		uid = user.Uid
	}
	return
}