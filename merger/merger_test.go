package main

import (
	"testing"
	"fmt"
)

func Test_compareDir(t *testing.T) {
	cls, err := compareDir(directory, remoteDir)
	if err != nil {
		t.Fatal("compareDir returns", err.Error())
	}

	for _, cl := range cls {
		fmt.Println(cl.StringNearest(0.8))
	}
}
