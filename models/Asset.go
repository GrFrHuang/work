package models

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"strconv"
	"strings"
	"encoding/base64"
	"io/ioutil"
)

var (
	AssetsPath  string
	ErrorFormat error = errors.New("格式不合法")
	ErrorSize   error = errors.New("图片尺寸不合法")
	//存放目录层级
	Level uint = 3
)

func init() {
	orm.RegisterModel(new(Asset))
	AssetsPath = beego.AppPath + "/" + beego.AppConfig.String("assets_save_path")
}

type Asset struct {
	Id      int64  `json:"id,omitempty"  orm:"column(id);null"`
	HashStr string `json:"hash_str,omitempty" orm:"column(hash_str);size(32);unique"`
	//jpg,png,txt,doc...
	ExtType string `json:"ext_type,omitempty" orm:"column(ext_type);size(10)"`
	//image/gif, text/plain
	MimeType   string `json:"mime_type,omitempty" orm:"column(mime_type);size(100)"`
	CreateTime int64  `json:"create_time,omitempty" orm:"column(create_time)"`
	Name       string `json:"name,omitempty" orm:"column(name);size(50)"`
	UserId     int    `json:"user_id,omitempty" orm:"column(user_id)"`
	Base              `orm:"-"`
}

func NewAsset() *Asset {
	mObj := &Asset{}
	mObj.Outer = mObj
	return mObj
}
func (m *Asset) getRowsContainer() interface{} {
	var t interface{}
	u := make([]*Asset, 0)
	t = &u
	return t
}

//保存资源文件到散列路径下
func (m *Asset) SaveAsset(mFile multipart.File,
	fHeader *multipart.FileHeader, uid int) (*Asset, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, mFile); err != nil {
		beego.Warn(err)
		return nil, err
	}
	hashByte := hash.Sum(nil)
	m.HashStr = fmt.Sprintf("%x", hashByte)
	m.Name = fHeader.Filename
	fnameByte := []byte(fHeader.Filename)
	m.ExtType = string(
		fnameByte[bytes.LastIndexByte(fnameByte, '.')+1:])
	m.MimeType = fHeader.Header["Content-Type"][0]
	fPath := m.GetFilePath()
	beego.Debug(fPath)
	mFile.Seek(0, io.SeekStart)
	fDir := filepath.Dir(fPath)
	if _, err := os.Stat(fDir); os.IsNotExist(err) {
		err = os.MkdirAll(fDir, 0755)
		if err != nil {
			beego.Warn(err)
			return nil, err
		}
	}
	nFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE, 0666)
	defer mFile.Close()
	defer nFile.Close()
	if err != nil {
		beego.Warn(err)
		return nil, err
	}
	if _, err := io.Copy(nFile, mFile); err != nil {
		beego.Warn(err)
		return nil, err
	}
	m.CreateTime = time.Now().Unix()
	m.UserId = uid
	beego.Warning(m)
	if id, err := m.ReadOrCreate("HashStr"); err != nil {
		beego.Warn(err, id)
		return nil, err
	}
	return m, nil
}

//保存资源文件到散列路径下
func (m *Asset) SaveAssetBase64Img(img string, uid int) (*Asset, error) {
	filename := strconv.Itoa(uid) + "_" + time.Now().Format("2006-01-02-15-04-05") + ".jpg"
	hash := md5.New()
	if _, err := io.Copy(hash, strings.NewReader(img)); err != nil {
		beego.Warn(err)
		return nil, err
	}
	hashByte := hash.Sum(nil)
	m.HashStr = fmt.Sprintf("%x", hashByte)
	m.Name = filename
	fnameByte := []byte(filename)
	m.ExtType = string(fnameByte[bytes.LastIndexByte(fnameByte, '.')+1:])
	m.MimeType = "image/jpg"
	fPath := m.GetFilePath()
	beego.Debug(fPath)
	strings.NewReader(img).Seek(0, io.SeekStart)
	fDir := filepath.Dir(fPath)
	if _, err := os.Stat(fDir); os.IsNotExist(err) {
		err = os.MkdirAll(fDir, 0755)
		if err != nil {
			beego.Warn(err)
			return nil, err
		}
	}
	d, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(fPath, d, os.ModePerm)
	if err != nil {
		return nil, err
	}
	m.CreateTime = time.Now().Unix()
	m.UserId = uid
	beego.Warning(m)
	if id, err := m.ReadOrCreate("HashStr"); err != nil {
		beego.Warn(err, id)
		return nil, err
	}
	return m, nil
}

//通过关键字生成hash 文件路径
func (m *Asset) GetPathByKeyword(extType string) (string, error) {
	fileAbsPath := AssetsPath + m.ParseHashStr() +
		"." + m.ExtType
	fDir := filepath.Dir(fileAbsPath)
	if _, err := os.Stat(fDir); os.IsNotExist(err) {
		err = os.MkdirAll(fDir, 0755)
		if err != nil {
			beego.Warn(err)
			return "", err
		}
	}
	return fileAbsPath, nil
}

func (m *Asset) GetFilePath() string {
	fileAbsPath := AssetsPath + m.ParseHashStr() +
		"." + m.ExtType
	return fileAbsPath
}

func (m *Asset) ParseHashStr() string {
	hvBytes := []byte(m.HashStr)
	nBytes := []byte("")
	level := Level
	for i, v := range hvBytes {
		if level > 0 && i%2 == 0 {
			nBytes = append(nBytes, '/')
			level--
		}
		nBytes = append(nBytes, v)
	}
	return string(nBytes)
}
func (m *Asset) DelAsset() (bool, error) {
	fPath := m.GetFilePath()
	err := os.Remove(fPath)
	errDb := m.Delete()
	if err != nil || errDb != nil {
		if err != nil {
			beego.Notice(err)
		} else if errDb != nil {
			beego.Notice(errDb)
		}
		return false, err
	} else {
		return true, nil
	}
}
func (m *Asset) CheckImage(expectSize image.Point, expFormat []string) (bool, error) {
	fPath := m.GetFilePath()
	iFile, err := os.Open(fPath)
	defer iFile.Close()
	if err != nil {
		beego.Warn(err)
		return false, err
	}
	img, format, err := image.Decode(iFile)
	allowFormat := false
	for _, iFormat := range expFormat {
		if iFormat == format {
			allowFormat = true
			break
		}
	}
	if allowFormat {
		zeroSize := image.Point{}
		if zeroSize == expectSize || img.Bounds().Size() == expectSize {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("图片大小必须为%dx%d像素", img.Bounds().Size().X, img.Bounds().Size().Y))
	} else {
		return false, ErrorFormat
	}
}

func (m *Asset) CheckFileType(format string, expFormat []string) (bool, error) {
	allowFormat := false
	for _, iFormat := range expFormat {
		if iFormat == format {
			allowFormat = true
			break
		}
	}
	if allowFormat {
		return true, nil
	} else {
		return false, ErrorFormat
	}
}
