package config

import "testing"

func TestConfig(t *testing.T) {
	key := "testKey"
	value := "testValue"
	SetEnv(key, value)
	tarvalue := GetEnv(key, "defaultValue")
	if tarvalue != value {
		t.Errorf("Expected environment variable")
	}
}
