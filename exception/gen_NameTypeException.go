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

var NameTypeExceptionName = reflect.TypeOf(NameTypeException{}).Name()

type NameTypeException struct {
	_ChainException
	Elog log.Messages
}

func NewNameTypeException(parent _ChainException, message log.Message) *NameTypeException {
	return &NameTypeException{parent, log.Messages{message}}
}

func (e NameTypeException) Code() int64 {
	return 3010001
}

func (e NameTypeException) Name() string {
	return NameTypeExceptionName
}

func (e NameTypeException) What() string {
	return "Invalid name"
}

func (e *NameTypeException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e NameTypeException) GetLog() log.Messages {
	return e.Elog
}

func (e NameTypeException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e NameTypeException) DetailMessage() string {
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

func (e NameTypeException) String() string {
	return e.DetailMessage()
}

func (e NameTypeException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3010001,
		Name: NameTypeExceptionName,
		What: "Invalid name",
	}

	return json.Marshal(except)
}

func (e NameTypeException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*NameTypeException):
		callback(&e)
		return true
	case func(NameTypeException):
		callback(e)
		return true
	default:
		return false
	}
}
