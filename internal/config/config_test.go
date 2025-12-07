package config

import (
	"testing"

	"github.com/lwmacct/251207-go-pkg-config/pkg/config"
)

func TestGenerateExample(t *testing.T) { config.RunGenerateExampleTest(t, DefaultConfig()) }
func TestConfigKeysValid(t *testing.T) { config.RunConfigKeysValidTest(t) }
