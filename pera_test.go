package main

import (
	"bytes"
	"testing"
)

func TestAutoPlay(t *testing.T) {

	stdout1 := new(bytes.Buffer)
	expected1 := []byte{27, 91, 50, 74,
		27, 91, 48, 59, 48, 72,
		27, 91, 50, 75, 13, 10,
		27, 91, 50, 75, 97, 13, 10,
		27, 91, 50, 74,
		27, 91, 48, 59, 48, 72,
		27, 91, 50, 75, 98, 13, 10,
		27, 91, 50, 74, 27, 91, 48, 59, 48, 72,
		27, 91, 50, 75, 99, 13, 10,
		27, 91, 50, 75, 13, 10}

	AutoPlay(stdout1, 80, 40, 10, sample1, false, false)
	if bytes.Compare(expected1, stdout1.Bytes()) != 0 {
		t.Fatal("not matched")
	}

}

var sample1 = `
a
---
b
---
c
`
