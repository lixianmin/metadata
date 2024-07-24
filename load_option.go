package metadata

import "reflect"

/********************************************************************
created:    2020-06-20
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type loadOptions struct {
	SheetName string           // 如果是空字符串""，则直接使用反射的类型
	Filter    func(v any) bool // 默认为nil，此时使用emptyFilter
}

type LoadOption func(*loadOptions)

func WithSheet(sheetName string) LoadOption {
	return func(opt *loadOptions) {
		opt.SheetName = sheetName
	}
}

func WithFilter(filter func(any) bool) LoadOption {
	return func(opt *loadOptions) {
		opt.Filter = filter
	}
}

func (my *loadOptions) complement(metaType reflect.Type) {
	if my.SheetName == "" {
		my.SheetName = metaType.Name()
	}

	//if my.Filter == nil {
	//	my.Filter = emptyFilter
	//}
}
