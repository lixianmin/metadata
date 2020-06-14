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

func GetTemplate(id interface{}, pTemplate interface{}) bool {
	return metadataManager.GetTemplate(id, pTemplate)
}

func GetTemplates(pTemplateList interface{}) bool {
	return metadataManager.GetTemplates(pTemplateList)
}

func GetConfig(pConfig interface{}) bool {
	return metadataManager.GetTemplates(pConfig)
}
