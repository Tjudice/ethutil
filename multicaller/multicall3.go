package multicaller

import (
	"gfx.cafe/open/ghost/abi"
	"gfx.cafe/open/ghost/abir"
	"github.com/ethereum/go-ethereum/common"
)

var _ MulticallCodec[Multicall3Args] = &Multicall3{}

type Multicall3 struct{}

type Multicall3Args struct {
	Target       common.Address
	AllowFailure bool
	CallData     []byte
}

type Result struct {
	Results []ResultItem `abi:"(bool,bytes)[]"`
}

type ResultItem struct {
	Success    bool   `abi:"bool"`
	ReturnData []byte `abi:"bytes"`
}

func NewMulticall3(address common.Address, maxBatchSize int) *Multicall[Multicall3Args] {
	return NewMultiCall[Multicall3Args](address, maxBatchSize, &Multicall3{})
}

func (m *Multicall3) EncodeBatch(calls []Multicall3Args) ([]byte, error) {
	builder := new(abi.Builder).EnterDynamicArray()
	for _, args := range calls {
		builder.EnterTuple().Address(args.Target).Bool(args.AllowFailure).Bytes(args.CallData).Exit()
	}
	builder = builder.Exit()
	return builder.Finish(abi.SIG("aggregate3", abi.SLICE(abi.TUPLE(abi.ADDRESS, abi.BOOL, abi.BYTES))).Fn()), nil
}

func (m *Multicall3) DecodeBatch(response []byte) (callResponses [][]byte, err error) {
	var res []ResultItem
	err = abir.DecodeInto(abi.NewDecoder(response), &res)
	if err != nil {
		return nil, err
	}
	resultBytes := make([][]byte, 0, len(res))
	for _, k := range res {
		resultBytes = append(resultBytes, k.ReturnData)
	}
	return resultBytes, nil
}
