/*
 * Copyright 2018-2019 The CovenantSQL Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package internal

import (
	"fmt"
	"strings"

	"github.com/CovenantSQL/CovenantSQL/client"
	"github.com/CovenantSQL/CovenantSQL/conf"
	"github.com/CovenantSQL/CovenantSQL/crypto"
	"github.com/CovenantSQL/CovenantSQL/types"
)

var (
	tokenName string // get specific token's balance of current account
)

// CmdWallet is cql wallet command entity.
var CmdWallet = &Command{
	UsageLine: "cql wallet [common params] [-token type]",
	Short:     "get the wallet address and the balance of current account",
	Long: `
Wallet gets the CovenantSQL wallet address and the token balances of the current account.
e.g.
    cql wallet

    cql wallet -token Particle
`,
}

func init() {
	CmdWallet.Run = runWallet

	addCommonFlags(CmdWallet)
	addConfigFlag(CmdWallet)
	CmdWallet.Flag.StringVar(&tokenName, "token", "", "Get specific token's balance of current account, e.g. Particle, Wave, All")
}

func walletGen() string {
	//TODO if config has wallet, print and return

	publicKey := getPublicFromConfig()

	keyHash, err := crypto.PubKeyHash(publicKey)
	if err != nil {
		ConsoleLog.WithError(err).Error("unexpected error")
		SetExitStatus(1)
		return ""
	}

	return keyHash.String()

	//TODO store in config.yaml
}

func runWallet(cmd *Command, args []string) {
	commonFlagsInit(cmd)
	configInit()

	if conf.GConf.WalletAddress != "" {
		fmt.Printf("\n\nwallet address: %s\n", conf.GConf.WalletAddress)
	} else {
		fmt.Printf("\n\nwallet address: %s\n", walletGen())
	}

	var err error
	if strings.ToLower(tokenName) == "" {
		var stableCoinBalance, covenantCoinBalance uint64

		if stableCoinBalance, err = client.GetTokenBalance(types.Particle); err != nil {
			if strings.Contains(err.Error(), "no such token balance") {
				fmt.Println("Your account is not created in the TestNet, please apply tokens from our faucet first.")
				SetExitStatus(1)
			} else {
				ConsoleLog.WithError(err).Error("get Particle balance failed")
				SetExitStatus(1)
			}
		}
		ExitIfErrors()
		if covenantCoinBalance, err = client.GetTokenBalance(types.Wave); err != nil {
			if strings.Contains(err.Error(), "no such token balance") {
				fmt.Println("Your account is not created in the TestNet, please apply tokens from our faucet first.")
				SetExitStatus(1)
			} else {
				ConsoleLog.WithError(err).Error("get Wave balance failed")
				SetExitStatus(1)
			}
		}
		ExitIfErrors()

		fmt.Printf("Particle balance is: %d\n", stableCoinBalance)
		fmt.Printf("Wave balance is: %d\n", covenantCoinBalance)
	} else {
		var tokenBalance uint64
		tokenType := types.FromString(tokenName)
		if !tokenType.Listed() {
			values := make([]string, len(types.TokenList))
			for i := types.Particle; i < types.SupportTokenNumber; i++ {
				values[i] = types.TokenList[i]
			}
			ConsoleLog.Errorf("no such token supporting in CovenantSQL (what we support: %s)",
				strings.Join(values, ", "))
			SetExitStatus(1)
			return
		}
		if tokenBalance, err = client.GetTokenBalance(tokenType); err != nil {
			if strings.Contains(err.Error(), "no such token balance") {
				fmt.Println("Your account is not created in the TestNet, please apply tokens from our faucet first.")
			} else {
				ConsoleLog.WithError(err).Error("get token balance failed")
				SetExitStatus(1)
				return
			}
		}
		fmt.Printf("%s balance is: %d\n", tokenType, tokenBalance)
	}
}
