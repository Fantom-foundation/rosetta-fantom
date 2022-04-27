// Copyright 2020 Coinbase, Inc.
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

package fantom

import (
	"context"
	"github.com/ethereum/go-ethereum/common"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	// NodeVersion is the version of opera we are using.
	NodeVersion = "1.1.0-rc.4"

	// Blockchain is Fantom.
	Blockchain string = "Fantom"

	// TestnetNetwork is the value of the network
	// in NetworkIdentifier.
	TestnetNetwork string = "Testnet"

	// MainnetNetwork is the value of the network
	// in MainnetNetworkIdentifier.
	MainnetNetwork string = "Mainnet"

	// Symbol is the symbol value
	// used in Currency.
	Symbol = "FTM"

	// Decimals is the decimals value
	// used in Currency.
	Decimals = 18

	// FeeOpType is used to represent fee operations.
	FeeOpType = "FEE"

	// CallOpType is used to represent CALL trace operations.
	CallOpType = "CALL"

	// CreateOpType is used to represent CREATE trace operations.
	CreateOpType = "CREATE"

	// Create2OpType is used to represent CREATE2 trace operations.
	Create2OpType = "CREATE2"

	// SelfDestructOpType is used to represent SELFDESTRUCT trace operations.
	SelfDestructOpType = "SELFDESTRUCT"

	// CallCodeOpType is used to represent CALLCODE trace operations.
	CallCodeOpType = "CALLCODE"

	// DelegateCallOpType is used to represent DELEGATECALL trace operations.
	DelegateCallOpType = "DELEGATECALL"

	// StaticCallOpType is used to represent STATICCALL trace operations.
	StaticCallOpType = "STATICCALL"

	// DestructOpType is a synthetic operation used to represent the
	// deletion of suicided accounts that still have funds at the end
	// of a transaction.
	DestructOpType = "DESTRUCT"

	// SuccessStatus is the status of any
	// Opera operation considered successful.
	SuccessStatus = "SUCCESS"

	// FailureStatus is the status of any
	// Opera operation considered unsuccessful.
	FailureStatus = "FAILURE"

	// HistoricalBalanceSupported is whether
	// historical balance is supported.
	HistoricalBalanceSupported = true

	// GenesisBlockIndex is the index of the
	// genesis block.
	GenesisBlockIndex = int64(0)

	// TransferGasLimit is the gas limit
	// of a transfer.
	TransferGasLimit = int64(21000) //nolint:gomnd

	// MainnetOperaArguments are the arguments to start a mainnet Opera instance.
	MainnetOperaArguments = `--config=/app/fantom/opera.toml --genesis=/data/mainnet.g`

	// TestnetOperaArguments are the arguments to start a testnet Opera instance.
	TestnetOperaArguments = `--config=/app/fantom/opera.toml --genesis=/data/testnet.g`

	// IncludeMempoolCoins does not apply to rosetta-fantom as it is not UTXO-based.
	IncludeMempoolCoins = false
)

var (
	// FantomMainnetGenesisHash represents the hash of the genesis block (block 0)
	FantomMainnetGenesisHash = common.HexToHash("0x00000000000003e83fddf1e9330f0a8691d9f0b2af57b38c3bb85488488a40df")

	// FantomTestnetGenesisHash represents the hash of the genesis block (block 0)
	FantomTestnetGenesisHash = common.HexToHash("0x00000000000003e8c717f00dc4306a6ff72eabc9a6ec6e4a46bf6ba044ca88d2")
)

var (
	// FantomMainnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the mainnet genesis block.
	FantomMainnetGenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  FantomMainnetGenesisHash.Hex(),
		Index: GenesisBlockIndex,
	}

	// FantomTestnetGenesisBlockIdentifier is the *types.BlockIdentifier
	// of the mainnet genesis block.
	FantomTestnetGenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  FantomTestnetGenesisHash.Hex(),
		Index: GenesisBlockIndex,
	}

	// Currency is the *types.Currency for all
	// Opera networks.
	Currency = &types.Currency{
		Symbol:   Symbol,
		Decimals: Decimals,
	}

	// OperationTypes are all suppoorted operation types.
	OperationTypes = []string{
		FeeOpType,
		CallOpType,
		CreateOpType,
		Create2OpType,
		SelfDestructOpType,
		CallCodeOpType,
		DelegateCallOpType,
		StaticCallOpType,
		DestructOpType,
	}

	// OperationStatuses are all supported operation statuses.
	OperationStatuses = []*types.OperationStatus{
		{
			Status:     SuccessStatus,
			Successful: true,
		},
		{
			Status:     FailureStatus,
			Successful: false,
		},
	}

	// CallMethods are all supported call methods.
	CallMethods = []string{
		"eth_getBlockByNumber",
		"eth_getTransactionReceipt",
		"eth_call",
		"eth_estimateGas",
	}
)

// JSONRPC is the interface for accessing go-ethereum's JSON RPC endpoint.
type JSONRPC interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
	Close()
}

// GraphQL is the interface for accessing go-ethereum's GraphQL endpoint.
type GraphQL interface {
	Query(ctx context.Context, input string) (string, error)
}

// CallType returns a boolean indicating
// if the provided trace type is a call type.
func CallType(t string) bool {
	callTypes := []string{
		CallOpType,
		CallCodeOpType,
		DelegateCallOpType,
		StaticCallOpType,
	}

	for _, callType := range callTypes {
		if callType == t {
			return true
		}
	}

	return false
}

// CreateType returns a boolean indicating
// if the provided trace type is a create type.
func CreateType(t string) bool {
	createTypes := []string{
		CreateOpType,
		Create2OpType,
	}

	for _, createType := range createTypes {
		if createType == t {
			return true
		}
	}

	return false
}
