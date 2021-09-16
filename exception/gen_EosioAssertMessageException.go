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

var EosioAssertMessageExceptionName = reflect.TypeOf(EosioAssertMessageException{}).Name()

type EosioAssertMessageException struct {
	_ActionValidateException
	Elog log.Messages
}

func NewEosioAssertMessageException(parent _ActionValidateException, message log.Message) *EosioAssertMessageException {
	return &EosioAssertMessageException{parent, log.Messages{message}}
}

func (e EosioAssertMessageException) Code() int64 {
	return 3050003
}

func (e EosioAssertMessageException) Name() string {
	return EosioAssertMessageExceptionName
}

func (e EosioAssertMessageException) What() string {
	return "eosio_assert_message assertion failure"
}

func (e *EosioAssertMessageException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e EosioAssertMessageException) GetLog() log.Messages {
	return e.Elog
}

func (e EosioAssertMessageException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e EosioAssertMessageException) DetailMessage() string {
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

func (e EosioAssertMessageException) String() string {
	return e.DetailMessage()
}

func (e EosioAssertMessageException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3050003,
		Name: EosioAssertMessageExceptionName,
		What: "eosio_assert_message assertion failure",
	}

	return json.Marshal(except)
}

func (e EosioAssertMessageException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*EosioAssertMessageException):
		callback(&e)
		return true
	case func(EosioAssertMessageException):
		callback(e)
		return true
	default:
		return false
	}
}
