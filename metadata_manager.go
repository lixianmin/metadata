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
			args.FilePath = localPath
			my.addLocalExcel(args)
		})
	} else {
		my.addLocalExcel(args)
	}
}

func (my *MetadataManager) addLocalExcel(args ExcelArgs) {
	var sheetNames = loadSheetNames(args.FilePath)
	for _, name := range sheetNames {
		my.routeTable.Store(name, args)
	}

	atomic.StorePointer(&my.templateManager, unsafe.Pointer(newTemplateManager()))
	atomic.StorePointer(&my.configManager, unsafe.Pointer(newConfigManager()))
	atomic.AddInt32(&my.excelCount, 1)

	logger.Warn("Excel file is added, args=%v", args)
	if args.OnAdded != nil {
		args.OnAdded(args.FilePath)
	}
}

func (my *MetadataManager) GetTemplate(pTemplate interface{}, id interface{}, sheetName ...string) bool {
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	if manager == nil || id == nil {
		return false
	}

	return manager.getTemplate(&my.routeTable, pTemplate, id, sheetName...)
}

func (my *MetadataManager) GetTemplates(pTemplateList interface{}, args ...Args) bool {
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	if manager == nil {
		return false
	}

	return manager.getTemplates(&my.routeTable, pTemplateList, args...)
}

func (my *MetadataManager) GetConfig(pConfig interface{}, sheetName ...string) bool {
	var manager = (*ConfigManager)(atomic.LoadPointer(&my.configManager))
	return manager != nil && manager.GetConfig(&my.routeTable, pConfig, sheetName...)
}

func (my *MetadataManager) GetExcelCount() int {
	var count = atomic.LoadInt32(&my.excelCount)
	return int(count)
}
