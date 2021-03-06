package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/rocksideio/rockside-sdk-go"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Rockside client",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return cmd.Usage()
	},
	SilenceUsage: true,
}

var (
	eoaCmd = &cobra.Command{
		Use:   "eoa",
		Short: "Manage EOA",
	}

	listEOACmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List EOA",
		RunE: func(cmd *cobra.Command, args []string) error {
			eoaList, err := RocksideClient().EOA.List()
			if err != nil {
				return err
			}

			return printJSON(eoaList)
		},
	}

	createEOACmd = &cobra.Command{
		Use:   "create",
		Short: "Create an EOA",
		RunE: func(cmd *cobra.Command, args []string) error {
			eoa, err := RocksideClient().EOA.Create()
			if err != nil {
				return err
			}

			return printJSON(eoa)
		},
	}
)

var (
	tokensCmd = &cobra.Command{
		Use:   "tokens",
		Short: "Manage Tokens",
	}

	createTokenCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a Token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing domain")
			}
			domain := args[0]
			contracts := []string{}

			for i := 1; i < len(args); i++ {
				contracts = append(contracts, args[i])
			}

			token, err := RocksideClient().Tokens.Create(domain, contracts)
			if err != nil {
				return err
			}

			return printJSON(token)
		},
	}
)

var (
	smartWalletsCmd = &cobra.Command{
		Use:   "smartwallets",
		Short: "Manage smart wallets",
	}

	listSmartWalletsCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List smart wallets",
		RunE: func(cmd *cobra.Command, args []string) error {
			smartWallets, err := RocksideClient().SmartWallets.List()
			if err != nil {
				return err
			}

			return printJSON(smartWallets)
		},
	}

	createSmartWalletCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy a smart wallet given the account address and forwarder address",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing public address of the account and/or forwarder address")
			}

			smartWallet, err := RocksideClient().SmartWallets.Create(args[0], args[1])
			if err != nil {
				return err
			}

			return printJSON(smartWallet)
		},
	}
)

var (
	forwarderCmd = &cobra.Command{
		Use:   "forwarder",
		Short: "Manage forwarders",
	}

	deployForwarder = &cobra.Command{
		Use:   "deploy",
		Short: "deploy a forwarder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing public address of the account")
			}

			forwarder, err := RocksideClient().Forwarder.Create(args[0])
			if err != nil {
				return err
			}

			return printJSON(forwarder)
		},
	}

	getNonceCmd = &cobra.Command{
		Use:   "nonce",
		Short: "get nonce of a smart wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing contract address as first param and account address as second param")
			}
			contractAddress := args[0]
			accountAddress := args[1]

			nonce, err := RocksideClient().Forwarder.GetRelayParams(contractAddress, accountAddress)
			if err != nil {
				return err
			}

			return printJSON(nonce)
		},
	}

	signCmd = &cobra.Command{
		Use:   "sign",
		Short: "sign transaction and parameters to be relayed",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing contract address and transaction payload {\"from\":\"\",\"to\":\"\", \"value\":\"\", \"data\":\"\" }")
			}

			contractAddress := args[0]
			txJSON := args[1]
			tx := &rockside.Transaction{}
			if err := json.Unmarshal([]byte(txJSON), tx); err != nil {
				return err
			}

			signResponse, err := RocksideClient().Forwarder.SignTxParams(privateKeyFlag, contractAddress, tx.From, tx.To, tx.Data, tx.Nonce)

			if err != nil {
				return err
			}

			return printJSON(signResponse)
		},
	}

	relayCmd = &cobra.Command{
		Use:   "relay",
		Short: "relay transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("missing contract address and transaction payload {\"from\":\"\", \"to\":\"\", \"value\":\"\", \"speed\":\"\", \"gas_price_limit\":\"\", \"data\":\"\", \"signature\":\"\"}")
			}

			contractAddress := args[0]
			txJSON := args[1]
			relayTx := &rockside.RelayExecuteTxRequest{}
			if err := json.Unmarshal([]byte(txJSON), relayTx); err != nil {
				return err
			}

			relayResponse, err := RocksideClient().Forwarder.Relay(contractAddress, *relayTx)
			if err != nil {
				return err
			}

			return printJSON(relayResponse)
		},
	}
)

var (
	transactionCmd = &cobra.Command{
		Use:   "transaction",
		Short: "Manage transaction",
	}

	sentTxCmd = &cobra.Command{
		Use:   "send",
		Short: "send transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing transaction payload {\"from\":\"\",\"to\":\"\", \"value\":\"\", gas:\"\", \"gasPrice\":\"\", \"nonce\":\"\"}")
			}

			txJSON := args[0]
			tx := &rockside.Transaction{}
			if err := json.Unmarshal([]byte(txJSON), tx); err != nil {
				return err
			}

			txResponse, err := RocksideClient().Transaction.Send(*tx)
			if err != nil {
				return err
			}

			return printJSON(txResponse)
		},
	}

	showTxCmd = &cobra.Command{
		Use:   "show",
		Short: "show transaction given a tx hash or tracking ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing transaction hash or tracking ID")
			}
			result, err := RocksideClient().Transaction.Show(args[0])
			if err != nil {
				return err
			}

			return printJSON(result)
		},
	}
)

func printJSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	return enc.Encode(v)
}
