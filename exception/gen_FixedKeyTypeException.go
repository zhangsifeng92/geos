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

var FixedKeyTypeExceptionName = reflect.TypeOf(FixedKeyTypeException{}).Name()

type FixedKeyTypeException struct {
	_ChainException
	Elog log.Messages
}

func NewFixedKeyTypeException(parent _ChainException, message log.Message) *FixedKeyTypeException {
	return &FixedKeyTypeException{parent, log.Messages{message}}
}

func (e FixedKeyTypeException) Code() int64 {
	return 3010013
}

func (e FixedKeyTypeException) Name() string {
	return FixedKeyTypeExceptionName
}

func (e FixedKeyTypeException) What() string {
	return "Invalid fixed key"
}

func (e *FixedKeyTypeException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e FixedKeyTypeException) GetLog() log.Messages {
	return e.Elog
}

func (e FixedKeyTypeException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e FixedKeyTypeException) DetailMessage() string {
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

func (e FixedKeyTypeException) String() string {
	return e.DetailMessage()
}

func (e FixedKeyTypeException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3010013,
		Name: FixedKeyTypeExceptionName,
		What: "Invalid fixed key",
	}

	return json.Marshal(except)
}

func (e FixedKeyTypeException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*FixedKeyTypeException):
		callback(&e)
		return true
	case func(FixedKeyTypeException):
		callback(e)
		return true
	default:
		return false
	}
}
