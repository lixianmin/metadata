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
	excelFiles      sync.Map // 原来使用的是loom.Map，在Range()的时候调用了Put()，死锁了
	onExcelChanged  delegateString

	watchLocalFileOnce sync.Once
}

func (my *Manager) AddExcel(opts ...ExcelOption) {
	// 默认值
	var options = excelOptions{}

	// 初始化
	for _, opt := range opts {
		opt(&options)
	}

	var rawFilePath = options.Uri
	var isUrl = strings.HasPrefix(rawFilePath, "http://") || strings.HasPrefix(rawFilePath, "https://")
	if isUrl {
		var web = NewWebFile(rawFilePath)
		web.Start(func(localPath string) {
			options.Uri = localPath
			my.addLocalExcel(rawFilePath, options)
		})
	} else {
		my.addLocalExcel(rawFilePath, options)

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
			my.excelFiles.Range(func(key, value interface{}) bool {
				var rawFilePath, args = key.(string), value.(excelOptions)
				var info, err = os.Stat(rawFilePath)
				if err != nil {
					return true
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

				return true
			})
		}
	}
}

func (my *Manager) addLocalExcel(rawFilePath string, args excelOptions) {
	var sheetNames = loadSheetNames(args.Uri)
	for _, name := range sheetNames {
		my.routeTable.Store(name, args)
	}

	atomic.StorePointer(&my.templateManager, unsafe.Pointer(newTemplateManager()))
	atomic.StorePointer(&my.configManager, unsafe.Pointer(newConfigManager()))
	my.excelFiles.Store(rawFilePath, args)

	logger.Info("Excel file is added, args=%v, excelCount=%d", args, my.GetExcelCount())
	my.onExcelChanged.Invoke(args.Uri)
}

func (my *Manager) OnExcelChanged(handler func(excelFilePath string)) {
	my.onExcelChanged.Add(handler)
}

func (my *Manager) GetTemplate(pTemplate any, id any, opts ...LoadOption) bool {
	var options = my.createLoadOptions(opts)
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && id != nil && manager.getTemplate(&my.routeTable, pTemplate, id, options.SheetName)
}

// GetTemplates 相同的参数每次返回的pTemplateList中的items的不保证顺序：这个是跟实现相关的，目前遍历基于map是不稳定的
func (my *Manager) GetTemplates(pTemplateList any, opts ...LoadOption) bool {
	var options = my.createLoadOptions(opts)
	var manager = (*TemplateManager)(atomic.LoadPointer(&my.templateManager))
	return manager != nil && manager.getTemplates(&my.routeTable, pTemplateList, options)
}

func (my *Manager) GetConfig(pConfig any, opts ...LoadOption) bool {
	var args = my.createLoadOptions(opts)
	var manager = (*ConfigManager)(atomic.LoadPointer(&my.configManager))
	return manager != nil && manager.getConfig(&my.routeTable, pConfig, args.SheetName)
}

func (my *Manager) createLoadOptions(opts []LoadOption) loadOptions {
	var args loadOptions
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
	var count = 0
	my.excelFiles.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	return count
}

// 用于fast fail 的判断
func (my *Manager) isSheetName(name string) bool {
	var _, ok = my.routeTable.Load(name)
	return ok
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
