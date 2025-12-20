package config

import (
	"testing"

	"github.com/lwmacct/251207-go-pkg-cfgm/pkg/cfgm"
)

var helper = cfgm.ConfigTestHelper[Config]{
	ExamplePath: "config/config.example.yaml",
	ConfigPath:  "config/config.yaml",
}

func TestWriteExample(t *testing.T)    { helper.WriteExampleFile(t, DefaultConfig()) }
func TestConfigKeysValid(t *testing.T) { helper.ValidateKeys(t) }
