package metadata

import (
	"github.com/lixianmin/metadata/logger"
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
	excelFilePath string
	configs       sync.Map
}

func newConfigManager(excelFilePath string) *ConfigManager {
	var manager = &ConfigManager{
		excelFilePath: excelFilePath,
	}

	return manager
}

func (manager *ConfigManager) GetConfig(pConfig interface{}) bool {
	if IsNil(pConfig) {
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
	var configName = configType.Name()
	var config, ok = manager.configs.Load(configName)
	if ok {
		return checkSetValue(configValue, config)
	}

	var err = manager.loadConfig(configType, configName)
	if err != nil {
		manager.configs.Store(configName, nil)
		return false
	}

	config, ok = manager.configs.Load(configName)
	return checkSetValue(configValue, config)
}

func (manager *ConfigManager) loadConfig(configType reflect.Type, sheetName string) error {
	return loadOneSheet(manager.excelFilePath, sheetName, func(reader excel.Reader) error {
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
