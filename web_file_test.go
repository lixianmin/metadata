package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"testing"
	"time"
)

func TestNewWebFile(t *testing.T) {
	var manager = &MetadataManager{}
	manager.AddExcel(testWebFileUrl)
	manager.AddExcel(testWebFileUrl2)

	var template TestTemplate
	var another AnotherTemplate

	for i := 0; i < 8; i++ {
		manager.GetTemplate(1, &template)
		manager.GetTemplate(1, &another)
		logger.Info("template=%v, another=%v", template, another)

		time.Sleep(30 * time.Second)
	}
	//time.Sleep(time.Hour)
}
