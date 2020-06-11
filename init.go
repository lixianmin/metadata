package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var templateManager unsafe.Pointer
var configManager unsafe.Pointer
var lock sync.Mutex

func Init(log logger.ILogger, excelFilePath string) {
	logger.Init(log)

	var isUrl = strings.HasPrefix(excelFilePath, "http://") || strings.HasPrefix(excelFilePath, "https://")
	if isUrl {
		var web = NewWebFile(excelFilePath)
		web.Start(func(filepath string) {
			logger.Warn("Metadata file is changed, excelFilePath=%q, filepath=%q", excelFilePath, filepath)
			atomic.StorePointer(&templateManager, unsafe.Pointer(newTemplateManager(filepath)))
			atomic.StorePointer(&configManager, unsafe.Pointer(newConfigManager(filepath)))
		})
	} else {
		logger.Warn("Metadata file is changed, excelFilePath=%q", excelFilePath)
		atomic.StorePointer(&templateManager, unsafe.Pointer(newTemplateManager(excelFilePath)))
		atomic.StorePointer(&configManager, unsafe.Pointer(newConfigManager(excelFilePath)))
	}
}

func GetTemplate(id int, pTemplate interface{}) bool {
	var manager = (*TemplateManager)(atomic.LoadPointer(&templateManager))
	return manager != nil && manager.GetTemplate(id, pTemplate)
}

func GetTemplates(pTemplateList interface{}) bool {
	var manager = (*TemplateManager)(atomic.LoadPointer(&templateManager))
	return manager != nil && manager.GetTemplates(pTemplateList)
}

func GetConfig(pConfig interface{}) bool {
	var manager = (*ConfigManager)(atomic.LoadPointer(&configManager))
	return manager != nil && manager.GetConfig(pConfig)
}
