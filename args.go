package metadata

import "reflect"

/********************************************************************
created:    2020-06-20
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var emptyFilter = func(v interface{}) bool {
	return true
}

type Args struct {
	SheetName string                   // 如果是空字符串""，则直接使用反射的类型
	Filter    func(v interface{}) bool // 默认为nil，此时使用emptyFilter
}

func (args *Args) complement(metaType reflect.Type) {
	if args.SheetName == "" {
		args.SheetName = metaType.Name()
	}

	if args.Filter == nil {
		args.Filter = emptyFilter
	}
}
