package plugin

import (
	"testing"
)

func TestPathsChanged(t *testing.T) {
	// assert that PathsChanged fulfills the converter.Plugin interface
	_ = NewRouter(WithConvertPlugins(NewPathsChanged()))
}
