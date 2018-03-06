// Copyright © 2016年 kuaifazs.com. All rights reserved.
// 
// @Author: wuchengshuang@kuaifzs.com
// @Date: 2017/1/16
// @Version: -
// @Desc: -

package appcache

import "github.com/astaxie/beego/cache"

var Bm cache.Cache
func init(){
    Bm, _ = cache.NewCache("memory", `{"key":"baole","conn":":6379","dbNum":"0"`)
}