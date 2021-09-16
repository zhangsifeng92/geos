package entity

import (
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto/ecc"
)

type PublicKeyHistoryObject struct {
	ID         common.IdType         `multiIndex:"id,increment,byPubKey,byAccountPermission"`
	PublicKey  ecc.PublicKey         `multiIndex:"byPubKey,orderedUnique"`            //c++ publicKey+id unique
	Name       common.AccountName    `multiIndex:"byAccountPermission,orderedUnique"` //c++ ByAccountPermission+id unique
	Permission common.PermissionName `multiIndex:"byAccountPermission,orderedUnique"` //c++ ByAccountPermission+id unique
}
