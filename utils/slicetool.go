package utils

import "reflect"

/** 
* Created by stydm on 2019/8/22. 
*/

//给任意切片类型插入元素
func SliceInsert(slice interface{}, index int, value interface{}) (interface{}, bool) {
	//判断是否是切片类型
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}

	//参数检查
	if index < 0 || index > v.Len() || reflect.TypeOf(slice).Elem() != reflect.TypeOf(value) {
		return nil, false
	}

	//尾部追加元素
	if index == v.Len() {
		return reflect.Append(v, reflect.ValueOf(value)).Interface(), true
	}

	v = reflect.AppendSlice(v.Slice(0, index + 1), v.Slice(index, v.Len()))
	v.Index(index).Set(reflect.ValueOf(value))
	return v.Interface(), true
}

//删除任意切片类型指定下标的元素
func SliceDelete(slice interface{}, index int) (interface{}, bool) {
	//判断是否是切片类型
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}
	//参数检查
	if v.Len() == 0 || index < 0 || index > v.Len() - 1 {
		return nil, false
	}

	return reflect.AppendSlice(v.Slice(0, index), v.Slice(index+1, v.Len())).Interface(), true
}


//修改任意切片类型指定下标的元素
func SliceUpdate(slice interface{}, index int, value interface{}) (interface{}, bool) {
	//判断是否是切片类型
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}

	//参数检查
	if index > v.Len() - 1 || reflect.TypeOf(slice).Elem() != reflect.TypeOf(value) {
		return nil, false
	}

	v.Index(index).Set(reflect.ValueOf(value))

	return v.Interface(), true
}

//查找指定元素在任意切片类型中的所有下标
func SliceSearch(slice interface{}, value interface{}) ([]int, bool) {
	//判断是否是切片类型
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}

	var index []int
	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == reflect.ValueOf(value).Interface() {
			index = append(index, i)
		}
	}
	return index, true
}

