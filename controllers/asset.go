package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"image"
	"kuaifa.com/kuaifa/work-together/models"
	"kuaifa.com/kuaifa/work-together/models/bean"
	"net/http"
	"strconv"
)

type AssetController struct {
	BaseController
}

// @Title 根据id获取资源文件
// @Description get Game by id
// @Param	id		path 	string	true		"file id"
// @Success 200 {object} models.Asset
// @Failure 403 :id is empty
// @router /:id [get]
func (c *AssetController) Get() {
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

// @Title 根据id获取资源文件名称
// @Description get Game by id
// @Param	id		path 	string	true		"file id"
// @Success 200
// @Failure 403 :id is empty
// @router /name/:id [get]
func (c *AssetController) GetName() {
	mAssert := models.NewAsset()
	idInt, _ := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	id, _ := c.GetInt64("id", idInt)
	if id == 0 {
		c.RespJSON(bean.CODE_Params_Err, "404")
		return
	}

	mAssert.Id = id
	mAssert.Read("Id")

	c.RespJSONData(mAssert.Name)
}

// Post ...
// @Title Post
// @Description upload file
// @Param	asset_name	form	file	true	"文件"
// @Success 201 {int} models.Asset
// @Failure 403 body is empty
// @router /upload [post]
func (c *AssetController) Post() {
	//图片限制尺寸为90x90,格式为png或者jpg
	//c.commonUpload(image.Point{90, 90},
	//暂不限制
	c.commonUpload(image.Point{},
		[]string{
			"png",
			"jpeg",
			"gif",
			"jpg",
			"docx",
			"pdf",
			"xlsx",
			"csv",
			"doc",
			"xls",
		})
}
func (c *AssetController) commonUpload(expectSize image.Point, expFormat []string) {
	f, h, err := c.GetFile("asset_name")
	defer f.Close()
	if err != nil {
		beego.Warn(err)
		c.RespJSON(bean.CODE_Params_Err, err.Error())
	} else {
		mAssert := models.NewAsset()
		m, err := mAssert.SaveAsset(f, h, c.Uid())
		if err != nil {
			beego.Warn(err)
			c.RespJSON(bean.CODE_Bad_Request, err.Error())
		} else {
			//图片限制尺寸或者格式
			retB, err := mAssert.CheckFileType(m.ExtType, expFormat)
			if retB {
				c.RespJSONData(*m)
			} else {
				c.RespJSON(bean.CODE_Params_Err, err.Error())
			}
		}
	}
}

func (c *AssetController) commonUploadPic(expectSize image.Point, expFormat []string) {
	f, h, err := c.GetFile("asset_name")
	defer f.Close()
	if err != nil {
		beego.Warn(err)
		c.RespJSON(bean.CODE_Params_Err, err.Error())
	} else {
		mAssert := models.NewAsset()
		m, err := mAssert.SaveAsset(f, h, c.Uid())
		if err != nil {
			beego.Warn(err)
			c.RespJSON(bean.CODE_Bad_Request, err.Error())
		} else {
			// 图片格式
			retB, err := mAssert.CheckImage(image.Point{}, expFormat)
			if retB {
				c.RespJSONData(*m)
			} else {
				c.RespJSON(bean.CODE_Params_Err, err.Error())
			}
		}
	}
}

// Post ...
// @Title Post
// @Description upload file
// @Param	asset_name	form	file	true	"文件"
// @Success 201 {int} models.Asset
// @Failure 403 body is empty
// @router /uploadFacePic [post]
func (c *AssetController) UploadFacePic() {
	c.commonUploadPic(image.Point{},
		[]string{
			"png",
			"jpeg",
			"gif",
			"jpg",
		})
}
