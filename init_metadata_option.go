package metadata

import "github.com/lixianmin/logo"

/********************************************************************
created:    2024-07-24
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type initOptions struct {
	Logger logo.ILogger // 自定义日志对象，默认只输出到控制台
}

type InitOption func(*initOptions)

func WithLogger(logger logo.ILogger) InitOption {
	return func(opt *initOptions) {
		opt.Logger = logger
	}
}
