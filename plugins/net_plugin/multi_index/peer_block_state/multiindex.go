package peer_block_state

import (
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto"
	"github.com/zhangsifeng92/geos/libraries/container"
	"github.com/zhangsifeng92/geos/plugins/net_plugin/multi_index"
)

//go:generate go install "github.com/zhangsifeng92/geos/libraries/multiindex/"
//go:generate go install "github.com/zhangsifeng92/geos/libraries/multiindex/multi_index_container/..."
//go:generate go install "github.com/zhangsifeng92/geos/libraries/multiindex/ordered_index/..."

//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/multiindex/multi_index_container" PeerBlockStateIndex(ById,ByIdNode,multi_index.PeerBlockState)
func (m *PeerBlockStateIndex) GetById() *ById             { return m.super }
func (m *PeerBlockStateIndex) GetByBlockNum() *ByBlockNum { return m.super.super }

//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/multiindex/ordered_index" ById(PeerBlockStateIndex,PeerBlockStateIndexNode,ByBlockNum,ByBlockNumNode,multi_index.PeerBlockState,common.BlockIdType,ByIdFunc,ByIdCompare,false)
var ByIdFunc = func(n multi_index.PeerBlockState) common.BlockIdType { return n.ID }
var ByIdCompare = crypto.Sha256Compare

//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/multiindex/ordered_index" ByBlockNum(PeerBlockStateIndex,PeerBlockStateIndexNode,PeerBlockStateIndexBase,PeerBlockStateIndexBaseNode,multi_index.PeerBlockState,uint32,ByBlockNumFunc,ByBlockNumCompare,false)
var ByBlockNumFunc = func(n multi_index.PeerBlockState) uint32 { return n.BlockNum }
var ByBlockNumCompare = container.UInt32Comparator

//go:generate go build
