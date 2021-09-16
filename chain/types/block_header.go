package types

import (
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto"
	"github.com/zhangsifeng92/geos/crypto/ecc"
)

type BlockHeader struct {
	Timestamp        BlockTimeStamp     `json:"timestamp"`
	Producer         common.AccountName `json:"producer"`
	Confirmed        uint16/* default=1 */ `json:"confirmed"`
	Previous         common.BlockIdType          `multiIndex:"byPrevious,orderedUnique" json:"previous"`
	TransactionMRoot common.CheckSum256Type      `json:"transaction_mroot"`
	ActionMRoot      common.CheckSum256Type      `json:"action_mroot"`
	ScheduleVersion  uint32                      `json:"schedule_version"`
	NewProducers     *SharedProducerScheduleType `json:"new_producers" eos:"optional"`
	HeaderExtensions []Extension                 `json:"header_extensions"`
}

func NewBlockHeader() *BlockHeader {
	return &BlockHeader{Confirmed: 1}
}

func (b *BlockHeader) Digest() crypto.Sha256 {
	return *crypto.Hash256(*b)
}

func (b *BlockHeader) BlockNumber() uint32 {
	return NumFromID(&b.Previous) + 1
}

func NumFromID(id *common.BlockIdType) uint32 {
	return common.EndianReverseU32(uint32(id.Hash[0]))
}

func (b *BlockHeader) BlockID() common.BlockIdType {
	// Do not include signed_block_header attributes in id, specifically exclude producer_signature.
	result := b.Digest()
	result.Hash[0] &= 0xffffffff00000000
	result.Hash[0] += uint64(common.EndianReverseU32(b.BlockNumber())) // store the block num in the ID, 160 bits is plenty for the hash
	return common.BlockIdType(result)
}

type SignedBlockHeader struct {
	BlockHeader       `multiIndex:"inline"`
	ProducerSignature ecc.Signature `json:"-"`
}

func NewSignedBlockHeader() *SignedBlockHeader {
	return &SignedBlockHeader{BlockHeader: *NewBlockHeader()}
}

type HeaderConfirmation struct {
	BlockId           common.BlockIdType `json:"block_id"`
	Producer          common.AccountName `json:"producer"`
	ProducerSignature ecc.Signature      `json:"producers_signature"`
}

func (b BlockHeader) IsEmpty() bool {
	return b.Timestamp == 0 && b.Producer.IsEmpty() && b.Confirmed == 0 && b.Previous.IsEmpty() &&
		b.TransactionMRoot.IsEmpty() && b.ActionMRoot.IsEmpty() && b.ScheduleVersion == 0 && b.NewProducers == nil &&
		len(b.HeaderExtensions) == 0
}
func (s SignedBlockHeader) IsEmpty() bool {
	return len(s.ProducerSignature.Content) == 0 && s.BlockHeader.IsEmpty()
}
