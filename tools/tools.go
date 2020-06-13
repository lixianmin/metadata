package tools

import (
	"reflect"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
