package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"testing"
	"time"
)

func TestNewWebFile(t *testing.T) {
	var manager = &Manager{}
	manager.AddExcel(WithFile(testWebFileUrl))
	manager.AddExcel(WithFile(testWebFileUrl2))

	var template TestTemplate
	var another AnotherTemplate

	for i := 0; i < 8; i++ {
		manager.GetTemplate(&template, 1, WithSheet("TestTemplate"))
		manager.GetTemplate(&another, 1)
		logger.Info("template=%v, another=%v", template, another)

		time.Sleep(2 * time.Second)
	}
	//time.Sleep(time.Hour)
}
