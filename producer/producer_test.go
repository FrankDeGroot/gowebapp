package producer

import (
	"testing"
)

func TestProduce(t *testing.T) {
	Connect()
	defer Close()
}
