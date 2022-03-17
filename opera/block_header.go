// Copyright 2022 Fantom Foundation, Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opera

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
