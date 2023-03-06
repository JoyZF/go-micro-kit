package ants

import (
	"github.com/panjf2000/ants/v2"
	"sync"
)

type Ants struct {
	pool *ants.Pool
}

var (
	once sync.Once
	a    Ants
	err  error
)

type fn func()

func NewAnts(size int, options ...ants.Option) (*Ants, error) {
	once.Do(func() {
		pool, _ := ants.NewPool(size, options...)
		a.pool = pool
	})
	return &a, err
}

// Submit 提交任务通过调用 ants.Submit(func())方法：
func (a *Ants) Submit(f fn) error {
	return a.pool.Submit(f)
}

// Tune 动态调整 goroutine 池容量
func (a *Ants) Tune(size int) {
	a.pool.Tune(size)
}

// Release 释放 Pool
func (a *Ants) Release() {
	a.pool.Release()
}

// Reboot 重启pool
func (a *Ants) Reboot() {
	a.pool.Reboot()
}
