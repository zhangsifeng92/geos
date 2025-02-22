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

var WalletExceptionName = reflect.TypeOf(WalletException{}).Name()

type WalletException struct {
	_WalletException
	Elog log.Messages
}

func NewWalletException(parent _WalletException, message log.Message) *WalletException {
	return &WalletException{parent, log.Messages{message}}
}

func (e WalletException) Code() int64 {
	return 3120000
}

func (e WalletException) Name() string {
	return WalletExceptionName
}

func (e WalletException) What() string {
	return "Invalid contract vm version"
}

func (e *WalletException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e WalletException) GetLog() log.Messages {
	return e.Elog
}

func (e WalletException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e WalletException) DetailMessage() string {
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

func (e WalletException) String() string {
	return e.DetailMessage()
}

func (e WalletException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3120000,
		Name: WalletExceptionName,
		What: "Invalid contract vm version",
	}

	return json.Marshal(except)
}

func (e WalletException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*WalletException):
		callback(&e)
		return true
	case func(WalletException):
		callback(e)
		return true
	default:
		return false
	}
}
