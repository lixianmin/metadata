package metadata

import (
	"testing"
	"time"

	"github.com/lixianmin/logo"
)

const testExcelFilePath = "res/metadata.xlsx"
const testExcelFilePath2 = "res/metadata2.xlsx"
const testWebFileUrl = "https://github.com/lixianmin/metadata/raw/master/res/metadata.xlsx"
const testWebFileUrl2 = "https://github.com/lixianmin/metadata/raw/master/res/metadata2.xlsx"

func TestNewWebFile(t *testing.T) {
	var manager = NewManager()
	manager.AddExcel(WithFile(testWebFileUrl))
	manager.AddExcel(WithFile(testWebFileUrl2))

	var template TestTemplate
	var another AnotherTemplate

	for i := 0; i < 8; i++ {
		manager.GetTemplate(&template, 1, WithSheet("TestTemplate"))
		manager.GetTemplate(&another, 1)
		logo.JsonI("template", template, "another", another)

		time.Sleep(2 * time.Second)
	}
	//time.Sleep(time.Hour)
}
