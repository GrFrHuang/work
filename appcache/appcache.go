// Copyright © 2016年 kuaifazs.com. All rights reserved.
// 
// @Author: wuchengshuang@kuaifzs.com
// @Date: 2017/1/16
// @Version: -
// @Desc: -

package appcache

import (
	"errors"
	"fmt"
	"kuaifa.com/kuaifa/work-together/models"
	"time"
)

//获取今日总流水
func GetDashboardBasicInfoCache(basictype string) (*models.DashboardBasicInfo, error) {
	var basicinfo *models.DashboardBasicInfo
	var today = time.Now().Format("2006-01-02")
	cachekey := fmt.Sprintf("Cache_basic_total_%s_%s", basictype, today)
	cacheinfo := Bm.Get(cachekey + "00")
	if cacheinfo == nil {
		var err error
		switch basictype {
		case "todayMoney":
			basicinfo, err = models.TodayTotal(today)
			break
		case "notReconMoney":
			basicinfo, err = models.NotReconTotal()
			break
		case "notRemitMoney":
			basicinfo, err = models.NotRemitTotal()
			break
		case "notContractCount":
			basicinfo, err = models.NotContractTotal()
			break
		}
		if err != nil {
			return basicinfo, errors.New("没有数据")
		}
		Bm.Put(cachekey, basicinfo, 600*time.Second)
	} else {
		basicinfo = cacheinfo.(*models.DashboardBasicInfo)
	}
	return basicinfo, nil
}
