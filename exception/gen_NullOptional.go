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

var NullOptionalName = reflect.TypeOf(NullOptional{}).Name()

type NullOptional struct {
	Exception
	Elog log.Messages
}

func NewNullOptional(parent Exception, message log.Message) *NullOptional {
	return &NullOptional{parent, log.Messages{message}}
}

func (e NullOptional) Code() int64 {
	return NullOptionalCode
}

func (e NullOptional) Name() string {
	return NullOptionalName
}

func (e NullOptional) What() string {
	return "null optionale"
}

func (e *NullOptional) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e NullOptional) GetLog() log.Messages {
	return e.Elog
}

func (e NullOptional) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e NullOptional) DetailMessage() string {
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

func (e NullOptional) String() string {
	return e.DetailMessage()
}

func (e NullOptional) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: NullOptionalCode,
		Name: NullOptionalName,
		What: "null optionale",
	}

	return json.Marshal(except)
}

func (e NullOptional) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*NullOptional):
		callback(&e)
		return true
	case func(NullOptional):
		callback(e)
		return true
	default:
		return false
	}
}
