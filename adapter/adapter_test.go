package adapter

import "testing"

func TestConfig(t *testing.T) {
	adapter := GetAdapter("sqlite:///Users/Amin/.doop/doop.db")
	//TODO: some tests
	adapter.Close()
}
