package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"strings"
	"sync"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var templateManager *TemplateManager
var configManager *ConfigManager
var lock sync.Mutex

func Init(log logger.ILogger, excelFilePath string) {
	logger.Init(log)

	var isUrl = strings.HasPrefix(excelFilePath, "http://") || strings.HasPrefix(excelFilePath, "http://")
	if isUrl {
		var web = NewWebFile(excelFilePath)
		web.Start(func(filepath string) {
			templateManager = newTemplateManager(filepath)
			configManager = newConfigManager(filepath)
		})
	} else {
		templateManager = newTemplateManager(excelFilePath)
		configManager = newConfigManager(excelFilePath)
	}
}

func GetTemplate(id int, pTemplate interface{}) bool {
	return templateManager != nil && templateManager.GetTemplate(id, pTemplate)
}

func GetTemplates(pTemplateList interface{}) bool {
	return templateManager != nil && templateManager.GetTemplates(pTemplateList)
}

func GetConfig(pConfig interface{}) bool {
	return configManager != nil && configManager.GetConfig(pConfig)
}
