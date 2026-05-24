package fibonacci

import (
	"errors"
	"math"
	"runtime"
	"sync/atomic"
)

var ErrOverflow = errors.New("fib generator overflow")

const (
	unlocked = iota
	locked
)

type Generator interface {
	Next() uint64
}

type spinlock struct {
	state atomic.Int64
}

func (s *spinlock) lock() {
	for !s.state.CompareAndSwap(unlocked, locked) {
		runtime.Gosched()
	}
}

func (s *spinlock) unlock() {
	s.state.Store(unlocked)
}

var _ Generator = (*generatorImpl)(nil)

type generatorImpl struct {
	first        uint64
	second       uint64
	isOverflowed bool
	locker       spinlock
}

func NewGenerator() *generatorImpl {
	return &generatorImpl{
		first:  1,
		second: 0,
	}
}

func (g *generatorImpl) Next() uint64 {
	g.locker.lock()
	defer g.locker.unlock()

	if g.isOverflowed {
		panic(ErrOverflow)
	}

	if g.first > math.MaxUint64-g.second {
		g.isOverflowed = true
	}
	g.first, g.second = g.second, g.first+g.second
	return g.first
}
