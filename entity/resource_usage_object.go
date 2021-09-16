package entity

import (
	"github.com/zhangsifeng92/geos/chain/types"
	"github.com/zhangsifeng92/geos/common"
)

type ResourceUsageObject struct {
	ID       common.IdType      `multiIndex:"id,increment"`
	Owner    common.AccountName `multiIndex:"byOwner,orderedUnique"`
	NetUsage types.UsageAccumulator
	CpuUsage types.UsageAccumulator
	RamUsage uint64
}
