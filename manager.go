package metadata

import (
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/metadata/logger"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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

	watchLocalFileOnce sync.Once
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

		my.watchLocalFileOnce.Do(func() {
			loom.Go(my.goWatchLocalExcel)
		})
	}
}

func (my *Manager) goWatchLocalExcel(later loom.Later) {
	var ticker = later.NewTicker(5 * time.Second)

	type FileInfo struct {
		Size    int64
		ModTime time.Time
	}
	var lastInfos = make(map[string]FileInfo)

	for {
		select {
		case <-ticker.C:
			my.excelFiles.Range(func(key interface{}, value interface{}) {
				var rawFilePath, args = key.(string), value.(ExcelArgs)
				var info, err = os.Stat(rawFilePath)
				if err != nil {
					return
				}

				if last, ok := lastInfos[rawFilePath]; !ok || !last.ModTime.Equal(info.ModTime()) || last.Size != info.Size() {
					lastInfos[rawFilePath] = FileInfo{
						Size:    info.Size(),
						ModTime: info.ModTime(),
					}

					// 如果ok，则意味着文件变了
					if ok {
						my.addLocalExcel(rawFilePath, args)
					}
				}
			})
		}
	}
}

// 本来使用 github.com/fsnotify/fsnotify ，测试了三次，结果分别是：op=Remove, op=Chmod, 以及没有收到event
//func (my *Manager) goWatchLocalExcel(rawFilePath string, args ExcelArgs) {
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		logo.JsonE("err", err)
//	}
//	defer watcher.Close()
//
//	err = watcher.Add(rawFilePath)
//	if err != nil {
//		logo.JsonE("err", err)
//	}
//
//	for {
//		select {
//		case event, ok := <-watcher.Events:
//			if !ok {
//				return
//			}
//
//			if event.Op&fsnotify.Write == fsnotify.Write {
//				my.addLocalExcel(rawFilePath, args)
//			}
//			logo.JsonI("event", event)
//		case err, ok := <-watcher.Errors:
//			if !ok {
//				return
//			}
//
//			logo.JsonW("err", err)
//		}
//	}
//}

func (my *Manager) addLocalExcel(rawFilePath string, args ExcelArgs) {
	var sheetNames = loadSheetNames(args.FilePath)
	for _, name := range sheetNames {
		my.routeTable.Store(name, args)
	}

	atomic.StorePointer(&my.templateManager, unsafe.Pointer(newTemplateManager()))
	atomic.StorePointer(&my.configManager, unsafe.Pointer(newConfigManager()))
	my.excelFiles.Put(rawFilePath, args)

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
