package types

import (
	"github.com/zhangsifeng92/geos/chain/types/generated_containers"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto"
	"github.com/zhangsifeng92/geos/crypto/rlp"
	. "github.com/zhangsifeng92/geos/exception"
	. "github.com/zhangsifeng92/geos/exception/try"
)

// Action
type Action struct {
	Account       common.AccountName       `json:"account"`
	Name          common.ActionName        `json:"name"`
	Authorization []common.PermissionLevel `json:"authorization,omitempty"`
	Data          common.HexBytes          `json:"data"`
}

func (a Action) DataAs(t interface{}) {
	err := rlp.DecodeBytes(a.Data, t)
	if err != nil {
		EosThrow(&ParseErrorException{}, "action data parse error: %s", err.Error())
	}
}

type ContractTypesInterface interface {
	GetAccount() common.AccountName
	GetName() common.ActionName
}

//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/common/container/treeset" AccountNameSet(common.AccountName,common.CompareName,false)
//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/common/container/treemap" AccountNameUint64Map(common.AccountName,uint64,common.CompareName,false)
type ActionReceipt struct {
	Receiver       common.AccountName             `json:"receiver"`
	ActDigest      crypto.Sha256                  `json:"act_digest"`
	GlobalSequence uint64                         `json:"global_sequence"`
	RecvSequence   uint64                         `json:"recv_sequence"`
	AuthSequence   generated.AccountNameUint64Map `json:"auth_sequence"`
	CodeSequence   common.Vuint32                 `json:"code_sequence"`
	AbiSequence    common.Vuint32                 `json:"abi_sequence"`
}

func (a *ActionReceipt) Digest() crypto.Sha256 {
	return *crypto.Hash256(a)
}
