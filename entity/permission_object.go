package entity

import (
	"github.com/zhangsifeng92/geos/chain/types"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/database"
)

type PermissionObject struct {
	Parent      common.IdType `multiIndex:"byParent,orderedUnique"`
	ID          common.IdType `multiIndex:"id,increment,byParent,byName"`
	UsageId     common.IdType
	Owner       common.AccountName    `multiIndex:"byOwner,orderedUnique"`
	Name        common.PermissionName `multiIndex:"byOwner,orderedUnique:byName,orderedUnique"`
	LastUpdated common.TimePoint
	Auth        types.SharedAuthority
}

func (po *PermissionObject) Satisfies(other PermissionObject, PermissionIndex *database.MultiIndex) bool {
	if po.Owner != other.Owner {
		return false
	}
	if po.ID == other.ID || po.ID == other.Parent {
		return true
	}
	itr, err := PermissionIndex.LowerBound(PermissionObject{ID: other.Parent})
	if err != nil {
		return false
	}
	parent := PermissionObject{}
	itr.Data(&parent)
	for {
		if po.ID == parent.Parent {
			return true
		}
		if parent.Parent == 0 {
			return false
		}

		itr, err = PermissionIndex.LowerBound(PermissionObject{ID: parent.Parent})
		if err != nil {
			break
		}
		itr.Data(&parent)
	}
	return false
}
