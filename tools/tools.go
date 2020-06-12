package tools

import (
	"os"
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

func EnsureDir(dirname string) error {
	if _, err := os.Stat(dirname); err != nil {
		err = os.MkdirAll(dirname, os.ModePerm)
		return err
	}

	return nil
}
