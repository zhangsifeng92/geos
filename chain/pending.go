package chain

import (
	"github.com/zhangsifeng92/geos/chain/types"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/database"
)

type PendingState struct {
	//MaybeSession      *database.Session `json:"db_session"`
	DbSession         *MaybeSession         `json:"db_session"`
	PendingBlockState *types.BlockState     `json:"pending_block_state"`
	Actions           []types.ActionReceipt `json:"actions"`
	BlockStatus       types.BlockStatus     `json:"block_status"`
	ProducerBlockId   common.BlockIdType
	PendingValid      bool
}

func NewPendingState(db database.DataBase) *PendingState {
	pending := PendingState{}
	pending.DbSession = NewMaybeSession(db)
	pending.BlockStatus = types.Incomplete
	pending.PendingValid = false
	return &pending
}

func NewDefaultPendingState() *PendingState {
	return &PendingState{}
}
func (p *PendingState) Reset() *PendingState {
	p.DbSession.Undo()
	return nil
}

func (p *PendingState) Push() {
	if p.DbSession != nil {
		p.DbSession.Push()
	}
}
