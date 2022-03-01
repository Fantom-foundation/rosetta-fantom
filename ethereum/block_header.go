package ethereum

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// blockHeader wraps block Header and adds block header hash obtained from RPC
type blockHeader struct {
	types.Header       `json:"-"`
	Hash   common.Hash `json:"hash"`
}

func (h *blockHeader) UnmarshalJSON(input []byte) error {
	// unmarshal hash
	var temp struct {
		Hash   common.Hash `json:"hash"`
	}
	if err := json.Unmarshal(input, &temp); err != nil {
		return err
	}
	h.Hash = temp.Hash

	// unmarshall nested Header
	if err := json.Unmarshal(input, &h.Header); err != nil {
		return err
	}
	return nil
}
