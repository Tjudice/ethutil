package sprint

import (
	"context"
	"fmt"

	"gfx.cafe/open/ghost/hexutil"
	"gfx.cafe/open/jrpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

type TransactionInfo struct {
	Hash                 common.Hash     `json:"hash"`
	From                 common.Address  `json:"from"`
	To                   *common.Address `json:"to"`
	SenderNonce          *uint256.Int    `json:"nonce"`
	Value                *uint256.Int    `json:"value"`
	Gas                  *uint256.Int    `json:"gas"`
	GasPrice             *uint256.Int    `json:"gasPrice"`
	MaxFeePerGas         *uint256.Int    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *uint256.Int    `json:"maxPriorityFeePerGas"`
}

type BlockInfo struct {
	Number       *uint256.Int                     `json:"number"`
	Timestamp    *uint256.Int                     `json:"timestamp"`
	Transactions []*TransactionInfo               `json:"transactions"`
	TxMap        map[common.Hash]*TransactionInfo `json:"-"`
}

type blockCache struct {
	c          jrpc.Conn
	workerChan chan struct{}
}

func NewBlockCache(ctx context.Context, c jrpc.Conn, workers int) *blockCache {
	o := &blockCache{
		c:          c,
		workerChan: make(chan struct{}, workers),
	}
	for i := 0; i < workers; i++ {
		o.workerChan <- struct{}{}
	}
	return o
}

func (t *blockCache) acquireWorker() {
	<-t.workerChan
}

func (t *blockCache) releaseWorker() {
	t.workerChan <- struct{}{}
}

// Meant to be used concurrently from caller
func (t *blockCache) BlockAt(ctx context.Context, block int) (*BlockInfo, error) {
	t.acquireWorker()
	defer t.releaseWorker()
	blk := &BlockInfo{}
	err := t.c.Do(ctx, blk, "eth_getBlockByNumber", []any{hexutil.Uint64(block), true})
	if err != nil {
		return nil, err
	}
	if blk.Number == nil {
		return nil, fmt.Errorf("nil block")
	}
	num := int(blk.Number.Uint64())
	if num != block {
		return nil, fmt.Errorf("block number doesnt match")
	}
	blk.TxMap = make(map[common.Hash]*TransactionInfo)
	for _, tx := range blk.Transactions {
		if tx.To == nil {
			tx.To = &common.Address{}
		}
		if tx.SenderNonce == nil {
			tx.SenderNonce = &uint256.Int{}
		}
		if tx.Value == nil {
			tx.Value = &uint256.Int{}
		}
		if tx.Gas == nil {
			tx.Gas = &uint256.Int{}
		}
		if tx.GasPrice == nil {
			tx.GasPrice = &uint256.Int{}
		}
		if tx.MaxFeePerGas == nil {
			tx.MaxFeePerGas = &uint256.Int{}
		}
		if tx.MaxPriorityFeePerGas == nil {
			tx.MaxPriorityFeePerGas = &uint256.Int{}
		}
		blk.TxMap[tx.Hash] = tx
	}
	return blk, nil
}
