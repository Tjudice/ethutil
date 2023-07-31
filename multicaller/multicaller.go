package multicaller

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"

	"gfx.cafe/open/ghost"
	"gfx.cafe/open/ghost/abi"
	"gfx.cafe/open/ghost/abir"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"golang.org/x/sync/errgroup"
)

type MulticallCodec[T any] interface {
	EncodeBatch(calls []T) ([]byte, error)
	DecodeBatch(response []byte) (callResponses [][]byte, err error)
}

type ABIDecoder interface {
	DecodeABI(data []byte) error
}

type Multicall[T any] struct {
	multicallContract common.Address
	maxSize           int
	calls             []callDecoder[T]
	codec             MulticallCodec[T]
	mut               sync.Mutex
}

type callDecoder[T any] struct {
	args   T
	result any
}

func NewMultiCall[T any](mcAddress common.Address, maxSize int, codec MulticallCodec[T]) *Multicall[T] {
	return &Multicall[T]{
		multicallContract: mcAddress,
		maxSize:           maxSize,
		codec:             codec,
	}
}

func (m *Multicall[T]) Call(callArgs T, result any) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("result must be a pointer")
	}
	m.calls = append(m.calls, callDecoder[T]{
		args:   callArgs,
		result: result,
	})
	return nil
}

func (m *Multicall[T]) Exec(ctx context.Context, cl ghost.Client, concurrent bool, block int64) error {
	m.mut.Lock()
	defer m.mut.Unlock()
	if len(m.calls) == 0 {
		return nil
	}
	divided := m.getBatches()
	wg := errgroup.Group{}
	if !concurrent {
		wg.SetLimit(1)
	}
	for _, batch := range divided {
		batch := batch
		wg.Go(func() error {
			return m.execBatch(ctx, cl, batch, block)
		})
	}
	err := wg.Wait()
	if err != nil {
		return err
	}
	m.calls = m.calls[:0]
	return wg.Wait()
}

func (m *Multicall[T]) execBatch(ctx context.Context, cl ghost.Client, batch []callDecoder[T], block int64) error {
	callArgs := make([]T, 0, len(batch))
	for _, batchItem := range batch {
		callArgs = append(callArgs, batchItem.args)
	}
	encodedBatch, err := m.codec.EncodeBatch(callArgs)
	if err != nil {
		return err
	}
	callMsg := ethereum.CallMsg{
		To:   &m.multicallContract,
		Data: encodedBatch,
	}
	result, err := cl.CallContract(ctx, callMsg, big.NewInt(block))
	if err != nil {
		return err
	}
	responses, err := m.codec.DecodeBatch(result)
	if err != nil {
		return err
	}
	if len(responses) != len(batch) {
		return fmt.Errorf("invalid number of responses: %d != %d", len(responses), len(batch))
	}
	for i, response := range responses {
		err := decode(response, batch[i].result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Multicall[T]) getBatches() [][]callDecoder[T] {
	divided := make([][]callDecoder[T], 0, len(m.calls)/m.maxSize+1)
	for i := 0; i < len(m.calls); i += m.maxSize {
		end := i + m.maxSize
		if end > len(m.calls) {
			end = len(m.calls)
		}
		divided = append(divided, m.calls[i:end])
	}
	return divided
}

func decode(data []byte, result any) (err error) {
	dec := abi.NewDecoder(data)
	switch x := result.(type) {
	case ABIDecoder:
		return x.DecodeABI(data)
	case *uint256.Int:
		r, err := dec.Uint256()
		if err != nil {
			return err
		}
		x.Set(r)
		return nil
	case *common.Address:
		r, err := dec.Address()
		if err != nil {
			return err
		}
		*x = r
		return nil
	default:
		return abir.DecodeInto(dec, result)
	}
}
