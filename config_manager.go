package metadata

import (
	"github.com/lixianmin/logo"
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
	var manager = &ConfigManager{}

	return manager
}

func (manager *ConfigManager) getConfig(routeTable *sync.Map, pConfig any, sheetName string) bool {
	if tools.IsNil(pConfig) {
		logo.Error("pConfig is nil")
		return false
	}

	var pConfigValue = reflect.ValueOf(pConfig)
	if pConfigValue.Kind() != reflect.Ptr {
		logo.Error("pConfig should be of type *Config")
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

	option, ok := routeTable.Load(sheetName2)
	if !ok {
		logo.Error("Can not find excelFilePath for sheetName=%q", sheetName2)
		return false
	}

	var err = manager.loadConfig(option.(excelOptions), configType, sheetName2)
	if err != nil {
		manager.configs.Store(sheetName2, nil)
		return false
	}

	config, ok = manager.configs.Load(sheetName2)
	return checkSetValue(configValue, config)
}

func (manager *ConfigManager) loadConfig(options excelOptions, configType reflect.Type, sheetName string) error {
	return loadOneSheet(options, sheetName, func(reader excel.Reader) error {
		// double check，如果已经被其它协程加载过了，则不再重复加载
		if _, ok := manager.configs.Load(sheetName); ok {
			return nil
		}

		var pConfigValue = reflect.New(configType)
		var pConfig = pConfigValue.Interface()
		var err = reader.Read(pConfig)
		if err != nil {
			logo.JsonW("sheet", sheetName, "err", err)
			return err
		}

		var config = pConfigValue.Elem().Interface()
		manager.configs.Store(sheetName, config)
		return nil
	})
}
