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

	"github.com/Fantom-foundation/rosetta-fantom/fantom"

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

	// OperaEnv is an optional environment variable
	// used to connect rosetta-fantom to an already
	// running Opera node.
	OperaEnv = "OPERA"

	// DefaultOperaURL is the default URL for
	// a running Opera node. This is used
	// when OperaEnv is not populated.
	DefaultOperaURL = "http://localhost:18545"

	// SkipAdminEnv is an optional environment variable
	// to skip RPC `admin` calls which are typically not supported
	// by hosted node services. When not set, defaults to true.
	SkipAdminEnv = "SKIP_ADMIN"

	// MiddlewareVersion is the version of rosetta-fantom.
	MiddlewareVersion = "0.0.4"
)

// Configuration determines how
type Configuration struct {
	Mode                   Mode
	Network                *types.NetworkIdentifier
	GenesisBlockIdentifier *types.BlockIdentifier
	OperaURL               string
	RemoteOpera            bool
	Port                   int
	OperaArguments         string
	SkipAdmin              bool
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
			Blockchain: fantom.Blockchain,
			Network:    fantom.MainnetNetwork,
		}
		config.GenesisBlockIdentifier = fantom.FantomMainnetGenesisBlockIdentifier
		config.ChainID = big.NewInt(0xFA)
		config.OperaArguments = fantom.MainnetOperaArguments
	case Testnet:
		config.Network = &types.NetworkIdentifier{
			Blockchain: fantom.Blockchain,
			Network:    fantom.TestnetNetwork,
		}
		config.GenesisBlockIdentifier = fantom.FantomTestnetGenesisBlockIdentifier
		config.ChainID = big.NewInt(0xFA2)
		config.OperaArguments = fantom.TestnetOperaArguments
	case "":
		return nil, errors.New("NETWORK must be populated")
	default:
		return nil, fmt.Errorf("%s is not a valid network", networkValue)
	}

	config.OperaURL = DefaultOperaURL
	envOperaURL := os.Getenv(OperaEnv)
	if len(envOperaURL) > 0 {
		config.RemoteOpera = true
		config.OperaURL = envOperaURL
	}

	config.SkipAdmin = false
	envSkipAdmin := os.Getenv(SkipAdminEnv)
	if len(envSkipAdmin) > 0 {
		val, err := strconv.ParseBool(envSkipAdmin)
		if err != nil {
			return nil, fmt.Errorf("%w: unable to parse SKIP_ADMIN %s", err, envSkipAdmin)
		}
		config.SkipAdmin = val
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
