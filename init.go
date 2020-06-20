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

var metadataManager = &MetadataManager{}

func Init(args InitArgs) {
	logger.Init(args.Logger)

	// 每次项目启动时，删除旧的下载文件
	_ = os.RemoveAll(downloadDirectory)
}

func AddExcel(args ExcelArgs) {
	metadataManager.AddExcel(args)
}

func GetTemplate(pTemplate interface{}, args Args) bool {
	return metadataManager.GetTemplate(pTemplate, args)
}

func GetConfig(pConfig interface{}, args Args) bool {
	return metadataManager.GetConfig(pConfig, args)
}
