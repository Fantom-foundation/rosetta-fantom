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
	"testing"

	"github.com/Fantom-foundation/rosetta-fantom/fantom"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	tests := map[string]struct {
		Mode      string
		Network   string
		Port      string
		Opera     string
		SkipAdmin string
		OperaArgs string

		cfg *Configuration
		err error
	}{
		"no envs set": {
			err: errors.New("MODE must be populated"),
		},
		"only mode set": {
			Mode: string(Online),
			err:  errors.New("NETWORK must be populated"),
		},
		"only mode and network set": {
			Mode:      string(Online),
			Network:   Mainnet,
			OperaArgs: "--",
			err:       errors.New("PORT must be populated"),
		},
		"all set (mainnet)": {
			Mode:      string(Online),
			Network:   Mainnet,
			Port:      "1000",
			SkipAdmin: "FALSE",
			OperaArgs: "--",
			cfg: &Configuration{
				Mode: Online,
				Network: &types.NetworkIdentifier{
					Network:    fantom.MainnetNetwork,
					Blockchain: fantom.Blockchain,
				},
				GenesisBlockIdentifier: fantom.FantomMainnetGenesisBlockIdentifier,
				Port:                   1000,
				OperaURL:               DefaultOperaURL,
				OperaArguments:         "--",
				SkipAdmin:              false,
				ChainID:                big.NewInt(0xFA),
			},
		},
		"all set (mainnet) + opera": {
			Mode:      string(Online),
			Network:   Mainnet,
			Port:      "1000",
			Opera:     "http://blah",
			SkipAdmin: "TRUE",
			OperaArgs: "--",
			cfg: &Configuration{
				Mode: Online,
				Network: &types.NetworkIdentifier{
					Network:    fantom.MainnetNetwork,
					Blockchain: fantom.Blockchain,
				},
				GenesisBlockIdentifier: fantom.FantomMainnetGenesisBlockIdentifier,
				Port:                   1000,
				OperaURL:               "http://blah",
				RemoteOpera:            true,
				OperaArguments:         "--",
				SkipAdmin:              true,
				ChainID:                big.NewInt(0xFA),
			},
		},
		"all set (testnet)": {
			Mode:      string(Online),
			Network:   Testnet,
			Port:      "1000",
			SkipAdmin: "TRUE",
			OperaArgs: "--",
			cfg: &Configuration{
				Mode: Online,
				Network: &types.NetworkIdentifier{
					Network:    fantom.TestnetNetwork,
					Blockchain: fantom.Blockchain,
				},
				GenesisBlockIdentifier: fantom.FantomTestnetGenesisBlockIdentifier,
				Port:                   1000,
				OperaURL:               DefaultOperaURL,
				OperaArguments:         "--",
				SkipAdmin:              true,
				ChainID:                big.NewInt(0xFA2),
			},
		},
		"invalid mode": {
			Mode:    "bad mode",
			Network: Testnet,
			Port:    "1000",
			err:     errors.New("bad mode is not a valid mode"),
		},
		"invalid network": {
			Mode:    string(Offline),
			Network: "bad network",
			Port:    "1000",
			err:     errors.New("bad network is not a valid network"),
		},
		"invalid port": {
			Mode:      string(Offline),
			Network:   Testnet,
			Port:      "bad port",
			OperaArgs: "--",
			err:       errors.New("unable to parse port bad port"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			os.Setenv(ModeEnv, test.Mode)
			os.Setenv(NetworkEnv, test.Network)
			os.Setenv(PortEnv, test.Port)
			os.Setenv(OperaEnv, test.Opera)
			os.Setenv(SkipAdminEnv, test.SkipAdmin)
			os.Setenv(OperaArgsEnv, test.OperaArgs)

			cfg, err := LoadConfiguration()
			if test.err != nil {
				assert.Nil(t, cfg)
				assert.Contains(t, err.Error(), test.err.Error())
			} else {
				fmt.Printf("%s / %s\n", cfg.Network.Network, test.cfg.Network.Network)
				assert.Equal(t, test.cfg, cfg)
				assert.NoError(t, err)
			}
		})
	}
}
