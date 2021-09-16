package entity

import (
	"github.com/zhangsifeng92/geos/common"
)

type BlockSummaryObject struct {
	Id      common.IdType `multiIndex:"id,increment"`
	BlockId common.BlockIdType
}
