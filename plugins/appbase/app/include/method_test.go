package include_test

import (
	"fmt"
	"github.com/zhangsifeng92/geos/chain/types"
	. "github.com/zhangsifeng92/geos/plugins/appbase/app"
	. "github.com/zhangsifeng92/geos/plugins/chain_interface"
	"testing"
)

type Gbi struct {
}

func (g Gbi) GetBlock(s *types.SignedBlock) {
	fmt.Println("getBlock")
	fmt.Println(s.Timestamp)
}

func Test_Method(t *testing.T) {
	gbi := App().GetMethod(GetBlockById)

	//register
	gbi.Register(&RejectedBlockCaller{Gbi{}.GetBlock})

	sb := new(types.SignedBlock)
	sb.Timestamp = types.NewBlockTimeStamp(100)
	//CallMethods
	gbi.CallMethods(sb)

}
