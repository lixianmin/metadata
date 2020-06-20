package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"testing"
	"time"
)

func TestNewWebFile(t *testing.T) {
	var manager = &MetadataManager{}
	manager.AddExcel(ExcelArgs{FilePath: testWebFileUrl})
	manager.AddExcel(ExcelArgs{FilePath: testWebFileUrl2})

	var template TestTemplate
	var another AnotherTemplate

	for i := 0; i < 8; i++ {
		manager.GetTemplate(&template, Args{Id: 1, SheetName: "TestTemplate"})
		manager.GetTemplate(&another, Args{Id: 1, SheetName: "FakeTemplate"})
		logger.Info("template=%v, another=%v", template, another)

		time.Sleep(30 * time.Second)
	}
	//time.Sleep(time.Hour)
}
