package config

import (
	"fmt"
	"testing"
)

var testCfgFile = "./cfg.json"

func TestParse(t *testing.T) {
	cfg := LoadConfigFile(testCfgFile)
	if cfg.GetFloat("port") != 8010 || cfg.GetString("role") != "cc" || cfg.GetBool("idgen") == false || cfg.GetFloat("sid") != 0 {
		t.Error("Fatal")
	}
}

func TestA(t *testing.T) {
	var ss string
	str := "https://www.ygdy8.net/html/gndy/dyzz/list_23_"
	for i := 24; i < 211; i++ {
		ss = ss + str + fmt.Sprintf("%v", i) + ".html,"
	}
	t.Log(ss)
}
