package rockside_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocksideio/rockside-sdk-go"
)

func Example_gaslessContractTransactor() {
	rocksideSmartWalletAddress := common.HexToAddress("my_rockside_smartwallet_hex_contract_address")
	contractAddress := common.HexToAddress("my_contract_hex_address")

	rocksideTransactor := rockside.NewTransactor(rocksideSmartWalletAddress, rocksideClient)

	// NewContractTransactor is typically in your contract binding GO file which was generated via `abigen`
	contract, err := NewContractTransactor(contractAddress, rocksideTransactor)
	if err != nil {
		panic(err)
	}

	tx, _ := contract.Write(rockside.TransactOpts(), [32]byte{})

	txHash := rocksideTransactor.ReturnRocksideTransactionHash(tx.Hash())
	fmt.Println(txHash)
}
