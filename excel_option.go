package metadata

/********************************************************************
created:    2024-07-24
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type excelOptions struct {
	Uri           string // 支持http, https的远程excel文件，也支持本地的excel文件
	TitleRowIndex int    // excel表格标题行索引，默认值0，即Excel中的第1行
	Skip          int    // 跳过标题后的n行，默认值0（不跳过）。空行不计算在内
}

type ExcelOption func(*excelOptions)

// WithFile 支持http, https的远程excel文件，也支持本地的excel文件
func WithFile(uri string) ExcelOption {
	return func(opt *excelOptions) {
		opt.Uri = uri
	}
}

// WithTitleRowIndex excel表格标题行索引，默认值0，即Excel中的第1行
func WithTitleRowIndex(titleRowIndex int) ExcelOption {
	return func(opt *excelOptions) {
		opt.TitleRowIndex = titleRowIndex
	}
}

// WithSkipRows 跳过标题后的n行，默认值0（不跳过）。空行不计算在内
func WithSkipRows(skip int) ExcelOption {
	return func(opt *excelOptions) {
		opt.Skip = skip
	}
}
