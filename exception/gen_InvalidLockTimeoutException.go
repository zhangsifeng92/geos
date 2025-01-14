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

var InvalidLockTimeoutExceptionName = reflect.TypeOf(InvalidLockTimeoutException{}).Name()

type InvalidLockTimeoutException struct {
	_WalletException
	Elog log.Messages
}

func NewInvalidLockTimeoutException(parent _WalletException, message log.Message) *InvalidLockTimeoutException {
	return &InvalidLockTimeoutException{parent, log.Messages{message}}
}

func (e InvalidLockTimeoutException) Code() int64 {
	return 3120011
}

func (e InvalidLockTimeoutException) Name() string {
	return InvalidLockTimeoutExceptionName
}

func (e InvalidLockTimeoutException) What() string {
	return "Wallet lock timeout is invalid"
}

func (e *InvalidLockTimeoutException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e InvalidLockTimeoutException) GetLog() log.Messages {
	return e.Elog
}

func (e InvalidLockTimeoutException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e InvalidLockTimeoutException) DetailMessage() string {
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

func (e InvalidLockTimeoutException) String() string {
	return e.DetailMessage()
}

func (e InvalidLockTimeoutException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3120011,
		Name: InvalidLockTimeoutExceptionName,
		What: "Wallet lock timeout is invalid",
	}

	return json.Marshal(except)
}

func (e InvalidLockTimeoutException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*InvalidLockTimeoutException):
		callback(&e)
		return true
	case func(InvalidLockTimeoutException):
		callback(e)
		return true
	default:
		return false
	}
}
