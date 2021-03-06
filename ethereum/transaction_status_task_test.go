// +build unit

package ethereum

import (
	"testing"

	"github.com/centrifuge/go-centrifuge/testingutils/identity"

	"github.com/centrifuge/go-centrifuge/transactions"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/stretchr/testify/assert"
)

func TestMintingConfirmationTask_ParseKwargs_success(t *testing.T) {
	task := TransactionStatusTask{}
	txHash := "0xd18036d7c1fe109af377e8ce1d9096e69a5df0741fba7e4f3507f8e6aa573515"
	txID := transactions.NewTxID().String()
	cid := testingidentity.GenerateRandomDID()

	kwargs := map[string]interface{}{
		transactions.TxIDParam:  txID,
		TransactionAccountParam: cid.String(),
		TransactionTxHashParam:  txHash,
	}

	decoded, err := utils.SimulateJSONDecodeForGocelery(kwargs)
	assert.Nil(t, err, "json decode should not thrown an error")
	err = task.ParseKwargs(decoded)
	assert.Nil(t, err, "parsing should be successful")

	assert.Equal(t, cid, task.accountID, "accountID should be parsed correctly")
	assert.Equal(t, txID, task.TxID.String(), "txID should be parsed correctly")
	assert.Equal(t, txHash, task.txHash, "txHash should be parsed correctly")

}

func TestMintingConfirmationTask_ParseKwargs_fail(t *testing.T) {
	task := TransactionStatusTask{}
	tests := []map[string]interface{}{
		{
			transactions.TxIDParam:  transactions.NewTxID().String(),
			TransactionAccountParam: testingidentity.GenerateRandomDID().String(),
		},
		{
			TransactionAccountParam: testingidentity.GenerateRandomDID().String(),
			TransactionTxHashParam:  "0xd18036d7c1fe109af377e8ce1d9096e69a5df0741fba7e4f3507f8e6aa573515",
		},
		{
			transactions.TxIDParam: transactions.NewTxID().String(),
			TransactionTxHashParam: "0xd18036d7c1fe109af377e8ce1d9096e69a5df0741fba7e4f3507f8e6aa573515",
		},
		{
			//empty map

		},
		{
			"dummy": "dummy",
		},
	}

	for i, test := range tests {
		decoded, err := utils.SimulateJSONDecodeForGocelery(test)
		assert.Nil(t, err, "json decode should not thrown an error")
		err = task.ParseKwargs(decoded)
		assert.Error(t, err, "test case %v: parsing should fail", i)
	}
}
