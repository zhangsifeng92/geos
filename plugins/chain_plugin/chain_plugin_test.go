package chain_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhangsifeng92/geos/chain/types"
	. "github.com/zhangsifeng92/geos/chain/types/generated_containers"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto"
	"github.com/zhangsifeng92/geos/crypto/ecc"
	"testing"
)

func TestApiParams(t *testing.T) {
	//get_currency_balance params
	getCurrencyBalanceParams := GetCurrencyBalanceParams{
		Code:    common.N("eosio.token"),
		Account: common.N("eosio"),
		Symbol:  "SYS",
	}

	var variant common.Variant
	common.ToVariant(getCurrencyBalanceParams, &variant)
	body, err := json.Marshal(variant)
	assert.NoError(t, err)

	var params GetCurrencyBalanceParams
	err = json.Unmarshal(body, &params)
	assert.NoError(t, err)
	assert.Equal(t, common.N("eosio.token"), params.Code)
	assert.Equal(t, common.N("eosio"), params.Account)
	assert.Equal(t, "SYS", params.Symbol)

	//push_transaction params
	action := []byte("\"code\":\"eosio\", \"action\":\"newaccount\", \"data\":\"eosio.token\"")
	prikey, _ := ecc.NewRandomPrivateKey()
	sign, err := prikey.Sign(crypto.Hash256(action).Bytes())
	assert.NoError(t, err)

	packed := types.PackedTransaction{
		Signatures:            []ecc.Signature{sign},
		Compression:           types.CompressionNone,
		PackedContextFreeData: action,
		PackedTrx:             action,
		UnpackedTrx:           nil,
	}

	common.ToVariant(packed, &variant)

	body, err = json.Marshal(variant)
	assert.NoError(t, err)

	var pushTrxParams PushTransactionParams
	err = json.Unmarshal(body, &pushTrxParams)
	assert.NoError(t, err)

	var prettyInput types.PackedTransaction
	common.FromVariant(pushTrxParams, &prettyInput)
	assert.Equal(t, sign, prettyInput.Signatures[0])
	assert.Equal(t, types.CompressionNone, prettyInput.Compression)
	assert.Equal(t, common.HexBytes(action), prettyInput.PackedTrx)
	assert.Equal(t, common.HexBytes(action), prettyInput.PackedContextFreeData)
	assert.Equal(t, (*types.Transaction)(nil), prettyInput.UnpackedTrx)
}

func TestJsonFormat(t *testing.T) {
	res := GetRequiredKeysResult{*NewPublicKeySet()}
	pri1, _ := ecc.NewRandomPrivateKey()
	pri2, _ := ecc.NewRandomPrivateKey()
	res.RequiredKeys.Add(pri1.PublicKey(), pri2.PublicKey())
	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func TestConvertUint64(t *testing.T) {
	a := convertToUint64("123", "a")
	fmt.Println(a)
	b := convertToUint64("eosio", "b")
	fmt.Println(b)
	c := convertToUint64("sys", "c")
	fmt.Println(c)
}
