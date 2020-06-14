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
	routeTable      sync.Map // 路由表，sheetName => ExcelArgs
}

func (my *MetadataManager) AddExcel(args ExcelArgs) {
	var isUrl = strings.HasPrefix(args.FilePath, "http://") || strings.HasPrefix(args.FilePath, "https://")
	if isUrl {
		var web = NewWebFile(args.FilePath)
		web.Start(func(localPath string) {
			logger.Warn("Metadata file is changed, args=%v, localPath=%q", args, localPath)
			my.onAddNewExcel(ExcelArgs{FilePath: localPath, TitleRowIndex: args.TitleRowIndex})
		})
	} else {
		logger.Warn("Metadata file is changed, args=%v", args)
		my.onAddNewExcel(args)
	}
}

func (my *MetadataManager) onAddNewExcel(args ExcelArgs) {
	var sheetNames = loadSheetNames(args.FilePath)
	for _, name := range sheetNames {
		my.routeTable.Store(name, args)
	}

	atomic.StorePointer(&my.templateManager, unsafe.Pointer(newTemplateManager()))
	atomic.StorePointer(&my.configManager, unsafe.Pointer(newConfigManager()))
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
