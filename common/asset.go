package common

import (
	"encoding/json"
	"fmt"
	"github.com/zhangsifeng92/geos/crypto/rlp"
	. "github.com/zhangsifeng92/geos/exception"
	. "github.com/zhangsifeng92/geos/exception/try"
	"github.com/zhangsifeng92/geos/log"
	"math"
	"strconv"
	"strings"
)

const maxAmount int64 = int64(1)<<62 - 1
const SizeofAsset int = 16

type Asset struct {
	Amount int64 `eos:"asset"`
	Symbol
}

func (s Asset) Pack() (re []byte, err error) {
	re = append(re, WriteInt64(s.Amount)...)
	reSymbol, err := s.Symbol.Pack()
	if err != nil {
		return nil, err
	}
	re = append(re, reSymbol...)
	return re, nil
}

func (s *Asset) Unpack(in []byte) (int, error) {
	decoder := rlp.NewDecoder(in)
	a, err := decoder.ReadInt64()
	if err != nil {
		return 0, err
	}
	s.Amount = a
	l, err := s.Symbol.Unpack(decoder.GetData()[decoder.GetPos():])
	return l + decoder.GetPos(), err

}

func NewAssetWithCheck(a int64, id Symbol) *Asset {
	re := &Asset{
		Amount: a,
		Symbol: id,
	}
	re.assert()
	return re
}
func (a *Asset) assert() {
	EosAssert(a.isAmountWithinRange(), &AssetTypeException{}, "magnitude of asset amount must be less than 2^62")
	EosAssert(a.Symbol.Valid(), &AssetTypeException{}, "invalid symbol")
}

func (a *Asset) isAmountWithinRange() bool {
	return -maxAmount <= a.Amount && a.Amount <= maxAmount
}

func (a *Asset) isValid() bool {
	return a.isAmountWithinRange() && a.Symbol.Valid()
}

func (a Asset) Add(b Asset) Asset {
	EosAssert(a.Symbol == b.Symbol, &AssetTypeException{}, "addition between two different asset is not allowed")
	return Asset{Amount: a.Amount + b.Amount, Symbol: a.Symbol}
}

func (a Asset) Sub(b Asset) Asset {
	EosAssert(a.Symbol == b.Symbol, &AssetTypeException{}, "subtraction between two different asset is not allowed")
	return Asset{Amount: a.Amount - b.Amount, Symbol: a.Symbol}
}

func (a Asset) String() string {
	sign := ""
	abs := a.Amount
	if a.Amount < 0 {
		sign = "-"
		abs = -1 * a.Amount
	}
	strInt := fmt.Sprintf("%d", abs)
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

	return fmt.Sprintf("%s %s", sign+result, a.Symbol.Symbol)
}

func (a Asset) FromString(from *string) Asset {
	spacePos := strings.Index(*from, " ")
	EosAssert(spacePos != -1, &AssetTypeException{}, "Asset's amount and symbol should be separated with space")
	symbolStr := string([]byte(*from)[spacePos+1:])
	amountStr := string([]byte(*from)[:spacePos])

	dotPos := strings.Index(amountStr, ".")
	if dotPos != -1 {
		EosAssert(dotPos != len(amountStr)-1, &AssetTypeException{}, "Missing decimal fraction after decimal point")
	}

	var precisionDigitStr string
	if dotPos != -1 {
		precisionDigitStr = strconv.Itoa(len(amountStr) - dotPos - 1)
	} else {
		precisionDigitStr = "0"
	}

	symbolPart := precisionDigitStr + "," + symbolStr
	sym := Symbol{}.FromString(&symbolPart)

	var intPart, fractPart int64
	if dotPos != -1 {
		intPartString := string([]byte(amountStr)[:dotPos])
		CheckParseInt64(intPartString)
		intPart, _ = strconv.ParseInt(intPartString, 10, 64)
		fractPart, _ = strconv.ParseInt(string([]byte(amountStr)[dotPos+1:]), 10, 64)
		if amountStr[0] == '-' {
			fractPart *= -1
		}
	} else {
		intPartString := amountStr
		CheckParseInt64(intPartString)
		intPart, _ = strconv.ParseInt(intPartString, 10, 64)
	}
	amount := intPart
	for i := uint8(0); i < sym.Precision; i++ {
		amount *= 10
		if intPart >= 0 {
			EosAssert(intPart <= amount, &OverflowException{}, "asset amount overflow")
		} else {
			EosAssert(intPart >= amount, &UnderflowException{}, "asset amount underflow")
		}
	}

	amount += fractPart
	if fractPart > 0 {
		EosAssert(fractPart <= amount, &OverflowException{}, "asset amount overflow")
	} else if fractPart < 0 {
		EosAssert(fractPart >= amount, &UnderflowException{}, "asset amount underflow")
	}
	asset := Asset{Amount: amount, Symbol: sym}
	asset.assert()
	return asset
}

func CheckParseInt64(s string) {
	MaxInt64String := strconv.FormatInt(math.MaxInt64, 10)
	if s[0] != '-' {
		EosAssert(len(s) < len(MaxInt64String) || (len(s) == len(MaxInt64String) && s <= MaxInt64String), &ParseErrorException{}, "Couldn't parse int64")
	} else {
		s = s[1:]
		EosAssert(len(s) < len(MaxInt64String) || (len(s) == len(MaxInt64String) && s <= MaxInt64String), &ParseErrorException{}, "Couldn't parse int64")
	}
}

type ExtendedAsset struct {
	Asset    Asset `json:"asset"`
	Contract AccountName
}

type SymbolCode = uint64

// NOTE: there's also a new ExtendedSymbol (which includes the contract (as AccountName) on which it is)
type Symbol struct {
	Precision uint8
	Symbol    string
}

func (s Symbol) Pack() (re []byte, err error) {
	symbol := make([]byte, 7, 7)
	copy(symbol[:], []byte(s.Symbol))

	re = append(re, byte(s.Precision))
	re = append(re, symbol...)
	return re, nil
}
func (s *Symbol) Unpack(in []byte) (int, error) {
	if len(in) < 8 {
		return 0, fmt.Errorf("asset symbol required [%d] bytes, remaining [%d]", 7, len(in))
	}
	s.Precision = uint8(in[0])
	s.Symbol = strings.TrimRight(string(in[1:8]), "\x00")
	return 8, nil

}

func StringToSymbol(precision uint8, str string) (result uint64) {
	Try(func() {
		len := uint32(len(str))
		for i := uint32(0); i < len; i++ {
			// All characters must be upper case alphabets
			EosAssert(str[i] >= 'A' && str[i] <= 'Z', &SymbolTypeException{}, "invalid character in symbol name")
			result |= uint64(str[i]) << (8 * (i + 1))
		}
		result |= uint64(precision)
	}).FcCaptureLogAndRethrow("str:%s", str)
	return
}

var MaxPrecision = uint8(18)

func (sym Symbol) FromString(from *string) Symbol {
	//TODO: unComplete
	EosAssert(!Empty(*from), &SymbolTypeException{}, "creating symbol from empty string")
	commaPos := strings.Index(*from, ",")
	EosAssert(commaPos != -1, &SymbolTypeException{}, "missing comma in symbol")
	precPart := string([]byte(*from)[:commaPos])
	p, _ := strconv.ParseInt(precPart, 10, 64)
	namePart := string([]byte(*from)[commaPos+1:])
	EosAssert(sym.ValidName(namePart), &SymbolTypeException{}, "invalid symbol: %s", namePart)
	EosAssert(uint8(p) <= MaxPrecision, &SymbolTypeException{}, "precision %v should be <= 18", p)
	return Symbol{Precision: uint8(p), Symbol: namePart}
}

func (sym Symbol) String() string {
	EosAssert(sym.Valid(), &SymbolTypeException{}, "symbol is not valid")
	v := sym.Precision
	ret := strconv.Itoa(int(v))
	ret += "," + sym.Symbol
	return ret
}

func (sym *Symbol) SymbolValue() uint64 {
	result := uint64(0)
	for i := len(sym.Symbol) - 1; i >= 0; i-- {
		if sym.Symbol[i] < 'A' || sym.Symbol[i] > 'Z' {
			log.Error("symbol cannot exceed A~Z")
		} else {
			result |= uint64(sym.Symbol[i])
		}
		result = result << 8
	}
	result |= uint64(sym.Precision)
	return result
}

func (sym *Symbol) ToSymbolCode() SymbolCode {
	return SymbolCode(sym.SymbolValue()) >> 8
}

func (sym *Symbol) Decimals() uint8 {
	return sym.Precision
}

func (sym *Symbol) Name() string {
	return sym.Symbol
}

func (sym *Symbol) Valid() bool {
	return sym.Decimals() <= MaxPrecision && len(sym.Symbol) != 0 && sym.ValidName(sym.Symbol)
}

func (sym *Symbol) ValidName(name string) bool {
	return -1 == strings.IndexFunc(name, func(r rune) bool {
		return !(r >= 'A' && r <= 'Z')
	})
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
	asset := Asset{Amount: amount, Symbol: EOSSymbol}
	asset.assert()
	return asset
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
	out.assert()
	return
}

func (a *Asset) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	var asset Asset
	invalid := false
	Try(func() {
		asset, err = NewAsset(s)
		if err != nil {
			invalid = true
		}
	}).Catch(func(e interface{}) {
		invalid = true
	}).End()

	if invalid {
		return err
	}

	*a = asset
	return nil
}

func (a Asset) MarshalJSON() (data []byte, err error) {
	return json.Marshal(a.String())
}
