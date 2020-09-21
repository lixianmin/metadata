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

type options struct {
	SheetName string                   // 如果是空字符串""，则直接使用反射的类型
	Filter    func(v interface{}) bool // 默认为nil，此时使用emptyFilter
}

type Option func(*options)

func WithSheetName(name string) Option {
	return func(opt *options) {
		opt.SheetName = name
	}
}

func WithFilter(filter func(interface{}) bool) Option {
	return func(opt *options) {
		opt.Filter = filter
	}
}

func (my *options) complement(metaType reflect.Type) {
	if my.SheetName == "" {
		my.SheetName = metaType.Name()
	}

	if my.Filter == nil {
		my.Filter = emptyFilter
	}
}
