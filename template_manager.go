package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"github.com/lixianmin/metadata/tools"
	"github.com/szyhf/go-excel"
	"reflect"
	"sync"
)

/********************************************************************
created:    2020-06-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type TemplateTable map[interface{}]interface{}

type TemplateManager struct {
	tables sync.Map
}

func newTemplateManager() *TemplateManager {
	var manager = &TemplateManager{

	}

	return manager
}

// template是一个结构体指针
func (manager *TemplateManager) getTemplate(routeTable *sync.Map, pTemplate interface{}, id interface{}, sheetName string) bool {
	if tools.IsNil(pTemplate) {
		logger.Error("pTemplate is nil")
		return false
	}

	var pTemplateValue = reflect.ValueOf(pTemplate)
	if pTemplateValue.Kind() != reflect.Ptr {
		logger.Error("pTemplate should be of type *Template")
		return false
	}

	var templateValue = pTemplateValue.Elem()
	var templateType = templateValue.Type()

	// 获取sheetName
	var sheetName2 = sheetName
	if sheetName2 == "" {
		sheetName2 = templateType.Name()
	}

	var table = manager.getTemplateTable(sheetName2)
	if table != nil {
		return checkSetValue(templateValue, table[id])
	}

	excelArgs, ok := routeTable.Load(sheetName2)
	if !ok {
		logger.Error("Can not find excelFilePath for sheetName=%q", sheetName2)
		return false
	}

	var err = manager.loadTemplateTable(excelArgs.(ExcelArgs), templateType, sheetName2)
	if err != nil {
		manager.tables.Store(sheetName2, make(TemplateTable))
		return false
	}

	table = manager.getTemplateTable(sheetName2)
	return checkSetValue(templateValue, table[id])
}

func (manager *TemplateManager) getTemplates(routeTable *sync.Map, pTemplateList interface{}, args options) bool {
	var pTemplateListValue = reflect.ValueOf(pTemplateList)
	if pTemplateListValue.Kind() != reflect.Ptr {
		logger.Error("pTemplateList should be a pointer")
		return false
	}

	var templateListValue = pTemplateListValue.Elem()
	if templateListValue.Kind() != reflect.Slice {
		logger.Error("pTemplateList should be a pointer of slice")
		return false
	}

	// 取得元素类型
	var elemType = templateListValue.Type().Elem()

	// 取得args
	args.complement(elemType)

	var sheetName = args.SheetName
	var table = manager.getTemplateTable(sheetName)
	if table != nil {
		var hasData = len(table) > 0
		if hasData {
			fillSliceByTable(args, pTemplateListValue, elemType, table)
		}
		return hasData
	}

	excelArgs, ok := routeTable.Load(sheetName)
	if !ok {
		logger.Error("Can not find excelFilePath for sheetName=%q", sheetName)
		return false
	}

	var err = manager.loadTemplateTable(excelArgs.(ExcelArgs), elemType, sheetName)
	if err != nil {
		manager.tables.Store(sheetName, make(TemplateTable))
		return false
	}

	table = manager.getTemplateTable(sheetName)
	fillSliceByTable(args, pTemplateListValue, elemType, table)
	return true
}

func (manager *TemplateManager) getTemplateTable(sheetName string) TemplateTable {
	var table, ok = manager.tables.Load(sheetName)
	if !ok {
		return nil
	}

	return table.(TemplateTable)
}

func (manager *TemplateManager) loadTemplateTable(args ExcelArgs, templateType reflect.Type, sheetName string) error {
	return loadOneSheet(args, sheetName, func(reader excel.Reader) error {
		// double check，如果已经被其它协程加载过了，则不再重复加载
		if _, ok := manager.tables.Load(sheetName); ok {
			return nil
		}

		var pSlice = makeSlice(templateType)
		var err = reader.ReadAll(pSlice.Interface())
		if err != nil {
			return logger.Dot(err)
		}

		// 原始的slice因为是一个struct，所以不能直接用，只能从pSlice.Elem()重新获取一个
		var slice = pSlice.Elem()
		var table = fillTemplateTable(slice)
		manager.tables.Store(sheetName, table)
		return nil
	})
}

func fillSliceByTable(args options, pTemplateListValue reflect.Value, elemType reflect.Type, table TemplateTable) {
	var slice = reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(table))
	var filter = args.Filter
	for _, item := range table {
		if filter(item) {
			slice = reflect.Append(slice, reflect.ValueOf(item))
		}
	}

	pTemplateListValue.Elem().Set(slice)
}

func checkSetValue(v reflect.Value, i interface{}) bool {
	if i != nil {
		v.Set(reflect.ValueOf(i))
		return true
	}

	return false
}

func makeSlice(elemType reflect.Type) reflect.Value {
	// reflect.SliceOf() --> 我们平时make([]int, 0, 8)的时候，这里传入的也是slice的type，而不是直接传入elemType
	var slice = reflect.MakeSlice(reflect.SliceOf(elemType), 0, 8)
	pSlice := reflect.New(slice.Type())
	// pSlice背后是对象指针，因此需要使用.Elem().Set()设值
	pSlice.Elem().Set(slice)
	return pSlice
}

func fillTemplateTable(slice reflect.Value) TemplateTable {
	var count = slice.Len()
	var table = make(TemplateTable, count)

	for i := 0; i < count; i++ {
		var item = slice.Index(i)

		// Id 或 ID 都是可以接受的，但必须有
		var fieldId = item.FieldByNameFunc(func(s string) bool {
			return s == "Id" || s == "ID" || s == "id"
		})

		if !fieldId.IsValid() {
			logger.Error("You must define an \"Id\" field in template struct.")
			continue
		}

		var id = fieldId.Interface()
		var newItem = item.Interface()
		var oldItem, ok = table[id]
		if ok {
			logger.Error("Found duplicate templates: id=%d, oldItem=%v, newItem=%v", id, oldItem, newItem)
		}

		table[id] = newItem
	}

	return table
}
