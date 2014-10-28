package core

import "testing"

func TestConfig(t *testing.T) {
	doop := GetDoop()
	doop.getConfig()
	if doop.config.Database.DSN == "" {
		t.Fatalf("No DSN set or Doop not installed.")
	}
}
