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

type Manager struct {
	templateManager unsafe.Pointer
	configManager   unsafe.Pointer
	excelCount      int32    // 成功加载的excel文件个数，用于判断初始化完成
	routeTable      sync.Map // 路由表，sheetName => ExcelArgs
}

func (my *Manager) AddExcel(args ExcelArgs) {
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

func (my *Manager) addLocalExcel(args ExcelArgs) {
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

func (my *Manager) GetTemplate(pTemplate interface{}, id interface{}, sheetName ...string) bool {
	// 判断是否传入了sheetName
	var sheetName2 = ""
	if len(sheetName) > 0 {
		sheetName2 = sheetName[0]
		if sheetName2 != "" && !my.isSheetName(sheetName2) {
			return false
		}
	}

	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && id != nil && manager.getTemplate(&my.routeTable, pTemplate, id, sheetName2)
}

func (my *Manager) GetTemplates(pTemplateList interface{}, args ...Args) bool {
	// 判断是否传入了合法的sheetName
	var args2 Args
	if len(args) > 0 {
		args2 = args[0]
		if args2.SheetName != "" && !my.isSheetName(args2.SheetName) {
			return false
		}
	}

	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && manager.getTemplates(&my.routeTable, pTemplateList, args2)
}

func (my *Manager) GetConfig(pConfig interface{}, sheetName ...string) bool {
	// 判断是否传入了sheetName
	var sheetName2 = ""
	if len(sheetName) > 0 {
		sheetName2 = sheetName[0]
		if sheetName2 != "" && !my.isSheetName(sheetName2) {
			return false
		}
	}

	var manager = (*ConfigManager)(atomic.LoadPointer(&my.configManager))
	return manager != nil && manager.GetConfig(&my.routeTable, pConfig, sheetName2)
}

func (my *Manager) GetExcelCount() int {
	var count = atomic.LoadInt32(&my.excelCount)
	return int(count)
}

// 用于fast fail 的判断
func (my *Manager) isSheetName(name string) bool {
	var _, ok = my.routeTable.Load(name)
	return ok
}
