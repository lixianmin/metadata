package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"os"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func Init(args InitArgs) {
	logger.Init(args.Logger)

	// 每次项目启动时，删除旧的下载文件
	_ = os.RemoveAll(downloadDirectory)
}
