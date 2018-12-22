package vlog

import (
	"testing"
)

func TestInit(t *testing.T) {
	err := Init("logs", "debug", 365)
	if err != nil {
		t.Log("xx:", err.Error())
	}
	Debug("这是debug信息%s", "xx")
}
