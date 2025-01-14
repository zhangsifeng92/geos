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

var SecureEnclaveExceptionName = reflect.TypeOf(SecureEnclaveException{}).Name()

type SecureEnclaveException struct {
	_WalletException
	Elog log.Messages
}

func NewSecureEnclaveException(parent _WalletException, message log.Message) *SecureEnclaveException {
	return &SecureEnclaveException{parent, log.Messages{message}}
}

func (e SecureEnclaveException) Code() int64 {
	return 3120012
}

func (e SecureEnclaveException) Name() string {
	return SecureEnclaveExceptionName
}

func (e SecureEnclaveException) What() string {
	return "Secure Enclave Exception"
}

func (e *SecureEnclaveException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e SecureEnclaveException) GetLog() log.Messages {
	return e.Elog
}

func (e SecureEnclaveException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e SecureEnclaveException) DetailMessage() string {
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

func (e SecureEnclaveException) String() string {
	return e.DetailMessage()
}

func (e SecureEnclaveException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3120012,
		Name: SecureEnclaveExceptionName,
		What: "Secure Enclave Exception",
	}

	return json.Marshal(except)
}

func (e SecureEnclaveException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*SecureEnclaveException):
		callback(&e)
		return true
	case func(SecureEnclaveException):
		callback(e)
		return true
	default:
		return false
	}
}
