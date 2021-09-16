package include_test

import (
	"fmt"
	"github.com/zhangsifeng92/geos/chain/types"
	"github.com/zhangsifeng92/geos/libraries/asio"
	. "github.com/zhangsifeng92/geos/plugins/appbase/app"
	. "github.com/zhangsifeng92/geos/plugins/appbase/app/include"
	. "github.com/zhangsifeng92/geos/plugins/chain_interface"
	"testing"
	"time"
)

//var  acceptedBlockHeader *Signal;
//var  acceptedBlock *Signal;
//var  irreversibleBlock *Signal;
//var  acceptedTransaction *Signal;
//var  appliedTransaction *Signal;
//var  acceptedConfirmation *Signal;
//var  badAlloc *Signal;

type blockAcceptor struct {
}

func (*blockAcceptor) doAccept(s *types.SignedBlock) {
	fmt.Println(s.Timestamp)
}

func doAccept(s *types.SignedBlock) {
	fmt.Println(s.Timestamp)
}

func (*blockAcceptor) doRejectedBlockFunc(s *types.SignedBlock) {
	fmt.Println(s.Timestamp)
}

func Test_Channel(t *testing.T) {

	//subscribe
	App().GetChannel(PreAcceptedBlock).Subscribe(&PreAcceptedBlockCaller{doAccept})
	App().GetChannel(PreAcceptedBlock).Subscribe(&PreAcceptedBlockCaller{new(blockAcceptor).doAccept})
	rbf := &RejectedBlockCaller{new(blockAcceptor).doRejectedBlockFunc}
	App().GetChannel(RejectedBlock).Subscribe(rbf)

	//call
	sb := &types.SignedBlock{}
	sb.Timestamp = types.NewBlockTimeStamp(100)
	//App().GetChannel(chain.PreAcceptedBlock).Publish(sb)
	App().GetChannel(RejectedBlock).Publish(sb)
	App().GetChannel(PreAcceptedBlock).Publish(sb)
	App().GetChannel(AcceptedBlockHeader).Publish(sb)

	timer := asio.NewDeadlineTimer(App().GetIoService())
	timer.ExpiresFromNow(time.Millisecond)
	timer.AsyncWait(func(err error) {
		App().Quit()
	})

	App().Exec()
}
