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
	routeTable      sync.Map // 路由表，sheetName => ExcelArgs

	m          sync.Mutex
	excelFiles map[string]struct{}
	excelCount int // 成功加载的excel文件个数，用于判断初始化完成
}

func (my *Manager) AddExcel(args ExcelArgs) {
	var rawFilePath = args.FilePath
	var isUrl = strings.HasPrefix(rawFilePath, "http://") || strings.HasPrefix(rawFilePath, "https://")
	if isUrl {
		var web = NewWebFile(rawFilePath)
		web.Start(func(localPath string) {
			args.FilePath = localPath
			my.addLocalExcel(rawFilePath, args)
		})
	} else {
		my.addLocalExcel(rawFilePath, args)
	}
}

func (my *Manager) addLocalExcel(rawFilePath string, args ExcelArgs) {
	var sheetNames = loadSheetNames(args.FilePath)
	for _, name := range sheetNames {
		my.routeTable.Store(name, args)
	}

	atomic.StorePointer(&my.templateManager, unsafe.Pointer(newTemplateManager()))
	atomic.StorePointer(&my.configManager, unsafe.Pointer(newConfigManager()))
	my.rememberExcelFiles(rawFilePath)

	logger.Info("Excel file is added, args=%v", args)
	if args.OnAdded != nil {
		args.OnAdded(args.FilePath)
	}
}

func (my *Manager) rememberExcelFiles(rawFilePath string) {
	my.m.Lock()
	if my.excelFiles == nil {
		my.excelFiles = make(map[string]struct{}, 4)
	}

	my.excelFiles[rawFilePath] = struct{}{}
	my.excelCount = len(my.excelFiles)
	my.m.Unlock()
}

func (my *Manager) GetTemplate(pTemplate interface{}, id interface{}, opts ...Option) bool {
	var args = my.createOptions(opts)
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && id != nil && manager.getTemplate(&my.routeTable, pTemplate, id, args.SheetName)
}

func (my *Manager) GetTemplates(pTemplateList interface{}, opts ...Option) bool {
	var args = my.createOptions(opts)
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && manager.getTemplates(&my.routeTable, pTemplateList, args)
}

func (my *Manager) GetConfig(pConfig interface{}, opts ...Option) bool {
	var args = my.createOptions(opts)
	var manager = (*ConfigManager)(atomic.LoadPointer(&my.configManager))
	return manager != nil && manager.GetConfig(&my.routeTable, pConfig, args.SheetName)
}

func (my *Manager) createOptions(opts []Option) options {
	var args options
	for _, opt := range opts {
		opt(&args)
	}

	// 判断sheetName是不合法
	if args.SheetName != "" && !my.isSheetName(args.SheetName) {
		args.SheetName = ""
	}

	return args
}

func (my *Manager) GetExcelCount() int {
	return my.excelCount
}

// 用于fast fail 的判断
func (my *Manager) isSheetName(name string) bool {
	var _, ok = my.routeTable.Load(name)
	return ok
}
