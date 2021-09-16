package transaction_state

import (
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto"
	"github.com/zhangsifeng92/geos/libraries/container"
	"github.com/zhangsifeng92/geos/plugins/net_plugin/multi_index"
)

//go:generate go install "github.com/zhangsifeng92/geos/libraries/multiindex/"
//go:generate go install "github.com/zhangsifeng92/geos/libraries/multiindex/multi_index_container/..."
//go:generate go install "github.com/zhangsifeng92/geos/libraries/multiindex/ordered_index/..."

//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/multiindex/multi_index_container" TransactionStateIndex(ById,ByIdNode,multi_index.TransactionState)
func (m *TransactionStateIndex) GetById() *ById             { return m.super }
func (m *TransactionStateIndex) GetByExpiry() *ByExpiry     { return m.super.super }
func (m *TransactionStateIndex) GetByBlockNum() *ByBlockNum { return m.super.super.super }

//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/multiindex/ordered_index" ById(TransactionStateIndex,TransactionStateIndexNode,ByExpiry,ByExpiryNode,multi_index.TransactionState,common.TransactionIdType,ByIdFunc,ByIdCompare,false)
var ByIdFunc = func(n multi_index.TransactionState) common.TransactionIdType { return n.ID }
var ByIdCompare = crypto.Sha256Compare

//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/multiindex/ordered_index" ByExpiry(TransactionStateIndex,TransactionStateIndexNode,ByBlockNum,ByBlockNumNode,multi_index.TransactionState,common.TimePointSec,ByExpiryFunc,ByExpiryCompare,true)
var ByExpiryFunc = func(n multi_index.TransactionState) common.TimePointSec { return n.Expires }
var ByExpiryCompare = func(a, b common.TimePointSec) int { return container.UInt32Comparator(uint32(a), uint32(b)) }

//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/multiindex/ordered_index" ByBlockNum(TransactionStateIndex,TransactionStateIndexNode,TransactionStateIndexBase,TransactionStateIndexBaseNode,multi_index.TransactionState,uint32,ByBlockNumFunc,ByBlockNumCompare,true)
var ByBlockNumFunc = func(n multi_index.TransactionState) uint32 { return n.BlockNum }
var ByBlockNumCompare = container.UInt32Comparator

//go:generate go build
