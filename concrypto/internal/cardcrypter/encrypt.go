package cardcrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

var errNegativeWorkers = errors.New("negative workers")

type CardNumber [16]byte

type Card struct {
	ID     string
	Number CardNumber
}

type Crypter interface {
	Encrypt(cards []Card, key []byte) ([]string, error)
}

type crypterImpl struct {
	workers int
}

func New(opts ...CrypterOption) *crypterImpl {
	c := &crypterImpl{
		workers: runtime.GOMAXPROCS(0),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type CrypterOption func(*crypterImpl)

func WithWorkers(workers int) CrypterOption {
	const maxWorkers = 1000
	return func(c *crypterImpl) {
		c.workers = min(workers, maxWorkers)
	}
}

func (c *crypterImpl) Encrypt(cards []Card, key []byte) ([]string, error) {
	if c.workers <= 0 {
		return nil, errNegativeWorkers
	}
	if len(cards) == 0 {
		return nil, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()

	wg := new(sync.WaitGroup)

	perWorker := (len(cards) + c.workers - 1) / c.workers
	out := make([]string, len(cards))

	for worker := 0; worker < c.workers; worker++ {
		wg.Go(func() {
			start := worker * perWorker
			end := min(start+perWorker, len(cards))
			for i := start; i < end; i++ {
				dataLen := nonceSize + len(cards[i].Number) + gcm.Overhead()
				encodedLen := hex.EncodedLen(dataLen)
				buf := make([]byte, dataLen+encodedLen)
				nonce := buf[:nonceSize]

				_, _ = rand.Read(nonce)

				additionalData := unsafe.Slice(unsafe.StringData(cards[i].ID), len(cards[i].ID))
				gcm.Seal(nonce, nonce, cards[i].Number[:], additionalData)

				hex.Encode(buf[dataLen:], buf[:dataLen])
				out[i] = unsafe.String(&buf[dataLen], encodedLen)
			}
		})
	}

	wg.Wait()

	return out, nil
}
