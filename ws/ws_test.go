package ws

import (
	"testing"
)

func TestInit(t *testing.T) {
	if Init() == nil {
		t.Fatal("Init should not return nil")
	}
}
