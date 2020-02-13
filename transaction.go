package rockside

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type transactionEndpoint endpoint

type Transaction struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Value    string `json:"value,omitempty"`
	Data     string `json:"data,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasprice,omitempty"`
}

type SendTxResponse struct {
	TransactionHash string `json:"transaction_hash"`
}

func (t *transactionEndpoint) Send(transaction Transaction) (SendTxResponse, error) {
	var result SendTxResponse

	if err := transaction.ValidateFields(); err != nil {
		return result, err
	}

	path := fmt.Sprintf("ethereum/%s/transaction", t.client.network)
	if _, err := t.client.post(path, transaction, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (t Transaction) ValidateFields() error {
	if !common.IsHexAddress(t.From) {
		return errors.New("invalid 'from' address")
	}
	// To can be empty for contract creation
	if t.To != "" && !common.IsHexAddress(t.To) {
		return errors.New("invalid 'to' address")
	}
	if len(t.Data) > 0 {
		if _, err := hexutil.Decode(t.Data); err != nil {
			return fmt.Errorf("invalid 'data' bytes: %w", err)
		}
	}
	if len(t.Value) > 0 {
		if _, err := hexutil.DecodeBig(t.Value); err != nil {
			return fmt.Errorf("invalid 'value' number: %w", err)
		}
	}
	return nil
}
