package metadata

/********************************************************************
created:    2020-06-14
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type ExcelArgs struct {
	FilePath      string // 支持http, https的远程excel文件，也支持本地的excel文件
	TitleRowIndex int    // excel表格标题行索引，默认值0，即Excel中的第1行
	Skip          int    // 跳过标题后的n行，默认值0（不跳过）。空行不计算在内
	OnAdded       func(localPath string)
}
