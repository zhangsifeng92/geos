package entity

import (
	"github.com/zhangsifeng92/geos/chain/types"
	"github.com/zhangsifeng92/geos/common"
)

var DefaultResourceLimitsConfigObject ResourceLimitsConfigObject

func init() {
	DefaultResourceLimitsConfigObject.ID = 0
}

type ResourceLimitsConfigObject struct {
	ID                           common.IdType `multiIndex:"id,increment"`
	CpuLimitParameters           types.ElasticLimitParameters
	NetLimitParameters           types.ElasticLimitParameters
	AccountCpuUsageAverageWindow uint32
	AccountNetUsageAverageWindow uint32
}

func NewResourceLimitsConfigObject() ResourceLimitsConfigObject {
	config := ResourceLimitsConfigObject{}
	config.CpuLimitParameters = types.ElasticLimitParameters{Target: common.EosPercent(uint64(common.DefaultConfig.MaxBlockCpuUsage), common.DefaultConfig.TargetBlockCpuUsagePct),
		Max:           uint64(common.DefaultConfig.MaxBlockCpuUsage),
		Periods:       uint32(common.DefaultConfig.BlockCpuUsageAverageWindowMs) / uint32(common.DefaultConfig.BlockIntervalMs),
		MaxMultiplier: 1000, ContractRate: types.Ratio{Numerator: 99, Denominator: 100}, ExpandRate: types.Ratio{Numerator: 1000, Denominator: 999},
	}

	config.NetLimitParameters = types.ElasticLimitParameters{Target: common.EosPercent(uint64(common.DefaultConfig.MaxBlockNetUsage), common.DefaultConfig.TargetBlockNetUsagePct),
		Max:           uint64(common.DefaultConfig.MaxBlockNetUsage),
		Periods:       uint32(common.DefaultConfig.BlockSizeAverageWindowMs) / uint32(common.DefaultConfig.BlockIntervalMs),
		MaxMultiplier: 1000, ContractRate: types.Ratio{Numerator: 99, Denominator: 100}, ExpandRate: types.Ratio{Numerator: 1000, Denominator: 999},
	}
	config.AccountCpuUsageAverageWindow = common.DefaultConfig.AccountCpuUsageAverageWindowMs / uint32(common.DefaultConfig.BlockIntervalMs)
	config.AccountNetUsageAverageWindow = common.DefaultConfig.AccountNetUsageAverageWindowMs / uint32(common.DefaultConfig.BlockIntervalMs)
	return config
}
