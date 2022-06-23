package plugin

import (
	"testing"
)

func TestPathsChangedFulfillsPluginInterface(t *testing.T) {
	_ = NewRouter(WithConvertPlugins(NewPathsChanged()))
}
