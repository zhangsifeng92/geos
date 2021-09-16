package entity

import "github.com/zhangsifeng92/geos/common"

type PermissionUsageObject struct {
	ID       common.IdType `multiIndex:"id,increment"`
	LastUsed common.TimePoint
}
