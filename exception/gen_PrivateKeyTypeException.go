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

var PrivateKeyTypeExceptionName = reflect.TypeOf(PrivateKeyTypeException{}).Name()

type PrivateKeyTypeException struct {
	_ChainException
	Elog log.Messages
}

func NewPrivateKeyTypeException(parent _ChainException, message log.Message) *PrivateKeyTypeException {
	return &PrivateKeyTypeException{parent, log.Messages{message}}
}

func (e PrivateKeyTypeException) Code() int64 {
	return 3010003
}

func (e PrivateKeyTypeException) Name() string {
	return PrivateKeyTypeExceptionName
}

func (e PrivateKeyTypeException) What() string {
	return "Invalid private key"
}

func (e *PrivateKeyTypeException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e PrivateKeyTypeException) GetLog() log.Messages {
	return e.Elog
}

func (e PrivateKeyTypeException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e PrivateKeyTypeException) DetailMessage() string {
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

func (e PrivateKeyTypeException) String() string {
	return e.DetailMessage()
}

func (e PrivateKeyTypeException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3010003,
		Name: PrivateKeyTypeExceptionName,
		What: "Invalid private key",
	}

	return json.Marshal(except)
}

func (e PrivateKeyTypeException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*PrivateKeyTypeException):
		callback(&e)
		return true
	case func(PrivateKeyTypeException):
		callback(e)
		return true
	default:
		return false
	}
}
