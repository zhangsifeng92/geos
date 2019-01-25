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

var InvalidOperationExceptionName = reflect.TypeOf(InvalidOperationException{}).Name()

type InvalidOperationException struct {
	Exception
	Elog log.Messages
}

func NewInvalidOperationException(parent Exception, message log.Message) *InvalidOperationException {
	return &InvalidOperationException{parent, log.Messages{message}}
}

func (e InvalidOperationException) Code() int64 {
	return InvalidOperationExceptionCode
}

func (e InvalidOperationException) Name() string {
	return InvalidOperationExceptionName
}

func (e InvalidOperationException) What() string {
	return "Invalid Operation"
}

func (e *InvalidOperationException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e InvalidOperationException) GetLog() log.Messages {
	return e.Elog
}

func (e InvalidOperationException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e InvalidOperationException) DetailMessage() string {
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

func (e InvalidOperationException) String() string {
	return e.DetailMessage()
}

func (e InvalidOperationException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: InvalidOperationExceptionCode,
		Name: InvalidOperationExceptionName,
		What: "Invalid Operation",
	}

	return json.Marshal(except)
}

func (e InvalidOperationException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*InvalidOperationException):
		callback(&e)
		return true
	case func(InvalidOperationException):
		callback(e)
		return true
	default:
		return false
	}
}
