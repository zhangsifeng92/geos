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

var WrongSigningKeyName = reflect.TypeOf(WrongSigningKey{}).Name()

type WrongSigningKey struct {
	_BlockValidateException
	Elog log.Messages
}

func NewWrongSigningKey(parent _BlockValidateException, message log.Message) *WrongSigningKey {
	return &WrongSigningKey{parent, log.Messages{message}}
}

func (e WrongSigningKey) Code() int64 {
	return 3030008
}

func (e WrongSigningKey) Name() string {
	return WrongSigningKeyName
}

func (e WrongSigningKey) What() string {
	return "Block is not signed with expected key"
}

func (e *WrongSigningKey) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e WrongSigningKey) GetLog() log.Messages {
	return e.Elog
}

func (e WrongSigningKey) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); len(msg) > 0 {
			return msg
		}
	}
	return e.String()
}

func (e WrongSigningKey) DetailMessage() string {
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

func (e WrongSigningKey) String() string {
	return e.DetailMessage()
}

func (e WrongSigningKey) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3030008,
		Name: WrongSigningKeyName,
		What: "Block is not signed with expected key",
	}

	return json.Marshal(except)
}

func (e WrongSigningKey) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*WrongSigningKey):
		callback(&e)
		return true
	case func(WrongSigningKey):
		callback(e)
		return true
	default:
		return false
	}
}
