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

var PackExceptionName = reflect.TypeOf(PackException{}).Name()

type PackException struct {
	_AbiException
	Elog log.Messages
}

func NewPackException(parent _AbiException, message log.Message) *PackException {
	return &PackException{parent, log.Messages{message}}
}

func (e PackException) Code() int64 {
	return 3150014
}

func (e PackException) Name() string {
	return PackExceptionName
}

func (e PackException) What() string {
	return "Pack data exception"
}

func (e *PackException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e PackException) GetLog() log.Messages {
	return e.Elog
}

func (e PackException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e PackException) DetailMessage() string {
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

func (e PackException) String() string {
	return e.DetailMessage()
}

func (e PackException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3150014,
		Name: PackExceptionName,
		What: "Pack data exception",
	}

	return json.Marshal(except)
}

func (e PackException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*PackException):
		callback(&e)
		return true
	case func(PackException):
		callback(e)
		return true
	default:
		return false
	}
}
