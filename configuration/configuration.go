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

package configuration

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/coinbase/rosetta-ethereum/ethereum"

	"github.com/coinbase/rosetta-sdk-go/types"
)

// Mode is the setting that determines if
// the implementation is "online" or "offline".
type Mode string

const (
	// Online is when the implementation is permitted
	// to make outbound connections.
	Online Mode = "ONLINE"

	// Offline is when the implementation is not permitted
	// to make outbound connections.
	Offline Mode = "OFFLINE"

	// Mainnet is the Fantom Mainnet.
	Mainnet string = "MAINNET"

	// Testnet is the Fantom Testnet.
	Testnet string = "TESTNET"

	// DataDirectory is the default location for all
	// persistent data.
	DataDirectory = "/data"

	// ModeEnv is the environment variable read
	// to determine mode.
	ModeEnv = "MODE"

	// NetworkEnv is the environment variable
	// read to determine network.
	NetworkEnv = "NETWORK"

	// PortEnv is the environment variable
	// read to determine the port for the Rosetta
	// implementation.
	PortEnv = "PORT"

	// GethEnv is an optional environment variable
	// used to connect rosetta-ethereum to an already
	// running geth node.
	GethEnv = "GETH"

	// DefaultGethURL is the default URL for
	// a running geth node. This is used
	// when GethEnv is not populated.
	DefaultGethURL = "http://localhost:8545"

	// SkipGethAdminEnv is an optional environment variable
	// to skip geth `admin` calls which are typically not supported
	// by hosted node services. When not set, defaults to false.
	SkipGethAdminEnv = "SKIP_GETH_ADMIN"

	// MiddlewareVersion is the version of rosetta-ethereum.
	MiddlewareVersion = "0.0.4"
)

// Configuration determines how
type Configuration struct {
	Mode                   Mode
	Network                *types.NetworkIdentifier
	GenesisBlockIdentifier *types.BlockIdentifier
	GethURL                string
	RemoteGeth             bool
	Port                   int
	GethArguments          string
	SkipGethAdmin          bool
	ChainID                *big.Int
}

// LoadConfiguration attempts to create a new Configuration
// using the ENVs in the environment.
func LoadConfiguration() (*Configuration, error) {
	config := &Configuration{}

	modeValue := Mode(os.Getenv(ModeEnv))
	switch modeValue {
	case Online:
		config.Mode = Online
	case Offline:
		config.Mode = Offline
	case "":
		return nil, errors.New("MODE must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid mode", modeValue)
	}

	networkValue := os.Getenv(NetworkEnv)
	switch networkValue {
	case Mainnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: ethereum.Blockchain,
			Network:    ethereum.MainnetNetwork,
		}
		config.GenesisBlockIdentifier = ethereum.FantomMainnetGenesisBlockIdentifier
		config.ChainID = big.NewInt(0xFA)
		config.GethArguments = ethereum.MainnetGethArguments
	case Testnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: ethereum.Blockchain,
			Network:    ethereum.TestnetNetwork,
		}
		config.GenesisBlockIdentifier = ethereum.FantomTestnetGenesisBlockIdentifier
		config.ChainID = big.NewInt(0xFA2)
		config.GethArguments = ethereum.TestnetGethArguments
	case "":
		return nil, errors.New("NETWORK must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid network", networkValue)
	}

	config.GethURL = DefaultGethURL
	envGethURL := os.Getenv(GethEnv)
	if len(envGethURL) > 0 {
		config.RemoteGeth = true
		config.GethURL = envGethURL
	}

	config.SkipGethAdmin = false
	envSkipGethAdmin := os.Getenv(SkipGethAdminEnv)
	if len(envSkipGethAdmin) > 0 {
		val, err := strconv.ParseBool(envSkipGethAdmin)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to parse SKIP_GETH_ADMIN %s", err, envSkipGethAdmin)
		}
		config.SkipGethAdmin = val
	}

	portValue := os.Getenv(PortEnv)
	if len(portValue) == 0 {
		return nil, errors.New("PORT must be populated")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || len(portValue) == 0 || port <= 0 {
		return nil, fmt.Errorf("%w: unable to parse port %s", err, portValue)
	}
	config.Port = port

	return config, nil
}
