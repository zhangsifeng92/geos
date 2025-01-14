// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/zhangsifeng92/geos/log"
)

// template type Exception(PARENT,CODE,WHAT)

var ContractBlacklistExceptionName = reflect.TypeOf(ContractBlacklistException{}).Name()

type ContractBlacklistException struct {
	_WhitelistBlacklistException
	Elog log.Messages
}

func NewContractBlacklistException(parent _WhitelistBlacklistException, message log.Message) *ContractBlacklistException {
	return &ContractBlacklistException{parent, log.Messages{message}}
}

func (e ContractBlacklistException) Code() int64 {
	return 3130004
}

func (e ContractBlacklistException) Name() string {
	return ContractBlacklistExceptionName
}

func (e ContractBlacklistException) What() string {
	return "Contract to execute is on the blacklist"
}

func (e *ContractBlacklistException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ContractBlacklistException) GetLog() log.Messages {
	return e.Elog
}

func (e ContractBlacklistException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e ContractBlacklistException) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteByte(' ')
	buffer.WriteString(e.Name())
	buffer.Write([]byte{':', ' '})
	buffer.WriteString(e.What())
	buffer.WriteByte('\n')
	for _, l := range e.Elog {
		buffer.WriteByte('[')
		buffer.WriteString(l.GetMessage())
		buffer.Write([]byte{']', ' '})
		buffer.WriteString(l.GetContext().String())
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func (e ContractBlacklistException) String() string {
	return e.DetailMessage()
}

func (e ContractBlacklistException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3130004,
		Name: ContractBlacklistExceptionName,
		What: "Contract to execute is on the blacklist",
	}

	return json.Marshal(except)
}

func (e ContractBlacklistException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ContractBlacklistException):
		callback(&e)
		return true
	case func(ContractBlacklistException):
		callback(e)
		return true
	default:
		return false
	}
}
