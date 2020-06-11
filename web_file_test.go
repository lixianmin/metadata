package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"testing"
	"time"
)

func TestNewWebFile(t *testing.T) {
	var manager = &MetadataManager{}
	manager.AddExcel(testWebFileUrl)

	var template TestTemplate
	for i := 0; i < 8; i++ {
		manager.GetTemplate(1, &template)
		logger.Info("template=%v", template)

		time.Sleep(30 * time.Second)
	}
	//time.Sleep(time.Hour)
}
