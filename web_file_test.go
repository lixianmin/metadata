package metadata

import (
	"github.com/lixianmin/metadata/logger"
	"testing"
	"time"
)

var testWebFileUrl = "https://dbp-media.bj.bcebos.com/0123456789/metadata.xlsx"

func TestNewWebFile(t *testing.T) {
	Init(nil, testWebFileUrl)

	var template TestTemplate
	for i := 0; i < 8; i++ {
		GetTemplate(1, &template)
		logger.Info("template=%v", template)

		time.Sleep(30 * time.Second)
	}
	//time.Sleep(time.Hour)
}
