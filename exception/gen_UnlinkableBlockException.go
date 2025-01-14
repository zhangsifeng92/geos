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

var UnlinkableBlockExceptionName = reflect.TypeOf(UnlinkableBlockException{}).Name()

type UnlinkableBlockException struct {
	_BlockValidateException
	Elog log.Messages
}

func NewUnlinkableBlockException(parent _BlockValidateException, message log.Message) *UnlinkableBlockException {
	return &UnlinkableBlockException{parent, log.Messages{message}}
}

func (e UnlinkableBlockException) Code() int64 {
	return 3030001
}

func (e UnlinkableBlockException) Name() string {
	return UnlinkableBlockExceptionName
}

func (e UnlinkableBlockException) What() string {
	return "Unlinkable block"
}

func (e *UnlinkableBlockException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e UnlinkableBlockException) GetLog() log.Messages {
	return e.Elog
}

func (e UnlinkableBlockException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e UnlinkableBlockException) DetailMessage() string {
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

func (e UnlinkableBlockException) String() string {
	return e.DetailMessage()
}

func (e UnlinkableBlockException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3030001,
		Name: UnlinkableBlockExceptionName,
		What: "Unlinkable block",
	}

	return json.Marshal(except)
}

func (e UnlinkableBlockException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*UnlinkableBlockException):
		callback(&e)
		return true
	case func(UnlinkableBlockException):
		callback(e)
		return true
	default:
		return false
	}
}
