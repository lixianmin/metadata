package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"sync"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var templateManager *TemplateManager
var lock sync.Mutex

func Init(log logger.ILogger, excelFilePath string) {
	logger.Init(log)
	templateManager = newTemplateManager(excelFilePath)
}

func GetTemplate(id int, pTemplate interface{}) bool {
	return templateManager.GetTemplate(id, pTemplate)
}

func GetTemplates(pTemplateList interface{}) bool {
	return templateManager.GetTemplates(pTemplateList)
}
