package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

/********************************************************************
created:    2020-06-11
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type MetadataManager struct {
	templateManager unsafe.Pointer
	configManager   unsafe.Pointer
	excelCount      int32    // 成功加载的excel文件个数，用于判断初始化完成
	routeTable      sync.Map // 路由表，sheetName => excelFilePath
}

func (my *MetadataManager) AddExcel(remotePath string) {
	var isUrl = strings.HasPrefix(remotePath, "http://") || strings.HasPrefix(remotePath, "https://")
	if isUrl {
		var web = NewWebFile(remotePath)
		web.Start(func(localPath string) {
			logger.Warn("Metadata file is changed, remotePath=%q, localPath=%q", remotePath, localPath)
			my.onAddNewExcel(localPath)
		})
	} else {
		logger.Warn("Metadata file is changed, remotePath=%q", remotePath)
		my.onAddNewExcel(remotePath)
	}
}

func (my *MetadataManager) onAddNewExcel(localPath string) {
	var sheetNames = loadSheetNames(localPath)
	for _, name := range sheetNames {
		my.routeTable.Store(name, localPath)
	}

	atomic.StorePointer(&my.templateManager, unsafe.Pointer(&TemplateManager{}))
	atomic.StorePointer(&my.configManager, unsafe.Pointer(&ConfigManager{}))
	atomic.AddInt32(&my.excelCount, 1)
}

func (my *MetadataManager) GetTemplate(id interface{}, pTemplate interface{}) bool {
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && manager.GetTemplate(&my.routeTable, id, pTemplate)
}

func (my *MetadataManager) GetTemplates(pTemplateList interface{}) bool {
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && manager.GetTemplates(&my.routeTable, pTemplateList)
}

func (my *MetadataManager) GetConfig(pConfig interface{}) bool {
	var manager = (*ConfigManager)(atomic.LoadPointer(&my.configManager))
	return manager != nil && manager.GetConfig(&my.routeTable, pConfig)
}

func (my *MetadataManager) GetExcelCount() int {
	var count = atomic.LoadInt32(&my.excelCount)
	return int(count)
}
