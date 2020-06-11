package metadata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testWebFileUrl = "https://dbp-media.bj.bcebos.com/0123456789/metadata.xlsx"

func TestNewWebFile(t *testing.T) {
	var web = NewWebFile(testWebFileUrl)
	assert.NotNil(t, web)
	//time.Sleep(time.Hour)
}
