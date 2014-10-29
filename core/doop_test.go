package core

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	doop := GetDoop()
	cfg := doop.getConfig()
	fmt.Println(cfg)
	if cfg.Database.DSN == "" {
		t.Fatalf("No DSN set or Doop not installed.")
	}
}
