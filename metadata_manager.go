package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"github.com/szyhf/go-excel"
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
	routeTable      sync.Map
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

func (my *MetadataManager) onAddNewExcel(excelFilePath string) {
	// 互斥加载excel文件
	lock.Lock()
	defer lock.Unlock()

	conn := excel.NewConnecter()
	err := conn.Open(excelFilePath)
	if err != nil {
		logger.Error(err)
		return
	}

	defer conn.Close()

	var sheetNames = conn.GetSheetNames()
	for _, name := range sheetNames {
		my.routeTable.Store(name, excelFilePath)
	}

	atomic.StorePointer(&my.templateManager, unsafe.Pointer(&TemplateManager{}))
	atomic.StorePointer(&my.configManager, unsafe.Pointer(&ConfigManager{}))
}

func (my *MetadataManager) GetTemplate(id int, pTemplate interface{}) bool {
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
