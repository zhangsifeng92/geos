package abi_serializer

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhangsifeng92/geos/chain/types"
	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/log"
	"io"
	"strconv"
	"strings"
)

var abiLog log.Logger

type typeName = string
type fieldName = string

func init() {
	abiLog = log.New("abi")
	//abiLog.SetHandler(log.TerminalHandler)
	abiLog.SetHandler(log.DiscardHandler())
}

type TypeDef struct {
	NewTypeName typeName `json:"new_type_name"`
	Type        typeName `json:"type"`
}

type FieldDef struct {
	Name fieldName `json:"name"`
	Type typeName  `json:"type"`
}

type StructDef struct {
	Name   typeName   `json:"name"`
	Base   typeName   `json:"base"`
	Fields []FieldDef `json:"fields,omitempty"`
}

type ActionDef struct {
	Name              common.ActionName `json:"name"`
	Type              typeName          `json:"type"`
	RicardianContract string            `json:"ricardian_contract"`
}

// TableDef defines a table. See libraries/chain/include/eosio/chain/contracts/types.hpp:78
type TableDef struct {
	Name      common.TableName `json:"name"`
	IndexType typeName         `json:"index_type"`
	KeyNames  []fieldName      `json:"key_names,omitempty"`
	KeyTypes  []typeName       `json:"key_types,omitempty"`
	Type      typeName         `json:"type"`
}

// ClausePair represents clauses, related to Ricardian Contracts.
type ClausePair struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type ErrorMessage struct {
	Code    uint64 `json:"error_code"`
	Message string `json:"error_msg"`
}

type VariantDef struct {
	Name  typeName   `json:"name"`
	Types []typeName `json:"types"`
}

func CommonTypeDefs() []TypeDef {
	types := make([]TypeDef, 7)
	types[0] = TypeDef{"account_name", "name"}
	types[1] = TypeDef{"permission_name", "name"}
	types[2] = TypeDef{"action_name", "name"}
	types[3] = TypeDef{"table_name", "name"}
	types[4] = TypeDef{"transaction_id_type", "checksum256"}
	types[5] = TypeDef{"block_id_type", "checksum256"}
	types[6] = TypeDef{"weight_type", "uint16"}
	return types
}

type AbiDef struct {
	Version          string             `json:"version"`
	Types            []TypeDef          `json:"types,omitempty"`
	Structs          []StructDef        `json:"structs,omitempty"`
	Actions          []ActionDef        `json:"actions,omitempty"`
	Tables           []TableDef         `json:"tables,omitempty"`
	RicardianClauses []ClausePair       `json:"ricardian_clauses,omitempty"`
	ErrorMessages    []ErrorMessage     `json:"error_messages,omitempty"`
	Extensions       []*types.Extension `json:"abi_extensions,omitempty"`
	Variants         []VariantDef       `json:"variants,omitempty"` // TODO may not exit
}

func NewABI(r io.Reader) (*AbiDef, error) {
	abi := &AbiDef{}
	err := json.NewDecoder(r).Decode(abi)
	if err != nil {
		return nil, fmt.Errorf("read abi: %s", err)
	}
	return abi, nil
}

func (a *AbiDef) ActionForName(name common.ActionName) *ActionDef {
	for _, a := range a.Actions {
		if a.Name == name {
			return &a
		}
	}
	return nil
}

func (a *AbiDef) StructForName(name typeName) *StructDef {
	for _, s := range a.Structs {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

func (a *AbiDef) TableForName(name common.TableName) *TableDef {
	for _, s := range a.Tables {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

func (a *AbiDef) TypeNameForNewTypeName(typeName string) string {
	for _, t := range a.Types {
		if t.NewTypeName == typeName {
			return t.Type
		}
	}
	return typeName
}

func (a AbiDef) IsEmpty() bool {
	return a.Version != "" && len(a.Types) == 0 && len(a.Structs) == 0 && len(a.Actions) == 0 &&
		len(a.Tables) == 0 && len(a.RicardianClauses) == 0 && len(a.ErrorMessages) == 0 &&
		len(a.Extensions) == 0 && len(a.Variants) == 0
}

//types for abi

type Int64 int64

func (i Int64) MarshalJSON() (data []byte, err error) {
	if i > 0xffffffff || i < -0xffffffff {
		encodedInt, err := json.Marshal(int64(i))
		if err != nil {
			return nil, err
		}
		data = append([]byte{'"'}, encodedInt...)
		data = append(data, '"')
		return data, nil
	}
	return json.Marshal(int64(i))
}

func (i *Int64) UnmarshalJSON(data []byte) error {
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

		*i = Int64(val)

		return nil
	}

	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Int64(v)

	return nil
}

type Uint64 uint64

func (i Uint64) MarshalJSON() (data []byte, err error) {
	if i > 0xffffffff {
		encodedInt, err := json.Marshal(uint64(i))
		if err != nil {
			return nil, err
		}
		data = append([]byte{'"'}, encodedInt...)
		data = append(data, '"')
		return data, nil
	}
	return json.Marshal(uint64(i))
}

func (i *Uint64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}

		*i = Uint64(val)

		return nil
	}

	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = Uint64(v)

	return nil
}

type Uint128 struct {
	Lo uint64
	Hi uint64
}

type Int128 Uint128

type Float128 Uint128

func (i Uint128) MarshalJSON() (data []byte, err error) {
	return json.Marshal(i.String())
}

func (i Int128) MarshalJSON() (data []byte, err error) {
	return json.Marshal(Uint128(i).String())
}

func (i Float128) MarshalJSON() (data []byte, err error) {
	return json.Marshal(Uint128(i).String())
}

func (i Uint128) String() string {
	// Same for Int128, Float128
	number := make([]byte, 16)
	binary.LittleEndian.PutUint64(number[:], i.Lo)
	binary.LittleEndian.PutUint64(number[8:], i.Hi)
	return fmt.Sprintf("0x%s%s", hex.EncodeToString(number[:8]), hex.EncodeToString(number[8:]))
}

func (i *Int128) UnmarshalJSON(data []byte) error {
	var el Uint128
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}

	out := Int128(el)
	*i = out

	return nil
}

func (i *Float128) UnmarshalJSON(data []byte) error {
	var el Uint128
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}

	out := Float128(el)
	*i = out

	return nil
}

func (i *Uint128) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if !strings.HasPrefix(s, "0x") && !strings.HasPrefix(s, "0X") {
		return fmt.Errorf("int128 expects 0x prefix")
	}

	truncatedVal := s[2:]
	if len(truncatedVal) != 32 {
		return fmt.Errorf("int128 expects 32 characters after 0x, had %d", len(truncatedVal))
	}

	loHex := truncatedVal[:16]
	hiHex := truncatedVal[16:]

	lo, err := hex.DecodeString(loHex)
	if err != nil {
		return err
	}

	hi, err := hex.DecodeString(hiHex)
	if err != nil {
		return err
	}

	loUint := binary.LittleEndian.Uint64(lo)
	hiUint := binary.LittleEndian.Uint64(hi)

	i.Lo = loUint
	i.Hi = hiUint

	return nil
}
