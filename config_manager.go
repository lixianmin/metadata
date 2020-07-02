package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"github.com/lixianmin/metadata/tools"
	"github.com/szyhf/go-excel"
	"reflect"
	"sync"
)

/********************************************************************
created:    2020-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type ConfigManager struct {
	configs sync.Map
}

func newConfigManager() *ConfigManager {
	var manager = &ConfigManager{

	}

	return manager
}

func (manager *ConfigManager) GetConfig(routeTable *sync.Map, pConfig interface{}, sheetName string) bool {
	if tools.IsNil(pConfig) {
		logger.Error("pConfig is nil")
		return false
	}

	var pConfigValue = reflect.ValueOf(pConfig)
	if pConfigValue.Kind() != reflect.Ptr {
		logger.Error("pConfig should be of type *Config")
		return false
	}

	var configValue = pConfigValue.Elem()
	var configType = configValue.Type()

	// 获取sheetName
	var sheetName2 = sheetName
	if sheetName2 == "" {
		sheetName2 = configType.Name()
	}

	var config, ok = manager.configs.Load(sheetName2)
	if ok {
		return checkSetValue(configValue, config)
	}

	excelArgs, ok := routeTable.Load(sheetName2)
	if !ok {
		logger.Error("Can not find excelFilePath for sheetName=%q", sheetName2)
		return false
	}

	var err = manager.loadConfig(excelArgs.(ExcelArgs), configType, sheetName2)
	if err != nil {
		manager.configs.Store(sheetName2, nil)
		return false
	}

	config, ok = manager.configs.Load(sheetName2)
	return checkSetValue(configValue, config)
}

func (manager *ConfigManager) loadConfig(args ExcelArgs, configType reflect.Type, sheetName string) error {
	return loadOneSheet(args, sheetName, func(reader excel.Reader) error {
		// double check，如果已经被其它协程加载过了，则不再重复加载
		if _, ok := manager.configs.Load(sheetName); ok {
			return nil
		}

		var pConfigValue = reflect.New(configType)
		var pConfig = pConfigValue.Interface()
		var err = reader.Read(pConfig)
		if err != nil {
			return logger.Dot(err)
		}

		var config = pConfigValue.Elem().Interface()
		manager.configs.Store(sheetName, config)
		return nil
	})
}
