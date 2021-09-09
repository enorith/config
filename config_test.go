package config_test

import (
	"testing"

	"github.com/enorith/config"
	"gopkg.in/yaml.v3"
)

var confFoo = `
name: test
maps:
 foo:
  type: foo
 bar:
  type: bar
  enabled: yes
  other: 42
`

type FooConf struct {
	Name string               `yaml:"name"`
	Maps map[string]yaml.Node `yaml:"maps"`
}

type BarConf struct {
	Type    string `yaml:"type"`
	Enabled string `yaml:"enabled"`
	Other   int    `yaml:"other"`
}

func TestUnmarshal(t *testing.T) {
	var fc FooConf
	e := config.UnmarshalBytes([]byte(confFoo), &fc)
	if e != nil {
		t.Fatal(e)
	}

	var bc BarConf
	x := fc.Maps["bar"]
	config.UnmarshalNode(x, &bc)
	t.Log(bc)
}
