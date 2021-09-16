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

var PublicKeyTypeExceptionName = reflect.TypeOf(PublicKeyTypeException{}).Name()

type PublicKeyTypeException struct {
	_ChainException
	Elog log.Messages
}

func NewPublicKeyTypeException(parent _ChainException, message log.Message) *PublicKeyTypeException {
	return &PublicKeyTypeException{parent, log.Messages{message}}
}

func (e PublicKeyTypeException) Code() int64 {
	return 3010002
}

func (e PublicKeyTypeException) Name() string {
	return PublicKeyTypeExceptionName
}

func (e PublicKeyTypeException) What() string {
	return "Invalid public key"
}

func (e *PublicKeyTypeException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e PublicKeyTypeException) GetLog() log.Messages {
	return e.Elog
}

func (e PublicKeyTypeException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e PublicKeyTypeException) DetailMessage() string {
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

func (e PublicKeyTypeException) String() string {
	return e.DetailMessage()
}

func (e PublicKeyTypeException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3010002,
		Name: PublicKeyTypeExceptionName,
		What: "Invalid public key",
	}

	return json.Marshal(except)
}

func (e PublicKeyTypeException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*PublicKeyTypeException):
		callback(&e)
		return true
	case func(PublicKeyTypeException):
		callback(e)
		return true
	default:
		return false
	}
}
