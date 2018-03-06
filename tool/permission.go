package tool

import (
	"kuaifa.com/kuaifa/work-together/utils"
	"strings"
)

// 根据前台传入的条件 和 权限控制的权限 生成查询条件
func InjectPermissionWhere(permission map[string][]interface{}, where *map[string][]interface{}) {
	for k, p := range permission {
		_, ok := (*where)[k]
		if !ok {
			// 添加在条件中但是不在权限中的条件
			(*where)[k] = p
		} else {
			// 交集
			if strings.Contains(k, "__in") {
				(*where)[k] = utils.IntersectionSlice((*where)[k], p)
			} else {
				// 直接覆盖
				(*where)[k] = p
			}
		}

	}
}
