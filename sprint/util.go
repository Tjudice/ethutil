package sprint

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

func toGethFilter(e ethereum.FilterQuery) map[string]any {
	m := map[string]any{
		"topics":    e.Topics,
		"fromBlock": uint256.MustFromBig(e.FromBlock).Hex(),
		"toBlock":   uint256.MustFromBig(e.ToBlock).Hex(),
	}
	if len(e.Addresses) > 0 {
		m["address"] = e.Addresses
	}
	return m
}

func onConflictDoNothing(s string) string {
	return s + " ON CONFLICT DO NOTHING"
}

func checkAllTxsExist(block *BlockInfo, txs []common.Hash) bool {
	for i := range txs {
		if _, ok := block.TxMap[txs[i]]; !ok {
			return false
		}
	}
	return true
}
