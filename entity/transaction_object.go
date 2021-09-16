package entity

import "github.com/zhangsifeng92/geos/common"

type TransactionObject struct {
	ID         common.IdType            `multiIndex:"id,increment,byExpiration"`
	Expiration common.TimePointSec      `multiIndex:"byExpiration,orderedUnique"`
	TrxID      common.TransactionIdType `multiIndex:"byTrxId,orderedUnique"`
}
