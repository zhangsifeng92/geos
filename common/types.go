package common

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/ecc"
	"math"
	"strconv"
	"strings"
	"time"
)

// For reference:
// https://github.com/mithrilcoin-io/EosCommander/blob/master/app/src/main/java/io/mithrilcoin/eoscommander/data/remote/model/types/EosByteWriter.java
type ChainIDType [4]uint64
type NodeIDType [4]uint64
type BlockIDType [4]uint64
type TransactionIDType [4]uint64
type CheckSum256Type [4]uint64
type Sha256 [4]uint64

func DecodeIDTypeString(str string) (id [4]uint64, err error) {
	b, err := hex.DecodeString(str)
	if err != nil {
		return
	}

	fmt.Println(b)
	for i := range id {
		id[i] = binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
	}

	return
}
func DecodeIDTypeByte(b []byte) (id [4]uint64, err error) {
	for i := range id {
		id[i] = binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
	}

	return id, nil
}
func (n ChainIDType) MarshalJSON() ([]byte, error) {
	b := make([]byte, 32)
	for i := range n {
		binary.LittleEndian.PutUint64(b[i*8:(i+1)*8], n[i])
	}
	return json.Marshal(hex.EncodeToString(b))
}

func (n *ChainIDType) UnmarshalJSON(data []byte) error {
	fmt.Println(47)
	fmt.Println(data)
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(s)

	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	fmt.Println(b)
	for i := range n {
		n[i] = binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
	}

	return nil
}

func (n NodeIDType) MarshalJSON() ([]byte, error) {
	b := make([]byte, 32)
	for i := range n {
		binary.LittleEndian.PutUint64(b[i*8:(i+1)*8], n[i])
	}
	return json.Marshal(hex.EncodeToString(b))
}
func (n Sha256) MarshalJSON() ([]byte, error) {
	b := make([]byte, 32)
	for i := range n {
		binary.LittleEndian.PutUint64(b[i*8:(i+1)*8], n[i])
	}
	return json.Marshal(hex.EncodeToString(b))
}
func (n BlockIDType) MarshalJSON() ([]byte, error) {
	b := make([]byte, 32)
	for i := range n {
		binary.LittleEndian.PutUint64(b[i*8:(i+1)*8], n[i])
	}
	return json.Marshal(hex.EncodeToString(b))
}

func (n *BlockIDType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	for i := range n {
		n[i] = binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
	}
	return nil
}

func (n TransactionIDType) MarshalJSON() ([]byte, error) {
	b := make([]byte, 32)
	for i := range n {
		binary.LittleEndian.PutUint64(b[i*8:(i+1)*8], n[i])
	}
	return json.Marshal(hex.EncodeToString(b))
}
func (n CheckSum256Type) MarshalJSON() ([]byte, error) {
	b := make([]byte, 32)
	for i := range n {
		binary.LittleEndian.PutUint64(b[i*8:(i+1)*8], n[i])
	}
	return json.Marshal(hex.EncodeToString(b))
}

type Name uint64
type AccountName uint64
type PermissionName uint64
type ActionName uint64
type TableName uint64
type ScopeName uint64

func (n AccountName) MarshalJSON() ([]byte, error) {
	return json.Marshal(NameToString(uint64(n)))
}

func (n *AccountName) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*n = AccountName(StringToName(s))
	return nil
}

func (n Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(NameToString(uint64(n)))
}
func (n PermissionName) MarshalJSON() ([]byte, error) {
	return json.Marshal(NameToString(uint64(n)))
}
func (n ActionName) MarshalJSON() ([]byte, error) {
	return json.Marshal(NameToString(uint64(n)))
}
func (n TableName) MarshalJSON() ([]byte, error) {
	return json.Marshal(NameToString(uint64(n)))
}
func (n ScopeName) MarshalJSON() ([]byte, error) {
	return json.Marshal(NameToString(uint64(n)))
}

type AccountResourceLimit struct {
	Used      JSONInt64 `json:"used"`
	Available JSONInt64 `json:"available"`
	Max       JSONInt64 `json:"max"`
}

type DelegatedBandwidth struct {
	From      AccountName `json:"from"`
	To        AccountName `json:"to"`
	NetWeight Asset       `json:"net_weight"`
	CPUWeight Asset       `json:"cpu_weight"`
}

type TotalResources struct {
	Owner     AccountName `json:"owner"`
	NetWeight Asset       `json:"net_weight"`
	CPUWeight Asset       `json:"cpu_weight"`
	RAMBytes  JSONInt64   `json:"ram_bytes"`
}

type VoterInfo struct {
	Owner             AccountName   `json:"owner"`
	Proxy             AccountName   `json:"proxy"`
	Producers         []AccountName `json:"producers"`
	Staked            JSONInt64     `json:"staked"`
	LastVoteWeight    JSONFloat64   `json:"last_vote_weight"`
	ProxiedVoteWeight JSONFloat64   `json:"proxied_vote_weight"`
	IsProxy           byte          `json:"is_proxy"`
}

type RefundRequest struct {
	Owner       AccountName `json:"owner"`
	RequestTime JSONTime    `json:"request_time"` //         {"name":"request_time", "type":"time_point_sec"},
	NetAmount   Asset       `json:"net_amount"`
	CPUAmount   Asset       `json:"cpu_amount"`
}

type CompressionType uint8

const (
	CompressionNone = CompressionType(iota)
	CompressionZlib
)

func (c CompressionType) String() string {
	switch c {
	case CompressionNone:
		return "none"
	case CompressionZlib:
		return "zlib"
	default:
		return ""
	}
}

func (c CompressionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *CompressionType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case "zlib":
		*c = CompressionZlib
	default:
		*c = CompressionNone
	}
	return nil
}

// CurrencyName

type CurrencyName string

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	var num int
	err := json.Unmarshal(data, &num)
	if err == nil {
		*b = Bool(num != 0)
		return nil
	}

	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err != nil {
		return fmt.Errorf("couldn't unmarshal bool as int or true/false: %s", err)
	}

	*b = Bool(boolVal)
	return nil
}

// Asset

// NOTE: there's also ExtendedAsset which is a quantity with the attached contract (AccountName)
type Asset struct {
	Amount int64
	Symbol
}

func (a Asset) Add(other Asset) Asset {
	if a.Symbol != other.Symbol {
		panic("Add applies only to assets with the same symbol")
	}
	return Asset{Amount: a.Amount + other.Amount, Symbol: a.Symbol}
}

func (a Asset) Sub(other Asset) Asset {
	if a.Symbol != other.Symbol {
		panic("Sub applies only to assets with the same symbol")
	}
	return Asset{Amount: a.Amount - other.Amount, Symbol: a.Symbol}
}

func (a Asset) String() string {
	strInt := fmt.Sprintf("%d", a.Amount)
	if len(strInt) < int(a.Symbol.Precision+1) {
		// prepend `0` for the difference:
		strInt = strings.Repeat("0", int(a.Symbol.Precision+uint8(1))-len(strInt)) + strInt
	}

	var result string
	if a.Symbol.Precision == 0 {
		result = strInt
	} else {
		result = strInt[:len(strInt)-int(a.Symbol.Precision)] + "." + strInt[len(strInt)-int(a.Symbol.Precision):]
	}

	return fmt.Sprintf("%s %s", result, a.Symbol.Symbol)
}

// NOTE: there's also a new ExtendedSymbol (which includes the contract (as AccountName) on which it is)
type Symbol struct {
	Precision uint8
	Symbol    string
}

// EOSSymbol represents the standard EOS symbol on the chain.  It's
// here just to speed up things.
var EOSSymbol = Symbol{Precision: 4, Symbol: "EOS"}

func NewEOSAssetFromString(amount string) (out Asset, err error) {
	if len(amount) == 0 {
		return out, fmt.Errorf("cannot be an empty string")
	}

	if strings.Contains(amount, " EOS") {
		amount = strings.Replace(amount, " EOS", "", 1)
	}
	if !strings.Contains(amount, ".") {
		val, err := strconv.ParseInt(amount, 10, 64)
		if err != nil {
			return out, err
		}
		return NewEOSAsset(val * 10000), nil
	}

	parts := strings.Split(amount, ".")
	if len(parts) != 2 {
		return out, fmt.Errorf("cannot have two . in amount")
	}

	if len(parts[1]) > 4 {
		return out, fmt.Errorf("EOS has only 4 decimals")
	}

	val, err := strconv.ParseInt(strings.Replace(amount, ".", "", 1), 10, 64)
	if err != nil {
		return out, err
	}
	return NewEOSAsset(val * int64(math.Pow10(4-len(parts[1])))), nil
}

func NewEOSAsset(amount int64) Asset {
	return Asset{Amount: amount, Symbol: EOSSymbol}
}

// NewAsset parses a string like `1000.0000 EOS` into a properly setup Asset
func NewAsset(in string) (out Asset, err error) {
	sec := strings.SplitN(in, " ", 2)
	if len(sec) != 2 {
		return out, fmt.Errorf("invalid format %q, expected an amount and a currency symbol", in)
	}

	if len(sec[1]) > 7 {
		return out, fmt.Errorf("currency symbol %q too long", sec[1])
	}

	out.Symbol.Symbol = sec[1]
	amount := sec[0]
	amountSec := strings.SplitN(amount, ".", 2)

	if len(amountSec) == 2 {
		out.Symbol.Precision = uint8(len(amountSec[1]))
	}

	val, err := strconv.ParseInt(strings.Replace(amount, ".", "", 1), 10, 64)
	if err != nil {
		return out, err
	}

	out.Amount = val

	return
}

func (a *Asset) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	asset, err := NewAsset(s)
	if err != nil {
		return err
	}

	*a = asset

	return nil
}

func (a Asset) MarshalJSON() (data []byte, err error) {
	return json.Marshal(a.String())
}

type Permission struct {
	PermName     string    `json:"perm_name"`
	Parent       string    `json:"parent"`
	RequiredAuth Authority `json:"required_auth"`
}

type PermissionLevel struct {
	Actor      AccountName    `json:"actor"`
	Permission PermissionName `json:"permission"`
}

// NewPermissionLevel parses strings like `account@active`,
// `otheraccount@owner` and builds a PermissionLevel struct. It
// validates that there is a single optional @ (where permission
// defaults to 'active'), and validates length of account and
// permission names.
func NewPermissionLevel(in string) (out PermissionLevel, err error) {
	parts := strings.Split(in, "@")
	if len(parts) > 2 {
		return out, fmt.Errorf("permission %q invalid, use account[@permission]", in)
	}

	if len(parts[0]) > 12 {
		return out, fmt.Errorf("account name %q too long", parts[0])
	}

	out.Actor = AccountName(StringToName(parts[0]))
	out.Permission = PermissionName(StringToName("active"))

	if len(parts) == 2 {
		if len(parts[1]) > 12 {
			return out, fmt.Errorf("permission %q name too long", parts[1])
		}

		out.Permission = PermissionName(StringToName("active"))
	}

	return
}

type PermissionLevelWeight struct {
	Permission PermissionLevel `json:"permission"`
	Weight     uint16          `json:"weight"` // weight_type
}

type Authority struct {
	Threshold uint32                  `json:"threshold"`
	Keys      []KeyWeight             `json:"keys,omitempty"`
	Accounts  []PermissionLevelWeight `json:"accounts,omitempty"`
	Waits     []WaitWeight            `json:"waits,omitempty"`
}

type KeyWeight struct {
	PublicKey ecc.PublicKey `json:"key"`
	Weight    uint16        `json:"weight"` // weight_type
}

type WaitWeight struct {
	WaitSec uint32 `json:"wait_sec"`
	Weight  uint16 `json:"weight"` // weight_type
}

// type GetCodeResp struct {
// 	AccountName AccountName `json:"account_name"`
// 	CodeHash    string      `json:"code_hash"`
// 	WASM        string      `json:"wasm"`
// 	ABI         ABI         `json:"abi"`
// }

// type GetABIResp struct {
// 	AccountName AccountName `json:"account_name"`
// 	ABI         ABI         `json:"abi"`
// }

// JSONTime

type JSONTime struct {
	time.Time
}

const JSONTimeFormat = "2006-01-02T15:04:05"

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", t.Format(JSONTimeFormat))), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}

	t.Time, err = time.Parse(`"`+JSONTimeFormat+`"`, string(data))
	return err
}

// HexBytes

type HexBytes []byte

func (t HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(t))
}

func (t *HexBytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = hex.DecodeString(s)
	return
}

// SHA256Bytes

type SHA256Bytes []byte // should always be 32 bytes

func (t SHA256Bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(t))
}

func (t *SHA256Bytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = hex.DecodeString(s)
	return
}

type Varuint32 uint32

type Tstamp uint64

const tstampFormate = "2006-01-02T15:04:05.000000"

func (t Tstamp) MarshalJSON() ([]byte, error) {
	slot := int64(t) * 1000
	tm := time.Unix(0, slot).UTC()
	return []byte(fmt.Sprintf("%q", tm.Format(tstampFormate))), nil
}
func (t *Tstamp) UnmarshalJSON(data []byte) (err error) {
	var unixNano int64
	if data[0] == '"' {
		var s string
		if err = json.Unmarshal(data, &s); err != nil {
			return
		}

		unixNano, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

	} else {
		unixNano, err = strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
	}
	*t = Tstamp(unixNano / 1000)

	return nil
}

type BlockTimeStamp uint32

const blockTimestampFormat = "2006-01-02T15:04:05.000"

func (t *BlockTimeStamp) Next() BlockTimeStamp {
	result := NewBlockTimeStamp(t.ToTimePoint())
	result += 1
	return result
}

func (t *BlockTimeStamp) ToTimePoint() time.Time {
	msec := int64(*t) * int64(BlockIntervalMs)
	msec += int64(BlockTimestampEpochMs)
	return time.Unix(0, msec*1e6)
}

func MaxBlockTime() BlockTimeStamp {
	return BlockTimeStamp(0xffff)
}

func MinBlockTime() BlockTimeStamp {
	return BlockTimeStamp(0)
}

func NewBlockTimeStamp(t time.Time) BlockTimeStamp {
	msecSinceEpoch := uint64(t.UnixNano() / 1e6)
	bt := (msecSinceEpoch - BlockTimestampEpochMs) / uint64(BlockIntervalMs)
	return BlockTimeStamp(bt)
}

func (t BlockTimeStamp) MarshalJSON() ([]byte, error) {

	//slot := int64(t)*block_interval_ms*1000000 + block_timestamp_epoch //为了显示0.5s
	tm := time.Unix(0, t).UTC()

	return []byte(fmt.Sprintf("%q", tm.Format(blockTimestampFormat))), nil
}

func (t *BlockTimeStamp) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}
	var temp time.Time
	temp, err = time.Parse(`"`+blockTimestampFormat+`"`, string(data))
	if err != nil {
		temp, err = time.Parse(`"`+blockTimestampFormat+`Z07:00"`, string(data))
		if err != nil {
			return err
		}
	}
	slot := (temp.UnixNano() - block_timestamp_epoch) / 1000000 / block_interval_ms
	*t = BlockTimeStamp(slot)
	return err
}

// func TimeNow() uint64 { //微妙 s*1000000   from 1970
// 	tm := time.Now().UTC()

// 	dur := tm.UnixNano() / 1000

// 	return uint64(dur)
// }
func TimeNow() Tstamp { //微妙 s*1000000   from 1970
	tm := time.Now().UTC()

	dur := tm.UnixNano() / 1000

	return Tstamp(dur)
}

type JSONFloat64 float64

func (f *JSONFloat64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}

		*f = JSONFloat64(val)

		return nil
	}

	var fl float64
	if err := json.Unmarshal(data, &fl); err != nil {
		return err
	}

	*f = JSONFloat64(fl)

	return nil
}

type JSONInt64 int64

func (i *JSONInt64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

		*i = JSONInt64(val)

		return nil
	}

	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = JSONInt64(v)

	return nil
}