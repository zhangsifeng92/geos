package types

import (
	// "bytes"
	// "crypto/sha256"
	// "encoding/json"
	// "fmt"
	// "github.com/eos-go/rlp"
	"github.com/eosspark/eos-go/common"
	// "reflect"
)

// See: libraries/chain/include/eosio/chain/contracts/types.hpp:203
// See: build/contracts/eosio.system/eosio.system.abi

// SetCode represents the hard-coded `setcode` action.
type SetCode struct {
	Account   common.AccountName `json:"account"`
	VMType    byte               `json:"vmtype"`
	VMVersion byte               `json:"vmversion"`
	Code      common.HexBytes    `json:"bytes"`
}

// SetABI represents the hard-coded `setabi` action.
type SetABI struct {
	Account common.AccountName `json:"account"`
	ABI     ABI                `json:"abi"`
}

// Action
type Action struct {
	Account       common.AccountName       `json:"account"`
	Name          common.ActionName        `json:"name"`
	Authorization []common.PermissionLevel `json:"authorization,omitempty"`
	ActionData
}

// func (a Action) Digest() SHA256Bytes {
// 	toEat := jsonActionToServer{
// 		Account:       a.Account,
// 		Name:          a.Name,
// 		Authorization: a.Authorization,
// 		Data:          a.ActionData.HexData,
// 	}
// 	bin, err := rlp.MarshalBinary(toEat)
// 	if err != nil {
// 		panic("this should never panic, we know it marshals properly all the time")
// 	}

// 	h := sha256.New()
// 	_, _ = h.Write(bin)
// 	return h.Sum(nil)
// }

type ActionData struct {
	HexData  common.HexBytes `json:"hex_data,omitempty"`
	Data     interface{}     `json:"data,omitempty" eos:"-"`
	abi      []byte          // TBD: we could use the ABI to decode in obj
	toServer bool
}

func NewActionData(obj interface{}) ActionData {
	return ActionData{
		HexData:  []byte{},
		Data:     obj,
		toServer: true,
	}
}
func (a *ActionData) SetToServer(toServer bool) {
	// FIXME: let's clarify this design, make it simpler..
	// toServer doesn't speak of the intent..
	a.toServer = toServer
}

//  jsonActionToServer represents what /v1/chain/push_transaction
//  expects, which isn't allllways the same everywhere.
type jsonActionToServer struct {
	Account       common.AccountName       `json:"account"`
	Name          common.ActionName        `json:"name"`
	Authorization []common.PermissionLevel `json:"authorization,omitempty"`
	Data          common.HexBytes          `json:"data,omitempty"`
}

type jsonActionFromServer struct {
	Account       common.AccountName       `json:"account"`
	Name          common.ActionName        `json:"name"`
	Authorization []common.PermissionLevel `json:"authorization,omitempty"`
	Data          interface{}              `json:"data,omitempty"`
	HexData       common.HexBytes          `json:"hex_data,omitempty"`
}

// func (a *Action) MarshalJSON() ([]byte, error) {
// 	if a.toServer {
// 		buf := new(bytes.Buffer)
// 		if err := rlp.Encode(buf, a.ActionData.Data); err != nil {
// 			return nil, err
// 		}

// 		data := buf.Bytes()
// 		println("MarshalJSON data length : ", len(data)) /**/

// 		return json.Marshal(&jsonActionToServer{
// 			Account:       a.Account,
// 			Name:          a.Name,
// 			Authorization: a.Authorization,
// 			Data:          data,
// 		})
// 	}

// 	return json.Marshal(&jsonActionFromServer{
// 		Account:       a.Account,
// 		Name:          a.Name,
// 		Authorization: a.Authorization,
// 		HexData:       a.HexData,
// 		Data:          a.Data,
// 	})
// }

// func (a *Action) MapToRegisteredAction() error {
// 	src, ok := a.ActionData.Data.(map[string]interface{})
// 	if !ok {
// 		return nil
// 	}

// 	actionMap := RegisteredActions[a.Account]

// 	var decodeInto reflect.Type
// 	if actionMap != nil {
// 		objType := actionMap[a.Name]
// 		if objType != nil {
// 			decodeInto = objType
// 		}
// 	}
// 	if decodeInto == nil {
// 		return nil
// 	}

// 	obj := reflect.New(decodeInto)
// 	objIface := obj.Interface()

// 	cnt, err := json.Marshal(src)
// 	if err != nil {
// 		return fmt.Errorf("marshaling data: %s", err)
// 	}
// 	err = json.Unmarshal(cnt, objIface)
// 	if err != nil {
// 		return fmt.Errorf("json unmarshal into registered actions: %s", err)
// 	}

// 	a.ActionData.Data = objIface

// 	return nil
// }