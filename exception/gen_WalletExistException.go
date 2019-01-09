// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/eosspark/eos-go/log"
)

// template type Exception(PARENT,CODE,WHAT)

var WalletExistExceptionName = reflect.TypeOf(WalletExistException{}).Name()

type WalletExistException struct {
	_WalletException
	Elog log.Messages
}

func NewWalletExistException(parent _WalletException, message log.Message) *WalletExistException {
	return &WalletExistException{parent, log.Messages{message}}
}

func (e WalletExistException) Code() int64 {
	return 3120001
}

func (e WalletExistException) Name() string {
	return WalletExistExceptionName
}

func (e WalletExistException) What() string {
	return "Wallet already exists"
}

func (e *WalletExistException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e WalletExistException) GetLog() log.Messages {
	return e.Elog
}

func (e WalletExistException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e WalletExistException) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(e.Name())
	buffer.WriteString(": ")
	buffer.WriteString(e.What())
	buffer.WriteString("\n")
	for _, l := range e.Elog {
		buffer.WriteString("[")
		buffer.WriteString(l.GetMessage())
		buffer.WriteString("] ")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e WalletExistException) String() string {
	return e.DetailMessage()
}

func (e WalletExistException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3120001,
		Name: WalletExistExceptionName,
		What: "Wallet already exists",
	}

	return json.Marshal(except)
}

func (e WalletExistException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*WalletExistException):
		callback(&e)
		return true
	case func(WalletExistException):
		callback(e)
		return true
	default:
		return false
	}
}
