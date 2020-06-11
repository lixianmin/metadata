package metadata

import (
	"github.com/lixianmin/metadata/logger"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var metadataManager = &MetadataManager{}

func Init(log logger.ILogger) {
	logger.Init(log)
}

func AddExcel(remotePath string) {
	metadataManager.AddExcel(remotePath)
}

func GetTemplate(id int, pTemplate interface{}) bool {
	return metadataManager.GetTemplate(id, pTemplate)
}

func GetTemplates(pTemplateList interface{}) bool {
	return metadataManager.GetTemplates(pTemplateList)
}

func GetConfig(pConfig interface{}) bool {
	return metadataManager.GetTemplates(pConfig)
}
