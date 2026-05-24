package fact

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"runtime"
	"strconv"
	"sync"
)

var (
	ErrFactorizationCancelled = errors.New("cancelled")
	ErrWriterInteraction      = errors.New("writer interaction")
)

type Factorizer interface {
	Factorize(ctx context.Context, numbers []int, writer io.Writer) error
}

type factorizerImpl struct {
	factorizationWorkers int
	writeWorkers         int
}

func New(opts ...FactorizeOption) (*factorizerImpl, error) {
	f := &factorizerImpl{
		factorizationWorkers: 1 + runtime.GOMAXPROCS(0)/2,
		writeWorkers:         1 + runtime.GOMAXPROCS(0)/2,
	}
	for _, opt := range opts {
		opt(f)
	}
	if f.factorizationWorkers <= 0 {
		return nil, fmt.Errorf("invalid factorization workers: %d", f.factorizationWorkers)
	}
	if f.writeWorkers <= 0 {
		return nil, fmt.Errorf("invalid write workers: %d", f.writeWorkers)
	}
	return f, nil
}

type FactorizeOption func(*factorizerImpl)

func WithFactorizationWorkers(workers int) FactorizeOption {
	return func(f *factorizerImpl) {
		f.factorizationWorkers = workers
	}
}

func WithWriteWorkers(workers int) FactorizeOption {
	return func(f *factorizerImpl) {
		f.writeWorkers = workers
	}
}

func (f *factorizerImpl) cancelledErr(ctx context.Context) error {
	return errors.Join(ErrFactorizationCancelled, context.Cause(ctx))
}

func (f *factorizerImpl) Factorize(
	ctx context.Context,
	numbers []int,
	writer io.Writer,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wgFactor := new(sync.WaitGroup)
	wgWriters := new(sync.WaitGroup)

	numbersCh := make(chan int)
	factorCh := make(chan []byte)

	var resErr error
	var once sync.Once

	onceErr := func(err error) {
		once.Do(func() {
			resErr = err
		})
	}

	go func() {
		defer close(numbersCh)
		f.addNumbersInChannel(ctx, numbersCh, numbers)
	}()

	for i := 0; i < f.factorizationWorkers; i++ {
		wgFactor.Go(func() {
			for {
				select {
				case <-ctx.Done():
					onceErr(f.cancelledErr(ctx))
					return
				case num, ok := <-numbersCh:
					if !ok {
						return
					}
					factor := f.getFactorization(num)
					select {
					case <-ctx.Done():
						onceErr(f.cancelledErr(ctx))
						return
					case factorCh <- factor:
					}
				}
			}
		})
	}

	for i := 0; i < f.writeWorkers; i++ {
		wgWriters.Go(func() {
			for {
				select {
				case <-ctx.Done():
					onceErr(f.cancelledErr(ctx))
					return
				case factor, ok := <-factorCh:
					if !ok {
						return
					}
					factor = append(factor, '\n')
					if _, err := writer.Write(factor); err != nil {
						onceErr(errors.Join(ErrWriterInteraction, err))
						cancel()
						return
					}
				}
			}
		})
	}

	wgFactor.Wait()
	close(factorCh)
	wgWriters.Wait()
	return resErr
}

func (f *factorizerImpl) addNumbersInChannel(ctx context.Context, numbersCh chan<- int, numbers []int) {
	for _, number := range numbers {
		select {
		case <-ctx.Done():
			return
		case numbersCh <- number:
		}
	}
}

func (f *factorizerImpl) getFactorization(value int) []byte {
	var builder []byte
	builder = append(builder, strconv.Itoa(value)...)
	builder = append(builder, " = "...)

	if value == -1 || value == 1 || value == math.MinInt {
		builder = append(builder, strconv.Itoa(value)...)
		return builder
	}

	if value < 0 {
		builder = append(builder, "-1 * "...)
		value *= -1
	}
	for i := 2; i <= int(math.Sqrt(float64(value))); i++ {
		for value%i == 0 {
			value /= i
			builder = append(builder, strconv.Itoa(i)...)
			if value != 1 {
				builder = append(builder, " * "...)
			}
		}
	}
	if value != 1 {
		builder = append(builder, strconv.Itoa(value)...)
	}
	return builder
}
