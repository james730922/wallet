package tools

import (
	"reflect"
)

func SliceBatch(eachCount int, slice interface{}) []interface{} {

	refVal := reflect.ValueOf(slice)

	if refVal.Kind() != reflect.Slice {
		return []interface{}{}
	}

	count := refVal.Len()
	num := count / eachCount

	var result []interface{}
	if count%eachCount == 0 {
		result = make([]interface{}, num)
	} else {
		result = make([]interface{}, num+1)
	}

	for i := 0; i < num; i++ {
		s := i * eachCount
		e := s + eachCount
		result[i] = refVal.Slice(s, e).Interface()
	}
	if count%eachCount > 0 {
		result[num] = refVal.Slice(num*eachCount, count).Interface()
	}

	return result
}

func IntInSlice(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func Int64InSlice(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func StringInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
