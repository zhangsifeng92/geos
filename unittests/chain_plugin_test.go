package unittests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhangsifeng92/geos/chain/abi_serializer"
	"github.com/zhangsifeng92/geos/chain/types"
	. "github.com/zhangsifeng92/geos/chain/types/generated_containers"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/entity"
	. "github.com/zhangsifeng92/geos/exception"
	. "github.com/zhangsifeng92/geos/exception/try"
	"github.com/zhangsifeng92/geos/log"
	"github.com/zhangsifeng92/geos/plugins/chain_plugin"
	"github.com/zhangsifeng92/geos/unittests/test_contracts"
	"testing"
)

func TestGetBlockWithInvalidAbi(t *testing.T) {
	tester := NewValidatingTesterTrustedProducers(NewAccountNameSet())
	defer tester.close()
	Try(func() {

		tester.ProduceBlocks(2, false)

		tester.CreateAccounts([]common.AccountName{common.N("asserter")}, false, true)
		tester.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

		// setup contract and abi
		tester.SetCode(common.N("asserter"), wast2wasm(test_contracts.AsserterWast), nil)
		tester.SetAbi(common.N("asserter"), test_contracts.AsserterAbi, nil)

		tester.ProduceBlocks(1, false)

		resolver := func(name common.AccountName) (r *abi_serializer.AbiSerializer) {
			Try(func() {
				accnt := entity.AccountObject{Name: name}
				if err := tester.Control.DataBase().Find("byName", accnt, &accnt); err != nil {
					EosThrow(&DatabaseException{}, err.Error())
				}
				var abi abi_serializer.AbiDef
				if abi_serializer.ToABI(accnt.Abi, &abi) {
					r = abi_serializer.NewAbiSerializer(&abi, tester.AbiSerializerMaxTime)
					return
				}
				r = nil
			}).FcRethrowExceptions(log.LvlError, "resolver failed at chain_plugin_tests::abi_invalid_type")
			return r
		}

		// abi should be resolved
		r := resolver(common.N("asserter"))
		assert.NotNil(t, r != nil)

		//prettyTrx := common.Variants{
		//	"actions": []common.Variants{{
		//		"account": "asserter",
		//		"name":    "procassert",
		//		"authorization": common.Variants{
		//			"actor":     "asserter",
		//			"permisson": common.DefaultConfig.ActiveName.String(),
		//		},
		//		"data": common.Variants{
		//			"condition": 1,
		//			"message":   "Should Not Assert!",
		//		},
		//	},
		//	}}

		trx := types.SignedTransaction{}
		trx.Actions = append(trx.Actions, &types.Action{
			Account: common.N("asserter"),
			Name:    common.N("procassert"),
			Authorization: []common.PermissionLevel{
				{
					Actor:      common.N("asserter"),
					Permission: common.DefaultConfig.ActiveName,
				},
			},
			Data: nil,
		})

		trx.Actions[0].Data = r.VariantToBinary("assertdef", &common.Variants{
			"condition": 1,
			"message":   "Should Not Assert!",
		}, tester.AbiSerializerMaxTime)

		//err := common.FromVariant(prettyTrx, &trx)
		//assert.NoError(t, err)
		tester.SetTransactionHeaders(&trx.Transaction, tester.DefaultExpirationDelta, 0)
		priKey, chainId := tester.getPrivateKey(common.N("asserter"), "active"), tester.Control.GetChainId()
		trx.Sign(&priKey, &chainId)
		tester.PushTransaction(&trx, common.MaxTimePoint(), tester.DefaultBilledCpuTimeUs)
		tester.ProduceBlocks(1, false)

		// retrieve block num
		headNum := tester.Control.HeadBlockNum()
		headNumStr := fmt.Sprintf("%d", headNum)
		param := chain_plugin.GetBlockParams{headNumStr}
		plugin := chain_plugin.NewReadOnly(tester.Control, common.MaxMicroseconds())

		// block should be decoded successfully
		blockStr, err := json.Marshal(plugin.GetBlock(param))
		fmt.Println(string(blockStr))
		assert.NoError(t, err)

		// block should be decoded successfully
		assert.Equal(t, true, bytes.Contains(blockStr, []byte("procassert")))
		//TODO show data with hex_data
		//assert.Equal(t, true, bytes.Contains(blockStr, []byte("condition")))
		//assert.Equal(t, true, bytes.Contains(blockStr, []byte("Should Not Assert!")))
		assert.Equal(t, true, bytes.Contains(blockStr, []byte("011253686f756c64204e6f742041737365727421"))) //action data

		// set an invalid abi (int8->xxxx)
		abi2 := test_contracts.AsserterAbi
		pos := bytes.Index(abi2, []byte("int8"))
		assert.Equal(t, true, pos > 0)
		copy(abi2[pos:pos+4], []byte("xxxx"))
		tester.SetAbi(common.N("asserter"), abi2, nil)
		tester.ProduceBlocks(1, false)

		// resolving the invalid abi result in exception
		//TODO check abi type
		//CheckThrowException(t, &InvalidTypeInsideAbi{}, func() { resolver(common.N("asserter")) })

		// get the same block as string, results in decode failed(invalid abi) but not exception
		blockStr2, err := json.Marshal(plugin.GetBlock(param))
		assert.Equal(t, true, bytes.Contains(blockStr2, []byte("procassert")))
		//TODO show data with hex_data
		//assert.Equal(t, false, bytes.Contains(blockStr2, []byte("condition")))
		//assert.Equal(t, false, bytes.Contains(blockStr2, []byte("Should Not Assert!")))
		assert.Equal(t, true, bytes.Contains(blockStr2, []byte("011253686f756c64204e6f742041737365727421"))) //action data

	}).Catch(func(e Exception) {
		t.Fatal(e.DetailMessage())
	}) // get_block_with_invalid_abi
}
