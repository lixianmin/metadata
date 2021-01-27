package metadata

import (
	"github.com/fsnotify/fsnotify"
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/logo"
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
	excelFiles      loom.Map
	onExcelChanged  delegateString
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
		go my.goWatchLocalExcel(rawFilePath, args)
	}
}

func (my *Manager) goWatchLocalExcel(rawFilePath string, args ExcelArgs) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logo.Error("err", err)
	}
	defer watcher.Close()

	err = watcher.Add(rawFilePath)
	if err != nil {
		logo.Error("err", err)
	}

	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				return
			}
			my.addLocalExcel(rawFilePath, args)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			logo.Warn("err", err)
		}
	}
}

func (my *Manager) addLocalExcel(rawFilePath string, args ExcelArgs) {
	var sheetNames = loadSheetNames(args.FilePath)
	for _, name := range sheetNames {
		my.routeTable.Store(name, args)
	}

	atomic.StorePointer(&my.templateManager, unsafe.Pointer(newTemplateManager()))
	atomic.StorePointer(&my.configManager, unsafe.Pointer(newConfigManager()))
	my.excelFiles.Put(rawFilePath, nil)

	logger.Info("Excel file is added, args=%v", args)
	my.onExcelChanged.Invoke(args.FilePath)
}

func (my *Manager) OnExcelChanged(handler func(excelFilePath string)) {
	my.onExcelChanged.Add(handler)
}

func (my *Manager) GetTemplate(pTemplate interface{}, id interface{}, opts ...Option) bool {
	var args = my.createOptions(opts)
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && id != nil && manager.getTemplate(&my.routeTable, pTemplate, id, args.SheetName)
}

// 相同的参数每次返回的pTemplateList中的items的不保证顺序：这个是跟实现相关的，目前遍历基于map是不稳定的
func (my *Manager) GetTemplates(pTemplateList interface{}, opts ...Option) bool {
	var args = my.createOptions(opts)
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && manager.getTemplates(&my.routeTable, pTemplateList, args)
}

func (my *Manager) GetConfig(pConfig interface{}, opts ...Option) bool {
	var args = my.createOptions(opts)
	var manager = (*ConfigManager)(atomic.LoadPointer(&my.configManager))
	return manager != nil && manager.getConfig(&my.routeTable, pConfig, args.SheetName)
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
	return my.excelFiles.Size()
}

// 用于fast fail 的判断
func (my *Manager) isSheetName(name string) bool {
	var _, ok = my.routeTable.Load(name)
	return ok
}
