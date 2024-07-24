package metadata

import (
	"github.com/lixianmin/logo"
	"github.com/lixianmin/metadata/logger"
	"os"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func Init(opts ...InitOption) {
	// 默认值
	var options = initOptions{
		Logger: logo.GetLogger(),
	}

	// 初始化
	for _, opt := range opts {
		opt(&options)
	}

	logger.Init(options.Logger)

	// 每次项目启动时，删除旧的下载文件
	_ = os.RemoveAll(downloadDirectory)
}
