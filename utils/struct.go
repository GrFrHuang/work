package utils

import (
	"reflect"
	"strconv"
)



func IsEmptyValue(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func Interface2String(value interface{}, strict bool) (v string, ok bool) {
	switch value.(type) {
	case string:
		v, ok = value.(string), true
	case []uint8:
		v, ok = string(value.([]uint8)), true
	}

	if ok {
		return
	}
	if strict {
		if !ok {
			return
		}
	}

	switch value.(type) {
	case int64, int8, int32, int:
		i, _ := Interface2Int(value, true)
		v, ok = strconv.FormatInt(i, 10), true
	case uint64, uint8, uint32, uint:
		i, _ := Interface2UInt(value, true)
		v, ok = strconv.FormatUint(i, 10), true
	case float64, float32:
		f, _ := Interface2Float(value, true)
		v, ok = strconv.FormatFloat(f, 'f', -1, 64), true
	case bool:
		v, ok = strconv.FormatBool(value.(bool)), true
	}
	return
}

func Interface2Float(value interface{}, strict bool) (v float64, ok bool) {
	switch value.(type) {
	case float32:
		v, ok = float64(value.(float32)), true
	case float64:
		v, ok = float64(value.(float64)), true
	}
	if ok {
		return
	}
	if strict {
		if !ok {
			return
		}
	}

	switch value.(type) {
	case int, int8, int32, int64:
		i, ok := Interface2Int(value, true)
		if ok {
			v, ok = float64(i), true
		}
	case uint, uint8, uint32, uint64:
		i, ok := Interface2UInt(value, true)
		if ok {
			v, ok = float64(i), true
		}
	case string:
		f, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	case []uint8:
		f, err := strconv.ParseFloat(string(value.([]uint8)), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return
}

func Interface2UInt(value interface{}, strict bool) (v uint64, ok bool) {
	switch value.(type) {
	case uint:
		v, ok = uint64(value.(uint)), true
	case uint8:
		v, ok = uint64(value.(uint8)), true
	case uint16:
		v, ok = uint64(value.(uint16)), true
	case uint32:
		v, ok = uint64(value.(uint32)), true
	case uint64:
		v, ok = uint64(value.(uint64)), true
	}
	if ok {
		return
	}
	if strict {
		if !ok {
			return
		}
	}

	switch value.(type) {
	case string, []uint8:
		s, _ := Interface2String(value, true)
		i, err := strconv.ParseUint(s, 10, 64)
		if err == nil {
			v, ok = i, true
		}
	case int, int8, int16, int32, int64:
		s, _ := Interface2Int(value, true)
		v, ok = uint64(s), true
	case float32, float64:
		f, _ := Interface2Float(value, true)
		v, ok = uint64(f), true
	case bool:
		if value.(bool) {
			v = 1
		} else {
			v = 0
		}
		ok = true
	}
	return
}

func Interface2Int(value interface{}, strict bool) (v int64, ok bool) {
	switch value.(type) {
	case int:
		v, ok = int64(value.(int)), true
	case int8:
		v, ok = int64(value.(int8)), true
	case int16:
		v, ok = int64(value.(int16)), true
	case int32:
		v, ok = int64(value.(int32)), true
	case int64:
		v, ok = int64(value.(int64)), true
	}
	if ok {
		return
	}
	if strict {
		if !ok {
			return
		}
	}

	switch value.(type) {
	case string, []uint8:
		s, _ := Interface2String(value, true)
		i, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			v, ok = i, true
		}
	case float32, float64:
		f, _ := Interface2Float(value, true)
		v, ok = int64(f), true
	case bool:
		if value.(bool) {
			v = 1
		} else {
			v = 0
		}
		ok = true
	}
	return
}

// 判断item是否在数组里
// 如果数组为空则返回false
func ItemInArray(item string, max []string) (has bool) {
	return ArrayStringIndex(item, max) != -1
}


func ArrayStringIndex(item string, max []string) (index int) {
	index = -1
	if max == nil || len(max) == 0 {
		return
	}
	for i, l := 0, len(max); i < l; i++ {
		if max[i] == item {
			index = i
			return
		}
	}
	return
}


// 获取不为空的在inFields中的 结构体中的字段
func GetNotEmptyFields(obj interface{}, inFields ...string) (fields []string) {
	fields = []string{}
	pointer := reflect.Indirect(reflect.ValueOf(obj))
	types := pointer.Type()
	fieldNum := pointer.NumField()
	for i := 0; i < fieldNum; i++ {
		v := pointer.Field(i)
		name := types.Field(i).Name
		if inFields != nil && len(inFields) != 0 {
			if !ItemInArray(name, inFields) {
				continue
			}
		}
		if IsEmptyValue(v.Interface()) {
			continue
		}
		fields = append(fields, name)
	}
	return
}

//删除指定字段
func RemoveFields(fields []string, delete ...string) (result []string){
	if len(delete) == 0 {
		result = fields
		return
	}

	for _, field := range fields{
		if !ItemInArray(field, delete){
			result = append(result, field)
		}
	}
	return
}

func MergeSlice(s1 []interface{}, s2 []interface{}) []interface{} {
	if s1 == nil {
		if s2 == nil {
			return []interface{}{}
		}
		return s2
	}
	if s2 == nil {
		if s1 == nil {
			return []interface{}{}
		}
		return s1
	}
	var temp []interface{}

	for _, v1 := range s1 {
		inS2 := false
		for _, v2 := range s2 {
			if v1 == v2 {
				inS2 = true
			}
		}
		if !inS2 {
			temp = append(temp, v1)
		}
	}
	temp = append(temp, s2...)
	return temp
}

// 求交集
func IntersectionSlice(s1 []interface{}, s2 []interface{}) []interface{} {
	var temp []interface{}
	for _, v1 := range s1 {
		inS2 := false
		for _, v2 := range s2 {
			if v1 == v2 {
				inS2 = true
			}
		}
		if inS2 {
			temp = append(temp, v1)
		}
	}

	UnDuplicatesSlice(&temp)
	return temp
}

func UnDuplicatesSlice(is *[]interface{}) {
	t := map[interface{}]bool{}
	var temp []interface{}
	for _, i := range *is {
		if t[i] == true {
			continue
		}
		t[i] = true
		temp = append(temp, i)
	}
	*is = temp
}
