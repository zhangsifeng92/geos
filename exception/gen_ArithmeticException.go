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

var ArithmeticExceptionName = reflect.TypeOf(ArithmeticException{}).Name()

type ArithmeticException struct {
	_ContractApiException
	Elog log.Messages
}

func NewArithmeticException(parent _ContractApiException, message log.Message) *ArithmeticException {
	return &ArithmeticException{parent, log.Messages{message}}
}

func (e ArithmeticException) Code() int64 {
	return 3230003
}

func (e ArithmeticException) Name() string {
	return ArithmeticExceptionName
}

func (e ArithmeticException) What() string {
	return "Arithmetic exception"
}

func (e *ArithmeticException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ArithmeticException) GetLog() log.Messages {
	return e.Elog
}

func (e ArithmeticException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e ArithmeticException) DetailMessage() string {
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

func (e ArithmeticException) String() string {
	return e.DetailMessage()
}

func (e ArithmeticException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3230003,
		Name: ArithmeticExceptionName,
		What: "Arithmetic exception",
	}

	return json.Marshal(except)
}

func (e ArithmeticException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ArithmeticException):
		callback(&e)
		return true
	case func(ArithmeticException):
		callback(e)
		return true
	default:
		return false
	}
}
