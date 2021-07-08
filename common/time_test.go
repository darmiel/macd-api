package common

import (
	"testing"
	"time"
)

func TestNormalizeTime(t *testing.T) {
	const (
		input  = 1625772149
		output = 1625702400
	)
	dt := time.Unix(input, 0)
	nm := NormalizeTime(dt)
	if nm.Unix() != output {
		t.Errorf("%d expected but got %d", output, nm.Unix())
	}
}
