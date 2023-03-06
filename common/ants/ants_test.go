package ants

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"testing"
	"time"
)

func TestNewAnts(t *testing.T) {
	a, _ := NewAnts(1000, ants.WithPreAlloc(true))
	for i := 0; i < 10000; i++ {
		a.Submit(func() {
			time.Sleep(3 * time.Second)
			fmt.Println("hello world")
		})
	}
	time.Sleep(5 * time.Second)
	a.Release()
}
