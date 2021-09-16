package testing

import (
	"github.com/zhangsifeng92/geos/chain/types"
	. "github.com/zhangsifeng92/geos/chain/types/generated_containers"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto"
	"github.com/zhangsifeng92/geos/crypto/ecc"
)

type ChainTester struct {
	Control  *Controller
	KeyPairs map[common.AccountName]common.Pair //[]<pubKey, priKey>
}

func NewChainTester(when types.BlockTimeStamp, names ...common.AccountName) *ChainTester {
	tester := new(ChainTester)
	priKey, _ := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	pubKey := priKey.PublicKey()

	tester.KeyPairs = make(map[common.AccountName]common.Pair)
	tester.KeyPairs[common.N("eosio")] = common.MakePair(pubKey, priKey)
	tester.KeyPairs[common.N("yuanc")] = common.MakePair(pubKey, priKey)

	hbs := tester.NewHeaderStateTester(when)
	sbk := tester.NewSignedBlockTester(hbs)
	sch := tester.NewProducerScheduleTester(names...)

	tester.Control = newController()
	tester.Control.head = types.NewBlockState(hbs)
	tester.Control.head.SignedBlock = sbk

	tester.Control.head.ActiveSchedule = sch
	tester.Control.head.PendingSchedule = sch
	tester.Control.head.PendingScheduleHash = *crypto.Hash256(sch)

	tester.Control.forkDb.add(tester.Control.head)

	return tester
}

func (t *ChainTester) NewProducerScheduleTester(names ...common.AccountName) types.ProducerScheduleType {
	if len(names) == 0 {
		names = append(names, common.N("eosio"))
	}

	initSchedule := types.ProducerScheduleType{Version: 0, Producers: []types.ProducerKey{}}

	for _, n := range names {
		pk := types.ProducerKey{ProducerName: n, BlockSigningKey: t.KeyPairs[n].First.(ecc.PublicKey)}
		initSchedule.Producers = append(initSchedule.Producers, pk)
	}

	return initSchedule
}

func (t *ChainTester) NewSignedBlockTester(bhs *types.BlockHeaderState) *types.SignedBlock {
	genSigned := new(types.SignedBlock)
	genSigned.SignedBlockHeader = bhs.Header
	return genSigned
}

func (t *ChainTester) NewHeaderStateTester(when types.BlockTimeStamp) *types.BlockHeaderState {
	if when == 0 {
		when = types.NewBlockTimeStamp(common.Now())
	}
	genHeader := new(types.BlockHeaderState)
	genHeader.Header.Timestamp = when
	genHeader.Header.Confirmed = 1
	genHeader.BlockId = genHeader.Header.BlockID()
	genHeader.BlockNum = genHeader.Header.BlockNumber()
	genHeader.ProducerToLastProduced = *NewAccountNameUint32Map()
	genHeader.ProducerToLastImpliedIrb = *NewAccountNameUint32Map()

	return genHeader
}
